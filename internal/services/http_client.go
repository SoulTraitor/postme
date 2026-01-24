package services

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/SoulTraitor/postme/internal/models"
	"github.com/andybalholm/brotli"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
	"golang.org/x/sys/windows/registry"
)

// Default User-Agent that mimics a modern browser
const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

// HTTPClient handles HTTP request execution
type HTTPClient struct {
	client         *http.Client
	useSystemProxy bool
	mu             sync.Mutex
}

// utlsTransport wraps http.Transport to use uTLS for TLS fingerprint spoofing
type utlsTransport struct {
	proxyFunc func(*http.Request) (*url.URL, error)
	dialer    *net.Dialer
}

func (t *utlsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// For HTTPS requests, use uTLS
	if req.URL.Scheme == "https" {
		return t.roundTripHTTPS(req)
	}
	
	// For HTTP requests, use standard transport
	transport := &http.Transport{
		Proxy:             t.proxyFunc,
		DialContext:       t.dialer.DialContext,
		ForceAttemptHTTP2: true,
	}
	return transport.RoundTrip(req)
}

func (t *utlsTransport) roundTripHTTPS(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	
	// Determine target address
	host := req.URL.Host
	if !strings.Contains(host, ":") {
		host = host + ":443"
	}
	
	// Check for proxy
	var conn net.Conn
	var err error
	
	if t.proxyFunc != nil {
		proxyURL, proxyErr := t.proxyFunc(req)
		if proxyErr == nil && proxyURL != nil {
			// Connect through proxy
			conn, err = t.dialThroughProxy(ctx, proxyURL, host)
		} else {
			// Direct connection
			conn, err = t.dialer.DialContext(ctx, "tcp", host)
		}
	} else {
		conn, err = t.dialer.DialContext(ctx, "tcp", host)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Create uTLS connection with Chrome fingerprint
	tlsConn := utls.UClient(conn, &utls.Config{
		ServerName: req.URL.Hostname(),
	}, utls.HelloChrome_120)
	
	// Perform TLS handshake
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return nil, err
	}
	
	// Check if ALPN negotiated HTTP/2
	alpn := tlsConn.ConnectionState().NegotiatedProtocol
	
	if alpn == "h2" {
		// Use HTTP/2 transport
		h2Transport := &http2.Transport{}
		h2Conn, err := h2Transport.NewClientConn(tlsConn)
		if err != nil {
			tlsConn.Close()
			return nil, err
		}
		return h2Conn.RoundTrip(req)
	}
	
	// Use HTTP/1.1
	transport := &http.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return tlsConn, nil
		},
	}
	
	return transport.RoundTrip(req)
}

func (t *utlsTransport) dialThroughProxy(ctx context.Context, proxyURL *url.URL, targetHost string) (net.Conn, error) {
	// Connect to proxy
	proxyHost := proxyURL.Host
	if !strings.Contains(proxyHost, ":") {
		if proxyURL.Scheme == "https" {
			proxyHost = proxyHost + ":443"
		} else {
			proxyHost = proxyHost + ":80"
		}
	}
	
	conn, err := t.dialer.DialContext(ctx, "tcp", proxyHost)
	if err != nil {
		return nil, err
	}
	
	// Send CONNECT request
	connectReq := &http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Opaque: targetHost},
		Host:   targetHost,
		Header: make(http.Header),
	}
	
	// Add proxy authentication if present
	if proxyURL.User != nil {
		password, _ := proxyURL.User.Password()
		connectReq.SetBasicAuth(proxyURL.User.Username(), password)
	}
	
	if err := connectReq.Write(conn); err != nil {
		conn.Close()
		return nil, err
	}
	
	// Read response
	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, connectReq)
	if err != nil {
		conn.Close()
		return nil, err
	}
	
	if resp.StatusCode != 200 {
		conn.Close()
		return nil, &net.OpError{Op: "dial", Err: &proxyError{resp.Status}}
	}
	
	return conn, nil
}

type proxyError struct {
	status string
}

