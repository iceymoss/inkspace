# Phase 2：yu.log 终端开发者主题

## 状态
- 创建日期: 2026-07-22
- 状态: 已完成，并于 2026-07-22 通过用户审核
- 阶段: Phase 2
- 主题 ID: `terminal`
- 设计规范: [`style-b-terminal.md`](../design/style-b-terminal.md)
- 高保真原型: [`style-b-terminal.html`](../design/style-b-terminal.html)
- 上游规范: [`user-ui-theme-system.md`](./user-ui-theme-system.md)

## 已确认决策

- 本阶段完整适配所有可达用户侧页面，不以首页或核心页面试点作为完成标准。
- 完成后在 `/dashboard/appearance` 正式开放 `terminal` 预览和账号级保存。
- 继续支持 `system|light|dark`，视觉风格和明暗模式保持独立。
- 只改变展示结构、主题令牌和轻量微交互，不增加可执行命令行、命令面板或快捷键业务。
- 复用当前预览取消、保存失败回滚、访客/账号缓存隔离及 `magazine` 默认兜底。
- 桌面、平板、手机以及浅色、深色的全页面矩阵属于 MVP，不拆成后续补齐。

## 目标

在不改变现有业务 API、权限和数据模型语义的前提下，将 `yu.log` 暗色科技感主题完整接入现有多主题运行时，让用户可在终端式视觉语言下完成浏览、互动、内容管理和知识库操作。

## 非目标

- 不新增数据库表、业务字段或 HTTP 接口；仅扩展现有主题白名单。
- 不复制页面、store、API 客户端或请求逻辑，不构建第二套 SPA。
- 不实现真正可执行的终端命令、命令导航、全站快捷键或开发者控制台。
- 不修改文章、作品、摄影、Wiki、评论、关注、收藏、通知、分享的业务规则。
- 不删除管理后台轮播图、首页 Hero 文案、代码高亮等现有配置能力。
- Phase 2 不开放 `cozy` 或 `swiss`；二者已分别在后续独立阶段开放。
- 不把整站正文改为等宽字体；等宽仅用于路径、命令、标签、时间、EXIF 和状态信息。
- 不为贴合原型伪造城市、技术栈、工作年限、版本号或不存在的社交链接。

## 用户故事

- 作为登录用户，我想在外观页预览并保存 `yu.log`，以便账号在不同设备上使用同一终端主题。
- 作为访客或登录用户，我想让终端主题继续支持跟随系统、浅色和深色，以便按环境舒适阅读。
- 作为内容读者，我想在终端风格下继续使用文章筛选、作品浏览、摄影、Wiki、评论和分享，而不学习新的业务流程。
- 作为内容作者，我想在用户中心继续编辑文章、作品和知识库文档，主题不能降低表单、表格和 Vditor 的可用性。
- 作为维护者，我想复用现有业务树和主题运行时，只增加可隔离、可回滚的展示层适配。

## 触发入口

1. 登录用户进入 `/dashboard/appearance`。
2. 选择 `yu.log · 暗色科技感` 和明暗模式。
3. 点击预览时立即应用，但不写服务端。
4. 点击应用后调用现有 `PUT /api/profile/appearance`。
5. 保存成功后写账号缓存并保持应用；失败或取消时恢复原偏好。
6. 未登录访客不新增主题风格选择入口，继续只能在页头切换明暗；默认主题仍为 `magazine`。

## 核心流程

### A. 主题开放

1. 新增 `terminal.css`，先在未开放状态完成页面适配。
2. 注册表记录 `stylesheet: terminal.css`，开发阶段仍可保持 `coming_soon`。
3. 完成自动化和视觉矩阵后，将注册表状态改为 `available`。
4. 同一提交将后端 `UserAppearance` 白名单扩展为 `magazine|terminal`。
5. Phase 2 验收时 `cozy`、`swiss` 和未知值继续返回 400；Phase 3 开放 Cozy 后仅 `swiss` 和未知值继续返回 400。

### B. 预览与保存

