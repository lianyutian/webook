package web

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"regexp"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
)

const emailRegexPattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
const passwordRegexPattern = "^[a-zA-Z0-9!@#\\$%\\^&\\*\\(\\)_\\+\\-=\\[\\]\\{\\};':\",.<>?/\\\\|]{6,10}$"
const nickNameRegexPattern = `^[a-zA-Z0-9_\p{Han}]{3,16}$`
const birthDayRegexPattern = `^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`
const aboutRegexPattern = `^[a-zA-Z0-9_\p{Han}]{3,16}$`
const phoneRegexPattern = `^1[3-9]\d{9}$`
const biz = "login"

type UserHandler struct {
	userSvc *service.UserService
	codeSvc *service.CodeService
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 自己要放入 token 的数据
	Uid int64
}

var JwtSecret = []byte("your-secret-key")

func NewUserHandler(userSvc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
		codeSvc: codeSvc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/users")
	group.POST("signup", u.signUp)
	group.POST("login", u.login)
	group.POST("login_sms/code/send", u.sendSmsLoginCode)
	group.POST("login_sms", u.loginSms)
	group.GET("profile", u.profile)
	group.POST("edit", u.edit)
}

// 校验电子邮件格式
func isValidEmail(email string) bool {
	re := regexp.MustCompile(emailRegexPattern)
	return re.MatchString(email)
}

// 校验电子邮件格式
func isValidPassword(password string) bool {
	re := regexp.MustCompile(passwordRegexPattern)
	return re.MatchString(password)
}

// 校验手机号码
func isValidPhone(phone string) bool {
	re := regexp.MustCompile(phoneRegexPattern)
	return re.MatchString(phone)
}

func (u *UserHandler) signUp(c *gin.Context) {
	type signUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req signUpReq

	if err := c.Bind(&req); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if !isValidEmail(req.Email) {
		c.String(http.StatusOK, "Email is invalid")
		return
	}
	if !isValidPassword(req.Password) {
		c.String(http.StatusOK, "Password is invalid")
		return
	}

	err := u.userSvc.SignUp(c, domain.User{Email: req.Email, Password: req.Password})
	if errors.Is(err, service.ErrUserDuplicate) {
		c.String(http.StatusOK, "Email is duplicated")
		return
	}

	c.String(http.StatusOK, "Sign up success")
}

func (u *UserHandler) login(c *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq

	if err := c.Bind(&req); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	user, err := u.userSvc.Login(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		c.String(http.StatusOK, "Invalid user or password")
		return
	}

	if err != nil {
		c.String(http.StatusInternalServerError, "系统找错误")
		return
	}

	err = u.setJWT(c, user)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
	}

	c.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) setJWT(c *gin.Context, user domain.User) error {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid: user.Id,
	})
	tokenString, err := claims.SignedString(JwtSecret)
	c.Header("x-jwt-token", tokenString)
	return err
}

func (u *UserHandler) sendSmsLoginCode(c *gin.Context) {
	type loginReq struct {
		Phone string `json:"phone"`
	}

	var req loginReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "入参错误",
		})
	}

	ok := isValidPhone(req.Phone)
	if !ok {
		c.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "无效的手机号",
		})
		return
	}

	err := u.codeSvc.Send(c, biz, req.Phone)
	switch {
	case err == nil:
		c.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "发送短信成功",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		c.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "发送短信太频繁",
		})
	default:
		c.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "系统错误",
		})
		return
	}
}

func (u *UserHandler) loginSms(c *gin.Context) {
	type loginReq struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	var req loginReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Code: -1,
			Msg:  "入参错误",
		})
	}
	ok := isValidPhone(req.Phone)
	if !ok {
		c.String(http.StatusOK, "phone number is invalid")
		return
	}
	// 验证验证码
	ok, err := u.codeSvc.Verify(c, biz, req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "验证码错误",
		})
		return
	}
	// 查询用户是否存在
	// 不存在注册用户
	user, err := u.userSvc.FindOrCreate(c, req.Phone)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "系统错误",
		})
		return
	}

	err = u.setJWT(c, user)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "系统错误",
		})
	}

	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "OK",
		Data: user,
	})
}

func (u *UserHandler) profile(c *gin.Context) {
	profile, err := u.userSvc.Profile(c)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		c.String(http.StatusUnauthorized, "请登录")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误")
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (u *UserHandler) edit(c *gin.Context) {
	type editReq struct {
		NickName string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	var req editReq

	if err := c.Bind(&req); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if !regexp.MustCompile(nickNameRegexPattern).MatchString(req.NickName) {
		c.String(http.StatusOK, "nickname is invalid")
		return
	}
	if !regexp.MustCompile(birthDayRegexPattern).MatchString(req.Birthday) {
		c.String(http.StatusOK, "birthday is invalid")
		return
	}
	if !regexp.MustCompile(aboutRegexPattern).MatchString(req.AboutMe) {
		c.String(http.StatusOK, "about is invalid")
		return
	}

	userId := sessions.Default(c).Get("userId")
	id := userId.(int64)
	err := u.userSvc.Edit(c, domain.User{
		Id:       id,
		Nickname: req.NickName,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "更新成功")
}
