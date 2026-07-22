# Phase 3：inkspace 温暖手作主题

## 状态

- 创建日期: 2026-07-22
- 状态: 已完成，并于 2026-07-22 通过用户审核
- 阶段: Phase 3
- 主题 ID: `cozy`
- 展示名称: `inkspace · 温暖手作感`
- 默认明暗: `light`
- 设计规范: [`style-c-cozy.md`](../design/style-c-cozy.md)
- 高保真原型: [`style-c-cozy.html`](../design/style-c-cozy.html)
- 上游规范: [`user-ui-theme-system.md`](./user-ui-theme-system.md)

## 已确认决策

- 本阶段完整适配所有可达用户侧页面，不以首页或公共页局部换肤作为完成标准。
- `cozy` 继续使用现有 `system|light|dark` 明暗轴，浅色是默认推荐而不是唯一可用模式。
- 复用同一套路由、API、store、权限、Markdown、Vditor 和 Element Plus，不复制业务视图。
- 首页签名结构为真实摄影组成的胶带拍立得堆；不足三张时按真实数量排版，无摄影时显示明确纸条空态，不生成假照片。
- Wiki 封面复用 `Workspace.Icon`；无图片时使用工作区图标或名称首字的笔记本封面，不新增封面字段。
- 三色系统表达内容类别和视觉节奏，不改变审核状态、发布状态、错误等级等业务语义。
- 手作质感只使用 CSS、SVG 和真实内容图片，不引入纹理位图、网络字体、Canvas 粒子或装饰素材包。
- 装饰旋转限制在 `±7deg`，正文、表单、表格、弹窗和编辑器保持正向，不牺牲管理效率。
- 浮动交互终端跨主题继续工作；在 Cozy 下只映射主题令牌和窗口外壳，不把全站改回终端语言。
- 自动化和用户审核已通过，前端注册表与后端 `cozy` 白名单已同步开放。

## 目标

在不改变 InkSpace 业务契约的前提下，将现有用户侧应用完整呈现为一间可阅读、可创作、可管理的内容小屋。公共页面通过暖纸、拍立得、笔记本和手绘线建立亲切感；用户中心通过更克制的纸张层次、清晰描边和物理按钮保留高密度操作效率。

## 非目标

- 不新增数据库表、业务字段或 Cozy 专用 HTTP 接口。
- 不新增摄影内容模型；拍立得继续使用 `Work.Type=photography` 的已发布作品。
- 不伪造姓名、城市、照片地点、拍摄参数、阅读时长、技术栈、作者身份或工作区更新时间。
- 不把 emoji 作为唯一图标或状态表达；现有 Element Plus 图标和文字状态继续保留。
- 不让随机旋转导致刷新后布局跳动；角度由稳定序号或资源 ID 映射。
- 不对用户上传图片叠加永久滤镜，不改变原图、裁剪结果或下载内容。
- 不改变 Vditor 工具栏、自动保存、版本、发布、分享和全屏行为。
- 不改变管理后台整体视觉；若增加 Cozy Hero 配置，只在现有公开设置表单中增加对应键。
- 不在本阶段开放 `swiss`，未知主题继续回退 `magazine`。
- 不引入 `vue-i18n`、新 UI 框架、CSS 框架或图形动画依赖。

## 用户故事

- 作为登录用户，我想预览并保存“inkspace · 温暖手作感”，并在其他设备登录后恢复同一主题。
- 作为读者，我想在温暖而轻松的页面中阅读长文、浏览作品、摄影和 Wiki，但正文仍然清晰稳定。
- 作为摄影浏览者，我想看到由真实作品组成的拍立得墙，并能正常进入作品详情和图集。
- 作为知识库读者，我想把公开工作区理解为一本本笔记本，同时继续使用目录树和文档导航。
- 作为作者，我想在手作主题下完成文章、作品和知识库管理，不因装饰效果降低表格、表单或编辑器效率。
- 作为低动态偏好用户，我想关闭旋转、弹性和缩放动效，但仍保留完整信息层次。
- 作为维护者，我想让 Cozy 只增加隔离的令牌和必要展示变体，并沿用已验证的主题预览、回滚与缓存机制。

