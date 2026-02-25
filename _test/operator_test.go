package test

import (
	"errors"
	"testing"

	"github.com/farseer-go/utils/operator"
	"github.com/stretchr/testify/assert"
)

// CopyFolder
func TestSum(t *testing.T) {
	num := operator.GetSum(34)
	assert.Equal(t, 7, num)
}

func TestGetNotNil(t *testing.T) {
	var err1, err2, err3 error
	err3 = errors.New("test")
	err := operator.GetNotNil(err1, err2, err3)
	assert.Equal(t, "test", err.Error())
}

func TestExistsFalse(t *testing.T) {
	val := operator.ExistsFalse(true, true, false)
	assert.True(t, val)
}

func TestExistsTrue(t *testing.T) {
	val := operator.ExistsTrue(true, true, false)
	assert.True(t, val)
}
