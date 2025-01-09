package executor

import (
	"errors"
	"strings"
	"time"

	"github.com/gorcon/rcon"
)

var (
	ErrPasswordEmpty = errors.New("password is empty")
)

type ExecuteCloser interface {
	Execute(command string) (string, error)
	Close() error
}

type Executor struct {
	skipErrors bool
	client     ExecuteCloser
}

func NewExecutor(address, password string, timeout int, skipErrors bool) (*Executor, error) {
	// 检查密码是否为空
	if password == "" {
		return nil, ErrPasswordEmpty
	}

	// 将超时时间转换为时间.Duration类型
	timeoutDuration := time.Duration(timeout) * time.Second

	// 连接到RCON服务器
	client, err := rcon.Dial(address, password, rcon.SetDialTimeout(timeoutDuration), rcon.SetDeadline(timeoutDuration))
	if err != nil {
		// 如果连接失败，则返回错误
		return nil, err
	}

	// 返回Executor结构体指针和nil错误
	return &Executor{client: client, skipErrors: skipErrors}, nil
}

func (e *Executor) Execute(command string) (string, error) {
	// 执行命令并获取响应和错误
	response, err := e.client.Execute(command)
	// 去除响应字符串前后的空白字符
	response = strings.TrimSpace(response)
	// 如果存在错误，并且设置了忽略错误，且响应不为空
	if err != nil && e.skipErrors && response != "" {
		// 返回响应和空错误
		return response, nil
	}
	// 返回响应和错误
	return response, err
}

func (e *Executor) Close() error {
	// 如果client不为空
	if e.client != nil {
		// 调用client的Close方法
		return e.client.Close()
	}
	// 如果client为空，则返回nil
	return nil
}
