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
	// 获取配置文件中的保存路径
	savCliPath := viper.GetString("save.decode_path")
	// 如果保存路径为空或者为默认值，则进行特殊处理
	if savCliPath == "" || savCliPath == "/path/to/your/sav_cli" {
		// 获取可执行文件的目录
		ed, err := system.GetExecDir()
		if err != nil {
			// 记录错误日志
			logger.Errorf("error getting exec directory: %s", err)
			return "", err
		}
		// 将可执行文件目录与"sav_cli"拼接，形成新的保存路径
		savCliPath = filepath.Join(ed, "sav_cli")
		// 如果操作系统为Windows，则在保存路径后添加".exe"
		if runtime.GOOS == "windows" {
			savCliPath += ".exe"
		}
	}
	// 检查保存路径是否存在
	if _, err := os.Stat(savCliPath); err != nil {
		return "", err
	}
	return savCliPath, nil
}

func Decode(file string) error {
	// 获取savCli客户端
	savCli, err := getSavCli()
	if err != nil {
		return errors.New("error getting executable path: " + err.Error())
	}

	// 从源文件中获取解码路径
	levelFilePath, err := getFromSource(file, "decode")
	if err != nil {
		return err
	}
	// 清理临时文件
	defer os.RemoveAll(filepath.Dir(levelFilePath))

	// 构造基础URL
	baseUrl := fmt.Sprintf("http://127.0.0.1:%d", viper.GetInt("web.port"))
	if viper.GetBool("web.tls") && !strings.HasSuffix(baseUrl, "/") {
		// 如果使用TLS且URL不以"/"结尾，则使用public_url
		baseUrl = viper.GetString("web.public_url")
	}

	// 构造请求URL
	requestUrl := fmt.Sprintf("%s/api/", baseUrl)
	tokenString, err := auth.GenerateToken()
	if err != nil {
		return errors.New("error generating token: " + err.Error())
	}
	// 构造执行参数
	execArgs := []string{"-f", levelFilePath, "--request", requestUrl, "--token", tokenString}
	// 创建执行命令
	cmd := exec.Command(savCli, execArgs...)
	// 设置命令输出
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

	// 从源路径获取备份文件路径
	levelFilePath, err := getFromSource(sourcePath, "backup")
	if err != nil {
		return "", err
	}
	// 在函数结束时删除备份文件所在的目录
	defer os.RemoveAll(filepath.Dir(levelFilePath))

	// 获取备份目录
	backupDir, err := GetBackupDir()
	if err != nil {
		return "", fmt.Errorf("failed to get backup directory: %s", err)
	}

	// 获取当前时间并格式化
	currentTime := time.Now().Format("2006-01-02-15-04-05")
	// 生成备份文件的路径
	backupZipFile := filepath.Join(backupDir, fmt.Sprintf("%s.zip", currentTime))
	// 将目录压缩为zip文件
	err = system.ZipDir(filepath.Dir(levelFilePath), backupZipFile)
	if err != nil {
		return "", fmt.Errorf("failed to create backup zip: %s", err)
	}
	// 返回备份文件的名称和nil错误
	return filepath.Base(backupZipFile), nil
}

func GetBackupDir() (string, error) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		// 如果获取工作目录失败，则返回错误
		return "", err
	}

	// 将工作目录和 "backups" 目录拼接，形成备份目录的路径
	backDir := filepath.Join(wd, "backups")

	// 检查并创建备份目录
	if err = system.CheckAndCreateDir(backDir); err != nil {
		// 如果检查或创建目录失败，则返回错误
		return "", err
	}

	// 返回备份目录的路径
	return backDir, nil
}

func getFromSource(file, way string) (string, error) {
	var levelFilePath string
	var err error

	// 处理http(s)://url的情况
	if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
		// http(s)://url
		levelFilePath, err = source.DownloadFromHttp(file, way)
		if err != nil {
			return "", errors.New("error downloading file: " + err.Error())
		}
	}
	// 处理k8s://namespace/pod/container:remotePath的情况
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
	// 处理docker://containerID(Name):remotePath的情况
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
	// 处理本地文件的情况
	else {
		// local file
		levelFilePath, err = source.CopyFromLocal(file, way)
		if err != nil {
			return "", errors.New("error copying file to temporary directory: " + err.Error())
		}
	}
	return levelFilePath, nil
}