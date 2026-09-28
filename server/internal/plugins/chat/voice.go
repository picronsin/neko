package chat

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/m1k1o/neko/server/pkg/auth"
	"github.com/m1k1o/neko/server/pkg/types"
	"github.com/m1k1o/neko/server/pkg/utils"
	"github.com/spf13/viper"
)

// Voice uses independent browser-to-browser audio connections. Only bounded
// signaling descriptions pass through this mailbox; audio is never recorded.
type voiceMember struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Muted   bool   `json:"muted"`
	seen    time.Time
	signals []voiceSignal
}
type voiceSignal struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}
type voiceRoom struct {
	sync.Mutex
	members map[string]*voiceMember
}

func (v *voiceRoom) remove(id string) {
	v.Lock()
	defer v.Unlock()
	delete(v.members, id)
}
func (m *Manager) voiceRoutes(r types.Router) {
	r.Post("/voice", m.voiceJoin)
	r.Get("/voice", m.voicePoll)
	r.Delete("/voice", m.voiceLeave)
	r.Post("/voice/signal", m.voiceSignal)
}
func (m *Manager) voiceSession(r *http.Request) (types.Session, error) {
	s, ok := auth.GetSession(r)
	if !ok || !s.State().IsConnected {
		return nil, utils.HttpUnauthorized("join the room first")
	}
	settings, err := m.settingsForSession(s)
	if err != nil || !settings.CanSend || !settings.CanReceive {
		return nil, utils.HttpForbidden("voice chat is not allowed")
	}
	return s, nil
}
func (v *voiceRoom) expire() {
	for id, member := range v.members {
		if time.Since(member.seen) > 15*time.Second {
			delete(v.members, id)
		}
	}
}
func (m *Manager) voiceJoin(w http.ResponseWriter, r *http.Request) error {
	s, err := m.voiceSession(r)
	if err != nil {
		return err
	}
	m.voice.Lock()
	defer m.voice.Unlock()
	m.voice.expire()
	if m.voice.members == nil {
		m.voice.members = make(map[string]*voiceMember)
	}
	if _, exists := m.voice.members[s.ID()]; exists {
		return utils.HttpError(http.StatusConflict, "already in voice chat")
	}
	if len(m.voice.members) >= 6 {
		return utils.HttpError(http.StatusConflict, "voice chat is full (6 participants)")
	}
	var servers []types.ICEServer
	if err := viper.UnmarshalKey("webrtc.iceservers.frontend", &servers, viper.DecodeHook(utils.JsonStringAutoDecode(servers))); err != nil {
		return utils.HttpInternalServerError().Msg("invalid frontend ICE servers")
	}
	m.voice.members[s.ID()] = &voiceMember{ID: s.ID(), Name: s.Profile().Name, seen: time.Now()}
	return utils.HttpSuccess(w, struct {
		ID      string            `json:"id"`
		Servers []types.ICEServer `json:"servers"`
	}{s.ID(), servers})
}
func (m *Manager) voicePoll(w http.ResponseWriter, r *http.Request) error {
	s, err := m.voiceSession(r)
	if err != nil {
		return err
	}
	m.voice.Lock()
	defer m.voice.Unlock()
	m.voice.expire()
	member := m.voice.members[s.ID()]
	if member == nil {
		return utils.HttpError(http.StatusConflict, "voice session expired; rejoin")
	}
	member.seen = time.Now()
	member.Muted = r.URL.Query().Get("muted") == "true"
	members := make([]voiceMember, 0, len(m.voice.members))
	for _, peer := range m.voice.members {
		members = append(members, *peer)
	}
	signals := member.signals
	member.signals = nil
	return utils.HttpSuccess(w, struct {
		Members []voiceMember `json:"members"`
		Signals []voiceSignal `json:"signals"`
	}{members, signals})
}
func (m *Manager) voiceLeave(w http.ResponseWriter, r *http.Request) error {
	s, ok := auth.GetSession(r)
	if !ok {
		return utils.HttpUnauthorized("session not found")
	}
	m.voice.remove(s.ID())
	return utils.HttpSuccess(w)
}
func (m *Manager) voiceSignal(w http.ResponseWriter, r *http.Request) error {
	s, err := m.voiceSession(r)
	if err != nil {
		return err
	}
	var signal voiceSignal
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536)).Decode(&signal); err != nil {
		return utils.HttpBadRequest("invalid voice signal")
	}
	if (signal.Type != "offer" && signal.Type != "answer") || len(signal.SDP) == 0 || signal.To == s.ID() {
		return utils.HttpBadRequest("invalid voice description")
	}
	m.voice.Lock()
	defer m.voice.Unlock()
	m.voice.expire()
	if m.voice.members[s.ID()] == nil {
		return utils.HttpForbidden("join voice chat first")
	}
	target := m.voice.members[signal.To]
	if target == nil {
		return utils.HttpNotFound("voice participant left")
	}
	if len(target.signals) >= 16 {
		return utils.HttpError(http.StatusTooManyRequests, "voice signaling queue full")
	}
	signal.From = s.ID()
	target.signals = append(target.signals, signal)
	return utils.HttpSuccess(w)
}
