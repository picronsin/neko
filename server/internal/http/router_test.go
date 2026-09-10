package http

import (
	"net/http"
	"testing"

	"github.com/m1k1o/neko/server/pkg/protocol"
)

func TestDefaultErrorCode(t *testing.T) {
	tests := []struct {
		status int
		code   string
	}{
		{http.StatusBadRequest, string(protocol.InvalidPayload)},
		{http.StatusForbidden, string(protocol.PermissionDenied)},
		{http.StatusNotFound, string(protocol.NotFound)},
		{http.StatusConflict, string(protocol.ControlConflict)},
		{http.StatusInternalServerError, string(protocol.InternalError)},
	}
	for _, test := range tests {
		if got := defaultErrorCode(test.status); got != test.code {
			t.Fatalf("defaultErrorCode(%d) = %q, want %q", test.status, got, test.code)
		}
	}
}
