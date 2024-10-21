package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	db := initDB()
	r := initWebServer()
	initUser(db, r)
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func initDB() *gorm.DB {
	dsn := "cloud:li.ming9518@tcp(116.198.217.158:3306)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initWebServer() *gin.Engine {
	r := gin.Default()

	// 处理跨域问题
	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                   // 允许的域
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的Header
		ExposeHeaders:    []string{"Content-Length", "x-jwt-token"},           // 公开的Header
		AllowCredentials: true,                                                // 允许携带凭证（如 cookies）
		MaxAge:           12 * time.Hour,                                      // 预检请求的缓存时间
	}))

	// 登录校验
	r.Use(
		middleware.NewLoginMiddlewareBuilder().
			IgnorePath("/users/login").
			IgnorePath("/users/signup").
			Build(),
	)

	return r
}

func initUser(db *gorm.DB, r *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	uh := web.NewUserHandler(us)
	uh.RegisterRoutes(r)
}
