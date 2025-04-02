package tool

import (
	"encoding/base64"

	"github.com/spf13/viper"
	"github.com/qycnet/palworld-server-tool-main/internal/executor"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
)

func executeCommand(command string) (*executor.Executor, string, error) {
	// 判断是否使用Base64编码
	useBase64 := viper.GetBool("rcon.use_base64")

	// 创建Executor实例
	exec, err := executor.NewExecutor(
		viper.GetString("rcon.address"),  // 获取rcon地址
		viper.GetString("rcon.password"), // 获取rcon密码
		viper.GetInt("rcon.timeout"),      // 获取rcon超时时间
		true)                              // 启用调试模式
	if err != nil {
		return nil, "", err
	}

	// 如果使用Base64编码，则对命令进行编码
	if useBase64 {
		command = base64.StdEncoding.EncodeToString([]byte(command))
	}

	// 执行命令
	response, err := exec.Execute(command)
	if err != nil {
		return nil, "", err
	}

	// 如果使用Base64编码，则对响应进行解码
	if useBase64 {
		decoded, err := base64.StdEncoding.DecodeString(response)
		if err != nil {
			logger.Warnf("decode base64 '%s' error: %v\n", response, err)
			return exec, response, nil
		}
		// 将解码后的字节数组转换为字符串
		response = string(decoded)
	}

	// 返回执行器实例、响应和错误
	return exec, response, nil
}

func CustomCommand(command string) (string, error) {
	// 执行命令
	exec, response, err := executeCommand(command)
	if err != nil {
		// 如果执行命令出错，则返回错误
		return "", err
	}
	// 关闭执行命令的资源
	defer exec.Close()

	// 返回命令执行结果
	return response, nil
}
