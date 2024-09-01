package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey"`
	Name       string         `gorm:"size:255;not null"`
	Mail       string         `gorm:"size:255;unique;not null" validate:"required,email"`
	Pass       string         `gorm:"size:255;not null" validate:"required,min=8"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	ModifiedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
