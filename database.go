package main

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

var (
	db *gorm.DB
)

func initDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})

	if err != nil {
		return err
	}

	return db.AutoMigrate(&Device{})
}
