package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
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
		// 不需要登录校验的接口路径
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		// 校验 token
		tokenString := c.GetHeader("Authorization")

		// 没有登录
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		slices := strings.Split(tokenString, " ")
		if len(slices) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := slices[1]
		userClaims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, userClaims, func(token *jwt.Token) (interface{}, error) {
			// 验证加密方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return web.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		expiresAt := userClaims.ExpiresAt
		now := time.Now()
		if expiresAt.Sub(now) < time.Second*50 {
			userClaims.ExpiresAt = jwt.NewNumericDate(now)
			tokenStr, err = token.SignedString(web.JwtSecret)
			if err != nil {
				log.Println("jwt signing error:", err)
			}
			c.Header("x-jwt-token", tokenStr)
		}

		c.Set("userId", userClaims.Uid)

	}
}
