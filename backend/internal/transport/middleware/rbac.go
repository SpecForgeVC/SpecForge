package middleware

import (
	"net/http"

	"github.com/scott/specforge/internal/domain"
)

func RequireRole(roles ...domain.Role) func(http.Handler) http.Handler {
	allowedRoles := make(map[domain.Role]bool)
	for _, r := range roles {
		allowedRoles[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := PrincipalFromContext(r.Context())
			if !ok {
				renderError(w, http.StatusUnauthorized, "AUTH_FAILED", "Authentication required")
				return
			}

			if !allowedRoles[principal.Role] {
				renderError(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
