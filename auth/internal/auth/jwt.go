package auth

import (
	"fmt"
	"time"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/Marattttt/new_new_portfolio/auth/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenInvalid = fmt.Errorf("Token is not valid")

type Claims struct {
	UserID   int    `json:"userId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
	jwt.RegisteredClaims
}

func ClaimsFromUser(conf *config.AppConfig, u models.User) Claims {
	provider := "none"

	if u.Google != nil {
		provider = "google"
	}

	expiry := time.Now().Add(
		time.Duration(
			conf.JWTAccessExprirySecs * 1e6))

	return Claims{
		UserID:   u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Provider: provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "coderunner",
		},
	}
}

func JWTGenerate(conf *config.AppConfig, u models.User) (string, error) {
	claims := ClaimsFromUser(conf, u)

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signed, err := token.SignedString(conf.JWTSecret)

	return signed, err
}

func JWTVerify(conf *config.AppConfig, tokenStr string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&claims,
		func(_ *jwt.Token) (any, error) {
			return conf.JWTSecret, nil
		},
	)

	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, ErrTokenInvalid
	}

	return claims, nil
}
