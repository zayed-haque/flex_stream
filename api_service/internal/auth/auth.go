package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(ctx context.Context) (*Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata provided")
	}

	tokenString := md["authorization"]
	if len(tokenString) == 0 {
		return nil, errors.New("no token provided")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString[0], claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
