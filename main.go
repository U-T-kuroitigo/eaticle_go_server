package main

import (
	"fmt"

	"github.com/U-T-kuroitigo/eaticle_go_server/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
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
