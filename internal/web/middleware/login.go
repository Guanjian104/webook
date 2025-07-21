package middleware

import (
    "encoding/gob"
    "net/http"
    "time"

    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
    // 注册时间类型
    gob.Register(time.Now())
    return func(ctx *gin.Context) {
        path := ctx.Request.URL.Path
        if path == "/users/signup" ||
            path == "/users/login" ||
            path == "/login_sms/code/send" ||
            path == "/login_sms" {
            return
        }
        sess := sessions.Default(ctx)
        userId := sess.Get("userId")
        if userId == nil {
            ctx.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        now := time.Now()
        const updateTimeKey = "update_time"
        val := sess.Get(updateTimeKey)
        lastUpdateTime, ok := val.(time.Time)
        if val == nil || !ok || now.Sub(lastUpdateTime) > time.Minute {
            // first checklogin
            sess.Set(updateTimeKey, now)
            sess.Set("userId", userId)
            err := sess.Save()
            if err != nil {
                // 记录日志
                println(err)
            }
        }
    }
}
