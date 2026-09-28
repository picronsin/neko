package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestUsesLegacyGlobalICEServers(t *testing.T) {
	tests := []struct {
		name   string
		value  any
		legacy bool
	}{
		{name: "unset", value: nil},
		{name: "separate endpoint servers", value: map[string]any{"frontend": []any{}, "backend": []any{}}},
		{name: "legacy list", value: []any{map[string]any{"urls": []any{"stun:example.test"}}}, legacy: true},
		{name: "legacy JSON environment value", value: `[{"urls":["stun:example.test"]}]`, legacy: true},
		{name: "unknown nested key", value: map[string]any{"public": []any{}}, legacy: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := usesLegacyGlobalICEServers(test.value); got != test.legacy {
				t.Fatalf("usesLegacyGlobalICEServers() = %v, want %v", got, test.legacy)
			}
		})
	}
}

func TestValidateUnsupportedLegacyConfig(t *testing.T) {
	t.Cleanup(viper.Reset)

	viper.Set("webrtc.epr", "52000-52100")
	if err := validateUnsupportedLegacyConfig(); err == nil {
		t.Fatal("validateUnsupportedLegacyConfig() accepted webrtc.epr")
	}

	viper.Reset()
	viper.Set("webrtc.iceservers", []any{map[string]any{"urls": []any{"stun:example.test"}}})
	if err := validateUnsupportedLegacyConfig(); err == nil {
		t.Fatal("validateUnsupportedLegacyConfig() accepted global webrtc.iceservers")
	}

	viper.Reset()
	viper.Set("webrtc.iceservers.frontend", []any{})
	viper.Set("webrtc.iceservers.backend", []any{})
	if err := validateUnsupportedLegacyConfig(); err != nil {
		t.Fatalf("validateUnsupportedLegacyConfig() rejected endpoint ICE servers: %v", err)
	}
}