## 触发入口

1. 登录用户进入 `/dashboard/appearance`。
2. 选择“inkspace · 温暖手作感”和 `system|light|dark`。
3. 点击“预览效果”后调用现有 `appearance.preview()`，立即设置 `data-ui-theme=cozy`，不写服务端。
4. 点击“应用主题”后调用现有 `PUT /api/profile/appearance`。
5. 保存成功后写账号缓存；失败、取消或未保存离页时恢复进入页面前的偏好。
6. 未登录访客不新增风格选择入口，仍只切换明暗；默认 UI 主题保持 `magazine`。

## 核心流程

### A. 主题开放

1. 新增 `web/blog/src/themes/cozy.css`，所有规则限定在 `html[data-ui-theme="cozy"]`。
2. 注册表在开发期间保持 `coming_soon`；当前实现与自动化验证完成后已切换为 `available` 并加载 `cozy.css`。
3. 用户审核通过后，前端状态已改为 `available`，后端白名单已同步扩展为 `magazine|terminal|cozy`。
4. `swiss` 和未知值继续被前端归一化、后端拒绝，不增加兼容分支。
5. 默认偏好仍是 `magazine + system`，开放 Cozy 不迁移或改写现有用户记录。

### B. 首页浏览

1. 公共页头使用不规则“屿”字标、中文口语导航和明暗切换；移动端使用真实菜单，不隐藏入口。
2. Hero 左侧使用可配置问候、标题、简介和双 CTA；右侧从已发布摄影中取最多三项组成拍立得堆。
3. 拍立得图片、标题和链接均来自真实作品；地点只在 `metadata.location` 存在时显示。
4. 首页文章使用随笔卡，作品使用手作卡，摄影使用照片墙，公开 Wiki 使用笔记本卡。
5. 首页各区块继续复用现有 API；空数据、失败和图片错误使用主题化纸条状态，不填充设计原型示例。
6. Hero 与首个内容区之间使用装饰性 SVG 波浪线，`aria-hidden=true`，不得增加键盘停靠点。

### C. 公共内容与互动

1. 博客列表保留搜索、分类、标签、排行、排序、分页和 URL query，同一数据仅改变卡片结构。
2. 文章详情以宽纸页呈现，目录、代码块、表格、图片、评论、点赞、收藏和上一篇/下一篇保持原行为。
3. 作品列表继续区分项目与摄影；项目使用手作卡，摄影使用拍立得卡，状态和筛选不由颜色猜测。
4. 摄影详情保留图集、EXIF、互动和评论，图片不因旋转裁掉必要内容，查看器保持正向全屏。
5. Wiki 列表使用笔记本封面；工作区使用书架式目录面板；文档详情使用纸页正文，但目录层级和公共权限不变。
6. 用户主页、About、Links、搜索、登录和分享页使用同一纸张、描边、按钮、空错态语言。

### D. 用户中心与创作

1. 用户中心侧栏使用克制的笔记本索引视觉，不旋转导航文字，不改变菜单顺序。
2. Dashboard、文章、作品、评论、收藏、通知、资料和外观页保持当前信息密度和操作入口。
3. 列表、表格和分页不模拟散落纸片；仅使用 `card -> card-soft` 层次、2px 描边和小范围暖阴影。
4. 创建/编辑表单、上传、裁剪、确认框、抽屉、下拉、日期控件和 loading 状态必须完整适配。
5. Workspace 管理继续支持公开开关、目录拖拽、文档状态、发布、版本和分享；书脊色不代表发布状态。
6. Vditor 外壳跟随 Cozy，Markdown 明暗跟随当前实际模式，代码高亮继续由管理后台 `code_theme` 决定。

### E. 跨主题终端

1. 已启动的 FloatingTerminal 在切换 Cozy 后继续保留窗口、历史、cwd 和 pending 状态。
2. 窗口映射 `--card|--card-soft|--line|--ink|--moss|--mustard`，保持文本对比度和错误语义。
3. 终端命令、虚拟文件系统和确认流程不改变；不得加入 Cozy 专属命令或把输出改成装饰文案。
4. 浮窗、移动端抽屉和焦点环继续遵循交互终端规范，层级高于 sticky header、低于模态框与 Vditor 全屏。

