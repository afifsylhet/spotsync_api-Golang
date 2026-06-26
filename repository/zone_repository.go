
package repository


import (
	"gorm.io/gorm",
	"parking_zone.go/models"
)

type ZoneRepository interface {
Create(zone *models.ParkingZone) error
FindAll() ([]models.ParkingZone, error)
FindByID(id uint) (*models.ParkingZone, error)
Update(zone *models.ParkingZone) error
Delete(id uint) error
CountActiveReservations(zoneID uint) (int64, error)
}
type zoneRepository struct {
db *gorm.DB
}
func NewZoneRepository(db *gorm.DB) ZoneRepository {
return &zoneRepository{db: db}
}
func (r *zoneRepository) Create(zone *models.ParkingZone) error {
return r.db.Create(zone).Error
}
func (r *zoneRepository) FindAll() ([]models.ParkingZone, error) {
var zones []models.ParkingZone
err := r.db.Find(&zones).Error
return zones, err
}
func (r *zoneRepository) FindByID(id uint) (*models.ParkingZone, error) {
var zone models.ParkingZone
err := r.db.First(&zone, id).Error
if err != nil {
return nil, err
}
return &zone, nil
}
repository/reservation_repository.go
func (r *zoneRepository) Update(zone *models.ParkingZone) error {
return r.db.Save(zone).Error
}
func (r *zoneRepository) Delete(id uint) error {
return r.db.Delete(&models.ParkingZone{}, id).Error
}
// Counts active reservations for available_spots calculation
func (r *zoneRepository) CountActiveReservations(zoneID uint) (int64, error) {
var count int64
    err := r.db.Model(&models.Reservation{}).
Where("zone_id = ? AND status = ?", zoneID, "active").
Count(&count).Error
return count, err
}