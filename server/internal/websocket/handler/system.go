package handler

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/m1k1o/neko/server/pkg/types"
	"github.com/m1k1o/neko/server/pkg/types/event"
	"github.com/m1k1o/neko/server/pkg/types/message"
)

func (h *MessageHandlerCtx) systemInit(session types.Session) error {
	snapshot := h.room.Snapshot(session)

	sessions := map[string]message.SessionData{}
	for _, current := range snapshot.Room.Participants {
		sessions[current.Session.ID] = message.SessionData{
			ID: current.Session.ID,
			Profile: types.MemberProfile{
				Name:                  current.Member.DisplayName,
				Avatar:                current.Member.Avatar,
				IsAdmin:               current.Member.Permission.Admin,
				CanLogin:              current.Member.Permission.Login,
				CanConnect:            current.Member.Permission.Connect,
				CanWatch:              current.Member.Permission.Watch,
				CanHost:               current.Member.Permission.Host,
				CanShareMedia:         current.Member.Permission.ShareMedia,
				CanAccessClipboard:    current.Member.Permission.AccessClipboard,
				SendsInactiveCursor:   current.Member.Permission.SendInactiveCursor,
				CanSeeInactiveCursors: current.Member.Permission.SeeInactiveCursors,
			},
			State: types.SessionState{
				IsConnected: current.Session.Connected,
				IsWatching:  current.Session.Watching,
			},
		}
	}

	session.Send(
		event.SYSTEM_INIT,
		message.SystemInit{
			SessionId: snapshot.SessionID,
			ControlHost: message.ControlHost{
				HasHost: snapshot.Room.Control.HolderID != "",
				HostID:  snapshot.Room.Control.HolderID,
				Epoch:   snapshot.Room.Control.Epoch,
			},
			ScreenSize:        snapshot.ScreenSize,
			Sessions:          sessions,
			Settings:          snapshot.Settings,
			TouchEvents:       snapshot.TouchEvents,
			ScreencastEnabled: snapshot.ScreencastEnabled,
			WebRTC: message.SystemWebRTC{
				Videos: snapshot.VideoIDs,
			},
		})

	return nil
}

func (h *MessageHandlerCtx) systemAdmin(session types.Session) error {
	session.Send(event.SYSTEM_ADMIN, message.SystemAdmin{})

	return nil
}

func (h *MessageHandlerCtx) systemLogs(session types.Session, payload *message.SystemLogs) error {
	for _, msg := range *payload {
		level, _ := zerolog.ParseLevel(msg.Level)

		if level < zerolog.DebugLevel || level > zerolog.ErrorLevel {
			level = zerolog.NoLevel
		}

		// do not use handler logger context
		log.WithLevel(level).
			Fields(msg.Fields).
			Str("module", "client").
			Str("session_id", session.ID()).
			Msg(msg.Message)
	}

	return nil
}
