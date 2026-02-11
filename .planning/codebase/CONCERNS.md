# Codebase Concerns

**Analysis Date:** 2026-02-11

## Tech Debt

**Unimplemented Request Execution in Frontend:**
- Issue: The core request sending functionality is stubbed with a placeholder that simulates responses instead of calling the backend
- Files: `frontend/src/composables/useRequest.ts` (lines 72-84, 108-109)
- Impact: Application cannot actually send HTTP requests; users see mocked 200 OK responses instead of real data
- Fix approach: Uncomment the Wails backend calls in `sendRequest()` and `cancelRequest()` functions and integrate with the Go backend's RequestHandler

**Unhandled JSON Marshal Errors:**
- Issue: Two utility functions silently ignore JSON marshaling errors using blank identifier `_`
- Files: `internal/services/http_client.go` (lines 478, 484 in `BuildRequestHeadersJSON` and `BuildResponseHeadersJSON`)
- Impact: If marshaling fails, empty JSON strings are returned, potentially causing silent data loss in history records
- Fix approach: Either handle errors explicitly or log them; at minimum, return error values instead of discarding them

**Missing Deflate Reader Error Handling:**
- Issue: Deflate decompression reader is not checked for errors, unlike gzip and brotli
- Files: `internal/services/http_client.go` (line 448)
- Impact: If response body is corrupted with deflate encoding, no error is raised; subsequent `io.ReadAll` may fail with unclear error context
- Fix approach: Wrap deflate reader in error checking similar to gzip (lines 441-445)

**Unbounded File Reading for Binary Body:**
- Issue: Binary request bodies read entire files into memory with `io.ReadAll()` without size limits
- Files: `internal/services/http_client.go` (lines 359-363)
- Impact: Large file uploads can cause memory exhaustion; no maximum size validation exists
- Fix approach: Implement maximum file size check before reading; consider streaming large files instead of buffering

**JSON Unmarshal Silently Fails for Form Data:**
- Issue: When parsing form-data or urlencoded body, JSON unmarshaling errors are silently ignored
- Files: `internal/services/http_client.go` (lines 300, 340 with `err == nil` checks)
- Impact: If form data is malformed, the entire body is ignored and request is sent empty; no error feedback to user
- Fix approach: Propagate unmarshal errors or log them; validate form data format before processing

## Known Bugs

**Tab State Dirty Flag Comparison Using JSON.stringify:**
- Symptoms: Deep object comparison using `JSON.stringify()` may incorrectly detect changes due to key ordering or object reference changes
- Files: `frontend/src/stores/tabs.ts` (lines 47, 50 in `computeDirty()`)
- Trigger: Modify tab content and switch tabs; dirty state may not update consistently
- Workaround: Manual tab save may be needed; monitor browser console for unexpected state changes
- Impact: Users may lose unsaved changes or see incorrect dirty indicators

**Memory Leak in Keyboard Action Listeners:**
- Symptoms: Keyboard action callbacks stored in module-level Map never fully cleaned up if components unmount without unsubscribing
- Files: `frontend/src/composables/useKeyboardActions.ts` (line 6, global listeners map)
- Trigger: Rapidly mount/unmount components using keyboard actions without properly calling unsubscribe function
- Workaround: Ensure all `onKeyboardAction()` calls capture and invoke the returned unsubscribe function
- Impact: Memory usage grows over time in long-running sessions

**Missing Variable Interpolation Validation:**
- Symptoms: Variable replacement using regex pattern `\{\{(\w+)\}\}` only matches word characters; variable names with special characters are ignored
- Files: `frontend/src/stores/environment.ts` (line 43 in `replaceVariables()`)
- Trigger: Try to use environment variables with underscores, hyphens, or dots in names
- Workaround: Limit variable names to alphanumeric and underscore only
- Impact: Environment variables with common naming patterns (e.g., `api_key`, `base-url`) won't interpolate

## Security Considerations

**User-Agent Spoofing and Fingerprint Masking:**
- Risk: Application implements browser fingerprint spoofing to evade detection, which may be used for malicious purposes
- Files: `internal/services/http_client.go` (lines 30, 93, 395-423)
- Current mitigation: Feature exists but no usage restrictions or logging
- Recommendations: Add user consent/warning when browser spoofing is enabled; log all spoofed requests; consider disabling this feature for general users

