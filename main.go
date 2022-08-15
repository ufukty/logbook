package main

import (
	"logbook/main/database"
	"logbook/main/router"
)

func main() {
	database.Init(database.TEST_DB_DSN) // os.Getenv("DATABASE_URL")
	defer database.CloseConnection()

	router.StartRouter()
}
