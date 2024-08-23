package test

import (
	"github.com/farseer-go/utils/http"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpDownload(t *testing.T) {
	_, err := http.Download("https://github.com/farseers/FOPS-Actions/releases/download/v1/gitProxy", "./gitProxy", nil, 0, "socks5://127.0.0.1:7890")
	assert.Nil(t, err)
}
