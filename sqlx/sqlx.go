package sqlx

import (
	"github.com/jmoiron/sqlx"
	"io"
	"log"
)

const connStr = "user=tailor dbname=tailor_core password=tailor sslmode=disable"

func Database() *sqlx.DB {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
