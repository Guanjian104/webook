package main

import (
    "github.com/Guanjian104/webook/config"
    "github.com/Guanjian104/webook/internal/repository"
    "github.com/Guanjian104/webook/internal/repository/cache"
    "github.com/Guanjian104/webook/internal/repository/dao"
    "github.com/Guanjian104/webook/internal/service"
    "github.com/Guanjian104/webook/internal/service/sms"
    "github.com/Guanjian104/webook/internal/service/sms/localsms"
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

    redisClient := initRedis()

    codeSvc := initCodeSvc(redisClient)

    server := initWebServer(redisClient)

    initUserHdl(db, redisClient, codeSvc, server)

    // server := gin.Default()
    server.GET("/hello", func(ctx *gin.Context) {
        ctx.String(http.StatusOK, "hello world")
    })
    server.Run(":8080")
}

func initUserHdl(db *gorm.DB, redisClient redis.Cmdable, codeSvc *service.CodeService, server *gin.Engine) {
    ud := dao.NewUserDAO(db)
    uc := cache.NewUserCache(redisClient)
    ur := repository.NewUserRepository(ud, uc)
    us := service.NewUserService(ur)
    hdl := web.NewUserHandler(us, codeSvc)
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

func initCodeSvc(redisClient redis.Cmdable) *service.CodeService {
    cc := cache.NewCodeCache(redisClient)
    crepo := repository.NewCodeRepository(cc)
    return service.NewCodeService(crepo, initSmsMemoryService())
}

func initSmsMemoryService() sms.Service {
    return localsms.NewService()
}

func initWebServer(redisClient redis.Cmdable) *gin.Engine {
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
