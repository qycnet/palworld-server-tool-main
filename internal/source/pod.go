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
	"time"

	"github.com/google/uuid"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
	"github.com/qycnet/palworld-server-tool-main/internal/system"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	ErrPodNotFound    = errors.New("未找到 Pod")
	ErrContainerEmpty = errors.New("容器为空")
	ErrAddressInvalid = errors.New("无效的 save.path, eg: k8s://namespace/podName:filePath")
)

func CopyFromPod(namespace, podName, container, remotePath, way string) (string, error) {
	// 记录日志，显示从哪个容器和远程路径复制
	logger.Infof("复制 savDir 自 %s:%s\n", container, remotePath)

	// 获取集群内的配置
	config, err := rest.InClusterConfig()
	if err != nil {
		// 如果获取配置失败，返回错误信息
		return "", errors.New("获取集群内配置时出错: " + err.Error())
	}

	// 创建clientset对象
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// 如果创建clientset失败，返回错误信息
		return "", errors.New("获取客户端集时出错: " + err.Error())
	}

	// 如果命名空间为空，则获取当前命名空间
	if namespace == "" {
		var err error
		namespace, err = getCurrentNamespace()
		if err != nil {
			// 如果获取当前命名空间失败，返回错误信息
			return "", errors.New("获取当前命名空间时出错: " + err.Error())
		}
	}

	// 如果容器名为空，返回错误
	if container == "" {
		return "", ErrContainerEmpty
	}

	// 构建查找命令
	findCmd := []string{"sh", "-c", fmt.Sprintf("find %s -maxdepth 4 -path '*/backup/*' -prune -o -name 'Level.sav' -print | xargs dirname", remotePath)}
	// 执行查找命令，获取保存目录
	savDir, err := execPodCommand(clientset, config, namespace, podName, container, findCmd)
	if err != nil {
		// 如果执行查找命令失败，返回错误信息
		return "", errors.New("执行查找命令时出错: " + err.Error())
	}
	// 去除保存目录前后的空白字符
	savDir = strings.TrimSpace(savDir)
	// 如果保存目录为空，返回错误信息
	if savDir == "" {
		return "", errors.New("包含 Level.sav 的目录在 Pod 中找不到")
	}
	// 记录日志，显示保存目录路径
	logger.Debugf("目录路径: %s\n", savDir)

	// 构建打包命令
	tarCmd := []string{"sh", "-c", fmt.Sprintf("cd \"%s\" && tar czf - ./*.sav ./Players/*.sav", savDir)}
	// 执行打包命令，获取打包流
	tarStream, err := execPodCommandStream(clientset, config, namespace, podName, container, tarCmd)
	if err != nil {
		// 如果执行打包命令失败，返回错误信息
		return "", errors.New("执行打包命令时出错: " + err.Error())
	}

	// 生成临时目录
	id := uuid.New().String()
	tempDir := filepath.Join(os.TempDir(), "palworldsav-pod-"+way+"-"+id)
	// 创建临时目录
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		// 如果创建临时目录失败，返回错误信息
		return "", err
	}

	// 解压打包流到临时目录
	err = system.UnTarGzDir(tarStream, tempDir)
	if err != nil {
		// 如果解压失败，返回错误信息
		return "", err
	}

	// 记录日志，显示从Pod复制的目录
	logger.Debugf("从 Pod 复制的目录: %s\n", tempDir)

	// 构建Level.sav文件路径
	levelFilePath := filepath.Join(tempDir, "Level.sav")
	return levelFilePath, nil
}

func getCurrentNamespace() (string, error) {
	// 读取文件 /var/run/secrets/kubernetes.io/serviceaccount/namespace
	ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		// 如果读取文件出错，则返回错误
		return "", err
	}
	// 去除读取内容的首尾空白字符，并转换为字符串类型
	return strings.TrimSpace(string(ns)), nil
}

func execPodCommand(clientset *kubernetes.Clientset, config *rest.Config, namespace, podName, container string, cmd []string) (string, error) {
	// 构建执行Pod命令的请求
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   cmd,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
			Container: container,
		}, scheme.ParameterCodec)

	// 设置请求的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// 创建一个SPDY执行器
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", err
	}

	// 初始化标准输出和标准错误缓冲区
	var stdout, stderr bytes.Buffer
	// 执行命令并处理输入输出
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return "", err
	}

	// 如果有错误输出，则返回错误
	if stderr.Len() > 0 {
		return "", errors.New(stderr.String())
	}

	// 返回标准输出内容
	return stdout.String(), nil
}

func execPodCommandStream(clientset *kubernetes.Clientset, config *rest.Config, namespace, podName, container string, cmd []string) (io.Reader, error) {
	// 构建执行命令的请求
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   cmd,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
			Container: container,
		}, scheme.ParameterCodec)

	// 创建SPDY执行器
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, err
	}

	// 创建管道用于读写数据
	reader, writer := io.Pipe()

	// 创建上下文并设置取消函数
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动协程处理流数据
	go func() {
		defer writer.Close()
		// 执行流操作
		err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdout: writer,
			Stderr: os.Stderr,
		})
		if err != nil {
			// 记录错误日志并取消上下文
			logger.Errorf("Stream to pod failed: %v", err)
			cancel()
		}
	}()

	// 返回读取器
	return reader, nil
}

func ParseK8sAddress(address string) (namespace, pod, container, filePath string, err error) {
	// 去除地址前缀 "k8s://"
	address = strings.TrimPrefix(address, "k8s://")

	// 按冒号分割地址
	parts := strings.SplitN(address, ":", 2)
	// 如果分割后的数组长度不等于2，则返回错误
	if len(parts) != 2 {
		return "", "", "", "", errors.New("地址格式无效")
	}

	// 按斜杠分割第一部分地址
	pathParts := strings.Split(parts[0], "/")
	// 根据路径部分的长度进行判断
	switch len(pathParts) {
	case 2: // podname  container
		// 提取 pod 和 container
		pod, container = pathParts[0], pathParts[1]
	case 3: // namespace  podname  container
		// 提取 namespace、pod 和 container
		namespace, pod, container = pathParts[0], pathParts[1], pathParts[2]
	default:
		// 如果路径格式不正确，则返回错误
		return "", "", "", "", errors.New("路径格式无效")
	}

	// 提取文件路径
	filePath = parts[1]
	return namespace, pod, container, filePath, nil
}
