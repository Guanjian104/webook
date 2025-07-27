//go:build wireinject

package main

import (
    "github.com/Guanjian104/webook/internal/repository"
    "github.com/Guanjian104/webook/internal/repository/cache"
    "github.com/Guanjian104/webook/internal/repository/dao"
    "github.com/Guanjian104/webook/internal/service"
    "github.com/Guanjian104/webook/internal/web"
    "github.com/Guanjian104/webook/ioc"
    "github.com/gin-gonic/gin"
)
import "github.com/google/wire"

func InitWebServer() *gin.Engine {
    wire.Build(
        ioc.InitDB,
        ioc.InitRedis,
        dao.NewUserDAO,
        cache.NewCodeCache,
        cache.NewUserCache,
        repository.NewCachedUserRepository,
        repository.NewCodeRepository,
        ioc.InitSMSService,
        service.NewUserService,
        service.NewCodeService,
        web.NewUserHandler,
        ioc.InitGinMiddlewares,
        ioc.InitWebServer,
    )
    return gin.Default()
}
