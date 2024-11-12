package api

import (
	"context"
	"errors"
	"github.com/Darkhackit/events/token"
	"net/http"
	"strings"
)

func AuthMiddleware(pasetoToken *token.PasetoToken) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				WriteResponse(w, http.StatusUnauthorized, "unauthorised")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Verify the token
			payload, err := pasetoToken.VerifyToken(tokenString)
			if err != nil {
				WriteResponse(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// Add the payload to the request context for downstream handlers
			ctx := context.WithValue(r.Context(), "auth", payload)
			r = r.WithContext(ctx)

			// Pass the request to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func getAuthenticatedUser(r *http.Request) (*token.Payload, error) {
	payload, ok := r.Context().Value("auth").(*token.Payload)
	if !ok || payload == nil {
		return nil, errors.New("no user authenticated")
	}
	return payload, nil
}
