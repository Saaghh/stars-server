package postgres

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"stars-server/app/models"
)

//create

func (p *Postgres) TxCreateGame(ctx context.Context, game models.Game) (models.Game, error) {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return models.Game{}, fmt.Errorf("p.getTXFromContext: %w", err)
	}

	query, args, err := psql.
		Insert("games").
		Columns("owner_id", "world_time", "name").
		Values(game.Owner.ID, game.WorldTime, game.Name).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return models.Game{}, fmt.Errorf("psql.Insert.ToSql: %w", err)
	}

	if err = tx.QueryRow(ctx, query, args...).Scan(&game.ID); err != nil {
		return models.Game{}, fmt.Errorf("tx.QueryRow.Scan: %w", err)
	}

	return game, nil
}

func (p *Postgres) TxCreateSystem(ctx context.Context, starSystem models.System) (models.System, error) {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return models.System{}, fmt.Errorf("p.getTXFromContext: %w", err)
	}

	query, args, err := psql.
		Insert("systems").
		Columns("game_id", "name").
		Values(starSystem.GameID, starSystem.Name).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return models.System{}, fmt.Errorf("psql.Insert.ToSql: %w", err)
	}

	if err = tx.QueryRow(ctx, query, args...).Scan(&starSystem.ID); err != nil {
		return models.System{}, fmt.Errorf("tx.QueryRow.Scan: %w", err)
	}

	return starSystem, nil
}

func (p *Postgres) TxCreateBody(ctx context.Context, body models.StellarBody) (models.StellarBody, error) {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return models.StellarBody{}, fmt.Errorf("p.getTXFromContext: %w", err)
	}

	query, args, err := psql.
		Insert("stellar_bodies").
		Columns(
			"system_id",
			"name",
			"type_id",
			"mass",
			"diameter",
			"parent_body_id",
			"orbital_radius",
			"angle",
			"angle_speed",
		).
		Values(
			body.SystemID,
			body.Name,
			body.TypeID,
			body.Mass,
			body.Diameter,
			body.ParentBodyID,
			body.OrbitalRadius,
			body.Angle,
			body.AngleSpeed,
		).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return models.StellarBody{}, fmt.Errorf("psql.Insert.ToSql: %w", err)
	}

	if err = tx.QueryRow(ctx, query, args...).Scan(&body.ID); err != nil {
		return models.StellarBody{}, fmt.Errorf("tx.QueryRow.Scan: %w", err)
	}

	return body, nil
}

// update

func (p *Postgres) TxUpdateWorldTime(ctx context.Context, duration time.Duration, gameID int) error {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return fmt.Errorf("p.getTXFromContext: %w", err)
	}

	query := `UPDATE games SET world_time = world_time + INTERVAL '1 second' * $1 WHERE id = $2`

	if _, err = tx.Exec(ctx, query, duration.Seconds(), gameID); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
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

// get

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

func (p *Postgres) TxGetStellarBodies(ctx context.Context, filter models.StellarBodyFilter) (map[uuid.UUID]*models.StellarBody, error) {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("p.getTXFromContext: %w", err)
	}

	builder := psql.
		Select("*").
		From("stellar_bodies as sb")

	if filter.Systems != nil {
		builder = builder.Where(squirrel.Eq{"system_id": filter.Systems})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql.Select: %w", err)
	}

	var stellarBodies []models.StellarBody

	err = pgxscan.Select(ctx, tx, &stellarBodies, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgxscan.Select: %w", err)
	}

	result := make(map[uuid.UUID]*models.StellarBody, len(stellarBodies))
	for _, sb := range stellarBodies {
		result[sb.ID] = &sb
	}

	return result, nil
}

func (p *Postgres) TxGetStellarBodiesStockpiles(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]models.Stockpile, error) {
	tx, err := p.getTXFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("p.getTXFromContext: %w", err)
	}

	query, args, err := psql.
		Select(
			"sbs.stellar_body_id",
			"s.id",
			"s.quantity",
			"rt.id",
			"rt.density",
			"rt.name",
		).
		From("stockpiles as s").
		Join("stellar_bodies_stockpiles sbs on sbs.stockpile_id = s.id").
		Join("resource_types rt on rt.id = s.resource_type_id").
		Where(squirrel.Eq{"sbs.stellar_body_id": ids}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql.Select.ToSql: %w", err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("tx.Query: %w", err)
	}

	result := make(map[uuid.UUID][]models.Stockpile, len(ids))

	for rows.Next() {
		var sbID uuid.UUID
		var stockpile models.Stockpile

		if err = rows.Scan(
			&sbID,
			&stockpile.ID,
			&stockpile.Quantity,
			&stockpile.ResourceType.ID,
			&stockpile.ResourceType.Density,
			&stockpile.ResourceType.Name,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		result[sbID] = append(result[sbID], stockpile)
	}

	return result, nil
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

// delete

func (p *Postgres) DeleteGame(ctx context.Context, id int) error {
	query, args, err := psql.Delete("games").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("psql.Delete.ToSql: %w", err)
	}

	if ct, err := p.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("p.db.Exec: %w", err)
	} else {
		zap.S().Debugf("DeleteGame: deleted %d lines", ct.RowsAffected())
	}

	return nil
}
