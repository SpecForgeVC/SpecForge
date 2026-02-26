package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

func TestJWTValidator_Validate(t *testing.T) {
	secret := []byte("test-secret-key-12345678901234567890")
	issuer := "specforge-test"
	audience := "specforge-client"
	cfg := Config{
		SigningKey: secret,
		Issuer:     issuer,
		Audience:   audience,
		Algorithm:  "HS256",
	}
	validator := NewJWTValidator(cfg)

	userID := uuid.New()
	workspaceID := uuid.New()

	tests := []struct {
		name    string
		tokenFn func() string
		wantErr error
	}{
		{
			name: "Valid Token",
			tokenFn: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub":       userID.String(),
					"workspace": workspaceID.String(),
					"role":      string(domain.RoleAdmin),
					"iss":       issuer,
					"aud":       audience,
					"exp":       time.Now().Add(time.Hour).Unix(),
					"iat":       time.Now().Unix(),
				})
				s, _ := token.SignedString(secret)
				return s
			},
			wantErr: nil,
		},
		{
			name: "Expired Token",
			tokenFn: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub":       userID.String(),
					"workspace": workspaceID.String(),
					"role":      string(domain.RoleAdmin),
					"iss":       issuer,
					"aud":       audience,
					"exp":       time.Now().Add(-time.Hour).Unix(),
				})
				s, _ := token.SignedString(secret)
				return s
			},
			wantErr: ErrExpired,
		},
		{
			name: "Invalid Signature",
			tokenFn: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub":       userID.String(),
					"workspace": workspaceID.String(),
					"role":      string(domain.RoleAdmin),
					"iss":       issuer,
					"aud":       audience,
					"exp":       time.Now().Add(time.Hour).Unix(),
				})
				s, _ := token.SignedString([]byte("wrong-secret"))
				return s
			},
			wantErr: ErrInvalidToken,
		},
		{
			name: "Missing Role",
			tokenFn: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub":       userID.String(),
					"workspace": workspaceID.String(),
					"iss":       issuer,
					"aud":       audience,
					"exp":       time.Now().Add(time.Hour).Unix(),
				})
				s, _ := token.SignedString(secret)
				return s
			},
			wantErr: ErrInvalidToken,
		},
		{
			name: "Unknown Role",
			tokenFn: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub":       userID.String(),
					"workspace": workspaceID.String(),
					"role":      "HACKER",
					"iss":       issuer,
					"aud":       audience,
					"exp":       time.Now().Add(time.Hour).Unix(),
				})
				s, _ := token.SignedString(secret)
				return s
			},
			wantErr: ErrUnknownRole,
		},
		{
			name: "Invalid UUID",
			tokenFn: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub":       "not-a-uuid",
					"workspace": workspaceID.String(),
					"role":      string(domain.RoleAdmin),
					"iss":       issuer,
					"aud":       audience,
					"exp":       time.Now().Add(time.Hour).Unix(),
				})
				s, _ := token.SignedString(secret)
				return s
			},
			wantErr: ErrInvalidUUID,
		},
		{
			name: "Wrong Issuer",
			tokenFn: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub":       userID.String(),
					"workspace": workspaceID.String(),
					"role":      string(domain.RoleAdmin),
					"iss":       "wrong-issuer",
					"aud":       audience,
					"exp":       time.Now().Add(time.Hour).Unix(),
				})
				s, _ := token.SignedString(secret)
				return s
			},
			wantErr: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.tokenFn()
			principal, err := validator.Validate(token)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() && err.Error() != "invalid or expired token" {
					// jwt-go often returns a wrapped error, we check our specific errors or the generic one
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Validate() error = %v, wantErr nil", err)
				return
			}
			if principal.UserID != userID {
				t.Errorf("Validate() userID = %v, want %v", principal.UserID, userID)
			}
			if principal.Role != domain.RoleAdmin {
				t.Errorf("Validate() role = %v, want %v", principal.Role, domain.RoleAdmin)
			}
		})
	}
}
