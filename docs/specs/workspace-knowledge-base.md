# 工作空间 / 私有知识库文档系统

## 状态
- 创建日期: 2026-07-10
- 状态: 草稿

## 目标
在 InkSpace 博客前端（`web/blog` + `cmd/server` :8081）内，为登录用户提供一套「类语雀」的私有内容管理能力：用户可创建多个**工作空间（知识库）**，在空间内用**多层树形目录**分类，撰写并发布**文档**，并可为私有文档创建**免登录分享链接**（永久 / 过期两种，且可随时撤销）。

## 非目标
- 不做多人协作 / 空间成员 / 权限角色体系（本期为**单人多空间**，空间仅归属创建者）。
- 不复用现有博客 `Article` 表；私有文档走**独立 Doc 表**，与博客文章体系解耦。
- 不做分享链接的密码保护（本期只做永久/过期 + 撤销）。
- 分享对象只到**单篇文档**；不分享整个目录 / 整个空间（后续迭代）。
- 不引入 `vue-i18n`、不引入新的 UI/Markdown 库（沿用 Element Plus + vditor + axios + Pinia）。
- 不做实时协同编辑（OT/CRDT）。

## 用户故事
- 作为登录用户，我想创建自己的工作空间，以便把不同主题的笔记分库管理。
- 作为用户，我想在空间里建立多层目录并把文档挂在目录下，以便像语雀一样分类整理。
- 作为用户，我想撰写 Markdown 文档并在草稿 / 发布间切换，以便未定稿的内容保持私有。
- 作为用户，我想给某篇文档生成一个分享链接（可设永久或过期时间），以便把私有内容发给站外的人只读查看。
- 作为用户，我想随时禁用或删除已发出的分享链接，以便及时收回访问权。
- 作为访问者，我想通过分享链接**免登录**直接阅读文档，以便无需注册就能查看对方分享的内容。

## 核心流程

### A. 空间与目录
1. 用户进入「我的知识库」→ 新建工作空间（名称、描述、可选图标）。
2. 进入空间，左侧为目录树；在根或任意目录下「新建子目录」。
3. 目录支持重命名、删除（删除目录时其下文档与子目录一并软删除，需二次确认）、拖拽调整父级与顺序。

### B. 文档撰写与发布
4. 在某目录（或空间根）下「新建文档」→ 打开 vditor 编辑器（复用 `VditorEditor.vue`）。
5. 编辑中每隔一段时间**自动保存**草稿（`autosave` 接口，不产生版本）。
6. 用户点「保存」→ 落库并生成一条**版本快照**（`DocVersion`）。
7. 用户点「发布」→ 状态置为已发布、渲染并存 `ContentHTML`、记录 `PublishedAt`。
8. 用户可在版本历史里查看/回滚到某个历史版本。

### C. 分享
9. 在文档页点「分享」→ 选择永久有效 or 设置过期时间 → 生成随机 token 短码链接（`/share/:token`）。
10. 同一文档可有多条分享链接；列表可查看状态、访问次数，可**禁用**或**删除**。
11. 站外访问者打开 `/share/:token`：后端校验链接启用中且未过期 → 返回文档只读内容渲染页；否则显示「链接已失效 / 已过期」。

## 异常处理
| 场景 | 处理方式 |
|------|---------|
| 访问不属于自己的空间/目录/文档（编辑接口） | service 层追加 `owner_id = 当前用户` 过滤，查不到返回 404（`utils.NotFound`），不泄露存在性 |
| 分享 token 不存在 | `utils.NotFound`，前端展示「链接不存在」 |
| 分享链接已禁用 | `utils.Forbidden` / 语义码，前端展示「链接已被作者关闭」 |
| 分享链接已过期 | 同上，前端展示「链接已过期」 |
| 移动目录形成环（把父目录拖进自己子孙下） | service 校验目标 parent 不是自身或后代，拒绝并 `utils.BadRequest` |
| 删除含子目录/文档的目录 | 事务内级联软删除；前端二次确认弹窗 |
| 自动保存与手动保存并发 | autosave 只写 `content`，不建版本；以 `updated_at` 为准，后写覆盖（本期不做乐观锁） |
| 回滚到历史版本 | 把该版本内容写回当前文档，并新建一条「回滚生成」的版本快照，保证可再回滚 |
| 空间/文档数量为空 | 前端空状态引导「新建」 |

