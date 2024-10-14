package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	UserID      string `json:"user_id" gorm:"type:uuid;primaryKey;not null" validate:"max=32"`
	MailAddress string `json:"mail_address" gorm:"index:,unique;type:varchar(255);not null"`
	GmailID     string `json:"gmail_id" gorm:"unique;type:varchar(255);not null;size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateUser(user *User) error {
	return validate.Struct(user)
}
