# 用户侧多 UI 主题系统

## 状态
- 创建日期: 2026-07-22
- 状态: Phase 1「屿刊」及附加功能已完成，并于 2026-07-22 通过用户验收
- 首期主题: A · 极简杂志风「屿刊」
- 当前阶段: Phase 2「yu.log」实现与自动化验证完成，人工视觉矩阵待验收
- 后续顺序: B · `yu.log` -> C ·「小屿的角落」-> D · `CHEN YU®`
- Phase 2 规范: [`terminal-ui-theme.md`](./terminal-ui-theme.md)

## 设计关联

本 Spec 以 `docs/design/README.md` 的共用内容模型和全局基线为准。每套主题的视觉实现、交互和验收必须同时关联对应 Markdown 规范与 HTML 原型，不以开发者自行发挥替代设计稿。

| 主题 ID | 展示名称 | 规范 | 高保真原型 | 实现阶段 |
|---|---|---|---|---|
| `magazine` | 屿刊 · 极简杂志风 | [`style-a-magazine.md`](../design/style-a-magazine.md) | [`style-a-magazine.html`](../design/style-a-magazine.html) | Phase 1 已完成，默认主题 |
| `terminal` | yu.log · 暗色科技感 | [`style-b-terminal.md`](../design/style-b-terminal.md) | [`style-b-terminal.html`](../design/style-b-terminal.html) | Phase 2 已开放，人工视觉矩阵待验收 |
| `cozy` | 小屿的角落 · 温暖手作感 | [`style-c-cozy.md`](../design/style-c-cozy.md) | [`style-c-cozy.html`](../design/style-c-cozy.html) | 第三阶段 |
| `swiss` | CHEN YU® · 瑞士网格风 | [`style-d-swiss.md`](../design/style-d-swiss.md) | [`style-d-swiss.html`](../design/style-d-swiss.html) | 第四阶段 |

共同基线：

- 视觉风格与明暗模式是两个独立维度：`data-ui-theme` 控制风格，`data-theme="light|dark"` 控制明暗。
- 所有主题共享业务路由、数据加载、表单、权限和 API 客户端；主题只替换布局组合、展示组件和设计令牌。
- 颜色必须来自 CSS 自定义属性；业务组件不得写死主题色。
- 响应式断点统一为 `900px` 和 `560px`；窄屏导航使用可操作的菜单，不直接隐藏入口。
- 支持键盘焦点、正文对比度、图片替代文本、懒加载和 `prefers-reduced-motion`。

## 目标

在不破坏现有文章、作品、评论、用户资料等 API 契约和业务逻辑的前提下，为 `web/blog` 建立可扩展的多 UI 主题架构；登录用户可在用户中心预览并保存账号级主题与明暗偏好，未登录访客默认使用「屿刊」并可在当前浏览器切换明暗模式。首期完整适配「屿刊」，并补齐设计信息架构中的摄影和公共知识库入口。

## 非目标

- 不重构 `web/admin` 的整体视觉；仅在既有设置页中增加首页静态 Hero 文案配置，沿用原保存契约。
- 不修改现有文章、作品、评论、关注、收藏、通知、用户资料等 API 的既有字段含义、路径或成功响应结构。
- 不为摄影新增数据表、上传流程或独立内容模型；摄影继续使用 `Work.Type = photography`。
- 不为公共知识库新增编辑器、协作、成员权限、全文检索、双向链接或新内容模型。
- 不允许主题包自行请求业务 API 或复制页面业务状态；主题不是四套独立 SPA。
- Phase 1 不实现 B/C/D 的完整页面样式；只注册其元数据和设计关联，外观页标记为“即将推出”。
- 不引入 `vue-i18n`、新的 UI 框架或 CSS 框架；继续使用 Vue 3、Pinia、Vue Router、Element Plus 和现有 Vditor。

## 用户故事

- 作为未登录访客，我打开网站时希望默认看到「屿刊」并跟随系统明暗模式，以获得稳定且符合设备偏好的阅读体验。
- 作为未登录访客，我想临时切换浅色、深色或跟随系统，以便在当前浏览器舒适阅读。
- 作为登录用户，我想在用户中心的“外观与主题”页面预览主题，再确认保存，以便主题偏好在不同设备登录后同步。
- 作为登录用户，我希望取消预览或保存失败时恢复原主题，以免误操作改变当前体验。
- 作为站点访问者，我想通过 `/photos` 浏览已发布的摄影作品，通过 `/wiki` 阅读公开知识库中的已发布文档。
- 作为知识库作者，我想决定某个工作区是否公开，并确保草稿或再次编辑中的内容不会被公共知识库泄露。
- 作为维护者，我想让 B/C/D 复用同一主题注册表、页面槽位和数据适配层，以便逐个适配时无需改后端契约或重写业务页面。

