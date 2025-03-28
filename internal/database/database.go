package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbConnection *gorm.DB

func setup() {
	// Read db path from `.env`
	dbPath := os.Getenv("DB_PATH")

	if dbPath == "" {
		dbPath = "iori_sync.db"
	}

	// Initialize the database connection
	var err error
	dbConnection, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	dbConnection.AutoMigrate(&WatchDirectory{}, &PathRule{})
	// Log the successful connection
	log.Println("Connected to the database successfully.")
}

func GetDB() *gorm.DB {
	if dbConnection == nil {
		setup()
	}

	return dbConnection
}
