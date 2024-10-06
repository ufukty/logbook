package database

import (
	"fmt"
	"log"
	"logbook/cmd/tags/service"
	"logbook/internal/utils/run"
)

func RunMigration(srvcfg *service.Config) error {
	output := run.ExitAfterStderr("psql",
		"-U", srvcfg.Database.User,
		"-d", srvcfg.Database.Default,
		"-c", "DROP DATABASE IF EXISTS "+srvcfg.Database.Name+";",
		"-c", "CREATE DATABASE "+srvcfg.Database.Name+";",
	)
	log.Println("dropping and recreating the application database:")
	fmt.Println(output)
	output = run.ExitAfterStderr("psql",
		"-U", srvcfg.Database.User,
		"-d", srvcfg.Database.Name,
		"-f", "../database/schema.sql", // working directory is not guaranteed
	)
	log.Println("building the application database:")
	fmt.Println(output)
	return nil
}