## 核心流程

### A. 首次访问与启动

1. 应用在 Vue 挂载前读取本地缓存，立即给 `<html>` 设置 `data-ui-theme` 和 `data-theme`，避免明显闪烁。
2. 无缓存时使用 `ui_theme=magazine`、`color_scheme=system`；`system` 根据 `prefers-color-scheme` 解析为实际 `light` 或 `dark`。
3. 同时读取 `GET /api/settings/public`：若管理员设置 `site_theme=holiday` 或 `mourning`，添加独立站点覆盖状态，但不改写用户保存的主题偏好。
4. 登录用户读取自己的主题偏好，以服务端账号值为准并写入当前用户专属缓存；未登录访客继续使用独立的访客缓存，两者不得互相覆盖。
5. 当偏好为 `system` 时监听系统明暗变化并即时更新；显式 `light`/`dark` 不跟随系统变化。

### B. 用户预览并确认主题

1. 登录用户进入 `/dashboard/appearance`，看到四套主题卡片、当前主题、明暗模式选择和设计预览。
2. Phase 1 中 `magazine` 可选择；`terminal`、`cozy`、`swiss` 展示名称、缩略预览、设计说明和“即将推出”，不可提交。后续主题仅在其独立阶段完成全矩阵验收后解锁。
3. 用户点“预览”后，仅在当前会话临时应用候选主题和明暗模式，不调用保存接口。
4. 用户点“取消预览”或离开页面且未确认，恢复进入页面前的已保存偏好。
5. 用户点“应用主题”，调用主题偏好接口；成功后更新 Pinia 状态和本地缓存，并保持全站立即生效。
6. 保存失败时恢复原主题并显示错误提示，不保留未同步的本地预览。

### C. 未登录访客切换明暗模式

1. 公共页头提供明暗模式快捷入口，支持 `system`、`light`、`dark`。
2. 选择只写本地缓存，不要求登录，不调用用户偏好接口。
3. 登录后若账号已有偏好，以账号偏好覆盖访客本地值；若账号使用默认值，则仍以账号返回值为准。

### D. 管理员特殊主题覆盖

1. 现有 `site_theme=day/night` 作为兼容值继续被后端和管理端保存，但新用户侧运行时忽略这两个值；明暗模式由用户 `color_scheme` 或访客本地值决定。本期按用户确认保持 `web/admin` 原样，管理端旧文案与新语义的差异作为已知兼容限制记录，不扩大本期范围。
2. `site_theme=holiday` 时保留当前节日颜色配置，并作为四套 UI 主题之上的临时色彩覆盖。
3. `site_theme=mourning` 时保留全站灰度覆盖。
4. 特殊覆盖期间，用户仍可调整和保存 UI 主题/明暗偏好；偏好不被覆盖值改写。打开页面时读取覆盖状态；页面保持打开时，在窗口重新获得焦点且距上次读取超过 5 分钟后刷新一次，覆盖结束后恢复个人偏好，不引入轮询或实时推送。

### E. 摄影浏览

1. 访问 `/photos`，使用现有作品列表能力固定筛选 `type=photography`，展示已发布摄影作品。
2. 点击摄影作品进入现有 `/works/:id` 详情，不新增摄影详情 API；页面根据作品类型使用主题化摄影布局。
3. 公共作品列表服务端强制 `status=1`，忽略或拒绝客户端传入的其他状态，避免草稿、审核中或驳回作品泄露。

### F. 公共知识库

1. 作者在现有工作区管理界面设置“公开工作区”；历史工作区迁移后默认不公开。
2. `/wiki` 只列出公开工作区，并按已发布文档数展示统计。
3. 进入工作区后，只展示 `status=1` 的文档及其所在目录和必要祖先目录；纯草稿目录分支不公开。
4. 文档详情只返回经过公共阅读白名单净化的已持久化 `content_html`，不返回 Markdown 原文或编辑关系字段。
5. 已发布文档执行手动保存、自动保存或版本回滚后自动转为草稿；重新发布后才重新进入公共知识库。
6. 文档编辑页必须展示草稿/已发布状态，并提供发布、取消发布和重新发布操作；否则作者无法完成公开流程。
7. 公开工作区的名称/简介修改、目录重命名/移动/删除、已发布文档移动会立即反映到公共导航，并在操作确认文案中明确提示；正文或标题修改仍按第 5 条转草稿。
8. token 分享继续保持现有“非公开链接”语义，与工作区公开状态相互独立。

