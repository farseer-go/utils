package encrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

// Md5 对字符串做MD5加密（32位）
// str:要加密的字符串
// return:加密后的字符串
func Md5(str string) string {
	sum := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", sum) //将[]byte转成16进制
}

// Sha1 加密码
func Sha1(str string) string {
	// The pattern for generating a hash is `sha1.New()`,
	// `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
	// Here we start with a new hash.
	h := sha1.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	h.Write([]byte(str))

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	bs := h.Sum(nil)

	// SHA1 values are often printed in hex, for example
	// in git commits. Use the `%x` format verb to convert
	// a hash results to a hex string.
	return fmt.Sprintf("%x", bs)
}