func (e *proxyError) Error() string {
	return "proxy connect failed: " + e.status
}

// NewHTTPClient creates a new HTTPClient with browser-like TLS fingerprint
func NewHTTPClient() *HTTPClient {
	c := &HTTPClient{
		useSystemProxy: true,
	}
	
	c.rebuildClient()
	return c
}

func (c *HTTPClient) rebuildClient() {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	
	var proxyFunc func(*http.Request) (*url.URL, error)
	if c.useSystemProxy {
		proxyFunc = c.getProxyFunc()
	}
	
	transport := &utlsTransport{
		proxyFunc: proxyFunc,
		dialer:    dialer,
	}
	
	c.client = &http.Client{
		Transport: transport,
		Timeout:   0, // We handle timeout via context
	}
}

// SetUseSystemProxy enables or disables system proxy usage
func (c *HTTPClient) SetUseSystemProxy(useProxy bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.useSystemProxy = useProxy
	c.rebuildClient()
}

// getProxyFunc returns a proxy function based on system settings
func (c *HTTPClient) getProxyFunc() func(*http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		// Try to get Windows system proxy
		if runtime.GOOS == "windows" {
			proxyURL := getWindowsSystemProxy()
			if proxyURL != "" {
				return url.Parse(proxyURL)
			}
		}
		// Fall back to environment variables
		return http.ProxyFromEnvironment(req)
	}
}

// getWindowsSystemProxy reads the Windows system proxy settings from registry
func getWindowsSystemProxy() string {
	if runtime.GOOS != "windows" {
		return ""
	}

	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer key.Close()

	// Check if proxy is enabled
	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil || proxyEnable == 0 {
		return ""
	}

	// Get proxy server
	proxyServer, _, err := key.GetStringValue("ProxyServer")
	if err != nil || proxyServer == "" {
		return ""
	}

	// If it doesn't have a scheme, add http://
	if !strings.HasPrefix(proxyServer, "http://") && !strings.HasPrefix(proxyServer, "https://") {
		proxyServer = "http://" + proxyServer
	}

	return proxyServer
}

// ExecuteRequest represents the request to execute
type ExecuteRequest struct {
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Headers  []models.KeyValue `json:"headers"`
	Body     string            `json:"body"`
	BodyType string            `json:"bodyType"`
	Timeout  float64           `json:"timeout"`
}

