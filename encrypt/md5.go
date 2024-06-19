package encrypt

import (
	"crypto/md5"
	"fmt"
)

// Md5 对字符串做MD5加密（32位）
// str:要加密的字符串
// return:加密后的字符串
func Md5(str string) string {
	sum := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", sum) //将[]byte转成16进制
}
