package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		//session := sessions.Default(c)
		//userId := session.Get("userId")
		tokenString := c.GetHeader("authorization")
		// 没有登录
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var jwtSecret = []byte("your-secret-key")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证加密方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
		//updateTime := session.Get("updateTime")
		//session.Options(sessions.Options{
		//	MaxAge: 10,
		//})
		//now := time.Now().UnixMilli()
		//// 刚登陆还未设置刷新时间
		//if updateTime == nil {
		//	session.Set("updateTime", now)
		//	err := session.Save()
		//	if err != nil {
		//		c.AbortWithStatus(http.StatusInternalServerError)
		//	}
		//	return
		//}
		//
		//updateTimeVal := updateTime.(int64)
		//// 刷新时间
		//if now-updateTimeVal > 5 {
		//	session.Set("updateTime", now)
		//	err := session.Save()
		//	if err != nil {
		//		c.AbortWithStatus(http.StatusInternalServerError)
		//		return
		//	}
		//}
	}
}
