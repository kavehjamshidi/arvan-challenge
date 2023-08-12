package postgres

import (
	"database/sql"
	"github.com/spf13/viper"
	"log"
)

func Setup() *sql.DB {
	pg, err := sql.Open("postgres", viper.GetString("DB_URI"))
	if err != nil {
		log.Panicf("could not connect to database: %v\n", err)
	}
	return pg
}

func TestSetup() *sql.DB {
	pg, err := sql.Open("postgres", viper.GetString("TEST_DB_URI"))
	if err != nil {
		log.Panicf("could not connect to database: %v\n", err)
	}
	return pg
}
