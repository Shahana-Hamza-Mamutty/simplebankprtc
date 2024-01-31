package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// func payloadToMapClaims(p *Payload) jwt.MapClaims {
// 	return jwt.MapClaims{
// 		"id":         p.ID,
// 		"username":   p.Username,
// 		"role":       p.Role,
// 		"issued_at":  p.IssuedAt,
// 		"expired_at": p.ExpiredAt,
// 		jwt.RegisteredClaims{
// 			// Also fixed dates can be used for the NumericDate
// 			ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
// 			Issuer:    "test",
// 		},
// 	}
// }

func (maker *JWTMaker) CreateToken(username string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration)

	if err != nil {
		return "", payload, err
	}
	// payloadToMapClaims(payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