### G. 后续主题适配

1. 主题注册表将对应主题状态从 `coming_soon` 改为 `available`，并加载该主题令牌和展示组件。
2. 复用已有页面数据、查询条件、分页、表单、权限与交互状态，只实现该主题的视觉变体和必要的布局组合。
3. 按 B -> C -> D 顺序逐个完成全页面矩阵和专项验收后开放选择，不允许只完成首页便标记可用。

## 异常处理

| 场景 | 处理方式 |
|---|---|
| 用户偏好读取失败 | 使用有效本地缓存；无缓存则回退 `magazine + system`，不阻塞页面 |
| 用户偏好保存失败 | 恢复保存前主题，保留服务端已保存值，Element Plus 提示同步失败 |
| 服务端或缓存返回未知主题 | 归一化为 `magazine`，不得生成未知动态类名或加载任意资源 |
| 服务端或缓存返回未知明暗偏好 | 归一化为 `system` |
| 用户预览后直接离开外观页 | 路由离开前恢复已保存偏好 |
| 系统明暗媒体查询不可用 | `system` 回退为 `light` |
| 特殊全站主题配置缺少颜色 | `holiday` 使用现有内置节日默认令牌；`mourning` 仍应用灰度 |
| 尚未完成的主题 | 展示“即将推出”，禁用预览和保存；接口也拒绝未开放主题值 |
| 公共工作区不存在或已关闭 | 统一返回 404，不泄露此前是否存在 |
| 文档不存在、未发布或工作区未公开 | 统一返回 404 |
| 公开目录仅包含草稿 | 目录及其分支不出现在公开树中 |
| 已发布文档发生保存/自动保存/回滚 | 同一事务内转为草稿；公共详情立即不可访问，直至重新发布 |
| 公开 HTML 含脚本、事件属性或危险 URL | 公共 Wiki 响应净化后移除；不得依赖前端 `v-html` 作为安全边界 |
| 公共树超过 2,000 个目录与文档节点 | 返回 422 和明确错误，管理界面提示拆分工作区；避免无界公开响应 |
| `/photos` 无摄影作品 | 展示主题化空状态，不回退展示项目作品 |
| 图片加载失败 | 保留比例占位并显示可读替代内容，不造成布局跳动 |

## 技术设计

### 设计原则

1. **单一业务树，多套展示皮肤**：路由和业务视图负责请求、权限、分页、提交及错误状态；共享展示组件接收规范化数据；主题通过令牌、变体组件和页面布局槽位表达。
2. **两个主题轴**：`uiTheme = magazine|terminal|cozy|swiss`，`colorScheme = system|light|dark`；实际 DOM 仅使用白名单归一化后的值。
3. **渐进式适配**：主题注册表记录可用状态和动态样式入口。未完成主题不可保存，已开放主题必须覆盖完整页面矩阵。
4. **旧契约兼容**：不重命名或删除旧 API 字段。新增偏好字段只出现在新的专用偏好接口，避免扩大 `UserResponse` 和公开用户资料。
5. **主题样式隔离**：所有主题令牌限定在 `html[data-ui-theme="..."]`；Element Plus 变量在主题作用域内映射，避免页面级 `!important` 相互覆盖。

### 前端模块

建议最小结构如下，具体命名可在实现时保持仓库风格，但职责边界不可反转：

```text
web/blog/src/
├── components/theme/
│   ├── AppearanceSelector.vue
│   ├── ColorSchemeToggle.vue
│   └── ThemePreview.vue
├── components/site/
│   ├── SiteHeader.vue
│   └── SiteFooter.vue
├── components/content/
│   ├── ArticleItem.vue
│   ├── WorkItem.vue
│   ├── PhotoItem.vue
│   └── WikiItem.vue
├── stores/appearance.js
├── themes/
│   ├── registry.js
│   ├── base.css
│   ├── magazine.css
│   └── terminal.css（Phase 2）
└── views/user/Appearance.vue
```

约束：

