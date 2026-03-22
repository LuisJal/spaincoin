package auth

import (
	"context"
	"net/http"
	"strings"
)

// contextKey is the unexported key type used to store values in request contexts.
type contextKey string

// ClaimsKey is the context key under which validated JWT Claims are stored.
const ClaimsKey contextKey = "claims"

// AuthMiddleware extracts the Bearer token from the Authorization header,
// validates it, and stores the resulting Claims in the request context.
// Requests that are missing a valid token receive a 401 response.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"authorization header required"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, `{"error":"authorization header must be Bearer <token>"}`, http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// GetClaims extracts the Claims stored in r's context by AuthMiddleware.
// Returns nil when no claims are present.
func GetClaims(r *http.Request) *Claims {
	v := r.Context().Value(ClaimsKey)
	if v == nil {
		return nil
	}
	claims, _ := v.(*Claims)
	return claims
}
