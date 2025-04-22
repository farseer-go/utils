package test

import (
	"testing"

	"github.com/farseer-go/utils/system"
	"github.com/stretchr/testify/assert"
)

func TestResourceResource(t *testing.T) {
	resource := system.GetResource("/home", "/")
	assert.Greater(t, resource.CpuCores, 0)
	resource.ToString()
}
