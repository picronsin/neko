package room

import (
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
	return &Service{sessions: sessions, desktop: desktop, capture: capture, control: controlService}
}

func (s *Service) Snapshot(session types.Session) Snapshot {
	status := s.control.Status()

	participants := make([]domain.Participant, 0)
	for _, current := range s.sessions.List() {
		profile := current.Profile()
		state := current.State()
		participants = append(participants, domain.Participant{
			Member: domain.Member{
				ID:          current.ID(),
				DisplayName: profile.Name,
				Avatar:      profile.Avatar,
				Permission: domain.Permission{
					Admin:              profile.IsAdmin,
					Login:              profile.CanLogin,
					Connect:            profile.CanConnect,
					Watch:              profile.CanWatch,
					Host:               profile.CanHost,
					ShareMedia:         profile.CanShareMedia,
					AccessClipboard:    profile.CanAccessClipboard,
					SendInactiveCursor: profile.SendsInactiveCursor,
					SeeInactiveCursors: profile.CanSeeInactiveCursors,
				},
			},
			Session: domain.Session{
				ID:        current.ID(),
				Connected: state.IsConnected,
				Watching:  state.IsWatching,
			},
		})
	}

	return Snapshot{
		SessionID: session.ID(),
		Room: domain.Room{
			Participants: participants,
			Control:      domain.ControlLease{HolderID: status.HostID, Epoch: status.Epoch},
		},
		ScreenSize:        s.desktop.GetScreenSize(),
		Settings:          s.sessions.Settings(),
		TouchEvents:       s.desktop.HasTouchSupport(),
		ScreencastEnabled: s.capture.Screencast().Enabled(),
		VideoIDs:          s.capture.Video().IDs(),
	}
}
