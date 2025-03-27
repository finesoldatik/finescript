package helpers

import "strconv"

func RemoveEscapeSigns(s string) string {
	if len(s) < 2 {
		return s
	}

	switch s[1] {
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
		if len(s) == 4 {
			hex := s[2:]
			if val, err := strconv.ParseUint(hex, 16, 8); err == nil {
				return string([]byte{byte(val)})
			}
		}
	case 'u':
		// Обработка Unicode-последовательностей (\uXXXX)
		if len(s) == 6 {
			hex := s[2:]
			if val, err := strconv.ParseUint(hex, 16, 16); err == nil {
				return string(rune(val))
			}
		}
	case 'U':
		// Обработка расширенных Unicode-последовательностей (\UXXXXXXXX)
		if len(s) == 10 {
			hex := s[2:]
			if val, err := strconv.ParseUint(hex, 16, 32); err == nil {
				return string(rune(val))
			}
		}
	}

	return s
}
