package models

import (
	"github.com/google/uuid"
	"time"
)

type UserRole string

const (
	UserRole_Admin UserRole = "admin"
	UserRole_User  UserRole = "user"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         UserRole  `gorm:"type:varchar(50);default:user;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (u *User) NewUser(email string, password string) {
	u.ID = uuid.New()
	u.Email = email
	//u.Role = UserRole_Admin
	u.PasswordHash = password
}
