package middleware

import (
    "github.com/Guanjian104/webook/internal/web"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "log"
    "net/http"
    "strings"
    "time"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        path := ctx.Request.URL.Path
        if path == "/users/signup" ||
            path == "/users/login" ||
            path == "/login_sms/code/send" ||
            path == "/login_sms" {
            return
        }
        authCode := ctx.GetHeader("Authorization")
        if authCode == "" {
            ctx.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        segs := strings.Split(authCode, " ")
        if len(segs) != 2 {
            ctx.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        tokenStr := segs[1]
        var uc web.UserClaims
        token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
            return web.JWTKey, nil
        })
        if err != nil {
            ctx.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        if token == nil || !token.Valid {
            ctx.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        if uc.UserAgent != ctx.GetHeader("User-Agent") {
            ctx.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        expireTime := uc.ExpiresAt
        // 剩余过期时间 < 50s 刷新 token
        if expireTime.Sub(time.Now()) < time.Second*50 {
            uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
            tokenStr, err = token.SignedString(web.JWTKey)
            ctx.Header("x-jwt-token", tokenStr)
            if err != nil {
                log.Println(err)
            }
        }

        ctx.Set("user", uc)
    }
}
