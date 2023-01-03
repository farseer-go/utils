package test

import (
	"github.com/farseer-go/utils/encrypt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// CopyFolder
func TestMd5(t *testing.T) {
	s := encrypt.Md5("123")
	assert.Equal(t, s, "202cb962ac59075b964b07152d234b70")
}
