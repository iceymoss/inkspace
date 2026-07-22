# inkspace.log 交互式终端模拟器

## 状态
- 创建日期: 2026-07-22
- 状态: 草稿，待确认
- 关联主题: `terminal` / `inkspace.log`
- 上游规范: [`terminal-ui-theme.md`](./terminal-ui-theme.md)

## 已确认决策

- 首页 Terminal Hero 中的终端是唯一启动入口；点击后从原位升起为可操作浮窗。
- 浮窗激活后挂载在应用根节点，跨公共页、登录页和用户中心路由持续存在，直到用户关闭。
- 桌面支持拖拽和缩放；手机端降级为近全屏底部抽屉，不支持自由拖拽。
- 窗口位置、尺寸、开关状态和命令历史只保存在当前标签页会话，不同步账号，不写长期本地缓存。
- 命令采用类 Unix 语法；不执行 shell、JavaScript、HTML 或任意远程代码。
- 导航、搜索和筛选优先复用 Vue Router 路径及查询参数，刷新、前进和后退后状态可复现。
- 查询复用现有业务 API；点赞、收藏和关注复用现有登录态与权限，不新增终端专用后端接口。
- 所有写操作先在终端输出操作摘要，用户输入 `yes` 后执行，输入 `no` 或按 `Escape` 取消。
- MVP 一次性交付浮窗、拖拽缩放、命令系统、查询、写操作确认、跨路由和移动端能力。
- 博客日志等级不随机；根据置顶状态、文章 ID、分类和标签稳定映射到 `FEAT|INFO|NOTE|WARN`。
- 终端提供映射真实路由和公开内容的安全虚拟文件系统；首页目录为 `/inkspace/index`，支持 `ls|cd|pwd|cat|grep` 路径语法。

## 目标

将首页装饰性终端升级为安全、可持续浮动的站点控制终端，让用户通过可发现的命令完成导航、主题切换、页面搜索筛选、内容查询和明确确认后的互动操作，同时保持现有页面、API、权限和浏览器导航语义。

## 非目标

- 不提供真实 shell、文件系统、JavaScript `eval`、HTML 注入、SQL、服务器命令或插件脚本执行。
- 不允许命令绕过登录、资源可见性、内容发布状态、作者权限或后端业务校验。
- 不新增 WebSocket、RPC、消息队列或实时协作。
- 不持久化完整命令历史到账号、数据库或长期 `localStorage`。
- 不让终端直接读取 token、Cookie、环境变量、管理端配置原文或私有草稿。
- 不实现公共 Wiki 全文检索；现有 API 只支持公开工作区列表、目录树和文档详情。
- 不在 MVP 支持管道、重定向、子 shell、命令替换、通配符、后台任务或用户自定义命令。
- 不把 Terminal Hero 和全局浮窗复制成两个同时运行的终端实例。
- 不改变 `magazine` 的首页 Hero、主题视觉和默认访客体验。

## 用户故事

- 作为 Terminal 主题访问者，我想点击首页终端让它从 Hero 中升起，以获得连续且可信的交互体验。
- 作为站点访问者，我想拖动和缩放终端，同时继续查看页面内容。
- 作为手机访问者，我想使用适配触控的底部终端抽屉，而不是操作难以控制的小浮窗。
- 作为读者，我想通过命令跳转、搜索和筛选文章、作品、摄影、Wiki 与用户，并让页面 URL 同步变化。
- 作为登录用户，我想在明确确认后通过终端点赞、收藏或关注，且页面状态能同步刷新。
- 作为访客，我希望需要登录的命令明确说明限制，并可前往登录，而不是静默失败或绕过权限。
- 作为键盘用户，我想使用命令历史、补全、焦点管理和可读输出完成全部终端操作。
- 作为维护者，我希望命令注册、解析、执行和 UI 分离，便于测试和后续增加安全命令。

## 触发入口

1. 用户使用 `terminal` 主题访问首页。
2. Terminal Hero 窗口显示“点击启动交互终端”的可操作状态；标题栏和正文点击均可启动，内部链接仍按原行为导航。
3. 启动时记录 Hero 终端的视口矩形，渲染全局浮窗并从该矩形过渡到会话中的目标位置和尺寸。
4. Hero 原位置显示轻量占位和“终端已在运行”状态，避免页面布局跳动，也避免出现第二套输入框。
5. 用户关闭浮窗后，它回到 Hero 原位；如果当前不在首页，则直接关闭，回到首页后恢复静态 Hero。
6. 浮窗打开期间可跨路由持续使用；切换到非 Terminal 主题时功能保留，但使用当前主题的通用令牌呈现，不强制改回 Terminal。

