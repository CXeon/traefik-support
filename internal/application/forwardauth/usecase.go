package forwardauth

import (
	"context"

	forwardauthdomain "github.com/CXeon/traefik-support/internal/domain/forwardauth"
)

// CheckRequest 是 Application 层的输入 DTO，与 Domain 层解耦。
type CheckRequest struct {
	Env     string
	Cluster string
	Company string
	Project string
	Color   string
	UserID  string
	JWT     string
	Path    string
}

// UseCase 处理 ForwardAuth 认证逻辑。
type UseCase interface {
	Check(ctx context.Context, req *CheckRequest) error
}

type useCase struct {
	domainSvc forwardauthdomain.AuthDomainService
}

func NewUseCase(domainSvc forwardauthdomain.AuthDomainService) UseCase {
	return &useCase{domainSvc: domainSvc}
}

func (uc *useCase) Check(ctx context.Context, req *CheckRequest) error {
	// TODO: OpenTelemetry span 注入
	accessCtx := &forwardauthdomain.AccessContext{
		Env:     req.Env,
		Cluster: req.Cluster,
		Company: req.Company,
		Project: req.Project,
		Color:   req.Color,
		UserID:  req.UserID,
		JWT:     req.JWT,
		Path:    req.Path,
	}
	return uc.domainSvc.Check(ctx, accessCtx)
}
