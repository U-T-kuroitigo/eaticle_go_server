package functions

import (
	"net/http"

	"github.com/U-T-kuroitigo/eaticle_go_server/configuration"
	"github.com/U-T-kuroitigo/eaticle_go_server/models"
	"github.com/labstack/echo"
)

// article_tagテーブルへの追加処理
func CreateArticleTag(c echo.Context) error {
	articleTag := new(models.ArticleTag)
	if err := c.Bind(articleTag); err != nil {
		return HandleInvalidRequestBody(c, err)
	}

	if err := models.ValidateArticleTag(articleTag); err != nil {
		return HandleInvalidRequestBody(c, err) // バリデーションエラーも400で処理
	}

	db := configuration.GetDB()
	if err := db.Create(&articleTag).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Article tag created successfully", articleTag, http.StatusCreated)
}

// article_tagテーブルの全件取得処理
func GetAllArticleTags(c echo.Context) error {
	articleTags := []models.ArticleTag{}
	db := configuration.GetDB()
	if err := db.Find(&articleTags).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Successfully retrieved article tags", articleTags, http.StatusOK)
}

// article_tagテーブルの一件取得処理
func GetArticleTag(c echo.Context) error {
	ati := c.QueryParam("article_tag_id")
	db := configuration.GetDB()

	var articleTag models.ArticleTag
	if err := db.Where("article_tag_id = ?", ati).First(&articleTag).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Successfully retrieved article tag", articleTag, http.StatusOK)
}

// article_tagテーブルの更新処理
func UpdateArticleTag(c echo.Context) error {
	ati := c.QueryParam("article_tag_id")
	db := configuration.GetDB()

	var articleTag models.ArticleTag
	if err := db.Where("article_tag_id = ?", ati).First(&articleTag).Error; err != nil {
		return HandleDBError(c, err)
	}

	var requestBody map[string]interface{}
	if err := c.Bind(&requestBody); err != nil {
		return HandleInvalidRequestBody(c, err)
	}

	// 許可されたフィールドのみ更新
	allowedUpdates := map[string]bool{
		"article_tag_name": true,
	}
	updates := FilterAllowedFields(requestBody, allowedUpdates)

	if err := db.Model(&models.ArticleTag{}).Where("article_tag_id = ?", ati).Updates(updates).Error; err != nil {
		return HandleDBError(c, err)
	}

	// 更新後のデータを取得して返却
	if err := db.Where("article_tag_id = ?", ati).First(&articleTag).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Article tag updated successfully", articleTag, http.StatusAccepted)
}

// article_tagテーブルの削除処理
func DeleteArticleTag(c echo.Context) error {
	ati := c.QueryParam("article_tag_id")
	db := configuration.GetDB()

	var articleTag models.ArticleTag
	if err := db.Where("article_tag_id = ?", ati).First(&articleTag).Error; err != nil {
		return HandleDBError(c, err)
	}

	if err := db.Delete(&articleTag).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Article tag deleted successfully", articleTag, http.StatusAccepted)
}