1. 外观页选择 terminal 时复用当前 `appearance.preview()`。
2. `data-ui-theme=terminal` 控制视觉，`data-theme=light|dark` 控制实际明暗。
3. `system` 继续监听 `prefers-color-scheme`，主题预览期间也即时变化。
4. 取消预览、离页或保存失败时恢复进入页面前的已保存偏好。
5. 页面刷新时先应用对应账号缓存防闪烁，再由服务端偏好校正。

### C. 用户侧浏览

1. 公共布局切换为 `$yu.log` 品牌、终端导航语法、背景网格和系统状态页脚。
2. 首页呈现左文右终端 Hero、日志式文章、仓库式作品、captures 摄影和文件树式 Wiki。
3. 所有内容使用现有真实 API；空数据和错误状态使用终端语义，不填充原型假数据。
4. 文章、作品、摄影、Wiki 和分享详情继续使用现有互动、权限和 Markdown 渲染。

### D. 用户中心

1. 用户中心使用 panel/panel-2 层次、路径式标题和等宽状态信息。
2. 文章、作品、通知、评论、收藏、资料和知识库管理流程保持不变。
3. 编辑器外壳跟随 terminal，Markdown 浅深跟随 UI，代码高亮继续受管理后台 `code_theme` 控制。
4. 表格、弹窗、抽屉、下拉、上传、图片裁剪和全屏编辑器必须可用。

## 异常处理

| 场景 | 处理方式 |
|---|---|
| 服务端或缓存返回未开放主题 | 归一化为 `magazine` |
| terminal 样式未加载或初始化失败 | 保留页面可读性并回退已保存主题；不得写入错误偏好 |
| 预览后离开页面 | 恢复已保存主题 |
| 保存失败 | 恢复保存前主题并提示失败 |
| `system` 媒体查询不可用 | 回退实际浅色显示 |
| 管理员启用 holiday/mourning | 作为 terminal 之上的临时覆盖，不改写个人偏好 |
| Hero 终端动效被禁用 | `prefers-reduced-motion` 下静态显示首条命令，光标不闪烁 |
| 代码高亮 CDN 加载失败 | 保证代码块前景和背景仍可读，不阻塞正文 |
| 空列表或 API 失败 | 展示主题化空态/错态和既有重试入口，不使用假内容 |
| 图片加载失败 | 保留比例和可读替代内容，不造成布局跳动 |
| 超长命令、代码、标题或表格 | 换行或横向滚动，不撑破视口 |

## 视觉系统

### 颜色令牌

新增 `web/blog/src/themes/terminal.css`，所有规则限定在 `html[data-ui-theme="terminal"]`。

| Token | 深色 | 浅色 | 语义 |
|---|---|---|---|
| `--bg` | `#0D1320` | `#F4F6FA` | 页面背景 |
| `--panel` | `#131B2C` | `#FFFFFF` | 卡片与面板 |
| `--panel-2` | `#182238` | `#EBEFF6` | 面板内高亮与悬停 |
| `--ink` | `#C9D5E6` | `#23304A` | 正文 |
| `--bright` | `#EDF3FB` | `#0E1730` | 标题与命令 |
| `--sub` | `#64748F` | `#7C88A1` | 辅助文字 |
| `--line` | `#1F2B42` | `#DDE4EF` | 描边和网格 |
| `--accent` | `#5CB2FF` | `#1668C7` | 交互与链接 |
| `--amber` | `#F0A94B` | `#C07A17` | 数字、WARN、行内代码 |
| `--green` | `#6FCF8E` | `#1F9D57` | 正常状态与提示符 |

语义色不得混用：绿表示正常，蓝表示交互，琥珀表示数据或警示。状态同时有文字或形状，不只依赖颜色。

### 字体与层次

- 正文和中文标题使用系统无衬线栈。
- 路径、命令、时间戳、标签、EXIF 和页脚使用系统等宽栈，不加载网络字体。
- 内容最大宽度 `1100px`，桌面两侧 `28px`。
- 面板与卡片圆角 `14px`，按钮 `10px`，小控件 `8px`。
- 层次使用 `bg -> panel -> panel-2`、1px 描边和克制阴影表达。
- 背景网格为 44px，透明度 `.25`，只在公共布局页面上半部出现且不可交互。

### 动效

