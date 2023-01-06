package test

import (
	"github.com/farseer-go/utils/http"
	"testing"
)

func Test_client_Post(t *testing.T) {
	_, _, _ = http.NewClient("https://httpbin.org/get").Body(nil).Head(nil).Post()
}
