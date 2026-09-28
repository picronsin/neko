package chat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appchat "github.com/m1k1o/neko/server/internal/application/chat"
	"github.com/m1k1o/neko/server/pkg/auth"
	"github.com/m1k1o/neko/server/pkg/types"
)

type voiceTestSession struct {
	types.Session
	id     string
	denied bool
}

func (s voiceTestSession) ID() string                { return s.id }
func (s voiceTestSession) State() types.SessionState { return types.SessionState{IsConnected: true} }
func (s voiceTestSession) Profile() types.MemberProfile {
	return types.MemberProfile{Name: s.id, Plugins: types.PluginSettings{"chat.can_send": !s.denied, "chat.can_receive": true}}
}

type voiceTestSessions struct{ types.SessionManager }

func (voiceTestSessions) Settings() types.Settings { return types.Settings{} }
func voiceRequest(id, body string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	return r.WithContext(auth.SetSession(r, voiceTestSession{id: id}))
}
func TestVoiceMembershipAndSignaling(t *testing.T) {
	m := &Manager{service: appchat.NewService(voiceTestSessions{}, true, "", 200)}
	if err := m.voiceJoin(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", nil)); err == nil {
		t.Fatal("anonymous join accepted")
	}
	denied := voiceRequest("denied", "")
	denied = denied.WithContext(auth.SetSession(denied, voiceTestSession{id: "denied", denied: true}))
	if err := m.voiceJoin(httptest.NewRecorder(), denied); err == nil {
		t.Fatal("muted member joined")
	}
	for i := 0; i < 6; i++ {
		if err := m.voiceJoin(httptest.NewRecorder(), voiceRequest(fmt.Sprint(i), "")); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.voiceJoin(httptest.NewRecorder(), voiceRequest("7", "")); err == nil {
		t.Fatal("capacity exceeded")
	}
	message := `{"from":"spoofed","to":"1","type":"offer","sdp":"test"}`
	if err := m.voiceSignal(httptest.NewRecorder(), voiceRequest("outsider", message)); err == nil {
		t.Fatal("nonmember signaled")
	}
	if err := m.voiceSignal(httptest.NewRecorder(), voiceRequest("0", message)); err != nil {
		t.Fatal(err)
	}
	if m.voice.members["1"].signals[0].From != "0" {
		t.Fatal("sender identity spoofed")
	}
	for i := 1; i < 16; i++ {
		if err := m.voiceSignal(httptest.NewRecorder(), voiceRequest("0", message)); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.voiceSignal(httptest.NewRecorder(), voiceRequest("0", message)); err == nil {
		t.Fatal("queue limit exceeded")
	}
	if err := m.voicePoll(httptest.NewRecorder(), voiceRequest("1", "")); err != nil {
		t.Fatal(err)
	}
	if len(m.voice.members["1"].signals) != 0 {
		t.Fatal("queue not drained")
	}
	m.voice.members["2"].seen = time.Now().Add(-time.Minute)
	m.voice.expire()
	if m.voice.members["2"] != nil {
		t.Fatal("expired session retained")
	}
	m.voice.remove("0")
	if m.voice.members["0"] != nil {
		t.Fatal("departed member retained")
	}
}
