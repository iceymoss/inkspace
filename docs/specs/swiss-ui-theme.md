# Phase 4：Swiss 瑞士网格主题

## 状态

- 创建日期: 2026-07-22
- 状态: 已完成，并于 2026-07-22 通过用户审核
- 阶段: Phase 4
- 主题 ID: `swiss`
- 展示名称: `{site_name} · 瑞士网格风`，缺省站点名为 `InkSpace`
- 默认明暗: `light`
- 设计规范: [`style-d-swiss.md`](../design/style-d-swiss.md)
- 高保真原型: [`style-d-swiss.html`](../design/style-d-swiss.html)
- 上游规范: [`user-ui-theme-system.md`](./user-ui-theme-system.md)

## 已确认决策

- Swiss 顶部品牌、页脚和外观页预览使用公开设置 `site_name`，缺省为 `InkSpace`；不固定使用原型人物名 `CHEN YU®`，不展示未经确认的注册商标符号。
- 本阶段完整适配全部可达用户侧页面，不以首页局部换肤作为完成标准。
- `swiss` 继续使用现有 `system|light|dark` 明暗轴，浅色是默认推荐而不是唯一模式。
- 复用同一套路由、API、store、权限、Markdown、Vditor、Element Plus 和 FloatingTerminal，不复制业务视图。
- 视觉保留可见十二列网格、克莱因蓝、编号系统、零圆角、零阴影和快速反色反馈。
- 原型内容只作为版式参考；人物身份、地点、坐标、成立年份、文章标题、项目、stars、摄影地点和 Wiki 分类均不得作为默认真实内容。
- 文章、作品、摄影和 Wiki 编号以真实资源 ID 生成，补零只改变展示，不伪装为全站连续排名。
- Hero 统计只使用公开 API 的真实总数；新增最小只读 Wiki 统计接口返回公开文档总数，请求失败显示不可用状态，不以 `0` 冒充真实统计。
- Hero 的地点、坐标和成立年份由管理员独立配置；字段为空时不渲染，不从用户资料、备案信息或内容中猜测。
- “全站五色”约束仅作用于主题基础界面；危险、警告、成功状态、用户摄影和代码语法高亮保留独立语义。
- holiday/mourning 继续作为高于个人主题的临时站点覆盖；覆盖期间不改写用户保存的 Swiss 偏好。
- 代码、自动化验证和用户审核均已完成，前端注册表与后端白名单已同步开放。

## 目标

在不改变 InkSpace 业务契约的前提下，把用户侧应用呈现为一套精确、可引用的数字档案系统。公共页面通过外露网格、真实编号和克莱因蓝建立强识别度；用户中心通过同一线系统和高密度排版保留内容管理效率，并确保所有真实业务状态清晰可辨。

## 非目标

- 不新增数据库表、作品 stars、开源状态、license 或 Wiki 分类码字段；地点、坐标和成立年份作为现有 Setting JSON 的可选字段保存。
- 除公开 Wiki 文档总数外，不新增 Swiss 专用业务 HTTP 接口；Hero 文案继续复用公开设置与现有批量保存接口。
- 不把原型中的 `CHEN YU®`、`SHANGHAI`、`31.23°N 121.47°E`、`EST. 2020` 写入产品默认值。
- 不伪造 `42+` 篇文章、`3.4k stars`、`258` 篇知识笔记、项目技术栈、照片地点或固定内容编号。
- 不把 `github_url` 推断成“开源”，也不生成无法验证的 GitHub star 数。
- 不用渐变色块冒充真实摄影；缺图时显示明确的等比例档案占位。
- 不为了严格五色删除 Element Plus 的危险/警告/成功语义，不覆盖管理后台设置的代码高亮 token。
- 不改变文章筛选、作品类型、评论、互动、关注、收藏、通知、知识库权限、发布、版本和分享语义。
- 不改变 Vditor 工具栏、自动保存、全屏、版本和发布行为。
- 不重做 `web/admin` 整体视觉，只在现有首页设置区域增加 Swiss Hero 配置。
- 不引入网络字体、CSS 框架、UI 框架、Canvas 动画或 `vue-i18n`。
- 不在实现和审核完成前开放 `swiss`；未知主题继续回退 `magazine`。

## 用户故事

- 作为登录用户，我想预览并保存“站点名 · 瑞士网格风”，并在其他设备恢复同一偏好。
- 作为读者，我想通过编号和清晰网格快速扫描文章、作品、摄影和公开知识库。
- 作为站点管理员，我想配置 Swiss 首页的真实定位文案，而不是发布原型作者的虚构身份。
- 作为内容作者，我想在 Swiss 下继续完成高密度表格、表单、上传、编辑、发布和版本管理。
- 作为键盘或触屏用户，我想获得与 hover 等价的反色、焦点和操作反馈。
- 作为低动态偏好用户，我想关闭箭头位移和图片缩放，但保留完整的信息层次。
- 作为维护者，我想让编号和统计由真实数据稳定生成，避免内容排序变化后产生误导。

