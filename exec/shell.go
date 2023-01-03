package exec

import (
	"bufio"
	"context"
	"github.com/farseer-go/utils/str"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

// RunShell 执行shell命令
// command：要执行的命令
// receiveOutput：输出流
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：exit code
func RunShell(command string, receiveOutput chan string, environment map[string]string, workingDirectory string) int {
	return RunShellContext(context.Background(), command, receiveOutput, environment, workingDirectory)
}

// RunShellContext 执行shell命令
// command：要执行的命令
// receiveOutput：输出流
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：exit code
func RunShellContext(ctx context.Context, command string, receiveOutput chan string, environment map[string]string, workingDirectory string) int {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Dir = workingDirectory
	// 如果设置了环境变量，则追回进来
	if environment != nil {
		cmd.Env = append(os.Environ(), str.MapToStringList(environment)...)
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		receiveOutput <- "执行失败：" + err.Error()
		return -1
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go readInputStream(stdout, receiveOutput, &waitGroup)
	go readInputStream(stderr, receiveOutput, &waitGroup)

	var res int
	err := cmd.Wait()
	waitGroup.Wait()

	if err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			res = ex.Sys().(syscall.WaitStatus).ExitStatus() // 获取命令执行返回状态
		}
		if !strings.Contains(err.Error(), "exit status") {
			receiveOutput <- "wait:" + err.Error()
		}
	}
	return res
}

func readInputStream(out io.ReadCloser, receiveOutput chan string, waitGroup *sync.WaitGroup) {
	defer func() {
		waitGroup.Done()
		_ = out.Close()
	}()
	reader := bufio.NewReader(out)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if io.EOF != err {
				receiveOutput <- err.Error()
			}
			break
		}
		receiveOutput <- str.CutRight(line, "\n")
	}
}
