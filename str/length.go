package str

// Length 字符串的字符长度
func Length(s string) int {
	return len([]rune(s))
}
