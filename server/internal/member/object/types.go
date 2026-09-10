package object

import (
	memberpassword "github.com/m1k1o/neko/server/internal/member/password"
	"github.com/m1k1o/neko/server/pkg/types"
)

type memberEntry struct {
	password string
	profile  types.MemberProfile
}

func (m *memberEntry) CheckPassword(password string) bool {
	return memberpassword.Verify(m.password, password, true)
}

type User struct {
	Username string
	Password string
	Profile  types.MemberProfile
}

type Config struct {
	Users []User
}
