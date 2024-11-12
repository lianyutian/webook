package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"webook/internal/domain"
	"webook/internal/repository"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("invalid user/password")
	ErrSystemError           = errors.New("system error")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(c *gin.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return svc.repo.Create(c, user)
}

func (svc *UserService) Login(c *gin.Context, user domain.User) (domain.User, error) {
	u, err := svc.repo.FindByEmail(c, user.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) Edit(c *gin.Context, user domain.User) error {
	err := svc.repo.Update(c, user)
	if err != nil {
		return err
	}
	return nil
}

func (svc *UserService) Profile(c *gin.Context) (domain.User, error) {
	userId, ok := c.Get("userId")
	if !ok {
		c.String(http.StatusOK, "系统错误")
		return domain.User{}, ErrSystemError
	}

	user, err := svc.repo.FindById(c, userId.(int64))

	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (svc *UserService) FindOrCreate(c *gin.Context, phone string) (domain.User, error) {
	user, err := svc.repo.FindByPhone(c, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		// 绝大部分请求都会从这返回
		return user, err
	}
	// 注册用户
	err = svc.repo.Create(c, domain.User{
		Phone: phone,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return user, ErrUserDuplicate
		}
		return user, err
	}

	// TODO 主从模式这里要在主库读取
	return svc.repo.FindByPhone(c, user.Phone)
}
