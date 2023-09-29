package web

import (
	"errors"
	"fmt"
	"gitee.com/webook/internal/domain"
	"gitee.com/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会直接写会一个 404 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//const (
	//	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	//	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	//)
	//emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		fmt.Printf("email error: (%s)", err)
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不对")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	//passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		fmt.Printf("passwrod error: (%s)", err)
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码格式不对")
		return
	}
	fmt.Printf("req: (%v)", req)
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	//if err != nil {
	//	ctx.String(http.StatusOK, "service 方法出错")
	//}
	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "账号或密码错误")
		return
	}
	if errors.Is(err, service.ErrUserNotFound) {
		ctx.String(http.StatusOK, "找不到相关账号")
		return
	}
	println("user", user)
	// 用 jwt 设置登录状态
	// 生成一个 jwt token
	claims := UserClaims{

		Uid: user.Id,
	}
	//token := jwt.New(jwt.SigningMethodHS512)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "token 生成错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)

	ctx.String(http.StatusOK, "登录成功")
	return

}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "账号或密码错误")
		return
	}
	if errors.Is(err, service.ErrUserNotFound) {
		ctx.String(http.StatusOK, "找不到相关账号")
		return
	}
	// 登录成功 设置 session
	s := sessions.Default(ctx)
	// 设置值
	s.Set("userId", user.Id)
	s.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		MaxAge:   10,
	})
	s.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	s := sessions.Default(ctx)
	s.Options(sessions.Options{
		MaxAge: -1,
	})
	s.Save()
	ctx.String(http.StatusOK, "退出登录成功")
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	s := sessions.Default(ctx)

	id := s.Get("userId")
	if id == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := u.svc.Profile(ctx, id.(int64))
	if err != nil {
		ctx.String(http.StatusOK, "找不到相关账号")
	}
	ctx.JSON(http.StatusOK, domain.User{
		Email:       user.Email,
		NickName:    user.NickName,
		Age:         user.Age,
		Description: user.Description,
	})
}

func (u *UserHandler) RegisterUserRoutes(server *gin.Engine) {

	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}
