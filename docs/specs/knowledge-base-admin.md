# 知识库后台管理模块

## 状态
- 创建日期: 2026-09-15
- 状态: 草稿

## 目标
在管理后台（`web/admin` + `cmd/admin` :8083）新增「知识库」一级模块，让平台管理员能够对私有知识库做**内容审核与合规处置**、**资源与配额治理**、**运营数据概览**，补齐当前知识库完全没有后台入口的空白。

## 非目标
- **不开放私有文档正文给管理员**。管理员只能看元数据 + 已公开内容（`is_public=true` 且 `status=1`），私有正文一律不可读。
- 不做举报工单系统（用户侧举报入口、工单流转）。本期管理员只能主动巡检处置，举报渠道后续迭代。
- 不做按用户细粒度配额（每人单独设额度）。本期只做**全局配额上限**，复用现有 `settings` 键值表。
- 不做管理员代编辑文档内容、代恢复误删文档。
- 不改动知识库现有的成员角色体系（`owner/admin/editor/viewer`），平台管理员与工作区角色是两套东西，不混用。
- 不引入新 UI 库、不引入 i18n（沿用 Element Plus + axios + 中文单语）。

## 用户故事
- 作为平台管理员，我想查看全站工作区与文档列表并按用户、可见性、状态筛选，以便巡检违规内容。
- 作为平台管理员，我想把违规的工作区或文档**下架**并填写理由，以便立刻切断其对外传播，同时让作者知道被处置的原因。
- 作为平台管理员，我想查看并强制删除任意免登录分享链接，以便在内容外泄时快速止血。
- 作为平台管理员，我想删除严重违规的工作区或文档，以便彻底清理平台内容。
- 作为平台管理员，我想查看每个用户占用的工作区数、文档数、存储空间，并设置全局上限，以便防止个别用户滥用存储。
- 作为平台管理员，我想看到知识库整体运营数据（工作区总数、文档总数、公开 wiki 流量），以便评估功能健康度。
- 作为平台管理员，我想查看历史处置记录，以便复核处置是否得当、追溯谁在什么时候做了什么。

## 核心流程

### A. 内容巡检与下架
1. 管理员进入「知识库 → 工作区」或「知识库 → 文档」列表，按所有者、可见性、审核状态、关键词筛选。
2. 点击某条记录查看详情：元数据（标题、所有者、大小、时间、统计）+ 若为公开内容则可预览正文。
3. 点「下架」→ 弹窗**必填处置理由** → 提交。
4. 后端在事务内：置 `AuditStatus=1`、记录理由/操作人/时间，同时写一条审计日志。
5. 下架生效后：公开 wiki 列表与详情不再返回该内容；其下分享链接访问直接失效；**作者本人登录仍可看到该内容，并看到「已被管理员下架 + 理由」提示**。
6. 管理员可对已下架内容执行「恢复」，同样需填理由并记审计日志。

### B. 分享链接治理
7. 管理员进入「知识库 → 分享链接」，可按所有者、文档、状态（生效中/已禁用/已过期）筛选。
8. 点「强制删除」→ 二次确认 + 填理由 → 软删除该链接并记审计日志，站外访问立即 404。

### C. 删除工作区 / 文档
9. 在列表或详情点「删除」→ **两步确认**（输入标题以确认）+ 必填理由。
10. 后端事务内软删除，工作区删除时级联软删除其目录、文档、分享链接，并写审计日志。

### D. 配额治理
11. 「知识库 → 资源占用」页展示按用户聚合的工作区数、文档数、附件存储占用，支持排序找出 Top N。
12. 「知识库 → 配额设置」页配置三个全局上限，存入 `settings` 表。
13. 博客侧创建工作区 / 创建文档 / 上传附件时校验上限，超限返回明确提示。

### E. 运营概览
14. 「知识库 → 概览」页展示汇总卡片：工作区总数（公开/私有）、文档总数（草稿/已发布/已下架）、分享链接数、附件总存储、近 7 日新增文档趋势、公开 wiki 浏览量 Top 文档。

