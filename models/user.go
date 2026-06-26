go

package models
import "time"


type User struct {
    ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string `gorm:"not null" json:"name"`
    Email     string `gorm:"unique;not null" json:"email"`
    Password  string `gorm:"not null" json:"-"` // json:"-" hides password f
    Role      string `gorm:"default:driver;not null" json:"role"` // "driver
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}