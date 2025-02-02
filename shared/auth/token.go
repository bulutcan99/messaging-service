package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret string
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{
		secret,
	}
}

func (t *TokenService) GenerateAccessToken(id, email string) (string, error) {
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    email,
		ID:        id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	ss, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (t *TokenService) GenerateRefreshToken(id, email string) (string, error) {
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(14 * 24 * time.Hour)),
		Issuer:    email,
		ID:        id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	ss, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (t *TokenService) VerifyAccessToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(t.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("access token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid access token")
}

func (t *TokenService) VerifyRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(t.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("refresh token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}
