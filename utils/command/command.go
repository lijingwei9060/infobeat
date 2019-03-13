package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var (
	NoCommand = fmt.Errorf("no command error")
)

// TimeOutCommand 可以执行command，定时超时
type TimeOutCommand struct {
	config Config         // config 配置信息
	wg     sync.WaitGroup // wg 控制command等待
	stop   chan struct{}  // stop 接收到信息就停掉所有执行的command
	pool   chan int       // pool 执行command的控制池
}

// New 根据参数创建一个TimeOutCommand执行期
func New(config Config) (*TimeOutCommand, error) {
	command := &TimeOutCommand{
		config: config,
		stop:   make(chan struct{}),
		pool:   make(chan int, config.Backoff),
	}
	return command, nil
}

// RunCommand 根据并发情况和超时时间执行
// 有错误并不表示有问题，比如ping就需要设置超时时间才退出
func (tc *TimeOutCommand) RunCommand(command ...string) ([]byte, error) {
	if len(command) < 1 {
		return nil, NoCommand
	}
	tc.pool <- 0 //确认pool还有空间处理command，没有则阻塞
	defer func() {
		<-tc.pool //任务执行完从poll里面取掉1个任务，准备开始执行
	}()

	ctx, cancel := context.WithTimeout(context.Background(), tc.config.TimeOut)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	var cmd *exec.Cmd
	if len(command) > 1 {
		cmd = exec.CommandContext(ctx, command[0], command[1:]...)
	} else {
		cmd = exec.CommandContext(ctx, command[0])
	}

	cmd.Env = append(os.Environ(), "lang=en_US.utf-8")
	// This time we can simply use CombinedOutput() to get the result.
	out, err := cmd.Output()

	// Check the context error to see if the timeout was executed.
	// The error returned by cmd.CombinedOutput() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		return out, context.DeadlineExceeded
	}

	// If there's no context error, we know the command completed (or errored).
	return out, err
}
