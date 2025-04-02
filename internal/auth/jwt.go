package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var SecretKey = []byte(viper.GetString("web.password"))

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Authorization
		authHeader := c.GetHeader("Authorization")

		// 定义前缀
		prefixBearer := "Bearer "
		prefixJWT := "JWT "

		// 初始化 token 字符串
		var tokenString string

		// 判断请求头是否以 Bearer 开头
		if strings.HasPrefix(authHeader, prefixBearer) {
			// 如果是，则去除前缀并赋值给 tokenString
			tokenString = strings.TrimPrefix(authHeader, prefixBearer)
		} else if strings.HasPrefix(authHeader, prefixJWT) {
			// 判断请求头是否以 JWT 开头
			// 如果是，则去除前缀并赋值给 tokenString
			tokenString = strings.TrimPrefix(authHeader, prefixJWT)
		} else {
			// 如果请求头既不是 Bearer 开头也不是 JWT 开头
			// 则返回未授权状态，并附带错误信息
			//c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - token missing"})
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权 - 令牌缺失"})
			return
		}

		// 解析 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 判断 token 的签名方法是否为 HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// 如果不是，则返回错误
				return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
			}
			// 返回密钥
			return SecretKey, nil
		})

		// 如果解析 token 时出现错误
		if err != nil {
			// 则返回未授权状态，并附带错误信息
			//c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - invalid token"})
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权 - 令牌无效"})
			return
		}

		// 判断 token 是否有效
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 如果有效，则将 claims 存入上下文
			c.Set("claims", claims)
		} else {
			// 如果无效，则返回未授权状态，并附带错误信息
			//c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - invalid claims"})
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权 - 无效的声明"})
			return
		}

		// 继续执行下一个中间件或处理器
		c.Next()
	}
}

func GenerateToken() (string, error) {
	// 创建一个新的JWT令牌，使用HS256签名方法，并设置过期时间为当前时间加上24小时
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// 使用SecretKey对令牌进行签名，并转换为字符串
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		// 如果签名失败，返回空字符串和错误
		return "", err
	}

	// 返回签名后的令牌字符串和nil错误
	return tokenString, nil
}
