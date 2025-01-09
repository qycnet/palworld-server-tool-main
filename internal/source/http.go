package source

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
	"github.com/qycnet/palworld-server-tool-main/internal/system"
)

func DownloadFromHttp(url, way string) (string, error) {
	// 打印日志信息，表示开始从指定URL下载文件
	logger.Infof("downloading sav.zip from %s\n", url)
	// 发送HTTP GET请求获取文件
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	// 确保响应体在使用完毕后关闭
	defer resp.Body.Close()

	// 生成一个UUID
	uuid := uuid.New().String()
	// 构建临时文件路径
	tempPath := filepath.Join(os.TempDir(), "palworldsav-http-"+way+"-"+uuid)
	// 获取临时文件的绝对路径
	absPath, err := filepath.Abs(tempPath)
	if err != nil {
		return "", err
	}

	// 清理并创建目录
	if err = system.CleanAndCreateDir(absPath); err != nil {
		return "", err
	}

	// 构建临时ZIP文件路径
	tempZipFilePath := filepath.Join(absPath, "sav.zip")
	// 确保临时ZIP文件在使用完毕后删除
	defer os.Remove(tempZipFilePath)

	// 创建临时ZIP文件
	zipOut, err := os.Create(tempZipFilePath)
	if err != nil {
		return "", err
	}
	// 确保ZIP文件在使用完毕后关闭
	defer zipOut.Close()

	// 将响应体复制到ZIP文件中
	_, err = io.Copy(zipOut, resp.Body)
	if err != nil {
		return "", err
	}

	// 解压ZIP文件到指定目录
	err = system.UnzipDir(tempZipFilePath, absPath)
	if err != nil {
		return "", err
	}
	// 构建Level.sav文件的路径
	levelFilePath := filepath.Join(absPath, "Level.sav")
	// 打印日志信息，表示文件下载和解压完成
	logger.Info("sav.zip downloaded and extracted\n")
	return levelFilePath, nil
}
