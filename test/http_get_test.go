package test

import (
	"github.com/farseer-go/utils/http"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/url"
	"testing"
)

type result struct {
	Args    string            `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	Url     string            `json:"url"`
}

func TestGetJson(t *testing.T) {
	res, err := http.GetJson[result]("https://httpbin.org/get", "", 5000)
	assert.NoError(t, err)
	expected := result{
		Args:    "",
		Headers: res.Headers,
		Origin:  res.Origin,
		Url:     "https://httpbin.org/get",
	}
	assert.Equal(t, expected, res)
}

func TestGet(t *testing.T) {
	_, statusCode, err := http.Get("https://httpbin.org/get", "", "application/json", 5000)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, statusCode)
}

func TestGetForm(t *testing.T) {
	params := url.Values{}
	params.Set("name", "zhaofan")
	params.Set("age", "23")
	_, statusCode, err := http.GetForm("https://httpbin.org/get", params.Encode(), 5000)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, statusCode)
}

func TestGetFormWithoutBody(t *testing.T) {
	_, statusCode, err := http.GetFormWithoutBody("https://httpbin.org/get", 5000)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, statusCode)
}