## 技术设计

### 数据模型
> 全部内嵌 `gorm.DeletedAt`（`gorm:"index" json:"-"`）软删除；主键 `gorm:"primarykey"`；新增后须加入 `internal/database/migrate.go` 的 `AutoMigrate` 列表。请求体 `XxxRequest`（`binding:"required"`）与 `ToResponse()` 与 model 同文件。

**Workspace（工作空间/知识库）** `internal/models/workspace.go`
| 字段 | 类型 / tag | 说明 |
|------|-----------|------|
| ID | uint primarykey | |
| OwnerID | uint `gorm:"index;not null"` | 归属用户 |
| Name | string `gorm:"size:100;not null"` | 空间名 |
| Description | string `gorm:"size:500"` | 简介 |
| Icon | string `gorm:"size:255"` | 图标/emoji（美化，可空） |
| Sort | int `gorm:"default:0"` | 排序 |
| DocCount | int `gorm:"default:0;not null"` | 文档数（事务内维护） |
| CreatedAt / UpdatedAt / DeletedAt | | |

**Catalog（目录节点，多层树）** `internal/models/catalog.go`
| 字段 | 类型 / tag | 说明 |
|------|-----------|------|
| ID | uint primarykey | |
| WorkspaceID | uint `gorm:"index;not null"` | 所属空间 |
| ParentID | *uint `gorm:"index"` | nil=空间根节点 |
| OwnerID | uint `gorm:"index;not null"` | 冗余，便于鉴权过滤 |
| Name | string `gorm:"size:100;not null"` | 目录名 |
| Sort | int `gorm:"default:0"` | 同级排序（拖拽） |
| CreatedAt / UpdatedAt / DeletedAt | | |

> 树由 `ParentID` 自关联；service 层一次性拉出空间内全部目录在内存组树，避免递归查询。参考 `Category.ParentID` / `Comment.ParentID` 先例。

**Doc（文档）** `internal/models/doc.go`
| 字段 | 类型 / tag | 说明 |
|------|-----------|------|
| ID | uint primarykey | |
| WorkspaceID | uint `gorm:"index;not null"` | 所属空间 |
| CatalogID | *uint `gorm:"index"` | nil=挂空间根 |
| OwnerID | uint `gorm:"index;not null"` | 归属用户 |
| Title | string `gorm:"size:200;not null"` | 标题 |
| Content | string `gorm:"type:longtext"` | Markdown 原文 |
| ContentHTML | string `gorm:"type:longtext"` | 渲染 HTML（发布时生成） |
| Status | int `gorm:"index;default:0;not null"` | 0=草稿 1=已发布 |
| WordCount | int `gorm:"default:0"` | |
| ViewCount | int `gorm:"default:0"` | 通过分享链接阅读累计 |
| Sort | int `gorm:"default:0"` | 目录内排序（拖拽） |
| PublishedAt | *time.Time | 首次发布时间 |
| CreatedAt / UpdatedAt / DeletedAt | | |

**DocVersion（版本快照）** `internal/models/doc_version.go`
| 字段 | 类型 / tag | 说明 |
|------|-----------|------|
| ID | uint primarykey | |
| DocID | uint `gorm:"index;not null"` | |
| Version | int `gorm:"not null"` | 文档内递增版本号 |
| Title | string `gorm:"size:200"` | 快照标题 |
| Content | string `gorm:"type:longtext"` | 快照原文 |
| OwnerID | uint `gorm:"index"` | |
| Remark | string `gorm:"size:100"` | 如「手动保存」「回滚自 v3」 |
| CreatedAt | | 无需软删除 |

