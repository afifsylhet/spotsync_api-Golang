package reservation

import (
	"time"

	"github.com/afifsylhet/spotsync-api/internal/user"
	"github.com/afifsylhet/spotsync-api/internal/zone"
)

type Reservation struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	ZoneID       uint      `gorm:"not null" json:"zone_id"`
	LicensePlate string    `gorm:"not null;size:15" json:"license_plate"`
	Status       string    `gorm:"default:active;not null" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         user.User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Zone         zone.ParkingZone `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
}
