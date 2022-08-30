package exec

import (
	"context"
	"github.com/farseer-go/fs/flog"
	"testing"
)

func TestRunShell(t *testing.T) {
	receiveOutput := make(chan string, 100)
	ctx, cancel := context.WithCancel(context.Background())
	env := map[string]string{
		"a": "b",
	}

	go func() {
		for output := range receiveOutput {
			flog.Println(output)
		}
		//for {
		//	select {
		//	case output := <-receiveOutput:
		//		flog.Println(output)
		//	case <-ctx.Done():
		//		return
		//	}
		//}
	}()

	exitCode := RunShellContext("go env", receiveOutput, env, "", ctx)
	flog.Println(exitCode)
	cancel()
}
