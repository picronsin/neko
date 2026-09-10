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

var (
	_ domain.RoomRepository    = (*SessionRoomRepository)(nil)
	_ domain.MemberRepository  = (*LiveMemberRepository)(nil)
	_ domain.SessionRepository = (*LiveSessionRepository)(nil)
)

func NewSessionRoomRepository(sessions types.SessionManager) *SessionRoomRepository {
	return &SessionRoomRepository{sessions: sessions}
}

func (r *SessionRoomRepository) Get(_ context.Context, roomID string) (domain.Room, error) {
	participants := make([]domain.Participant, 0)
	for _, current := range r.sessions.List() {
		participants = append(participants, domain.Participant{
			Member:  memberFromSession(current),
			Session: sessionFromLegacy(current),
		})
	}

	host, hasHost := r.sessions.GetHost()
	control := domain.ControlLease{Epoch: r.sessions.ControlEpoch()}
	if hasHost {
		control.HolderID = host.ID()
	}

	return domain.Room{ID: roomID, Participants: participants, Control: control}, nil
}

// LiveMemberRepository adapts member/profile reads from the current session
// manager. A database-backed member repository can replace it without
// changing domain or application service code.
type LiveMemberRepository struct {
	sessions types.SessionManager
}

func NewLiveMemberRepository(sessions types.SessionManager) *LiveMemberRepository {
	return &LiveMemberRepository{sessions: sessions}
}

func (r *LiveMemberRepository) Get(_ context.Context, memberID string) (domain.Member, error) {
	session, ok := r.sessions.Get(memberID)
	if !ok {
		return domain.Member{}, types.ErrSessionNotFound
	}
	return memberFromSession(session), nil
}

func (r *LiveMemberRepository) List(_ context.Context) ([]domain.Member, error) {
	members := make([]domain.Member, 0)
	for _, session := range r.sessions.List() {
		members = append(members, memberFromSession(session))
	}
	return members, nil
}

// LiveSessionRepository adapts live session state to the domain session port.
type LiveSessionRepository struct {
	sessions types.SessionManager
}

func NewLiveSessionRepository(sessions types.SessionManager) *LiveSessionRepository {
	return &LiveSessionRepository{sessions: sessions}
}

func (r *LiveSessionRepository) Get(_ context.Context, sessionID string) (domain.Session, error) {
	session, ok := r.sessions.Get(sessionID)
	if !ok {
		return domain.Session{}, types.ErrSessionNotFound
	}
	return sessionFromLegacy(session), nil
}

func (r *LiveSessionRepository) List(_ context.Context) ([]domain.Session, error) {
	result := make([]domain.Session, 0)
	for _, session := range r.sessions.List() {
		result = append(result, sessionFromLegacy(session))
	}
	return result, nil
}

func memberFromSession(session types.Session) domain.Member {
	profile := session.Profile()
	return domain.Member{
		ID:          session.ID(),
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
	}
}

func sessionFromLegacy(session types.Session) domain.Session {
	state := session.State()
	return domain.Session{ID: session.ID(), Connected: state.IsConnected, Watching: state.IsWatching}
}