## 触发入口

1. 登录用户进入 `/dashboard/appearance`。
2. 选择“{site_name} · 瑞士网格风”和 `system|light|dark`。
3. 点击“预览效果”后调用现有 `appearance.preview()`，立即设置 `data-ui-theme=swiss`，不请求服务端。
4. 点击“应用主题”后调用 `PUT /api/profile/appearance`。
5. 保存成功后写账号缓存；保存失败、取消或未保存离页时恢复进入页面前偏好。
6. 未登录访客不新增风格选择入口，继续只切换明暗；默认 UI 主题保持 `magazine`。

## 核心流程

### A. 主题开放

1. 新增 `web/blog/src/themes/swiss.css`，规则限定在 `html[data-ui-theme="swiss"]`。
2. 开发期间注册表保持 `coming_soon` 和 `stylesheet:null`，后端继续拒绝 Swiss 保存。
3. 全量自动化和用户审核通过后，注册表改为 `available`、stylesheet 改为 `swiss.css`，后端白名单同步扩展为 `magazine|terminal|cozy|swiss`。
4. 默认仍为 `magazine + system`，开放 Swiss 不迁移或改写任何已有用户偏好。
5. `theme swiss` 和 Tab 补全只在 Swiss 正式开放后自动进入 FloatingTerminal 的可用主题集合。

### B. 首页浏览

1. 公共页头使用真实 `site_name`、英文大写导航和方形明暗入口；手机端使用现有可操作菜单，不隐藏业务入口。
2. Swiss Hero 必须优先于 `home_carousel` 呈现，避免站点配置轮播后退回 Magazine 结构。
3. Hero 左侧显示可配置档案标签、地点、坐标、成立年份、英文标题、强调词、中文说明和 CTA；三个档案字段为空时分别隐藏，不提供虚构默认值。
4. Hero 右侧统计使用公开文章总数、已发布作品总数和公开文档总数；请求失败显示 `--` 与“暂不可用”，不显示虚假零值。
5. 首页 Writing 使用真实文章索引行；Works 展示接口返回的全部作品并将真实封面铺在文字信息层下方；Photos 使用真实摄影版面；Wiki 使用带真实封面或档案占位的公开工作区表格行。
6. 数据不足时按实际数量排版；空态、错误态和图片失败保留完整网格单元，不填充原型示例。

### C. 公共内容与互动

1. `/blog` 保留关键词、分类、标签、排行、排序、分页和 URL query；桌面仍保持文章索引与热门标签区域，不因原型删业务入口。
2. 文章列表行显示真实编号、标题、分类、日期和箭头；详情页保留目录、Markdown、代码、评论、点赞、收藏和前后篇。
3. `/works` 继续区分 project/photography；只显示真实技术栈、互动数和外链，不生成 stars 或“Open Source”。
4. `/photos` 和摄影详情使用真实图片、标题、地点和 EXIF；默认降低饱和或灰度只作为视觉效果，不表达状态。
5. `/wiki` 将公开工作区呈现为档案行；工作区和文档详情继续遵守公开/已发布过滤和 HTML 净化契约。
6. About、Links、用户搜索、个人主页、分享和登录页使用同一零圆角线系统，并覆盖正常、空、错误、失效和过期状态。

### D. 用户中心与创作

1. 用户中心侧栏使用精确索引视觉和真实 `site_name`，不改变菜单顺序、权限或通知入口。
2. Dashboard、文章、作品、评论、收藏、通知、资料和外观页保持当前信息密度与完整操作。
3. 表格、筛选器、分页和操作列以一像素线分区；按钮不因反色或大写文案换行、遮挡或丢失危险语义。
4. 创建/编辑表单、上传、裁剪、确认框、抽屉、下拉、日期控件、loading 和校验错误完整适配。
5. Workspace 管理保留公开开关、目录拖拽、文档状态、发布、版本和分享；蓝色不得替代这些业务状态。
6. Vditor 外壳使用零圆角和线框；Markdown 明暗跟随实际模式，代码高亮继续由 `code_theme` 控制。

### E. 跨主题终端

