package test

import (
	"github.com/farseer-go/utils/str"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToDateTime(t *testing.T) {
	data := time.Date(2023, 1, 4, 0, 2, 0, 0, time.UTC)
	s := str.ToDateTime(data)
	assert.Equal(t, "2023-01-04 00:02:00", s)
}
