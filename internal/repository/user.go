package repository

import (
	"context"
	"gitee.com/webook/internal/domain"
	"gitee.com/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

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

func (r *UserRepository) FindByEmail(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := r.dao.FindByEmail(ctx, u.Email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