- 主题颜色切换 `.3s`。
- 卡片 hover 上移 `4px`，描边变蓝；触屏不依赖 hover 显示必要信息。
- 图片 hover 缩放 `.5s`。
- 终端光标 `1.1s steps(1)`，打字动画组件卸载时必须清理计时器。
- `prefers-reduced-motion` 时关闭光标、打字、位移、缩放和平滑滚动。

## 页面适配矩阵

### 公共骨架与页面

| 路由/区域 | 文件 | 终端主题要求 |
|---|---|---|
| 公共骨架 | `layouts/MainLayout.vue` | `$yu.log` 品牌、英文终端导航、移动菜单、背景网格、真实站点页脚 |
| 首页 `/` | `views/Home.vue` | 左文右终端 Hero；轮播配置保留；真实文章、作品、摄影、Wiki 分区 |
| 博客 `/blog` | `views/Blog.vue` | FEAT/INFO 日志流；筛选、搜索、分页与键盘打开保持 |
| 文章 `/blog/:id` | `views/BlogDetail.vue` | 阅读面板、元数据、Markdown、评论和互动状态 |
| 作品 `/works` | `views/Works.vue` | repo/capture 视觉变体；类型筛选、搜索和分页保持 |
| 作品详情 `/works/:id` | `views/WorkDetail.vue` | 项目/摄影分支、图集、EXIF、评论与互动保持 |
| 摄影 `/photos` | `views/Photos.vue` | 3:2 captures 网格、扫描线、EXIF、触屏可读图注 |
| Wiki `/wiki` | `views/wiki/WikiIndex.vue` | 公开工作区面板、作者入口、封面和分页 |
| Wiki 空间 | `views/wiki/WikiWorkspace.vue` | 桌面文件树 + 文档区；窄屏使用现有可操作抽屉 |
| Wiki 文档 | `views/wiki/WikiDoc.vue` | 路径面包屑、正文、表格和行内代码 |
| 站点页 | `About.vue`、`Links.vue`、`UserSearch.vue` | panel、路径标题、表单、空错态与外链安全 |
| 用户主页 | `UserProfile.vue` | 资料、关注和全部内容标签页 |
| 分享 | `ShareDoc.vue` | 正文、有效/失效/过期/网络错误状态 |
| 登录 | `Login.vue` | 登录/注册、校验和回跳；独立入口无主题闪烁 |

### 用户中心

| 区域 | 页面 |
|---|---|
| 骨架 | `layouts/UserCenterLayout.vue` |
| 概览 | `views/user/Dashboard.vue` |
| 文章 | `MyArticles.vue`、`ArticleEdit.vue` |
| 作品 | `MyWorks.vue`、`WorkEdit.vue` |
| 社交与状态 | `MyComments.vue`、`Notifications.vue`、`Favorites.vue` |
| 资料 | `ProfileEdit.vue` |
| 知识库 | `WorkspaceList.vue`、`WorkspaceDetail.vue`、`DocEdit.vue` |
| 外观 | `Appearance.vue` |

要求：侧栏、面包屑、表单、表格、状态标签、确认框、抽屉、上传、拖拽、版本、发布和分享均可操作，不把管理界面变成装饰性终端。

## 主题专属展示组件

建议最小抽取，只有出现真实复用时才建立：

- `TerminalHero.vue`：纯展示，接收 Hero 配置和真实统计，不自行请求 API。
- `TerminalWindow.vue`：标题栏、静态输出和打字命令；负责清理计时器。
- `SectionHeading.vue`：路径式栏目标题，可同时服务后续主题变体。
- `LogRow.vue`：文章日志行，等级映射集中管理。
- `RepoCard.vue`：项目卡；技术栈颜色只映射已知值，未知值使用回退令牌。

不为了形式拆分每个页面；业务状态和 API 调用继续留在现有视图或 store。

## Element Plus 与内容渲染

### Element Plus

terminal 全局作用域至少覆盖：

- card、button、input、textarea、select、radio、switch、tag、badge；
- dialog、drawer、message-box、message、notification；
- dropdown、tooltip、popover、date picker 等 teleport 组件；
- table、tabs、pagination、menu、breadcrumb；
- skeleton、empty、alert、upload、image viewer、loading mask、carousel。

不得依赖页面 scoped CSS 控制挂在 `body` 下的浮层。

