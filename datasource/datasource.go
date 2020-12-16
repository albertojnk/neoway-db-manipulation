package datasource

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// StartDB simply start the connection with postgres
func StartDB() *gorm.DB {
	dsn := "host=database port=5432 user=admin sslmode=disable dbname=neoway password=password"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		log.Fatal("failed to connect database")
	}

	db.AutoMigrate(&ClientInfo{})

	log.Println("Successfully connected to database")
	return db
}

// GetDB returns an open instance of *gorm.DB
func GetDB() *gorm.DB {
	return db
}

// BulkCreateClientInfo insert all rows in a single INSERT
func BulkCreateClientInfo(clients []ClientInfo) error {
	chunks := sliceClients(clients, 500)

	for _, clients := range chunks {
		db.Create(&clients)
	}

	return nil
}

// sliceClients divide a slice into a series of slices bases on chunkSize
func sliceClients(slice []ClientInfo, chunkSize int) [][]ClientInfo {
	var chunks [][]ClientInfo

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
