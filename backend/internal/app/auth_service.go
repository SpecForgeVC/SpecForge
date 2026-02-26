package app

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (*AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error)
}

type authService struct {
	userRepo  UserRepository
	jwtSecret string
	issuer    string
	audience  string
}

func NewAuthService(userRepo UserRepository, jwtSecret, issuer, audience string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		issuer:    issuer,
		audience:  audience,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.verifyPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// For demo/dev purposes, we use the default workspace ID from seeding
	// In a real app, this would be looked up or passed in
	workspaceID := "00000000-0000-0000-0000-000000000001"

	return s.generateTokenResponse(user.ID.String(), workspaceID, string(user.Role))
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Simple validation for demo - in production, you should store refresh tokens in DB
	// and validate them against the store, implementing rotation and revocation.
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid refresh token claims")
	}

	// For simplicity, we assume the refresh token has all the necessary info
	// In production, you'd re-verify the user still exists and is active.
	sub := claims["sub"].(string)
	workspace := claims["workspace"].(string)
	role := claims["role"].(string)

	return s.generateTokenResponse(sub, workspace, role)
}

func (s *authService) generateTokenResponse(userID, workspaceID, role string) (*AuthResponse, error) {
	now := time.Now()
	accessTokenExp := now.Add(15 * time.Minute)
	refreshTokenExp := now.Add(7 * 24 * time.Hour)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       userID,
		"workspace": workspaceID,
		"role":      role,
		"iss":       s.issuer,
		"aud":       s.audience,
		"exp":       accessTokenExp.Unix(),
		"iat":       now.Unix(),
	})

	accessTokenStr, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       userID,
		"workspace": workspaceID,
		"role":      role,
		"iss":       s.issuer,
		"aud":       s.audience,
		"exp":       refreshTokenExp.Unix(),
		"iat":       now.Unix(),
	})

	refreshTokenStr, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(15 * 60),
	}, nil
}

func (s *authService) verifyPassword(password, hash string) bool {
	// Standard Argon2id format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false
	}

	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	keyLen := uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1
}
