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
	"github.com/farseer-go/utils/str"
)

type ShellWait struct {
	receiveOutput chan string
	wait          func() int
	waitGroup     *sync.WaitGroup
	once          *sync.Once
	exitCode      int
}

func NewExitShellWait(exitCode int, message string) ShellWait {
	receiveOutput := make(chan string, 1)
	receiveOutput <- message
	close(receiveOutput) // 关闭 channel，让消费者能正常退出
	return ShellWait{
		exitCode:      exitCode,
		receiveOutput: receiveOutput,
		waitGroup:     new(sync.WaitGroup),
		once:          new(sync.Once),
		wait: func() int {
			return exitCode
		},
	}
}

// 等待命令执行完成,不需要返回结果
func (receiver *ShellWait) Wait() int {
	receiver.once.Do(func() {
		// 异步转存消息
		receiver.waitGroup.Add(1)
		go func() {
			defer receiver.waitGroup.Done()
			for range receiver.receiveOutput {
			}
		}()

		// 等待命令执行完成
		receiver.exitCode = receiver.wait()
		// 等待消息转存完成
		receiver.waitGroup.Wait()
	})
	return receiver.exitCode
}

// 将结果写入到progress中
func (receiver *ShellWait) WaitToChan(progress chan string) int {
	receiver.once.Do(func() {
		// 异步转存消息
		receiver.waitGroup.Add(1)
		go func() {
			defer receiver.waitGroup.Done()
			for output := range receiver.receiveOutput {
				progress <- output
			}
		}()

		// 等待命令执行完成
		receiver.exitCode = receiver.wait()
		// 等待消息转存完成
		receiver.waitGroup.Wait()
	})
	return receiver.exitCode
}

// 将结果写入到progress中
func (receiver *ShellWait) WaitToFunc(f func(progress string)) int {
	receiver.once.Do(func() {
		// 异步转存消息
		receiver.waitGroup.Add(1)
		go func() {
			defer receiver.waitGroup.Done()
			for output := range receiver.receiveOutput {
				f(output)
			}
		}()

		// 等待命令执行完成
		receiver.exitCode = receiver.wait()
		// 等待消息转存完成
		receiver.waitGroup.Wait()
	})
	return receiver.exitCode
}

// 将结果写入到collections.List[string]中
func (receiver *ShellWait) WaitToList() (collections.List[string], int) {
	lstResult := collections.NewList[string]()
	receiver.once.Do(func() {
		// 异步转存消息
		receiver.waitGroup.Add(1)
		go func() {
			defer receiver.waitGroup.Done()
			for output := range receiver.receiveOutput {
				lstResult.Add(output)
			}
		}()

		// 等待命令执行完成
		receiver.exitCode = receiver.wait()
		// 等待消息转存完成
		receiver.waitGroup.Wait()
	})
	return lstResult, receiver.exitCode
}

// 将结果写入到string中
func (receiver *ShellWait) WaitToFirstResult() (string, int) {
	var firstResult string
	receiver.once.Do(func() {
		// 异步转存消息
		receiver.waitGroup.Add(1)
		go func() {
			defer receiver.waitGroup.Done()
			lstResult := collections.NewList[string]()
			for output := range receiver.receiveOutput {
				lstResult.Add(output)
			}
			firstResult = lstResult.First()
		}()

		// 等待命令执行完成
		receiver.exitCode = receiver.wait()
		// 等待消息转存完成
		receiver.waitGroup.Wait()
	})
	return firstResult, receiver.exitCode
}

// RunShell 执行shell命令（异步版本，立即返回）
// command：要执行的命令
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：输出流 channel（可实时接收）和 wait 函数（调用后阻塞等待命令完成并返回 exit code）
func RunShell(command string, args []string, environment map[string]string, workingDirectory string, outputCmd bool) ShellWait {
	return runCmdContext(context.Background(), command, args, environment, workingDirectory, outputCmd)
}

// RunShellContext 执行shell命令（支持上下文控制）
// ctx：上下文
// command：要执行的命令
// environment：环境变量
// workingDirectory：当前工作目录位置
// return：输出流 channel（可实时接收）和 wait 函数（调用后阻塞等待命令完成并返回 exit code）
func RunShellContext(ctx context.Context, command string, args []string, environment map[string]string, workingDirectory string, outputCmd bool) ShellWait {
	//return runCmdContext(ctx, "bash", []string{"-c", command}, environment, workingDirectory, outputCmd)
	return runCmdContext(ctx, command, args, environment, workingDirectory, outputCmd)
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
func runCmdContext(ctx context.Context, command string, args []string, environment map[string]string, workingDirectory string, outputCmd bool) ShellWait {
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
		return NewExitShellWait(-1, "执行失败："+err.Error())
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go readInputStream(stdout, receiveOutput, &waitGroup)
	go readInputStream(stderr, receiveOutput, &waitGroup)

	// 返回 wait 函数，由外部决定何时调用（阻塞等待命令完成）
	return ShellWait{
		waitGroup:     new(sync.WaitGroup),
		once:          new(sync.Once),
		receiveOutput: receiveOutput,
		wait: func() int {
			var res int
			err := cmd.Wait()
			waitGroup.Wait()
			close(receiveOutput)

			if err != nil {
				// 优先检查 context 超时
				if ctx.Err() == context.DeadlineExceeded {
					res = 124 // timeout 退出码
				} else if ctx.Err() == context.Canceled {
					res = 130 // SIGINT 退出码
				} else if ex, ok := err.(*exec.ExitError); ok {
					if status, ok := ex.Sys().(syscall.WaitStatus); ok {
						res = status.ExitStatus()
					} else {
						res = -1
					}
				} else {
					res = -1
				}
			}
			return res
		},
	}
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
