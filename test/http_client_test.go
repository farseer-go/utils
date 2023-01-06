package test

import (
	"github.com/farseer-go/utils/http"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func Test_client_Post(t *testing.T) {
	_, _ = http.NewClient("https://httpbin.org/get").Body(nil).Head(nil).Post()
}

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
	_, err := http.Get("https://httpbin.org/get", "", "application/json", 5000)
	assert.NoError(t, err)
}

func TestGetForm(t *testing.T) {
	params := url.Values{}
	params.Set("name", "zhaofan")
	params.Set("age", "23")
	_, err := http.GetForm("https://httpbin.org/get", params.Encode(), 5000)
	assert.NoError(t, err)
}

func TestGetFormWithoutBody(t *testing.T) {
	res, err := http.GetFormWithoutBody("https://httpbin.org/get", 5000)
	assert.NoError(t, err)
	t.Log(res)
}
