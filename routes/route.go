package routes

import (
	"github.com/U-T-kuroitigo/eaticle/tables/user"
	"github.com/labstack/echo"
)

func userRoutes(e *echo.Echo) {
	e.GET("api/v2/users", user.GetAll)   //GetAll Users
	e.GET("api/v2/user", user.Get)       //GET one user
	e.POST("api/v2/user", user.Create)   //CREATE
	e.PUT("api/v2/user", user.Update)    //UPDATE
	e.DELETE("api/v2/user", user.Delete) //DELETE
}

func StartRoutes(e *echo.Echo) {
	userRoutes(e)
}
