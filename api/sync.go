package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"github.com/qycnet/palworld-server-tool-main/internal/task"
)

type From string

const (
	FromRest From = "rest"
	FromSav  From = "sav"
)

// syncData godoc
//
//	@Summary		Sync Data
//	@Description	Sync Data
//	@Tags			Sync
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			from	query		From	true	"from"	enum(rest,sav)
//
//	@Success		200		{object}	SuccessResponse
//	@Failure		401		{object}	ErrorResponse
//	@Router			/api/sync [post]
func syncData(c *gin.Context) {
	// 从请求中获取from参数
	from := c.Query("from")

	// 如果from参数为"rest"
	if from == "rest" {
		// 异步启动玩家数据同步任务
		go task.PlayerSync(database.GetDB())
		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	} else if from == "sav" {
		// 如果from参数为"sav"
		// 异步启动sav数据同步任务
		go task.SavSync()
		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	// 如果from参数既不是"rest"也不是"sav"，返回错误响应
	c.JSON(http.StatusOK, gin.H{"error": "invalid from"})
}