## 异常处理
| 场景 | 处理方式 |
|------|---------|
| 管理员尝试读取私有文档正文 | 详情接口对 `is_public=false` 或 `status!=1` 的文档**不返回 content/content_html 字段**，前端显示「私有内容，仅展示元数据」 |
| 下架 / 删除未填理由 | `binding:"required,min=2,max=200"`，`utils.BadRequest` |
| 下架已下架的内容 / 恢复未下架的内容 | 幂等处理：状态已是目标值时直接返回成功，不重复写审计日志 |
| 删除不存在或已软删的工作区/文档 | `utils.NotFound`，不泄露存在性 |
| 工作区被下架但其下文档未下架 | 公开可见性取 **AND** 语义：工作区被下架则其下所有文档对外均不可见，无需逐条下架 |
| 分享链接指向的文档已被下架 | `share_service.validateShareLink` 增加下架校验，返回 `ErrShareDisabled` 同类语义码 |
| 作者访问自己被下架的内容 | 正常返回内容，并在响应里带 `audit_status` + `audit_reason`，前端展示告警条 |
| 配额上限被调小，存量用户已超限 | 只在**新增时**校验，不回溯清理存量；后台资源占用页对超限用户标红提示 |
| 配额值非法（负数、非数字） | service 校验并 `utils.BadRequest`；读取时解析失败按「不限制」兜底并记 warn 日志 |
| 并发下架同一文档 | 用条件更新 `WHERE audit_status = 0`，受影响行数为 0 时视为已被他人处置，返回当前状态 |
| 审计日志写入失败 | 与业务变更同事务，写失败则整体回滚，保证处置必留痕 |
| 站内通知发送失败 | **不回滚**处置（与审计日志相反）。记 error 日志，管理员侧仍提示处置成功——通知表故障不应阻塞合规处置 |
| 幂等场景下的重复通知 | 状态已是目标值时直接返回，既不写审计日志也不发通知 |

## 技术设计

### 数据模型

#### 1. 新增审核字段（`Workspace` 与 `Doc` 各加一组）

`internal/models/workspace.go` / `internal/models/doc.go`：

| 字段 | 类型 | GORM tag | 说明 |
|------|------|----------|------|
| `AuditStatus` | `int8` | `gorm:"default:0;index"` | 0=正常，1=已下架 |
| `AuditReason` | `string` | `gorm:"size:255"` | 处置理由 |
| `AuditedBy` | `uint` | `gorm:"index"` | 处置管理员 user_id，0 表示未处置 |
| `AuditedAt` | `*time.Time` | — | 处置时间，NULL 表示未处置 |

> 均为新增可空列，`AutoMigrate` 可安全处理，无需迁移 SQL。

常量定义（放 `doc.go` 与 `workspace.go`，与现有 `DocStatusDraft` 风格一致）：
```go
AuditStatusNormal  int8 = 0
AuditStatusBlocked int8 = 1
```

#### 2. 新增表 `admin_audit_logs`（`internal/models/admin_audit_log.go`）

追加型表，**无软删除**（审计记录不可删）：

| 字段 | 类型 | GORM tag | 说明 |
|------|------|----------|------|
| `ID` | `uint` | `gorm:"primarykey"` | |
| `AdminID` | `uint` | `gorm:"index;not null"` | 操作管理员 |
| `AdminUsername` | `string` | `gorm:"size:50"` | 冗余存储，避免关联查询 |
| `Action` | `string` | `gorm:"size:32;index;not null"` | `block`/`unblock`/`delete` |
| `TargetType` | `string` | `gorm:"size:32;index;not null"` | `workspace`/`doc`/`share_link` |
| `TargetID` | `uint` | `gorm:"index;not null"` | 目标主键 |
| `TargetTitle` | `string` | `gorm:"size:255"` | 冗余标题，目标删除后仍可读 |
| `OwnerID` | `uint` | `gorm:"index"` | 被处置内容的所有者 |
| `Reason` | `string` | `gorm:"size:255"` | 处置理由 |
| `IP` | `string` | `gorm:"size:45"` | 操作来源 IP |
| `CreatedAt` | `time.Time` | | |

需加入 `internal/database/migrate.go` 的 AutoMigrate 列表。

#### 3. 通知支持关联文档（`internal/models/notification.go`）

`Notification` 现有 `ArticleID/WorkID/CommentID` 三个关联字段，但没有文档关联。新增一列：

| 字段 | 类型 | GORM tag | 说明 |
|------|------|----------|------|
| `DocID` | `*uint` | `gorm:"index"` | 相关知识库文档 ID |

