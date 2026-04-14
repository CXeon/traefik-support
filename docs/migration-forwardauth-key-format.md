# ForwardAuth Middleware Key 格式迁移说明

## 背景

`tiles/gateway/traefik/handler.go` 中 ForwardAuth 中间件名称的参数顺序存在 bug。
旧格式（错误）：`{Company}.{Project}.{Env}.{Cluster}.ForwardAuth`
新格式（正确）：`{Env}.{Cluster}.{Company}.{Project}.ForwardAuth`

## 影响

- 已写入 KV Store 的旧格式中间件 key 将与路由引用产生不匹配
- 路由在 Traefik 中将处于 disabled 状态，直到迁移完成

## 迁移步骤

### 1. 确认旧格式 key（etcd 示例）

```bash
# 列出所有 ForwardAuth 中间件 key���旧格式以 Company 开头）
etcdctl get --prefix traefik/http/middlewares/ --keys-only

# 找出旧格式 key（格式：traefik/http/middlewares/{Company}.{Project}.{Env}.{Cluster}.ForwardAuth/）
```

### 2. 删除旧格式 key

```bash
# 对每个旧格式中间件，删除其所有子 key
# 示例：Company=acme, Project=billing, Env=prod, Cluster=china
etcdctl del --prefix "traefik/http/middlewares/acme.billing.prod.china.ForwardAuth/"
```

### 3. 部署新版 traefik-support

新版 `traefik-support` 启动时会自动调用 `provisionMiddleware()`，以新格式
（`{Env}.{Cluster}.{Company}.{Project}.ForwardAuth`）写入 ForwardAuth 中间件定义。

### 4. 验证

```bash
# 确认新格式 key 已写入
etcdctl get --prefix "traefik/http/middlewares/prod.china.acme.billing.ForwardAuth/"

# 预期输出：
# traefik/http/middlewares/prod.china.acme.billing.ForwardAuth/forwardAuth/address
# traefik/http/middlewares/prod.china.acme.billing.ForwardAuth/forwardAuth/trustForwardHeader
```

### 5. 回滚方案

如新版出现问题：
1. 回滚 `traefik-support` 至旧版本
2. 手动写入旧格式中间件 key 恢复路由：

```bash
OLD_NAME="acme.billing.prod.china.ForwardAuth"
ADDRESS="http://<traefik-support-ip>:<port>"
etcdctl put "traefik/http/middlewares/${OLD_NAME}/forwardAuth/address" "${ADDRESS}"
etcdctl put "traefik/http/middlewares/${OLD_NAME}/forwardAuth/trustForwardHeader" "true"
```

## Consul 环境（如适用）

```bash
# 删除旧格式 key
consul kv delete -recurse "traefik/http/middlewares/acme.billing.prod.china.ForwardAuth/"

# 验证新格式 key
consul kv get -recurse "traefik/http/middlewares/prod.china.acme.billing.ForwardAuth/"
```
