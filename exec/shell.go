package exec

import (
	"bufio"
	"context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/utils/str"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

// RunShellCommand 执行shell命令
// command：要执行的命令
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：exit code
func RunShellCommand(command string, environment map[string]string, workingDirectory string, outputCmd bool) (int, []string) {
	receiveOutput := make(chan string, 1000)
	result := runCmdContext(context.Background(), "bash", []string{"-c", command}, receiveOutput, environment, workingDirectory, outputCmd)
	return result, collections.NewListFromChan(receiveOutput).ToArray()
}

// RunShell 执行shell命令
// command：要执行的命令
// receiveOutput：输出流
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：exit code
func RunShell(command string, receiveOutput chan string, environment map[string]string, workingDirectory string, outputCmd bool) int {
	return runCmdContext(context.Background(), "bash", []string{"-c", command}, receiveOutput, environment, workingDirectory, outputCmd)
}

// RunShellContext 执行shell命令
// command：要执行的命令
// receiveOutput：输出流
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：exit code
func RunShellContext(ctx context.Context, command string, receiveOutput chan string, environment map[string]string, workingDirectory string, outputCmd bool) int {
	return runCmdContext(ctx, "bash", []string{"-c", command}, receiveOutput, environment, workingDirectory, outputCmd)
}

// RunShellContext 执行shell命令
// command：要执行的命令
// receiveOutput：输出流
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：exit code
func runCmdContext(ctx context.Context, command string, args []string, receiveOutput chan string, environment map[string]string, workingDirectory string, outputCmd bool) int {
	if outputCmd {
		receiveOutput <- command + " " + strings.Join(args, " ")
	}
	cmd := exec.CommandContext(ctx, command, args...)
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
			if err.Error() != "EOF" && err.Error() != "read |0: file already closed" {
				receiveOutput <- err.Error()
			}
			break
		}
		receiveOutput <- str.CutRight(line, "\n")
	}
}
