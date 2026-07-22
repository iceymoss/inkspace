# 个人站点 UI 设计规范 · 总览

> 项目：博客 + 作品（项目/摄影）+ 知识库 的个人站点
> 四个风格方向的设计规范，与 HTML 高保真原型一一对应。
> 版本：v1.0 · 2026-07

## 文档结构

| 文档 | 风格 | 原型文件 | 一句话定位 |
|---|---|---|---|
| [style-a-magazine.md](style-a-magazine.md) | A · 极简杂志风「屿刊」 | `style-a-magazine.html` | 纸感留白，衬线排版，为长文与摄影让路 |
| [style-b-terminal.md](style-b-terminal.md) | B · 暗色科技感「yu.log」 | `style-b-terminal.html` | 终端美学，等宽点缀，开发者身份表达 |
| [style-c-cozy.md](style-c-cozy.md) | C · 温暖手作感「小屿的角落」 | `style-c-cozy.html` | 暖纸色与拍立得，亲切的生活气息 |
| [style-d-swiss.md](style-d-swiss.md) | D · 瑞士网格风「CHEN YU®」 | `style-d-swiss.html` | 外露网格与克莱因蓝，理性国际范 |

每份规范包含：设计原则、色彩令牌（明/暗双模式）、字体系统、间距与栅格、组件规范、签名元素、动效参数、响应式断点、无障碍要点、Vue 落地建议。

## 四风格速览对比

| 维度 | A 杂志风 | B 科技感 | C 手作感 | D 瑞士风 |
|---|---|---|---|---|
| 默认模式 | 浅色 | 深色 | 浅色 | 浅色 |
| 主色调 | 纸白/墨黑/黛绿 | 深蓝黑/冰蓝/琥珀 | 暖纸/苔绿/芥黄 | 纯白/纯黑/克莱因蓝 |
| 标题字 | 衬线（宋体系） | 无衬线 + 等宽 | 无衬线加粗 | 无衬线超粗大写 |
| 圆角 | 0（仅标签胶囊） | 8–14px | 16–20px + 不规则 | 0 |
| 边框 | 1px 发丝线 | 1px 面板描边 | 2px 描边 + 虚线 | 1px 线 + 主结构黑线 |
| 阴影 | 无 | 大而深 | 软而暖 | 无 |
| 签名元素 | 竖排文字 + 期刊号 | 打字终端 Hero | 胶带拍立得墙 | 网格坐标编号系统 |
| 情绪关键词 | 安静、克制 | 极客、精密 | 温暖、亲切 | 理性、有力 |

## 共用内容模型

四版原型使用同一套内容，便于纯粹对比风格。落地时的信息架构：

```
站点
├── 首页          Hero + 各板块精选
├── 文章 /blog    列表页 → 详情页（正文、目录、上一篇/下一篇）
├── 作品 /works   项目卡片（名称、角色、简介、技术栈、链接）
├── 摄影 /photos  相册列表 → 相册详情（大图浏览、EXIF 可选）
└── 知识库 /wiki  分类入口 → 笔记页（双向链接、更新时间）
```

核心实体字段：

| 实体 | 字段 |
|---|---|
| 文章 Post | title, excerpt, date, tag, readingTime, cover? |
| 项目 Project | name, role(开源/模板…), description, lang, stars?, link |
| 照片 Photo | title, location, year, exif?(focal/aperture/iso) |
| 笔记分类 WikiCategory | name, description, noteCount, updatedAt |

## 全局基线（四版通用）

**明暗模式**：以 `<html data-theme="light|dark">` 切换，所有颜色经 CSS 自定义属性引用，禁止组件内写死色值。首次加载跟随 `prefers-color-scheme`，用户手动切换后以用户选择为准（落地时可持久化到 localStorage，原型未做持久化）。

**响应式**：桌面优先编写，断点统一为 `900px`（平板）与 `560px`（手机）。导航在窄屏收起为汉堡菜单（原型简化为隐藏链接，落地需补齐）。

**无障碍**：所有交互元素提供 `:focus-visible` 样式（描边色用各风格主强调色）；`prefers-reduced-motion: reduce` 时关闭全部过渡与动画；图片类内容提供替代文本；颜色对比度正文不低于 4.5:1，大字号不低于 3:1。

**性能**：字体全部使用系统字体栈（零字体加载成本）；照片占位在原型中为 CSS 渐变，落地时替换为 `<img loading="lazy">` + `aspect-ratio` 占位防抖动。

## 落地技术建议

原型为纯 HTML/CSS/原生 JS。推荐落地栈：Vue 3 + Vite（或 Nuxt 3 做 SSG），样式可选 UnoCSS/Tailwind（按规范映射 token）或直接沿用 CSS 变量方案。内容层建议 Markdown + frontmatter（博客/知识库），照片元数据用 JSON/YAML 清单。

组件拆分建议（四风格通用骨架，仅样式层不同）：

```
components/
├── SiteHeader.vue      导航 + 主题切换
├── SiteFooter.vue
├── ThemeToggle.vue
├── home/HeroSection.vue
├── blog/PostList.vue · PostCard.vue
├── works/ProjectCard.vue
├── photos/PhotoGrid.vue · PhotoItem.vue
└── wiki/WikiCategoryCard.vue
```