// Execute executes an HTTP request
func (c *HTTPClient) Execute(ctx context.Context, req ExecuteRequest) (*models.Response, error) {
	// Create context with timeout if specified
	if req.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(req.Timeout*float64(time.Second)))
		defer cancel()
	}

	// Create body reader and content type
	var bodyReader io.Reader
	var contentType string

	if req.BodyType == "form-data" && req.Body != "" {
		// Parse form data from JSON array
		var formItems []models.KeyValue
		if err := json.Unmarshal([]byte(req.Body), &formItems); err == nil {
			// Build multipart form
			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			for _, item := range formItems {
				if item.Enabled && item.Key != "" {
					if item.Type == "file" && item.Value != "" {
						// File upload
						file, err := os.Open(item.Value)
						if err != nil {
							return nil, fmt.Errorf("failed to open form file %s: %w", item.Value, err)
						}
						// Get just the filename for the form
						fileName := item.Value
						if idx := strings.LastIndexAny(fileName, "/\\"); idx >= 0 {
							fileName = fileName[idx+1:]
						}
						part, err := writer.CreateFormFile(item.Key, fileName)
						if err != nil {
							file.Close()
							return nil, fmt.Errorf("failed to create form file: %w", err)
						}
						_, err = io.Copy(part, file)
						file.Close()
						if err != nil {
							return nil, fmt.Errorf("failed to write form file: %w", err)
						}
					} else {
						// Text field
						writer.WriteField(item.Key, item.Value)
					}
				}
			}
			writer.Close()
			bodyReader = &buf
			contentType = writer.FormDataContentType()
		}
	} else if req.BodyType == "x-www-form-urlencoded" && req.Body != "" {
		// Parse form data from JSON array and encode as URL encoded form
		var formItems []models.KeyValue
		if err := json.Unmarshal([]byte(req.Body), &formItems); err == nil {
			formData := url.Values{}
			for _, item := range formItems {
				if item.Enabled && item.Key != "" {
					formData.Add(item.Key, item.Value)
				}
			}
			bodyReader = strings.NewReader(formData.Encode())
			contentType = "application/x-www-form-urlencoded"
		}
	} else if req.BodyType == "binary" && req.Body != "" {
		// Body contains file path - read the file
		file, err := os.Open(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to open binary file: %w", err)
		}
		defer file.Close()
		
		// Read file into buffer
		fileContent, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read binary file: %w", err)
		}
		bodyReader = bytes.NewReader(fileContent)
		contentType = "application/octet-stream"
	} else if req.Body != "" && req.BodyType != "none" {
		bodyReader = strings.NewReader(req.Body)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	for _, h := range req.Headers {
		if h.Enabled && h.Key != "" {
			httpReq.Header.Set(h.Key, h.Value)
		}
	}

	// Set content-type based on body type (if not already set by user)
	if httpReq.Header.Get("Content-Type") == "" {
		if contentType != "" {
			httpReq.Header.Set("Content-Type", contentType)
		} else if req.BodyType == "json" {
			httpReq.Header.Set("Content-Type", "application/json")
		} else if req.BodyType == "xml" {
			httpReq.Header.Set("Content-Type", "application/xml")
		} else if req.BodyType == "text" {
			httpReq.Header.Set("Content-Type", "text/plain")
		}
	}

	// Set default User-Agent if not already set (helps bypass Cloudflare detection)
	if httpReq.Header.Get("User-Agent") == "" {
		httpReq.Header.Set("User-Agent", defaultUserAgent)
	}

	// Set Accept header if not set (browser-like behavior)
	if httpReq.Header.Get("Accept") == "" {
		httpReq.Header.Set("Accept", "*/*")
	}

	// Set Accept-Language if not set
	if httpReq.Header.Get("Accept-Language") == "" {
		httpReq.Header.Set("Accept-Language", "en-US,en;q=0.9")
	}
	
	// Set Accept-Encoding for compression (browser-like)
	if httpReq.Header.Get("Accept-Encoding") == "" {
		httpReq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	}
	
	// Set Sec-Fetch headers (modern browser behavior)
	if httpReq.Header.Get("Sec-Fetch-Dest") == "" {
		httpReq.Header.Set("Sec-Fetch-Dest", "empty")
	}
	if httpReq.Header.Get("Sec-Fetch-Mode") == "" {
		httpReq.Header.Set("Sec-Fetch-Mode", "cors")
	}
	if httpReq.Header.Get("Sec-Fetch-Site") == "" {
		httpReq.Header.Set("Sec-Fetch-Site", "cross-site")
	}

	// Execute request
	startTime := time.Now()
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	duration := time.Since(startTime).Milliseconds()

	// Read response body with decompression
	var respBodyReader io.Reader = resp.Body
	contentEncoding := resp.Header.Get("Content-Encoding")
	switch strings.ToLower(contentEncoding) {
	case "gzip":
		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress gzip: %w", err)
		}
		defer gr.Close()
		respBodyReader = gr
	case "deflate":
		respBodyReader = flate.NewReader(resp.Body)
	case "br":
		respBodyReader = brotli.NewReader(resp.Body)
	}

	body, err := io.ReadAll(respBodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Extract headers
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &models.Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    headers,
		Body:       string(body),
		Size:       int64(len(body)),
		Duration:   duration,
	}, nil
}

// BuildRequestHeadersJSON builds JSON string from headers
func BuildRequestHeadersJSON(headers []models.KeyValue) string {
	data, _ := json.Marshal(headers)
	return string(data)
}

// BuildResponseHeadersJSON builds JSON string from response headers
func BuildResponseHeadersJSON(headers map[string]string) string {
	data, _ := json.Marshal(headers)
	return string(data)
}
