// Package domain contains transport-neutral room concepts. Adapters convert
// these values to the existing session, REST, WebSocket, and media types.
package domain

// Permission is the room-facing permission set. It deliberately contains no
// authentication provider, HTTP, WebSocket, or runtime manager types.
type Permission struct {
	Admin              bool
	Login              bool
	Connect            bool
	Watch              bool
	Host               bool
	ShareMedia         bool
	AccessClipboard    bool
	SendInactiveCursor bool
	SeeInactiveCursors bool
}

type Member struct {
	ID          string
	DisplayName string
	Avatar      string
	Permission  Permission
}

type Session struct {
	ID        string
	Connected bool
	Watching  bool
}

type Participant struct {
	Member  Member
	Session Session
}

type ControlLease struct {
	HolderID string
	Epoch    uint64
}

type Room struct {
	ID           string
	Participants []Participant
	Control      ControlLease
}

// CanRequestControl is the domain permission rule shared by adapters. The
// caller remains responsible for queueing and lease mutation.
func CanRequestControl(member Member, locked bool, privateMode bool, alreadyHolder bool) bool {
	if !member.Permission.Host || privateMode || alreadyHolder {
		return false
	}
	return !locked || member.Permission.Admin
}