**ShareLink（分享链接）** `internal/models/share_link.go`
| 字段 | 类型 / tag | 说明 |
|------|-----------|------|
| ID | uint primarykey | |
| Token | string `gorm:"uniqueIndex;size:32;not null"` | 随机短码 |
| DocID | uint `gorm:"index;not null"` | 分享的文档（本期单篇） |
| OwnerID | uint `gorm:"index;not null"` | 创建者 |
| ExpiresAt | *time.Time | nil=永久有效；否则为过期时刻 |
| Enabled | bool `gorm:"default:true;not null"` | 撤销=置 false |
| ViewCount | int `gorm:"default:0"` | 访问计数 |
| CreatedAt / UpdatedAt / DeletedAt | | |

> 有效性判定：`Enabled == true && (ExpiresAt == nil || now < *ExpiresAt)`。token 用 `crypto/rand` 生成的 URL-safe 短码（约 22~32 字符）。

### API 接口
> 全部走 `SetupUserRouter`（`internal/router/blog.go`）。除注明 public 外均在 `protected`（`AuthMiddleware`）分组下；响应统一用 `internal/utils/response.go`。

**工作空间**
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/workspaces | 创建空间 |
| GET | /api/workspaces | 我的空间列表 |
| GET | /api/workspaces/:id | 空间详情 |
| PUT | /api/workspaces/:id | 更新（名称/简介/图标/排序） |
| DELETE | /api/workspaces/:id | 删除（级联软删目录与文档） |

**目录**
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/workspaces/:id/catalogs | 新建目录 |
| GET | /api/workspaces/:id/catalogs | 拉取该空间目录树 |
| PUT | /api/catalogs/:id | 重命名 |
| DELETE | /api/catalogs/:id | 删除（级联软删） |
| PUT | /api/catalogs/:id/move | 移动/排序（parent_id + sort，校验环） |

**文档**
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/docs | 新建文档 |
| GET | /api/workspaces/:id/docs | 文档列表（可按 catalog_id 过滤） |
| GET | /api/docs/:id/edit | 作者取原文编辑 |
| PUT | /api/docs/:id | 保存（建版本快照） |
| PUT | /api/docs/:id/autosave | 自动保存草稿（不建版本） |
| POST | /api/docs/:id/publish | 发布（渲染 HTML + PublishedAt） |
| DELETE | /api/docs/:id | 删除文档 |
| PUT | /api/docs/:id/move | 移动到其它目录 / 排序 |
| GET | /api/docs/:id/versions | 版本列表 |
| GET | /api/docs/:id/versions/:version | 单版本内容 |
| POST | /api/docs/:id/versions/:version/rollback | 回滚 |
| GET | /api/workspaces/:id/search?q= | 空间内标题/内容搜索 |

**分享**
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/docs/:id/shares | 创建分享链接（body: 永久 or 过期时间） |
| GET | /api/docs/:id/shares | 该文档的分享链接列表 |
| PUT | /api/shares/:id | 启用/禁用 / 修改过期时间 |
| DELETE | /api/shares/:id | 删除分享链接 |
| GET | /api/share/:token | **public，免登录**：校验有效性并返回只读文档，累加 ViewCount |