## 内容映射规则

### 首页拍立得

- 数据源：`GET /api/works?type=photography&status=1&page=1&page_size=3`。
- 图片优先级：`cover` -> `images[0].url|string`；均缺失时显示标题首字和“暂无照片”纸面占位。
- 图注：只显示真实 `title`；真实 `metadata.location` 可作为第二行，不从描述中猜地点。
- 角度：按卡片索引稳定使用 `-6deg|3deg|-2deg`，不足三项时按实际数量居中。
- 图片失败后切换同尺寸占位，不移除卡片，不造成 Hero 抖动。

### 分类三色

- `moss`：设计、产品、体验、UI/UX 等已知分类或标签。
- `sky`：代码、开发、前端、后端、工程等已知分类或标签。
- `blush`：摄影、生活、旅行、随笔等已知分类或标签。
- 无已知语义时使用资源 ID 的稳定取模映射，不使用 `Math.random()`，同一内容跨首页、列表和刷新保持一致。
- 颜色只用于标签、图标底和书脊；审核、发布、错误和危险操作继续使用现有状态文字及语义色。

### Wiki 笔记本

- 封面图片来自 `workspace.icon`，只接受现有安全 URL 规则。
- 无图片时显示 `workspace.icon` 文本值或名称首字；描述、文档数、作者、更新时间均读取真实公共 DTO。
- 书脊色按稳定序号或 ID 在 moss/sky/blush 间轮换，不表达分类或权限。
- 工作区关闭、无已发布文档和 API 错误沿用现有不可见与错误处理，不保留旧封面缓存伪装可访问状态。

## 异常处理

| 场景 | 处理方式 |
|---|---|
| 服务端或缓存返回未开放主题 | 归一化为 `magazine` |
| Cozy 样式加载失败 | 保持基础样式可读并恢复已保存主题，不写错误偏好 |
| 预览取消、离页或保存失败 | 复用 appearance store 回滚到进入页面前偏好 |
| `system` 媒体查询不可用 | 实际模式回退 `light` |
| 首页摄影不足三张 | 按真实数量重排，不复制图片 |
| 首页没有摄影 | 显示可操作空态和 `/photos` 入口，不生成渐变假照片 |
| 图片或 Wiki 封面失败 | 保留比例并显示文字占位与 alt，不造成布局跳动 |
| 标题、描述或标签过长 | 两行截断或安全换行；详情页展示完整内容 |
| 分类无法识别 | 使用稳定 ID 映射的回退色，不新增后端分类 |
| holiday/mourning 生效 | 作为 Cozy 之上的临时覆盖，不改写账号偏好 |
| reduced-motion | 取消旋转变化、过冲、位移、缩放和主题按钮转动 |
| 触摸设备无 hover | 标题、图注和操作始终可见，不依赖回正动画 |
| 宽表格、长代码、长 URL | 容器横向滚动或安全换行，不旋转、不撑破纸页 |
| Teleport 浮层 | 使用 Cozy 全局作用域变量，不能残留 Element Plus 默认蓝色 |

## 视觉系统

### 颜色令牌

新增 `web/blog/src/themes/cozy.css`。设计稿令牌通过兼容别名映射到现有页面 `--theme-*`，避免业务组件写主题分支。

