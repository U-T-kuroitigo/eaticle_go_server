package main

import (
	"fmt"
	"log"

	"github.com/U-T-kuroitigo/eaticle_go_server/configuration"
	"github.com/U-T-kuroitigo/eaticle_go_server/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// データベースの初期化
	configuration.InitDB()

	// データベース接続を取得
	db := configuration.GetDB()
	defer func() {
		// データベース接続を安全に閉じる
		// もしエラーが発生しても終了を阻害しない
		sqlDB, err := db.DB()
		if err != nil {
			// データベース接続を閉じる際のエラー
			log.Printf("Failed to close database connection: %v", err)
		} else {
			sqlDB.Close()
		}
	}()

	// Echoフレームワークのインスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Recover()) // パニックから復帰するためのミドルウェア
	e.Use(middleware.Logger())  // リクエスト・レスポンスのログを記録
	e.Use(middleware.CORS())    // CORSの設定

	// ルートを初期化
	routes.StartRoutes(e)

	// サーバーを起動
	err := e.Start(":5000")
	if err != nil {
		// サーバー起動エラー
		fmt.Printf("Error, could not run server: %v", err)
	}
}