同步在 `NotificationResponse` 与 `ToResponse()` 中透出 `doc_id`。新增可空列，`AutoMigrate` 安全。

> 工作区被下架时 `DocID` 为空，通知内容里写明工作区名称即可；本期不为工作区单独加关联字段。

#### 4. 配额配置（复用现有 `settings` 表，不建新表）

| key | 默认值 | 说明 |
|-----|--------|------|
| `knowledge.max_workspaces_per_user` | `10` | 每用户最大工作区数，0=不限 |
| `knowledge.max_docs_per_workspace` | `500` | 每工作区最大文档数，0=不限 |
| `knowledge.max_storage_per_user_mb` | `1024` | 每用户知识库附件总存储 MB，0=不限 |

### 可见性判定的变更

现有公开判定是 `workspaces.is_public = true AND docs.status = 1 AND workspaces.owner_id = docs.owner_id`（`public_wiki_service.go:24-27/44-46/62-64/135-137`）。本期在**每一处**追加：

```
AND workspaces.audit_status = 0 AND docs.audit_status = 0
```

分享链接侧（`share_service.go:141-149` 的 `validateShareLink`）追加文档与工作区的 `audit_status` 校验。

> ⚠️ 调研发现 `share_service.Public()` 当前**不校验 `doc.Status`**，草稿文档也能通过分享链接访问——这是与公开 wiki 不同的既有语义。本期只补 `audit_status` 校验，**不改动**草稿可分享这一既有行为，避免破坏用户已发出的链接。

### 权限模型

管理员**不绕过** `workspace_permission.go` 的 owner/成员校验。admin 侧全部走**独立的只读聚合查询**（仿 `public_wiki_service.go` 的写法），处置类写操作也是独立方法，不复用博客侧 service 的 `(id, userID)` 签名。这样：
- 博客侧 10+ 个现有调用点零改动；
- 「管理员不可读私有正文」这条约束由查询本身保证，而非靠调用方自觉。

### API 接口

