package reservation

import (
	"errors"

	"github.com/afifsylhet/spotsync-api/internal/zone"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrZoneFull = errors.New("zone is at full capacity")

type ReservationRepository interface {
	CreateWithTransaction(userID, zoneID uint, licensePlate string) (*Reservation, error)
	FindByUserID(userID uint) ([]Reservation, error)
	FindByID(id uint) (*Reservation, error)
	Cancel(id uint) error
	FindAll() ([]Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateWithTransaction(userID, zoneID uint, licensePlate string) (*Reservation, error) {
	var newReservation Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var parkingZone zone.ParkingZone

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&parkingZone, zoneID).Error; err != nil {
			return err
		}

		var activeCount int64

		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		if activeCount >= int64(parkingZone.TotalCapacity) {
			return ErrZoneFull
		}

		newReservation = Reservation{
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

func (r *reservationRepository) FindByUserID(userID uint) ([]Reservation, error) {
	var reservations []Reservation

	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Find(&reservations).Error

	return reservations, err
}

func (r *reservationRepository) FindByID(id uint) (*Reservation, error) {
	var reservation Reservation

	err := r.db.First(&reservation, id).Error
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *reservationRepository) Cancel(id uint) error {
	return r.db.Model(&Reservation{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

func (r *reservationRepository) FindAll() ([]Reservation, error) {
	var reservations []Reservation

	err := r.db.Preload("User").
		Preload("Zone").
		Find(&reservations).Error

	return reservations, err
}
