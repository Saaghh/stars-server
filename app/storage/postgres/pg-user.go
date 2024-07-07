package postgres

import (
	"context"
	"fmt"

	"stars-server/app/models"
)

func (p *Postgres) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`

	err := p.db.QueryRow(ctx, query, user.Username, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("p.db.QueryRow: %w", err)
	}

	return user, nil
}
