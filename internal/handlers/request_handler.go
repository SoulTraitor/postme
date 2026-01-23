package handlers

import (
	"context"
	"sync"

	"github.com/SoulTraitor/postme/internal/database"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/SoulTraitor/postme/internal/services"
)

// RequestHandler handles request-related operations for the frontend
type RequestHandler struct {
	service    *services.RequestService
	httpClient *services.HTTPClient
	history    *services.HistoryService

	// For request cancellation
	mu           sync.Mutex
	cancelFuncs  map[string]context.CancelFunc
}

// NewRequestHandler creates a new RequestHandler
func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

// Init initializes the handler with database connection
func (h *RequestHandler) Init() {
	db := database.GetDB()
	h.service = services.NewRequestService(db)
	h.httpClient = services.NewHTTPClient()
	h.history = services.NewHistoryService(db)
}

// Create creates a new request
func (h *RequestHandler) Create(req models.Request) (*models.Request, error) {
	if err := h.service.Create(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// GetByID retrieves a request by ID
func (h *RequestHandler) GetByID(id int64) (*models.Request, error) {
	return h.service.GetByID(id)
}

// Update updates a request
func (h *RequestHandler) Update(req models.Request) error {
	return h.service.Update(&req)
}

// Delete deletes a request
func (h *RequestHandler) Delete(id int64) error {
	return h.service.Delete(id)
}

// ExecuteRequestParams represents the parameters for executing a request
type ExecuteRequestParams struct {
	TabID    string            `json:"tabId"`
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Headers  []models.KeyValue `json:"headers"`
	Body     string            `json:"body"`
	BodyType string            `json:"bodyType"`
	Timeout  float64           `json:"timeout"`
}

// Execute executes an HTTP request
func (h *RequestHandler) Execute(params ExecuteRequestParams) (*models.Response, error) {
	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	h.mu.Lock()
	h.cancelFuncs[params.TabID] = cancel
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.cancelFuncs, params.TabID)
		h.mu.Unlock()
	}()

	// Execute request
	resp, err := h.httpClient.Execute(ctx, services.ExecuteRequest{
		Method:   params.Method,
		URL:      params.URL,
		Headers:  params.Headers,
		Body:     params.Body,
		BodyType: params.BodyType,
		Timeout:  params.Timeout,
	})

	// Save to history
	historyEntry := &models.History{
		Method:          params.Method,
		URL:             params.URL,
		RequestHeaders:  services.BuildRequestHeadersJSON(params.Headers),
		RequestBody:     params.Body,
	}

	if resp != nil {
		statusCode := resp.StatusCode
		historyEntry.StatusCode = &statusCode
		historyEntry.ResponseHeaders = services.BuildResponseHeadersJSON(resp.Headers)
		historyEntry.ResponseBody = resp.Body
		historyEntry.DurationMs = &resp.Duration
	}

	h.history.Create(historyEntry)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CancelRequest cancels a running request
func (h *RequestHandler) CancelRequest(tabID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if cancel, ok := h.cancelFuncs[tabID]; ok {
		cancel()
		delete(h.cancelFuncs, tabID)
	}
}
