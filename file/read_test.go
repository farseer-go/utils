package file

import (
	"github.com/farseer-go/fs/flog"
	"testing"
)

func TestReadString(t *testing.T) {
	file := "/Users/steden/Desktop/code/project/Farseer.Go/go.mod"
	flog.Println(ReadString(file))
}
