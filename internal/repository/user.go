package repository

import (
	"context"
	"github.com/Guanjian104/webook/internal/domain"
	"github.com/Guanjian104/webook/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
	ErrEditFailure    = dao.ErrEditFailure
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) Edit(ctx context.Context, u domain.User) error {
	return repo.dao.Update(ctx, dao.User{
		Id:          u.Id,
		Nickname:    u.Nickname,
		Birthday:    u.Birthday,
		Description: u.Description,
	})
}

func (repo *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx *gin.Context, Id int64) (domain.UserProfile, error) {
	u, err := repo.dao.FindById(ctx, Id)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return repo.toProfile(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (repo *UserRepository) toProfile(u dao.User) domain.UserProfile {
	return domain.UserProfile{
		Id:          u.Id,
		Email:       u.Email,
		Nickname:    u.Nickname,
		Birthday:    u.Birthday,
		Description: u.Description,
	}
}
