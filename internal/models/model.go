package models

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrLoginOrPassIsEmpty = errors.New("login and password_hash required")
	ErrSyncPullChanges    = errors.New("sync svc pull changes")
	ErrSyncPushChanges    = errors.New("sync svc push changes")
	ErrSyncService        = errors.New("sync error")
	ErrDateSecretServer   = errors.New("small date on server")
)

// User структура
type User struct {
	ID           int64  `json:"id"`
	Login        string `json:"login"`
	PasswordHash []byte `json:"password_hash"`
	Token        string
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int64  `json:"uid"`
	Login  string `json:"login"`
}

// contextKey
type contextKey string

const (
	UserIDKey contextKey = "userID"
)

type ListSecrets map[string]Secret
