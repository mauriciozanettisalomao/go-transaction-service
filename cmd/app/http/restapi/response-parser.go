package restapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// response represents a response body format
type response struct {
	Data     any      `json:"data,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(data any, metadata Metadata) response {
	return response{
		Data:     data,
		Metadata: metadata,
	}
}

// Metadata represents metadata for a paginated response
type Metadata struct {
	Limit int `json:"limit" example:"10"`
}

// newMeta is a helper function to create metadata for a paginated response
func newMeta(limit int) Metadata {
	return Metadata{
		Limit: limit,
	}
}

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	// domain.ErrDataNotFound:               http.StatusNotFound,
	// domain.ErrConflictingData:            http.StatusConflict,
	// domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	// domain.ErrUnauthorized:               http.StatusUnauthorized,
	// domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	// domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	// domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	// domain.ErrInvalidToken:               http.StatusUnauthorized,
	// domain.ErrExpiredToken:               http.StatusUnauthorized,
	// domain.ErrForbidden:                  http.StatusForbidden,
	// domain.ErrNoUpdatedData:              http.StatusBadRequest,
	// domain.ErrInsufficientStock:          http.StatusBadRequest,
	// domain.ErrInsufficientPayment:        http.StatusBadRequest,
}

// validationError sends an error response for some specific request validation error
func validationError(ctx *gin.Context, err error) {
	errMsgs := parseError(err)
	errRsp := newErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// handleError determines the status code of an error and returns a JSON response with the error message and status code
func handleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// handleAbort sends an error response and aborts the request with the specified status code and error message
func handleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

// parseError parses error messages from the error object and returns a slice of error messages
func parseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

// errorResponse represents an error response body format
type errorResponse struct {
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Messages: errMsgs,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any, metadata Metadata) {
	rsp := newResponse(data, metadata)
	ctx.JSON(http.StatusOK, rsp)
}

// handleSuccess sends a success response with the specified status code and optional data
func handleCreatedSuccess(ctx *gin.Context, data any, metadata Metadata) {
	rsp := newResponse(data, metadata)
	ctx.JSON(http.StatusCreated, rsp)
}
