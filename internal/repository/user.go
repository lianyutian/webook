package repository

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
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

func (rep *UserRepository) Create(c context.Context, user domain.User) error {
	return rep.dao.Insert(c, rep.domainToEntity(user))
}

func (rep *UserRepository) FindByEmail(c context.Context, email string) (domain.User, error) {
	ud, err := rep.dao.FindByEmail(c, email)
	if err != nil {
		return domain.User{}, err
	}
	return rep.entityToDomain(ud), nil
}

func (rep *UserRepository) FindByPhone(c context.Context, phone string) (domain.User, error) {
	ud, err := rep.dao.FindByPhone(c, phone)
	if err != nil {
		return domain.User{}, err
	}
	return rep.entityToDomain(ud), nil
}

func (rep *UserRepository) Update(c context.Context, u domain.User) error {
	return rep.dao.Update(c, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		About:    u.AboutMe,
		Utime:    time.Now().UnixMilli(),
	})
}

func (rep *UserRepository) FindById(c *gin.Context, id int64) (domain.User, error) {
	user, err := rep.cache.Get(c, id)
	// 缓存有数据
	if err == nil {
		return user, nil
	}
	// 缓存无数据
	ud, err := rep.dao.FindById(c, id)
	if err != nil {
		return domain.User{}, err
	}
	user = rep.entityToDomain(ud)
	err = rep.cache.Set(c, user.Id, user)
	if err != nil {
		log.Println(err)
	}
	return user, nil
}

func (rep *UserRepository) entityToDomain(ud dao.User) domain.User {
	return domain.User{
		Id:       ud.Id,
		Email:    ud.Email.String,
		Phone:    ud.Phone.String,
		Nickname: ud.Nickname,
		Birthday: ud.Birthday,
		AboutMe:  ud.About,
	}
}

func (rep *UserRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		About:    user.AboutMe,
	}
}
