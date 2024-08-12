package password_hash

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hasher interface {
	GetHashPassword(password string) (string, error)
}

type SHA256Hasher struct{}

func (h *SHA256Hasher) GetHashPassword(password string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}
	hashBytes := hash.Sum(nil)
	hashPass := hex.EncodeToString(hashBytes)
	return hashPass, nil
}
