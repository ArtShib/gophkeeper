package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"

	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"golang.org/x/crypto/pbkdf2"
)

const (
	IterCount = 100000
	KeyLen    = 32
)

type CryptoService struct {
	key    []byte
	logger *slog.Logger
}

func NewCryptoService(pass string, salt []byte, log *slog.Logger) *CryptoService {
	return &CryptoService{
		key:    pbkdf2.Key([]byte(pass), salt, IterCount, KeyLen, sha256.New),
		logger: log,
	}
}

func (c *CryptoService) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	log := loghelper.New(c.logger, "crypto.Encrypt")
	aesblock, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "aes.NewCipher", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "aes.NewGCM", err)
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, log.LogAndReturnError(ctx, "rand nonce", err)
	}
	dst := aesgcm.Seal(nonce, nonce, data, nil)

	return dst, nil
}

func (c *CryptoService) Decrypt(ctx context.Context, encryptedData []byte) ([]byte, error) {
	log := loghelper.New(c.logger, "crypto.Decrypt")
	aesblock, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "aes.NewCipher", err)
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "aes.NewGCM", err)
	}

	if len(encryptedData) < aesgcm.NonceSize()+aesgcm.Overhead() {
		return nil, log.LogAndReturnError(ctx, "invalid encryptedData", fmt.Errorf("encrypted data too short"))
	}

	nonce, encrypted := encryptedData[:aesgcm.NonceSize()], encryptedData[aesgcm.NonceSize():]
	decrypted, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "aesgcm.Open", err)
	}
	return decrypted, nil
}
