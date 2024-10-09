package routes

import (
	"github.com/U-T-kuroitigo/eaticle_go_server/functions"
	"github.com/labstack/echo"
)

func userRoutes(e *echo.Echo) {
	e.GET("api/v2/users", functions.GetAllUsers)  // GetAll Users
	e.GET("api/v2/user", functions.GetUser)       // GET one user
	e.POST("api/v2/user", functions.CreateUser)   // CREATE
	e.PUT("api/v2/user", functions.UpdateUser)    // UPDATE
	e.DELETE("api/v2/user", functions.DeleteUser) // DELETE
}

func StartRoutes(e *echo.Echo) {
	userRoutes(e)
}
