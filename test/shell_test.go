package test

import (
	"sync"
	"testing"

	"github.com/farseer-go/utils/exec"
	"github.com/stretchr/testify/assert"
)

func TestRunShellContext(t *testing.T) {
	receiveOutput, exitCode := exec.RunShellCommand("Sleep 1", nil, "", false)
	receiveOutput.Foreach(func(output *string) {
		assert.True(t, *output == "执行失败：context canceled" || *output == "Sleep 1")
	})

	assert.Equal(t, 0, exitCode)
}

func TestRunShell(t *testing.T) {
	t.Run("env test", func(t *testing.T) {
		env := map[string]string{
			"a": "b",
		}
		receiveOutput, wait := exec.RunShell("env", env, "", false)

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
		assert.Equal(t, 0, wait())
		waitGroup.Wait()
	})

	t.Run("error test", func(t *testing.T) {
		receiveOutput, wait := exec.RunShell("commandError", nil, "", false)

		var waitGroup sync.WaitGroup
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			var res string
			for output := range receiveOutput {
				res = output
			}
			assert.Contains(t, res, "commandError: command not found")
		}()

		assert.Equal(t, 127, wait())
		waitGroup.Wait()
	})
}
