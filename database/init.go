package database

import (
	"log"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitGORM(DSN string) {
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

func CloseDatabaseConnections() {
	sqlDb, err := Db.DB()
	if err != nil {
		log.Fatalln("Couldn't close database connections because failed to retrieve database object from GORM")
	}
	sqlDb.Close()
}
