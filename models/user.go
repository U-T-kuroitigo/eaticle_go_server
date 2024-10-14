package models

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	UserID      string `json:"user_id" gorm:"type:uuid;primaryKey;not null" validate:"uuid"`
	MailAddress string `json:"mail_address" gorm:"index:,unique;type:varchar(255);not null" validate:"required,email"`
	GmailID     string `json:"gmail_id" gorm:"unique;type:varchar(255);not null;size:255" validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateUser(user *User) error {
	err := validate.Struct(user)
	if err != nil {
		// バリデーションエラーが発生したフィールドをログに記録
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("Validation failed for field '%s' with tag '%s'", err.Field(), err.Tag())
		}
		return err
	}
	return nil
}