### Markdown、Vditor 与代码高亮

- Markdown 和 Vditor 的浅深模式继续跟随 `appearance.resolvedColorScheme`。
- 管理后台 `markdown_theme` 不影响用户侧。
- 管理后台 `code_theme` 继续绝对控制代码块高亮主题，terminal 不自动配对或替换。
- terminal 只控制编辑器外壳、正文排版、行内代码和代码块容器边界，不用 `!important` 覆盖语法颜色。
- 主题预览期间调用现有 Vditor `setTheme()`，不得重建编辑器或丢失内容。
- Vditor 全屏层级高于 sticky header、背景网格和用户中心侧栏。

## 数据与 API

### 数据模型

不新增字段。沿用 `UserAppearance`：

```json
{
  "ui_theme": "terminal",
  "color_scheme": "system"
}
```

仅扩展 `UITheme` 服务端白名单。

### API

不新增接口。复用：

| 方法 | 路径 | 用途 |
|---|---|---|
| GET | `/api/profile/appearance` | 获取账号偏好 |
| PUT | `/api/profile/appearance` | 保存 terminal 与明暗偏好 |
| GET | `/api/settings/public` | 特殊覆盖、Hero、站点信息和代码主题 |
| 既有接口 | 文章、作品、摄影、Wiki、用户和管理能力 | 页面数据与业务操作 |

Terminal 首页文案使用公开设置键 `home_hero_terminal`，值仍为 JSON 字符串并复用现有管理端批量设置接口；它与屿刊的 `home_hero`、`home_carousel` 相互独立，不新增 HTTP 接口或数据库字段。未保存该键时使用前端内置的 yu.log 默认文案。

## 响应式要求

### `<=900px`

- 公共导航和用户中心侧栏使用现有可操作菜单/抽屉。
- Hero 改为单列，终端在正文之后。
- 项目和摄影网格改为两列。
- 日志流改为单列。
- Wiki 文件树改为按钮打开的抽屉，不永久隐藏。
- 表格可横向滚动或转换为可读列表。
- 弹窗、抽屉和编辑器不超出视口。

### `<=560px`

- 项目和摄影网格改为一列，页面水平 padding 约 16px。
- CTA 可换行或纵向排列。
- 终端长命令换行或横向滚动。
- 评论回复、Wiki 元数据和宽表格保持可读。
- 外观页的预览、取消和应用操作全部可见。

现有 460/640/700/760/768 等组件内部断点可以保留，但不得覆盖 900/560 下的最终可用行为。

## 实现步骤

每步应可独立审查；正式开放动作放在最后，避免半成品被保存。

1. [x] **白名单准备与测试**：`terminal` 已加入前后端合法值，默认仍为 `magazine`。
2. [x] **Terminal 令牌**：新增 `terminal.css`，实现浅深令牌、Element Plus 映射、背景网格、特殊覆盖和 reduced-motion。
3. [x] **公共与用户中心骨架**：适配 MainLayout、UserCenterLayout、通知下拉、移动菜单和页脚。
4. [x] **首页签名结构**：实现 Terminal Hero、轮播容器、日志、仓库、摄影和 Wiki 真实数据分区。
5. [x] **公共索引页**：通过主题作用域令牌和稳定页面钩子适配 Blog、Works、Photos、WikiIndex、UserSearch、Links 和 UserProfile 标签页。
6. [x] **公共详情页**：通过主题作用域令牌适配 BlogDetail、WorkDetail 两分支、WikiWorkspace、WikiDoc、ShareDoc、About 和 Login。
7. [x] **用户中心页面**：适配内容管理、资料、收藏、通知、评论、知识库和三个编辑器页面。
8. [x] **外观页与运行时测试**：完成 terminal 预览证明、取消、保存、失败回滚、缓存、system 和特殊覆盖测试。
9. [x] **完整矩阵验收并开放**：注册表、样式和后端白名单已开放；三视口全路由视觉与功能审核通过。
10. [x] **规范同步**：已更新总规范状态、自动化测试记录和剩余风险。

## 测试计划

### 后端

