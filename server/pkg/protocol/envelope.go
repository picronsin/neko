package protocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

// Envelope is the canonical real-time message shape. Payload is omitted for
// events that do not carry data; null and flat top-level fields are invalid.
type Envelope struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func (message *Envelope) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	for key := range fields {
		if key != "event" && key != "payload" {
			return errors.New("websocket message contains unsupported top-level fields")
		}
	}

	type plain Envelope
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if strings.TrimSpace(value.Event) == "" {
		return errors.New("websocket message event is required")
	}
	if payload, ok := fields["payload"]; ok {
		trimmed := bytes.TrimSpace(payload)
		if bytes.Equal(trimmed, []byte("null")) {
			return errors.New("websocket message payload must be omitted when empty")
		}
		if len(trimmed) == 0 || (trimmed[0] != '{' && trimmed[0] != '[') {
			return errors.New("websocket message payload must be an object or array")
		}
	}

	*message = Envelope(value)
	return nil
}

// ValidateEvent is intentionally separate from envelope parsing. Plugins may
// register private events through the WebSocket handler extension point, while
// the public protocol can still validate that an event belongs to the catalog.
func ValidateEvent(event string) error {
	if !IsKnownEvent(event) {
		return errors.New("unknown websocket event: " + event)
	}
	return nil
}
