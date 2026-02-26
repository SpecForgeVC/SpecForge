package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/auth"
)

func TestAuthMiddleware(t *testing.T) {
	secret := []byte("test-secret-key-12345678901234567890")
	issuer := "specforge-test"
	audience := "specforge-client"
	cfg := auth.Config{
		SigningKey: secret,
		Issuer:     issuer,
		Audience:   audience,
		Algorithm:  "HS256",
	}
	validator := auth.NewJWTValidator(cfg)
	middleware := AuthMiddleware(validator)

	userID := uuid.New()
	workspaceID := uuid.New()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := PrincipalFromContext(r.Context())
		if !ok {
			t.Error("Principal not found in context")
		}
		if p.UserID != userID {
			t.Errorf("UserID mismatch: got %v, want %v", p.UserID, userID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("Valid Token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":       userID.String(),
			"workspace": workspaceID.String(),
			"role":      string(domain.RoleAdmin),
			"iss":       issuer,
			"aud":       audience,
			"exp":       time.Now().Add(time.Hour).Unix(),
		})
		s, _ := token.SignedString(secret)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+s)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}
	})

	t.Run("Missing Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
		}

		var resp ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Error.Code != "AUTH_FAILED" {
			t.Errorf("Expected AUTH_FAILED error code, got %v", resp.Error.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
		}
	})
}

func TestRequireRole(t *testing.T) {
	middleware := RequireRole(domain.RoleAdmin, domain.RoleOwner)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("Allowed Role (Admin)", func(t *testing.T) {
		p := &domain.Principal{Role: domain.RoleAdmin}
		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(ContextWithPrincipal(req.Context(), p))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}
	})

	t.Run("Disallowed Role (Engineer)", func(t *testing.T) {
		p := &domain.Principal{Role: domain.RoleEngineer}
		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(ContextWithPrincipal(req.Context(), p))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusForbidden)
		}
	})

	t.Run("Missing Principal", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
		}
	})
}