1. 已打开的 FloatingTerminal 在切换 Swiss 后继续保留 cwd、历史、窗口位置和 pending confirmation。
2. Swiss 下窗口映射 `--surface|--bg-soft|--hairline|--ink|--accent`，使用零圆角、无阴影和克莱因蓝焦点。
3. 命令、虚拟文件系统、拖拽、缩放、移动抽屉和确认流程不变，不增加 Swiss 专属命令。
4. 层级继续低于 Element Plus 模态框和 Vditor 全屏；切换主题不卸载全局终端。

## 真实内容映射

### Hero 配置与统计

- 文案键：`home_hero_swiss`，JSON 字符串，复用 `/admin/settings/batch` 和 `/settings/public`。
- 配置字段：`eyebrow`、`location`、`coordinates`、`established`、`title`、`accent`、`description`、`primary_text`、`primary_link`、`secondary_text`、`secondary_link`。
- 默认眉头：`A1 / INDEX`；默认标题：`WRITE, BUILD & ARCHIVE`；强调词：`ARCHIVE`。
- 默认简介：`文章 × 作品 × 摄影 × 知识库`；主按钮：`START READING` -> `/blog`；次按钮：`VIEW WORKS` -> `/works`。
- `location`、`coordinates`、`established` 默认均为空；管理端分别提供“地点”“坐标”“成立年份”输入，不校验为特定城市或坐标格式，最大长度分别为 80、80、40 个字符，并只做安全文本渲染。
- 品牌始终读取 `site_name`，不放入 Hero JSON，避免后台两处站点名称互相冲突。
- 统计来源：`/articles?page=1&page_size=1` 的 `total`、`/works?page=1&page_size=1` 的 `total`、`GET /api/wiki/stats` 的 `public_doc_count`。
- `GET /api/wiki/stats` 在数据库中统计 `is_public=true` 工作区下 `status=1` 且未软删除的文档，返回全站公开文档总数；不把分页工作区的 `doc_count` 相加。
- 统计接口是通用公共 Wiki 能力，不读取私人工作区、草稿文档、分享 token 文档或前端传入的筛选条件。

### 稳定编号

- 新增纯函数 `formatSwissCode(prefix, id, minLength)`；以真实正整数 ID 生成，补零但不截断。
- 文章：`W–${article.id}`，最少三位，例如 ID 42 -> `W–042`。
- 作品：`P–${work.id}`；摄影：`PL.${work.id}`；公开工作区：`K–${workspace.id}`。
- ID 缺失时显示 `W–—` 等明确占位，不使用循环索引伪造资源身份。
- 当前版面位置可另显示 `01/02/03`，但必须标为展示序号，不能作为资源永久编号或 URL 参数。

### 文章索引

- 标题：`article.title`；分类：`article.category.name`；日期：`article.created_at`；链接：`/blog/:id`。
- 首页现有 hot 列表不是时间归档，不使用“总数减循环索引”生成倒序档案号。
- 分类缺失时显示 `UNCATEGORIZED`，不从标签或标题猜测分类。
- 超长标题允许两行并安全截断；详情页显示完整标题。

### 作品单元格

- 标题、类型、描述、技术栈、点赞/浏览数、作者和链接均来自真实 `WorkResponse`。
- `tech_stack` 按现有逗号分隔值展示；空值时省略整段，不显示示例技术栈。
- `github_url` 只决定 GitHub 链接是否出现；不推断 license、开源状态或 star 数。
- photography 类型不套用 project 技术栈语言，改用真实地点、拍摄日期或照片数；字段缺失时省略。

### 摄影版面

- 数据源：`GET /api/works?type=photography&status=1&page=1&page_size=6`。
- 图片优先级：`cover` -> `images[0].url|string`；均缺失或加载失败时显示 `NO IMAGE / PL.id` 等比例占位。
- 图注只显示真实 `title` 和存在的 `metadata.location`；不从描述推断地点。
- 版面按实际位置稳定使用 `6|3|3|3|3|6` 的十二列跨度节奏；不足六项时不复制内容。
- 图片灰度/缩放只作用于页面展示，不修改原图、查看器、下载或正文图片。

### Wiki 档案行

- 标题、描述、`doc_count`、作者和更新时间来自 `PublicWorkspaceResponse`。
- 分类码使用真实 ID 生成 `K–xxx`；不使用原型 `K–FE/K–DS/K–RD`，因为模型没有分类码字段。
- `doc_count` 以 `${n} NOTES` 呈现；这是单个公开工作区的已发布文档数，不表示全部私人或草稿文档。
- 图标只作为辅助内容，不用其文本猜测分类；路由保持 `/wiki/:workspaceId`。

## 异常处理