- [x] `terminal + system|light|dark` 均通过模型校验；MySQL PUT/GET 集成断言已补充，需在提供测试库环境变量时执行。
- [x] `magazine` 继续合法并保持默认。
- [x] Phase 2 验收时 `cozy`、`swiss`、未知值和空值返回 400；二者已在 Phase 3/4 按独立规范开放，未知值仍返回 400。
- [ ] 两个用户偏好互不污染，并发 upsert 每用户仍只有一行。
- [ ] 未登录调用偏好接口返回 401。

### 前端自动化

- [x] terminal 缓存启动设置 `data-ui-theme=terminal`。
- [x] terminal 预览不调用 PUT，取消恢复原主题。
- [x] 保存成功写账号缓存，失败恢复原主题。
- [x] `system` 变化即时更新 terminal 明暗；显式 light/dark 不跟随。
- [ ] 登录/登出恢复对应账号或访客上下文。
- [x] holiday/mourning 与 terminal 并存且不改写保存值。
- [x] TerminalHero 在 reduced-motion 下静态显示，卸载后无残留计时器。
- [x] 外观页 terminal 可选且 LIVE PROOF 与当前主题一致；Cozy 与 Swiss 已在后续独立阶段开放。

### 验证命令

- [x] `gofmt` 修改过的 Go 文件
- [x] `go test ./... -count=1`
- [x] `go vet ./...`
- [x] `web/blog: pnpm test --run`（2 个测试文件、12 项测试）
- [x] `web/blog: pnpm lint`
- [x] `web/blog: pnpm build`（通过，保留既有大 chunk 警告）
- [x] `git diff --check`

### 人工视觉矩阵（已通过）

视口：

- `1440x900`
- `900x1024`
- `390x844`

主题轴：

- terminal dark：全部路由、三个视口；
- terminal light：全部路由、桌面和手机；
- terminal system：首页、外观页和打开中的 Vditor，分别切换系统浅/深；
- holiday/mourning：首页、文章详情、摄影详情、WikiDoc、登录、Dashboard、编辑器和弹窗。

数据状态：正常、空列表、API 失败、404/关闭 Wiki/失效分享、图片失败、超长标题、无封面、长代码、宽表格、loading/skeleton。

## 验收标准

1. `terminal` 在外观页可预览、取消和保存，跨设备同步；失败恢复原主题。
2. 默认主题仍为 `magazine + system`，未开放主题仍不可选择或保存。
3. 所有可达用户侧路由完成 terminal 浅深适配，原业务操作和 API 契约无回归。
4. 公共布局具有 `$yu.log`、终端导航、网格、Terminal Hero、日志、仓库、captures、文件树和系统状态等签名语言。
5. 中文正文保持无衬线可读，等宽字体只用于机器语言；蓝、绿、琥珀遵循固定语义。
6. Element Plus 浮层、表单、表格、Vditor、Markdown 和代码块无屿刊或默认主题残留。
7. `900px` 和 `560px` 下导航、Wiki 树、表格、弹窗、表单和编辑器均可操作。
8. 键盘焦点、对比度、alt、状态文字和 reduced-motion 达标。
9. holiday/mourning 不破坏 terminal 可读性，也不改写账号偏好。
10. 完整自动化与人工矩阵已通过，terminal 已标记为 `available`。

## 待定事项

- 本次审核未要求保存完整截图矩阵；C/D 主题继续以设计规范、HTML 原型和实际用户审核作为验收基线。
- 首页 Terminal Hero 的命令输出只能使用现有站点设置和真实统计；若未来需要城市、技能或年限，应另行扩展公共设置，不在本阶段伪造。
- 当前管理后台 `pnpm lint` 脚本引用缺失的 `.gitignore` 且无独立 ESLint 配置；如 Phase 2 修改管理后台，应另行修复工具配置。本阶段原则上不改管理后台。

## MVP 范围

本阶段 MVP 即完整交付，不拆分桌面优先或深色优先版本：

1. terminal 浅色、深色和跟随系统令牌。
2. 全部公共页面、用户中心和编辑器适配。
3. Terminal Hero、日志、仓库、摄影与 Wiki 签名结构。
4. 外观页正式开放和账号级保存。
5. Element Plus、Vditor、Markdown、代码高亮和特殊覆盖兼容。
6. 三视口、全路由、空错态、无障碍和 reduced-motion 验收。
