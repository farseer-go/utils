package test

import (
	"github.com/farseer-go/utils/file"
	"github.com/farseer-go/utils/str"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

// GetFiles
func TestGetFiles(t *testing.T) {
	path := "Farseer.Go"
	files := file.GetFiles(path, "*.txt", true)
	assert.Equal(t, []string{"Farseer.Go/subDir/sub.txt", "Farseer.Go/test.txt"}, files)

	files = file.GetFiles(path, "*.txt", false)
	assert.Equal(t, []string{"Farseer.Go/test.txt"}, files)

}

// ClearFile
func TestClearFile(t *testing.T) {
	path := "Farseer.Go/create"
	file.CreateDir766(path)
	file.CreateDir766(path + "/1")
	assert.True(t, file.IsExists(path))
	file.ClearFile(path)
}

// IsExists
func TestIsExists(t *testing.T) {
	path := "Farseer.Go/create"
	file.CreateDir(path, 0766)
	assert.True(t, file.IsExists(path))
}

func TestOther(t *testing.T) {
	git := "https://github.com/FarseerNet/farseer.go.git"
	git = filepath.Base(git)
	git = str.CutRight(git, ".git")
	assert.Equal(t, "farseer.go", git)
}
