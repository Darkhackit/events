package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Darkhackit/events/sessions"
	"github.com/Darkhackit/events/token"
	"net/http"
	"strings"
)

func AuthMiddleware(pasetoToken *token.PasetoToken, redisClient *sessions.RedisClient) func(http.Handler) http.Handler {
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
			var retrievedPayload token.Payload
			result, err := redisClient.GetSession(r.Context(), payload.ID.String())
			if err != nil {
				WriteResponse(w, http.StatusUnauthorized, "unauthorized")
			}
			err = json.Unmarshal([]byte(result), &retrievedPayload)
			if err != nil {
				WriteResponse(w, http.StatusUnauthorized, "unauthorized")
			}
			fmt.Println(retrievedPayload.IssuedAt)
			// Pass the request to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func GetAuthenticatedUser(r *http.Request) (*token.Payload, error) {
	payload, ok := r.Context().Value("auth").(*token.Payload)
	if !ok || payload == nil {
		return nil, errors.New("no user authenticated")
	}
	return payload, nil
}
