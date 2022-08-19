package str

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCutRight(t *testing.T) {
	assert.Equal(t, CutRight("aaaacbb", "bb"), "aaaac")
}
