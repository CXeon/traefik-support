package forwardauth

// AccessContext 描述一次需要被验证的访问请求的上下文。
type AccessContext struct {
	Env     string // 环境
	Cluster string // 集群
	Company string // 公司
	Project string // 项目
	Color   string // 染色标识（用于流量分流，是路由 key 的组成元素）
	UserID  string // 用户唯一标识（从 X-UserID header 获取）
	JWT     string // 原始 JWT token（从 Authorization header 提取）
	Path    string // 请求路径（从 X-Forwarded-Uri 获取）
}
