package handler

import (
	"testing"

	"github.com/m1k1o/neko/server/internal/control"
	"github.com/m1k1o/neko/server/pkg/protocol"
)

func TestProtocolErrorMapsControlFailures(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code protocol.ErrorCode
	}{
		{name: "conflict", err: control.ErrConflict, code: protocol.ControlConflict},
		{name: "stale epoch", err: control.ErrStaleEpoch, code: protocol.StaleEpoch},
		{name: "not allowed", err: control.ErrNotAllowed, code: protocol.PermissionDenied},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := protocolError(test.err).Code; got != test.code {
				t.Fatalf("protocolError(%v) = %q, want %q", test.err, got, test.code)
			}
		})
	}
}
