package functions

import (
	"net/http"

	"github.com/U-T-kuroitigo/eaticle_go_server/configuration"
	"github.com/U-T-kuroitigo/eaticle_go_server/models"
	"github.com/labstack/echo"
)

// userテーブルへの追加処理
func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return HandleInvalidRequestBody(c, err)
	}

	if err := models.ValidateUser(user); err != nil {
		return HandleInvalidRequestBody(c, err) // バリデーションエラーも400で処理
	}

	db := configuration.GetDB()
	if err := db.Create(&user).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Created successfully", user, http.StatusCreated)
}

// userテーブルの全件取得処理
func GetAllUsers(c echo.Context) error {
	users := []models.User{}
	db := configuration.GetDB()
	if err := db.Find(&users).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Successfully retrieved users", users, http.StatusOK)
}

// userテーブルの一件取得処理
func GetUser(c echo.Context) error {
	ui := c.QueryParam("user_id")
	db := configuration.GetDB()

	var user models.User
	if err := db.Where("user_id = ?", ui).First(&user).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "Successfully retrieved user", user, http.StatusOK)
}

// userテーブルの更新処理
func UpdateUser(c echo.Context) error {
	ui := c.QueryParam("user_id")
	db := configuration.GetDB()

	var user models.User
	if err := db.Where("user_id = ?", ui).First(&user).Error; err != nil {
		return HandleDBError(c, err)
	}

	var requestBody map[string]interface{}
	if err := c.Bind(&requestBody); err != nil {
		return HandleInvalidRequestBody(c, err)
	}

	// 許可されたフィールドのみ更新
	allowedUpdates := map[string]bool{
		"provider_name": true,
		"provider_id":   true,
		"eaticle_id":    true,
		"user_name":     true,
		"user_img":      true,
	}
	updates := FilterAllowedFields(requestBody, allowedUpdates)

	if err := db.Model(&models.User{}).Where("user_id = ?", ui).Updates(updates).Error; err != nil {
		return HandleDBError(c, err)
	}

	// 更新後のデータを取得して返却
	if err := db.Where("user_id = ?", ui).First(&user).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "User updated successfully", user, http.StatusAccepted)
}

// userテーブルの削除処理
func DeleteUser(c echo.Context) error {
	ui := c.QueryParam("user_id")
	db := configuration.GetDB()

	var user models.User
	if err := db.Where("user_id = ?", ui).First(&user).Error; err != nil {
		return HandleDBError(c, err)
	}

	if err := db.Delete(&user).Error; err != nil {
		return HandleDBError(c, err)
	}

	return HandleSuccess(c, "User deleted successfully", user, http.StatusAccepted)
}
