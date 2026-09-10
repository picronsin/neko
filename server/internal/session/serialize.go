package session

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/m1k1o/neko/server/pkg/types"
)

func (manager *SessionManagerCtx) save() {
	if manager.config.File == "" {
		return
	}
	manager.persistMu.Lock()
	defer manager.persistMu.Unlock()

	// Snapshot sessions while holding the map lock. save is called after the
	// mutating operation releases sessionsMu, so it must protect this read
	// itself from concurrent Create/Update/Delete calls.
	manager.sessionsMu.Lock()
	sessions := make([]types.SessionProfile, 0, len(manager.sessions))
	for _, session := range manager.sessions {
		sessions = append(sessions, types.SessionProfile{
			Id:      session.id,
			Token:   session.token,
			Profile: session.profile,
		})
	}
	manager.sessionsMu.Unlock()

	// convert to json
	data, err := json.Marshal(sessions)
	if err != nil {
		manager.logger.Error().Err(err).Msg("failed to marshal sessions")
		return
	}

	if err := os.MkdirAll(filepath.Dir(manager.config.File), 0750); err != nil {
		manager.logger.Error().Err(err).
			Str("file", manager.config.File).
			Msg("failed to create session directory")
		return
	}

	// Write to a private temporary file and replace the target atomically. This
	// keeps authentication tokens private and avoids leaving a partially-written
	// session file after a crash.
	temporary := manager.config.File + ".tmp"
	err = os.WriteFile(temporary, data, 0600)
	if err != nil {
		manager.logger.Error().Err(err).
			Str("file", manager.config.File).
			Msg("failed to write sessions to a temporary file")
		return
	}
	if err := os.Rename(temporary, manager.config.File); err != nil {
		manager.logger.Error().Err(err).
			Str("file", manager.config.File).
			Msg("failed to replace sessions file")
	}
}

func (manager *SessionManagerCtx) load() {
	if manager.config.File == "" {
		return
	}

	// read file
	data, err := os.ReadFile(manager.config.File)
	if err != nil {
		// if file does not exist
		if errors.Is(err, os.ErrNotExist) {
			manager.logger.Info().
				Str("file", manager.config.File).
				Msg("sessions file does not exist")
			return
		}
		manager.logger.Error().Err(err).
			Str("file", manager.config.File).
			Msg("failed to read sessions from a file")
		return
	}
	if err := os.Chmod(manager.config.File, 0600); err != nil {
		manager.logger.Warn().Err(err).
			Str("file", manager.config.File).
			Msg("failed to restrict sessions file permissions")
	}

	// if file is empty
	if len(data) == 0 {
		manager.logger.Info().
			Str("file", manager.config.File).
			Msg("sessions file is empty")
		return
	}

	// deserialize sessions
	sessions := make([]types.SessionProfile, 0)
	err = json.Unmarshal(data, &sessions)
	if err != nil {
		manager.logger.Error().Err(err).Msg("failed to unmarshal sessions")
		return
	}

	// create sessions
	manager.sessionsMu.Lock()
	for _, session := range sessions {
		profile := manager.withStoredAvatar(session.Profile, session.Id)
		manager.tokens[session.Token] = session.Id
		manager.sessions[session.Id] = &SessionCtx{
			id:      session.Id,
			token:   session.Token,
			manager: manager,
			logger:  manager.logger.With().Str("session_id", session.Id).Logger(),
			profile: profile,
		}
	}
	manager.sessionsMu.Unlock()

	manager.logger.Info().
		Int("sessions", len(sessions)).
		Str("file", manager.config.File).
		Msg("loaded sessions from a file")
}
