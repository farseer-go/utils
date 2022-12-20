package test

import (
	"github.com/farseer-go/utils/http"
	"testing"
)

func Test_client_Post(t *testing.T) {
	_, _ = http.NewClient("https://www.fsgit.cc").Body(nil).Head(nil).Post()
}
