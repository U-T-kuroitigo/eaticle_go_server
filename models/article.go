package models

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Article struct {
	ArticleID            string `json:"article_id" gorm:"type:varchar(255);primaryKey;not null" validate:"uuid"`
	UserID               string `json:"user_id" gorm:"type:varchar(255);not null" validate:"required,uuid"`
	ArticleThumbnailPath string `json:"article_thumbnail_path" gorm:"type:varchar(255);not null" validate:"required,url"`
	ArticleTitle         string `json:"article_title" gorm:"type:varchar(255);not null" validate:"required"`
	ArticleBody          string `json:"article_body" gorm:"type:text;not null" validate:"required"`
	Public               bool   `json:"public" gorm:"type:boolean;not null;default:false"` // デフォルト値を設定
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

func ValidateArticle(article *Article) error {
	err := Validate.Struct(article)
	if err != nil {
		// バリデーションエラーが発生したフィールドをログに記録
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("Validation failed for field '%s' with tag '%s'", err.Field(), err.Tag())
		}
		return err
	}
	return nil
}
