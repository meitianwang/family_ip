# 家庭 IP 代理售卖系统 — 设计规格

**日期**: 2026-04-14  
**基础项目**: fast-frame  
**方案**: 混合方案 — 新增代理专属模块，复用现有用户/支付/Admin 框架

---

## 一、业务背景

从 lisahost.com 购买带家庭 IP 的 VPS，在 VPS 上手动配置 HTTP 代理（3proxy 等）和 VLESS（Xray）。通过本系统向用户销售这些 IP 的访问权限。

**计费模型**: 用户选购具体 IP + 套餐（天数 + 流量上限），独享使用。  
**协议支持**: HTTP 代理（用户名/密码认证）+ VLESS（兼容 Shadowrocket）。  
**开通方式**: 半自动 — 系统自动生成凭证，管理员手动在 VPS 上同步账号。

---

## 二、数据模型

### 新增 5 个 Ent Schema

#### ProxyNode — 代理节点

```
id              uuid, PK
ip_address      string, unique, not null
country         string
country_code    string        // 2位ISO码，用于显示国旗
city            string
isp             string
http_port       int
vless_port      int
vless_network   enum(tcp, ws, grpc)
vless_tls       bool
vless_sni       string        // TLS SNI 域名，空=直接IP
vless_ws_path   string        // WebSocket路径，仅network=ws时有效
tags            []string      // 如 ["Netflix解锁", "低延迟", "游戏"]
status          enum(available, rented, offline)
description     text          // 管理员备注
created_at      timestamptz
updated_at      timestamptz
deleted_at      timestamptz   // 软删除
```

索引: `(status)`, `(country_code)`, `(deleted_at)` where null

#### ProxyProduct — 租用套餐

```
id                uuid, PK
name              string
description       text
duration_days     int           // 有效天数：1 / 7 / 30
traffic_limit_gb  int           // 流量上限 GB，0 表示不限
price             decimal(10,2) // 人民币
sort_order        int           // 展示排序
is_active         bool
created_at        timestamptz
updated_at        timestamptz
```

#### ProxyRental — 租用订单

```
id                   uuid, PK
user_id              FK → User
node_id              FK → ProxyNode
product_id           FK → ProxyProduct
payment_order_id     FK → PaymentOrder (nullable，取消订单时为空)
status               enum(pending_payment, active, expired, cancelled)
started_at           timestamptz   // 激活时设置
expires_at           timestamptz   // started_at + duration_days
traffic_used_bytes   int64         // 管理员累计更新
traffic_limit_bytes  int64         // 创建时从 product 快照
created_at           timestamptz
updated_at           timestamptz
```

索引: `(user_id, status)`, `(node_id, status)`, `(expires_at)` where status=active

并发控制: 创建租用时，在数据库事务内 `SELECT FOR UPDATE` ProxyNode，确认 status=available 后再创建，防止同一 IP 被并发购买。

#### ProxyCredential — 访问凭证

```
id             uuid, PK
rental_id      FK → ProxyRental, unique
http_username  string    // 随机 8 位字母数字
http_password  string    // 随机 16 位
vless_uuid     uuid v4   // 随机生成
vless_link     text      // 预计算的完整 vless:// 链接
created_at     timestamptz
```

vless:// 链接格式:
```
vless://{uuid}@{ip}:{port}?encryption=none&security={none|tls}&sni={sni}&type={network}&path={ws_path}&host={sni}#{备注名}
```

示例(ws+tls):
```
vless://550e8400-e29b-41d4-a716-446655440000@1.2.3.4:443?encryption=none&security=tls&sni=example.com&type=ws&path=/vless&host=example.com#HK-家庭IP
```

#### ProxyTrafficLog — 流量更新记录

```
id            uuid, PK
rental_id     FK → ProxyRental
delta_bytes   int64    // 本次新增流量（非累计）
operator_id   FK → User (管理员)
note          text
created_at    timestamptz
```

### 复用现有实体（零改动）

- `User` — 用户账户
- `PaymentOrder` — 支付订单，通过 `payment_order_id` 关联到 ProxyRental
- `PaymentChannel` — 支付渠道（Stripe / 支付宝 / 微信）

---

## 三、支付流程

```
用户选 IP + 套餐
      ↓
POST /api/proxy/rentals
      ↓ 事务内:
  SELECT FOR UPDATE ProxyNode（确认 available）
  创建 ProxyRental(status=pending_payment)
  创建 PaymentOrder（复用现有支付模块）
      ↓
前端跳转支付（现有支付页面，零改动）
      ↓
支付 Webhook 回调（现有 webhook handler）
      ↓ 新增钩子:
  查找关联的 ProxyRental
  事务内激活:
    ProxyRental.status = active
    ProxyRental.started_at = now()
    ProxyRental.expires_at = now() + duration_days
    ProxyNode.status = rented
    创建 ProxyCredential（生成凭证）
      ↓
用户在"我的租用"页查看凭证
```

支付取消/超时: PaymentOrder 超时时，将 ProxyRental 设为 cancelled，ProxyNode 状态回滚为 available。

---

## 四、过期处理

后端启动一个定时 goroutine（每小时执行一次）:

1. 查询 `status=active AND expires_at < now()` 的 ProxyRental
2. 批量设置 `status=expired`
3. 将对应 ProxyNode 状态改回 `available`

