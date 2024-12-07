package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/U-T-kuroitigo/eaticle_go_server/configuration"
	"github.com/U-T-kuroitigo/eaticle_go_server/functions/common"
	"github.com/U-T-kuroitigo/eaticle_go_server/models"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// 記事タグの削除と再登録を行う関数
func processArticleTags(tx *gorm.DB, articleID string, tagNames []string) error {
	// 既存のタグを削除
	if err := tx.Where("article_id = ?", articleID).Delete(&models.ArticleTag{}).Error; err != nil {
		return err
	}

	// 新しいタグを登録
	for _, tagName := range tagNames {
		articleTag := models.ArticleTag{
			ArticleTagID:   common.GenerateUUIDv7(),
			ArticleID:      articleID,
			ArticleTagName: tagName,
		}

		// モデルのバリデーションを実行
		if err := models.ValidateArticleTag(&articleTag); err == nil {
			// バリデーションエラーがなければタグを登録
			if err := tx.Create(&articleTag).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// 記事を保存する関数
func SaveArticle(c echo.Context) error {
	// リクエストボディをパース
	requestData := struct {
		ArticleID            string   `json:"article_id"`
		UserID               string   `json:"user_id"`
		ArticleThumbnailPath string   `json:"article_thumbnail_path"`
		ArticleTitle         string   `json:"article_title"`
		ArticleBody          string   `json:"article_body"`
		Public               bool     `json:"public"`
		ArticleTagNameList   []string `json:"article_tag_name_list"`
	}{}

	if err := c.Bind(&requestData); err != nil {
		return common.HandleInvalidRequestBody(c, err)
	}

	db := configuration.GetDB()
	tx := db.Begin() // トランザクションを開始

	// トランザクションの終了を確実に行うためのdefer
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // エラー発生時にロールバック
		}
	}()

	// 記事の初期データを作成
	article := models.Article{
		ArticleID:            requestData.ArticleID,
		UserID:               requestData.UserID,
		ArticleThumbnailPath: requestData.ArticleThumbnailPath,
		ArticleTitle:         requestData.ArticleTitle,
		ArticleBody:          requestData.ArticleBody,
		Public:               requestData.Public,
	}

	// モデルのバリデーションを実行
	if err := models.ValidateArticle(&article); err != nil {
		return common.HandleInvalidRequestBody(c, err)
	}

	// 記事の存在確認と更新または作成
	if err := tx.Where("article_id = ?", article.ArticleID).First(&models.Article{}).Error; err == nil {
		// 記事が存在する場合は更新
		updates := map[string]interface{}{
			"article_thumbnail_path": article.ArticleThumbnailPath,
			"article_title":          article.ArticleTitle,
			"article_body":           article.ArticleBody,
			"public":                 article.Public,
		}
		if err := tx.Model(&models.Article{}).Where("article_id = ?", article.ArticleID).Updates(updates).Error; err != nil {
			tx.Rollback()
			return common.HandleDBError(c, err)
		}
	} else {
		// 記事が存在しない場合は新規作成
		if err := tx.Create(&article).Error; err != nil {
			tx.Rollback()
			return common.HandleDBError(c, err)
		}
	}

	// タグを削除して再登録
	if err := processArticleTags(tx, article.ArticleID, requestData.ArticleTagNameList); err != nil {
		tx.Rollback()
		return common.HandleDBError(c, err)
	}

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		return common.HandleDBError(c, err)
	}

	// 成功レスポンスを返す
	responseData := map[string]string{
		"article_id": article.ArticleID,
	}
	return common.HandleSuccess(c, "Article saved successfully", responseData, http.StatusCreated)
}

// 記事のサーチクエリ関数
func ArticleSearchQuery(query *gorm.DB, searchQuery string) *gorm.DB {
	if strings.TrimSpace(searchQuery) != "" {
		likePattern := "%" + strings.TrimSpace(searchQuery) + "%"
		return query.Where(`
			articles.article_title ILIKE ? OR
			articles.article_id IN (
				SELECT article_id FROM article_tags WHERE article_tag_name ILIKE ?
			)`,
			likePattern, likePattern,
		)
	}
	return query
}

