package models

import (
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// User структура
type User struct {
	ID           int64  `json:"id"`
	Login        string `json:"login"`
	PasswordHash []byte `json:"password_hash"` //был string
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int64  `json:"uid"`
	Login  string `json:"login"`
}

// SecretType алиас
type SecretType string

const (
	TypeCredentials SecretType = "credentials"
	TypeText        SecretType = "text_data"
	TypeBinary      SecretType = "binary_data"
	TypeCard        SecretType = "credit_card"
)

// Secret структура секрета
type Secret struct {
	ID        string          `json:"id"`
	UserID    int64           `json:"-"`
	Type      SecretType      `json:"type"`
	Data      []byte          `json:"data"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt int64           `json:"created_at"`
	UpdatedAt int64           `json:"updated_at"`
	IsDeleted bool            `json:"is_deleted"`
}

// ArraySecret список секретов
type ArraySecret []Secret

type contextKey string

// contextKey
const (
	UserIDKey contextKey = "userID"
)
