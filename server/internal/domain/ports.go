package domain

import "context"

// RoomRepository is the application-facing read boundary for room state.
// Implementations may read an in-memory session manager today or a persistent
// room store later without leaking that choice into domain code.
type RoomRepository interface {
	Get(ctx context.Context, roomID string) (Room, error)
}

// MemberRepository is the persistence boundary for room members. The current
// deployment still uses the legacy member providers; keeping this port in the
// domain package makes the migration target explicit for M3 adapters.
type MemberRepository interface {
	Get(ctx context.Context, memberID string) (Member, error)
	List(ctx context.Context) ([]Member, error)
}

// SessionRepository is the live-session boundary used by application services.
// It intentionally exposes domain values rather than transport/session
// manager interfaces.
type SessionRepository interface {
	Get(ctx context.Context, sessionID string) (Session, error)
	List(ctx context.Context) ([]Session, error)
}
