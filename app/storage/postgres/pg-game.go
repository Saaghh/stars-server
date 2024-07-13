package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"stars-server/app/models"
)

func (p *Postgres) TxGetStellarBodies(ctx context.Context, filter models.StellarBodyFilter) ([]models.StellarBody, error) {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("p.getTXFromContext: %w", err)
	}

	query, args, err := psql.
		Select("*").
		From("stellar_bodies as sb").
		// Join("stellar_bodies_types sbt ON sbt.id = sb.type").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql.Select: %w", err)
	}

	var stellarBodies []models.StellarBody

	err = pgxscan.Select(ctx, tx, &stellarBodies, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgxscan.Select: %w", err)
	}

	return stellarBodies, nil
}

func (p *Postgres) TxGetSystems(ctx context.Context, filter models.StellarBodyFilter) ([]models.System, error) {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("p.getTXFromContext: %w", err)
	}

	query, args, err := psql.
		Select("*").
		From("systems").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql.Select: %w", err)
	}

	var systems []models.System

	err = pgxscan.Select(ctx, tx, &systems, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgxscan.Select: %w", err)
	}

	return systems, nil
}

func (p *Postgres) TxGetStellarBodyTypes(ctx context.Context) ([]models.StellarBodyType, error) {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("p.getTXFromContext: %w", err)
	}

	query, args, err := psql.
		Select("*").
		From("stellar_bodies_types").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql.Select: %w", err)
	}

	var types []models.StellarBodyType

	err = pgxscan.Select(ctx, tx, &types, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgxscan.Select: %w", err)
	}

	return types, nil
}

func (p *Postgres) TxUpdateStellarBodiesMovement(ctx context.Context, duration time.Duration, gameID int) error {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return fmt.Errorf("p.getTXFromContext: %w", err)
	}

	k := duration.Hours() / 24

	query := `
	UPDATE stellar_bodies 
	SET angle = MOD(angle + (angle_speed * $1), 360)
	WHERE system_id IN 
	      (SELECT id FROM systems WHERE game_id = $2)`

	_, err = tx.Exec(ctx, query, k, gameID)
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
}

func (p *Postgres) TxGetGames(ctx context.Context) ([]models.DBGame, error) {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("p.getTXFromContext: %w", err)
	}

	query, args, err := psql.Select("*").From("games").ToSql()
	if err != nil {
		return nil, fmt.Errorf("pgxscan.Select: %w", err)
	}

	var games []models.DBGame

	err = pgxscan.Select(ctx, tx, &games, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgxscan.Select: %w", err)
	}

	return games, nil
}
