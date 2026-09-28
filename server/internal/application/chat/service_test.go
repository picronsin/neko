package chat

import (
	"path/filepath"
	"testing"
)

func TestHistorySurvivesRestartAndRespectsLimit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "history.json")
	service := NewService(nil, true, path, 2)
	for _, text := range []string{"first", "second", "third"} {
		service.append(Message{ID: "member", Name: "Alice", Content: Content{Text: text}})
	}
	restored := NewService(nil, true, path, 2).messages()
	if len(restored) != 2 || restored[0].Content.Text != "second" || restored[1].Content.Text != "third" || restored[1].Name != "Alice" {
		t.Fatalf("unexpected persisted history: %#v", restored)
	}
	trimmed := NewService(nil, true, path, 1).messages()
	if len(trimmed) != 1 || trimmed[0].Content.Text != "third" {
		t.Fatalf("limit on startup was not applied: %#v", trimmed)
	}
}
