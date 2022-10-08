package http

import "testing"

func Test_client_Post(t *testing.T) {
	_, _ = NewClient("https://www.fsgit.cc").Body(nil).Head(nil).Post()
}
