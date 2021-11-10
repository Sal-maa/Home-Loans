package contract

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique" json:"username" validate:"required"`
	Password string `gorm:"not null" json:"password" validate:"required,gte=6"`
	LoginAs  uint   `gorm:"not null" json:"login_as"`
}

type UserReturn struct {
	Username string `gorm:"not null" json:"username" validate:"required"`
	LoginAs  uint   `gorm:"not null" json:"login_as"`
}
