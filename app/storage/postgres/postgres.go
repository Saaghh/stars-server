package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"net/url"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"stars-server/app/config"
	"stars-server/app/models"
)

type Postgres struct {
	db  *pgxpool.Pool
	dsn string
}

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type txCtxKey struct{}

var Migrations embed.FS

func New(ctx context.Context, cfg config.PostgresConfig) (*Postgres, error) {
	urlScheme := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:     cfg.DataBase,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	dsn := urlScheme.String()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	zap.L().Info("successfully connected to database")

	return &Postgres{
		db:  db,
		dsn: dsn,
	}, nil
}

func (p *Postgres) MigrateUp(cfg config.PostgresConfig) error {
	db, err := sql.Open("pgx", p.dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			zap.L().With(zap.Error(err)).Error("conn.Close")
		}
	}()

	goose.SetBaseFS(Migrations)

	if err = goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect(postgres): %w", err)
	}

	if err = goose.Up(db, "migrations/postgres"); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	if err = db.Close(); err != nil {
		return fmt.Errorf("goose.Close: %w", err)
	}

	return nil
}

func (p *Postgres) getTXFromContext(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txCtxKey{}).(pgx.Tx)
	if !ok {
		return nil, models.ErrNoTx
	}

	return tx, nil
}

func (p *Postgres) WithTx(ctx context.Context, f func(context.Context) error) error {
	tx, err := p.getTXFromContext(ctx)
	if errors.Is(err, models.ErrNoTx) {
		tx, err = p.db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("p.pool.Begin(ctx): %w", err)
		}
		ctx = context.WithValue(ctx, txCtxKey{}, tx)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			zap.L().With(zap.Error(err)).Error("WithTx/tx.Rollback(ctx)")
		}
	}()

	if err = f(ctx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return errors.Join(fmt.Errorf("f(ctx): %w", err), fmt.Errorf("tx.Rollback(ctx): %w", err))
		}
		return fmt.Errorf("f(ctx): %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit(ctx): %w", err)
	}

	return nil
}