## 核心流程

### A. 启动与浮动

1. 用户点击首页 Terminal Hero。
2. `TerminalHero` 发出带源矩形的 `activate` 事件，不自行创建全局状态或请求 API。
3. 根级 terminal store 标记为打开，并将源矩形转换为受视口约束的初始浮窗位置。
4. 桌面端聚焦命令输入框；输出欢迎语、当前路径和 `help` 提示。
5. 用户通过标题栏拖动，通过右下角缩放手柄调整尺寸；位置与尺寸写入 `sessionStorage`。
6. 页面导航后根组件不卸载，浮窗、历史和输入保持。

### B. 命令输入

1. 用户输入命令并按 Enter。
2. 解析器执行引号感知的分词，只接受命令注册表中的白名单语法。
3. 终端先追加输入记录，再执行对应命令；异步命令期间显示 pending 状态并禁止重复提交。
4. 成功时输出结构化文本和可聚焦结果链接；失败时输出错误消息，不插入不可信 HTML。
5. `ArrowUp/ArrowDown` 浏览当前会话命令历史；`Tab` 在已知命令、子命令和静态参数间补全。
6. `Ctrl+L` 或 `clear` 清空可见输出但保留命令历史；`Escape` 取消待确认操作，否则最小化终端。

### C. 页面导航和筛选

1. `open`、`search`、`filter` 命令转换成白名单路由和查询参数。
2. Router 更新 URL，目标页面从 route query 初始化并监听变化。
3. 页面使用原有 API 加载流程，不由终端复制业务列表状态。
4. 瞬时操作如 `scroll top`、`scroll bottom`、`focus search` 通过受控页面动作执行，不写 URL。
5. 浏览器前进/后退恢复搜索和筛选状态，终端输出不重复执行历史命令。

### C.1 虚拟文件系统

1. 根目录固定为 `/inkspace`，不允许通过 `..` 逃逸。
2. 首页路由 `/` 映射到 `/inkspace/index`；一级目录包括 `index|blog|works|photos|wiki|users|about|links`。
3. `ls` 列出当前虚拟目录；静态一级目录不调用 API，内容目录按需调用现有公开 API。
4. `cd ../blog` 解析规范化路径并调用 `router.push('/blog')`；路由变化也反向更新终端当前目录。
5. 文章映射为 `blog/<id>-<slug>.md`，作品映射为 `works/<id>-<slug>.json`，摄影映射为 `photos/<id>-<slug>.jpg`，用户映射为 `users/<id>-<slug>/`。
6. Wiki 工作区映射为 `wiki/<workspaceId>-<slug>/`；公开文档映射为工作区目录中的 `<docId>-<slug>.md`，目录树按公开 API 结果生成。
7. `cat` 对内容文件读取对应公开详情 API，并同步导航到详情页；默认输出标题、元数据和安全截断后的纯文本摘要，不在终端渲染 `content_html`。
8. `grep <keyword> blog/|works/|photos/|users/` 调用现有关键词 API，并同步进入带 query 的列表页。
9. `favorite|like|follow` 可接收资源类型加 ID，也可接收虚拟路径；路径最终只解析为白名单资源类型和正整数 ID。
10. 所有名称 slug 只用于可读显示，路由与 API 以 ID 为准；重名、改名或非法字符不影响资源定位。

### D. 数据查询

1. `list` 命令调用现有公开 API，并将精简结果输出到终端。
2. 每项结果包含真实 ID、标题/名称和可操作 `open` 提示。
3. 默认返回前 5 项，最大 20 项；禁止无界查询。
4. 查询文章、作品、摄影和用户时使用现有 keyword 参数。
5. Wiki 仅列出公开工作区，或在给定工作区 ID 后读取公开树；不声称提供全文搜索。

### E. 写操作与确认

