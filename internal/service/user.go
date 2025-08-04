package service

import (
    "context"
    "errors"
    "github.com/Guanjian104/webook/internal/domain"
    "github.com/Guanjian104/webook/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

var (
    ErrDuplicateEmail        = repository.ErrDuplicateUser
    ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
    ErrEditFailure           = repository.ErrEditFailure
    ErrInvalidUser           = errors.New("用户不存在")
)

type UserService interface {
    Signup(ctx context.Context, u domain.User) error
    Login(ctx context.Context, email string, password string) (domain.User, error)
    Edit(ctx context.Context, u domain.User) error
    Profile(ctx context.Context, Id int64) (domain.UserProfile, error)
    FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{
        repo: repo,
    }
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hash)
    return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
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

func (svc *userService) Edit(ctx context.Context, u domain.User) error {
    return svc.repo.Edit(ctx, u)
}

func (svc *userService) Profile(ctx context.Context, Id int64) (domain.UserProfile, error) {
    u, err := svc.repo.FindById(ctx, Id)
    if errors.Is(err, repository.ErrUserNotFound) {
        return domain.UserProfile{}, ErrInvalidUser
    }
    if err != nil {
        return domain.UserProfile{}, err
    }
    return u, nil
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
    u, err := svc.repo.FindByPhone(ctx, phone)
    if !errors.Is(err, repository.ErrUserNotFound) {
        return u, err
    }
    err = svc.repo.Create(ctx, domain.User{
        Phone: phone,
    })

    if err != nil && !errors.Is(err, repository.ErrDuplicateUser) {
        return domain.User{}, err
    }
    // TODO: 主从延迟，强制走主库
    return svc.repo.FindByPhone(ctx, phone)
}