### 实现步骤（每步可独立 commit）
1. [ ] **数据模型**：新增 `workspace.go`、`catalog.go`、`doc.go`、`doc_version.go`、`share_link.go`（含 Request/Response），全部加入 `migrate.go` 的 AutoMigrate。
2. [ ] **workspace service + handler + 路由**：空间 CRUD，缓存 key `workspace:%d`。
3. [ ] **catalog service + handler + 路由**：目录 CRUD + 组树 + move（环校验、级联软删）。
4. [ ] **doc service + handler + 路由**：文档 CRUD + autosave + publish（Markdown→HTML 渲染，复用现有文章渲染方式）。
5. [ ] **doc version service + handler + 路由**：保存建快照 + 版本列表 + 回滚。
6. [ ] **share service + handler + 路由**：创建/列表/启用禁用/删除 + public `/api/share/:token` 只读访问（token 用 crypto/rand 生成，有效性判定 + ViewCount）。
7. [ ] **搜索接口**：空间内 `LIKE` 标题/内容（MVP 简单实现，非全文索引）。
8. [ ] **前端-空间与目录**：`views/user/WorkspaceList.vue`、`WorkspaceDetail.vue`（左侧 `el-tree` 目录 + 文档列表 + 拖拽），Pinia store，路由挂 `/dashboard/*` 需登录。
9. [ ] **前端-文档编辑**：`views/user/DocEdit.vue`（复用 `VditorEditor.vue` + 自动保存定时器 + 版本抽屉 + 分享弹窗）。
10. [ ] **前端-公开阅读页**：`views/ShareDoc.vue`，路由 `/share/:token`，走 `MainLayout` 免登录，失效/过期空状态。

### 参考的现有模式
- `internal/models/category.go` — `ParentID *uint` 树形先例、`Slug` 唯一生成、`XxxCount` 冗余计数在 model 的写法。
- `internal/models/comment.go` — 完整的父子层级（ParentID/RootID）实现，可参考目录树查询。
- `internal/models/article.go` — Status（草稿/发布/私有）、Content + ContentHTML 双存、`ToResponse()` 隐藏字段的模式，Doc 直接照搬字段风格。
- `internal/service/article_service.go` — 事务 `database.DB.Transaction`、`gorm.Expr` 原子计数、`SetCache/GetCache/DeleteCachePattern` 缓存、service 层 `author_id` 鉴权过滤。
- `internal/middleware/auth.go` — `AuthMiddleware`（强制登录）与 `OptionalAuthMiddleware`（可选登录，用于 public 分享页）。
- `internal/router/blog.go` `SetupUserRouter` — public / publicWithOptionalAuth / protected 三分组注册方式。
- `web/blog/src/components/VditorEditor.vue` — v-model 编辑器组件，直接复用于 DocEdit。
- `web/blog/src/views/user/ArticleEdit.vue` — 新建/编辑判断、状态切换、提交后跳转的前端范式。
- `web/blog/src/utils/api.js`、`src/stores/user.js`、`src/router/index.js` — axios 实例/拦截器、Pinia、路由守卫复用。

## 测试计划
- [ ] service 层：目录 move 环校验（把父拖进子孙）应报错；级联软删应连带子目录与文档。
- [ ] service 层：分享有效性判定（永久 / 未过期 / 已过期 / 已禁用）四态。
- [ ] service 层：非 owner 访问他人空间/文档返回 404，不泄露存在性。
- [ ] service 层：回滚后能再次回滚（回滚也生成新版本）。
- [ ] handler 层：`GET /api/share/:token` 免登录可访问，失效链接返回对应语义码；成功响应 `code==0`。
- [ ] handler 层：受保护接口缺 token 返回 401。

## 待定事项
- 目录/文档删除是「软删进回收站可恢复」还是「直接软删不可见」？（本期按不可见处理，回收站列后续迭代）
- 自动保存频率与并发覆盖策略是否需要乐观锁（`updated_at` 版本校验）？本期先「后写覆盖」。
- 搜索是否需要高亮 / 分页？本期先返回标题+摘要列表。
- 分享链接域名前缀（前端 `window.location.origin + /share/:token`）确认。
- Markdown → HTML 的服务端渲染库确认（沿用现有文章发布所用方式）。

## MVP 范围
**本期全部纳入**（用户确认）：
1. 工作空间 CRUD（单人多空间）。
2. 多层树形目录 CRUD + 拖拽移动/排序。
3. 文档 CRUD + 草稿/发布 + **自动保存**。
4. **版本历史** + 回滚。
5. 分享链接：永久 / 过期两种 + 撤销禁用 + token 免登录只读阅读页。
6. 空间内搜索（简单 LIKE）。

后续迭代：分享整个目录/空间、分享密码保护、回收站恢复、全文检索、多人协作与权限。
