package main

import (
    "github.com/Guanjian104/webook/config"
    "github.com/Guanjian104/webook/internal/repository"
    "github.com/Guanjian104/webook/internal/repository/cache"
    "github.com/Guanjian104/webook/internal/repository/dao"
    "github.com/Guanjian104/webook/internal/service"
    "github.com/Guanjian104/webook/internal/web"
    "github.com/Guanjian104/webook/internal/web/middleware"
    "github.com/Guanjian104/webook/pkg/ginx/middleware/ratelimit"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "net/http"
    "strings"
    "time"
)

func main() {
    db := initDB()

    cmd := initRedis()

    server := initWebServer()

    initUserHdl(db, cmd, server)

    // server := gin.Default()
    server.GET("/hello", func(ctx *gin.Context) {
        ctx.String(http.StatusOK, "hello world")
    })
    server.Run(":8080")
}

func initUserHdl(db *gorm.DB, cmd redis.Cmdable, server *gin.Engine) {
    ud := dao.NewUserDAO(db)
    uc := cache.NewUserCache(cmd)
    ur := repository.NewUserRepository(ud, uc)
    us := service.NewUserService(ur)
    hdl := web.NewUserHandler(us)
    hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
    db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
    if err != nil {
        panic(err)
    }

    err = dao.InitTables(db)
    if err != nil {
        panic(err)
    }

    return db
}

func initRedis() redis.Cmdable {
    return redis.NewClient(&redis.Options{
        Addr: config.Config.Redis.Addr,
    })
}

func initWebServer() *gin.Engine {
    server := gin.Default()

    server.Use(cors.New(cors.Config{
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
    }))

    redisClient := redis.NewClient(&redis.Options{
        Addr: config.Config.Redis.Addr,
    })
    server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

    useJWT(server)
    // useSession(server)

    return server
}

func useJWT(server *gin.Engine) {
    login := middleware.LoginJWTMiddlewareBuilder{}
    server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
    login := &middleware.LoginMiddlewareBuilder{}
    store := cookie.NewStore([]byte("secret"))
    server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
