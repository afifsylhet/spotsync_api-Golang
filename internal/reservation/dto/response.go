package dto

import (
	"time"

	userdto "github.com/afifsylhet/spotsync-api/internal/user/dto"
	zonedto "github.com/afifsylhet/spotsync-api/internal/zone/dto"
)

type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MyReservationResponse struct {
	ID           uint              `json:"id"`
	LicensePlate string            `json:"license_plate"`
	Status       string            `json:"status"`
	Zone         zonedto.ZoneBasic `json:"zone"`
	CreatedAt    time.Time         `json:"created_at"`
}

type AdminReservationResponse struct {
	ID           uint              `json:"id"`
	LicensePlate string            `json:"license_plate"`
	Status       string            `json:"status"`
	User         userdto.UserBasic `json:"user"`
	Zone         zonedto.ZoneBasic `json:"zone"`
	CreatedAt    time.Time         `json:"created_at"`
}
