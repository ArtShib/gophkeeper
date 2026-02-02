package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/ArtShib/gophkeeper/internal/lib/jwt"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	models2 "github.com/ArtShib/gophkeeper/internal/server/models"
	"golang.org/x/crypto/bcrypt"
)

type StoreUser interface {
	AddUser(ctx context.Context, login string, passHash []byte, createdAT int64) (*models2.User, error)
	GetUser(ctx context.Context, login string) (*models2.User, error)
}

type Auth struct {
	log      *slog.Logger
	store    StoreUser
	tokenTTL time.Duration
	config   *models2.ConfigJWT
}

func New(log *slog.Logger, store StoreUser, config *models2.ConfigJWT) *Auth {
	return &Auth{
		log:    log,
		store:  store,
		config: config,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, login string, passHash string) (int64, error) {
	log := loghelper.New(a.log, "Auth.RegisterNewUser")

	log.LogDebug(ctx, "register user", slog.String("login", login))

	passHashSrv, err := bcrypt.GenerateFromPassword([]byte(passHash), bcrypt.DefaultCost)
	if err != nil {
		return 0, log.LogAndReturnError(ctx, "bcrypt.GenerateFromPassword", err)
	}

	user, err := a.store.AddUser(ctx, login, passHashSrv, time.Now().Unix())
	if err != nil {
		return 0, log.LogAndReturnError(ctx, "store.AddUser", err)
	}

	return user.ID, nil
}

func (a *Auth) Login(ctx context.Context, login string, passHash string) (string, error) {
	log := loghelper.New(a.log, "Auth.Login")

	log.LogDebug(ctx, "login user", slog.String("login", login))

	user, err := a.store.GetUser(ctx, login)
	if err != nil {
		return "", log.LogAndReturnError(ctx, "store.GetUser", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(passHash)); err != nil {
		return "", log.LogAndReturnError(ctx, "bcrypt.CompareHashAndPassword", models2.ErrInvalidCredentials)
	}

	log.LogDebug(ctx, "login success", slog.String("login", login))

	token, err := jwt.NewToken(user, a.tokenTTL, a.config.SecretKey)
	if err != nil {
		return "", log.LogAndReturnError(ctx, "jwt.NewToken", err)
	}
	return token, nil
}

func (a *Auth) ParseToken(ctx context.Context, tokenString string) (int64, error) {
	log := loghelper.New(a.log, "Auth.ParseToken")
	log.LogDebug(ctx, "parse token", slog.String("token", tokenString))

	token, err := jwt.ParseToken(tokenString, a.config.SecretKey)

	if err != nil {
		return 0, log.LogAndReturnError(ctx, "jwt.ParseToken", err)
	}
	log.LogDebug(ctx, "login success", slog.String("token", tokenString))
	return token.UserID, nil
}