1. 用户输入写命令，例如 `favorite article 12`。
2. 执行器验证语法、登录态和目标 ID，不满足时不调用 API。
3. 终端进入单一 pending-confirmation 状态，输出请求方法语义、资源类型与 ID，并提示 `[yes/no]`。
4. 只有精确输入 `yes` 才执行；`no`、`Escape` 或其他命令取消确认。
5. 执行时使用统一 `api` 客户端；成功后输出服务端结果。
6. 若当前路由正显示同一资源，发出受控 refresh 信号，让详情页重新读取状态和计数；不直接篡改页面局部 ref。
7. 401/403/404 和业务错误沿用现有拦截器，同时在终端打印 `error.message`；不得重复弹出第二个 Element Plus 错误提示。

### F. 主题与明暗

1. `theme magazine|terminal` 调用现有 appearance store。
2. 访客只允许切换当前浏览器的可用主题与明暗；主题风格切换是否开放给访客遵循现有产品约束，MVP 中终端命令仅允许登录用户保存主题，访客可预览到会话结束。
3. `scheme system|light|dark` 复用现有访客保存或账号预览流程。
4. `status` 输出当前路由、登录状态、UI 主题、明暗偏好和实际明暗，不输出用户敏感字段。

## 命令规范

### 基础命令

| 命令 | 行为 |
|---|---|
| `help [command]` | 显示命令列表或单个命令帮助 |
| `status` | 当前路由、登录状态、主题与明暗状态 |
| `pwd` | 输出当前路由路径 |
| `clear` | 清空可见输出 |
| `history` | 输出当前会话命令历史，最多 50 条 |
| `close` | 关闭并回收终端 |
| `minimize` | 最小化为可恢复的状态条 |
| `ls [path]` | 列出虚拟目录；内容目录读取真实公开数据 |
| `cd <path>` | 切换虚拟目录并联动前端路由 |
| `cat <path>` | 读取虚拟内容文件并联动详情页 |
| `grep <keyword> <path>` | 查询映射目录并联动列表页筛选 |

### 导航命令

| 命令 | 路由 |
|---|---|
| `open home` | `/` |
| `open blog` | `/blog` |
| `open works` | `/works` |
| `open photos` | `/photos` |
| `open wiki` | `/wiki` |
| `open about` | `/about` |
| `open links` | `/links` |
| `open dashboard` | `/dashboard`，需登录 |
| `open article <id>` | `/blog/:id` |
| `open work <id>` | `/works/:id` |
| `open user <id>` | `/users/:id` |
| `open workspace <id>` | `/wiki/:workspaceId` |
| `open doc <id>` | `/wiki/docs/:id` |

### 搜索与筛选

| 命令 | 页面/query |
|---|---|
| `search blog "<keyword>"` | `/blog?keyword=...` |
| `search works "<keyword>"` | `/works?keyword=...` |
| `search photos "<keyword>"` | `/photos?keyword=...` |
| `search users "<keyword>"` | `/user-search?keyword=...` |
| `filter blog category <id>` | `/blog?category_id=...` |
| `filter blog tag <id>` | `/blog?tag_id=...` |
| `filter blog rank hot|week|month|year` | `/blog?rank_type=...` |
| `filter works type project|photography` | `/works?type=...` |
| `sort blog time|hot|view_count|like_count|comment_count` | `/blog?sort_by=...` |
| `sort works latest|popular|likes` | `/works?sort=...` |
| `page <number>` | 更新当前可分页页面的 `page` query |
| `reset filters` | 清除当前页面支持的筛选 query |

### 查询命令

| 命令 | API |
|---|---|
| `list articles [keyword]` | `GET /api/articles` |
| `list works [keyword]` | `GET /api/works` |
| `list photos [keyword]` | `GET /api/works?type=photography` |
| `list users <keyword>` | `GET /api/users/search` |
| `list wiki [page]` | `GET /api/wiki/workspaces` |
| `tree workspace <id>` | `GET /api/wiki/workspaces/:id/tree` |

### 互动命令

| 命令 | API 语义 |
|---|---|
| `like article <id>` | 检查状态后按需 `POST /api/articles/:id/like`；后端为 toggle |
| `unlike article <id>` | 检查状态后按需调用现有 toggle，避免反向操作 |
| `like work <id>` | 检查状态后按需 `POST /api/works/:id/like`；后端为 toggle |
| `unlike work <id>` | 检查状态后按需调用现有 toggle |
| `favorite article <id>` | `POST /api/articles/:id/favorite` |
| `unfavorite article <id>` | `DELETE /api/articles/:id/favorite` |
| `favorite work <id>` | `POST /api/works/:id/favorite` |
| `unfavorite work <id>` | `DELETE /api/works/:id/favorite` |
| `follow user <id>` | `POST /api/users/:id/follow` |
| `unfollow user <id>` | `DELETE /api/users/:id/follow` |

