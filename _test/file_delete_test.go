package test

import (
	"github.com/farseer-go/utils/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelete(t *testing.T) {
	path := "Farseer.Go/create"
	file.CreateDir766(path)
	file.Delete(path)
	assert.False(t, file.IsExists(path))
}
