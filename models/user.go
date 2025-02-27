package models

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	UserID       string `json:"user_id" gorm:"type:varchar(255);primaryKey;not null" validate:"uuid"`
	ProviderName string `json:"provider_name" gorm:"type:varchar(255);not null;uniqueIndex:provider_name_provider_id" validate:"required"`
	ProviderID   string `json:"provider_id" gorm:"type:varchar(255);not null;uniqueIndex:provider_name_provider_id" validate:"required"`
	EaticleID    string `json:"eaticle_id" gorm:"type:varchar(255);unique;not null" validate:"required"`
	UserName     string `json:"user_name" gorm:"type:varchar(255);not null" validate:"required"`
	UserImg      string `json:"user_img" gorm:"type:varchar(255)" validate:"omitempty,url"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Articles     []Article `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"articles"`
}

func ValidateUser(user *User) error {
	err := Validate.Struct(user)
	if err != nil {
		// バリデーションエラーが発生したフィールドをログに記録
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("Validation failed for field '%s' with tag '%s'", err.Field(), err.Tag())
		}
		return err
	}
	return nil
}
