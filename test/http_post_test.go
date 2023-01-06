package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farseer-go/utils/http"
	"github.com/stretchr/testify/assert"
	"io"
	goHttp "net/http"
	"testing"
)

func TestPost(t *testing.T) {
	data := make(map[string]interface{})
	data["name"] = "zhaofan"
	data["age"] = "23"
	res, err := http.Post("https://httpbin.org/post", nil, data, "application/json", 5000)
	assert.NoError(t, err)
	var val = make(map[string]interface{}, 0)
	err = json.Unmarshal([]byte(res), &val)
	assert.NoError(t, err)
	assert.Equal(t, data, val["json"])
}

func TestPostForm(t *testing.T) {
	data := make(map[string]interface{})
	data["name"] = "zhaofan"
	data["age"] = "23"
	res, err := http.PostForm("https://httpbin.org/post", nil, data, 5000)
	assert.NoError(t, err)
	var val = make(map[string]interface{}, 0)
	err = json.Unmarshal([]byte(res), &val)
	assert.NoError(t, err)
	t.Log(val)
	t.Logf("%+v", data)
	t.Logf("%+v", val["form"])
	// for k, v := range val["form"] {
	// 	t.Log(k, v)
	// }
	// assert.Equal(t, data, val["form"])
}

func TestA(t *testing.T) {
	client := &goHttp.Client{}
	data := make(map[string]interface{})
	data["name"] = "zhaofan"
	data["age"] = "23"
	bytesData, _ := json.Marshal(data)
	req, _ := goHttp.NewRequest("POST", "http://httpbin.org/post", bytes.NewReader(bytesData))
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
