package service

import (
	"context"
	"fmt"
	"time"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

// UserService handles user business logic
type UserService struct {
	users store.UserStore
}

func NewUserService(users store.UserStore) *UserService {
	return &UserService{users: users}
}

func (s *UserService) FindAll(ctx context.Context) ([]*domain.User, error) {
	return s.users.FindAll(ctx)
}

func (s *UserService) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.users.FindByUsername(ctx, username)
}

func (s *UserService) Create(ctx context.Context, user *domain.User) (string, error) {
	if user.Username == "" {
		return "", fmt.Errorf("username is required")
	}
	if user.UserType == "" {
		return "", fmt.Errorf("userType is required")
	}

	// Hash password (default "pw" if not provided)
	if user.PasswordHash == "" {
		hash, err := HashPassword("pw")
		if err != nil {
			return "", fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hash
	}

	user.CreatedAt = time.Now()
	if user.FirmAccounts == nil {
		user.FirmAccounts = []domain.FirmAccount{}
	}
	if user.AssignedAccounts == nil {
		user.AssignedAccounts = []string{}
	}
	if user.ReadOnlyAccounts == nil {
		user.ReadOnlyAccounts = []string{}
	}

	return s.users.Insert(ctx, user)
}

func (s *UserService) Update(ctx context.Context, id string, user *domain.User) error {
	existing, err := s.users.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Preserve fields that shouldn't change
	user.ID = existing.ID
	user.Username = existing.Username
	user.PasswordHash = existing.PasswordHash
	user.CreatedAt = existing.CreatedAt

	return s.users.Update(ctx, id, user)
}
