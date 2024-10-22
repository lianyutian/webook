package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
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
	user, err := r.cache.Get(c, id)
	// 缓存有数据
	if err == nil {
		return user, nil
	}
	// 缓存无数据
	u, err := r.dao.FindById(c, id)
	if err != nil {
		return domain.User{}, err
	}
	user = domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.About,
	}
	err = r.cache.Set(c, user.Id, user)
	if err != nil {
		log.Println(err)
	}
	return user, nil
}
