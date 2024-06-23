package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const userIDKey contextKey = "userID"

// WithAuth validates JWT tokens from HTTP requests and adds the user ID to the request's context.
// If the token is invalid or missing, and the path is not blacklisted, then it returns a 401
// Unauthorized response.
func WithAuth(jwtSecret string, blacklistedPaths []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// Skip JWT authentication for paths that do not require it.
			for _, p := range blacklistedPaths {
				if strings.HasPrefix(path, p) {
					next.ServeHTTP(w, r)
					return
				}
			}

			handleUnauthorized := func(err error) {
				slog.Error("Unauthorized", slog.Any("error", err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			token, err := jwtFromRequest(r, jwtSecret)
			if err != nil {
				handleUnauthorized(err)
				return
			}

			userID, err := userIDFromJWT(*token)
			if err != nil {
				handleUnauthorized(err)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// jwtFromRequest extracts a JWT token from an HTTP request's Authorization header and validates it.
func jwtFromRequest(r *http.Request, jwtSecret string) (*jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			// Make sure to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			// Return the key for validating the token. For example, a shared secret:
			return []byte(jwtSecret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// userIDFromJWT extracts a user ID from a JWT token's claims.
func userIDFromJWT(token jwt.Token) (uuid.UUID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("missing claims")
	}

	subID, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("missing sub in claims")
	}

	return uuid.Parse(subID)
}

// GetUserID retrieves a user ID from a context, assuming such context comes from an HTTP request
// with JWT authentication.
func GetUserID(ctx context.Context) uuid.UUID {
	return ctx.Value(userIDKey).(uuid.UUID)
}
