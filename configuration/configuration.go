package configuration

import (
	"fmt"
	"log"
	"os"

	"github.com/U-T-kuroitigo/eaticle_go_server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// データベース接続を保持するグローバル変数
var db *gorm.DB

// Configuration構造体: 環境変数を保持するための構造体
type Configuration struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
}

// 環境変数を取得してConfiguration構造体に格納する関数
func GetConfiguration() Configuration {
	var c Configuration

	// 環境変数を取得
	c.Server = os.Getenv("Server")
	c.Port = os.Getenv("Port")
	c.User = os.Getenv("User")
	c.Password = os.Getenv("Password")
	c.Database = os.Getenv("Database")

	return c
}

// データベース接続を初期化し、マイグレーションを実行する関数
func InitDB() {
	var err error
	c := GetConfiguration()

	// データベース接続文字列（DSN）を作成
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Server, c.Port, c.User, c.Password, c.Database)

	// データベース接続を初期化
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// データベースへの接続に失敗しました
		log.Fatal("Failed to connect to database: ", err)
	}

	// 接続確認を実行
	sqlDB, err := db.DB()
	if err != nil {
		// データベース初期化中にエラーが発生しました
		log.Fatal("Failed to initialize database: ", err)
	}
	if err := sqlDB.Ping(); err != nil {
		// データベースへの接続確認に失敗しました
		log.Fatal("Failed to ping database: ", err)
	}

	// モデルのマイグレーションを実行
	if err := db.AutoMigrate(&models.User{}, &models.Article{}, &models.ArticleTag{}); err != nil {
		// モデルのマイグレーションに失敗しました
		log.Fatal("Failed to auto-migrate models: ", err)
	}
}

// 初期化されたデータベース接続を取得する関数
func GetDB() *gorm.DB {
	// データベース接続が初期化されていない場合はエラーを表示して終了
	if db == nil {
		// データベース接続が初期化されていません
		log.Fatal("Database connection is not initialized. Call InitDB first.")
	}
	return db
}
