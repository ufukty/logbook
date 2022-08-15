package database

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	TEST_DB_DSN = "host=localhost user=ufuktan password=password dbname=logbook_dev port=5432 sslmode=disable TimeZone=utc"
)

var Db *gorm.DB

func initGORM(DSN string) {
	var err error
	Db, err = gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  DSN,
				PreferSimpleProtocol: false, // when True: it disables implicit prepared statement usage
			},
		),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Error on connecting to database trough GORM"))
	}
}

func CloseConnection() {
	sqlDb, err := Db.DB()
	if err != nil {
		log.Fatalln("Couldn't close database connections because failed to retrieve database object from GORM")
	}
	sqlDb.Close()
}

func ReloadTestDatabase() {
	cmd := exec.Command("make", "migrate")
	cmd.Dir = "../"

	// cmd.Stdin = strings.NewReader("and old falcon")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	err := cmd.Run()

	if err != nil {
		log.Fatalln(stdout.String(), errors.Wrap(err, "ReloadTestDatabase()"))
	}
	fmt.Println(stdout.String())
}

var regularExpressionForSQLStateNumber *regexp.Regexp

func initRegex() {
	var err error
	regularExpressionForSQLStateNumber, err = regexp.Compile(`\(SQLSTATE ([0-9]*)\)$`)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "initRegex() Could not compile regex"))
	}
}

// Use this function to strip SQLSTATE error code at the end
// of error message sent by PostgreSQL
// Meaning of error codes:
// https://www.postgresql.org/docs/14/errcodes-appendix.html
func StripSQLState(errorMessage string) string {
	results := regularExpressionForSQLStateNumber.FindStringSubmatch(errorMessage)
	if len(results) < 2 {
		return "-1"
	} else {
		return results[1]
	}
}

func Init(DSN string) {
	initRegex()
	initGORM(DSN)
}
