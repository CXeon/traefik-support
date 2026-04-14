package forwardauth

import (
	"context"

	forwardauthdomain "github.com/CXeon/traefik-support/internal/domain/forwardauth"
)

// stubUserClient 是 UserClient 的占位实现，待 RPC 客户端完成后替换。
// 当前行为：放行所有请求（仅用于开发阶段）。
type stubUserClient struct{}

// NewStubUserClient 返回 UserClient 的 stub 实现。
// TODO: 替换为真实的 RPC 调用 user 服务实现。
func NewStubUserClient() forwardauthdomain.UserClient {
	return &stubUserClient{}
}

// Verify TODO: RPC 调用 user 服务，验证 JWT 与 UserID 是否匹配。
func (c *stubUserClient) Verify(_ context.Context, _, _ string) error {
	return nil
}