- `registry.js` 是主题 ID、名称、状态、设计文档链接、默认模式和资源入口的唯一来源。
- `appearance` store 维护 `savedPreference`、`previewPreference`、实际解析后的明暗模式和特殊全站覆盖。
- `MainLayout.vue` 与 `UserCenterLayout.vue` 复用主题状态；公共布局按主题呈现，用户中心保留完整管理可用性并使用同一设计令牌。
- `Home.vue`、`Blog.vue`、`BlogDetail.vue`、`Works.vue`、`WorkDetail.vue`、`UserProfile.vue` 等保留现有请求和业务逻辑，逐步抽取共享内容项，避免复制四份 API 调用。
- `VditorEditor.vue` 的编辑能力不变；公开正文和编辑器外壳映射到当前主题令牌。`code_theme` 仍沿用现有全站设置，不改变其 API。
- Markdown 正文和 Vditor 跟随用户 UI 的实际浅深模式，包括 `system` 与外观预览；管理后台 `markdown_theme` 保留兼容但不再影响博客前端。代码高亮继续严格使用管理后台 `code_theme`。
- 移除 `main.css` 的重复引入，只保留一个全局入口。
- 首次运行新主题系统时清理旧 `localStorage.site_theme`、`body.theme-*` 和旧工具写入的内联 `--theme-*` 变量，再由新运行时统一设置；迁移只执行一次且不把旧 `day/night` 猜测为个人偏好。

### 页面适配矩阵

“完整用户侧”验收要求每个开放主题覆盖以下页面，不以首页局部换肤视为完成。Phase 1 已完成 `magazine` 适配：

| 区域 | 路由/页面 | 屿刊适配要求 |
|---|---|---|
| 公共骨架 | `MainLayout.vue` | 杂志式页头、窄屏菜单、页脚、明暗切换、节日/哀悼覆盖 |
| 首页 | `/` | Hero、精选文章、作品、摄影、知识库分区；复用真实 API 数据和空状态 |
| 文章 | `/blog`、`/blog/:id` | 杂志索引、筛选分页、正文排版、目录、评论及互动状态 |
| 作品 | `/works`、`/works/:id` | 项目/摄影筛选、详情、图集、评论及互动状态 |
| 摄影 | `/photos` | 摄影专属列表，详情复用 `/works/:id` |
| 知识库 | `/wiki`、`/wiki/:workspaceId`、`/wiki/docs/:id` | 公开空间、目录树、已发布文档阅读 |
| 站点页 | `/about`、`/links`、`/user-search`、`/share/:token` | 统一杂志令牌、排版、表单、空态和错误态 |
| 账号 | `/login`、`/users/:id` | 登录/注册、公开个人主页及其内容标签页 |
| 用户中心骨架 | `/dashboard/*`、`/favorites`、`/profile/edit` | 导航、面包屑、通知、下拉菜单、移动端可用性 |
| 用户内容管理 | 文章、作品、评论、收藏、通知 | 列表、筛选、分页、确认框、编辑入口和状态标签 |
| 知识库管理 | 工作区、目录、文档编辑 | 公开开关、树、拖拽、编辑、草稿状态、发布/取消发布/重新发布、版本、分享完整可用 |
| 外观设置 | `/dashboard/appearance` | 四主题卡、预览/取消/确认、明暗偏好、异常回退 |

### 屿刊主题验收映射

实现必须以 `docs/design/style-a-magazine.md` 和 `docs/design/style-a-magazine.html` 为视觉基准：

- 使用纸张背景、墨色正文、黛绿强调色和发丝分隔线，提供规范中的明/暗两套令牌。
- 标题和长文使用系统衬线栈；UI 控件、元数据使用无衬线栈。
- 内容宽度、留白、字号层级、零阴影和近零圆角遵循设计规范；仅标签等小型状态可使用胶囊形。
- 首页保留期刊号、竖排签名、文章索引、摄影图注等签名元素；真实数据缺字段时使用前端可推导值，不伪造后端字段。
- 文章行、摄影图、项目箭头和知识库卡片使用设计稿规定的克制动效；`prefers-reduced-motion` 时关闭。
- Element Plus 表单、弹窗、分页、下拉菜单等功能组件必须映射到同一令牌体系，不能保留默认蓝色和通用阴影。

### 数据模型

#### UserAppearance（新增）

使用独立一对一表保存账号偏好，避免给 `users` 表持续堆叠展示字段，也不改变现有用户资料请求/响应。

| 字段 | 类型 / tag | 说明 |
|---|---|---|
| ID | uint `primarykey` | |
| UserID | uint `uniqueIndex;not null` | 每个用户至多一条 |
| UITheme | string `size:20;default:'magazine';not null` | 白名单主题 ID |
| ColorScheme | string `size:10;default:'system';not null` | `system|light|dark` |
| CreatedAt / UpdatedAt | time.Time | |

约束：

