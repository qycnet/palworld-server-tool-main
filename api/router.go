package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/qycnet/palworld-server-tool-main/internal/auth"
)

type SuccessResponse struct {
	Success bool `json:"success"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type EmptyResponse struct{}

func ignoreLogPrefix(path string) bool {
	// 定义要忽略的前缀列表
	prefixes := []string{"/swagger/", "/assets/", "/favicon.ico", "/map"}
	for _, prefix := range prefixes {
		// 检查路径是否以某个前缀开头
		if strings.HasPrefix(path, prefix) {
			// 如果路径以某个前缀开头，则返回true
			return true
		}
	}
	// 如果路径不以任何前缀开头，则返回false
	return false
}

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 判断是否忽略日志前缀
		if !ignoreLogPrefix(param.Path) {
			// 获取状态码颜色
			statusColor := param.StatusCodeColor()
			// 获取方法颜色
			methodColor := param.MethodColor()
			// 获取重置颜色
			resetColor := param.ResetColor()
			// 格式化日志输出
			return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				param.ErrorMessage,
			)
		}
		// 如果忽略日志前缀，则返回空字符串
		return ""
	})
}

func RegisterRouter(r *gin.Engine) {
	// 使用中间件
	r.Use(Logger(), gin.Recovery())

	// 注册登录路由
	r.POST("/api/login", loginHandler)
	// 注册Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 创建一个新的API组
	apiGroup := r.Group("/api")

	// 创建匿名访问的路由组
	anonymousGroup := apiGroup.Group("")
	{
		// 获取服务器信息
		anonymousGroup.GET("/server", getServer)
		// 获取服务器工具信息
		anonymousGroup.GET("/server/tool", getServerTool)
		// 获取服务器性能指标
		anonymousGroup.GET("/server/metrics", getServerMetrics)
		// 获取玩家列表
		anonymousGroup.GET("/player", listPlayers)
		// 获取指定玩家的信息
		anonymousGroup.GET("/player/:player_uid", getPlayer)
		// 获取在线玩家列表
		anonymousGroup.GET("/online_player", listOnlinePlayers)
		// 获取公会列表
		anonymousGroup.GET("/guild", listGuilds)
		// 获取指定管理员玩家的公会信息
		anonymousGroup.GET("/guild/:admin_player_uid", getGuild)
	}

	// 创建需要认证的路由组
	authGroup := apiGroup.Group("")
	authGroup.Use(auth.JWTAuthMiddleware())
	{
		// 获取玩家列表
		//authGroup.GET("/player", listPlayers)
		// 获取指定玩家的信息
		//authGroup.GET("/player/:player_uid", getPlayer)
		// 获取在线玩家列表
		//authGroup.GET("/online_player", listOnlinePlayers)
		// 获取公会列表
		//authGroup.GET("/guild", listGuilds)
		// 获取指定管理员玩家的公会信息
		//authGroup.GET("/guild/:admin_player_uid", getGuild)
		// 发布广播消息
		authGroup.POST("/server/broadcast", publishBroadcast)
		// 关闭服务器
		authGroup.POST("/server/shutdown", shutdownServer)
		// 更新玩家信息
		authGroup.PUT("/player", putPlayers)
		// 踢出指定玩家
		authGroup.POST("/player/:player_uid/kick", kickPlayer)
		// 封禁指定玩家
		authGroup.POST("/player/:player_uid/ban", banPlayer)
		// 解封指定玩家
		authGroup.POST("/player/:player_uid/unban", unbanPlayer)
		// 更新公会信息
		authGroup.PUT("/guild", putGuilds)
		// 同步数据
		authGroup.POST("/sync", syncData)
		// 获取白名单列表
		authGroup.GET("/whitelist", listWhite)
		// 添加白名单
		authGroup.POST("/whitelist", addWhite)
		// 删除白名单
		authGroup.DELETE("/whitelist", removeWhite)
		// 更新白名单
		authGroup.PUT("/whitelist", putWhite)
		// 获取RCON命令列表
		authGroup.GET("/rcon", listRconCommand)
		// 添加RCON命令
		authGroup.POST("/rcon", addRconCommand)
		// 导入RCON命令
		authGroup.POST("/rcon/import", importRconCommands)
		// 发送RCON命令
		authGroup.POST("/rcon/send", sendRconCommand)
		// 更新指定UUID的RCON命令
		authGroup.PUT("/rcon/:uuid", putRconCommand)
		// 删除指定UUID的RCON命令
		authGroup.DELETE("/rcon/:uuid", removeRconCommand)
		// 获取备份列表
		authGroup.GET("/backup", listBackups)
		// 下载指定备份
		authGroup.GET("/backup/:backup_id", downloadBackup)
		// 删除指定备份
		authGroup.DELETE("/backup/:backup_id", deleteBackup)
	}
}
