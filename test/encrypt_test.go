package test

import (
	"testing"

	"github.com/farseer-go/utils/encrypt"
	"github.com/stretchr/testify/assert"
)

// CopyFolder
func TestMd5(t *testing.T) {
	s := encrypt.Md5("123")
	assert.Equal(t, "202cb962ac59075b964b07152d234b70", s)
}

// CopyFolder
func TestSha1(t *testing.T) {
	s := encrypt.Sha1("abc")
	assert.Equal(t, "a9993e364706816aba3e25717850c26c9cd0d89d", s)
}

// CopyFolder
func TestAES(t *testing.T) {
	// 加解密
	s := encrypt.AesEncryptByECB("hello world", "sqNWnLg7Y7dQbH6Y")
	content := encrypt.AesDecryptByECB(s, "sqNWnLg7Y7dQbH6Y")
	assert.Equal(t, "hello world", content)
}

// CopyFolder
func TestDES(t *testing.T) {
	// 加解密
	s := encrypt.DESEncrypt("hello world", "sqNWnLg7")
	content := encrypt.DESDecrypt(s, "sqNWnLg7")
	assert.Equal(t, "hello world", content)

	content = encrypt.DESDecrypt("GB3x+HDVotbZ3aAXvAB6mQ==", "sqNWnLg7")
	assert.Equal(t, "hello world", content)
}
