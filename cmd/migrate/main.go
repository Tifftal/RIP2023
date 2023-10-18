package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"MSRM/internal/app/ds"
	"MSRM/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Миграция таблицы Samples
	err = db.AutoMigrate(&ds.Samples{})
	if err != nil {
		panic("failed to migrate Samples table")
	}

	// Миграция таблицы Missions
	err = db.AutoMigrate(&ds.Missions{})
	if err != nil {
		panic("failed to migrate Missions table")
	}

	// Миграция таблицы Users
	err = db.AutoMigrate(&ds.Users{})
	if err != nil {
		panic("failed to migrate Users table")
	}

	// Миграция таблицы Mission_samples
	err = db.AutoMigrate(&ds.Mission_samples{})
	if err != nil {
		panic("failed to migrate Mission_samples table")
	}
}
