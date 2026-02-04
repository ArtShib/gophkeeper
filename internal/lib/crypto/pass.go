package crypto

import (
	"context"
	"crypto/sha256"
	"log/slog"

	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

type CryptoPassword struct {
	key    []byte
	logger *slog.Logger
}

func NewCryptoPassword(pass string, salt []byte, log *slog.Logger) *CryptoPassword {
	return &CryptoPassword{
		key:    pbkdf2.Key([]byte(pass), salt, IterCount, KeyLen, sha256.New),
		logger: log,
	}
}

func (c *CryptoPassword) GenerateFromPassword(ctx context.Context) ([]byte, error) {
	log := loghelper.New(c.logger, "crypto.GenerateFromPassword")

	hashKey, err := bcrypt.GenerateFromPassword(c.key, bcrypt.DefaultCost)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "bcrypt.GenerateFromPassword", err)
	}
	return hashKey, nil
}

func (c *CryptoPassword) CompareHashAndPassword(ctx context.Context, passHash []byte) error {
	log := loghelper.New(c.logger, "crypto.CompareHashAndPassword")
	if err := bcrypt.CompareHashAndPassword(passHash, c.key); err != nil {
		return log.LogAndReturnError(ctx, "bcrypt.CompareHashAndPassword", err)
	}
	return nil
}
