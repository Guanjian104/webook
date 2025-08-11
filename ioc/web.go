package ioc

import (
    "strings"
    "time"

    "github.com/Guanjian104/webook/internal/web"
    "github.com/Guanjian104/webook/internal/web/middleware"
    "github.com/Guanjian104/webook/pkg/ginx/middleware/ratelimit"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
    server := gin.Default()
    server.Use(mdls...)
    userHdl.RegisterRoutes(server)
    return server
}

func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
    return []gin.HandlerFunc{
        cors.New(cors.Config{
            AllowCredentials: true,
            AllowHeaders:     []string{"Content-Type", "Authorization"},
            ExposeHeaders:    []string{"x-jwt-token"},
            AllowOriginFunc: func(origin string) bool {
                if strings.HasPrefix(origin, "http://localhost") {
                    return true
                }
                return strings.Contains(origin, "your_company.com")
            },
            MaxAge: 12 * time.Hour,
        }),
        func(ctx *gin.Context) {
            println("这是我的 Middleware")
        },
        ratelimit.NewBuilder(redisClient, time.Second, 1000).Build(),
        (&middleware.LoginJWTMiddlewareBuilder{}).CheckLogin(),
    }
}