| 场景 | 处理方式 |
|---|---|
| 服务端或缓存返回未开放 Swiss | 归一化为 `magazine`，不写回错误偏好 |
| Swiss 样式加载失败 | 基础令牌保证页面可读；取消预览或保存失败时恢复旧主题 |
| `system` 媒体查询不可用 | 实际模式回退 `light` |
| `site_name` 为空 | 显示 `InkSpace`，不显示 `CHEN YU®` |
| Hero JSON 无效或字段为空 | 按字段保留内置默认值，不让空字符串覆盖默认文案 |
| 地点、坐标或成立年份为空 | 分别隐藏对应标签，不显示分隔符或占位值 |
| 某项 Hero 统计失败 | 对应统计显示 `--` 和不可用说明；其他统计继续显示 |
| 文章/作品/Wiki ID 缺失 | 显示破折号占位，不使用数组索引伪造 ID |
| 摄影不足六张 | 按真实数量和跨度顺序布局，不复制照片 |
| 图片缺失或加载失败 | 保留 aspect-ratio，显示文字档案占位和真实 alt |
| 技术栈、地点、作者或描述缺失 | 省略可选字段，不显示原型内容 |
| 超长中英文标题 | 安全换行并调整字距；不得撑破 Hero 或行网格 |
| holiday/mourning 生效 | 临时覆盖 Swiss 色彩，不改写账号偏好；结束后恢复 Swiss |
| danger/warning/success 状态 | 保留语义色、图标和文字，不强制改成蓝色 |
| 触摸设备无 hover | 操作、编号、图注和箭头始终可见，不依赖悬停恢复信息 |
| reduced-motion | 取消箭头位移、图片缩放、颜色 transition 和平滑滚动 |
| 宽表格、代码和长 URL | 横向滚动或安全换行，不压缩至不可读宽度 |
| Teleport 浮层 | 使用 Swiss 全局主题选择器，不残留默认圆角、阴影或旧主题色 |

## 视觉系统

### 颜色令牌

`swiss.css` 通过兼容别名映射到现有 `--theme-*` 和 Element Plus 变量。

| Token | 浅色 | 深色 | 用途 |
|---|---|---|---|
| `--bg` | `#FFFFFF` | `#0C0C0C` | 页面背景 |
| `--bg-soft` | `#F7F7F7` | `#171717` | 表头、次级单元格、禁用底 |
| `--surface` | `#FFFFFF` | `#0C0C0C` | 卡片、浮层、编辑器 |
| `--ink` | `#111111` | `#F2F2F2` | 正文、主结构线、反色底 |
| `--sub` | `#767676` | `#A0A0A0` | 辅助文字；深色较原型提高小字对比度 |
| `--hairline` | `#E4E4E4` | `#2E2E2E` | 单元格分隔线 |
| `--accent` | `#002FA7` | `#5B7FF0` | 激活、焦点、编号、CTA |
| `--accent-hover` | `#111111` | `#F2F2F2` | 反色交互 |
| `--accent-soft` | `#F2F5FC` | `#141A2B` | 行悬停；作为 `color-mix()` 回退 |
| `--accent-ink` | `#FFFFFF` | `#FFFFFF` | 蓝底文字 |
| `--shadow-color` | `transparent` | `transparent` | 禁止主题阴影 |

语义状态继续使用 Element Plus 现有 success/warning/danger 变量，并必须同时显示文字或图标。摄影、头像和代码高亮不受基础界面五色约束。

### 字体与排版

- 主字体：`'Helvetica Neue', Helvetica, Arial, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif`。
- 等宽编号：`'SF Mono', 'IBM Plex Mono', Consolas, monospace`，不加载网络字体。
- Hero H1：`clamp(42px, 7.2vw, 96px)`，800，行高 1.04；手机最小值可降至 36px。
- Hero 中文说明：`clamp(18px, 2.6vw, 34px)`；字距使用 `clamp(.08em, 1vw, .3em)`，防止长文案溢出。
- 栏目 H2 22px/800；条目标题 17px/600；正文 15px/1.7；编号与标签不得低于 12px。
- 英文可大写并收紧字距；中文不强制大写或逐字断开。
- Markdown 正文保持约 760px 阅读宽度和稳定行高，不扩展到十二列全宽。

### 网格与形状

- 公共内容最大宽度 1240px，桌面两侧 32px；用户中心沿用现有可用宽度。
- 首页签名区显式使用十二列 CSS Grid；普通业务页只在能改善信息结构时使用网格，不为装饰重排 DOM。
- header、section、footer 主结构使用 1px `--ink`；内部单元使用 1px `--hairline`，避免双线叠加。
- 卡片、按钮、输入、标签、弹窗、抽屉、图片和编辑器均为零圆角；不使用 box-shadow 或渐变背景。
- 相邻首页区块共享边线，不增加卡片式大间距；详情阅读页仍保留必要正文留白。
- 图片网格使用 `gap:1px` 和 `background:var(--hairline)`，图片本身不增加装饰边框。

