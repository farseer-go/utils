package test

import (
	"github.com/farseer-go/utils/encrypt"
	"github.com/stretchr/testify/assert"
	"testing"
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
