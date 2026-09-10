package room

import (
	"context"

	"github.com/m1k1o/neko/server/internal/control"
	"github.com/m1k1o/neko/server/internal/domain"
	"github.com/m1k1o/neko/server/pkg/types"
)

// Snapshot is the transport-neutral room state needed when a client joins.
type Snapshot struct {
	SessionID         string
	Room              domain.Room
	ScreenSize        types.ScreenSize
	Settings          types.Settings
	TouchEvents       bool
	ScreencastEnabled bool
	VideoIDs          []string
}

type Service struct {
	rooms    domain.RoomRepository
	sessions types.SessionManager
	desktop  types.DesktopManager
	capture  types.CaptureManager
	control  *control.Service
}

func NewService(
	sessions types.SessionManager,
	desktop types.DesktopManager,
	capture types.CaptureManager,
	controlService *control.Service,
) *Service {
	return &Service{
		rooms:    NewSessionRoomRepository(sessions),
		sessions: sessions,
		desktop:  desktop,
		capture:  capture,
		control:  controlService,
	}
}

func (s *Service) Snapshot(session types.Session) Snapshot {
	room, err := s.rooms.Get(context.Background(), "default")
	if err != nil {
		room = domain.Room{ID: "default"}
	}

	return Snapshot{
		SessionID:         session.ID(),
		Room:              room,
		ScreenSize:        s.desktop.GetScreenSize(),
		Settings:          s.sessions.Settings(),
		TouchEvents:       s.desktop.HasTouchSupport(),
		ScreencastEnabled: s.capture.Screencast().Enabled(),
		VideoIDs:          s.capture.Video().IDs(),
	}
}
