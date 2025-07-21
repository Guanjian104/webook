package repository

import (
    "context"
    "github.com/Guanjian104/webook/internal/domain"
    "github.com/Guanjian104/webook/internal/repository/cache"
    "github.com/Guanjian104/webook/internal/repository/dao"
)

var (
    ErrDuplicateEmail = dao.ErrDuplicateEmail
    ErrUserNotFound   = dao.ErrRecordNotFound
    ErrEditFailure    = dao.ErrEditFailure
)

type UserRepository struct {
    dao   *dao.UserDAO
    cache *cache.UserCache
}

func NewUserRepository(d *dao.UserDAO, c *cache.UserCache) *UserRepository {
    return &UserRepository{
        dao:   d,
        cache: c,
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

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
    u, err := repo.dao.FindByEmail(ctx, email)
    if err != nil {
        return domain.User{}, err
    }
    return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx context.Context, Id int64) (domain.UserProfile, error) {
    u, err := repo.cache.Get(ctx, Id)
    if err == nil {
        return domainToProfile(u), err
    }
    ue, err := repo.dao.FindById(ctx, Id)
    if err != nil {
        return domain.UserProfile{}, err
    }
    up := domain.UserProfile{
        Nickname:    ue.Nickname,
        Birthday:    ue.Birthday,
        Description: ue.Description,
    }

    // 设置缓存
    ud := domain.User{
        Id:          ue.Id,
        Email:       ue.Email,
        Password:    ue.Password,
        Nickname:    ue.Nickname,
        Birthday:    ue.Birthday,
        Description: ue.Description,
    }
    _ = repo.cache.Set(ctx, ud)

    return up, nil
}

func domainToProfile(u domain.User) domain.UserProfile {
    // Domain To Profile
    return domain.UserProfile{
        Nickname:    u.Nickname,
        Birthday:    u.Birthday,
        Description: u.Description,
    }
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
        Nickname:    u.Nickname,
        Birthday:    u.Birthday,
        Description: u.Description,
    }
}
