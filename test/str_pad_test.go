package test

import (
	"github.com/farseer-go/utils/str"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPadRight(t *testing.T) {
	assert.Equal(t, str.PadRight("1", 5, []rune("0")[0]), "10000")
	assert.Equal(t, str.PadRight("哈哈", 5, []rune("哼")[0]), "哈哈哼哼哼")
}

func TestPadLeft(t *testing.T) {
	assert.Equal(t, str.PadLeft("1", 5, []rune("0")[0]), "00001")
	assert.Equal(t, str.PadLeft("哈哈", 5, []rune("哼")[0]), "哼哼哼哈哈")
}
