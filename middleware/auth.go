package middleware

import (
	R "KillShopping/response"
	"KillShopping/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Auth 使用JWT用来进行用户信息验证的, 检测是否登录
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) < 7 {
			// 阻止调用后面的处理函数
			c.Abort()
			R.Response(c, http.StatusUnauthorized, "未登录", nil, http.StatusUnauthorized)
			return
		}
		jwtUserInfo := utils.JwtUserInfo{}
		// 检验JWT Token是否正确
		err := jwtUserInfo.ParseToken(token[7:])
		if err != nil {
			R.Response(c, http.StatusUnauthorized, "未登录", nil, http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Set("jwtUserInfo", jwtUserInfo)
		// 检测成功便是已经登录, 继续调用后面的处理函数
		c.Next()
		return
	}
}
