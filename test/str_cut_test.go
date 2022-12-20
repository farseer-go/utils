package test

import (
	"github.com/farseer-go/utils/str"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCutRight(t *testing.T) {
	assert.Equal(t, str.CutRight("aaaacbb", "bb"), "aaaac")
}
