package common

import (
	"log"

	"github.com/gofrs/uuid"
)

// GenerateUUIDv7: UUID v7を生成する関数
func GenerateUUIDv7() string {
	// UUID v7の生成
	newUUID, err := uuid.NewV7()
	if err != nil {
		// 生成中にエラーが発生した場合、ログに出力
		log.Fatalf("Failed to generate UUID v7: %v", err)
	}
	return newUUID.String()
}

// FilterAllowedFields: 許可されたフィールドのみをフィルタリングする関数
// input: クライアントから送られたリクエストボディ
// allowed: 許可されたフィールドのマップ
// 戻り値: 許可されたフィールドのみを含むマップ
func FilterAllowedFields(input map[string]interface{}, allowed map[string]bool) map[string]interface{} {
	filtered := make(map[string]interface{})
	// 入力フィールドを確認し、許可されたフィールドのみ追加
	for key, value := range input {
		if allowed[key] {
			filtered[key] = value
		}
	}
	return filtered
}