| Token | 浅色 | 深色 | 语义 |
|---|---|---|---|
| `--bg` | `#FBF6ED` | `#26201A` | 页面暖纸/暖夜 |
| `--card` | `#FFFDF8` | `#2E2721` | 主要纸张和卡片 |
| `--card-soft` | `#F7EFE2` | `#352D25` | 内层、悬停、表头 |
| `--ink` | `#43362A` | `#EFE6D9` | 正文和标题 |
| `--sub` | `#7D6C5A` | `#B8A58F` | 辅助文字；较设计稿原值提高正文对比度 |
| `--line` | `#E2D4C0` | `#4A3E33` | 2px 主描边与虚线 |
| `--moss` | `#66764E` | `#A3B385` | 主交互、激活和焦点 |
| `--moss-soft` | `#EDEFE2` | `#333A29` | 主色软底 |
| `--mustard` | `#D9A441` | `#E0B25C` | 主 CTA、装饰圆点和胶带 |
| `--mustard-ink` | `#4A3814` | `#31240C` | 芥黄按钮文字 |
| `--sky` | `#8FAABF` | `#93B0C6` | 辅助分类色 |
| `--blush` | `#E5B8A5` | `#D9A891` | 辅助分类色 |
| `--danger` | `#B65345` | `#E08878` | 删除、失败和危险操作 |
| `--shadow` | `0 6px 24px rgb(122 98 70 / 10%)` | `0 6px 24px rgb(0 0 0 / 35%)` | 暖软阴影 |

约束：芥黄背景使用 `--mustard-ink`，不用白字；`--sub` 不承载关键正文；状态不能只依赖 moss/sky/blush 三色。

### 字体与排版

- 正文和 UI：`'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif`。
- 品牌与公共页大标题可优先 `'Yuanti SC'`，缺失时自然回退，不加载网络字体。
- Hero H1：`clamp(32px, 5vw, 48px)`，字重 800，行高 1.4。
- 页面 H1：`clamp(30px, 4vw, 44px)`；栏目 H2 24px；卡片标题 17px；正文 15.5px/1.85。
- Markdown 长文保持阅读宽度约 `760px`，不为手作感降低字号或增加随机字距。
- 元数据最小 12.5px；手机端表单和输入文字不得低于 16px，避免浏览器自动缩放。

### 形状与空间

- 公共内容最大宽度 `1020px`，桌面两侧 28px；用户中心沿用现有可用宽度。
- 公共卡片圆角 20px，用户中心卡片 14px，按钮/输入 12px，小标签使用胶囊。
- 笔记本卡圆角 `8px 20px 20px 8px`，左侧 10-12px 书脊。
- 公共卡片使用 2px 描边；表单、表格和浮层使用 1-2px 清晰边界。
- 区块垂直间距 56px，Hero 桌面 72px/64px；网格 gap 18-20px。
- 旋转只用于品牌块、装饰纸条、拍立得、邮票和公共展示卡；高密度操作区域保持 `transform:none`。

### 签名元素

- **胶带拍立得堆**：Hero 最多三张真实摄影，胶带为半透明 CSS 色块。
- **手绘波浪线**：内联 SVG 使用 `currentColor`，纯装饰并隐藏于辅助技术。
- **果酱按钮**：芥黄胶囊和 4px 实体底边，按下移动 2px；禁用态取消实体按压感。
- **照片墙**：真实图片使用拍立得白边，奇偶稳定角度，hover 回正。
- **笔记本**：Wiki 卡片用书脊、虚线页数徽章和真实封面。
- **邮票页脚**：显示站点真实品牌简称；链接和版权沿用公共设置，不写死原型姓名。

### 动效

- 主题颜色切换 `.35s`。
- 卡片 hover `.25s`：最多上移 5px、旋转 `.5deg`、描边变化。
- 拍立得 `.3s cubic-bezier(.3,1.4,.5,1)`：回正并最多放大 1.05。
- 果酱按钮 `.15s`：按压位移和底边高度互换。
- 明暗按钮 hover 最多旋转 20deg；激活后不持续动画。
- 动效仅由用户输入触发，不增加自动漂浮、无限摇摆或滚动视差。
- `prefers-reduced-motion: reduce` 下禁用 transition/animation 和平滑滚动，并使所有卡片保持正向。

## 页面适配矩阵

### 公共骨架与页面