- 不使用软删除；记录可直接 upsert。
- 新增 `UserAppearanceRequest` 和 `UserAppearanceResponse`，服务端验证白名单。
- 未创建记录等价于 `{ui_theme: "magazine", color_scheme: "system"}`；GET 可返回默认值但不强制落库。
- Phase 1 接口只允许 `magazine`。每个后续主题完整验收并开放时，再同步加入服务端可用白名单。

#### Workspace（新增字段）

| 字段 | 类型 / tag | 说明 |
|---|---|---|
| IsPublic | bool `gorm:"default:false;index;not null"` | 是否允许公共知识库发现 |

同步扩展现有 `WorkspaceRequest` 和 `WorkspaceResponse` 的 `is_public` 字段。该变更是向后兼容的新增字段；历史数据默认 `false`，不得自动公开。

#### Doc（沿用字段，调整状态转换）

- 沿用 `Status=0` 草稿、`Status=1` 已发布和发布时生成的 `ContentHTML`。
- 对 `status=1` 文档执行手动保存、自动保存或版本回滚时，同一事务将 `status` 置为 `0`。
- 公共读取只使用上次明确发布生成的 `ContentHTML`，且必须同时满足当前 `status=1`；原始 `Content` 不进入公共 DTO。
- 发布和重新发布沿用现有 `PublishedAt` 的首次发布时间语义；GORM `UpdatedAt` 记录最近一次明确发布操作。公共列表按 `updated_at DESC, id DESC` 排序，页面分别显示“首次发布”和“最近发布”时必须使用对应字段。
- 取消发布继续复用现有发布接口的 `status=0`，不清空 `ContentHTML` 或 `PublishedAt`，但公共查询必须立即不可见。

### API 接口

所有接口继续使用 `{code, message, data}`，成功 `code=0`。旧接口不重命名、不删除字段、不改变路径。

#### 用户主题偏好（新增，需登录）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/profile/appearance` | 获取当前用户偏好；无记录返回默认值 |
| PUT | `/api/profile/appearance` | 原子保存 `ui_theme` 与 `color_scheme` |

请求与响应 `data`：

```json
{
  "ui_theme": "magazine",
  "color_scheme": "system"
}
```

校验失败返回 400；未开放主题也按无效值处理。独立接口避免改变现有 `PUT /api/profile` 的字段和逻辑。

#### 工作区管理（兼容扩展）

| 方法 | 路径 | 变更 |
|---|---|---|
| POST | `/api/workspaces` | 可选接收 `is_public`，默认 false |
| PUT | `/api/workspaces/:id` | 可更新 `is_public`，仍校验 owner |
| GET | `/api/workspaces` | 响应新增 `is_public` |
| GET | `/api/workspaces/:id` | 响应新增 `is_public` |

#### 公共知识库（新增，免登录）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/wiki/workspaces?page=1&page_size=20` | 公开工作区列表，`doc_count` 为已发布文档数；`page_size` 最大 100 |
| GET | `/api/wiki/workspaces/:id/tree` | 已发布文档和必要目录祖先组成的只读树 |
| GET | `/api/wiki/docs/:id` | 公开文档详情，只返回公共 DTO 和 `content_html` |

公共查询必须在 service 查询条件中同时约束工作区公开、文档已发布、软删除状态和归属一致性。不存在、已关闭、草稿和越权均返回同一种 404。

工作区列表使用现有分页响应：

```json
{
  "list": [
    {
      "id": 1,
      "name": "前端笔记",
      "description": "可公开阅读的前端知识库",
      "icon": "",
      "doc_count": 12,
      "author_id": 7,
      "author_name": "小屿",
      "author_avatar": "/uploads/avatar/example.jpg",
      "updated_at": "2026-07-22T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 20
}
```

工作区按 `updated_at DESC, id DESC` 排序。`doc_count` 必须实时统计已发布文档，不能复用包含草稿的 `Workspace.DocCount`。

树接口固定返回工作区与两组节点；目录和文档分离可直接复用现有目录树构建逻辑，并明确表示根文档：

```json
{
  "workspace": {
    "id": 1,
    "name": "前端笔记",
    "description": "可公开阅读的前端知识库",
    "icon": ""
  },
  "catalogs": [
    {
      "id": 10,
      "parent_id": null,
      "name": "Vue",
      "sort": 0,
      "children": []
    }
  ],
  "docs": [
    {
      "id": 100,
      "catalog_id": 10,
      "title": "响应式原理",
      "sort": 0,
      "published_at": "2026-07-22T00:00:00Z",
      "updated_at": "2026-07-22T00:00:00Z"
    }
  ]
}
```

