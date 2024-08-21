package main

import (
	"embed"

	"stars-server/app/storage/postgres"
	"stars-server/cmd"
)

//go:embed migrations/postgres/*.sql
var embedPGMigrations embed.FS

func main() {
	postgres.Migrations = embedPGMigrations

	cmd.Execute()
}