| 路由/区域 | 文件 | Cozy 要求 |
|---|---|---|
| 公共骨架 | `layouts/MainLayout.vue` | 不规则品牌、中文胶囊导航、移动菜单、纸张页脚、明暗入口 |
| 首页 `/` | `views/Home.vue` | Cozy Hero、真实拍立得堆、随笔、手作、照片墙、Wiki 笔记本 |
| 博客 `/blog` | `views/Blog.vue` | 三色分类卡、全部 query 筛选和分页、空错态 |
| 文章 `/blog/:id` | `views/BlogDetail.vue` | 纸页正文、目录、Markdown、互动、评论、上一篇/下一篇 |
| 作品 `/works` | `views/Works.vue` | 项目手作卡/摄影拍立得变体、筛选搜索分页 |
| 作品详情 `/works/:id` | `views/WorkDetail.vue` | 项目与摄影分支、图集、EXIF、外链、互动评论 |
| 摄影 `/photos` | `views/Photos.vue` | 稳定角度照片墙、图注、触屏可读操作 |
| Wiki `/wiki` | `views/wiki/WikiIndex.vue` | 笔记本卡、真实封面/占位、作者与分页 |
| Wiki 空间 | `views/wiki/WikiWorkspace.vue` | 书架式目录、文档列表、窄屏可操作抽屉 |
| Wiki 文档 | `views/wiki/WikiDoc.vue` | 纸页正文、面包屑、代码表格图片、阅读元数据 |
| 站点页 | `About.vue`、`Links.vue`、`UserSearch.vue` | 纸卡、邮签式标签、表单、空错态与安全外链 |
| 用户主页 | `UserProfile.vue` | 资料卡、关注、文章/作品/收藏标签页完整可用 |
| 分享 | `ShareDoc.vue` | 正文以及有效、失效、过期、网络错误状态 |
| 登录 | `Login.vue` | 登录注册纸卡、验证、回跳与启动无闪烁 |
| 浮动终端 | `components/terminal/FloatingTerminal.vue` | Cozy 令牌外壳，命令、拖拽、抽屉、历史和确认不变 |

### 用户中心

| 区域 | 页面 | Cozy 要求 |
|---|---|---|
| 骨架 | `layouts/UserCenterLayout.vue` | 笔记本索引侧栏、面包屑、通知、移动菜单 |
| 概览 | `views/user/Dashboard.vue` | 数据卡与快捷入口，不用装饰色代替指标含义 |
| 文章 | `MyArticles.vue`、`ArticleEdit.vue` | 筛选、状态、分页、编辑、发布与封面 |
| 作品 | `MyWorks.vue`、`WorkEdit.vue` | 类型、状态、图片、技术栈、外链和提交 |
| 社交与状态 | `MyComments.vue`、`Notifications.vue`、`Favorites.vue` | 列表、未读、删除、跳转和确认 |
| 资料 | `ProfileEdit.vue` | 表单、头像上传裁剪和失败反馈 |
| 知识库 | `WorkspaceList.vue`、`WorkspaceDetail.vue`、`DocEdit.vue` | 封面、公开开关、目录树、拖拽、状态、版本、发布、分享 |
| 外观 | `Appearance.vue` | Cozy 卡可选、Live Proof、预览取消保存和明暗模式 |

## 主题专属展示与复用

遵循最小抽取原则；仅在首页和列表真实复用后建立组件：

- `CozyHero.vue`：接收 Hero 配置、真实摄影和统计，不自行请求 API。
- `PolaroidItem.vue`：接收图片、alt、图注、角度和链接，负责图片失败占位。
- `PolaroidStack.vue`：只负责 0-3 项布局，不复制或补造数据。
- `WavyDivider.vue`：无状态 SVG 装饰组件。
- `NotebookCard.vue`：公开 Wiki 展示，接收封面、书脊色和真实元数据。
- `cozyContentTone.js`：纯函数，集中处理稳定三色映射和角度映射并提供测试。

若当前页面只使用一次且模板简短，保留在视图中，不为了设计稿目录创建空壳组件。业务请求、权限、提交和分页继续属于现有视图或 store。

## Element Plus 与内容渲染

### Element Plus

Cozy 全局主题至少覆盖：

- button、card、input、textarea、select、radio、checkbox、switch、tag、badge；
- dialog、drawer、message-box、message、notification、loading mask；
- dropdown、tooltip、popover、date picker、image viewer 等 Teleport 浮层；
- table、tabs、pagination、menu、breadcrumb、steps；
- skeleton、empty、alert、upload、carousel、avatar 和 form validation。

