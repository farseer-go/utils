package test

import (
	"testing"

	"github.com/farseer-go/utils/str"
	"github.com/stretchr/testify/assert"
)

func TestPadRight(t *testing.T) {
	assert.Equal(t, str.PadRight("1", 5, "0"), "10000")
	assert.Equal(t, str.PadRight("哈哈", 5, "哼"), "哈哈哼哼哼")
}

func TestPadLeft(t *testing.T) {
	assert.Equal(t, str.PadLeft("1", 5, "0"), "00001")
	assert.Equal(t, str.PadLeft("哈哈", 5, "哼"), "哼哼哼哈哈")
}

func TestRandInt64(t *testing.T) {
	assert.Equal(t, len(str.RandInt64(999999999)), 9)
}
