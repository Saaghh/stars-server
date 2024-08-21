/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"stars-server/app/config"
	"stars-server/app/logger"
	"stars-server/app/storage/postgres"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		migrateUp()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func migrateUp() {
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
	pg, err := postgres.New(context.Background(), cfg.Postgres)
	if err != nil {
		zap.L().With(zap.Error(err)).Panic("migrateUp/postgres.New")
	}

	if err = pg.MigrateUp(cfg.Postgres); err != nil {
		zap.L().With(zap.Error(err)).Panic("migrateUp/migrate.Migrate")
	}

	zap.L().Info("migration complete")
}
