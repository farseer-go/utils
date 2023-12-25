package test

import (
	"github.com/farseer-go/utils/str"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLength(t *testing.T) {
	assert.Equal(t, str.Length(""), 0)
	assert.Equal(t, str.Length(" "), 1)
	assert.Equal(t, str.Length("1"), 1)
	assert.Equal(t, str.Length("a"), 1)
	assert.Equal(t, str.Length("中"), 1)
	assert.Equal(t, str.Length("中"), 1)
	//测试全角
	assert.Equal(t, str.Length("　"), 1)
	assert.Equal(t, str.Length("１"), 1)
	assert.Equal(t, str.Length("ａ"), 1)
	assert.Equal(t, str.Length("，"), 1)
	assert.Equal(t, str.Length("中"), 1)
}
