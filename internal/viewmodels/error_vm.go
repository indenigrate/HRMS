package viewmodels

// ErrorResponse represents a standard error response.
type ErrorResponse struct {
	Error string `json:"error" example:"a description of the error"`
}