全部挂在 `api.Group("/admin/knowledge")` + `middleware.AdminAuthMiddleware()` 下（`internal/router/admin.go`）。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/knowledge/overview` | 运营概览统计 |
| GET | `/admin/knowledge/workspaces` | 工作区列表（筛选：owner_id/is_public/audit_status/keyword + 分页） |
| GET | `/admin/knowledge/workspaces/:id` | 工作区详情（元数据 + 目录/文档计数） |
| PUT | `/admin/knowledge/workspaces/:id/audit` | 下架/恢复工作区（body: `audit_status`, `reason`） |
| DELETE | `/admin/knowledge/workspaces/:id` | 删除工作区（级联软删，body: `reason`） |
| GET | `/admin/knowledge/docs` | 文档列表（筛选：owner_id/workspace_id/status/audit_status/kind/keyword + 分页） |
| GET | `/admin/knowledge/docs/:id` | 文档详情（公开文档含正文，私有仅元数据） |
| PUT | `/admin/knowledge/docs/:id/audit` | 下架/恢复文档 |
| DELETE | `/admin/knowledge/docs/:id` | 删除文档 |
| GET | `/admin/knowledge/shares` | 分享链接列表（筛选：owner_id/doc_id/state + 分页，state 含"已过期"计算态） |
| DELETE | `/admin/knowledge/shares/:id` | 强制删除分享链接 |
| GET | `/admin/knowledge/usage` | 按用户聚合的资源占用（工作区数/文档数/存储 MB，支持排序） |
| GET | `/admin/knowledge/quota` | 读取全局配额配置 |
| PUT | `/admin/knowledge/quota` | 更新全局配额配置 |
| GET | `/admin/knowledge/audit-logs` | 审计日志列表（筛选：admin_id/target_type/action/时间范围 + 分页） |

响应一律用 `utils.Success` / `utils.PageResponse`，成功 `code=0`。

### 下架通知

处置生效后给作者发**系统通知**（`FromUserID = nil`），实现方式完全对齐现有的 `CreateWorkAuditNotification`（`notification_service.go:288-332`）。

新增两个方法（`internal/service/notification_service.go`）：

```go
CreateDocAuditNotification(docID uint, blocked bool, reason string) error
CreateWorkspaceAuditNotification(workspaceID uint, blocked bool, reason string) error
```

| 场景 | Type | 文案 |
|------|------|------|
| 文档下架 | `doc_audit` | `你的文档《{标题}》已被管理员下架。原因：{理由}` |
| 文档恢复 | `doc_audit` | `你的文档《{标题}》已恢复正常访问` |
| 工作区下架 | `workspace_audit` | `你的知识库《{名称}》已被管理员下架，其中的文档将不再对外可见。原因：{理由}` |
| 工作区恢复 | `workspace_audit` | `你的知识库《{名称}》已恢复正常访问` |

接收者取内容的 `OwnerID`。文档类通知带上 `DocID` 便于前端跳转。

**事务边界**：通知与业务变更**不同事务**。审计日志必须与处置同事务（写失败则回滚，保证留痕），而通知发送失败只记 error 日志、不回滚处置——避免通知表故障阻塞合规处置。

> 删除工作区/文档时**不发通知**（内容已不存在，通知点进去是 404）。仅下架/恢复发通知。

### 存储占用统计口径

`Attachment` 表未区分知识库附件与博客图片，且本期**不加 `scope` 字段**（属破坏性变更）。按用户聚合存储时，通过 `Doc.AttachmentID` 反查：

```sql
SELECT d.owner_id, COALESCE(SUM(a.file_size), 0) AS bytes
FROM docs d
JOIN attachments a ON a.id = d.attachment_id
WHERE d.deleted_at IS NULL AND a.deleted_at IS NULL
GROUP BY d.owner_id
```

用 GORM 的 `Model().Joins().Group().Select()` 表达，不拼裸 SQL。

两点已知偏差，在页面上以脚注说明，不做精确化：
- `PublishedAttachmentID` 指向的发布态附件不计入（通常与 `AttachmentID` 同源）；
- 同一附件被多个文档引用时会被重复计数（`Attachment.UsageCount > 1` 的情况）。

### 概览趋势查询口径

近 7 日新增文档趋势实时查询，不预聚合、不加定时任务：

```sql
SELECT DATE(created_at) AS day, COUNT(*) AS count
FROM docs
WHERE deleted_at IS NULL AND created_at >= ?
GROUP BY day ORDER BY day
```

`docs.created_at` 无索引，数据量增长后此查询会变慢——若后续概览页出现明显延迟，再考虑加索引或转预聚合（本期不做）。返回时需**补齐零值日期**（SQL 不会返回无数据的日子），保证前端图表有连续 7 个点。

### 实现步骤（每步可独立 commit）

1. [ ] **数据模型**：`Workspace`/`Doc` 加审核四字段 + 常量；`Notification` 加 `DocID`；新增 `models/admin_audit_log.go`；加入 `migrate.go` AutoMigrate 列表。
2. [ ] **可见性收口**：`public_wiki_service.go` 四处查询追加 `audit_status = 0`；`share_service.validateShareLink` 增加下架校验 + 新错误 `ErrContentBlocked`（`knowledge_errors.go`）。
3. [ ] **审计日志 service**：`internal/service/admin_audit_service.go`，提供事务内 `Log(tx, entry)` 方法供其他 admin service 调用。
4. [ ] **下架通知**：`notification_service.go` 新增 `CreateDocAuditNotification` / `CreateWorkspaceAuditNotification`。
5. [ ] **知识库管理 service**：`internal/service/knowledge_admin_service.go`，实现列表/详情/下架/删除/资源聚合/概览统计，写操作走 `database.DB.Transaction` 并在同事务写审计日志，事务提交后异步发通知。
6. [ ] **配额 service**：`knowledge_quota.go` 读写 settings；在 `workspace_service.Create`、`doc_service.Create`、`knowledge_file_service` 上传处接入校验。
7. [ ] **handler + 路由**：`internal/handler/knowledge_admin_handler.go`；在 `router/admin.go` 注册 `/admin/knowledge` 分组。
8. [ ] **前端 - 概览与列表页**：`web/admin/src/views/admin/Knowledge{Overview,Workspaces,Docs,Shares}.vue`。
9. [ ] **前端 - 配额与审计页**：`Knowledge{Usage,Quota,AuditLogs}.vue`。
10. [ ] **前端 - 菜单与路由**：`router/index.js` 加 children；`layouts/AdminLayout.vue` 新增 `el-sub-menu`「知识库」（本项目首次引入子菜单）。
11. [ ] **博客前端下架提示**：作者查看被下架内容时展示 `el-alert` 告警条 + 理由；通知列表支持 `doc_audit`/`workspace_audit` 类型的图标与跳转。

### 参考的现有模式
- `internal/handler/user_handler.go:129-148`（GetUserList）— 分页三段式：`ShouldBindQuery` → page/page_size 兜底（>100 截回）→ `utils.PageResponse`
- `internal/handler/work_handler.go:325-378`（SetRecommend / UpdateWorkStatus）— 状态类写接口：`strconv.ParseUint(c.Param("id"))` + 匿名 inline struct + `binding:"required,oneof=..."`
- `internal/service/user_service.go:155-207` — service 侧列表查询：`Model()` → 逐个 `if` 拼 `Where` → `Count` → `Order/Offset/Limit/Find`，返回 `(list, total, error)`
- `internal/service/public_wiki_service.go:24-27,44-46` — 绕开权限层的独立聚合查询写法；注意它手写 `deleted_at IS NULL`（因用了 `Joins`）
- `internal/service/work_service.go` 的 quota 逻辑 — 配额校验的既有先例
- `internal/middleware/admin_auth.go:46-48` — 校验后 set `user_id`(uint)/`username`/`role`，审计日志的操作人从这里取
- `web/admin/src/views/admin/Users.vue` — 标准列表页范本（el-card 筛选 + el-table + el-pagination + `ElMessageBox.confirm`）
- `web/admin/src/views/admin/Works.vue:662-679` — `ElMessageBox.prompt` 收集审核意见，本期「填写处置理由」可直接复用
- `web/admin/src/utils/adminApi.js` — 唯一 axios 实例，响应拦截已剥出 `data`、按 `code===0` 判成功；页面直接 `adminApi.get('/admin/knowledge/...')`

## 测试计划
- [ ] `public_wiki_service_test.go` 补充：工作区/文档 `audit_status=1` 时不出现在公开列表、树、详情中
- [ ] `share_service_test.go` 补充：文档被下架时 `validateShareLink` 返回 `ErrContentBlocked`；工作区被下架时同样失效
- [ ] `knowledge_admin_service_test.go`：列表筛选条件组合、私有文档详情不返回正文字段、下架幂等（重复下架不重复记日志）
- [ ] `knowledge_quota_test.go`：table-driven 覆盖 0=不限、正常值、超限、配置值非法兜底
- [ ] 通知文案：下架/恢复 × 文档/工作区 四种组合的 content 拼装；系统通知 `FromUserID` 为 nil
- [ ] 概览趋势：无数据日期补零，返回连续 7 个点
- [ ] `handler` 层：`httptest` 验证下架接口缺 `reason` 返回 400、删除接口返回 `code=0`
- [ ] `router` 层：`/admin/knowledge/*` 全部挂在 `AdminAuthMiddleware` 下，无鉴权访问返回 401
- [ ] 保持纯单元测试，不依赖真实 MySQL/Redis

## 已决事项
以下四点已与需求方确认（2026-09-15）：

1. **审计日志永久保留** — 不做定时清理任务，`admin_audit_logs` 表只增不删。
2. **下架必须发站内通知** — 复用 `notification_service`，见上文「下架通知」小节。
3. **存储占用走 `Doc.AttachmentID` 反查** — 不给 `Attachment` 加 `scope` 字段，避免破坏性变更。
4. **概览趋势实时查询** — `GROUP BY DATE(created_at)` 实时算，不做预聚合、不加定时任务。

## MVP 范围

**包含**（步骤 1–9）：
- 工作区 / 文档 / 分享链接三张列表 + 筛选 + 分页
- 下架 / 恢复 / 删除三类处置，全部必填理由并落审计日志
- 审计日志查询页
- 资源占用聚合页 + 全局配额配置与创建时校验
- 运营概览统计页
- 公开 wiki 与分享链接的下架可见性收口

**二期**：
- 步骤 10（博客侧作者告警条）可与 MVP 并行，也可后置
- 用户侧举报入口与工单流转
- 按用户细粒度配额
- 审计日志定时清理
- 下架站内通知