`catalog_id=null` 表示根文档。目录按 `sort ASC, id ASC`，文档按同一目录内 `sort ASC, id ASC`；只返回承载公开文档的目录及其必要祖先。目录与文档节点总数上限 2,000。

公共文档详情字段固定为 `id`、`workspace_id`、`catalog_id`、`title`、`content_html`、`view_count`、`published_at`、`updated_at`。不返回 `content`、`owner_id`、`article_id`、版本或分享链接。

`content_html` 在公共 Wiki service 返回前使用服务端 HTML 白名单净化：允许常用 Markdown 结构、代码块、表格和安全图片属性；移除 `script/style/iframe/object`、所有 `on*` 事件属性及 `javascript:`/`data:text/html` 等危险协议。净化不修改数据库原文，也不改变现有 token 分享接口语义。

#### 摄影（复用并加固）

- `/photos` 继续调用现有 `GET /api/works?type=photography`。
- 修改公共 `GET /api/works` 的服务端约束：无论客户端是否传 `status`，公共查询只返回 `status=1`。
- 作者管理 `GET /api/works/my` 和管理后台接口继续按原逻辑支持状态筛选。

### 缓存与优先级

偏好最终优先级：

1. 管理员 `holiday` / `mourning` 特殊覆盖，仅影响最终呈现。
2. 登录用户服务端偏好。
3. 当前浏览器本地偏好。
4. 默认 `magazine + system`。

本地缓存使用两个命名空间：访客键 `inkspace_guest_appearance_v1`，账号键 `inkspace_user_appearance_v1:<userId>`；内容只允许白名单字段。登录只更新对应账号键，不覆盖访客键；登出恢复访客偏好。再次登录时先用账号缓存防闪烁，再由服务端值校正。工作区公开状态更新必须沿用现有缓存失效逻辑，公共读取不得因旧缓存绕过 `is_public` 校验。

### 数据库迁移与部署

- `UserAppearance` 和 `Workspace.IsPublic` 均为增量 AutoMigrate；`is_public` 使用 `NOT NULL DEFAULT FALSE`，迁移后执行查询确认历史记录均为 false。
- 部署顺序必须先由单实例完成数据库迁移并验证列、默认值和唯一索引，再放量包含公共 Wiki 路由的新服务实例，避免混合版本查询缺失列。
- 上线前记录 `workspaces` 行数；上线后校验 `is_public=true` 为 0（除非管理员在发布窗口明确设置），并验证唯一索引不会因历史数据失败。
- 本期无破坏性列变更，不需要回填主题偏好记录；回滚应用版本时保留新增表和列，旧版本会忽略它们。

### 实现步骤（每步可独立 commit）

1. [x] **主题偏好模型与 API**：新增 `UserAppearance`、AutoMigrate、service、handler、受保护路由及白名单测试。
2. [x] **主题运行时基础设施**：新增主题注册表、appearance store、启动防闪烁、系统模式监听、特殊覆盖兼容和本地缓存。
3. [x] **外观与主题页**：新增 `/dashboard/appearance`、用户中心入口、四套设计关联卡片、预览/取消/确认流程；管理后台仅扩展首页静态 Hero 文案配置。
4. [x] **共享页面骨架**：统一公共/用户中心布局、明暗开关和内容展示令牌；保持现有请求、字段和交互逻辑。
5. [x] **屿刊公共页面**：完成首页、文章、作品、摄影、知识库、站点页、登录和公开个人主页。
6. [x] **屿刊用户中心**：完成布局、内容管理、编辑器、通知、收藏、资料设置和移动端适配。
7. [x] **摄影入口与安全加固**：新增 `/photos` 页面/导航，复用 photography 作品；公共作品接口强制已发布。
8. [x] **公共知识库模型与 API**：增加 `Workspace.IsPublic`、公开 DTO/查询、草稿回退规则、净化和基础权限测试。
9. [x] **公共知识库前端**：新增 `/wiki`、工作区树和文档详情；管理界面增加公开开关、风险说明、文档状态及发布/取消发布/重新发布操作。
10. [x] **屿刊全矩阵验收**：桌面/平板/手机、浅色/深色、访客/用户、空态/错误态、节日/哀悼、无障碍和构建验证；2026-07-22 用户验收通过。
11. [ ] **B · yu.log**：实现、自动化验证、注册表与后端白名单已完成并开放；三视口全路由人工视觉矩阵待验收。
12. [ ] **C · 小屿的角落**：关联 C 设计规范完成全矩阵适配，验收后开放白名单。
13. [ ] **D · CHEN YU®**：关联 D 设计规范完成全矩阵适配，验收后开放白名单。

