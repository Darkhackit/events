package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Darkhackit/events/dto"
	"github.com/Darkhackit/events/sessions"
	"github.com/Darkhackit/events/token"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func AuthMiddleware(pasetoToken *token.PasetoToken, redisClient *sessions.RedisClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from the Authorization header
			currentRoute := mux.CurrentRoute(r)
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
				return
			}
			err = json.Unmarshal([]byte(result), &retrievedPayload)
			if err != nil {
				WriteResponse(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			var retrievePermissionPayload []dto.PermissionResponse
			permissions, err := redisClient.GetSession(r.Context(), "permissions_"+retrievedPayload.ID.String())
			err = json.Unmarshal([]byte(permissions), &retrievePermissionPayload)
			if err != nil {
				fmt.Println(err)
				WriteResponse(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			for _, permission := range retrievePermissionPayload {
				if strings.TrimSpace(strings.ToLower(permission.Name)) == strings.TrimSpace(strings.ToLower(currentRoute.GetName())) {
					next.ServeHTTP(w, r)
					return
				}
			}
			// Pass the request to the next handler
			WriteResponse(w, http.StatusForbidden, "permission denied")
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