主按钮可使用果酱按压感；危险按钮必须使用 `--danger`，不能使用芥黄。弹窗、表格、输入和编辑器不旋转。浮层样式必须写在主题全局作用域，不能依赖页面 scoped CSS。

### Markdown、Vditor 与代码高亮

- Markdown 和 Vditor 使用 `appearance.resolvedColorScheme` 选择浅深模式。
- Cozy 控制正文纸页、标题、引用、分隔线、表格边界、行内代码和代码容器，不覆盖语法 token。
- `code_theme` 继续由管理后台设置绝对控制，不根据 Cozy 三色系统替换。
- 引用块可使用 moss 左边线和 `card-soft` 背景；警告类内容仍使用明确语义和文字。
- 图片保持原始比例，可使用 12px 圆角和说明文字，不应用拍立得旋转到正文内容。
- Vditor 预览切换使用现有 `setTheme()`，不得重建编辑器或丢失未保存内容。
- Vditor 全屏层级高于页头、用户中心侧栏、装饰元素和浮动终端。

## 数据与 API

### 数据模型

不新增数据库字段。沿用：

```json
{
  "ui_theme": "cozy",
  "color_scheme": "system"
}
```

仅在正式开放时扩展 `UserAppearanceRequest` 的服务端白名单。`Workspace.Icon` 继续承载 Wiki 封面，`Work` 继续承载项目和摄影数据。

### API

不新增 Cozy 专用接口，复用：

| 方法 | 路径 | 用途 |
|---|---|---|
| GET/PUT | `/api/profile/appearance` | 读取和保存 Cozy 偏好 |
| GET | `/api/settings/public` | Hero、站点信息、特殊覆盖和代码主题 |
| GET | `/api/articles/*` | 首页随笔、博客列表与详情 |
| GET | `/api/works/*` | 手作项目、摄影、Hero 拍立得与详情 |
| GET | `/api/wiki/*` | 公开笔记本、目录与文档 |
| 既有接口 | 用户、评论、互动及用户中心 | 原业务操作 |

Cozy 首页文案使用公开设置键 `home_hero_cozy`，值为 JSON 字符串并复用现有管理端批量设置接口；它与屿刊的 `home_hero`、Terminal 的 `home_hero_terminal` 和轮播配置相互独立。未保存、字段为空或 JSON 无效时使用前端内置默认值：问候语“你好呀，欢迎来坐坐”、标题“把喜欢的事，慢慢做成日常。”、简介“这里收录文章、作品、照片，以及持续整理中的知识。”，按钮为“读读随笔”和“看看照片”。

## 响应式要求

### `<=900px`

- Hero 改为单列，拍立得堆在正文后居中，最大宽度 420px。
- 公共卡片、作品、照片墙和 Wiki 改为两列。
- 导航进入可操作菜单；主题与登录入口保持可达。
- 用户中心侧栏使用现有移动菜单/抽屉，不直接隐藏。
- Wiki 树使用按钮打开抽屉；表格可横向滚动或变为结构化列表。
- 弹窗宽度、图片查看器和编辑器不超出视口。

### `<=560px`

- 页面水平 padding 约 16px，全部内容卡和 Wiki 改为一列。
- Hero 拍立得堆降低高度并按实际数量错位，不能横向溢出或遮挡 CTA。
- 照片墙改为一列或两列窄卡，以图片和图注可读为准。
- 主/次 CTA 可纵向排列并提供至少 44px 触摸高度。
- 装饰旋转幅度减半；表单、列表、评论和编辑器保持正向。
- 外观页预览、取消和应用按钮始终可见；输入字号至少 16px。
- 安全区使用 `env(safe-area-inset-bottom)`，浮动终端抽屉不遮挡页面关键操作。

## 可访问性与性能

