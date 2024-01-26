package database

import (
	"fmt"
	"logbook/internal/utilities/run"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../test.env")
	run.ExitAfterStderr("psql", "-U", "ufuktan", "-d", "postgres", "-f", "migration.sql")
	os.Exit(m.Run())
}

func ExampleTestMain() {
	fmt.Println("Hello world.") // Output: Hello world.
}
