package auth

import (
	"context"
	"log/slog"

	"github.com/ArtShib/gophkeeper/internal/lib/jwt"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage[T any] interface {
	AddUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, login string) (*models.User, error)
}

type Auth struct {
	log    *slog.Logger
	store  UserStorage[models.User]
	config *models.ConfigJWT
}

func New(log *slog.Logger, store UserStorage[models.User], config *models.ConfigJWT) *Auth {
	return &Auth{
		log:    log,
		store:  store,
		config: config,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, login string, passHash []byte) (int64, error) {
	log := loghelper.New(a.log, "Auth.RegisterNewUser")

	log.LogDebug(ctx, "register user", slog.String("login", login))

	passHashSrv, err := bcrypt.GenerateFromPassword(passHash, bcrypt.DefaultCost)
	if err != nil {
		return 0, log.LogAndReturnError(ctx, "bcrypt.GenerateFromPassword", err)
	}
	user := &models.User{
		Login:        login,
		PasswordHash: passHashSrv,
	}

	id, err := a.store.AddUser(ctx, user)
	if err != nil {
		return 0, log.LogAndReturnError(ctx, "store.AddUser", err)
	}

	return id, nil
}

func (a *Auth) Login(ctx context.Context, login string, passHash []byte) (string, error) {
	log := loghelper.New(a.log, "Auth.Login")

	log.LogDebug(ctx, "login user", slog.String("login", login))

	user, err := a.store.GetUser(ctx, login)
	if err != nil {
		return "", log.LogAndReturnError(ctx, "store.GetUser", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, passHash); err != nil {
		return "", log.LogAndReturnError(ctx, "bcrypt.CompareHashAndPassword", err) // models.ErrInvalidCredentials)
	}

	log.LogDebug(ctx, "login success", slog.String("login", login))

	token, err := jwt.NewToken(user, a.config.TokenTTLMIN, a.config.SecretKey)
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
