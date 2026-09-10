package room

import (
	"context"

	"github.com/m1k1o/neko/server/internal/domain"
	"github.com/m1k1o/neko/server/pkg/types"
)

// SessionRoomRepository adapts the current live session manager to the domain
// room read port. It is deliberately kept in the application adapter package;
// the domain package does not know about authentication or transport types.
type SessionRoomRepository struct {
	sessions types.SessionManager
}

func NewSessionRoomRepository(sessions types.SessionManager) *SessionRoomRepository {
	return &SessionRoomRepository{sessions: sessions}
}

func (r *SessionRoomRepository) Get(_ context.Context, roomID string) (domain.Room, error) {
	participants := make([]domain.Participant, 0)
	for _, current := range r.sessions.List() {
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

	host, hasHost := r.sessions.GetHost()
	control := domain.ControlLease{Epoch: r.sessions.ControlEpoch()}
	if hasHost {
		control.HolderID = host.ID()
	}

	return domain.Room{ID: roomID, Participants: participants, Control: control}, nil
}
