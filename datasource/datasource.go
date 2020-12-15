package datasource

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

// StartDB simply start the connection with postgres
func StartDB() *gorm.DB {
	db, err = gorm.Open("postgres", "host=database port=5432 user=admin sslmode=disable dbname=neoway password=password")
	if err != nil {
		log.Fatal(err)
		log.Fatal("failed to connect database")
	}
	log.Println("Successfully connected to database")
	return db
}
