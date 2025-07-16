package helpers

import "strconv"

/*
Возвращает укороченную до заданной длины искомую строку
*/
func Ellipsis(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

/*
Обработка escape-последовательностей в строках

Примеры использования:

- Ellipsis("Hello, world!", 5) → "Hello..."

- RemoveEscapeSigns("\\n") → "\n"

- RemoveEscapeSigns("\\u263A") → "☺"
*/
func RemoveEscapeSigns(escapeSeq string) string {
	if len(escapeSeq) < 2 || escapeSeq[0] != '\\' {
		return escapeSeq
	}

	switch escapeSeq[1] {
	case 'n':
		return "\n"
	case 't':
		return "\t"
	case 'r':
		return "\r"
	case '\\':
		return "\\"
	case '"':
		return "\""
	case '\'':
		return "'"
	case 'x':
		// Обработка шестнадцатеричных escape-последовательностей (\xXX)
		if len(escapeSeq) == 4 { // \x + 2 символа
			if val, err := strconv.ParseUint(escapeSeq[2:], 16, 8); err == nil {
				return string(byte(val))
			}
		}
	case 'u':
		// Обработка Unicode-последовательностей (\uXXXX)
		if len(escapeSeq) == 6 { // \u + 4 символа
			if val, err := strconv.ParseUint(escapeSeq[2:], 16, 16); err == nil {
				return string(rune(val))
			}
		}
	case 'U':
		// Обработка расширенных Unicode-последовательностей (\UXXXXXXXX)
		if len(escapeSeq) == 10 { // \U + 8 символов
			if val, err := strconv.ParseUint(escapeSeq[2:], 16, 32); err == nil {
				return string(rune(val))
			}
		}
	}

	return escapeSeq
}
