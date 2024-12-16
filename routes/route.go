package routes

import (
	"github.com/U-T-kuroitigo/eaticle_go_server/functions/api"  // API用関数
	"github.com/U-T-kuroitigo/eaticle_go_server/functions/crud" // CRUD用関数
	"github.com/labstack/echo"
)

// ユーザー関連CRUDルートの定義
func userCRUDRoutes(e *echo.Echo) {
	e.POST("crud/v2/user", crud.CreateUser)   // ユーザー作成
	e.GET("crud/v2/users", crud.GetAllUsers)  // 全ユーザー取得
	e.GET("crud/v2/user", crud.GetUser)       // 特定ユーザー取得
	e.PUT("crud/v2/user", crud.UpdateUser)    // ユーザー更新
	e.DELETE("crud/v2/user", crud.DeleteUser) // ユーザー削除
}

// 記事関連CRUDルートの定義
func articleCRUDRoutes(e *echo.Echo) {
	e.POST("crud/v2/article", crud.CreateArticle)   // 記事作成
	e.GET("crud/v2/articles", crud.GetAllArticles)  // 全記事取得
	e.GET("crud/v2/article", crud.GetArticle)       // 特定記事取得
	e.PUT("crud/v2/article", crud.UpdateArticle)    // 記事更新
	e.DELETE("crud/v2/article", crud.DeleteArticle) // 記事削除
}

// 記事タグ関連CRUDルートの定義
func articleTagCRUDRoutes(e *echo.Echo) {
	e.POST("crud/v2/article_tag", crud.CreateArticleTag)   // 記事タグ作成
	e.GET("crud/v2/article_tags", crud.GetAllArticleTags)  // 全記事タグ取得
	e.GET("crud/v2/article_tag", crud.GetArticleTag)       // 特定記事タグ取得
	e.PUT("crud/v2/article_tag", crud.UpdateArticleTag)    // 記事タグ更新
	e.DELETE("crud/v2/article_tag", crud.DeleteArticleTag) // 記事タグ削除
}

func articleAPIRoutes(e *echo.Echo) {
	e.POST("api/v2/article/save", api.SaveArticle)                   // 記事保存
	e.GET("api/v2/article/list", api.GetArticles)                    // 記事一覧取得
	e.GET("api/v2/article/:article_id/detail", api.GetArticleDetail) // 記事詳細取得
	e.DELETE("api/v2/article/:article_id/delete", api.DeleteArticle) // 記事削除
}

// ルートの初期化
func StartRoutes(e *echo.Echo) {
	userCRUDRoutes(e)       // ユーザー関連CRUDルートを登録
	articleCRUDRoutes(e)    // 記事関連CRUDルートを登録
	articleTagCRUDRoutes(e) // 記事タグ関連CRUDルートを登録

	articleAPIRoutes(e) // 記事APIルートを登録
}
