package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"project/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
	);
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) Add(ctx context.Context, user models.User) error {
	const query = `INSERT INTO users (name,email, password_hash,role) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	const query = `
			SELECT
				id,
				name,
				email,
				password_hash,
				role,
				created_at
			FROM users
			WHERE email = $1 AND deleted_at IS NULL
`
	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
