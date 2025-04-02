package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
	"github.com/qycnet/palworld-server-tool-main/internal/tool"
)

type ServerInfo struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

type ServerMetrics struct {
	ServerFps        int     `json:"server_fps"`
	CurrentPlayerNum int     `json:"current_player_num"`
	ServerFrameTime  float64 `json:"server_frame_time"`
	MaxPlayerNum     int     `json:"max_player_num"`
	Uptime           int     `json:"uptime"`
	Days             int     `json:"days"`
}

type BroadcastRequest struct {
	Message string `json:"message"`
}

type ShutdownRequest struct {
	Seconds int    `json:"seconds"`
	Message string `json:"message"`
}

type ServerToolResponse struct {
	Version string `json:"version"`
	Latest  string `json:"latest"`
}

// getServerTool godoc
//
//	@Summary		Get PalWorld Server Tool
//	@Description	Get PalWorld Server Tool
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ServerToolResponse
//	@Router			/api/server/tool [get]
func getServerTool(c *gin.Context) {
	// 从上下文中获取版本号，如果不存在则默认为"Unknown"
	version, exists := c.Get("version")
	if !exists {
		// 如果版本号不存在，则设置为"Unknown"
		version = "Unknown"
	}

	// 获取最新版本号
	latest, err := tool.GetLatestTag()
	if err != nil {
		// 如果获取最新版本号失败，则记录错误日志
		logger.Errorf("%v\n", err)
	}

	// 如果最新版本号为空，则从Gitee获取最新版本号
	if latest == "" {
		latest, err = tool.GetLatestTagFromGitee()
		if err != nil {
			// 如果从Gitee获取最新版本号失败，则记录错误日志
			logger.Errorf("%v\n", err)
		}
	}

	// 将版本号和最新版本号以JSON格式返回给客户端
	c.JSON(http.StatusOK, gin.H{"version": version, "latest": latest})
}

// getServer godoc
//
//	@Summary		Get Server Info
//	@Description	Get Server Info
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ServerInfo
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/server [get]
func getServer(c *gin.Context) {
	// 获取系统信息
	info, err := tool.Info()
	if err != nil {
		// 如果获取信息时发生错误，返回错误状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: add system psutil info
	// 返回系统信息，状态码为200
	c.JSON(http.StatusOK, &ServerInfo{info["version"], info["name"]})
}

// getServerMetrics godoc
//
//	@Summary		Get Server Metrics
//	@Description	Get Server Metrics
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ServerMetrics
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/server/metrics [get]
func getServerMetrics(c *gin.Context) {
	// 调用 tool.Metrics() 获取服务器指标
	metrics, err := tool.Metrics()
	if err != nil {
		// 如果发生错误，则返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 将获取到的指标封装成 ServerMetrics 结构体并返回
	c.JSON(http.StatusOK, &ServerMetrics{
		// 服务器帧率
		ServerFps:        metrics["server_fps"].(int),
		// 当前玩家数量
		CurrentPlayerNum: metrics["current_player_num"].(int),
		// 服务器帧时间
		ServerFrameTime:  metrics["server_frame_time"].(float64),
		// 最大玩家数量
		MaxPlayerNum:     metrics["max_player_num"].(int),
		// 服务器运行时间
		Uptime:           metrics["uptime"].(int),
		// 服务器运行天数
		Days:             metrics["days"].(int),
	})
}

// publishBroadcast godoc
//
//	@Summary		Publish Broadcast
//	@Description	Publish Broadcast
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			broadcast	body		BroadcastRequest	true	"Broadcast"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/api/server/broadcast [post]
func publishBroadcast(c *gin.Context) {
	// 定义一个BroadcastRequest变量
	var req BroadcastRequest
	// 尝试将请求体绑定到req变量中
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果绑定失败，则返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 验证消息内容
	if err := validateMessage(req.Message); err != nil {
		// 如果验证失败，则返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 广播消息
	if err := tool.Broadcast(req.Message); err != nil {
		// 如果广播失败，则返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 如果所有操作都成功，则返回成功响应
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// shutdownServer godoc
//
//	@Summary		Shutdown Server
//	@Description	Shutdown Server
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			shutdown	body		ShutdownRequest	true	"Shutdown"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/api/server/shutdown [post]
func shutdownServer(c *gin.Context) {
	// 定义一个 ShutdownRequest 类型的变量 req
	var req ShutdownRequest
	// 将请求体绑定到 req 变量上，如果绑定失败则返回错误信息
	if err := c.ShouldBindJSON(&req); err != nil {
		// 返回 HTTP 状态码 400 和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证 req.Message 是否合法，如果不合法则返回错误信息
	if err := validateMessage(req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果 req.Seconds 为 0，则将其设置为 60
	if req.Seconds == 0 {
		req.Seconds = 60
	}

	// 调用 tool.Shutdown 函数进行关机操作，如果失败则返回错误信息
	if err := tool.Shutdown(req.Seconds, req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回 HTTP 状态码 200 和成功信息
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func validateMessage(message string) error {
	// 判断消息是否为空
	if message == "" {
		// 如果消息为空，则返回错误
		return errors.New("message cannot be empty")
	}
	// 如果消息不为空，则返回nil
	return nil
}
