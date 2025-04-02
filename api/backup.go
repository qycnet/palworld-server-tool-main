package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"github.com/qycnet/palworld-server-tool-main/internal/tool"
	"github.com/qycnet/palworld-server-tool-main/service"
)

// listBackups godoc
//
//	@Summary		List backups within a specified time range
//	@Description	List all backups or backups within a specific time range.
//	@Tags			backup
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			startTime	query		int	false	"Start time of the backup range in timestamp"
//	@Param			endTime		query		int	false	"End time of the backup range in timestamp"
//	@Success		200			{array}		database.Backup
//	@Failure		400			{object}	ErrorResponse
//	@Router			/api/backup [get]
func listBackups(c *gin.Context) {
	var startTimestamp, endTimestamp int64
	var startTime, endTime time.Time
	var err error

	// 从请求中获取起始时间和结束时间字符串
	startTimeStr, endTimeStr := c.Query("startTime"), c.Query("endTime")

	// 处理起始时间
	if startTimeStr != "" {
		// 将起始时间字符串转换为整数
		startTimestamp, err = strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			// 如果转换失败，返回错误响应
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间"})
			return
		}
		// 将起始时间整数转换为时间类型
		startTime = time.Unix(0, startTimestamp*int64(time.Millisecond))
	}

	// 处理结束时间
	if endTimeStr != "" {
		// 将结束时间字符串转换为整数
		endTimestamp, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			// 如果转换失败，返回错误响应
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间"})
			return
		}
		// 将结束时间整数转换为时间类型
		endTime = time.Unix(0, endTimestamp*int64(time.Millisecond))
	}

	// 调用服务方法获取备份列表
	backups, err := service.ListBackups(database.GetDB(), startTime, endTime)
	if err != nil {
		// 如果获取备份列表失败，返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回备份列表
	c.JSON(http.StatusOK, backups)
}

// downloadBackup godoc
//
//	@Summary		Download Backup
//	@Description	Download a backup
//	@Tags			backup
//	@Accept			json
//	@Produce		application/octet-stream
//	@Security		ApiKeyAuth
//	@Param			backup_id	path		string	true	"Backup ID"
//	@Success		200			{file}		"Backupfile"
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/backup/{backup_id} [get]
func downloadBackup(c *gin.Context) {
	// 获取URL参数中的backup_id
	backupId := c.Param("backup_id")

	// 获取备份信息
	backup, err := service.GetBackup(database.GetDB(), backupId)
	if err != nil {
		// 如果错误类型为没有找到记录，则返回404状态码
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		// 如果出现其他错误，则返回400状态码
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取备份目录
	backupDir, err := tool.GetBackupDir()
	if err != nil {
		// 如果出现错误，则返回500状态码
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置响应头中的Content-Disposition，以便浏览器下载文件
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", backup.Path))
	// 发送文件到客户端
	c.File(filepath.Join(backupDir, backup.Path))
}

// deleteBackup godoc
//
//	@Summary		Delete Backup
//	@Description	Delete a backup
//	@Tags			backup
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			backup_id	path		string	true	"Backup ID"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Router			/api/backup/{backup_id} [delete]
func deleteBackup(c *gin.Context) {
	// 获取请求参数中的backup_id
	backupId := c.Param("backup_id")

	// 定义并初始化一个Backup类型的变量backup
	var backup database.Backup

	// 调用service.GetBackup函数获取指定ID的备份信息
	backup, err := service.GetBackup(database.GetDB(), backupId)
	if err != nil {
		// 如果错误是service.ErrNoRecord，表示没有找到记录
		if err == service.ErrNoRecord {
			// 返回404状态码和空JSON对象
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		// 返回400状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用service.DeleteBackup函数删除指定ID的备份
	if err := service.DeleteBackup(database.GetDB(), backupId); err != nil {
		// 返回400状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用tool.GetBackupDir函数获取备份目录
	backupDir, err := tool.GetBackupDir()
	if err != nil {
		// 返回500状态码和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 拼接备份文件的完整路径并尝试删除该文件
	err = os.Remove(filepath.Join(backupDir, backup.Path))
	if err != nil {
		// 返回400状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回200状态码和成功信息
	c.JSON(http.StatusOK, gin.H{"success": true})
}
