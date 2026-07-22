# 风格 B · 暗色科技感「yu.log」设计规范

> 原型：`style-b-terminal.html` · 默认深色模式
> 关键词：终端、等宽、深蓝黑、冰蓝、日志

## 1. 设计原则

这一版把个人站做成一台「正在运行的机器」：Hero 是一个活的终端，博客是日志流（`FEAT/INFO/WARN`），项目是仓库卡片，知识库是文件树。所有隐喻都来自开发者的真实工作环境，而非泛泛的「科技感」贴皮。深色底选深蓝黑而非纯黑，长时间阅读更舒适；等宽字体只用于「机器语言」出现的地方（路径、命令、时间戳、EXIF），正文仍是无衬线中文——等宽是点缀，不是基调。光效克制：只有状态点有 glow，没有大面积霓虹。

## 2. 色彩令牌

| Token | 深色模式（默认） | 浅色模式 | 用途 |
|---|---|---|---|
| `--bg` | `#0D1320` | `#F4F6FA` | 页面背景 |
| `--panel` | `#131B2C` | `#FFFFFF` | 卡片/终端/面板 |
| `--panel-2` | `#182238` | `#EBEFF6` | 面板内高亮、悬停底 |
| `--ink` | `#C9D5E6` | `#23304A` | 正文 |
| `--bright` | `#EDF3FB` | `#0E1730` | 标题、命令高亮 |
| `--sub` | `#64748F` | `#7C88A1` | 次要文字、输出文本 |
| `--line` | `#1F2B42` | `#DDE4EF` | 描边、分隔线、背景网格 |
| `--accent` | `#5CB2FF` | `#1668C7` | 主强调（冰蓝）：链接、光标、选中 |
| `--amber` | `#F0A94B` | `#C07A17` | 次强调（琥珀）：数字、WARN、code |
| `--green` | `#6FCF8E` | `#1F9D57` | 终端提示符、在线状态、INFO |
| `--shadow` | `0 12px 40px rgba(0,0,0,.45)` | `0 10px 30px rgba(30,50,90,.10)` | 面板阴影 |

语义色规则：绿 = 系统正常/提示符；蓝 = 交互/链接；琥珀 = 数据/警示。三色不混用场景。

## 3. 字体系统

| 角色 | 字体栈 | 使用位置 |
|---|---|---|
| 等宽 `.mono` | `'SF Mono', 'JetBrains Mono', 'Cascadia Code', Consolas, 'Courier New', monospace` | 终端、导航、时间戳、标签、路径、EXIF、页脚 |
| 正文 | `'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif` | 标题、正文、描述 |

字号阶梯：

| 层级 | 字号 | 行高 | 字重 |
|---|---|---|---|
| H1 Hero | `clamp(34px, 5vw, 52px)` | 1.3 | 700 |
| 栏目标题 | 15px（等宽，路径式 `~/blog`） | — | 600 |
| H3 条目 | 16.5px | — | 600 |
| 正文 | 15px | 1.75 | 400 |
| 终端文本 | 13.5px | 2.05 | 400 |
| 辅助/标签 | 12–13px | — | 400 |

栏目标题采用路径写法：`~` 灰、`/` 蓝、栏目名亮色、后缀说明灰——这是本风格的标题语法。

## 4. 间距、圆角与层次

| Token | 值 |
|---|---|
| 内容最大宽度 | 1100px，两侧 28px |
| 区块垂直 padding | 68px（Hero 84px/72px） |
| 圆角 | 面板/终端 14px · 卡片 14px · 按钮 10px · 小元素 8px |
| 网格 gap | 18px |
| 背景装饰 | 44px 网格线（`--line`，opacity .25），径向 mask 只在页面上半部显示 |

层次由三级底色（bg → panel → panel-2）+ 描边表达，阴影只给悬浮态与终端/知识库大面板。

## 5. 组件规范

**导航**：sticky，高 60px，82% 透明 blur(12px)，底部 1px `--line`。品牌为等宽 `$yu.log` + 闪烁光标块（8×16px，`steps(1)` 1.1s）。链接小写英文，12px→13px，圆角 8px，悬停 `--panel-2` 底；激活项前缀 `// `（蓝色）。

