package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{dao: dao}
}

func (r *UserRepository) Create(c context.Context, u domain.User) error {
	return r.dao.Insert(c, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(c context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(c, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Update(c context.Context, u domain.User) error {
	return r.dao.Update(c, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		About:    u.AboutMe,
		Utime:    time.Now().UnixMilli(),
	})
}

func (r *UserRepository) FindById(c *gin.Context, id int64) (domain.User, error) {
	user, err := r.dao.FindById(c, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		AboutMe:  user.About,
	}, nil
}
