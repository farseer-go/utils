package test

import (
	"context"
	"github.com/farseer-go/utils/exec"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestRunShellContext(t *testing.T) {
	receiveOutput := make(chan string, 100)
	ctx, cancel := context.WithCancel(context.Background())
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		for output := range receiveOutput {
			assert.True(t, output == "执行失败：context canceled" || output == "Sleep 1")
		}
	}()

	go func() {
		defer waitGroup.Done()
		exitCode := exec.RunShellContext(ctx, "Sleep 1", receiveOutput, nil, "", false)
		close(receiveOutput)
		assert.Equal(t, -1, exitCode)
	}()
	cancel()

	waitGroup.Wait()
}

func TestRunShell(t *testing.T) {
	t.Run("env test", func(t *testing.T) {
		receiveOutput := make(chan string, 100)
		env := map[string]string{
			"a": "b",
		}
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			exist := false
			for output := range receiveOutput {
				if output == "a=b" {
					exist = true
				}
			}
			assert.True(t, exist)
		}()
		exitCode := exec.RunShell("env", receiveOutput, env, "", false)
		close(receiveOutput)
		assert.Equal(t, 0, exitCode)
		waitGroup.Wait()
	})
	t.Run("error test", func(t *testing.T) {
		receiveOutput := make(chan string, 100)
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)
		command := "commandError"
		go func() {
			defer waitGroup.Done()
			var res string
			for output := range receiveOutput {
				res = output
			}
			assert.Contains(t, res, "commandError: command not found")
		}()

		_ = exec.RunShell(command, receiveOutput, nil, "", false)
		close(receiveOutput)
		waitGroup.Wait()
		// assert.Equal(t, 0, exitCode)

	})
}
