package str

import "strings"

// PadLeft 为str向左填充paddingChar，直到长度=totalWidth
func PadLeft(str string, totalWidth int, paddingChar string) string {
	if Length(str) >= totalWidth {
		return str
	}
	return strings.Repeat(paddingChar, totalWidth-Length(str)) + str
}

// PadRight 为str向左填充paddingChar，直到长度=totalWidth
func PadRight(str string, totalWidth int, paddingChar string) string {
	if Length(str) >= totalWidth {
		return str
	}
	return str + strings.Repeat(paddingChar, totalWidth-Length(str))
}
