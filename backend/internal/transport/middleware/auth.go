package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/scott/specforge/internal/infra/auth"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func AuthMiddleware(validator *auth.JWTValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				renderError(w, http.StatusUnauthorized, "AUTH_FAILED", "Missing authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				renderError(w, http.StatusUnauthorized, "AUTH_FAILED", "Invalid authorization format")
				return
			}

			principal, err := validator.Validate(tokenString)
			if err != nil {
				renderError(w, http.StatusUnauthorized, "AUTH_FAILED", "Invalid or expired token")
				return
			}

			// Inject principal into context
			ctx := ContextWithPrincipal(r.Context(), principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func renderError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	json.NewEncoder(w).Encode(resp)
}
