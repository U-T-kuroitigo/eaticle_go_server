package common

import (
	"fmt"
	"net/http"

	"github.com/U-T-kuroitigo/eaticle_go_server/response"
	"github.com/jackc/pgx/v5/pgconn" // PostgreSQL用
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// データベース関連のエラーを処理する
func HandleDBError(c echo.Context, err error) error {
	var r response.Model

	// データが見つからない場合
	if err == gorm.ErrRecordNotFound {
		r = response.Model{
			Code:    "404",                // HTTPステータスコード: 404 Not Found
			Message: "Resource not found", // 資源が見つからない
			Data:    err.Error(),
		}
		return c.JSON(http.StatusNotFound, r)
	}

	// PostgreSQLの一意性制約違反 (エラーコード: 23505)
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		r = response.Model{
			Code:    "409",                       // HTTPステータスコード: 409 Conflict
			Message: "Conflict: Duplicate entry", // 衝突: 重複エントリ
			Data:    pgErr.Detail,                // エラーメッセージ
		}
		return c.JSON(http.StatusConflict, r)
	}

	// // MySQLの一意性制約違反 (エラーコード: 1062)
	// if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
	// 	r = response.Model{
	// 		Code:    "409",                       // HTTPステータスコード: 409 Conflict
	// 		Message: "Conflict: Duplicate entry", // 衝突: 重複エントリ
	// 		Data:    mysqlErr.Message,            // エラーメッセージ
	// 	}
	// 	return c.JSON(http.StatusConflict, r)
	// }

	// その他のデータベースエラー
	r = response.Model{
		Code:    "500",            // HTTPステータスコード: 500 Internal Server Error
		Message: "Database error", // データベースエラー
		Data:    err.Error(),
	}
	return c.JSON(http.StatusInternalServerError, r)
}

// リクエストボディのエラーを処理する
func HandleInvalidRequestBody(c echo.Context, err error) error {
	r := response.Model{
		Code:    "400",                  // HTTPステータスコード: 400 Bad Request
		Message: "Invalid request body", // 無効なリクエストボディ
		Data:    err.Error(),            // エラーメッセージ
	}
	return c.JSON(http.StatusBadRequest, r)
}

// 正常に操作が完了した場合のレスポンスを処理する
func HandleSuccess(c echo.Context, message string, data interface{}, code int) error {
	r := response.Model{
		Code:    fmt.Sprintf("%d", code), // HTTPステータスコード
		Message: message,                 // 操作が成功したメッセージ
		Data:    data,                    // クライアントに返すデータ（構造体や特定フィールドなど柔軟に対応）
	}
	return c.JSON(code, r)
}
