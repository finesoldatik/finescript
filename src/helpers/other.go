package helpers

// Функция для обрезания длинных строк
func Ellipsis(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
