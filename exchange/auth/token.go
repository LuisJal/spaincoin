package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims contains the custom fields embedded in the JWT alongside standard claims.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// devSecret is used when SPC_JWT_SECRET is not set.  It is intentionally weak
// and must never be used in production.
const devSecret = "spaincoin-dev-secret-change-me-in-production"

// GetJWTSecret reads SPC_JWT_SECRET from the environment.  Falls back to a
// development placeholder when the variable is absent.
// IMPORTANT: in production SPC_JWT_SECRET MUST be set to a long random string.
func GetJWTSecret() []byte {
	if s := os.Getenv("SPC_JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte(devSecret)
}

// GenerateToken creates a signed JWT valid for 7 days that carries userID and
// email as custom claims.
func GenerateToken(userID, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			Issuer:    "spaincoin-exchange",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// ValidateToken parses and validates tokenString, returning the embedded Claims.
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
