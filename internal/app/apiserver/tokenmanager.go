package apiserver

import (
	"errors"
	"redditclone/internal/app/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager struct {
	signingKey string
	ttl        time.Duration
}

type AuthResponse map[string]string

func NewManager(signingKey string, tokenTTL string) (*TokenManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}
	ttl, err := time.ParseDuration(tokenTTL)
	if err != nil {
		return nil, err
	}
	return &TokenManager{
		signingKey: signingKey,
		ttl:        ttl,
	}, nil
}

func (tm *TokenManager) NewJWT(u *model.User) (string, error) {
	claims := &jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tm.ttl).Unix(),
		"user": &model.User{
			ID:       u.ID,
			Username: u.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tm.signingKey))
}

func (tm *TokenManager) Parse(accessToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tm.signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
