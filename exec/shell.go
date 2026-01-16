package exec

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/async"
	"github.com/farseer-go/utils/str"
)

// 当调用完RunShell后，希望将结果写入到progress中
func SaveToChan(progress chan string, receiveOutput chan string, wait func() int) int {
	// 异步转存消息
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		for output := range receiveOutput {
			progress <- output
		}
	}()

	// 等待命令执行完成
	exitCode := wait()
	// 等待消息转存完成
	waitGroup.Wait()

	return exitCode
}

// RunShellCommand 执行shell命令（同步版本，等待命令执行完成）
// command：要执行的命令
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：输出行数组 和 exit code
func RunShellCommand(command string, environment map[string]string, workingDirectory string, outputCmd bool) (collections.List[string], int) {
	receiveOutput, wait := runCmdContext(context.Background(), "bash", []string{"-c", command}, environment, workingDirectory, outputCmd)

	lstResult := collections.NewList[string]()
	// 异步接收消息
	worker := async.New()
	worker.Add(func() {
		for output := range receiveOutput {
			lstResult.Add(output)
		}
	})
	defer worker.Wait()

	return lstResult, wait()
}

// RunShell 执行shell命令（异步版本，立即返回）
// command：要执行的命令
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：输出流 channel（可实时接收）和 wait 函数（调用后阻塞等待命令完成并返回 exit code）
func RunShell(command string, environment map[string]string, workingDirectory string, outputCmd bool) (chan string, func() int) {
	return runCmdContext(context.Background(), "bash", []string{"-c", command}, environment, workingDirectory, outputCmd)
}

// RunShellContext 执行shell命令（支持上下文控制）
// ctx：上下文
// command：要执行的命令
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：输出流 channel（可实时接收）和 wait 函数（调用后阻塞等待命令完成并返回 exit code）
func RunShellContext(ctx context.Context, command string, environment map[string]string, workingDirectory string, outputCmd bool) (chan string, func() int) {
	return runCmdContext(ctx, "bash", []string{"-c", command}, environment, workingDirectory, outputCmd)
}

// runCmdContext 执行命令（内部方法）
// ctx：上下文
// command：要执行的命令
// args：命令参数
// environment：环境变量
// workingDirectory：当前工作目录位置
// outputCmd：是否输出命令本身
// return：输出流 channel（实时输出）和 wait 函数（由外部调用以阻塞等待命令完成）
// 设计说明：
//   - 函数立即返回，不阻塞
//   - 命令已启动，输出会实时发送到 channel
//   - 调用 wait() 函数会阻塞直到命令执行完成，并关闭 channel
//   - 外部可以选择何时调用 wait()，从而控制同步/异步行为
func runCmdContext(ctx context.Context, command string, args []string, environment map[string]string, workingDirectory string, outputCmd bool) (chan string, func() int) {
	receiveOutput := make(chan string, 100)

	if outputCmd {
		receiveOutput <- command + " " + strings.Join(args, " ")
	}

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = workingDirectory
	// 如果设置了环境变量，则追加进来
	if environment != nil {
		cmd.Env = append(os.Environ(), str.MapToStringList(environment)...)
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		receiveOutput <- "执行失败：" + err.Error()
		return receiveOutput, func() int {
			close(receiveOutput)
			return -1
		}
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go readInputStream(stdout, receiveOutput, &waitGroup)
	go readInputStream(stderr, receiveOutput, &waitGroup)

	// 返回 wait 函数，由外部决定何时调用（阻塞等待命令完成）
	wait := func() int {
		var res int
		err := cmd.Wait()
		waitGroup.Wait()
		close(receiveOutput)

		if err != nil {
			if ex, ok := err.(*exec.ExitError); ok {
				res = ex.Sys().(syscall.WaitStatus).ExitStatus() // 获取命令执行返回状态
			}
		}
		return res
	}

	return receiveOutput, wait
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
