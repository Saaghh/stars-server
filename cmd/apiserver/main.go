package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
	"stars-server/app/config"
	"stars-server/app/logger"
	"stars-server/app/models"
	"stars-server/app/processor"
	"stars-server/app/server"
	"stars-server/app/storage/postgres"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// config
	cfgFile := os.Getenv("STARS_CONF")
	if cfgFile != "" {
		err := godotenv.Load(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cfg := config.NewFromEnv("STARS")

	// logger
	logger.InitLogger(logger.Config{Level: cfg.LogLevel})

	// no error handling for now
	// check https://github.com/uber-go/zap/issues/991
	//nolint: errcheck
	defer zap.L().Sync()

	// store layer
	pg, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("main/postgres.New")
	}

	if err = pg.Migrate(migrate.Up); err != nil {
		zap.L().With(zap.Error(err)).Fatal("main/migrate.Migrate")
	}

	// logic layer
	proc := processor.New(pg)

	// server layer
	srv, err := server.NewServer(cfg, proc)
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("main/server.NewServer")
	}

	// running server
	srv.RunServer()

	// debug
	_, err = proc.GetStellarBodies(ctx, models.StellarBodyFilter{})
	if err != nil {
		zap.L().With(zap.Error(err)).Error("main/proc.GetStellarBodies")
	}

	_, err = proc.GetSystems(ctx, models.StellarBodyFilter{})
	if err != nil {
		zap.L().With(zap.Error(err)).Error("main/proc.GetSystems")
	}

	_, err = proc.GetStellarBodyTypes(ctx)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("main/proc.GetStellarBodyTypes")
	}

	// shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	cancel()
}