**Hero（左文右终端）**：左侧：在线状态胶囊（绿点带 glow + `available · city (UTC+8)`）→ H1（关键词蓝色）→ 简介 → 双按钮（实心蓝 `tail -f blog` / 描边 `ls projects/`，按钮文案用命令语）。右侧终端窗口见 §6。

**日志列表（博客）**：grid `150px 1fr auto`。时间戳前置等级标签：`FEAT` 蓝 / `INFO` 绿 / `WARN` 琥珀（映射文章类型：新作/常规/观点）。整行悬停浮出 `--panel` 底 + 描边，圆角 12px。

**仓库卡片（项目）**：3 列，`--panel` 底 + 描边 + 14px 圆角。结构：`❯ name`（蓝）→ 描述 → meta 行（语言色点 + star 数 + 版本号）。语言点色：TS `#4E7CBF`、Rust `#C0714F`、Vue `#41B883`。悬停上移 4px + 蓝描边 + 阴影。

**摄影（captures）**：3 列 `aspect-ratio: 3/2`，14px 圆角 + 描边。悬停 `scale(1.05)`，底部浮出 EXIF 行（左标题右参数，11px 等宽）。叠加 3px 扫描线层（`repeating-linear-gradient`，白 2.5% 透明度）。

**知识库（文件树）**：大面板 grid `280px 1fr`。左侧目录树：`▸/▾` 蓝色折叠符，文件项 `· ` 前缀，选中项蓝字 + 10% 蓝底。右侧预览：面包屑（`/` 蓝色分隔）→ 标题 → 正文（`<code>` 为琥珀色 panel-2 底胶囊）→ 底部虚线分隔的统计行（数字蓝色 20px）。

**页脚**：1px 顶线，等宽 13px，含 `■ all systems normal`（绿色方块）状态语。

## 6. 签名元素:终端 Hero

结构：标题栏（红黄绿三点 11px + 居中 `chenyu@island: ~`）→ 正文区。内容脚本：`whoami` → 身份一行；`cat skills.json | jq .top` → 蓝色字符串数组；`uptime` → 琥珀数字年限；最后一行为打字循环。

打字动画参数：逐字 70–160ms 随机 → 停留 1600ms → 逐字删除 30ms → 停 500ms 换下一条。命令池：`ls photos/2026/`、`npm run build`、`git push origin main`、`vim wiki/today.md`。`prefers-reduced-motion` 时静态显示第一条，光标不闪。

## 7. 动效

| 场景 | 参数 |
|---|---|
| 主题切换 | `.3s` |
| 卡片悬停 | `.25s`，上移 4px + 描边变色 |
| 图片缩放 | `.5s cubic-bezier(.2,.6,.2,1)` |
| 光标闪烁 | `1.1s steps(1) infinite` |
| 按钮悬停 | 实心：`brightness(1.12)` + 上移 1px |

## 8. 响应式

| 断点 | 变化 |
|---|---|
| ≤900px | Hero 单列（终端在下）；卡片/照片 2 列；文件树隐藏只留预览；导航收起；日志单列 |
| ≤560px | 卡片/照片 1 列 |

## 9. 无障碍

焦点：`outline: 2px solid var(--accent)` + 4px 圆角。深色模式正文 `#C9D5E6` 对比度 ≈ 11:1；`--sub` `#64748F` ≈ 4.6:1，仅用于辅助文字。等级标签（FEAT/INFO/WARN）同时有文字，不单靠颜色传义。终端为装饰性内容，加 `aria-label`；打字动画尊重 reduced-motion。

## 10. Vue 落地提示

`TerminalHero.vue`（props: lines, typingCommands；内部 `useTyping()` composable）、`LogRow.vue`（level → 颜色映射表）、`RepoCard.vue`（lang → 色点映射）、`WikiTree.vue`（递归树组件，选中态受控）。背景网格封装为 `BgGrid.vue` 挂在布局层。等级/语言的颜色映射集中在一个 `tokens.ts` 导出，避免散落。
