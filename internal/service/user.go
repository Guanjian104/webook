package service

import (
	"context"
	"errors"
	"github.com/Guanjian104/webook/internal/domain"
	"github.com/Guanjian104/webook/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
	ErrEditFailure           = repository.ErrEditFailure
	ErrInvalidUser           = errors.New("用户不存在")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx *gin.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 检查密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) Edit(ctx context.Context, u domain.User) error {
	return svc.repo.Edit(ctx, u)
}

func (svc *UserService) Profile(ctx *gin.Context, Id int64) (domain.UserProfile, error) {
	u, err := svc.repo.FindById(ctx, Id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.UserProfile{}, ErrInvalidUser
	}
	if err != nil {
		return domain.UserProfile{}, err
	}
	return u, nil
}
