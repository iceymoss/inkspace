package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Error type constants
const (
	ErrBusiness = "business" // Business error type
	ErrInternal = "internal" // Internal error type
)

// Error represents a structured application error
type Error struct {
	Type       string // Error type: business/internal
	Code       int    // Business error code
	Message    string // Error message
	Err        error  // Original error (for internal errors)
	HTTPStatus int    // HTTP status code
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Type == ErrInternal && e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Error registry
var (
	errorMessages     = make(map[int]string) // Map of error codes to messages
	httpStatusMapping = make(map[int]int)    // Map of error codes to HTTP statuses
	lock              sync.RWMutex           // Registry lock
)

// RegisterMessage registers a message for an error code
func RegisterMessage(code int, message string) {
	lock.Lock()
	defer lock.Unlock()
	errorMessages[code] = message
}

// RegisterHTTPStatus registers an HTTP status for an error code
func RegisterHTTPStatus(code, httpStatus int) {
	lock.Lock()
	defer lock.Unlock()
	httpStatusMapping[code] = httpStatus
}

// getMessage retrieves the message for an error code
func getMessage(code int) string {
	lock.RLock()
	defer lock.RUnlock()

	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "Unknown error"
}

// getHTTPStatus retrieves the HTTP status for an error code
func getHTTPStatus(code int) int {
	lock.RLock()
	defer lock.RUnlock()

	if status, ok := httpStatusMapping[code]; ok {
		return status
	}

	// Default mapping rules
	switch {
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest
	case code >= 2000 && code < 3000:
		return http.StatusForbidden
	case code >= 3000 && code < 4000:
		return http.StatusConflict
	case code >= 9000:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

/**********************************
 * Error Creation Functions
 **********************************/

// New creates an error with default HTTP status
func New(typ string, code int, message string) *Error {
	return &Error{
		Type:       typ,
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatus(code),
	}
}

// NewWithStatus creates an error with custom HTTP status
func NewWithStatus(typ string, code int, message string, httpStatus int) *Error {
	return &Error{
		Type:       typ,
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

/**********************************
 * Convenience Functions
 **********************************/

// Business creates a business error
func Business(code int, message ...string) *Error {
	msg := getMessage(code)
	if len(message) > 0 {
		msg = message[0]
	}
	return New(ErrBusiness, code, msg)
}

// BusinessWithStatus creates a business error with custom HTTP status
func BusinessWithStatus(code, httpStatus int, message ...string) *Error {
	msg := getMessage(code)
	if len(message) > 0 {
		msg = message[0]
	}
	return NewWithStatus(ErrBusiness, code, msg, httpStatus)
}

// Internal creates an internal error
func Internal(err error, message string) *Error {
	return &Error{
		Type:       ErrInternal,
		Message:    message,
		Err:        err,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// Wrap wraps an existing error
func Wrap(err error, message string) *Error {
	if e, ok := err.(*Error); ok {
		return &Error{
			Type:       e.Type,
			Code:       e.Code,
			Message:    fmt.Sprintf("%s: %s", message, e.Message),
			Err:        e.Err,
			HTTPStatus: e.HTTPStatus,
		}
	}
	return Internal(err, message)
}

/**********************************
 * Error Handling
 **********************************/

// Handle processes an error response
func Handle(c *gin.Context, err error) {
	// Convert to internal error if not our error type
	appErr, ok := err.(*Error)
	if !ok {
		appErr = Internal(err, "Unhandled error type")
	}

	// Log internal errors
	if appErr.Type == ErrInternal {
		logError(c, appErr)
	}

	// Return response
	c.JSON(appErr.HTTPStatus, gin.H{
		"code":    appErr.Code,
		"message": appErr.Message,
	})
}

// logError logs internal errors with stack trace
func logError(c *gin.Context, err *Error) {
	// Get request information
	path := c.Request.URL.Path
	method := c.Request.Method
	ip := c.ClientIP()

	// Get call stack
	stack := getStack(3) // Skip 3 call frames

	// Build log message
	logMsg := fmt.Sprintf("[ERROR] %s %s %s - %s", method, path, ip, err.Message)
	if err.Err != nil {
		logMsg += fmt.Sprintf(": %v", err.Err)
	}
	logMsg += "\nStack:\n" + stack

	// In production, use a proper logger instead of fmt
	fmt.Println(logMsg)
}

// getStack retrieves the call stack
func getStack(skip int) string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	// Split lines
	lines := strings.Split(stack, "\n")

	// Skip unnecessary frames
	if skip*2+1 < len(lines) {
		return strings.Join(lines[skip*2:], "\n")
	}

	return stack
}