### 签名元素

- **可见十二列网格**：Hero、摄影和关键首页区域以结构线暴露栅格。
- **真实档案编号**：W/P/PL/K 前缀全部由资源 ID 生成。
- **克莱因蓝响应**：只用于主要交互、激活、焦点和档案编号，不大面积装饰正文。
- **反色工作单元**：hover 与 focus-within 使用 `ink/bg` 反色，触屏保持静态可读。
- **摄影 Plates**：真实照片轻度灰度，hover/focus 恢复原色并轻微放大。
- **Colophon**：显示真实版权；右侧可写 `GRID 12 COL / KLEIN BLUE 002FA7`，作为主题说明而非业务数据。

### 动效

- 行反馈 `.15s`；箭头位移 `.2s`，最多 6px。
- 反色单元 `.2s`，无弹性、旋转、入场或自动播放动画。
- 摄影 `.5s cubic-bezier(.2,.6,.2,1)`，最多 scale 1.04。
- 明暗切换 `.3s`；不对布局尺寸做过渡。
- `prefers-reduced-motion: reduce` 下移除全部位移、缩放、平滑滚动和非必要 transition。

## 页面适配矩阵

### 公共骨架与页面

| 路由/区域 | 文件 | Swiss 要求 |
|---|---|---|
| 公共骨架 | `layouts/MainLayout.vue` | 真实站点品牌、方形导航、移动菜单、网格页脚、明暗入口 |
| 首页 `/` | `views/Home.vue` | Swiss Hero、真实统计、Writing 行、Works 单元、Plates、Wiki 行 |
| 博客 `/blog` | `views/Blog.vue` | 索引行、完整 query 筛选、热门标签、分页、空错态 |
| 文章 `/blog/:id` | `views/BlogDetail.vue` | 编号页头、正文、目录、Markdown、互动、评论、前后篇 |
| 作品 `/works` | `views/Works.vue` | project/photography 单元、真实 metadata、筛选搜索分页 |
| 作品详情 `/works/:id` | `views/WorkDetail.vue` | 项目/摄影分支、图集、EXIF、外链、互动评论 |
| 摄影 `/photos` | `views/Photos.vue` | 真实 Plate 网格、分页、图片失败和触屏操作 |
| Wiki `/wiki` | `views/wiki/WikiIndex.vue` | K-ID 档案行、真实文档数、作者和分页 |
| Wiki 空间 | `views/wiki/WikiWorkspace.vue` | 线框目录/文档索引、移动抽屉 |
| Wiki 文档 | `views/wiki/WikiDoc.vue` | 正文、面包屑、代码表格图片和阅读元数据 |
| 站点页 | `About.vue`、`Links.vue`、`UserSearch.vue` | 网格信息块、表单、空错态和安全外链 |
| 用户主页 | `UserProfile.vue` | 资料、文章、作品、收藏、关注和粉丝完整可用 |
| 分享 | `ShareDoc.vue` | 正文及有效、失效、过期、错误状态 |
| 登录 | `Login.vue` | 独立页面令牌、方形表单、校验与回跳 |
| 浮动终端 | `components/terminal/FloatingTerminal.vue` | Swiss 外壳，命令和窗口状态不变 |

### 用户中心

| 区域 | 页面 | Swiss 要求 |
|---|---|---|
| 骨架 | `layouts/UserCenterLayout.vue` | 索引侧栏、面包屑、通知和移动菜单 |
| 概览 | `views/user/Dashboard.vue` | 数据格与快捷入口，真实状态不只依赖颜色 |
| 文章 | `MyArticles.vue`、`ArticleEdit.vue` | 筛选、状态、操作单行、编辑、发布和封面 |
| 作品 | `MyWorks.vue`、`WorkEdit.vue` | 类型、状态、图片、技术栈、外链和提交 |
| 社交与状态 | `MyComments.vue`、`Notifications.vue`、`Favorites.vue` | 列表、未读、删除、跳转和确认 |
| 资料 | `ProfileEdit.vue` | 表单、头像上传裁剪和失败反馈 |
| 知识库 | `WorkspaceList.vue`、`WorkspaceDetail.vue`、`DocEdit.vue` | 公开开关、目录树、拖拽、状态、版本、发布和分享 |
| 外观 | `Appearance.vue` | Swiss 卡、真实站点名 Live Proof、预览取消保存和明暗模式 |

## Element Plus 与内容渲染

### Element Plus

