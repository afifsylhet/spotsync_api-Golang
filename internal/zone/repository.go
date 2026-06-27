package zone

import "gorm.io/gorm"

type ZoneRepository interface {
	Create(zone *ParkingZone) error
	FindAll() ([]ParkingZone, error)
	FindByID(id uint) (*ParkingZone, error)
	Update(zone *ParkingZone) error
	Delete(id uint) error
	CountActiveReservations(zoneID uint) (int64, error)
}

type zoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db: db}
}

func (r *zoneRepository) Create(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *zoneRepository) FindAll() ([]ParkingZone, error) {
	var zones []ParkingZone
	err := r.db.Find(&zones).Error
	return zones, err
}

func (r *zoneRepository) FindByID(id uint) (*ParkingZone, error) {
	var zone ParkingZone
	err := r.db.First(&zone, id).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *zoneRepository) Update(zone *ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *zoneRepository) Delete(id uint) error {
	return r.db.Delete(&ParkingZone{}, id).Error
}

func (r *zoneRepository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64

	err := r.db.Table("reservations").
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error

	return count, err
}
