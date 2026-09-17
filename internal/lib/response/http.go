package response

import (
	"net/http"

	"github.com/prabhatlabs/go-mail-server/internal/lib/logger"
)

// OK sends a 200 success envelope.
func OK(w http.ResponseWriter, message string, data any) {
	SendJsonResponse(w, http.StatusOK, SuccessResponse{Message: message, Data: data})
}

// Created sends a 201 success envelope.
func Created(w http.ResponseWriter, message string, data any) {
	SendJsonResponse(w, http.StatusCreated, SuccessResponse{Message: message, Data: data})
}

// ErrorText maps status codes to their canonical error label.
func ErrorText(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusConflict:
		return "Conflict"
	default:
		return "Internal Server Error"
	}
}

// Error sends an error envelope with an explicit status.
func Error(w http.ResponseWriter, status int, message string) {
	SendJsonResponse(w, status, ErrorResponse{Error: ErrorText(status), Message: message})
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, message)
}

func InternalServerError(w http.ResponseWriter, err error, message string) {
	logger.ServerLog(err, message)
	Error(w, http.StatusInternalServerError, message)
}