// 記事一覧取得API
func GetArticles(c echo.Context) error {
	// デフォルトのリミット
	const DefaultLimit = 15

	// クエリパラメータの取得
	searchQuery := c.QueryParam("search")
	sortQuery := c.QueryParam("sort")
	offset := c.QueryParam("offset")

	// ソート条件の解析
	validSorts := map[string]string{
		"created_at_desc": "articles.created_at DESC",
		"created_at_asc":  "articles.created_at ASC",
		"updated_at_desc": "articles.updated_at DESC",
		"updated_at_asc":  "articles.updated_at ASC",
	}
	sortOrder := "articles.created_at DESC" // デフォルト: 作成日降順
	if sortQuery != "" && validSorts[sortQuery] != "" {
		sortOrder = validSorts[sortQuery]
	}

	// Offsetを整数に変換
	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		offsetInt = 0 // デフォルトは0
	}

	// トランザクションを開始
	db := configuration.GetDB()
	tx := db.Begin()

	// トランザクション終了時にロールバック
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ベースクエリ作成とサーチクエリ適用
	baseQuery := tx.Model(&models.Article{}).
		Select(`
			articles.article_id,
			articles.article_thumbnail_path,
			articles.article_title,
			articles.created_at,
			users.eaticle_id,
			users.user_name,
			users.user_img
		`).
		Joins("JOIN users ON articles.user_id = users.user_id").
		Where("articles.public = ?", true)
	baseQuery = ArticleSearchQuery(baseQuery, searchQuery)

	// 総件数の取得
	var totalCount int64
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		tx.Rollback()
		return common.HandleDBError(c, err)
	}

	// 記事データの取得
	var articles []map[string]interface{}
	if err := baseQuery.Order(sortOrder).Offset(offsetInt).Limit(DefaultLimit).Find(&articles).Error; err != nil {
		tx.Rollback()
		return common.HandleDBError(c, err)
	}

	// トランザクションのコミット
	if err := tx.Commit().Error; err != nil {
		return common.HandleDBError(c, err)
	}

	// レスポンスデータの構築
	responseData := map[string]interface{}{
		"article_list": articles,
		"pagination": map[string]interface{}{
			"offset":   offsetInt,
			"limit":    DefaultLimit,
			"has_more": offsetInt+DefaultLimit < int(totalCount),
		},
	}

	return common.HandleSuccess(c, "Successfully retrieved articles", responseData, http.StatusOK)
}

// 記事の詳細情報を取得する関数
func GetArticleDetail(c echo.Context) error {
	// パスパラメータからarticle_idを取得
	articleID := c.Param("article_id")
	if articleID == "" {
		return common.HandleInvalidRequestBody(c, echo.NewHTTPError(http.StatusBadRequest, "Missing article_id parameter"))
	}

	db := configuration.GetDB()
	var article models.Article
	var articleTags []models.ArticleTag
	var user models.User

	// 記事情報を取得
	if err := db.Where("article_id = ?", articleID).First(&article).Error; err != nil {
		return common.HandleDBError(c, err)
	}

	// 紐付いたタグ情報を取得
	if err := db.Where("article_id = ?", articleID).Find(&articleTags).Error; err != nil {
		return common.HandleDBError(c, err)
	}

	// 作成者情報を取得
	if err := db.Where("user_id = ?", article.UserID).First(&user).Error; err != nil {
		return common.HandleDBError(c, err)
	}

	// レスポンスデータの整形
	responseData := map[string]interface{}{
		"article_id":             article.ArticleID,
		"article_thumbnail_path": article.ArticleThumbnailPath,
		"article_title":          article.ArticleTitle,
		"created_at":             article.CreatedAt,
		"article_body":           article.ArticleBody,
		"article_tag_list": map[string]interface{}{
			"list":        articleTags,
			"total_count": len(articleTags),
		},
		"eaticle_id": user.EaticleID,
		"user_name":  user.UserName,
		"user_img":   user.UserImg,
	}

	return common.HandleSuccess(c, "Article details retrieved successfully", responseData, http.StatusOK)
}

// 記事削除用API
func DeleteArticle(c echo.Context) error {
	// パスパラメータからarticle_idを取得
	articleID := c.Param("article_id")

	// データベース接続の取得
	db := configuration.GetDB()

	// 該当記事の検索
	var article models.Article
	if err := db.Where("article_id = ?", articleID).First(&article).Error; err != nil {
		return common.HandleDBError(c, err)
	}

	// 該当記事の削除
	if err := db.Delete(&article).Error; err != nil {
		return common.HandleDBError(c, err)
	}

	// 成功レスポンスを返却
	return common.HandleSuccess(c, "Article deleted successfully", struct{}{}, http.StatusOK)
}
