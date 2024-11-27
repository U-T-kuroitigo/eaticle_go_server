package models

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

type ArticleTag struct {
	ArticleTagID   string `json:"article_tag_id" gorm:"type:varchar(255);primaryKey;not null" validate:"uuid"`
	ArticleID      string `json:"article_id" gorm:"type:varchar(255);not null" validate:"required,uuid"`
	ArticleTagName string `json:"article_tag_name" gorm:"type:varchar(255);not null" validate:"required"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func ValidateArticleTag(articleTag *ArticleTag) error {
	err := Validate.Struct(articleTag)
	if err != nil {
		// バリデーションエラーが発生したフィールドをログに記録
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("Validation failed for field '%s' with tag '%s'", err.Field(), err.Tag())
		}
		return err
	}
	return nil
}
