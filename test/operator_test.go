package test

import (
	"testing"

	"github.com/farseer-go/utils/operator"
	"github.com/stretchr/testify/assert"
)

// CopyFolder
func TestSum(t *testing.T) {
	num := operator.GetSum(34)
	assert.Equal(t, 7, num)
}
