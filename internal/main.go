package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/service/sms/memory"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	db := initDB()
	r := initWebServer()
	client := initRedis()
	initUser(db, r, client)
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func initDB() *gorm.DB {
	dsn := "cloud:li.ming9518@tcp(116.198.217.158:3306)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			LogLevel:      logger.Info,            // 日志级别
			SlowThreshold: 200 * time.Millisecond, // 慢查询阈值
			// IgnoreRecordNotFoundError: true, // 忽略未找到记录的错误
			Colorful: true, // 颜色输出
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
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
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").
			Build(),
	)

	return r
}

func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "116.198.217.158:6379",
		Password: "MIIEowIBAAKCAQEAwG90ULRHmAXFXQzZSwleoYts2+bCzUvqhhqtGiv/F5kUsETY", // no password set
		DB:       0,                                                                  // use default DB
	})
	return rdb
}

func initUser(db *gorm.DB, r *gin.Engine, client *redis.Client) {
	ud := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(client)
	ur := repository.NewUserRepository(ud, userCache)
	us := service.NewUserService(ur)

	codeCache := cache.NewCodeCache(client)
	codeRepository := repository.NewCodeRepository(codeCache)
	memService := memory.NewService()
	codeService := service.NewCodeService(codeRepository, memService, "")

	uh := web.NewUserHandler(us, codeService)
	uh.RegisterRoutes(r)
}
