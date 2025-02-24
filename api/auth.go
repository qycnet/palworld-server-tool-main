package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/qycnet/palworld-server-tool-main/internal/auth"
)

type LoginInfo struct {
	Password string `json:"password"`
}

// loginHandler godoc
// @Summary		Login
// @Description	Login
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			login_info	body		LoginInfo	true	"Login Info"
// @Success		200			{object}	SuccessResponse
// @Failure		400			{object}	ErrorResponse
// @Failure		401			{object}	ErrorResponse
// @Router			/api/login [post]
func loginHandler(c *gin.Context) {
	// 定义LoginInfo变量
	var loginInfo LoginInfo
	// 绑定JSON请求体到loginInfo变量
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		// 如果绑定失败，返回400错误码
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从配置文件中获取正确的密码
	correctPassword := viper.GetString("web.password")
	// 检查传入的密码是否正确
	if loginInfo.Password != correctPassword {
		// 如果密码不正确，返回401错误码
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 创建一个JWT令牌，设置过期时间为24小时后
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	// 签署JWT令牌
	tokenString, err := token.SignedString(auth.SecretKey)
	// 如果签署失败，返回400错误码
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法生成令牌"})
		return
	}
	// 返回生成的JWT令牌
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
