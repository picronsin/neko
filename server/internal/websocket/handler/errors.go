package handler

import (
	"errors"
	"strings"

	desktopapp "github.com/m1k1o/neko/server/internal/application/desktop"
	"github.com/m1k1o/neko/server/internal/control"
	"github.com/m1k1o/neko/server/pkg/protocol"
)

func protocolError(err error) protocol.ErrorPayload {
	code := protocol.InternalError
	switch {
	case err == nil:
		code = protocol.InternalError
	case errors.Is(err, control.ErrConflict), errors.Is(err, control.ErrUnavailable), errors.Is(err, ErrIsAlreadyHosted):
		code = protocol.ControlConflict
	case errors.Is(err, control.ErrStaleEpoch):
		code = protocol.StaleEpoch
	case errors.Is(err, control.ErrNotHolder), errors.Is(err, control.ErrNotHost), errors.Is(err, control.ErrNotAllowed),
		errors.Is(err, ErrIsNotAllowedToHost), errors.Is(err, ErrIsNotTheHost),
		errors.Is(err, desktopapp.ErrNotHost), errors.Is(err, desktopapp.ErrClipboardForbidden):
		code = protocol.PermissionDenied
	default:
		message := strings.ToLower(err.Error())
		switch {
		case strings.Contains(message, "unmarshal"), strings.Contains(message, "invalid character"), strings.Contains(message, "payload"):
			code = protocol.InvalidPayload
		case strings.Contains(message, "ice"), strings.Contains(message, "webrtc"):
			code = protocol.IceFailed
		case strings.Contains(message, "peer does not exist"), strings.Contains(message, "not ready"):
			code = protocol.RoomNotReady
		case strings.Contains(message, "not found"):
			code = protocol.NotFound
		}
	}
	return protocol.NewError(code, err.Error())
}
