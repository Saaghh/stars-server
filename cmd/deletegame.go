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
	"stars-server/app/processor"
	"stars-server/app/storage/postgres"
)

// deletegameCmd represents the deletegame command
var deletegameCmd = &cobra.Command{
	Use:   "deletegame",
	Short: "Delete whole game",
	Long:  `Delete game by id forever`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteGame()
	},
}

var gameID int

func init() {
	rootCmd.AddCommand(deletegameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deletegameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deletegameCmd.Flags().IntVar(&gameID, "id", 0, "ID of a game to delete")
}

func deleteGame() {
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

	if gameID == 0 {
		zap.S().Error("Please specify game id for deletion")
		return
	}

	// store layer
	pg, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("main/postgres.New")
	}

	// logic layer
	proc := processor.New(pg)

	if err = proc.DeleteWholeGame(ctx, gameID); err != nil {
		zap.S().Panicf("deletegameCmd/proc.DeleteWholeGame: %v", err)
		return
	}

	zap.S().Infof("Successfully deleted game id: %d", gameID)
}
