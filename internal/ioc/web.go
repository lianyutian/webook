package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func InitGin(modules []gin.HandlerFunc, userHandler *web.UserHandler) *gin.Engine {
	engine := gin.Default()
	engine.Use(modules...)
	userHandler.RegisterRoutes(engine)
	return engine
}

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 处理跨域问题
		// 配置 CORS 中间件
		corsMiddleware(),
		// 登录中间件
		loginMiddleware(),
	}
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                   // 允许的域
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的Header
		ExposeHeaders:    []string{"Content-Length", "x-jwt-token"},           // 公开的Header
		AllowCredentials: true,                                                // 允许携带凭证（如 cookies）
		MaxAge:           12 * time.Hour,                                      // 预检请求的缓存时间
	})
}

func loginMiddleware() gin.HandlerFunc {
	return middleware.NewLoginMiddlewareBuilder().
		IgnorePath("/users/login").
		IgnorePath("/users/signup").
		IgnorePath("/users/login_sms/code/send").
		IgnorePath("/users/login_sms").
		Build()
}
