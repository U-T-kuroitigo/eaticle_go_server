package routes

import (
	"github.com/U-T-kuroitigo/eaticle_go_server/functions"
	"github.com/labstack/echo"
)

// userエンドポイントのルート定義
func userRoutes(e *echo.Echo) {
	e.GET("api/v2/users", functions.GetAllUsers)  // 全ユーザー取得
	e.GET("api/v2/user", functions.GetUser)       // 特定ユーザー取得
	e.POST("api/v2/user", functions.CreateUser)   // ユーザー作成
	e.PUT("api/v2/user", functions.UpdateUser)    // ユーザー更新
	e.DELETE("api/v2/user", functions.DeleteUser) // ユーザー削除
}

// articleエンドポイントのルート定義
func articleRoutes(e *echo.Echo) {
	e.GET("api/v2/articles", functions.GetAllArticles)  // 全記事取得
	e.GET("api/v2/article", functions.GetArticle)       // 特定記事取得
	e.POST("api/v2/article", functions.CreateArticle)   // 記事作成
	e.PUT("api/v2/article", functions.UpdateArticle)    // 記事更新
	e.DELETE("api/v2/article", functions.DeleteArticle) // 記事削除
}

// ルートを初期化
func StartRoutes(e *echo.Echo) {
	userRoutes(e)    // ユーザー関連ルートを設定
	articleRoutes(e) // 記事関連ルートを設定
}
