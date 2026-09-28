package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/m1k1o/neko/server/pkg/types"
)

type Settings struct {
	CanSend    bool `mapstructure:"can_send"`
	CanReceive bool `mapstructure:"can_receive"`
}

type Content struct {
	Text string `json:"text"`
}

type Message struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	Content Content   `json:"content"`
	Name    string    `json:"name,omitempty"`
	Avatar  string    `json:"avatar,omitempty"`
}

type History struct {
	Messages []Message `json:"messages"`
}

type Service struct {
	sessions types.SessionManager
	enabled  bool
	file     string
	limit    int
	mu       sync.RWMutex
	history  []Message
}

func NewService(sessions types.SessionManager, enabled bool, historyFile string, historyLimit int) *Service {
	service := &Service{
		sessions: sessions,
		enabled:  enabled,
		file:     historyFile,
		limit:    historyLimit,
		history:  make([]Message, 0),
	}
	service.load()
	return service
}

func (s *Service) Settings(session types.Session) (Settings, error) {
	global := Settings{CanSend: true, CanReceive: true}
	if err := s.sessions.Settings().Plugins.Unmarshal("chat", &global); err != nil && !errors.Is(err, types.ErrPluginSettingsNotFound) {
		return Settings{}, fmt.Errorf("unable to unmarshal chat settings: %w", err)
	}
	profile := Settings{CanSend: true, CanReceive: true}
	if err := session.Profile().Plugins.Unmarshal("chat", &profile); err != nil && !errors.Is(err, types.ErrPluginSettingsNotFound) {
		return Settings{}, fmt.Errorf("unable to unmarshal chat profile settings: %w", err)
	}
	return Settings{
		CanSend:    s.enabled && (global.CanSend || session.Profile().IsAdmin) && profile.CanSend,
		CanReceive: s.enabled && (global.CanReceive || session.Profile().IsAdmin) && profile.CanReceive,
	}, nil
}

func (s *Service) Send(session types.Session, content Content) error {
	settings, err := s.Settings(session)
	if err != nil {
		return err
	}
	if !settings.CanSend {
		return errors.New("not allowed to send chat messages")
	}

	now := time.Now()
	message := Message{
		ID:      session.ID(),
		Created: now,
		Content: content,
		Name:    session.Profile().Name,
		Avatar:  session.Profile().Avatar,
	}
	s.append(message)
	s.sessions.Range(func(target types.Session) bool {
		settings, err := s.Settings(target)
		if err == nil && settings.CanReceive {
			target.Send("chat/message", message)
		}
		return true
	})
	return nil
}

func (s *Service) Initialize(session types.Session) {
	settings, err := s.Settings(session)
	var history []Message
	if err == nil && settings.CanReceive {
		history = s.messages()
	}
	session.Send("chat/init", struct {
		Enabled bool      `json:"enabled"`
		History []Message `json:"history,omitempty"`
	}{Enabled: s.enabled, History: history})
}

func (s *Service) append(message Message) {
	s.mu.Lock()
	s.history = append(s.history, message)
	if s.limit > 0 && len(s.history) > s.limit {
		s.history = append([]Message(nil), s.history[len(s.history)-s.limit:]...)
	}
	history := History{Messages: append([]Message(nil), s.history...)}
	// Keep the history lock while persisting the snapshot. Otherwise concurrent
	// sends can save snapshots out of order and an older snapshot can overwrite
	// a newer one on disk.
	s.save(history)
	s.mu.Unlock()
}

func (s *Service) messages() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Message(nil), s.history...)
}

func (s *Service) load() {
	if s.file == "" {
		return
	}

	data, err := os.ReadFile(s.file)
	if errors.Is(err, os.ErrNotExist) {
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("unable to read chat history")
		return
	}

	var history History
	if err := json.Unmarshal(data, &history); err != nil {
		log.Error().Err(err).Msg("unable to decode chat history")
		return
	}
	if s.limit > 0 && len(history.Messages) > s.limit {
		history.Messages = history.Messages[len(history.Messages)-s.limit:]
	}
	s.history = append([]Message(nil), history.Messages...)
}

func (s *Service) save(history History) {
	if s.file == "" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(s.file), 0750); err != nil {
		log.Error().Err(err).Msg("unable to create chat history directory")
		return
	}
	data, err := json.Marshal(history)
	if err != nil {
		return
	}

	temporary := s.file + ".tmp"
	if err := os.WriteFile(temporary, data, 0600); err != nil {
		log.Error().Err(err).Msg("unable to write chat history")
		return
	}
	if err := os.Rename(temporary, s.file); err != nil {
		log.Error().Err(err).Msg("unable to replace chat history")
	}
}
