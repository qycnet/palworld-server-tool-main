package api

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"github.com/qycnet/palworld-server-tool-main/service"
)

// putGuilds godoc
//
//	@Summary		Put Guilds
//	@Description	Put Guilds Only For SavSync
//	@Tags			Guild
//	@Accept			json
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Param			guilds	body		[]database.Guild	true	"Guilds"
//
//	@Success		200		{object}	SuccessResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/guild [put]
func putGuilds(c *gin.Context) {
	// 定义一个存储多个数据库公会的切片
	var guilds []database.Guild

	// 尝试将请求体绑定到guilds切片中
	if err := c.ShouldBindJSON(&guilds); err != nil {
		// 如果绑定失败，返回400 Bad Request错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 将公会信息保存到数据库中
	if err := service.PutGuilds(database.GetDB(), guilds); err != nil {
		// 如果保存失败，返回400 Bad Request错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存成功，返回200 OK状态码和成功信息
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// listGuilds godoc
//
//	@Summary		List Guilds
//	@Description	List Guilds
//	@Tags			Guild
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]database.Guild
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/guild [get]
func listGuilds(c *gin.Context) {
	// 从数据库中获取公会列表
	guilds, err := service.ListGuilds(database.GetDB())
	if err != nil {
		// 如果获取公会列表失败，则返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 默认按大本营等级排序
	// default sort by base_camp_level
	sort.Slice(guilds, func(i, j int) bool {
		return guilds[i].BaseCampLevel > guilds[j].BaseCampLevel
	})

	// 返回排序后的公会列表
	c.JSON(http.StatusOK, guilds)
}

// getGuild godoc
//
//	@Summary		Get Guild
//	@Description	Get Guild
//	@Tags			Guild
//	@Accept			json
//	@Produce		json
//	@Param			admin_player_uid	path		string	true	"Admin Player UID"
//	@Success		200					{object}	database.Guild
//	@Failure		400					{object}	ErrorResponse
//	@Failure		404					{object}	EmptyResponse
//	@Router			/api/guild/{admin_player_uid} [get]
func getGuild(c *gin.Context) {
	// 从数据库中获取公会信息
	guild, err := service.GetGuild(database.GetDB(), c.Param("admin_player_uid"))
	if err != nil {
		// 如果错误为没有记录的错误
		if err == service.ErrNoRecord {
			// 返回404状态码和空响应体
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		// 返回400状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 返回200状态码和公会信息
	c.JSON(http.StatusOK, guild)
}
