package config

import (
	"fmt"
	"log"
	"os"

	"github.com/afifsylhet/spotsync-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // logs all SQL queries
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pooling (important for production)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	sqlDB.SetMaxOpenConns(25) // max simultaneous connections
	sqlDB.SetMaxIdleConns(10) // idle connections kept open

	// Auto-migrate: creates/updates tables based on models
	err = db.AutoMigrate(
		&models.User{},
		&models.ParkingZone{},
		&models.Reservation{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("✅ Database connected and migrated successfully")

	return db
}
