package functions

import (
	"net/http"

	"github.com/U-T-kuroitigo/eaticle_go_server/configuration"
	"github.com/U-T-kuroitigo/eaticle_go_server/models"
	"github.com/labstack/echo"
)

// articleテーブルへの追加処理
func CreateArticle(c echo.Context) error {
	article := new(models.Article)
	if err := c.Bind(article); err != nil {
		return HandleInvalidRequestBody(c, err)
	}

	if err := models.ValidateArticle(article); err != nil {
		return HandleInvalidRequestBody(c, err) // バリデーションエラーも400で処理
	}

	db := configuration.GetDB()
	if err := db.Create(&article).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Created successfully", article, http.StatusCreated)
}

// articleテーブルの全件取得処理
func GetAllArticles(c echo.Context) error {
	articles := []models.Article{}
	db := configuration.GetDB()
	if err := db.Find(&articles).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Successfully retrieved articles", articles, http.StatusOK)
}

// articleテーブルの一件取得処理
func GetArticle(c echo.Context) error {
	ai := c.QueryParam("article_id")
	db := configuration.GetDB()

	var article models.Article
	if err := db.Where("article_id = ?", ai).First(&article).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Successfully retrieved article", article, http.StatusOK)
}

// articleテーブルの更新処理
func UpdateArticle(c echo.Context) error {
	ai := c.QueryParam("article_id")
	db := configuration.GetDB()

	var article models.Article
	if err := db.Where("article_id = ?", ai).First(&article).Error; err != nil {
		return HandleDBError(c, err)
	}

	var requestBody map[string]interface{}
	if err := c.Bind(&requestBody); err != nil {
		return HandleInvalidRequestBody(c, err)
	}

	// 許可されたフィールドのみ更新
	allowedUpdates := map[string]bool{
		"article_thumbnail_path": true,
		"article_title":          true,
		"article_body":           true,
		"public":                 true,
	}
	updates := FilterAllowedFields(requestBody, allowedUpdates)

	if err := db.Model(&models.Article{}).Where("article_id = ?", ai).Updates(updates).Error; err != nil {
		return HandleDBError(c, err)
	}

	// 更新後のデータを取得して返却
	if err := db.Where("article_id = ?", ai).First(&article).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Article updated successfully", article, http.StatusAccepted)
}

// articleテーブルの削除処理
func DeleteArticle(c echo.Context) error {
	ai := c.QueryParam("article_id")
	db := configuration.GetDB()

	var article models.Article
	if err := db.Where("article_id = ?", ai).First(&article).Error; err != nil {
		return HandleDBError(c, err)
	}

	if err := db.Delete(&article).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Article deleted successfully", article, http.StatusAccepted)
}
