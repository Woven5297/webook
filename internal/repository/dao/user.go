package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.CreatedTime = now
	u.UpdatedTime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

// User 直接对应数据库表 相当于 entity
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 创建时间 毫秒数
	CreatedTime int64
	// 更新时间 毫秒数
	UpdatedTime int64
}