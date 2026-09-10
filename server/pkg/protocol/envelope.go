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
	requireArray := func() ([]json.RawMessage, error) {
		if len(bytes.TrimSpace(raw)) == 0 {
			return nil, errors.New("payload is required")
		}
		var array []json.RawMessage
		if err := json.Unmarshal(raw, &array); err != nil || array == nil {
			return nil, errors.New("payload must be an array")
		}
		return array, nil
	}
	requireNoPayload := func() error {
		if len(bytes.TrimSpace(raw)) != 0 {
			return errors.New("payload must be omitted")
		}
		return nil
	}
	rejectUnknown := func(object map[string]json.RawMessage, fields ...string) error {
		allowed := make(map[string]struct{}, len(fields))
		for _, field := range fields {
			allowed[field] = struct{}{}
		}
		for field := range object {
			if _, ok := allowed[field]; !ok {
				return fmt.Errorf("payload contains unsupported field %q", field)
			}
		}
		return nil
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
	requireBool := func(object map[string]json.RawMessage, field string) error {
		value, ok := object[field]
		if !ok {
			return fmt.Errorf("payload field %q is required", field)
		}
		var result bool
		if err := json.Unmarshal(value, &result); err != nil {
			return fmt.Errorf("payload field %q must be a boolean", field)
		}
		return nil
	}
	requireStringArray := func(object map[string]json.RawMessage, field string) error {
		value, ok := object[field]
		if !ok {
			return nil
		}
		var values []string
		if err := json.Unmarshal(value, &values); err != nil {
			return fmt.Errorf("payload field %q must be a string array", field)
		}
		for _, value := range values {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("payload field %q must not contain empty strings", field)
			}
		}
		return nil
	}
	requireInteger := func(object map[string]json.RawMessage, field string) error {
		if err := requireNumber(object, field); err != nil {
			return err
		}
		var number int64
		if err := json.Unmarshal(object[field], &number); err != nil || number < 0 {
			return fmt.Errorf("payload field %q must be a non-negative integer", field)
		}
		return nil
	}
	validatePeerVideo := func(value json.RawMessage) error {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(value, &object); err != nil || object == nil {
			return errors.New("payload video must be an object")
		}
		if err := rejectUnknown(object, "disabled", "auto", "selector"); err != nil {
			return err
		}
		for _, field := range []string{"disabled", "auto"} {
			if _, ok := object[field]; ok {
				if err := requireBool(object, field); err != nil {
					return err
				}
			}
		}
		return nil
	}
	validatePeerAudio := func(value json.RawMessage) error {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(value, &object); err != nil || object == nil {
			return errors.New("payload audio must be an object")
		}
		if err := rejectUnknown(object, "disabled"); err != nil {
			return err
		}
		if _, ok := object["disabled"]; ok {
			return requireBool(object, "disabled")
		}
		return nil
	}

	switch event {
	case ClientHeartbeat, SignalRestart, ControlRelease:
		return requireNoPayload()
	case SignalRequest:
		object, err := requireObject()
		if err != nil {
			return err
		}
		if err := rejectUnknown(object, "video_codecs", "video", "audio"); err != nil {
			return err
		}
		if err := requireStringArray(object, "video_codecs"); err != nil {
			return err
		}
		if value, ok := object["video"]; ok {
			if err := validatePeerVideo(value); err != nil {
				return err
			}
		}
		if value, ok := object["audio"]; ok {
			if err := validatePeerAudio(value); err != nil {
				return err
			}
		}
		return nil
	case SignalOffer, SignalAnswer:
		object, err := requireObject()
		if err != nil {
			return err
		}
		if err := rejectUnknown(object, "sdp"); err != nil {
			return err
		}
		return requireString(object, "sdp", false)
	case SignalCandidate:
		object, err := requireObject()
		if err != nil {
			return err
		}
		if err := rejectUnknown(object, "candidate", "sdpMid", "sdpMLineIndex", "usernameFragment"); err != nil {
			return err
		}
		return requireString(object, "candidate", false)
	case SignalVideo:
		_, err := requireObject()
		if err != nil {
			return err
		}
		return validatePeerVideo(json.RawMessage(raw))
	case SignalAudio:
		_, err := requireObject()
		if err != nil {
			return err
		}
		return validatePeerAudio(json.RawMessage(raw))
	case ControlRequest:
		return requireNoPayload()
	case ControlRenew:
		object, err := requireObject()
		if err != nil {
			return err
		}
		if err := rejectUnknown(object, "epoch"); err != nil {
			return err
		}
		return requireInteger(object, "epoch")
	case ScreenSet:
		object, err := requireObject()
		if err != nil {
			return err
		}
		if err := rejectUnknown(object, "width", "height", "rate"); err != nil {
			return err
		}
		for _, field := range []string{"width", "height", "rate"} {
			if err := requireInteger(object, field); err != nil {
				return err
			}
		}
	case ClipboardSet, ChatMessage:
		object, err := requireObject()
		if err != nil {
			return err
		}
		if err := rejectUnknown(object, "text"); err != nil {
			return err
		}
		return requireString(object, "text", event == ClipboardSet)
	case KeyboardMap:
		object, err := requireObject()
		if err != nil {
			return err
		}
		if err := rejectUnknown(object, "layout", "variant"); err != nil {
			return err
		}
		return requireString(object, "layout", false)
	case KeyboardModifiers:
		_, err := requireObject()
		return err
	case SystemLogs:
		array, err := requireArray()
		if err != nil {
			return err
		}
		for _, item := range array {
			object := map[string]json.RawMessage{}
			if err := json.Unmarshal(item, &object); err != nil || object == nil {
				return errors.New("payload log entry must be an object")
			}
			if err := rejectUnknown(object, "level", "fields", "message"); err != nil {
				return err
			}
			for _, field := range []string{"level", "message"} {
				if err := requireString(object, field, field == "message"); err != nil {
					return err
				}
			}
			if fields, ok := object["fields"]; ok {
				var value map[string]any
				if err := json.Unmarshal(fields, &value); err != nil || value == nil {
					return errors.New("payload field \"fields\" must be an object")
				}
			}
		}
		return nil
	case SendUnicast:
		object, err := requireObject()
		if err != nil {
			return err
		}
		if err := rejectUnknown(object, "sender", "receiver", "subject", "body"); err != nil {
			return err
		}
		for _, field := range []string{"receiver", "subject"} {
			if err := requireString(object, field, true); err != nil {
				return err
			}
		}
		return nil
	case SendBroadcast:
		object, err := requireObject()
		if err != nil {
			return err
		}
		if err := rejectUnknown(object, "sender", "subject", "body"); err != nil {
			return err
		}
		return requireString(object, "subject", true)
	}
	return nil
}
