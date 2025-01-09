package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
	"github.com/qycnet/palworld-server-tool-main/internal/source"
	"github.com/qycnet/palworld-server-tool-main/internal/system"
)

var (
	port     int
	savedDir string
)

func main() {
	// 设置命令行参数
	flag.IntVar(&port, "port", 8081, "port")
	flag.StringVar(&savedDir, "d", "", "Directory containing Level.sav file")
	flag.Parse()

	// 使用 viper 绑定环境变量和设置默认值
	viper.BindEnv("saved_dir", "SAVED_DIR")
	viper.SetDefault("port", port)
	viper.SetDefault("saved_dir", savedDir)
	savedDir = viper.GetString("saved_dir")

	// 设置 Gin 的运行模式为发布模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 处理 GET 请求 /sync
	r.GET("/sync", func(c *gin.Context) {
		// 从本地复制 Level.sav 文件
		levelFile, err := source.CopyFromLocal(savedDir, "agent")
		if err != nil {
			logger.Errorf("Failed to get directory include Level.sav: %v\n", err)
			os.Exit(1)
		}
		// 获取缓存目录
		cacheDir := filepath.Dir(levelFile)
		defer os.RemoveAll(cacheDir)

		// 创建 zip 文件
		cacheFile := cacheDir + ".zip"
		err = system.ZipDir(cacheDir, cacheFile)
		if err != nil {
			logger.Errorf("Failed to create zip: %v\n", err)
			c.Redirect(http.StatusFound, "/404")
			return
		}
		defer os.Remove(cacheFile)

		// 设置响应头并发送 zip 文件
		c.Header("Content-Disposition", "attachment; filename=sav.zip")
		c.File(cacheFile)
	})

	// 打印日志，表示 PST-Agent 正在监听端口
	logger.Infof("PST-Agent Listening on port %d\n", port)

	// 创建一个信号通道，监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 运行 Gin 服务器
	go func() {
		if err := r.Run(":" + strconv.Itoa(port)); err != nil {
			logger.Errorf("Failed to start agent: %v\n", err)
		}
	}()

	// 等待接收到系统信号
	<-sigChan

	// 打印日志，表示 PST-Agent 已优雅停止
	logger.Info("PST-Agent gracefully stopped\n")
}
