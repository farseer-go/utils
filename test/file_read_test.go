package test

import (
	"github.com/farseer-go/fs/flog"
	file2 "github.com/farseer-go/utils/file"
	"testing"
)

func TestReadString(t *testing.T) {
	file := "/Users/steden/Desktop/code/project/Farseer.Go/go.mod"
	flog.Println(file2.ReadString(file))
}