### 页面和外观命令

| 命令 | 行为 |
|---|---|
| `theme magazine|terminal` | 预览/应用可用 UI 主题，遵循登录规则 |
| `scheme system|light|dark` | 调整明暗偏好 |
| `scroll top|bottom` | 当前页面滚动 |
| `focus search` | 聚焦当前页面已注册的搜索框 |

## 博客日志等级稳定映射

1. `article.is_top=true` 固定为 `FEAT`。
2. 非置顶文章使用稳定散列，不使用 `Math.random()`：输入依次包含文章 ID、分类 ID/名称以及排序后的标签 ID/名称。
3. 散列结果映射到 `INFO|NOTE|WARN`；同一篇文章在首页、博客列表和刷新后保持一致。
4. 已知语义标签可覆盖散列结果：包含 `观点|思考|随笔|opinion` 映射 `NOTE`；包含 `踩坑|故障|复盘|warning` 映射 `WARN`。
5. `ERROR` 不用于普通文章，避免暗示系统或内容故障；仅终端命令执行失败时作为输出级别。
6. 映射函数放在共享纯工具模块，首页和博客列表不得各自维护一份逻辑。

## 异常处理

| 场景 | 处理方式 |
|---|---|
| 未知命令 | 输出 `command not found` 和最接近的一个建议，不执行任何动作 |
| 引号未闭合或参数缺失 | 输出用法，不猜测用户意图 |
| ID 非正整数 | 拒绝命令，不请求 API |
| 异步命令重复提交 | pending 期间忽略 Enter，保留输入并提示正在执行 |
| 用户未登录执行写命令 | 输出需要登录，提供 `open login`；不调用受保护 API |
| token 过期 | 沿用统一 401 登出行为；终端会话可从 `sessionStorage` 恢复非敏感历史 |
| 写操作等待确认时输入其他命令 | 先取消确认，再解析新命令，并明确输出取消状态 |
| toggle 状态查询失败 | 不执行 like/unlike，避免把目标状态反转 |
| API 业务错误 | 输出 `error.message`，不重复触发第二个 toast |
| 路由不支持某个 query | 命令注册表不生成该组合；实现页面 query 支持后才开放 |
| Wiki 搜索请求 | 明确提示不支持全文搜索，建议 `list wiki` 或 `tree workspace <id>` |
| 拖拽超出视口 | 每次 pointer move 和窗口 resize 后夹紧到可见区域，标题栏始终可访问 |
| 恢复的窗口尺寸大于视口 | 按当前视口重新夹紧；手机直接使用底部抽屉布局 |
| 浏览器不支持 Pointer Events | 保持固定浮窗可用，禁用拖拽，不影响命令输入 |
| 页面切到非 Terminal 主题 | 浮窗继续工作并映射当前主题令牌；不自动执行主题切换 |
| `prefers-reduced-motion` | 禁止原位升起动画、光标闪烁和惯性过渡，直接显示最终位置 |
| 输出过长 | 单条截断并提示；历史最多 200 条输出、50 条命令 |
| 终端结果含用户文本 | 仅文本插值，不使用 `v-html` |

## 技术设计

### 前端模块

```text
web/blog/src/
├── components/terminal/
│   ├── FloatingTerminal.vue
│   ├── TerminalOutput.vue
│   └── TerminalLauncher.vue
├── stores/terminal.js
├── utils/terminal/
│   ├── parser.js
│   ├── commands.js
│   ├── executor.js
│   ├── virtualFileSystem.js
│   └── articleLogLevel.js
└── App.vue
```

职责：