**File Upload Path Disclosure:**
- Risk: File upload errors expose full file paths in error messages and request bodies
- Files: `internal/services/http_client.go` (line 310, 354)
- Current mitigation: None
- Recommendations: Sanitize error messages to not expose sensitive paths; only send filename, not full path, in form data

**Unvalidated File Paths for Binary Bodies:**
- Risk: Application will open and read any file path provided by user without path traversal validation
- Files: `internal/services/http_client.go` (lines 308, 352)
- Current mitigation: None (relies on OS file permissions)
- Recommendations: Implement path validation to prevent `../` traversal; validate against allowed directories if needed

**No HTTPS Verification Bypass Detection:**
- Risk: uTLS TLS fingerprinting could potentially be used with invalid certificates
- Files: `internal/services/http_client.go` (lines 91-99)
- Current mitigation: Standard Go certificate validation still applies
- Recommendations: Ensure certificate validation cannot be bypassed; audit uTLS integration for security implications

**System Proxy Password Exposure:**
- Risk: Windows system proxy configuration may contain plaintext passwords in registry
- Files: `internal/services/http_client.go` (lines 240-272)
- Current mitigation: Password is part of URL but not logged
- Recommendations: Never log proxy URLs; consider masking passwords in debug output

## Performance Bottlenecks

**Entire Response Body Loaded into Memory:**
- Problem: All response bodies are read entirely into memory before returning
- Files: `internal/services/http_client.go` (line 453)
- Cause: Design uses `io.ReadAll()` on potentially large responses
- Improvement path: Implement streaming responses for large payloads; add configurable response size limits

**Form-Data Buffering Entire Multipart in Memory:**
- Problem: Large multipart form uploads are fully buffered in memory before sending
- Files: `internal/services/http_client.go` (lines 302-335)
- Cause: Using `bytes.Buffer` for multipart writer instead of streaming
- Improvement path: Stream form data directly to network without intermediate buffering

**Synchronous Response Decompression:**
- Problem: All three decompression algorithms (gzip, deflate, brotli) block the thread during decompression
- Files: `internal/services/http_client.go` (lines 440-450)
- Cause: Linear processing of response before returning
- Improvement path: For large responses, consider async decompression or streaming

**Deep Object Cloning with JSON.stringify:**
- Problem: Tab duplication and saving uses `JSON.parse(JSON.stringify())` for deep cloning
- Files: `frontend/src/stores/tabs.ts` (line 315)
- Cause: Inefficient compared to structured cloning
- Improvement path: Use `structuredClone()` instead; or implement efficient deep clone utility

**No Lazy Loading of API Module:**
- Problem: Environment store lazily loads API module to avoid circular dependencies, adding async overhead
- Files: `frontend/src/stores/appState.ts` (lines 6-12)
- Cause: Circular dependency between stores and API layer
- Improvement path: Restructure imports to eliminate circular dependencies; use explicit dependency injection

## Fragile Areas

**HTTP Client Transport Configuration:**
- Files: `internal/services/http_client.go` (lines 45-123, 194-214)
- Why fragile: Complex custom uTLS transport creates new transport for every HTTP and HTTP/2 request, losing connection reuse; proxy configuration is rebuilt on every SetUseSystemProxy call
- Safe modification: Test proxy scenarios (SOCKS, HTTP, HTTPS, authenticated); verify connection pooling still works; test HTTP/2 with push enabled
- Test coverage: No unit tests visible for transport layer; proxy configuration untested

**Tab State Management Complex Logic:**
- Files: `frontend/src/stores/tabs.ts` (lines 31-58, 94-185)
- Why fragile: Multiple overlapping concerns: preview tabs, dirty state detection, original state tracking, and requestId associations; many similar functions with subtle differences (openRequest vs previewRequest)
- Safe modification: Add comprehensive unit tests for dirty detection with various input combinations; test preview -> pinned -> closed flow; verify requestId synchronization with backend
- Test coverage: No test files found for tab store; dirty detection logic untested

**Variable Interpolation Regex:**
- Files: `frontend/src/stores/environment.ts` (line 43)
- Why fragile: Single regex pattern handles all variable replacement; no escaping support; no way to include literal `{{` in values
- Safe modification: Add optional escape sequences (e.g., `\{\{` for literal); consider lexer-based approach for complex cases
- Test coverage: No tests for variable replacement edge cases

