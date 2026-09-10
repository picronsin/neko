package domain

import "testing"

func TestCanRequestControl(t *testing.T) {
	member := Member{Permission: Permission{Host: true}}
	if !CanRequestControl(member, false, false, false) {
		t.Fatal("host-capable member should be allowed")
	}
	if CanRequestControl(member, true, false, false) {
		t.Fatal("locked room should reject non-admin member")
	}
	member.Permission.Admin = true
	if !CanRequestControl(member, true, false, false) {
		t.Fatal("admin should bypass locked controls")
	}
	if CanRequestControl(member, false, true, false) || CanRequestControl(member, false, false, true) {
		t.Fatal("private mode and existing holder should reject control requests")
	}
}
