package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("B4VI09nHXgd0isxp7ph11pze0Qdy8O5x")

type JWTClaims struct {
	jwt.RegisteredClaims
}

func GenerateJWTToken(claims JWTClaims) (*string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func ValidateJWTToken(token string) (*JWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(jwtToken *jwt.Token) (any, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrUnknownClaimsType
	}

	return claims, nil
}
