// internal/model/error.go
package model

// ErrorResponse represents an API error response
type ErrorResponse struct {
    Error string `json:"error" example:"Error message"`
}