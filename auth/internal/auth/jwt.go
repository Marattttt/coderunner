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
	TokType  string `json:"tokType"`
	jwt.RegisteredClaims
}

type TokType int

const (
	TokTypeAccess TokType = iota
	TokTypeRefresh
)

func (t TokType) String() string {
	switch t {
	case TokTypeAccess:
		return "access"
	case TokTypeRefresh:
		return "refresh"
	default:
		panic(fmt.Sprintf("%d is not a valid value for enum auth.TokType", t))
	}
}

func ClaimsFromUser(conf *config.AppConfig, u *models.User, toktype TokType) Claims {
	provider := "none"

	if u.Google != nil {
		provider = "google"
	}

	expiry := time.Now().Add(
		time.Duration(
			conf.JWTAccessExprirySecs * 1e9))

	return Claims{
		UserID:   u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Provider: provider,
		TokType:  toktype.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "coderunner",
		},
	}
}

func JWTGenerateAccess(conf *config.AppConfig, u *models.User) (string, error) {
	claims := ClaimsFromUser(conf, u, TokTypeAccess)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(conf.JWTSecret))
	return signed, err
}

func JWTGenerateRefresh(conf *config.AppConfig, u *models.User) (string, error) {
	claims := ClaimsFromUser(conf, u, TokTypeRefresh)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(conf.JWTSecret))
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
