package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// generateToken は十分なエントロピーを持つセッショントークン（生）を返す。
// この生トークンだけが Cookie に載り、DBには hashToken した値のみを保存する。
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("auth: generate token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// hashToken は生トークンを SHA-256 でハッシュ化する（DB保存・照合用）。
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
