package restapi

// Envelope is the shared REST success response shape for newly added APIs.
type Envelope struct {
	Data any `json:"data,omitempty"`
	Meta any `json:"meta,omitempty"`
}

// ErrorEnvelope is the shared REST error response shape for newly added APIs.
type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody describes a transport-level API error.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}
