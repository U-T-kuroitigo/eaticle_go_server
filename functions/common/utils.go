package common

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
