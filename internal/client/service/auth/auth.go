package auth

import (
	"context"
	"errors"
	"log/slog"

	authgrpc "github.com/ArtShib/gophkeeper/internal/client/grpc"
	"github.com/ArtShib/gophkeeper/internal/lib/crypto"
	"github.com/ArtShib/gophkeeper/internal/lib/jwt"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
	models2 "github.com/ArtShib/gophkeeper/internal/models"
)

type StoreUser interface {
	AddUser(ctx context.Context, id int64, login string, passHash []byte) error
	GetUser(ctx context.Context, login string) (*models.User, error)
}

type Auth struct {
	log             *slog.Logger
	store           StoreUser
	cryptoPwdSever  *crypto.CryptoPassword
	cryptoPwdClient *crypto.CryptoPassword
	authGRPC        *authgrpc.AuthClient
	login           string
}

func New(log *slog.Logger, store StoreUser, password string, authGRPC *authgrpc.AuthClient, cfg models.ConfigCrypto) *Auth {
	cryptoPwdSever := crypto.NewCryptoPassword(password, cfg.SaltServer, log)
	cryptoPwdClient := crypto.NewCryptoPassword(password, cfg.SaltClient, log)
	return &Auth{
		log:             log,
		store:           store,
		cryptoPwdSever:  cryptoPwdSever,
		cryptoPwdClient: cryptoPwdClient,
		authGRPC:        authGRPC,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context) error {
	log := loghelper.New(a.log, "Auth.RegisterNewUser")
	log.LogDebug(ctx, "register user", slog.String("login", a.login))

	hashKeyServer, err := a.cryptoPwdSever.GenerateFromPassword(ctx)
	if err != nil {
		return log.LogAndReturnError(ctx, "cryptoPwdSever.GenerateFromPassword", err)
	}
	id, err := a.authGRPC.Register(ctx, a.login, hashKeyServer)
	if err != nil {
		return log.LogAndReturnError(ctx, "authGRPC.Register", err)
	}
	hashKeyClient, err := a.cryptoPwdClient.GenerateFromPassword(ctx)
	if err != nil {
		return log.LogAndReturnError(ctx, "cryptoPwdClient.GenerateFromPassword", err)
	}

	if err = a.store.AddUser(ctx, id, a.login, hashKeyClient); err != nil {
		return log.LogAndReturnError(ctx, "store.AddUser", err)
	}

	return nil
}

func (a *Auth) Login(ctx context.Context) (*models.User, error) {
	log := loghelper.New(a.log, "Auth.Login")
	log.LogDebug(ctx, "login user", slog.String("login", a.login))

	hashKeyServer, err := a.cryptoPwdSever.GenerateFromPassword(ctx)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "cryptoPwdSever.GenerateFromPassword", err)
	}

	token, err := a.authGRPC.Login(ctx, a.login, hashKeyServer)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "authGRPC.Register", err)
	}

	user, err := a.store.GetUser(ctx, a.login)
	if err != nil {
		if errors.Is(models.ErrUserNotFound, err) {
			//return nil, log.LogAndReturnError(ctx, "store.GetUser", err)
			hashKeyClient, err := a.cryptoPwdClient.GenerateFromPassword(ctx)
			if err != nil {
				return nil, log.LogAndReturnError(ctx, "cryptoPwdClient.GenerateFromPassword", err)
			}
			id, err := jwt.ParseTokenUserID(token)
			if err != nil {
				return nil, log.LogAndReturnError(ctx, "jwt.ParseTokenUserID", err)
			}

			if err = a.store.AddUser(ctx, id, a.login, hashKeyClient); err != nil {
				return nil, log.LogAndReturnError(ctx, "store.AddUser", err)
			}
		}
		return nil, log.LogAndReturnError(ctx, "store.GetUser", err)
	}

	if err := a.cryptoPwdClient.CompareHashAndPassword(ctx, user.PasswordHash); err != nil {
		return nil, log.LogAndReturnError(ctx, "bcrypt.CompareHashAndPassword", models2.ErrInvalidCredentials)
	}
	user.Token = token
	log.LogDebug(ctx, "login success", slog.String("login", a.login))
	return user, nil
}