- `App.vue` 只挂载全局终端宿主，不包含命令业务。
- `TerminalHero.vue` 保持首页展示职责，只发出 `activate` 和源矩形，不持有全局历史。
- `terminal` store 保存窗口、会话历史、pending confirmation 和页面 refresh 序号。
- `parser.js` 只处理安全分词和语法，不访问 Vue、DOM、router 或 API。
- `commands.js` 是命令名称、别名、参数、帮助和权限要求的唯一注册表。
- `virtualFileSystem.js` 负责安全路径归一化、路由双向映射、资源文件名生成和 ID 提取；禁止读取真实文件系统。
- `executor.js` 接收显式依赖 `{router, route, api, userStore, appearanceStore, terminalStore}`，不得使用 `eval` 或动态导入任意模块。
- `articleLogLevel.js` 提供稳定纯函数，供 Home/Blog 共用。

### 根级挂载与主题

- `FloatingTerminal` 挂在 `App.vue` 的根 `<RouterView />` 旁，因此跨布局持续存在。
- 未激活时不渲染大型面板，只保留 store 和可选最小恢复按钮。
- 终端主题使用 `--panel|--panel-2|--line|--accent|--green|--amber`；其他主题使用 `--theme-*` 兼容令牌。
- `z-index` 高于 sticky header 和用户中心侧栏，低于 Element Plus modal/message-box 与 Vditor 全屏层。
- holiday/mourning 继续通过根属性覆盖，不另存终端颜色。

### 浮窗与会话状态

store 最小状态：

```js
{
  open: false,
  minimized: false,
  bounds: { x, y, width, height },
  outputs: [],
  commandHistory: [],
  historyCursor: -1,
  pending: false,
  confirmation: null,
  refreshSignals: { article: 0, work: 0, user: 0 },
  sourceRect: null
}
```

- `sessionStorage` 键：`inkspace_terminal_session_v1`。
- 只存 `open|minimized|bounds|outputs|commandHistory`；不存 token、完整用户对象或 pending 写操作。
- 页面刷新后 pending confirmation 一律取消。
- 拖拽使用 Pointer Events、`setPointerCapture()` 和标题栏手柄。
- 缩放使用显式右下角 pointer handle，便于夹紧和测试；不只依赖浏览器 `resize: both`。
- 桌面最小 `420x280`，最大为视口减安全边距；手机 `<=560px` 固定为底部近全屏抽屉。

### 路由 query 统一

需要扩展现有页面：

- `Blog.vue`：读取并监听 `keyword|category_id|tag_id|rank_type|sort_by|page`，表单操作写回 query。
- `Works.vue`：读取并监听 `keyword|type|sort|page`，表单操作写回 query。
- `Photos.vue`：读取并监听 `keyword|sort|page`。
- `WikiIndex.vue`：读取并监听 `page`；不增加不存在的全文搜索参数。
- `UserSearch.vue`：监听 `keyword`，表单搜索写回 query。

约束：

- query 值必须经过白名单与正整数归一化。
- 页面内操作和终端命令生成相同 URL。
- watcher 不得造成重复请求或 `router.push` 循环。
- 终端不直接修改页面组件的内部筛选 ref。

### 页面刷新信号

- 详情页继续拥有文章/作品/用户业务状态。
- 终端写操作成功后增加对应资源 refresh signal，并携带目标 ID。
- `BlogDetail.vue`、`WorkDetail.vue`、`UserProfile.vue` 只在当前目标 ID 匹配时重新调用已有加载/状态检查函数。
- 不建立第二套全局文章、作品或关注状态模型。

### 数据模型

不新增数据库表或字段。会话状态只存浏览器 `sessionStorage`。

### API 接口

不新增终端专用接口，复用以下现有接口：

| 方法 | 路径 | 用途 |
|---|---|---|
| GET | `/api/articles` | 文章列表与关键词查询 |
| GET | `/api/works` | 作品与摄影查询 |
| GET | `/api/users/search` | 用户查询 |
| GET | `/api/wiki/workspaces` | 公开 Wiki 列表 |
| GET | `/api/wiki/workspaces/:id/tree` | 公开 Wiki 树 |
| GET | `/api/articles/:id/is-liked` | 写前检查文章点赞状态 |
| GET | `/api/articles/:id/is-favorited` | 写前检查文章收藏状态 |
| GET | `/api/works/:id/liked` | 写前检查作品点赞状态 |
| GET | `/api/works/:id/favorited` | 写前检查作品收藏状态 |
| GET | `/api/users/:id/follow-stats` | 写前检查关注状态 |
| POST | `/api/articles/:id/like` | toggle 文章点赞 |
| POST | `/api/works/:id/like` | toggle 作品点赞 |
| POST/DELETE | `/api/articles/:id/favorite` | 收藏/取消收藏文章 |
| POST/DELETE | `/api/works/:id/favorite` | 收藏/取消收藏作品 |
| POST/DELETE | `/api/users/:id/follow` | 关注/取消关注用户 |

