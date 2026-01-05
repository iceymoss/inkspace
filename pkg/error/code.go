package errors

import "net/http"

// System-level error codes (1000-1999 range), code >= 2000 it s a business error，customized by the module implementer
const (
	// General client errors (1000-1099)
	ErrBadRequest       = 1000 // HTTP 400
	ErrInvalidInput     = 1001 // HTTP 400
	ErrMissingParameter = 1002 // HTTP 400
	ErrInvalidJSON      = 1003 // HTTP 400

	// Authentication errors (1100-1199)
	ErrUnauthenticated  = 1100 // HTTP 401
	ErrInvalidToken     = 1101 // HTTP 401
	ErrTokenExpired     = 1102 // HTTP 401
	ErrInvalidSignature = 1103 // HTTP 401

	// Permission errors (1200-1299)
	ErrForbidden        = 1200 // HTTP 403
	ErrInsufficientPriv = 1201 // HTTP 403
	ErrReadOnlyMode     = 1202 // HTTP 403

	// Resource errors (1300-1399)
	ErrNotFound         = 1300 // HTTP 404
	ErrResourceNotFound = 1301 // HTTP 404
	ErrEndpointRemoved  = 1302 // HTTP 410

	// Conflict errors (1400-1499) - Focus area
	ErrConflict        = 1400 // HTTP 409
	ErrResourceExists  = 1401 // HTTP 409 - Resource already exists
	ErrVersionConflict = 1402 // HTTP 409 - Version conflict
	ErrResourceLocked  = 1403 // HTTP 423 - Resource locked

	// Rate limiting (1500-1599)
	ErrTooManyRequests   = 1500 // HTTP 429
	ErrRateLimitExceeded = 1501 // HTTP 429

	// Server errors (1700-1799)
	ErrInternalServer     = 1700 // HTTP 500
	ErrDatabaseError      = 1701 // HTTP 500
	ErrServiceUnavailable = 1702 // HTTP 503
	ErrTimeout            = 1703 // HTTP 504

	// code >= 2000 it s a business error，customized by the module implementer
)

// Initialize all base error codes
func init() {
	registerMessages()
	registerHTTPStatusMappings()
}

func registerMessages() {
	// General client errors
	RegisterMessage(ErrBadRequest, "Invalid request")
	RegisterMessage(ErrInvalidInput, "Invalid input parameters")
	RegisterMessage(ErrMissingParameter, "Missing required parameter: %s")
	RegisterMessage(ErrInvalidJSON, "Request body is not valid JSON")

	// Authentication errors
	RegisterMessage(ErrUnauthenticated, "Authentication required")
	RegisterMessage(ErrInvalidToken, "Invalid access token")
	RegisterMessage(ErrTokenExpired, "Access token expired")
	RegisterMessage(ErrInvalidSignature, "Signature verification failed")

	// Permission errors
	RegisterMessage(ErrForbidden, "Access denied")
	RegisterMessage(ErrInsufficientPriv, "Insufficient privileges")
	RegisterMessage(ErrReadOnlyMode, "System under maintenance (read-only mode)")

	// Resource errors
	RegisterMessage(ErrNotFound, "Resource not found")
	RegisterMessage(ErrResourceNotFound, "Requested resource not found: %s")
	RegisterMessage(ErrEndpointRemoved, "API endpoint deprecated")

	// Conflict errors (focus)
	RegisterMessage(ErrConflict, "Resource state conflict")
	RegisterMessage(ErrResourceExists, "Resource already exists") // When creating new resource
	RegisterMessage(ErrVersionConflict, "Resource version conflict")
	RegisterMessage(ErrResourceLocked, "Resource is locked")

	// Rate limiting
	RegisterMessage(ErrTooManyRequests, "Too many requests")
	RegisterMessage(ErrRateLimitExceeded, "API rate limit exceeded")

	// Server errors
	RegisterMessage(ErrInternalServer, "Internal server error")
	RegisterMessage(ErrDatabaseError, "Database operation failed")
	RegisterMessage(ErrServiceUnavailable, "Service temporarily unavailable")
	RegisterMessage(ErrTimeout, "Processing timeout")
}

func registerHTTPStatusMappings() {
	// General client errors -> 400
	RegisterHTTPStatus(ErrBadRequest, http.StatusBadRequest)
	RegisterHTTPStatus(ErrInvalidInput, http.StatusBadRequest)
	RegisterHTTPStatus(ErrMissingParameter, http.StatusBadRequest)
	RegisterHTTPStatus(ErrInvalidJSON, http.StatusBadRequest)

	// Authentication errors -> 401
	RegisterHTTPStatus(ErrUnauthenticated, http.StatusUnauthorized)
	RegisterHTTPStatus(ErrInvalidToken, http.StatusUnauthorized)
	RegisterHTTPStatus(ErrTokenExpired, http.StatusUnauthorized)
	RegisterHTTPStatus(ErrInvalidSignature, http.StatusUnauthorized)

	// Permission errors -> 403
	RegisterHTTPStatus(ErrForbidden, http.StatusForbidden)
	RegisterHTTPStatus(ErrInsufficientPriv, http.StatusForbidden)
	RegisterHTTPStatus(ErrReadOnlyMode, http.StatusForbidden)

	// Resource errors -> 404/410
	RegisterHTTPStatus(ErrNotFound, http.StatusNotFound)
	RegisterHTTPStatus(ErrResourceNotFound, http.StatusNotFound)
	RegisterHTTPStatus(ErrEndpointRemoved, http.StatusGone)

	// Conflict errors (focus)
	RegisterHTTPStatus(ErrConflict, http.StatusConflict)
	RegisterHTTPStatus(ErrResourceExists, http.StatusConflict)  // 409 for resource exists
	RegisterHTTPStatus(ErrVersionConflict, http.StatusConflict) // 409
	RegisterHTTPStatus(ErrResourceLocked, http.StatusLocked)    // 423

	// Rate limiting -> 429
	RegisterHTTPStatus(ErrTooManyRequests, http.StatusTooManyRequests)
	RegisterHTTPStatus(ErrRateLimitExceeded, http.StatusTooManyRequests)

	// Server errors -> 500/503/504
	RegisterHTTPStatus(ErrInternalServer, http.StatusInternalServerError)
	RegisterHTTPStatus(ErrDatabaseError, http.StatusInternalServerError)
	RegisterHTTPStatus(ErrServiceUnavailable, http.StatusServiceUnavailable)
	RegisterHTTPStatus(ErrTimeout, http.StatusGatewayTimeout)
}