### 参考的现有模式

- `docs/design/README.md` — 四主题共用内容模型、明暗、响应式、无障碍和组件边界。
- `docs/design/style-a-magazine.md`、`style-a-magazine.html` — MVP 屿刊的令牌、排版、签名元素和动效基准。
- `docs/design/style-b-terminal.md`、`style-b-terminal.html` — 第二阶段终端主题基准。
- `docs/design/style-c-cozy.md`、`style-c-cozy.html` — 第三阶段手作主题基准。
- `docs/design/style-d-swiss.md`、`style-d-swiss.html` — 第四阶段瑞士主题基准。
- `web/blog/src/stores/appearance.js` — 当前主题运行时、缓存、预览、回滚、系统明暗与特殊覆盖的唯一状态源。
- `web/blog/src/utils/codeTheme.js` — 保留现有代码主题设置契约。
- `web/blog/src/layouts/MainLayout.vue`、`UserCenterLayout.vue` — 公共和用户中心布局入口。
- `web/blog/src/views/Home.vue`、`Blog.vue`、`Works.vue` — 复用现有内容请求与分页筛选逻辑。
- `web/blog/src/views/BlogDetail.vue`、`WorkDetail.vue` — 复用详情、Vditor 渲染、评论和互动逻辑。
- `web/blog/src/stores/user.js`、`utils/api.js` — 登录态和统一响应拦截器；不得改变 `response.data` 的现有含义。
- `internal/models/setting.go`、`internal/service/setting_service.go` — 保留管理员全站特殊主题设置。
- `internal/models/workspace.go`、`internal/service/workspace_service.go` — 工作区 owner 鉴权、缓存和兼容扩展点。
- `internal/models/doc.go`、`internal/service/doc_service.go` — 文档发布状态与 `ContentHTML`。
- `internal/service/catalog_service.go` — 复用目录树构建逻辑，公开查询先过滤再组树。
- `docs/specs/workspace-knowledge-base.md` — 现有私有知识库、发布、版本和 token 分享语义；本 Spec 仅增补公开发现能力。

## 测试计划

### 后端

- [x] `UserAppearance` 无记录时返回 `magazine + system`，PUT 后再次 GET 一致。
- [x] 主题和明暗白名单表驱动测试覆盖合法值、未知值、未开放主题和空值。
- [ ] 用户只能读写自己的偏好；并发 upsert 唯一性已通过 MySQL 集成测试，跨用户 handler 鉴权仍待接口级测试。
- [x] 新建工作区未传 `is_public` 时保持 false；历史库上线后仍需执行部署校验 SQL。
- [ ] 非 owner 无法修改工作区公开状态，统一返回 404。
- [x] 公共工作区列表只返回 `is_public=true`，公开计数不包含草稿。
- [x] 公共目录树只包含已发布文档、其目录和必要祖先；不泄露纯草稿目录名。
- [x] 公共文档详情必须同时满足工作区公开与文档已发布，并且不返回 Markdown 原文。
- [ ] 已发布文档手动保存、自动保存、回滚后原子转草稿，公共详情随即 404；重新发布后恢复。
- [ ] 文档取消发布后公共详情 404；重新发布保持首次 `published_at`，并更新 `updated_at`。
- [x] 公共 HTML 净化移除脚本、事件属性、危险 URL，同时保留常用 Markdown、代码块和安全图片。
- [ ] 公开目录/文档移动、目录重命名和工作区关闭立即反映到公共树；关闭后全部公共接口不可见。
- [ ] token 分享不受工作区公开开关影响，保持现有行为。
- [x] 公共 `/api/works` 即使传 `status=0/2/3` 也固定使用发布状态；`/works/my` 代码路径保持独立，数据库结果级回归测试后续补充。

### 前端

