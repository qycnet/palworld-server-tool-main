package api

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"github.com/qycnet/palworld-server-tool-main/internal/tool"
	"github.com/qycnet/palworld-server-tool-main/service"
)

type PlayerOrderBy string

const (
	OrderByLastOnline PlayerOrderBy = "last_online"
	OrderByLevel      PlayerOrderBy = "level"
)

// listOnlinePlayers godoc
//
//	@Summary		List Online Players
//	@Description	List Online Players
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//
//	@Success		200	{object}	[]database.OnlinePlayer
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/online_player [get]
func listOnlinePlayers(c *gin.Context) {
	// 获取在线玩家列表
	onlinePLayers, err := tool.ShowPlayers()
	if err != nil {
		// 如果出现错误，返回错误状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 将在线玩家列表存入数据库
	service.PutPlayersOnline(database.GetDB(), onlinePLayers)
	// 返回成功状态码和在线玩家列表
	c.JSON(http.StatusOK, onlinePLayers)
}

// putPlayers godoc
//
//	@Summary		Put Players
//	@Description	Put Players Only For SavSync,PlayerSync
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Param			players	body		[]database.Player	true	"Players"
//
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Router			/api/player [put]
func putPlayers(c *gin.Context) {
	// 定义一个空的Player切片
	var players []database.Player

	// 将请求体绑定到players变量中
	if err := c.ShouldBindJSON(&players); err != nil {
		// 如果绑定失败，返回400错误码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用service.PutPlayers函数，将players插入数据库
	if err := service.PutPlayers(database.GetDB(), players); err != nil {
		// 如果插入失败，返回400错误码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果插入成功，返回200状态码和成功信息
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// listPlayers godoc
//
//	@Summary		List Players
//	@Description	List Players
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//
//	@Param			order_by	query		PlayerOrderBy	false	"order by field"	enum(last_online,level)
//	@Param			desc		query		bool			false	"order by desc"
//
//	@Success		200			{object}	[]database.TersePlayer
//	@Failure		400			{object}	ErrorResponse
//	@Router			/api/player [get]
func listPlayers(c *gin.Context) {
	// 获取查询参数 order_by
	orderBy := c.Query("order_by")
	// 获取查询参数 desc
	desc := c.Query("desc")

	// 从数据库中获取玩家列表
	players, err := service.ListPlayers(database.GetDB())
	if err != nil {
		// 如果获取玩家列表失败，则返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据 order_by 参数对玩家列表进行排序
	if orderBy == "level" {
		// 按等级排序
		sort.Slice(players, func(i, j int) bool {
			if desc == "true" {
				// 如果 desc 参数为 true，则降序排序
				return players[i].Level > players[j].Level
			}
			// 否则升序排序
			return players[i].Level < players[j].Level
		})
	}

	// 根据 order_by 参数对玩家列表进行排序
	if orderBy == "last_online" {
		// 按最后上线时间排序
		sort.Slice(players, func(i, j int) bool {
			if desc == "true" {
				// 如果 desc 参数为 true，则降序排序
				return players[i].LastOnline.Sub(players[j].LastOnline) > 0
			}
			// 否则升序排序
			return players[i].LastOnline.Sub(players[j].LastOnline) < 0
		})
	}

	// 返回排序后的玩家列表
	c.JSON(http.StatusOK, players)
}

// getPlayer godoc
//
//	@Summary		Get Player
//	@Description	Get Player
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//
//	@Param			player_uid	path		string	true	"Player UID"
//
//	@Success		200			{object}	database.Player
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	EmptyResponse
//	@Router			/api/player/{player_uid} [get]
func getPlayer(c *gin.Context) {
	// 从数据库中获取玩家信息
	player, err := service.GetPlayer(database.GetDB(), c.Param("player_uid"))
	if err != nil {
		// 检查错误类型
		if err == service.ErrNoRecord {
			// 如果没有记录，返回404状态码
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		// 如果有其他错误，返回400状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 返回玩家信息，状态码为200
	c.JSON(http.StatusOK, player)
}

// kickPlayer godoc
//
//	@Summary		Kick Player
//	@Description	Kick Player
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			player_uid	path		string	true	"Player UID"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/player/{player_uid}/kick [post]
func kickPlayer(c *gin.Context) {
	// 获取URL参数中的玩家UID
	playerUid := c.Param("player_uid")
	// 根据玩家UID获取玩家信息
	player, err := service.GetPlayer(database.GetDB(), playerUid)
	if err != nil {
		// 如果错误为没有记录的错误
		if err == service.ErrNoRecord {
			// 返回404状态码，提示玩家未找到
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			return
		}
		// 返回400状态码，提示错误详情
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 调用工具函数踢出玩家
	err = tool.KickPlayer(fmt.Sprintf("steam_%s", player.SteamId))
	if err != nil {
		// 如果踢出玩家失败，返回400状态码，提示错误详情
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 踢出玩家成功，返回200状态码，提示操作成功
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// banPlayer godoc
//
//	@Summary		Ban Player
//	@Description	Ban Player
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			player_uid	path		string	true	"Player UID"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/player/{player_uid}/ban [post]
func banPlayer(c *gin.Context) {
	// 从请求中获取玩家UID
	playerUid := c.Param("player_uid")

	// 从数据库中获取玩家信息
	player, err := service.GetPlayer(database.GetDB(), playerUid)
	if err != nil {
		// 如果错误是未找到记录
		if err == service.ErrNoRecord {
			// 返回404错误，表示未找到玩家
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			return
		}
		// 返回400错误，并输出错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 封禁玩家
	err = tool.BanPlayer(fmt.Sprintf("steam_%s", player.SteamId))
	if err != nil {
		// 如果封禁失败，返回400错误，并输出错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回200成功状态，表示封禁成功
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// unbanPlayer godoc
//
//	@Summary		Unban Player
//	@Description	Unban Player
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			player_uid	path		string	true	"Player UID"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/player/{player_uid}/unban [post]
func unbanPlayer(c *gin.Context) {
	// 获取请求中的玩家UID
	playerUid := c.Param("player_uid")

	// 从数据库中获取玩家信息
	player, err := service.GetPlayer(database.GetDB(), playerUid)
	if err != nil {
		// 如果没有找到玩家记录
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			return
		}
		// 如果出现其他错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解除对玩家的封禁
	err = tool.UnBanPlayer(fmt.Sprintf("steam_%s", player.SteamId))
	if err != nil {
		// 如果出现错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 封禁解除成功，返回成功信息
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// addWhite godoc
//
//	@Summary		Add White List
//	@Description	Add White List
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			player_uid	path		string	true	"Player UID"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/api/whitelist [post]
func addWhite(c *gin.Context) {
	var player database.PlayerW
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddWhitelist(database.GetDB(), player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// listWhite godoc
//
//	@Summary		List White List
//	@Description	List White List
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]database.PlayerW
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/whitelist [get]
func listWhite(c *gin.Context) {
	// 从数据库中获取白名单玩家列表
	players, err := service.ListWhitelist(database.GetDB())
	if err != nil {
		// 如果获取白名单玩家列表时出错，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 如果获取成功，返回玩家列表
	c.JSON(http.StatusOK, players)
}

// removeWhite godoc
//
//	@Summary		Remove White List
//	@Description	Remove White List
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			player_uid	path		string	true	"Player UID"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/api/whitelist [delete]
func removeWhite(c *gin.Context) {
	// 定义一个 database.PlayerW 类型的变量 player
	var player database.PlayerW
	// 尝试将请求体绑定到 player 变量
	if err := c.ShouldBindJSON(&player); err != nil {
		// 如果绑定失败，返回状态码 400 和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 service.RemoveWhitelist 方法从白名单中移除 player
	if err := service.RemoveWhitelist(database.GetDB(), player); err != nil {
		// 如果移除失败，返回状态码 400 和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果成功移除，返回状态码 200 和成功信息
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// putWhite godoc
//
//	@Summary		Put White List
//	@Description	Put White List
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			players	body		[]database.PlayerW	true	"Players"
//
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Router			/api/whitelist [put]
func putWhite(c *gin.Context) {
	// 定义一个空的players切片
	var players []database.PlayerW

	// 绑定JSON数据到players切片
	if err := c.ShouldBindJSON(&players); err != nil {
		// 如果绑定失败，返回错误状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用service.PutWhitelist方法，将players数据存入白名单
	if err := service.PutWhitelist(database.GetDB(), players); err != nil {
		// 如果存入白名单失败，返回错误状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果一切顺利，返回成功状态码和成功信息
	c.JSON(http.StatusOK, gin.H{"success": true})
}
