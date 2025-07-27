package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/Guanjian104/webook/internal/repository"
    "github.com/Guanjian104/webook/internal/service/sms"
    "math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany

type CodeService interface {
    Send(ctx context.Context, biz, phone string) error
    Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type codeService struct {
    repo repository.CodeRepository
    sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
    return &codeService{
        repo: repo,
        sms:  smsSvc,
    }
}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
    code := svc.generate()
    err := svc.repo.Set(ctx, biz, phone, code)
    if err != nil {
        return err
    }
    const codeTplId = "1877556"
    return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
    ok, err := svc.Verify(ctx, biz, phone, inputCode)
    if errors.Is(err, repository.ErrCodeVerifyTooMany) {
        return false, nil
    }
    return ok, err
}

func (svc *codeService) generate() string {
    // 0-999999
    code := rand.Intn(1000000)
    return fmt.Sprintf("%06d", code)
}
