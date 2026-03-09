package service

import (
	"context"
	"fmt"
	"time"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication
type AuthService struct {
	users store.UserStore
}

func NewAuthService(users store.UserStore) *AuthService {
	return &AuthService{users: users}
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, username, password string) (*domain.User, error) {
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if user.Disabled {
		return nil, fmt.Errorf("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	now := time.Now()
	err = s.users.UpdateFields(ctx, user.ID, map[string]interface{}{
		"is_logged_in": true,
		"last_login_at": now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update login state: %w", err)
	}

	user.IsLoggedIn = true
	user.LastLoginAt = &now
	return user, nil
}

// Logout logs out a user
func (s *AuthService) Logout(ctx context.Context, username string) error {
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	return s.users.UpdateFields(ctx, user.ID, map[string]interface{}{
		"is_logged_in": false,
	})
}

// HashPassword hashes a plain-text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