- 交互焦点使用 2px moss 虚线和至少 3px offset；危险操作保持 danger 焦点语义。
- 正文对比度不低于 4.5:1，大字号和非文字边界不低于 3:1。
- 拍立得、封面和正文图片都有来自真实标题的 alt；纯胶带、波浪线和邮票装饰正确隐藏。
- 卡片整体可点击时避免内部重复嵌套交互；作者、收藏等独立入口使用合法 DOM 和明确名称。
- 键盘顺序遵循 DOM，不跟随视觉旋转位置；hover 信息必须在 focus 和触屏下可见。
- 图片设置 `loading=lazy`（首屏 Hero 关键图按实际性能决定）、尺寸或 `aspect-ratio`，失败后不抖动。
- 不加载网络字体和纹理图片；CSS 装饰不创建高频 repaint 动画。
- `color-mix()` 需有基础 token 回退，旧浏览器至少保持可读边界和背景。

## 实现步骤

每步可独立审查；正式开放放在最后。

1. [x] **状态与测试基线**：保留 Swiss 不可保存断言，记录当前全量测试和构建基线。
2. [x] **Cozy 令牌**：新增 `cozy.css`，实现浅深令牌、Element Plus 映射、特殊覆盖和 reduced-motion。
3. [x] **稳定视觉映射**：新增分类色、旋转角度和图片回退纯函数及单测，禁止随机布局。
4. [x] **公共与用户中心骨架**：适配 MainLayout、UserCenterLayout、移动菜单、通知、页脚和浮动终端。
5. [x] **首页签名结构**：实现真实摄影拍立得 Hero、波浪线、随笔、手作、照片墙和 Wiki 笔记本。
6. [x] **公共索引页**：通过 Cozy 全局令牌和稳定页面钩子适配 Blog、Works、Photos、WikiIndex、UserSearch、Links 和 UserProfile 标签页。
7. [x] **公共详情页**：通过 Cozy 全局令牌适配 BlogDetail、WorkDetail 两分支、WikiWorkspace、WikiDoc、ShareDoc、About 和 Login。
8. [x] **用户中心与编辑器**：适配全部列表、表单、浮层、上传、知识库和 Vditor 页面。
9. [x] **外观 Live Proof**：为 Cozy 增加独立预览，并验证预览、取消、保存失败回滚和 system 切换基础流程。
10. [x] **自动化与构建**：前后端测试、lint、build、vet 和 diff 检查通过。
11. [x] **完整矩阵审核并开放**：三视口、浅深、全路由、空错态和无障碍完成用户审核，前后端白名单已同步开放。
12. [x] **规范同步**：已更新本规范、总主题规范的 Phase 3 状态、验证记录和剩余风险。

## 测试计划

### 后端

- [x] `cozy + system|light|dark` 均通过模型校验。
- [x] `magazine`、`terminal` 继续合法，默认仍为 `magazine + system`。
- [x] `swiss`、未知值、空主题和未知明暗继续返回 400。
- [ ] Cozy PUT 后 GET 一致，用户之间偏好不污染，未登录请求返回 401。

### 前端自动化

- [x] 注册表为 `available` 且 stylesheet 为 `cozy.css`，Swiss 仍为 `coming_soon`。
- [x] 启动缓存可应用合法 Cozy；未知或未开放状态正确回退。
- [x] Cozy 复用现有预览、取消和保存失败回滚流程。
- [x] `system` 实时跟随浅深变化；显式 light/dark 不跟随。
- [x] 分类色和旋转映射稳定，相同 ID 跨刷新一致。
- [x] Hero 不复制摄影数据，最多展示三张，图片失败有同尺寸占位。
- [x] Wiki 图片封面和无图笔记本占位均可读。
- [x] holiday/mourning 与个人主题并存且不改写保存偏好。
- [x] FloatingTerminal 使用 Cozy 令牌，主题命令和补全支持 Cozy。
- [x] reduced-motion 规则关闭旋转变化、过冲、缩放和按钮转动。
- [x] Appearance 的 Cozy 卡和独立 Live Proof 已开放，Swiss 仍禁用。

### 验证命令

- [x] `gofmt` 修改过的 Go 文件
- [x] `go test ./... -count=1`
- [x] `go vet ./...`
- [x] `web/blog: pnpm test --run`（11 个测试文件、73 项测试）
- [x] `web/blog: pnpm lint`
- [x] `web/blog: pnpm build`（通过，保留既有大 chunk 警告）
- [x] `git diff --check`

### 人工审核矩阵（已通过）

视口：

