package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateUUID генерация uuid
func GenerateUUID() (string, error) {
	lenUUID := 8
	bytes := make([]byte, lenUUID)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:lenUUID], nil
}
