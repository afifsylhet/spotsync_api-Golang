go
package dto
import "time"
// --- REQUEST DTOs --
type CreateReservationRequest struct {
    ZoneID       uint `json:"zone_id" validate:"required"`
    LicensePlate string `json:"license_plate" validate:"required,max=15"`
}
// --- RESPONSE DTOs --
type ReservationResponse struct {
    ID           uint `json:"id"`
    UserID       uint `json:"user_id"`
    ZoneID       uint `json:"zone_id"`
    LicensePlate string `json:"license_plate"`
    Status       string `json:"status"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
// Used in "My Reservations" - includes zone info
type MyReservationResponse struct {
    ID           uint `json:"id"`
    LicensePlate string `json:"license_plate"`
    Status       string `json:"status"`
    Zone         ZoneBasic `json:"zone"`
    CreatedAt    time.Time `json:"created_at"`
}
type ZoneBasic struct {
    ID   uint `json:"id"`
    Name string `json:"name"`
    Type string `json:"type"`
}
// Used in Admin "Get All Reservations"
type AdminReservationResponse struct {
    ID           uint `json:"id"`
    LicensePlate string `json:"license_plate"`
    Status       string `json:"status"`
8. Repository Layer (
repository/
 folder)
Repositories ONLY talk to the database. No HTTP, no business logic.
repository/user_repository.go
    User         UserBasic     `json:"user"`
    Zone         ZoneBasic     `json:"zone"`
    CreatedAt    time.Time     `json:"created_at"`
}
type UserBasic struct {
    ID    uint `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}