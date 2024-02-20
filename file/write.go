package file

import (
	"bufio"
	"os"
	"strings"
)

// WriteString 写入文件
// filePath：文件路径
// content：文件内容
func WriteString(filePath string, content string) {
	_ = os.WriteFile(filePath, []byte(content), 0766)
}

// WriteByte 写入文件
// filePath：文件路径
// content：文件内容
func WriteByte(filePath string, content []byte) {
	_ = os.WriteFile(filePath, content, 0766)
}

// AppendString 追加文件
// filePath：文件路径
// content：文件内容
func AppendString(filePath string, content string) {
	// 写日志到文件
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	// 写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	_, _ = write.WriteString(content)
	//Flush将缓存的文件真正写入到文件中
	_ = write.Flush()
	//及时关闭file句柄
	_ = file.Close()
}

// AppendLine 换行追加文件
// filePath：文件路径
// content：文件内容
func AppendLine(filePath string, content string) {
	AppendString(filePath, content+"\n")
}

// AppendAllLine 换行追加文件
// filePath：文件路径
// contents：文件内容
func AppendAllLine(filePath string, contents []string) {
	AppendString(filePath, strings.Join(contents, "\n")+"\n")
}
