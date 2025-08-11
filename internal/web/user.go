package web

import (
    "net/http"

    "time"
    "unicode/utf8"

    "github.com/Guanjian104/webook/internal/domain"
    "github.com/Guanjian104/webook/internal/service"
    regexp "github.com/dlclark/regexp2"
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    jwt "github.com/golang-jwt/jwt/v5"
)

const (
    emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
    passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
    bizLogin             = "login"
)

type UserHandler struct {
    emailRexExp    *regexp.Regexp
    passwordRexExp *regexp.Regexp
    svc            service.UserService
    codeSvc        service.CodeService
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
    return &UserHandler{
        emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
        passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
        svc:            svc,
        codeSvc:        codeSvc,
    }
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
    ug := server.Group("/users")
    ug.POST("/signup", h.SignUp)
    // ug.POST("/login", h.Login)
    ug.POST("/login", h.LoginJWT)
    ug.POST("/edit", h.Edit)
    ug.GET("/profile/", h.Profile)

    // 手机验证码登录相关功能
    ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
    ug.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
    type SignUpReq struct {
        Email           string `json:"email"`
        Password        string `json:"password"`
        ConfirmPassword string `json:"confirmPassword"`
    }

    var req SignUpReq
    if err := ctx.Bind(&req); err != nil {
        return
    }

    isEmail, err := h.emailRexExp.MatchString(req.Email)
    if err != nil {
        ctx.String(http.StatusOK, "系统错误")
        return
    }
    if !isEmail {
        ctx.String(http.StatusOK, "非法邮箱格式")
        return
    }

    if req.Password != req.ConfirmPassword {
        ctx.String(http.StatusOK, "两次输入密码不对")
        return
    }

    isPassword, err := h.passwordRexExp.MatchString(req.Password)
    if err != nil {
        ctx.String(http.StatusOK, "系统错误")
        return
    }
    if !isPassword {
        ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
        return
    }

    err = h.svc.Signup(ctx, domain.User{
        Email:    req.Email,
        Password: req.Password,
    })
    switch err {
    case nil:
        ctx.String(http.StatusOK, "注册成功")
    case service.ErrDuplicateEmail:
        ctx.String(http.StatusOK, "邮箱冲突，请换一个")
    default:
        ctx.String(http.StatusOK, "系统错误")
    }
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
    type LoginReq struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    var req LoginReq
    if err := ctx.Bind(&req); err != nil {
        return
    }

    u, err := h.svc.Login(ctx, req.Email, req.Password)
    switch err {
    case nil:
        h.setJWTToken(ctx, u.Id)
        ctx.String(http.StatusOK, "登录成功")
    case service.ErrInvalidUserOrPassword:
        ctx.String(http.StatusOK, "用户名或密码不对")
    default:
        ctx.String(http.StatusOK, "系统错误")
    }
}

func (h *UserHandler) Login(ctx *gin.Context) {
    type LoginReq struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    var req LoginReq
    if err := ctx.Bind(&req); err != nil {
        return
    }

    u, err := h.svc.Login(ctx, req.Email, req.Password)
    switch err {
    case nil:
        sess := sessions.Default(ctx)
        sess.Set("userId", u.Id)
        sess.Options(sessions.Options{
            MaxAge: 900,
        })
        err := sess.Save()
        if err != nil {
            ctx.String(http.StatusOK, "系统错误")
        }

        ctx.String(http.StatusOK, "登陆成功")
    case service.ErrInvalidUserOrPassword:
        ctx.String(http.StatusOK, "用户名或密码不对")
    default:
        ctx.String(http.StatusOK, "系统错误")
    }
}

func (h *UserHandler) Edit(ctx *gin.Context) {
    type EditReq struct {
        Nickname    string `json:"nickname"`
        Birthday    string `json:"birthday"`
        Description string `json:"description"`
    }
    var req EditReq
    if err := ctx.Bind(&req); err != nil {
        return
    }

    if nlen := utf8.RuneCountInString(req.Nickname); nlen > 30 || nlen <= 0 {
        ctx.String(http.StatusOK, "昵称长度要在1~30之间")
        return
    }

    _, err := time.Parse("2006-01-02", req.Birthday)
    if err != nil {
        ctx.String(http.StatusOK, "生日格式出错，需为 YYYY-MM-DD 格式")
        return
    }

    if dlen := utf8.RuneCountInString(req.Description); dlen > 500 {
        ctx.String(http.StatusOK, "个人简介长度不能大于500")
        return
    }

    us := ctx.MustGet("user").(UserClaims)
    uid := us.Uid
    err = h.svc.Edit(ctx, domain.User{
        Id:          uid,
        Nickname:    req.Nickname,
        Birthday:    req.Birthday,
        Description: req.Description,
    })
    switch err {
    case nil:
        ctx.String(http.StatusOK, "编辑成功")
    case service.ErrEditFailure:
        ctx.String(http.StatusOK, "编辑失败")
    default:
        ctx.String(http.StatusOK, "系统错误")
    }
}

func (h *UserHandler) Profile(ctx *gin.Context) {
    us := ctx.MustGet("user").(UserClaims)
    uid := us.Uid
    u, err := h.svc.Profile(ctx, uid)
    switch err {
    case nil:
        ctx.JSON(200, u)
    case service.ErrInvalidUser:
        ctx.String(http.StatusOK, "用户不存在")
    default:
        ctx.String(http.StatusOK, "系统错误")
    }
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
    type Req struct {
        Phone string `json:"phone"`
    }
    var req Req
    if err := ctx.Bind(&req); err != nil {
        return
    }
    // 你这边可以校验 Req
    if req.Phone == "" {
        ctx.JSON(http.StatusOK, Result{
            Code: 4,
            Msg:  "请输入手机号码",
        })
        return
    }
    err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
    switch err {
    case nil:
        ctx.JSON(http.StatusOK, Result{
            Msg: "发送成功",
        })
    case service.ErrCodeSendTooMany:
        ctx.JSON(http.StatusOK, Result{
            Code: 4,
            Msg:  "短信发送太频繁，请稍后再试",
        })
    default:
        ctx.JSON(http.StatusOK, Result{
            Code: 5,
            Msg:  "系统错误",
        })
        // 补充日志
    }
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
    type Req struct {
        Phone string `json:"phone"`
        Code  string `json:"code"`
    }
    var req Req
    if err := ctx.Bind(&req); err != nil {
        return
    }

    ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
    if err != nil {
        ctx.JSON(http.StatusOK, Result{
            Code: 5,
            Msg:  "系统异常",
        })
        return
    }
    if !ok {
        ctx.JSON(http.StatusOK, Result{
            Code: 4,
            Msg:  "验证码不对，请重新输入",
        })
        return
    }
    u, err := h.svc.FindOrCreate(ctx, req.Phone)
    if err != nil {
        ctx.JSON(http.StatusOK, Result{
            Code: 5,
            Msg:  "系统错误",
        })
        return
    }
    h.setJWTToken(ctx, u.Id)
    ctx.JSON(http.StatusOK, Result{
        Msg: "登录成功",
    })
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
    uc := UserClaims{
        Uid:       uid,
        UserAgent: ctx.GetHeader("User-Agent"),
        RegisteredClaims: jwt.RegisteredClaims{
            // 1 分钟过期
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
    tokenStr, err := token.SignedString(JWTKey)
    if err != nil {
        ctx.String(http.StatusOK, "系统错误")
    }
    ctx.Header("x-jwt-token", tokenStr)
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

type UserClaims struct {
    jwt.RegisteredClaims
    Uid       int64
    UserAgent string
}
