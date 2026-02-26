package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrExpired      = errors.New("token is expired")
	ErrInvalidUUID  = errors.New("invalid UUID format in claims")
	ErrUnknownRole  = errors.New("unknown role in claims")
)

type JWTValidator struct {
	cfg Config
}

func NewJWTValidator(cfg Config) *JWTValidator {
	return &JWTValidator{cfg: cfg}
}

func (v *JWTValidator) Validate(tokenString string) (*domain.Principal, error) {
	parseOptions := []jwt.ParserOption{
		jwt.WithTimeFunc(time.Now),
	}
	if v.cfg.Issuer != "" {
		parseOptions = append(parseOptions, jwt.WithIssuer(v.cfg.Issuer))
	}
	if v.cfg.Audience != "" {
		parseOptions = append(parseOptions, jwt.WithAudience(v.cfg.Audience))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		alg := v.cfg.Algorithm
		if alg == "" {
			alg = "HS256"
		}

		if token.Method.Alg() != alg {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.cfg.SigningKey, nil
	}, parseOptions...)

	if err != nil {
		fmt.Printf("[AUTH_DEBUG] Token parse error: %v\n", err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpired
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		fmt.Printf("[AUTH_DEBUG] Token is invalid\n")
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Printf("[AUTH_DEBUG] Token claims are invalid\n")
		return nil, ErrInvalidToken
	}

	// Debug logs for claims
	fmt.Printf("[AUTH_DEBUG] Token sub: %v\n", claims["sub"])
	fmt.Printf("[AUTH_DEBUG] Token iss: %v\n", claims["iss"])
	fmt.Printf("[AUTH_DEBUG] Token aud: %v\n", claims["aud"])
	fmt.Printf("[AUTH_DEBUG] Token exp: %v\n", claims["exp"])

	// Extract and validate UserID (sub)
	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	// Extract and validate WorkspaceID (workspace)
	ws, ok := claims["workspace"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	workspaceID, err := uuid.Parse(ws)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	// Extract and validate Role
	roleStr, ok := claims["role"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	role := domain.Role(roleStr)
	switch role {
	case domain.RoleOwner, domain.RoleAdmin, domain.RoleReviewer, domain.RoleEngineer, domain.RoleAIAgent:
		// Valid role
	default:
		return nil, ErrUnknownRole
	}

	return &domain.Principal{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Role:        role,
	}, nil
}
