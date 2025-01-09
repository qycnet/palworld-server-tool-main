package api

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"github.com/qycnet/palworld-server-tool-main/internal/tool"
	"github.com/qycnet/palworld-server-tool-main/service"
)

type SendRconCommandRequest struct {
	UUID    string `json:"uuid"`
	Content string `json:"content"`
}

// sendRconCommand godoc
//
//	@Summary		Send Rcon Command
//	@Description	Send Rcon Command
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			command	body		SendRconCommandRequest	true	"Rcon Command"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Router			/api/rcon/send [post]
func sendRconCommand(c *gin.Context) {
	// 定义请求结构体变量
	var req SendRconCommandRequest
	// 绑定请求数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果绑定失败，返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从数据库中获取Rcon命令
	rcon, err := service.GetRconCommand(database.GetDB(), req.UUID)
	if err != nil {
		// 如果命令不存在，返回404错误响应
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rcon command not found"})
			return
		}
		// 如果获取命令失败，返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构造执行命令
	execCommand := fmt.Sprintf("%s %s", rcon.Command, req.Content)
	// 执行命令
	response, err := tool.CustomCommand(execCommand)
	if err != nil {
		// 如果执行命令失败，返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回命令执行结果
	c.JSON(http.StatusOK, gin.H{"message": response})
}

// importRconCommands godoc
//
//	@Summary		Import Rcon Commands
//	@Description	Import Rcon Commands from a TXT file
//	@Tags			Rcon
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file	formData	file	true	"Upload txt file"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon/import [post]
func importRconCommands(c *gin.Context) {
	// 从请求中获取上传的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		// 如果获取文件失败，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close() // 确保文件在操作完成后关闭

	// 使用bufio.Scanner逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行
		parts := strings.Split(line, ",") // 使用逗号分隔行内容
		if len(parts) < 2 {
			// 如果分隔后的部分少于2个，返回文件格式错误
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
			return
		}
		placeholder := ""
		if len(parts) >= 3 {
			placeholder = parts[2] // 如果分隔后的部分大于等于3个，取第三个部分作为占位符
		}
		rconCommand := database.RconCommand{
			Command:     parts[0], // 命令
			Remark:      parts[1], // 备注
			Placeholder: placeholder, // 占位符
		}
		if err := service.AddRconCommand(database.GetDB(), rconCommand); err != nil {
			// 如果添加命令失败，返回错误信息
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := scanner.Err(); err != nil {
		// 如果读取文件出错，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true}) // 操作成功，返回成功信息
}

// listRconCommand godoc
//
//	@Summary		List Rcon Commands
//	@Description	List Rcon Commands
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	[]database.RconCommandList
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/rcon [get]
func listRconCommand(c *gin.Context) {
	// 调用服务层函数获取所有Rcon命令
	rcons, err := service.ListRconCommands(database.GetDB())
	if err != nil {
		// 如果出现错误，返回400状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 如果没有错误，返回200状态码和Rcon命令列表
	c.JSON(http.StatusOK, rcons)
}

// addRconCommand godoc
//
//	@Summary		Add Rcon Command
//	@Description	Add Rcon Command
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			command	body		database.RconCommand	true	"Rcon Command"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon [post]
func addRconCommand(c *gin.Context) {
	// 声明一个 RconCommand 类型的变量 rcon
	var rcon database.RconCommand
	// 尝试将请求体中的 JSON 数据绑定到 rcon 变量中
	if err := c.ShouldBindJSON(&rcon); err != nil {
		// 如果绑定失败，返回 HTTP 400 错误和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 调用 service.AddRconCommand 方法将 rcon 添加到数据库中
	err := service.AddRconCommand(database.GetDB(), rcon)
	// 如果添加失败，返回 HTTP 400 错误和错误信息
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 如果添加成功，返回 HTTP 200 状态和成功信息
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// putRconCommand godoc
//
//	@Summary		Put Rcon Command
//	@Description	Put Rcon Command
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			uuid	path		string					true	"UUID"
//	@Param			command	body		database.RconCommand	true	"Rcon Command"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon/{uuid} [put]
func putRconCommand(c *gin.Context) {
	// 从URL参数中获取uuid
	uuid := c.Param("uuid")

	// 定义一个RconCommand类型的变量
	var rcon database.RconCommand

	// 将请求的JSON数据绑定到rcon变量中
	if err := c.ShouldBindJSON(&rcon); err != nil {
		// 如果绑定失败，返回400 Bad Request响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层的PutRconCommand方法，将命令保存到数据库
	err := service.PutRconCommand(database.GetDB(), uuid, rcon)
	if err != nil {
		// 如果保存失败，返回400 Bad Request响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果保存成功，返回200 OK响应
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// removeRconCommand godoc
//
//	@Summary		Remove Rcon Command
//	@Description	Remove Rcon Command
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			uuid	path		string	true	"UUID"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon/{uuid} [delete]
func removeRconCommand(c *gin.Context) {
	// 从请求参数中获取uuid
	uuid := c.Param("uuid")
	// 调用service层的RemoveRconCommand函数，传入数据库连接和uuid
	err := service.RemoveRconCommand(database.GetDB(), uuid)
	if err != nil {
		// 如果发生错误，返回400状态码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 如果执行成功，返回200状态码和成功信息
	c.JSON(http.StatusOK, gin.H{"success": true})
}
