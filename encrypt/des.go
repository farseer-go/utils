package encrypt

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
)

// 加密函数
func DESEncrypt(plainText string, key string) string {
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	// PKCS7 填充
	paddedText := pad([]byte(plainText), block.BlockSize())
	cipherText := make([]byte, len(paddedText))

	// 使用 ECB 模式
	for i := 0; i < len(paddedText); i += block.BlockSize() {
		block.Encrypt(cipherText[i:i+block.BlockSize()], paddedText[i:i+block.BlockSize()])
	}

	// 返回 Base64 编码的密文
	return base64.StdEncoding.EncodeToString(cipherText)
}

// 解密函数
func DESDecrypt(cipherText string, key string) string {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return ""
	}

	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	if len(ciphertext)%block.BlockSize() != 0 {
		//return "", errors.New("ciphertext is not a multiple of the block size")
		return ""
	}

	plainText := make([]byte, len(ciphertext))

	// 使用 ECB 模式解密
	for i := 0; i < len(ciphertext); i += block.BlockSize() {
		block.Decrypt(plainText[i:i+block.BlockSize()], ciphertext[i:i+block.BlockSize()])
	}

	// 去掉填充
	return string(unpad(plainText))
}

// PKCS7 填充
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// 去除填充
func unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
