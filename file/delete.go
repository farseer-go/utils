package file

import "os"

// Delete 删除文件
// filePath：文件路径
func Delete(filePath string) {
	os.Remove(filePath)
}
