package util

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

// RandBytes returns n cryptographically secure random bytes.
func RandBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

// RandID returns a 32-char hex string (128-bit) suitable as an ID.
func RandID() string {
	return hex.EncodeToString(RandBytes(16))
}

// RandToken returns a 64-char hex string (256-bit) suitable for tokens.
func RandToken() string {
	return hex.EncodeToString(RandBytes(32))
}

// HashToken returns hex-encoded SHA-256 of the token string.
func HashToken(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// ConstantTimeEquals compares two strings in constant time.
func ConstantTimeEquals(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(strings.ToLower(a)), []byte(strings.ToLower(b))) == 1
}
