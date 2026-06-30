package repository

import (
	"context"
	"errors"

	"project/internal/models"
)

func (r *UserRepository) AddRefreshToken(ctx context.Context, refreshToken models.RefreshToken) error {
	const query = `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	rows, err := r.db.Exec(ctx, query, refreshToken.UserID, refreshToken.TokenHash, refreshToken.ExpiresAt)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return errors.New("could not insert refresh token")
	}

	return nil
}

func (r *UserRepository) GetRefreshToken(ctx context.Context, tokenHash string) (models.RefreshToken, error) {
	const query = `SELECT user_id, token_hash, expires_at FROM refresh_tokens WHERE token_hash = $1`
	
	var t models.RefreshToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&t.UserID, &t.TokenHash, &t.ExpiresAt)
	if err != nil {
		return models.RefreshToken{}, err
	}
	return t, nil
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	const query = `DELETE FROM refresh_tokens WHERE token_hash = $1`
	
	_, err := r.db.Exec(ctx, query, tokenHash)
	if err != nil {
		return err
	}
	return nil
}