全局 Swiss 主题至少覆盖 button、card、input、textarea、select、radio、checkbox、switch、tag、badge、table、tabs、pagination、menu、breadcrumb、steps、skeleton、empty、alert、upload、carousel、avatar 和 form validation。

Teleport 节点必须覆盖 dialog、drawer、message-box、message、notification、loading mask、dropdown、tooltip、popover、date picker 和 image viewer。所有表面零圆角、无阴影、1px 边界；主操作使用 accent，危险/警告/成功保留语义色和文字。

### Markdown、Vditor 与代码高亮

- Markdown 和 Vditor 继续按 `appearance.resolvedColorScheme` 选择浅深模式。
- Swiss 控制正文背景、标题、引用、分隔线、表格、行内代码和代码容器边界，不覆盖语法 token。
- `code_theme` 继续由管理后台绝对控制；用户内容中的颜色和图片不强制转换为五色。
- Vditor 主题切换调用现有 `setTheme()`，不重建编辑器或丢失未保存内容。
- 全屏层级保持高于页头、侧栏、FloatingTerminal 和所有网格装饰。

## 数据与 API

### 数据模型

不新增数据库字段。偏好继续使用：

```json
{
  "ui_theme": "swiss",
  "color_scheme": "system"
}
```

正式开放时仅扩展 `UserAppearanceRequest` 白名单。新增 `SettingHomeHeroSwiss = "home_hero_swiss"` 常量属于公开设置契约扩展，不需要 schema 迁移；地点、坐标和成立年份保存在该设置的 JSON 中，也不新增列。

新增只读响应 DTO，不持久化，字段 `public_doc_count` 为非负整数；实际值必须由公开文档查询实时统计，统一响应外层仍为 `{code,message,data}`。

### API

| 方法 | 路径 | 用途 |
|---|---|---|
| GET/PUT | `/api/profile/appearance` | 读取和保存 Swiss 偏好 |
| GET | `/api/settings/public` | 站点名、Swiss Hero、特殊覆盖和代码主题 |
| PUT | `/api/admin/settings/batch` | 管理员保存 `home_hero_swiss` |
| GET | `/api/articles/*` | 文章总数、首页索引、列表与详情 |
| GET | `/api/works/*` | 作品总数、项目、摄影和详情 |
| GET | `/api/wiki/stats` | 公开文档总数；免登录，不接受筛选参数 |
| GET | `/api/wiki/*` | 公开工作区、目录和文档 |
| 既有接口 | 用户、评论、互动及用户中心 | 原业务操作 |

`home_hero_swiss` 必须在 Setting service 新建和更新路径中都强制归类为公开 `carousel` 设置，避免已有私有同名记录无法被 `/settings/public` 返回。

`GET /api/wiki/stats` 必须在 Gin 路由中的 `/api/wiki/:workspaceId` 之前注册，避免静态段 `stats` 被动态工作区路由捕获。

## 响应式要求

### `<=900px`

- Hero 变单列；统计改为三列横排，过长说明安全换行。
- Writing 和 Wiki 隐藏非必要 meta/描述列，但编号、标题和箭头保留。
- Works 改为单列，列分隔线改为底线；照片统一 span 6。
- 页脚改为单列；导航使用现有移动菜单，所有路由入口可达。
- 用户中心侧栏使用现有移动抽屉；表格横向滚动或转结构化列表。
- 弹窗、图片查看器和编辑器不得超出视口。

### `<=560px`

- 页面水平 padding 16px；Hero H1 和中文字距按视口收缩，不出现横向滚动。
- Hero 统计允许纵向排列；CTA 提供至少 44px 触摸高度。
- Writing/Wiki 行使用 `56px 1fr 36px`；超长标题最多两行。
- 摄影全部 span 12；图注始终可见，不依赖 hover。
- 表单输入字号至少 16px；按钮组允许在不改变顺序时合理换行。
- 外观页预览、取消和应用始终可见；FloatingTerminal 保持底部抽屉与安全区。

## 可访问性与性能

- 焦点使用 2px `--accent` 实线和至少 2px offset；反色 hover 必须有等价 `:focus-visible` 或 `:focus-within`。
- 关键正文对比度至少 4.5:1，边界与大文字至少 3:1；深色小编号使用提亮 accent。
- 蓝色不是唯一状态线索；激活项同时使用底色/边界/文字，业务状态同时有文字或图标。
- 首页和博客中 click-only 卡片在结构改造时改为真实链接或按钮，支持 Enter、Space 和正确名称。
- 摄影、封面、头像和正文图片使用真实 alt；结构线、网格水印和装饰编号按语义决定是否 `aria-hidden`。
- 图片设置 `loading=lazy`、尺寸或 `aspect-ratio`；首屏关键图按实际性能决定 eager/fetchpriority。
- 不加载网络字体，不创建自动动画，不通过大面积滤镜影响页面滚动性能。
- `color-mix()` 必须有静态 `--accent-soft` 回退；基础边界在旧浏览器仍可读。
- 页面层级保持：sticky header < 用户侧栏 < FloatingTerminal < Element Plus 模态层 < Vditor 全屏。

