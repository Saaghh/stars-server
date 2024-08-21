/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"stars-server/app/config"
	"stars-server/app/logger"
	"stars-server/app/models"
	"stars-server/app/services/names"
	"stars-server/app/storage/postgres"
	"stars-server/app/worldgen"
)

// newgameCmd represents the newgame command
var newgameCmd = &cobra.Command{
	Use:   "newgame",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		generate()
	},
}

func init() {
	rootCmd.AddCommand(newgameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newgameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newgameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generate() {
	//context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	//logger
	logger.InitLogger(logger.Config{Level: cfg.LogLevel})

	//namesGen lib
	namesGen, err := names.New()
	if err != nil {
		zap.S().Panicf("newgameCmd/namesGen.New: %v", err)
	}

	// store layer
	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		zap.S().Panicf("newgameCmd/postgres.New: %v", err)
	}

	//generator
	worldGenerator := worldgen.New(db, namesGen)
	var game models.Game
	if game, err = worldGenerator.AutoGenerateGame(ctx); err != nil {
		zap.S().Panicf("newgameCmd/worldGenerator.AutoGenerateGame: %v", err)
	}

	zap.S().Infof("successfully created game id: %d", game.ID)
}
