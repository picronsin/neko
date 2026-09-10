package protocol

import (
	"encoding/json"
	"testing"
)

func TestValidatePayload(t *testing.T) {
	tests := []struct {
		name    string
		event   string
		payload string
		valid   bool
	}{
		{name: "valid offer", event: SignalOffer, payload: `{"sdp":"v=0"}`, valid: true},
		{name: "missing offer sdp", event: SignalOffer, payload: `{}`, valid: false},
		{name: "valid renew", event: ControlRenew, payload: `{"epoch":4}`, valid: true},
		{name: "missing renew epoch", event: ControlRenew, payload: `{}`, valid: false},
		{name: "valid modifiers", event: KeyboardModifiers, payload: `{}`, valid: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidatePayload(test.event, json.RawMessage(test.payload))
			if (err == nil) != test.valid {
				t.Fatalf("ValidatePayload(%s, %s) = %v, valid=%v", test.event, test.payload, err, test.valid)
			}
		})
	}
}
