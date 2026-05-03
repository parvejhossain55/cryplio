package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	gjwt "github.com/golang-jwt/jwt/v5"
)

const (
	ClaimUserID    = "user_id"
	ClaimTokenType = "token_type"
)

// Claims is the token claim set used by the application.
type Claims map[string]any

// Issue signs an HS256 JWT with the provided claims and timestamps.
func Issue(secret string, expiry time.Duration, claims Claims) (string, error) {
	if secret == "" {
		return "", errors.New("jwt secret is required")
	}

	mapClaims := gjwt.MapClaims{}
	for key, value := range claims {
		mapClaims[key] = value
	}

	now := time.Now()
	mapClaims["iat"] = now.Unix()
	mapClaims["exp"] = now.Add(expiry).Unix()

	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(secret))
}

// Parse validates a signed HS256 JWT and returns its claims.
func Parse(secret, tokenString string) (gjwt.MapClaims, error) {
	token, err := gjwt.Parse(tokenString, func(token *gjwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gjwt.SigningMethodHMAC); !ok {
			return nil, gjwt.ErrTokenInvalidClaims
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(gjwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// FromRequest extracts a bearer token or cookie token from an HTTP request.
func FromRequest(cookieValue, authorizationHeader string) (string, error) {
	if cookieValue != "" {
		return cookieValue, nil
	}
	if authorizationHeader == "" {
		return "", errors.New("authorization required")
	}

	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
	if tokenString == authorizationHeader {
		return "", fmt.Errorf("bearer token required")
	}
	return tokenString, nil
}
