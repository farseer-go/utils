package str

import "strings"

// PadLeft 为str向左填充paddingChar，直到长度=totalWidth
func PadLeft(str string, totalWidth int, paddingChar rune) string {
	char := string(paddingChar)
	if Length(str) >= totalWidth {
		return str
	}
	return strings.Repeat(char, totalWidth-Length(str)) + str
}

// PadRight 为str向左填充paddingChar，直到长度=totalWidth
func PadRight(str string, totalWidth int, paddingChar rune) string {
	char := string(paddingChar)
	if Length(str) >= totalWidth {
		return str
	}
	return str + strings.Repeat(char, totalWidth-Length(str))
}
