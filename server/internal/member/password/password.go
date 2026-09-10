package password

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const bcryptSHA256Prefix = "$bcrypt-sha256$"

// Hash creates a password hash suitable for storage. The password is
// pre-hashed so bcrypt's 72-byte input limit does not silently truncate long
// passwords.
func Hash(value string) (string, error) {
	digest := sha256.Sum256([]byte(value))
	hash, err := bcrypt.GenerateFromPassword(digest[:], bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return bcryptSHA256Prefix + string(hash), nil
}

// Verify accepts the current bcrypt format and legacy values. Legacy support
// allows existing installations to migrate without forcing a password reset.
func Verify(stored, value string, allowPlaintext bool) bool {
	digest := sha256.Sum256([]byte(value))

	switch {
	case strings.HasPrefix(stored, bcryptSHA256Prefix):
		return bcrypt.CompareHashAndPassword([]byte(strings.TrimPrefix(stored, bcryptSHA256Prefix)), digest[:]) == nil
	case strings.HasPrefix(stored, "$2a$") || strings.HasPrefix(stored, "$2b$") || strings.HasPrefix(stored, "$2y$"):
		// Accept hashes produced by older integrations that used bcrypt directly.
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(value)) == nil
	}

	legacy := sha256.Sum256([]byte(value))
	legacyEncoded := base64.StdEncoding.EncodeToString(legacy[:])
	if subtle.ConstantTimeCompare([]byte(stored), []byte(legacyEncoded)) == 1 {
		return true
	}

	return allowPlaintext && subtle.ConstantTimeCompare([]byte(stored), []byte(value)) == 1
}
