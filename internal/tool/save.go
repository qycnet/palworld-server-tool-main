package tool

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/qycnet/palworld-server-tool-main/internal/auth"
	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
	"github.com/qycnet/palworld-server-tool-main/internal/source"
	"github.com/qycnet/palworld-server-tool-main/internal/system"
)

type Sturcture struct {
	Players []database.Player `json:"players"`
	Guilds  []database.Guild  `json:"guilds"`
}

func getSavCli() (string, error) {
	// 从配置文件中获取savCli的路径
	savCliPath := viper.GetString("save.decode_path")
	if savCliPath == "" || savCliPath == "/path/to/your/sav_cli" {
		// 如果路径为空或为默认路径，则获取执行目录
		ed, err := system.GetExecDir()
		if err != nil {
			// 如果获取执行目录失败，记录错误并返回
			logger.Errorf("error getting exec directory: %s", err)
			return "", err
		}
		// 将执行目录与"sav_cli"拼接，得到savCli的实际路径
		savCliPath = filepath.Join(ed, "sav_cli")
		if runtime.GOOS == "windows" {
			// 如果是Windows系统，则在路径后添加".exe"后缀
			savCliPath += ".exe"
		}
	}
	// 检查savCli路径是否存在
	if _, err := os.Stat(savCliPath); err != nil {
		// 如果路径不存在，返回错误
		return "", err
	}
	// 返回savCli路径和nil表示成功
	return savCliPath, nil
}

func Decode(file string) error {
	// 获取可执行文件的路径
	savCli, err := getSavCli()
	if err != nil {
		return errors.New("error getting executable path: " + err.Error())
	}

	// 从源文件中获取解码文件路径
	levelFilePath, err := getFromSource(file, "decode")
	if err != nil {
		return err
	}
	// 删除临时目录
	defer os.RemoveAll(filepath.Dir(levelFilePath))

	// 构建基础URL
	baseUrl := fmt.Sprintf("http://127.0.0.1:%d", viper.GetInt("web.port"))
	// 判断是否使用TLS且URL没有以"/"结尾
	if viper.GetBool("web.tls") && !strings.HasSuffix(baseUrl, "/") {
		baseUrl = viper.GetString("web.public_url")
	}

	// 构建请求URL
	requestUrl := fmt.Sprintf("%s/api/", baseUrl)
	// 生成令牌
	tokenString, err := auth.GenerateToken()
	if err != nil {
		return errors.New("error generating token: " + err.Error())
	}
	// 构建执行命令的参数
	execArgs := []string{"-f", levelFilePath, "--request", requestUrl, "--token", tokenString}
	// 创建执行命令
	cmd := exec.Command(savCli, execArgs...)
	// 设置命令的标准输出和标准错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 启动命令
	err = cmd.Start()
	if err != nil {
		return errors.New("error starting command: " + err.Error())
	}
	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		return errors.New("error waiting for command: " + err.Error())
	}

	return nil
}

func Backup() (string, error) {
	// 从配置文件中获取保存路径
	sourcePath := viper.GetString("save.path")

	// 从源路径中获取文件路径
	levelFilePath, err := getFromSource(sourcePath, "backup")
	if err != nil {
		return "", err
	}
	// 在函数结束时删除 levelFilePath 所在目录
	defer os.RemoveAll(filepath.Dir(levelFilePath))

	// 获取备份目录
	backupDir, err := GetBackupDir()
	if err != nil {
		return "", fmt.Errorf("failed to get backup directory: %s", err)
	}

	// 获取当前时间，并格式化为字符串
	currentTime := time.Now().Format("2006-01-02-15-04-05")
	// 拼接备份文件的路径
	backupZipFile := filepath.Join(backupDir, fmt.Sprintf("%s.zip", currentTime))
	// 将 levelFilePath 所在目录压缩到备份文件
	err = system.ZipDir(filepath.Dir(levelFilePath), backupZipFile)
	if err != nil {
		return "", fmt.Errorf("failed to create backup zip: %s", err)
	}
	// 返回备份文件的文件名
	return filepath.Base(backupZipFile), nil
}

func GetBackupDir() (string, error) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		// 如果获取工作目录出错，则返回错误
		return "", err
	}

	// 拼接当前工作目录和备份目录名，得到备份目录的完整路径
	backDir := filepath.Join(wd, "backups")
	// 检查并创建备份目录
	if err = system.CheckAndCreateDir(backDir); err != nil {
		// 如果创建备份目录出错，则返回错误
		return "", err
	}

	// 返回备份目录的路径和nil错误
	return backDir, nil
}

func getFromSource(file, way string) (string, error) {
	var levelFilePath string
	var err error

	// 处理以http(s)://开头的URL
	if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
		// http(s)://url
		levelFilePath, err = source.DownloadFromHttp(file, way)
		if err != nil {
			return "", errors.New("error downloading file: " + err.Error())
		}
	}
	// 处理以k8s://开头的Kubernetes地址
	else if strings.HasPrefix(file, "k8s://") {
		// k8s://namespace/pod/container:remotePath
		namespace, podName, container, remotePath, err := source.ParseK8sAddress(file)
		if err != nil {
			return "", errors.New("error parsing k8s address: " + err.Error())
		}
		levelFilePath, err = source.CopyFromPod(namespace, podName, container, remotePath, way)
		if err != nil {
			return "", errors.New("error copying file from pod: " + err.Error())
		}
	}
	// 处理以docker://开头的Docker地址
	else if strings.HasPrefix(file, "docker://") {
		// docker://containerID(Name):remotePath
		containerId, remotePath, err := source.ParseDockerAddress(file)
		if err != nil {
			return "", errors.New("error parsing docker address: " + err.Error())
		}
		levelFilePath, err = source.CopyFromContainer(containerId, remotePath, way)
		if err != nil {
			return "", errors.New("error copying file from container: " + err.Error())
		}
	}
	// 处理本地文件
	else {
		// local file
		levelFilePath, err = source.CopyFromLocal(file, way)
		if err != nil {
			return "", errors.New("error copying file to temporary directory: " + err.Error())
		}
	}
	return levelFilePath, nil
}