## Scaling Limits

**In-Memory Tab State:**
- Current capacity: Tested up to ~100 tabs (estimate)
- Limit: With each tab storing full request state, hundreds of tabs will cause memory issues; no pagination/virtualization implemented
- Scaling path: Implement virtual scrolling for tab bar; lazy-load tab content; persist to database

**Response History Unlimited Growth:**
- Current capacity: Grows with every request executed
- Limit: No size limits or archival; database could grow unbounded
- Scaling path: Implement request history pruning (e.g., keep only last 1000 items); add database indexes on timestamps; implement pagination

**Single HTTP Client Instance:**
- Current capacity: Handles sequential requests with connection reuse
- Limit: No connection pooling limits; concurrent requests share single HTTP client without rate limiting
- Scaling path: Consider implementing request queue with max concurrency; add per-domain connection limits

## Dependencies at Risk

**uTLS (Refraction Networking):**
- Risk: Maintenance uncertainty; used for fingerprinting which could violate terms of service of many APIs
- Impact: If library unmaintained, vulnerabilities in TLS handling could affect security
- Migration plan: Document requirement clearly; consider making feature opt-in; prepare for removal if legal concerns arise

**Brotli Compression Support:**
- Risk: Additional decompression library adds dependency; not all platforms may have it readily available
- Impact: If brotli dependency breaks, responses with brotli encoding fail
- Migration plan: Wrap brotli in try-catch; fallback to raw body if decompression fails

## Missing Critical Features

**No Response Size Limits:**
- Problem: Can attempt to download infinitely large responses, causing memory exhaustion or UI freeze
- Blocks: Safe handling of large files; protection against DOS-like scenarios

**No Request Rate Limiting:**
- Problem: User can send unlimited requests without throttling
- Blocks: Protecting against accidental DOS of target servers; managing API rate limits

**No Request Validation:**
- Problem: Malformed URLs, invalid HTTP methods, invalid headers are not validated before sending
- Blocks: Early error detection; helpful error messages to users

**No Timeout Configuration per Request:**
- Problem: Only global timeout setting; cannot customize per-request timeouts
- Blocks: Handling slow APIs that need higher timeouts; implementing request priorities

**No SSL/TLS Certificate Management:**
- Problem: No UI for managing client certificates or CA certificates
- Blocks: Testing APIs requiring mutual TLS authentication

## Test Coverage Gaps

**Frontend Store Unit Tests:**
- What's not tested: Tab state management (dirty detection, preview/pin logic), environment variable interpolation, keyboard action emitter
- Files: `frontend/src/stores/tabs.ts`, `frontend/src/stores/environment.ts`, `frontend/src/composables/useKeyboardActions.ts`
- Risk: Tab management bugs could go unnoticed; variable interpolation silently fails; memory leaks in keyboard listeners
- Priority: High

**HTTP Client Integration Tests:**
- What's not tested: Multipart form handling, decompression (gzip/deflate/brotli), proxy configuration, error scenarios
- Files: `internal/services/http_client.go`
- Risk: Form uploads corrupted, responses corrupted, proxy connections fail silently
- Priority: High

**Backend Request Handler Edge Cases:**
- What's not tested: Request cancellation mid-flight, timeout behavior, history recording with failed requests, concurrent request handling
- Files: `internal/handlers/request_handler.go`
- Risk: Race conditions in cancellation, incomplete history records, resource leaks
- Priority: Medium

**Frontend-Backend Integration:**
- What's not tested: The actual useRequest().sendRequest() call with real Wails backend (currently stubbed)
- Files: `frontend/src/composables/useRequest.ts`
- Risk: Blocker for functional testing; stub masks backend integration issues
- Priority: Critical

**Keyboard Shortcuts:**
- What's not tested: Keyboard action emitter/listener pattern, unsubscribe cleanup
- Files: `frontend/src/composables/useKeyboardActions.ts`
- Risk: Memory leaks from unreleased listeners; unexpected keyboard behavior
- Priority: Medium

**Error Handling Edge Cases:**
- What's not tested: Malformed responses, network errors, timeout handling, file not found errors
- Files: `internal/services/http_client.go`, `frontend/src/composables/useRequest.ts`
- Risk: Users see confusing error messages; error recovery doesn't work
- Priority: Medium

---

*Concerns audit: 2026-02-11*
