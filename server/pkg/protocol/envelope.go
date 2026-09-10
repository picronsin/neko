package protocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

// ValidatePayload applies the runtime checks for command payloads whose
// schema has required fields. The JSON Schema remains the documentation and
// exchange contract; this small validator keeps the server independent from
// a third-party schema runtime in the hot WebSocket path.
func ValidatePayload(event string, raw json.RawMessage) error {
	requireObject := func() (map[string]json.RawMessage, error) {
		if len(bytes.TrimSpace(raw)) == 0 {
			return nil, errors.New("payload is required")
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(raw, &object); err != nil || object == nil {
			return nil, errors.New("payload must be an object")
		}
		return object, nil
	}
	requireString := func(object map[string]json.RawMessage, field string, allowEmpty bool) error {
		value, ok := object[field]
		if !ok {
			return fmt.Errorf("payload field %q is required", field)
		}
		var text string
		if err := json.Unmarshal(value, &text); err != nil || (!allowEmpty && strings.TrimSpace(text) == "") {
			return fmt.Errorf("payload field %q must be a non-empty string", field)
		}
		return nil
	}
	requireNumber := func(object map[string]json.RawMessage, field string) error {
		value, ok := object[field]
		if !ok {
			return fmt.Errorf("payload field %q is required", field)
		}
		var number json.Number
		if err := json.Unmarshal(value, &number); err != nil || number.String() == "" {
			return fmt.Errorf("payload field %q must be a number", field)
		}
		return nil
	}

	switch event {
	case SignalOffer, SignalAnswer:
		object, err := requireObject()
		if err != nil {
			return err
		}
		return requireString(object, "sdp", false)
	case SignalCandidate:
		object, err := requireObject()
		if err != nil {
			return err
		}
		return requireString(object, "candidate", false)
	case ControlRenew:
		object, err := requireObject()
		if err != nil {
			return err
		}
		return requireNumber(object, "epoch")
	case ScreenSet:
		object, err := requireObject()
		if err != nil {
			return err
		}
		for _, field := range []string{"width", "height", "rate"} {
			if err := requireNumber(object, field); err != nil {
				return err
			}
		}
	case ClipboardSet, ChatMessage:
		object, err := requireObject()
		if err != nil {
			return err
		}
		return requireString(object, "text", event == ClipboardSet)
	case KeyboardMap:
		object, err := requireObject()
		if err != nil {
			return err
		}
		return requireString(object, "layout", false)
	case SignalRequest, KeyboardModifiers:
		_, err := requireObject()
		return err
	}
	return nil
}
