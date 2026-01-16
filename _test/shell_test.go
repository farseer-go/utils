package test

import (
	"testing"

	"github.com/farseer-go/fs/async"
	"github.com/farseer-go/utils/exec"
	"github.com/stretchr/testify/assert"
)

func TestRunShellContext(t *testing.T) {
	receiveOutput, exitCode := exec.RunShellCommand("Sleep 1", nil, "", true)
	assert.Equal(t, "bash -c Sleep 1", receiveOutput.First())
	assert.Equal(t, 0, exitCode)
}

func TestRunShell(t *testing.T) {
	t.Run("env test", func(t *testing.T) {
		env := map[string]string{
			"a": "b",
		}
		receiveOutput, wait := exec.RunShell("env", env, "", false)

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
		assert.Equal(t, 0, wait())
		worker.Wait()
	})

	t.Run("error test", func(t *testing.T) {
		receiveOutput, wait := exec.RunShell("commandError", nil, "", false)

		worker := async.New()
		worker.Add(func() {
			var res string
			for output := range receiveOutput {
				res = output
			}
			assert.Contains(t, res, "commandError: command not found")
		})

		assert.Equal(t, 127, wait())
		worker.Wait()
	})
}
