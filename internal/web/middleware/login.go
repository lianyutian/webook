package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
		session := sessions.Default(c)
		userId := session.Get("userId")
		// 没有登录
		if userId == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := session.Get("updateTime")
		session.Options(sessions.Options{
			MaxAge: 10,
		})
		now := time.Now().UnixMilli()
		// 刚登陆还未设置刷新时间
		if updateTime == nil {
			session.Set("updateTime", now)
			err := session.Save()
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		updateTimeVal := updateTime.(int64)
		// 刷新时间
		if now-updateTimeVal > 5 {
			session.Set("updateTime", now)
			err := session.Save()
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
	}
}
