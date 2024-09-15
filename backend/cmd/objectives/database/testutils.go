package database

import (
	"fmt"
	"log"
	"logbook/cmd/objectives/service"
	"logbook/internal/utilities/run"
	"os"
	"path/filepath"
)

var migrationfile = filepath.Join(os.Getenv("WORKSPACE"), "backend/cmd/objectives/database/schema.sql")

func RunMigration(cfg *service.Config) error {
	output := run.ExitAfterStderr("psql",
		"-U", cfg.Database.User,
		"-d", cfg.Database.Default,
		"-c", "DROP DATABASE IF EXISTS "+cfg.Database.Name+";",
		"-c", "CREATE DATABASE "+cfg.Database.Name+";",
	)
	log.Println("dropping and recreating the application database:")
	fmt.Println(output)
	output = run.ExitAfterStderr("psql",
		"-U", cfg.Database.User,
		"-d", cfg.Database.Name,
		"-f", migrationfile,
	)
	log.Println("building the application database:")
	fmt.Println(output)
	return nil
}
