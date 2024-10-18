package web

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"webook/internal/domain"
	"webook/internal/service"
)

const emailRegexPattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
const passwordRegexPattern = "^[a-zA-Z0-9!@#\\$%\\^&\\*\\(\\)_\\+\\-=\\[\\]\\{\\};':\",.<>?/\\\\|]{6,10}$"
const nickNameRegexPattern = `^[a-zA-Z0-9_\p{Han}]{3,16}$`
const birthDayRegexPattern = `^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`
const aboutRegexPattern = `^[a-zA-Z0-9_\p{Han}]{3,16}$`

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/users")
	group.POST("/signup", u.signUp)
	group.POST("/login", u.login)
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

	err := u.svc.SignUp(c, domain.User{Email: req.Email, Password: req.Password})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
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

	user, err := u.svc.Login(c, domain.User{
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

	session := sessions.Default(c)
	session.Set("userId", user.Id)
	session.Options(sessions.Options{
		MaxAge: 10,
	})
	err = session.Save()
	if err != nil {
		c.String(http.StatusInternalServerError, "系统找错误")
		return
	}

	c.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) profile(c *gin.Context) {
	profile, err := u.svc.Profile(c)
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
	err := u.svc.Edit(c, domain.User{
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
