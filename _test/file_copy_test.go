package test

import (
	"github.com/farseer-go/utils/file"
	"testing"
)

// CopyFolder
func TestCopyFolder(t *testing.T) {
	path1 := "./Farseer.Go"
	path2 := "./Farseer.Go2"

	file.CopyFolder(path1, path2)
}
