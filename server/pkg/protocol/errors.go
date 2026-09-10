package protocol

// ErrorPayload is sent as system/error when a command cannot be completed.
// Its fields are generated from protocol/payloads.schema.json.
type ErrorPayload = ProtocolErrorPayload

func NewError(code ErrorCode, message string) ErrorPayload {
	return ErrorPayload{
		Title:     "Protocol error",
		Code:      code,
		Message:   message,
		Retryable: IsRetryable(code),
	}
}
