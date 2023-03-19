package payloads

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse is used as a payload for all errors
type ErrorResponse struct {
	// HTTP status code
	HTTPStatusCode int `json:"-"`
	// user facing error message
	Message string `json:"msg"`
	// full root cause
	Error string `json:"error"`
}

func (e ErrorResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func newErrorResponse(ctx context.Context, status int, userMsg string, err error) ErrorResponse {
	return ErrorResponse{
		HTTPStatusCode: status,
		Message:        userMsg,
		Error:          err.Error(),
	}
}

func NewInvalidRequestError(ctx context.Context, message string, err error) ErrorResponse {
	message = fmt.Sprintf("Invalid request: %s", message)
	return newErrorResponse(ctx, http.StatusBadRequest, message, err)
}

func NewNotFoundError(ctx context.Context, message string, err error) ErrorResponse {
	message = fmt.Sprintf("Not found: %s", message)
	return newErrorResponse(ctx, http.StatusNotFound, message, err)
}

func NewDAOError(ctx context.Context, message string, err error) ErrorResponse {
	message = fmt.Sprintf("DAO error: %s", message)
	return newErrorResponse(ctx, http.StatusInternalServerError, message, err)
}

func NewRenderError(ctx context.Context, message string, err error) ErrorResponse {
	message = fmt.Sprintf("Rendering error: %s", message)
	return newErrorResponse(ctx, http.StatusInternalServerError, message, err)
}