## 实现步骤

每步可独立审查；正式开放放在最后。

1. [x] **状态与测试基线**：记录当前 Go、blog、admin 测试和构建基线。
2. [x] **Swiss 令牌**：新增 `swiss.css`，实现浅深令牌、线系统、Element Plus、Teleport、Markdown、Vditor、特殊覆盖和 reduced-motion。
3. [x] **真实映射工具**：实现 `formatSwissCode` 与摄影跨度映射纯函数，覆盖大 ID、缺失 ID 和稳定性测试。
4. [x] **公共与用户中心骨架**：适配 MainLayout、UserCenterLayout、真实站点名、移动菜单、页脚和 FloatingTerminal。
5. [x] **Swiss Hero 设置**：新增 `home_hero_swiss` 后端公开设置及管理端表单，支持可选地点、坐标、成立年份、默认值、解析和保存。
6. [x] **公开 Wiki 统计**：新增免登录 `GET /api/wiki/stats`，在 service 中按公开工作区、已发布文档和软删除条件统计，补充 handler、路由和测试。
7. [x] **首页签名结构**：实现真实统计 Hero、Writing、Works、最多六张 Plates 和 Wiki 档案行，覆盖独立失败与空态。
8. [x] **公共索引与详情**：通过全局主题令牌适配 Blog、BlogDetail、Works、WorkDetail、Photos、Wiki、站点页、用户主页、分享和登录。
9. [x] **用户中心与编辑器**：适配全部列表、表格、表单、浮层、上传、知识库和 Vditor 页面。
10. [x] **语义交互与外观页**：首页签名结构使用语义链接，增加 Swiss Live Proof，验证预览、保存和 system 基础流程。
11. [x] **自动化与构建**：前后端测试、lint、build、vet 和 diff 检查通过。
12. [x] **完整矩阵审核**：三视口、浅深、全路由、空错态和无障碍完成用户审核，Swiss 保持开放。
13. [x] **规范同步**：已更新本规范、总主题规范、前阶段历史说明、验证记录和剩余风险。

## 测试计划

### 后端

- [x] `swiss + system|light|dark` 均通过模型校验。
- [x] `magazine`、`terminal`、`cozy` 继续合法，未知值、空主题和未知明暗继续返回 400。
- [ ] Swiss PUT 后 GET 一致，用户之间偏好不污染，未登录请求返回 401。
- [x] `home_hero_swiss` 新建和更新后均为公开 `carousel` 设置，可由 `/settings/public` 返回。
- [ ] `GET /api/wiki/stats` 只统计公开工作区内已发布且未软删除的文档，不包含私人工作区、草稿和已删除文档。
- [ ] Wiki stats 免登录可读，返回统一成功结构，不接受客户端 owner/status 条件绕过公开范围。

### 前端自动化

- [x] 注册表 stylesheet 为 `swiss.css` 且可用集合包含 Swiss。
- [x] 启动缓存、预览、取消、访客缓存、账号保存和失败回滚复用现有主题运行时并正确处理 Swiss。
- [ ] `system` 跟随浅深变化；显式 light/dark 不跟随。
- [x] 编号基于真实 ID，补零不截断，缺失 ID 不回退循环索引。
- [ ] Hero 三项统计独立处理成功/失败，不把失败当 0。
- [x] 地点、坐标、成立年份按配置独立显示；空字段不渲染标签和多余分隔符。
- [x] Plates 最多六项、不复制数据，跨度稳定，图片失败保留尺寸。
- [x] Wiki 使用真实 ID 和 `doc_count`，不生成原型分类码。
- [x] FloatingTerminal 在 Swiss 下可读，命令与补全支持 Swiss。
- [x] Appearance 使用真实 `site_name` 的 Swiss Live Proof 并可保存。
- [ ] holiday/mourning 与 Swiss 并存且不改写保存偏好。
- [ ] reduced-motion 禁用箭头位移、图片缩放和非必要 transition。

### 验证命令

