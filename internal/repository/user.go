package repository

import (
	"context"
	"gitee.com/webook/internal/domain"
	"gitee.com/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindById() {
	// 先从 cache 里找
	// 再从 dao 里面找
	// 找到了回写 cache
}
