package test

import (
	"testing"

	"github.com/farseer-go/fs/async"
	"github.com/farseer-go/utils/exec"
	"github.com/stretchr/testify/assert"
)

func TestRunShellContext(t *testing.T) {
	wait := exec.RunShell("Sleep", []string{"1"}, nil, "", true)
	receiveOutput, exitCode := wait.WaitToFirstResult()
	assert.Equal(t, "bash -c Sleep 1", receiveOutput)
	assert.Equal(t, 0, exitCode)
}

func TestRunShell(t *testing.T) {
	t.Run("env test", func(t *testing.T) {
		env := map[string]string{
			"a": "b",
		}
		wait := exec.RunShell("env", nil, env, "", false)
		receiveOutput := make(chan string, 10000)
		code := wait.WaitToChan(receiveOutput)

		worker := async.New()
		worker.Add(func() {
			exist := false
			for output := range receiveOutput {
				if output == "a=b" {
					exist = true
				}
			}
			assert.True(t, exist)
		})
		assert.Equal(t, 0, code)
		worker.Wait()
	})

	t.Run("error test", func(t *testing.T) {
		wait := exec.RunShell("commandError", nil, nil, "", false)
		receiveOutput := make(chan string, 10000)
		code := wait.WaitToChan(receiveOutput)

		worker := async.New()
		worker.Add(func() {
			var res string
			for output := range receiveOutput {
				res = output
			}
			assert.Contains(t, res, "commandError: command not found")
		})

		assert.Equal(t, 127, code)
		worker.Wait()
	})
}
