package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type Payload struct {
	jwt.StandardClaims
	UserID    string `json:"sub"`
	SessionID string `json:"session_id"`
	Type      string `json:"type"`
	Refresh   bool   `json:"refresh"`
	Role      int    `json:"role"`
}

func Verify(token string, key string) (Payload, error) {
	if token == "" {
		return Payload{}, ErrInvalidToken
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Printf("jwt.ParseWithClaims: %v", ErrInvalidToken)
			return nil, ErrInvalidToken
		}
		return []byte(key), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		log.Printf("jwt.ParseWithClaims: %v", err)
		return Payload{}, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		log.Printf("Parsing to Payload: %v", ErrInvalidToken)
		return Payload{}, ErrInvalidToken
	}

	return *payload, nil
}

func Sign(payload Payload, expirationTime time.Duration, key string) (string, error) {
	payload.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expirationTime).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	return signedToken, nil
}
