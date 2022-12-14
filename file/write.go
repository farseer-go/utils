package file

import (
	"os"
	"strings"
)

// WriteString 写入文件
// filePath：文件路径
// content：文件内容
func WriteString(filePath string, content string) {
	_ = os.WriteFile(filePath, []byte(content), 0766)
}

// AppendString 追加文件
// filePath：文件路径
// content：文件内容
func AppendString(filePath string, content string) {
	oldContent := ReadString(filePath)
	_ = os.WriteFile(filePath, []byte(oldContent+content), 0766)
}

// AppendLine 换行追加文件
// filePath：文件路径
// content：文件内容
func AppendLine(filePath string, content string) {
	oldContent := ReadString(filePath)
	_ = os.WriteFile(filePath, []byte(oldContent+"\n"+content), 0766)
}

// AppendAllLine 换行追加文件
// filePath：文件路径
// contents：文件内容
func AppendAllLine(filePath string, contents []string) {
	oldContent := ReadString(filePath)
	_ = os.WriteFile(filePath, []byte(oldContent+"\n"+strings.Join(contents, "\n")), 0766)
}
