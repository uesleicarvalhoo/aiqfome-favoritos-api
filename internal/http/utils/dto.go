package utils

// APIError data, all errors of our HTTP API returns this format
// @Description default error API format
type APIError struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Method    string         `json:"method,omitempty"`
	Path      string         `json:"path,omitempty"`
	Timestamp string         `json:"timestamp,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
}
