package str

import "strings"

// PadLeft 为str向左填充paddingChar，直到长度=totalWidth
func PadLeft(str string, totalWidth int, paddingChar string) string {
	if len(str) >= totalWidth {
		return str
	}
	return strings.Repeat(paddingChar, totalWidth-len(str)) + str
}
