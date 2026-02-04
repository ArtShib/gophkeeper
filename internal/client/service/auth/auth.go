package auth

import (
	"context"
	"log/slog"

	"github.com/ArtShib/gophkeeper/internal/client/models"
	"github.com/ArtShib/gophkeeper/internal/lib/crypto"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	models2 "github.com/ArtShib/gophkeeper/internal/server/models"
)

type StoreUser interface {
	AddUser(ctx context.Context, id int64, login string, passHash []byte) error
	GetUser(ctx context.Context, login string) (*models.User, error)
}

type Auth struct {
	log       *slog.Logger
	store     StoreUser
	cryptoPwd *crypto.CryptoPassword
}

func New(log *slog.Logger, store StoreUser, cryptoPwd *crypto.CryptoPassword) *Auth {
	return &Auth{
		log:       log,
		store:     store,
		cryptoPwd: cryptoPwd,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, id int64, login string) error {
	log := loghelper.New(a.log, "Auth.RegisterNewUser")
	log.LogDebug(ctx, "register user", slog.String("login", login))

	hashKey, err := a.cryptoPwd.GenerateFromPassword(ctx)
	if err != nil {
		return log.LogAndReturnError(ctx, "cryptoPwd.GenerateFromPassword", err)
	}

	if err = a.store.AddUser(ctx, id, login, hashKey); err != nil {
		return log.LogAndReturnError(ctx, "store.AddUser", err)
	}

	return nil
}

func (a *Auth) Login(ctx context.Context, login string) (*models.User, error) {
	log := loghelper.New(a.log, "Auth.Login")

	log.LogDebug(ctx, "login user", slog.String("login", login))

	user, err := a.store.GetUser(ctx, login)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "store.GetUser", err)
	}
	if err := a.cryptoPwd.CompareHashAndPassword(ctx, user.PasswordHash); err != nil {
		return nil, log.LogAndReturnError(ctx, "bcrypt.CompareHashAndPassword", models2.ErrInvalidCredentials)
	}

	log.LogDebug(ctx, "login success", slog.String("login", login))

	return user, nil
}
