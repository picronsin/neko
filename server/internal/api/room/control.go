package room

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/m1k1o/neko/server/internal/control"
	"github.com/m1k1o/neko/server/pkg/auth"
	"github.com/m1k1o/neko/server/pkg/protocol"
	"github.com/m1k1o/neko/server/pkg/utils"
)

type ControlStatusPayload struct {
	HasHost bool   `json:"has_host"`
	HostId  string `json:"host_id,omitempty"`
	Epoch   uint64 `json:"epoch"`
}

func (h *RoomHandler) controlStatus(w http.ResponseWriter, r *http.Request) error {
	status := h.control.Status()

	return utils.HttpSuccess(w, ControlStatusPayload{
		HasHost: status.HasHost,
		HostId:  status.HostID,
		Epoch:   status.Epoch,
	})
}

func (h *RoomHandler) controlRequest(w http.ResponseWriter, r *http.Request) error {
	session, _ := auth.GetSession(r)
	result, err := h.control.Request(session)
	if errors.Is(err, control.ErrNotAllowed) {
		return utils.HttpForbidden("controls are locked or unavailable").WithErrorCode(string(protocol.PermissionDenied))
	}
	if err != nil {
		return utils.HttpError(http.StatusConflict, err.Error()).WithErrorCode(string(protocol.ControlConflict))
	}
	if result.Granted {
		return utils.HttpSuccess(w)
	}
	if result.Queued {
		return utils.HttpError(http.StatusAccepted, "control request sent").WithErrorCode(string(protocol.ControlConflict))
	}

	return utils.HttpError(http.StatusConflict, "control lease is unavailable").WithErrorCode(string(protocol.ControlConflict))
}

func (h *RoomHandler) controlRelease(w http.ResponseWriter, r *http.Request) error {
	session, _ := auth.GetSession(r)
	if err := h.control.Release(session); err != nil {
		if errors.Is(err, control.ErrNotHost) {
			return utils.HttpUnprocessableEntity("session is not the host").WithErrorCode(string(protocol.PermissionDenied))
		}
		return utils.HttpUnprocessableEntity(err.Error()).WithErrorCode(string(protocol.PermissionDenied))
	}

	return utils.HttpSuccess(w)
}

func (h *RoomHandler) controlTake(w http.ResponseWriter, r *http.Request) error {
	session, _ := auth.GetSession(r)
	h.control.Take(session)

	return utils.HttpSuccess(w)
}

func (h *RoomHandler) controlGive(w http.ResponseWriter, r *http.Request) error {
	session, _ := auth.GetSession(r)
	sessionId := chi.URLParam(r, "sessionId")

	target, ok := h.sessions.Get(sessionId)
	if !ok {
		return utils.HttpNotFound("target session was not found").WithErrorCode(string(protocol.NotFound))
	}

	if err := h.control.Give(session, target); err != nil {
		if errors.Is(err, control.ErrNotAllowed) {
			return utils.HttpBadRequest("target session is not allowed to host").WithErrorCode(string(protocol.PermissionDenied))
		}
		return utils.HttpUnprocessableEntity(err.Error()).WithErrorCode(string(protocol.PermissionDenied))
	}

	return utils.HttpSuccess(w)
}

func (h *RoomHandler) controlReset(w http.ResponseWriter, r *http.Request) error {
	session, _ := auth.GetSession(r)
	h.control.Reset(session)

	return utils.HttpSuccess(w)
}
