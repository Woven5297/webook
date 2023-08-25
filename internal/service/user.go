package service

import (
	"context"
	"errors"
	"gitee.com/webook/internal/domain"
	"gitee.com/webook/internal/repository"
	"gitee.com/webook/internal/repository/dao"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound
var ErrInvalidUserOrPassword = errors.New("账号或密码错误")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 要考虑加密放在哪里
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	// 存起来
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, u domain.User) (domain.User, error) {

	user, err := svc.repo.FindByEmail(ctx, u)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, err
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return domain.User{}, err
	}
	return user, err

}
