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
	migrate "github.com/rubenv/sql-migrate"
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

//go:embed migrations
var migrations embed.FS

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

func (p *Postgres) Migrate(direction migrate.MigrationDirection) error {
	conn, err := sql.Open("pgx", p.dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			zap.L().With(zap.Error(err)).Error("conn.Close")
		}
	}()

	assetDir := func() func(string) ([]string, error) {
		return func(path string) ([]string, error) {
			dirEntry, err := migrations.ReadDir(path)
			if err != nil {
				return nil, fmt.Errorf("migrations.ReadDir: %w", err)
			}

			entries := make([]string, 0)

			for _, e := range dirEntry {
				entries = append(entries, e.Name())
			}

			return entries, nil
		}
	}()

	asset := migrate.AssetMigrationSource{
		Asset:    migrations.ReadFile,
		AssetDir: assetDir,
		Dir:      "migrations",
	}

	_, err = migrate.Exec(conn, "postgres", asset, direction)
	if err != nil {
		return fmt.Errorf("migrate.Exec: %w", err)
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
		var err error
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

	if err := f(ctx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return errors.Join(fmt.Errorf("f(ctx): %w", err), fmt.Errorf("tx.Rollback(ctx): %w", err))
		}
		return fmt.Errorf("f(ctx): %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit(ctx): %w", err)
	}

	return nil
}
