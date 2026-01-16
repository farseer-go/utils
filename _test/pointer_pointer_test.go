package test

import (
	"github.com/farseer-go/utils/pointer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Of(t *testing.T) {

}

func Unwrap(t *testing.T) {
	i := 100
	ptr := &i
	assert.Equal(t, pointer.Unwrap(ptr), 100)
	assert.Equal(t, pointer.Unwrap(ptr, -1), 100)
	assert.Equal(t, pointer.Unwrap[int](nil), 0)
	assert.Equal(t, pointer.Unwrap[int](nil, -1), -1)
}

func Extract(t *testing.T) {
	i := 100
	p1 := &i
	p2 := &p1
	p3 := &p2
	assert.Equal(t, pointer.Extract(p3), 100)
	assert.Equal(t, pointer.Extract(nil), nil)
}
