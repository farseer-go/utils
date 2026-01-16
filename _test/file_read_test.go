package test

import (
	selfFile "github.com/farseer-go/utils/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadString(t *testing.T) {
	file := "Farseer.Go/test.txt"
	expected := `test1
test2`
	assert.Equal(t, expected, selfFile.ReadString(file))
}

func TestReadAllLines(t *testing.T) {
	file := "Farseer.Go/test.txt"
	expected := []string{
		"test1",
		"test2",
	}
	assert.Equal(t, expected, selfFile.ReadAllLines(file))
}
