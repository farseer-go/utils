package str

import (
	"math/rand"
	"strings"
	"time"

	"github.com/farseer-go/fs/parse"
)

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

// RandString 随机字符串
func RandInt64(max int) string {
	val := parse.ToString(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max))
	return PadLeft(val, len(parse.ToString(max)), "0")
}
