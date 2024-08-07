package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
}

func GenerateJWTToken(secretKey string, claims JWTClaims) (*string, error) {
	fmt.Println(secretKey)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func ValidateJWTToken(secretKey string, token string) (*JWTClaims, error) {
	fmt.Println(secretKey)

	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(jwtToken *jwt.Token) (any, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}
		return []byte(secretKey), nil
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
