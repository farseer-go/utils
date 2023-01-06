package test

import (
	"encoding/json"
	"github.com/farseer-go/utils/http"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/url"
	"testing"
)

func TestGetJson(t *testing.T) {
	type result struct {
		Args    string            `json:"args"`
		Headers map[string]string `json:"headers"`
		Origin  string            `json:"origin"`
		Url     string            `json:"url"`
		Json    string            `json:"json"`
	}

	res, err := http.GetJson[result]("https://httpbin.org/get", nil, 5000)
	assert.NoError(t, err)
	expected := result{
		Args:    "",
		Headers: res.Headers,
		Origin:  res.Origin,
		Url:     res.Url,
	}
	assert.Equal(t, expected, res)
}

func TestGet(t *testing.T) {
	params := url.Values{}
	params.Set("name", "zhaofan")
	params.Set("age", "23")
	res, statusCode, err := http.Get("https://httpbin.org/get", params, 5000)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, statusCode)
	type result struct {
		Args map[string]string `json:"args"`
	}
	var val = result{}
	err = json.Unmarshal([]byte(res), &val)
	assert.NoError(t, err)
	assert.Equal(t, result{
		Args: map[string]string{"age": "23", "name": "zhaofan"},
	}, val)

}

func TestGetFormWithoutBody(t *testing.T) {
	_, statusCode, err := http.GetFormWithoutBody("https://httpbin.org/get", 5000)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, statusCode)
}
