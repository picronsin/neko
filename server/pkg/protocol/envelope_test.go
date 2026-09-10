package protocol

import (
	"encoding/json"
	"testing"
)

func TestEnvelopeRejectsScalarPayload(t *testing.T) {
	var message Envelope
	if err := json.Unmarshal([]byte(`{"event":"chat/message","payload":"flat"}`), &message); err == nil {
		t.Fatal("expected scalar payload to be rejected")
	}
}

func TestEnvelopeAcceptsGeneratedEventCatalog(t *testing.T) {
	if !IsKnownEvent(SystemInit) || !IsKnownEvent(ChatMessage) {
		t.Fatal("generated event catalog is incomplete")
	}
	if err := ValidateEvent("missing/event"); err == nil {
		t.Fatal("expected unknown event to be rejected by catalog validation")
	}
}

func TestErrorPayloadUsesStableCode(t *testing.T) {
	payload := NewError(ControlConflict, "control is already held")
	if payload.Code != ControlConflict || !payload.Retryable {
		t.Fatalf("unexpected error payload: %+v", payload)
	}
	if got := string(payload.Code); got != "CONTROL_CONFLICT" {
		t.Fatalf("unexpected error code %q", got)
	}
}