- [x] 增加 Vitest + jsdom 和 `pnpm test` 脚本；主题运行时共 7 个测试通过。
- [x] `system` 模式响应 `matchMedia` 变化，显式 light/dark 不响应。
- [x] 外观预览不发保存请求，取消恢复；保存失败回滚并提示，账号切换中的响应不会污染新会话缓存。
- [x] 访客和账号缓存相互隔离；登录后使用账号偏好，登出恢复此前访客偏好。
- [x] 带有旧 `site_theme` 缓存、`body.theme-*` 和内联变量的浏览器升级后只应用新运行时状态。
- [x] `magazine` 覆盖页面适配矩阵，既有 API 请求参数和响应字段读取不变。
- [x] `/photos` 固定 photography，空态和详情跳转正确。
- [x] `/wiki` 对空工作区、关闭工作区、草稿文档和 404 有明确状态。
- [x] `900px` 与 `560px` 两个断点下导航、列表、正文、表单、目录树、弹窗均可操作。
- [x] 键盘遍历、`:focus-visible`、图片 alt、颜色对比度和 reduced-motion 通过 Phase 1 用户验收。
- [x] 浅/深两模式下 Element Plus、Vditor、代码块、通知下拉和确认弹窗无已知旧主题硬编码残留。

### 验证命令

- [x] `go test ./... -count=1`
- [x] `go vet ./...`
- [x] 在 `web/blog` 执行 `pnpm build`
- [x] 在 `web/blog` 执行 `pnpm test --run`（1 个测试文件、7 个测试）
- [x] 在 `web/blog` 执行 `pnpm lint`
- [x] Phase 1 完成人工视觉与功能回归并通过用户验收；本阶段未要求保留完整截图矩阵。

后端数据库集成测试使用独立 MySQL 测试库（可由 Docker Compose 启动），每个测试事务回滚或清理自身 fixture，严禁连接开发/生产库；Redis 使用独立 DB 编号并在用例后清理命名空间。纯白名单、树构建、净化和主题解析逻辑保持为无数据库单元测试。

已执行条件式 MySQL 集成测试：`internal/service/theme_system_integration_test.go`。运行必须同时提供 `INKSPACE_TEST_MYSQL_DSN`、`INKSPACE_TEST_ALLOW_DROP=1`，且数据库名必须以 `_test` 结尾；验证完成后临时容器已删除。

## 验收标准

1. 新用户、无缓存访客和未知偏好均默认进入「屿刊 · 极简杂志风」。
2. 登录用户的主题与明暗偏好可跨设备同步；预览、取消、确认和失败回滚符合核心流程。
3. 管理员节日/哀悼模式可临时覆盖呈现，但不破坏个人偏好。
4. `magazine` 完成页面适配矩阵，桌面和移动端均可正常完成原有业务操作。
5. 既有文章、作品、评论、用户资料等 API 路径、字段与逻辑无破坏性变更。
6. `/photos` 只展示已发布摄影作品，非发布作品不能通过公共状态参数读取。
7. `/wiki` 只展示公开工作区的已发布文档，草稿、未发布修改和纯草稿目录不泄露。
8. B/C/D 可通过同一注册表、设计令牌和共享页面契约继续适配，不需要复制业务视图或修改业务 API。
9. 屿刊视觉与对应 Markdown/HTML 设计稿在色彩、排版、留白、签名元素、动效和响应式行为上保持一致。

## 待定事项

- 已知兼容限制：`web/admin` 仍显示原有 `day/night` 全站主题文案，但新用户侧只消费 `holiday/mourning` 特殊覆盖；本期遵循用户确认不修改管理端，后续可独立清理旧设置体验。
- 待部署校验：生产 AutoMigrate 后确认历史 `workspaces.is_public` 均为 false。
- 后端专项测试缺口仍以“测试计划 / 后端”中的未勾选项为准，不因视觉验收自动视为完成。

## Phase 1 范围（已完成）

**本期纳入：**

1. 账号级 `ui_theme` / `color_scheme` 偏好模型与专用 API。
2. 可扩展主题注册表、运行时、缓存、系统明暗监听和管理员特殊覆盖兼容。
3. 用户中心“外观与主题”页，含四套设计关联、屿刊预览/取消/确认；未完成主题显示“即将推出”。
4. 「屿刊」完整用户侧页面适配矩阵，而非仅首页换肤。
5. `/photos` 摄影入口，复用现有作品模型和详情，并修复公共作品状态过滤边界。
6. 工作区公开状态、公共 `/wiki` 只读 API 和页面；只公开已发布文档。
7. 已发布知识库文档再次编辑后自动转草稿，防止未发布修改泄露。
8. 浅色/深色、响应式、无障碍、构建和关键权限测试。

**后续阶段：**

1. B · `yu.log` 按 [`terminal-ui-theme.md`](./terminal-ui-theme.md) 完整适配并开放选择。
2. C ·「小屿的角落」完整适配并开放选择。
3. D · `CHEN YU®` 完整适配并开放选择。
4. 公共知识库全文搜索、双向链接、目录独立公开策略或发布快照，如后续另行立项。
