package forwardauth

import "context"

// UserClient 调用 user 服务验证用户身份，由基础设施层实现。
type UserClient interface {
	// Verify 验证用户身份。
	// user 服务负责：解析 JWT，比对 token 中 UserID 与 header UserID 是否一致。
	// 通过返回 nil，失败返回 tileerrors 定义的错误类型。
	Verify(ctx context.Context, userID, jwt string) error
}

// AuthDomainService 验证访问请求是否被允许。
type AuthDomainService interface {
	// Check 返回 nil 表示放行，返回 error 表示拒绝。
	Check(ctx context.Context, accessCtx *AccessContext) error
}

type authDomainService struct {
	userClient UserClient
}

func NewAuthDomainService(userClient UserClient) AuthDomainService {
	return &authDomainService{userClient: userClient}
}

func (s *authDomainService) Check(ctx context.Context, accessCtx *AccessContext) error {
	return s.userClient.Verify(ctx, accessCtx.UserID, accessCtx.JWT)
}