流量耗尽处理: 管理员更新流量时，若 `traffic_used_bytes >= traffic_limit_bytes`（且 limit > 0），同步将 rental 设为 expired，node 改为 available。

---

## 五、后端 API

### 用户端

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/proxy/nodes` | 节点列表（支持 country/tags 筛选，仅返回 available） |
| GET | `/api/proxy/nodes/:id` | 节点详情 |
| GET | `/api/proxy/products` | 套餐列表（仅 active） |
| POST | `/api/proxy/rentals` | 创建租用 + 支付订单 |
| GET | `/api/proxy/rentals` | 我的租用列表 |
| GET | `/api/proxy/rentals/:id` | 租用详情（含凭证，仅 active 时返回） |
| POST | `/api/proxy/rentals/:id/cancel` | 取消待支付订单 |

### 管理员端

| 方法 | 路径 | 说明 |
|------|------|------|
| GET/POST | `/api/admin/proxy/nodes` | 节点列表 / 新增节点 |
| PUT/DELETE | `/api/admin/proxy/nodes/:id` | 更新 / 删除节点 |
| GET/POST | `/api/admin/proxy/products` | 套餐列表 / 新增套餐 |
| PUT/DELETE | `/api/admin/proxy/products/:id` | 更新 / 删除套餐 |
| GET | `/api/admin/proxy/rentals` | 所有租用（支持状态/用户/节点筛选） |
| GET | `/api/admin/proxy/rentals/:id` | 租用详情 |
| POST | `/api/admin/proxy/rentals/:id/traffic` | 更新流量使用量 |
| POST | `/api/admin/proxy/rentals/:id/expire` | 强制过期（提前终止） |

---

## 六、后端代码结构

遵循项目现有分层架构:

```
backend/
├── ent/schema/
│   ├── proxy_node.go
│   ├── proxy_product.go
│   ├── proxy_rental.go
│   ├── proxy_credential.go
│   └── proxy_traffic_log.go
├── internal/
│   ├── handler/
│   │   ├── proxy/              # 用户端
│   │   │   ├── node.go         # 节点列表/详情
│   │   │   └── rental.go       # 租用创建/查看/取消
│   │   └── admin/
│   │       ├── proxy_node.go
│   │       ├── proxy_product.go
│   │       └── proxy_rental.go
│   ├── service/
│   │   └── proxy/
│   │       ├── node_service.go
│   │       ├── rental_service.go     # 核心：创建/激活/过期
│   │       ├── credential_service.go # 凭证生成逻辑
│   │       └── traffic_service.go    # 流量更新
│   ├── repository/
│   │   ├── proxy_node_repo.go
│   │   ├── proxy_rental_repo.go
│   │   ├── proxy_credential_repo.go
│   │   └── proxy_traffic_log_repo.go
│   └── server/
│       └── routes/               # 注册新路由（修改现有文件）
└── migrations/
    └── 0031_proxy_tables.sql
```

支付 webhook 集成点: 在现有 `service/payment/` 的支付成功回调中，新增调用 `proxy/rental_service.go` 的 `ActivateByPaymentOrder(orderID)` 方法。

---

## 七、前端页面

### 用户端（新增路由）

| 路由 | 页面 | 说明 |
|------|------|------|
| `/proxy` | `MarketplacePage.vue` | IP 市场：节点卡片列表，按国家/标签筛选，选择套餐后购买 |
| `/proxy/rentals` | `MyRentalsPage.vue` | 我的租用：列表展示状态、到期时间、流量进度 |
| `/proxy/rentals/:id` | `RentalDetailPage.vue` | 租用详情：HTTP 凭证、VLESS 链接（一键复制）、流量图表 |

### 管理员端（新增路由）

| 路由 | 页面 | 说明 |
|------|------|------|
| `/admin/proxy/nodes` | `ProxyNodesPage.vue` | 节点管理：CRUD，含 VLESS 配置参数 |
| `/admin/proxy/products` | `ProxyProductsPage.vue` | 套餐管理：CRUD，上下架 |
| `/admin/proxy/rentals` | `ProxyRentalsPage.vue` | 租用管理：查看全部，更新流量，强制过期 |

---

## 八、凭证生成规则

**HTTP 账号**: 8 位随机字母数字（大小写+数字，排除易混淆字符 0/O/l/1）  
**HTTP 密码**: 16 位随机字符串（字母+数字+特殊字符）  
**VLESS UUID**: 标准 UUID v4 随机生成  
**凭证备注名（VLESS 链接 `#` 后缀）**: `{city}-家庭IP`，如 `HongKong-家庭IP`

---

## 九、不在范围内（本期不做）

- VPS 上代理软件的自动化配置（SSH/API 调用）
- 流量自动采集（管理员手动更新）
- IP 续费（到期重新购买即可）
- 多用户共享同一 IP

---

## 十、实现顺序建议

1. 数据库迁移（5 个 Schema + SQL migration）
2. Repository 层
3. Service 层（rental_service 最核心）
4. 支付 webhook 集成
5. Admin API + 前端管理页（先让管理员能录入节点/套餐）
6. 用户端 API + 前端页面
7. 定时过期 goroutine
