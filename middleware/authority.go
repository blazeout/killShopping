package middleware

import (
	R "KillShopping/response"
	"KillShopping/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Admin 检测是否有权限操作
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, exists := c.Get("jwtUserInfo")
		if exists {
			info := userInfo.(utils.JwtUserInfo)
			if info.Authority == 2 {
				c.Next()
			} else {
				R.Response(c, http.StatusUnauthorized, "无权限", nil, http.StatusUnauthorized)
				c.Abort()
			}
		} else {
			R.Response(c, http.StatusUnauthorized, "无权限", nil, http.StatusUnauthorized)
			c.Abort()
		}
		return
	}
}
