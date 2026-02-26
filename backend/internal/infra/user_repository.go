package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) app.UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.queries.DB().QueryRowContext(ctx, "SELECT id, email, full_name, role, password_hash, created_at, updated_at FROM users WHERE email = $1", email)
	var u domain.User
	var role string
	err := row.Scan(&u.ID, &u.Email, &u.FullName, &role, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	u.Role = domain.Role(role)
	return &u, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row := r.queries.DB().QueryRowContext(ctx, "SELECT id, email, full_name, role, password_hash, created_at, updated_at FROM users WHERE id = $1", id)
	var u domain.User
	var role string
	err := row.Scan(&u.ID, &u.Email, &u.FullName, &role, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	u.Role = domain.Role(role)
	return &u, nil
}
