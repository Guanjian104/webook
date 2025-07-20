package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
	ErrEditFailure    = errors.New("编辑失败")
)

type UserDAO struct {
	db *gorm.DB
}

type User struct {
	Id          int64  `gorm:"primaryKey,autoIncrement"`
	Email       string `gorm:"unique"`
	Password    string
	Nickname    string
	Birthday    string
	Description string

	Ctime int64
	Utime int64
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) Update(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	result := dao.db.WithContext(ctx).Model(&u).Updates(User{Nickname: u.Nickname, Birthday: u.Birthday, Description: u.Description})
	if result.RowsAffected == 0 {
		return ErrEditFailure
	}
	return result.Error
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, Id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, Id).Error
	return u, err
}
