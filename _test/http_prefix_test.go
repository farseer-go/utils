package test

import (
	"github.com/farseer-go/utils/http"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddHttpPrefix(t *testing.T) {
	h := http.AddHttpPrefix("baidu.com")
	assert.Equal(t, "http://baidu.com", h)
}

func TestAddHttpsPrefix(t *testing.T) {
	h := http.AddHttpsPrefix("baidu.com")
	assert.Equal(t, "https://baidu.com", h)
}

func TestClearHttpPrefix(t *testing.T) {
	h := http.ClearHttpPrefix("https://baidu.com")
	assert.Equal(t, "baidu.com", h)

	h = http.ClearHttpPrefix("http://baidu.com")
	assert.Equal(t, "baidu.com", h)
}
