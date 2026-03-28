package models

import (
	"crypto/rand"
	"encoding/hex"
)

// NewID generates a random 16-byte hex string (similar length to UUID) without external dependencies
func NewID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
