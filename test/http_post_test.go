package test

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/farseer-go/utils/http"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestPost(t *testing.T) {
	data := make(map[string]string)
	data["name"] = "zhaofan"
	data["age"] = "23"
	type result struct {
		Json map[string]string `json:"json"`
		Form map[string]string `json:"form"`
	}
	t.Run("json", func(t *testing.T) {
		res, statusCode, err := http.Post("https://httpbin.org/post", nil, data, "application/json", 5000)
		assert.NoError(t, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)

		var val = result{}
		err = sonic.Unmarshal([]byte(res), &val)
		assert.NoError(t, err)
		assert.Equal(t, data, val.Json)
	})

	t.Run("form", func(t *testing.T) {
		res, statusCode, err := http.Post("https://httpbin.org/post", nil, data, "application/x-www-form-urlencoded", 5000)
		assert.NoError(t, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)
		var val = result{}
		err = sonic.Unmarshal([]byte(res), &val)
		assert.NoError(t, err)
		assert.Equal(t, data, val.Form)
	})

}

func TestPostForm(t *testing.T) {
	data := make(map[string]any)
	data["name"] = "zhaofan"
	data["age"] = "23"
	res, statusCode, err := http.PostForm("https://httpbin.org/post", nil, data, 5000)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, statusCode)
	var val = make(map[string]any, 0)
	err = sonic.Unmarshal([]byte(res), &val)
	assert.NoError(t, err)
	assert.Equal(t, data, val["form"])
}

func TestPostFormWithoutBody(t *testing.T) {
	res, statusCode, err := http.PostFormWithoutBody("https://httpbin.org/post", nil, 5000)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, statusCode)
	var val = make(map[string]any, 0)
	err = sonic.Unmarshal([]byte(res), &val)
	assert.NoError(t, err)
}

func TestPostJson(t *testing.T) {
	data := make(map[string]any)
	data["name"] = "zhaofan"
	data["age"] = "23"

	type result struct {
		Headers map[string]string `json:"headers"`
		Origin  string            `json:"origin"`
		Url     string            `json:"url"`
		Json    map[string]any    `json:"json"`
	}
	res, statusCode, err := http.PostJson[result]("https://httpbin.org/post", nil, data, 5000)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, statusCode)
	assert.Equal(t, data, res.Json)

}
