package main

import (
	"fmt"

	"github.com/U-T-kuroitigo/eaticle_go_server/configuration"
	"github.com/U-T-kuroitigo/eaticle_go_server/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// データベースの初期化
	db := configuration.InitDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	routes.StartRoutes(e)

	err := e.Start(":5000")
	if err != nil {
		fmt.Printf("Error, could not run server: %v", err)
	}
}
