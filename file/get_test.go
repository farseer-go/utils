package file

import (
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/str"
	"path/filepath"
	"testing"
)

// GetFiles
func TestGetFiles(t *testing.T) {
	path := "/Users/steden/Desktop/code/project/Farseer.Go"
	files := GetFiles(path, "*.md", true)
	for _, file := range files {
		flog.Println(file)
	}
}

// ClearFile
func TestClearFile(t *testing.T) {
	path := "/Users/steden/Desktop/code/project/Farseer.Go2"
	ClearFile(path)
}

// IsExists
func TestIsExists(t *testing.T) {
	path := "/Users/steden/Desktop/code/project/Farseer.Go3"
	flog.Println(IsExists(path))
}

func TestOther(t *testing.T) {
	git := "https://github.com/FarseerNet/farseer.go.git"
	git = filepath.Base(git)
	git = str.CutRight(git, ".git")
	flog.Println(git)
}
