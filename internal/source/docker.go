package source

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/qycnet/palworld-server-tool-main/internal/system"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
)

func getDockerClient() (*client.Client, error) {
	// 获取环境变量 DOCKER_API_VERSION 的值
	dockerAPIVersion := os.Getenv("DOCKER_API_VERSION")
	if dockerAPIVersion == "" {
		// 如果 DOCKER_API_VERSION 为空，则使用默认配置创建 Docker 客户端
		return client.NewClientWithOpts(client.FromEnv)
	} else {
		// 如果 DOCKER_API_VERSION 不为空，则使用指定的 API 版本创建 Docker 客户端
		return client.NewClientWithOpts(client.FromEnv, client.WithVersion(dockerAPIVersion))
	}
}

func CopyFromContainer(containerID, remotePath, way string) (string, error) {
	// 记录日志，记录从远程路径复制的操作
	logger.Infof("copying savDir from %s\n", remotePath)

	// 获取Docker客户端
	cli, err := getDockerClient()
	if err != nil {
		return "", err
	}
	defer cli.Close()

	// 取得Level.sav所在目录
	// 执行find命令查找Level.sav所在的目录
	findCmd := []string{"sh", "-c", fmt.Sprintf("find %s -maxdepth 4 -path '*/backup/*' -prune -o -name 'Level.sav' -print | xargs dirname", remotePath)}
	savDir, err := execCommand(containerID, findCmd, cli)
	if err != nil {
		return "", err
	}
	// 去除字符串两端的空白字符
	savDir = strings.TrimSpace(savDir)
	if savDir == "" {
		return "", errors.New("directory containing Level.sav not found in container")
	}

	// 压缩
	// 创建tar命令，压缩savDir目录下的.sav和Players目录下的.sav文件
	tarCmd := []string{"sh", "-c", fmt.Sprintf("cd \"%s\" && tar czf - ./*.sav ./Players/*.sav", savDir)}
	tarReader, err := execCommandStream(containerID, tarCmd, cli)
	if err != nil {
		return "", err
	}

	// 创建临时目录
	// 生成一个UUID作为临时目录的一部分
	id := uuid.New().String()
	// 拼接出临时目录的路径
	tempDir := filepath.Join(os.TempDir(), "palworldsav-docker-"+way+"-"+id)
	// 创建临时目录
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// 解压文件
	// 使用system.UnTarGzDir函数解压tarReader中的数据到临时目录中
	err = system.UnTarGzDir(tarReader, tempDir)
	if err != nil {
		return "", err
	}

	// 返回Level.sav文件的路径
	levelFilePath := filepath.Join(tempDir, "Level.sav")
	return levelFilePath, nil
}

func execCommandStream(containerID string, command []string, cli *client.Client) (io.Reader, error) {
	// 创建一个背景上下文
	ctx := context.Background()

	// 创建执行配置
	execConfig := types.ExecConfig{
		Cmd:          command,        // 执行的命令
		AttachStdout: true,           // 附加标准输出
		AttachStderr: true,           // 附加标准错误输出
	}

	// 在容器中执行命令并创建执行实例
	ir, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return nil, err
	}

	// 附加到执行实例
	hr, err := cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}

	// 创建一个管道
	reader, writer := io.Pipe()

	// 启动一个goroutine来处理管道输出
	go func() {
		defer writer.Close()     // 延迟关闭写入器
		defer hr.Close()         // 延迟关闭附加实例
		_, err = stdcopy.StdCopy(writer, os.Stderr, hr.Reader) // 复制标准输出和标准错误输出到管道
		if err != nil {
			logger.Errorf("Stream to docker failed: %v", err) // 记录错误日志
		}
	}()

	return reader, nil
}

func execCommand(containerID string, command []string, cli *client.Client) (string, error) {
	ctx := context.Background()

	// 创建Exec配置
	execConfig := types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
	}

	// 创建Exec实例
	ir, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", err
	}

	// 附加到Exec实例
	hr, err := cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}

	// 关闭连接
	defer hr.Close()

	// 创建一个缓冲区来存储输出
	var outBuf bytes.Buffer

	// 将Exec实例的输出复制到缓冲区和标准错误输出
	_, err = stdcopy.StdCopy(&outBuf, os.Stderr, hr.Reader)
	if err != nil {
		return "", err
	}

	// 返回输出字符串和nil错误
	return outBuf.String(), nil
}

func ParseDockerAddress(address string) (containerID, filePath string, err error) {
	// 去除地址前缀 "docker://"
	address = strings.TrimPrefix(address, "docker://")

	// 使用冒号将地址分割为两部分
	parts := strings.SplitN(address, ":", 2)
	// 如果分割后的部分数量不等于2，则返回错误
	if len(parts) != 2 {
		// 返回错误信息 "invalid Docker address format"
		return "", "", errors.New("invalid Docker address format")
	}

	// 将分割后的两部分分别赋值给 containerID 和 filePath
	containerID, filePath = parts[0], parts[1]
	// 返回 containerID, filePath 和 nil
	return containerID, filePath, nil
}
