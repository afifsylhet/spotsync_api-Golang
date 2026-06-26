package repository

import (
	"errors"

	"github.com/afifsylhet/spotsync-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Sentinel error – returned when zone is full
var ErrZoneFull = errors.New("zone is at full capacity")

type ReservationRepository interface {
	CreateWithTransaction(userID, zoneID uint, licensePlate string) (*models.Reservation, error)
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	Cancel(id uint) error
	FindAll() ([]models.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

// CreateWithTransaction – THE CONCURRENCY-SAFE RESERVATION
func (r *reservationRepository) CreateWithTransaction(userID, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var newReservation models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone

		// Step 1: Lock the parking zone row
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, zoneID).Error; err != nil {
			return err
		}

		// Step 2: Count active reservations
		var activeCount int64

		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// Step 3: Check capacity
		if activeCount >= int64(zone.TotalCapacity) {
			return ErrZoneFull
		}

		// Step 4: Create reservation
		newReservation = models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		return tx.Create(&newReservation).Error
	})

	if err != nil {
		return nil, err
	}

	return &newReservation, nil
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation

	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Find(&reservations).Error

	return reservations, err
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation

	err := r.db.First(&reservation, id).Error
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *reservationRepository) Cancel(id uint) error {
	return r.db.Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation

	err := r.db.Preload("User").
		Preload("Zone").
		Find(&reservations).Error

	return reservations, err
}
