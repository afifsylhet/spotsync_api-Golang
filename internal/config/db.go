package config

import (
	"fmt"
	"log"
	"os"

	"github.com/afifsylhet/spotsync-api/internal/reservation"
	"github.com/afifsylhet/spotsync-api/internal/user"
	"github.com/afifsylhet/spotsync-api/internal/zone"
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
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	err = db.AutoMigrate(
		&user.User{},
		&zone.ParkingZone{},
		&reservation.Reservation{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("✅ Database connected and migrated successfully")

	return db
}
