package crypto

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

const (
	IterCount = 100000
	KeyLen    = 32
)

func GenerateKeys(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, IterCount, KeyLen, sha256.New)
}