- `1440x900`
- `900x1024`
- `390x844`

主题轴：

- cozy light：全部路由、三个视口；
- cozy dark：全部路由、桌面和手机；
- cozy system：首页、外观页、打开中的 Vditor 和 FloatingTerminal，切换系统浅深；
- holiday/mourning：首页、文章详情、摄影详情、WikiDoc、登录、Dashboard、编辑器、弹窗和终端。

数据轴：正常、空列表、API 失败、404、失效分享、0/1/2/3 张摄影、图片失败、无 Wiki 封面、超长标题、长代码、宽表格、loading、上传和表单校验。

交互轴：键盘、触屏、hover/focus 等价、菜单、抽屉、弹窗、图片查看器、Vditor 全屏、预览取消保存、终端拖拽与移动抽屉。

## 验收标准

1. `cozy` 可在外观页预览、取消和保存，账号同步、缓存、失败回滚与明暗模式行为正确。
2. 默认仍为 `magazine + system`；`magazine`、`terminal` 无回归，`swiss` 仍不可保存。
3. 所有可达用户侧路由完成 Cozy 浅深适配，原业务流程、权限、URL query 和 API 契约不变。
4. 首页具有真实拍立得 Hero、手绘波浪线、随笔卡、手作架、照片墙和 Wiki 笔记本等签名结构。
5. 不伪造内容：照片、图注、地点、封面、作者、统计和更新时间均来自真实数据或明确空态。
6. 手作装饰受控：旋转稳定且不超过 7deg；正文、表单、表格、弹窗和编辑器保持正向可用。
7. 浅深色正文、辅助文字、按钮和焦点对比度达标，芥黄按钮使用深色文字。
8. Element Plus Teleport 浮层、Markdown、Vditor、代码块、图片查看器和 FloatingTerminal 无旧主题残留。
9. `900px` 和 `560px` 下导航、拍立得、Wiki 树、表格、弹窗、表单和编辑器均可操作且不溢出。
10. 键盘顺序、alt、focus-visible、触摸目标、图片防抖和 reduced-motion 达标。
11. holiday/mourning 可覆盖 Cozy 呈现但不改写偏好，覆盖结束后恢复 Cozy。
12. 自动化、构建和人工审核均已通过，`cozy` 已标记为 `available` 并加入服务端白名单。

## 待定事项

- 本次审核未要求保存完整截图矩阵；Phase 4 继续以设计规范、HTML 原型和实际用户审核作为验收基线。
- Cozy Hero 已接入管理后台首页轮播设置页，可独立配置问候语、标题、强调词、简介和两个按钮。
- 首页拍立得排序当前建议沿用摄影公开列表顺序；若未来需要人工精选，应独立扩展配置，不复用作品排序猜测精选语义。
- 设计稿的口语导航仅改变展示文本，不改变路由；最终文案需在实现审核中确认与全站内容调性一致。

## 最终审核调整

- 品牌统一为 `InkSpace`，手作方形标使用 `Ink`，不把主题名当作项目名称。
- `/blog` 保持左侧紧凑单列文章归档与右侧热门标签汇总；文章项宽度与左侧列表区域一致。
- 个人主页文章和收藏使用等宽模块卡片；关注/粉丝列表收紧垂直间距。
- 用户中心“我的文章”操作列保持单行，不让查看、编辑、状态和删除按钮换行。
- 首页 Hero 默认文案已接入 `home_hero_cozy`，可在管理后台独立配置并保留安全默认值。

## MVP 范围

本阶段 MVP 为完整交付，不拆分“只做浅色”“只做首页”或“先开放再补用户中心”：

1. Cozy 浅色、深色、跟随系统令牌和特殊覆盖。
2. 全部公共页面、用户中心、编辑器、Element Plus 浮层和跨主题终端适配。
3. 真实摄影拍立得 Hero、随笔、手作、照片墙和 Wiki 笔记本签名结构。
4. 稳定三色/旋转映射、真实数据回退、无图片和空错态。
5. 外观页预览、取消、保存、回滚、账号缓存与正式白名单开放。
6. 三视口、全路由、浅深、无障碍、reduced-motion、自动化与用户审核。
