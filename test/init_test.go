package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/modules"
)

func init() {
	fs.Initialize[modules.FarseerKernelModule]("unit test")
}
