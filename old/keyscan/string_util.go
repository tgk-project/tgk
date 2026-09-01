package keyscan

// ErrInvalidNumber は無効な数値形式が指定された場合のエラーです
type ErrInvalidNumber struct {
	Value string
}

func (e ErrInvalidNumber) Error() string {
	return "invalid number format: " + e.Value
}

// stringToInt は文字列を整数に変換します
func stringToInt(s string) (int, error) {
	// 符号の処理
	negative := false
	start := 0
	if len(s) > 0 && s[0] == '-' {
		negative = true
		start = 1
	} else if len(s) > 0 && s[0] == '+' {
		start = 1
	}

	// 数値変換
	var result int
	for i := start; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			digit := int(s[i] - '0')
			result = result*10 + digit
		} else {
			return 0, ErrInvalidNumber{Value: s}
		}
	}

	if negative {
		result = -result
	}

	return result, nil
}