- [x] `gofmt` 修改过的 Go 文件
- [x] `go test ./... -count=1`
- [x] `go vet ./...`
- [x] `web/blog: vitest run`（13 个测试文件、78 项测试）
- [x] `web/blog: eslint src --ext .vue,.js,.jsx,.cjs,.mjs`
- [x] `web/blog: vite build`（通过，保留既有大 chunk 警告）
- [x] `web/admin: vite build`（通过，保留既有大 chunk 警告）
- [x] `git diff --check`

### 人工审核矩阵（已通过）

视口：`1440x900`、`900x1024`、`390x844`。

主题轴：Swiss light 全路由三视口；Swiss dark 全路由桌面和手机；Swiss system 在首页、外观页、Vditor 和 FloatingTerminal 切换系统浅深；holiday/mourning 覆盖关键公共页、登录、Dashboard、编辑器、弹窗和终端。

数据轴：正常、空列表、API 独立失败、404、失效分享、0-6 张摄影、图片失败、无分类/技术栈/地点、无 Wiki 描述、超大 ID、超长中英文标题、长代码、宽表格、loading、上传和校验错误。

交互轴：键盘、触屏、hover/focus 等价、移动菜单、抽屉、弹窗、图片查看器、Vditor 全屏、预览取消保存、终端拖拽和移动抽屉。

## 验收标准

1. Swiss 可在外观页预览、取消和保存，账号同步、缓存、失败回滚与三种明暗模式正确。
2. 默认仍为 `magazine + system`；Magazine、Terminal、Cozy 无回归。
3. 全部可达用户侧路由完成 Swiss 浅深适配，原业务、权限、URL query 和 API 契约不变。
4. 首页具有十二列 Hero、真实统计、Writing 行、Works 单元、真实 Plates 和 Wiki 档案行。
5. 品牌读取真实 `site_name`；地点、坐标和成立年份只显示管理员配置值，不出现 `CHEN YU®`、虚构档案字段、stars、内容或统计。
6. W/P/PL/K 编号均基于真实 ID，统计均来自公开 API；失败和缺失显示明确占位。
7. 零圆角、零阴影、零界面渐变、结构线和克莱因蓝语言一致；语义状态、摄影和代码高亮不被错误压成五色。
8. Element Plus Teleport、Markdown、Vditor、代码块、图片查看器和 FloatingTerminal 无旧主题残留。
9. `900px` 和 `560px` 下导航、Hero、索引行、Plates、Wiki 树、表格、弹窗、表单和编辑器可操作且不溢出。
10. 语义链接、键盘顺序、alt、focus-visible、触摸目标、对比度、图片防抖和 reduced-motion 达标。
11. holiday/mourning 可覆盖 Swiss 呈现但不改写偏好，覆盖结束后恢复 Swiss。
12. 自动化、构建和用户审核均已通过，Swiss 已标记为 `available` 并加入服务端白名单。

## 已确认的范围决策

- Swiss Hero 需要管理员可配置的地点、坐标和成立年份；三个字段默认空并独立隐藏。
- 新增全站公开文档总数聚合接口，作为 Hero 第三项真实统计；必须严格遵守公开工作区和已发布文档边界。
- 不新增作品 license、开源状态或 GitHub stars 同步能力；当前只展示真实现有字段。

## 最终审核调整

- 首页 Works 展示现有热门作品接口返回的全部作品，不截断前三项。
- 作品真实封面铺满单元格底层，标题、编号、元数据和简介沉淀在底部半透明纯色信息层；不使用渐变。
- 信息层降低高度和不透明度并移除模糊，让封面主体清晰可见；hover/focus 保留可读反色反馈。
- 作品缺图或加载失败时使用真实 `P–xxx` 编号保持同尺寸档案占位。
- 首页 Wiki 使用 `workspace.icon` 中的安全真实图片作为封面；无图或加载失败时使用图标、名称首字或 `K` 档案占位。

## 待定事项

- 是否保留完整截图矩阵；不阻塞实现，但建议至少保留首页、Blog、作品详情、摄影、Wiki、Dashboard、编辑器和外观页的浅深/手机基线。

## MVP 范围

本阶段按完整主题交付，不拆分“只做浅色”“只做首页”或“先开放再补用户中心”：

1. Swiss 浅色、深色、跟随系统令牌和 holiday/mourning 覆盖。
2. 全部公共页面、用户中心、编辑器、Element Plus 浮层和跨主题终端适配。
3. 真实站点品牌、含可选档案字段的 Hero、公开文档聚合接口、真实三项统计、稳定资源编号和真实内容版面。
4. Writing、Works、Plates、Wiki 四种签名结构及完整空错态。
5. 外观页预览、取消、保存、回滚、账号缓存和最终前后端白名单开放。
6. 三视口、全路由、浅深、无障碍、reduced-motion、自动化与用户审核。
