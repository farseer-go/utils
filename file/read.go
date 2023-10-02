package file

import (
	"github.com/farseer-go/fs"
	"os"
	"strings"
)

// ReadString 读文件内容
// filePath：文件路径
func ReadString(filePath string) string {
	file, _ := os.ReadFile(filePath)
	return string(file)
}

// ReadAllLines 读文件内容，按行返回数组
// filePath：文件路径
func ReadAllLines(filePath string) []string {
	file, _ := os.ReadFile(filePath)
	return strings.Split(string(file), fs.Newline)
}
