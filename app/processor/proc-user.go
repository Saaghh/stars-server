package processor

import (
	"context"
	"fmt"

	"stars-server/app/models"
)

type user interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

func (p *Processor) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var err error

	user, err = p.db.CreateUser(ctx, user)
	if err != nil {
		return models.User{}, fmt.Errorf("p.db.CreateUser: %w", err)
	}

	return user, nil
}
