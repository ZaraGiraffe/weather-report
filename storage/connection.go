// This package is responsible for creating a connection to the database
package storage

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"example.com/weather-report/config"
)

func createConnectionString(conf config.StorageConfig) string {
	return "host=" + conf.Host +
		" port=" + conf.Port +
		" user=" + conf.User +
		" password=" + conf.Password +
		" dbname=" + conf.Database +
		" sslmode=disable"
}

func newStorageConnection(conf config.StorageConfig) *sql.DB {
	connectionString := createConnectionString(conf)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	return db
}

func NewStorageConnection() *sql.DB {
	conf := config.GetConfig()
	conn := newStorageConnection(conf.StorageConfig)
	log.Println("INFO: Connected to database")
	return conn
}