## 实现步骤（每步可独立 commit）

1. [ ] **稳定日志等级**：新增共享 `articleLogLevel` 纯函数及测试，Home/Blog 统一使用 `FEAT|INFO|NOTE|WARN`。
2. [ ] **命令语法基础**：实现引号感知 parser、命令注册表、帮助文本、参数校验和纯单元测试。
3. [ ] **虚拟文件系统**：实现 `/inkspace` 路径归一化、路由双向映射、内容资源文件和 `ls|cd|cat|grep` 测试。
4. [ ] **会话 store**：实现窗口状态、历史、pending confirmation、sessionStorage 白名单持久化和恢复夹紧。
5. [ ] **全局终端外壳**：在 App 挂载 FloatingTerminal，实现输出、输入、历史、补全、焦点、最小化和关闭。
6. [ ] **Hero 原位升起**：TerminalHero 增加可访问启动入口、源矩形交接、占位状态和 reduced-motion 降级。
7. [ ] **拖拽缩放与移动端**：Pointer Events、视口夹紧、窗口 resize 处理和 `<=560px` 底部抽屉。
8. [ ] **路由命令**：实现 open、pwd、status、scroll、focus，并保持跨路由终端状态。
9. [ ] **页面 query 统一**：让 Blog、Works、Photos、WikiIndex、UserSearch 以 URL 作为筛选状态源。
10. [ ] **只读查询命令**：接入文章、作品、摄影、用户和 Wiki 现有 API，输出安全文本和结果链接。
11. [ ] **主题命令**：接入 appearance store，区分访客会话预览和登录用户保存语义。
12. [ ] **写操作确认**：实现 yes/no 状态机、登录门禁、状态预检、点赞/收藏/关注 API 和防重复提交。
13. [ ] **页面状态刷新**：详情页监听匹配资源的 refresh signal，更新计数和互动状态。
14. [ ] **完整验证**：测试、lint、build、三视口拖拽/抽屉、主题切换、路由后退、过期 token、空错态和无障碍。
15. [ ] **规范同步**：更新 Terminal Phase 2 状态、命令帮助和剩余风险。

## 参考的现有模式

- `web/blog/src/App.vue` — 跨所有布局持续存在的根级挂载点。
- `web/blog/src/components/theme/TerminalHero.vue` — 首页终端视觉、reduced-motion 和计时器清理。
- `web/blog/src/stores/appearance.js` — setup store、会话切换、DOM 主题应用和安全归一化。
- `web/blog/src/stores/workspace.js` — API 型 Pinia store 的现有结构。
- `web/blog/src/router/index.js` — 公共/用户中心路由和登录守卫。
- `web/blog/src/utils/api.js` — JWT、统一响应解包、业务错误与 HTTP 错误处理。
- `web/blog/src/views/Blog.vue` — 文章筛选参数和列表加载。
- `web/blog/src/views/Works.vue`、`Photos.vue` — 作品/摄影查询参数。
- `web/blog/src/views/UserSearch.vue` — 用户搜索 API。
- `web/blog/src/views/BlogDetail.vue` — 文章点赞与收藏状态检查及加载函数。
- `web/blog/src/views/WorkDetail.vue` — 作品点赞、收藏、关注的状态检查与刷新模式。
- `web/blog/src/views/UserProfile.vue` — 用户关注操作和资料刷新。
- `internal/router/blog.go` — 查询与互动 API 的实际路由契约。

## 测试计划

### 纯函数

- [ ] parser 支持空格、单双引号、转义、空参数，并拒绝未闭合引号。
- [ ] 未知命令返回建议，不执行动态代码。
- [ ] ID、页码、主题、明暗和筛选值执行白名单验证。
- [ ] 文章日志等级对同一数据稳定，置顶固定 FEAT，语义标签覆盖正确，首页/列表结果一致。

### Store 和组件

