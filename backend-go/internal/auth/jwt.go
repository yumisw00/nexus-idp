package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("nexus-super-secret-key-2026") // In production, use ENV
type Claims struct {
	UserID string `json:"user_id"`
	OrgID  string `json:"org_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, orgID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		OrgID:  orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
