package protocol

// ErrorPayload is sent as system/error when a command cannot be completed.
// Code is stable for clients; Message is diagnostic and may change.
type ErrorPayload struct {
	Title     string    `json:"title"`
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Retryable bool      `json:"retryable"`
}

func NewError(code ErrorCode, message string) ErrorPayload {
	return ErrorPayload{
		Title:     "Protocol error",
		Code:      code,
		Message:   message,
		Retryable: IsRetryable(code),
	}
}
