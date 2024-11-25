//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/ioc"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 基础的三方依赖
		ioc.InitDB, ioc.InitRedis,
		// 初始化DAO
		dao.NewUserDAO,

		cache.NewCodeCache,
		cache.NewUserCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,

		ioc.InitSMSService,

		web.NewUserHandler,

		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
