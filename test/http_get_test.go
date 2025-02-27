package test

import (
	"net/url"
	"testing"

	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/utils/http"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestGetJson(t *testing.T) {
	type result struct {
		Args    string            `json:"args"`
		Headers map[string]string `json:"headers"`
		Origin  string            `json:"origin"`
		Url     string            `json:"url"`
		Json    string            `json:"json"`
	}
	params := url.Values{}
	params.Set("name", "zhaofan")
	params.Set("age", "23")
	res, err := http.GetJson[result]("https://httpbin.org/get", params, 5000)
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
	t.Run("", func(t *testing.T) {
		res, statusCode, err := http.Get("https://httpbin.org/get", params, 5000)
		assert.NoError(t, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)
		type result struct {
			Args map[string]string `json:"args"`
		}
		var val = result{}
		err = snc.Unmarshal([]byte(res), &val)
		assert.NoError(t, err)
		assert.Equal(t, result{
			Args: map[string]string{"age": "23", "name": "zhaofan"},
		}, val)
	})

	t.Run("", func(t *testing.T) {
		res, statusCode, err := http.Get("https://httpbin.org/get?p=123", params, 5000)
		assert.NoError(t, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)
		type result struct {
			Args map[string]string `json:"args"`
		}
		var val = result{}
		err = snc.Unmarshal([]byte(res), &val)
		assert.NoError(t, err)
		assert.Equal(t, result{
			Args: map[string]string{"age": "23", "name": "zhaofan", "p": "123"},
		}, val)
	})

}

func TestGetFormWithoutBody(t *testing.T) {
	_, statusCode, err := http.GetFormWithoutBody("https://httpbin.org/get", 5000)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, statusCode)
}
