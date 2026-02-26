package jwt

import (
	"fmt"
	"time"

	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user *models.User, duration time.Duration, secretKey []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string, secretKey []byte) (*models.UserClaims, error) {
	claims := &models.UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

type tokenClaimsParse struct {
	UID   int64
	Login string
	EXP   int64
	jwt.RegisteredClaims
}

func ParseTokenUserID(tokenString string) (int64, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &tokenClaimsParse{})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*tokenClaimsParse); ok {
		return claims.UID, nil
	}
	return 0, fmt.Errorf("invalid token")
}