- [ ] 点击 Terminal Hero 只创建一个全局终端，并保持 Hero 占位高度。
- [ ] 跨 MainLayout、Login、UserCenterLayout 路由终端不卸载。
- [ ] sessionStorage 只保存允许字段，刷新后取消 pending confirmation。
- [ ] 打开后聚焦输入；ArrowUp/Down、Tab、Ctrl+L、Escape 行为正确。
- [ ] 拖拽和缩放始终保留可见标题栏，window resize 后重新夹紧。
- [ ] `<=560px` 使用底部抽屉并禁用拖拽缩放。
- [ ] reduced-motion 下无升起、闪烁和位移动画。
- [ ] 输出用户/API 文本时不使用 `v-html`。

### 路由和查询

- [ ] open 命令只生成注册表中的路由，数字 ID 验证正确。
- [ ] Blog 的 keyword/category/tag/rank/sort/page 可由 URL 恢复并响应后退。
- [ ] Works 的 keyword/type/sort/page 可由 URL 恢复。
- [ ] Photos 的 keyword/sort/page 可由 URL 恢复。
- [ ] WikiIndex page 与 UserSearch keyword 可由 URL 恢复。
- [ ] 页面表单操作与终端生成相同 query，不产生 watcher 循环。

### API 和权限

- [ ] list 命令使用正确 endpoint、params 和最大结果数。
- [ ] Wiki 不提供虚假的全文搜索命令。
- [ ] 未登录写命令不请求 API，输出登录提示。
- [ ] 所有写命令必须先 yes/no 确认；no/Escape 不请求 API。
- [ ] article/work like 在 toggle 前检查状态，目标状态已满足时不反向操作。
- [ ] favorite/follow 使用正确 POST/DELETE，不把显式接口当 toggle。
- [ ] pending 阶段不能重复提交。
- [ ] API 错误输出 `error.message`，不额外触发重复 toast。
- [ ] 写成功后只有匹配 ID 的详情页刷新状态。

### 验证命令

- [ ] `web/blog: pnpm test --run`
- [ ] `web/blog: pnpm lint`
- [ ] `web/blog: pnpm build`
- [ ] `go test ./... -count=1`（若未修改 Go，仅作回归）
- [ ] `go vet ./...`
- [ ] `git diff --check`

### 人工矩阵

- [ ] `1440x900`：启动动画、自由拖拽、八方向边界、最小/最大尺寸、跨路由。
- [ ] `900x1024`：浮窗不遮死导航、编辑器或关键操作。
- [ ] `390x844`：底部抽屉、软键盘、输入历史、关闭和安全区。
- [ ] terminal light/dark/system、magazine、holiday、mourning 下可读。
- [ ] 登录/未登录/过期 token，正常/空/错误查询，yes/no 写操作。
- [ ] 键盘、屏幕阅读输出、焦点回收、reduced-motion。

## 待定事项

- `theme terminal|magazine` 对访客是仅当前会话预览，还是允许写访客缓存；当前建议保持既有“访客不新增风格选择入口”语义，仅登录用户可通过终端保存风格。
- 从非首页关闭终端后是否保留一个全站恢复按钮；当前建议关闭即完全关闭，只有最小化才显示恢复条。
- 是否为命令执行提供全局快捷键；当前 MVP 不设置，避免劫持编辑器和页面输入。
- 是否为终端命令结果增加可点击链接；建议支持，但链接必须是注册路由生成的 `router-link`，不能渲染 API 返回 HTML。
- 现有 API 拦截器会先弹 toast 再 reject；终端可输出同一错误文字，但无法阻止全局 toast，后续可为命令请求增加静默选项。

## MVP 范围

1. 首页原位升起、全站持续的单实例终端。
2. 桌面拖拽缩放、手机底部近全屏抽屉、会话级恢复。
3. 类 Unix parser、帮助、历史、补全、清屏、状态和焦点管理。
4. open/search/filter/sort/page 与页面 URL/query 联动。
5. 文章、作品、摄影、用户、公开 Wiki 的现有 API 查询。
6. 主题与明暗命令，遵循现有偏好权限和缓存规则。
7. 点赞、收藏、关注的登录门禁、状态预检、终端 yes/no 确认和详情刷新。
8. `FEAT|INFO|NOTE|WARN` 稳定文章日志等级。
9. terminal/magazine、浅深模式、三视口、键盘、reduced-motion 和安全输出验收。
