package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("该邮箱已注册")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

type User struct {
	Id       int64  `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"type:varchar(100);unique"`
	Password string `gorm:"type:varchar(100)"`
	Nickname string `gorm:"type:varchar(100)"`
	Birthday string `gorm:"type:varchar(10)"`
	About    string `gorm:"type:varchar(100)"`

	Ctime int64
	Utime int64
}

func (dao *UserDAO) Insert(c context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := dao.db.WithContext(c).Create(&u).Error

	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const uniqueIndexErrNo uint16 = 1062
		if me.Number == uniqueIndexErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(c context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(c).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) Update(ctx context.Context, user User) error {
	return dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", user.Id).Updates(&user).Error
}

func (dao *UserDAO) FindById(c *gin.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(c).Where("id = ?", id).First(&u).Error
	return u, err
}
