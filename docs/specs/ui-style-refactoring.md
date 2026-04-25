# UI Style Refactoring Prompt for AI Agent

> 本文档是给 AI Agent 的完整指令，用于在不改变任何交互逻辑的前提下，基于 `ui-ux-pro-max` skill 重构 InkSpace 项目的前端视觉样式。

---

## 0. 绝对红线（Read-Only Zones）

**以下内容严禁修改，任何情况下不得触碰：**

| 类别 | 说明 |
|------|------|
| `<script setup>` / `<script>` 块 | 所有 JS/TS 逻辑、响应式状态、计算属性、方法、watch、生命周期钩子 |
| 事件绑定 | `@click`、`@submit`、`@change`、`@input` 等所有 `@` 事件处理函数及其参数 |
| 条件渲染 | `v-if`、`v-else`、`v-else-if`、`v-show` 的条件和结构 |
| 列表渲染 | `v-for` 的数据源、`key` 绑定 |
| 路由 | `router/` 目录下所有文件、`<router-view>`、`<router-link>` 的 `to` 属性 |
| 状态管理 | Pinia stores（`stores/` 目录）的所有逻辑 |
| API 调用 | `utils/api.js`、`utils/adminApi.js` 及所有 `axios` 调用 |
| 组件 Props & Emits | 组件间通信接口（props 定义、emit 定义） |
| DOM 结构层级 | 元素的父子嵌套关系、`v-for` 包裹的元素层级 |
| Element Plus 组件的选择 | 不得将 `el-button` 换成 `el-link`，不得将 `el-table` 换成 `el-list` 等——只改样式，不改组件类型 |
| 表单验证规则 | `rules` 对象、`required`、`trigger`、`validator` 等 |

**唯一允许修改的是：`<style>` / `<style scoped>` 块、CSS 变量定义、`<template>` 中的 `class` 和 `style` 属性（仅限添加/修改 CSS 类名和内联样式值，不得改变绑定逻辑）。**

---

## 1. 项目上下文

### 项目概况

- **项目**：InkSpace — Go + Vue 3 博客平台
- **两个前端应用**：
  - `web/blog/` — 博客前台（用户面向），约 37 个 Vue 组件
  - `web/admin/` — 管理后台，约 22 个 Vue 组件

### 技术栈

- Vue 3.3 + Composition API + Pinia + Vue Router 4
- Element Plus 2.5（全量引入）
- Vite 5 构建
- **无 Tailwind**，纯 CSS + CSS Custom Properties
- Sass 作为 devDependency 存在但未使用（0 个 `.scss` 文件）

### 当前样式体系

- 设计令牌分散在三处：`main.css` `:root`、`theme.js`（运行时）、`App.vue`（admin）
- Blog 有 4 主题系统（day/night/holiday/mourning），Admin 无主题
- Element Plus 通过全局 CSS `!important` 覆盖，约 30+ 处
- 部分组件使用硬编码颜色值而非 CSS 变量
- 无统一设计语言，视觉风格不一致

### 页面清单

**Blog 前台（20 个路由）：**

| 路由 | 组件 | 说明 |
|------|------|------|
| `/` | Home.vue | 首页 |
| `/blog` | Blog.vue | 文章列表 |
| `/blog/:id` | BlogDetail.vue | 文章详情 |
| `/works` | Works.vue | 作品列表 |
| `/works/:id` | WorkDetail.vue | 作品详情 |
| `/about` | About.vue | 关于页 |
| `/links` | Links.vue | 友链页 |
| `/login` | Login.vue | 登录页 |
| `/user-search` | UserSearch.vue | 用户搜索 |
| `/users/:id` | UserProfile.vue | 用户主页 |
| `/favorites` | Favorites.vue | 收藏夹 |
| `/follow-list` | FollowList.vue | 关注列表 |
| `/dashboard` | user/Dashboard.vue | 用户仪表盘 |
| `/dashboard/articles` | user/MyArticles.vue | 我的文章 |
| `/dashboard/articles/create` | user/ArticleEdit.vue | 创建文章 |
| `/dashboard/articles/:id/edit` | user/ArticleEdit.vue | 编辑文章 |
| `/dashboard/works` | user/MyWorks.vue | 我的作品 |
| `/dashboard/works/create` | user/WorkEdit.vue | 创建作品 |
| `/dashboard/works/:id/edit` | user/WorkEdit.vue | 编辑作品 |
| `/dashboard/notifications` | user/Notifications.vue | 通知 |
| `/dashboard/comments` | user/MyComments.vue | 我的评论 |
| `/profile/edit` | ProfileEdit.vue | 编辑资料 |

**Admin 后台（14 个路由）：**

| 路由 | 组件 | 说明 |
|------|------|------|
| `/login` | admin/Login.vue | 管理员登录 |
| `/` | admin/Dashboard.vue | 仪表盘 |
| `/articles` | admin/Articles.vue | 文章管理 |
| `/articles/create` | admin/ArticleEdit.vue | 新建文章 |
| `/articles/:id/edit` | admin/ArticleEdit.vue | 编辑文章 |
| `/articles/:id` | admin/ArticleView.vue | 查看文章 |
| `/works` | admin/Works.vue | 作品管理 |
| `/categories` | admin/Categories.vue | 分类管理 |
| `/tags` | admin/Tags.vue | 标签管理 |
| `/comments` | admin/Comments.vue | 评论管理 |
| `/links` | admin/Links.vue | 友链管理 |
| `/settings` | admin/Settings.vue | 系统设置 |
| `/users` | admin/Users.vue | 用户管理 |
| `/ads` | admin/Ads.vue | 广告管理 |

---

## 2. 工作流（分阶段执行）

### Phase 1: 设计系统生成（必须首先执行）

**目标**：使用 `ui-ux-pro-max` skill 为 Blog 前台和 Admin 后台分别生成设计系统。

**步骤**：

1. 为 Blog 前台生成设计系统：

```bash
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "blog platform content publishing editorial" --design-system --persist -p "InkSpace Blog"
```

2. 为 Admin 后台生成设计系统：

```bash
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "admin dashboard CMS content management" --design-system --persist -p "InkSpace Admin"
```

3. 补充搜索 Vue 技术栈指南：

```bash
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "component layout form table" --stack vue
```

4. 补充搜索 UX 最佳实践：

```bash
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "animation accessibility loading" --domain ux
```

5. 审查生成的设计系统，确认与项目定位一致（博客平台应为编辑/阅读向风格，非 SaaS 风格）。如果风格推荐不合适，用更精确的关键词重新搜索。

**产出物**：`design-system/MASTER.md`（Blog）和 `design-system/pages/admin.md`（Admin 覆盖规则）

---

### Phase 2: 设计令牌体系重建

**目标**：统一和重构 CSS 自定义属性，建立单一样式真相源。

**步骤**：

1. **备份现有样式文件**：
   - `web/blog/src/assets/main.css`
   - `web/blog/src/utils/theme.js`
   - `web/admin/src/assets/main.css`

2. **在 `web/blog/src/assets/main.css` 中重建 `:root` 变量体系**，根据 Phase 1 的设计系统定义：

```css
:root {
  /* === 颜色系统 === */
  --color-primary: ...;
  --color-primary-hover: ...;
  --color-primary-light: ...;
  --color-secondary: ...;
  --color-accent: ...;
  --color-success: ...;
  --color-warning: ...;
  --color-danger: ...;
  --color-info: ...;

  /* === 文本色阶 === */
  --color-text-primary: ...;
  --color-text-regular: ...;
  --color-text-secondary: ...;
  --color-text-tertiary: ...;
  --color-text-placeholder: ...;
  --color-text-inverse: ...;

  /* === 背景色阶 === */
  --color-bg-primary: ...;
  --color-bg-secondary: ...;
  --color-bg-tertiary: ...;
  --color-bg-card: ...;
  --color-bg-hover: ...;
  --color-bg-active: ...;

  /* === 边框 === */
  --color-border: ...;
  --color-border-light: ...;
  --color-border-lighter: ...;

  /* === 阴影 === */
  --shadow-sm: ...;
  --shadow-md: ...;
  --shadow-lg: ...;

  /* === 圆角 === */
  --radius-sm: ...;
  --radius-md: ...;
  --radius-lg: ...;
  --radius-full: ...;

  /* === 间距 === */
  --spacing-xs: ...;
  --spacing-sm: ...;
  --spacing-md: ...;
  --spacing-lg: ...;
  --spacing-xl: ...;

  /* === 字体 === */
  --font-sans: ...;
  --font-mono: ...;
  --font-size-xs: ...;
  --font-size-sm: ...;
  --font-size-base: ...;
  --font-size-lg: ...;
  --font-size-xl: ...;
  --font-size-2xl: ...;
  --font-size-3xl: ...;
  --line-height-tight: ...;
  --line-height-base: ...;
  --line-height-relaxed: ...;

  /* === 过渡 === */
  --transition-fast: 150ms ease;
  --transition-base: 200ms ease;
  --transition-slow: 300ms ease;
}
```

3. **保留并映射主题系统**：Blog 的 4 主题（day/night/holiday/mourning）必须继续工作。在 `theme.js` 中将运行时变量映射到新的令牌体系。主题变量用 `--theme-*` 前缀，静态设计令牌用 `--color-*` / `--spacing-*` 等前缀。

4. **Admin 端同样建立令牌体系**：在 `web/admin/src/assets/main.css` 中建立 Admin 专属的 `:root` 变量。

5. **Element Plus CSS 变量覆盖**：使用 Element Plus 的 CSS 变量命名空间进行覆盖，替代 `!important` 方式：

```css
:root {
  --el-color-primary: var(--color-primary);
  --el-color-primary-light-3: var(--color-primary-light);
  --el-color-primary-light-5: ...;
  --el-color-primary-light-7: ...;
  --el-color-primary-light-9: ...;
  --el-color-primary-dark-2: ...;
  --el-bg-color: var(--color-bg-primary);
  --el-bg-color-overlay: var(--color-bg-card);
  --el-text-color-primary: var(--color-text-primary);
  --el-text-color-regular: var(--color-text-regular);
  --el-text-color-secondary: var(--color-text-secondary);
  --el-border-color: var(--color-border);
  --el-border-color-light: var(--color-border-light);
  --el-border-color-lighter: var(--color-border-lighter);
  --el-font-size-base: var(--font-size-base);
  --el-border-radius-base: var(--radius-md);
}
```

**产出物**：重构后的 `main.css`（blog + admin）、更新后的 `theme.js`

---

### Phase 3: 全局基础样式重构

**目标**：重写全局 CSS，建立统一的视觉基础。

**步骤**：

1. **重写 `web/blog/src/assets/main.css`** 中的全局样式部分：
   - body 基础样式（字体、颜色、背景）
   - 链接样式（颜色、hover 状态、过渡）
   - 标题层级样式（h1-h6）
   - 滚动条样式
   - 选区样式
   - 全局过渡/动画基础

2. **重写 Element Plus 全局覆盖**：用 CSS 变量方式替代所有 `!important` 覆盖

3. **重写 `web/admin/src/assets/main.css`** 中的全局样式

4. **定义工具类**（按需，不多于当前已有数量）：
   - 间距：`.mt-sm`、`.mb-md` 等
   - 文本：`.text-truncate`、`.text-secondary` 等
   - 布局：`.container`、`.page-section` 等

**产出物**：重构后的全局 CSS 文件

---

### Phase 4: 布局组件重构

**目标**：重构 Layout 组件的视觉样式。

**步骤**（逐文件，按顺序执行）：

| 序号 | 文件 | 重构重点 |
|------|------|----------|
| 1 | `web/blog/src/layouts/MainLayout.vue` | Header、导航、Footer、主体布局间距和视觉 |
| 2 | `web/blog/src/layouts/UserCenterLayout.vue` | 侧边栏样式、内容区域布局 |
| 3 | `web/admin/src/layouts/AdminLayout.vue` | 侧边栏（替换硬编码颜色）、Header、内容区域 |

**每个文件的操作规范**：
- 读取文件 → 理解当前 `<template>` 结构 → 仅修改 `<style scoped>` 和 `class`/`style` 属性
- 硬编码颜色值替换为 CSS 变量引用
- 统一间距、圆角、阴影使用设计令牌
- 确保 hover/active/focus 状态都有平滑过渡
- 确保 4 主题兼容（Blog 布局）

---

### Phase 5: 页面组件逐页重构

**目标**：按优先级逐页重构所有页面和组件的视觉样式。

**执行顺序**（按用户可见度从高到低排序）：

#### 第一批：Blog 高频页面

| 序号 | 文件 | 重构重点 |
|------|------|----------|
| 1 | `web/blog/src/views/Home.vue` | Hero 区域、文章卡片、统计区 |
| 2 | `web/blog/src/views/Blog.vue` | 文章列表、筛选、分页 |
| 3 | `web/blog/src/views/BlogDetail.vue` | 文章排版、代码块、目录、评论区 |
| 4 | `web/blog/src/views/Login.vue` | 登录表单视觉 |

#### 第二批：Blog 内容页面

| 序号 | 文件 | 重构重点 |
|------|------|----------|
| 5 | `web/blog/src/views/Works.vue` | 作品网格布局 |
| 6 | `web/blog/src/views/WorkDetail.vue` | 作品详情排版 |
| 7 | `web/blog/src/views/About.vue` | 关于页排版 |
| 8 | `web/blog/src/views/Links.vue` | 友链卡片 |
| 9 | `web/blog/src/views/UserProfile.vue` | 用户主页布局 |
| 10 | `web/blog/src/views/UserSearch.vue` | 搜索结果布局 |

#### 第三批：Blog 用户中心

| 序号 | 文件 | 重构重点 |
|------|------|----------|
| 11 | `web/blog/src/views/user/Dashboard.vue` | 仪表盘卡片 |
| 12 | `web/blog/src/views/user/MyArticles.vue` | 文章管理列表 |
| 13 | `web/blog/src/views/user/MyWorks.vue` | 作品管理列表 |
| 14 | `web/blog/src/views/user/MyComments.vue` | 评论列表 |
| 15 | `web/blog/src/views/user/Notifications.vue` | 通知列表 |
| 16 | `web/blog/src/views/user/ArticleEdit.vue` | 编辑器布局 |
| 17 | `web/blog/src/views/user/WorkEdit.vue` | 编辑器布局 |
| 18 | `web/blog/src/views/ProfileEdit.vue` | 资料编辑表单 |
| 19 | `web/blog/src/views/Favorites.vue` | 收藏列表 |
| 20 | `web/blog/src/views/FollowList.vue` | 关注列表 |

#### 第四批：公共组件

| 序号 | 文件 | 重构重点 |
|------|------|----------|
| 21 | `web/blog/src/components/NotificationDropdown.vue` | 通知下拉、硬编码颜色替换 |
| 22 | `web/blog/src/components/ImageCropUpload.vue` | 上传区域视觉 |
| 23 | `web/blog/src/components/VditorEditor.vue` | 编辑器主题适配 |
| 24 | `web/blog/src/components/MarkdownEditor.vue` | 编辑器主题适配 |

#### 第五批：Admin 后台

| 序号 | 文件 | 重构重点 |
|------|------|----------|
| 25 | `web/admin/src/views/admin/Login.vue` | 管理员登录页 |
| 26 | `web/admin/src/views/admin/Dashboard.vue` | 仪表盘 |
| 27 | `web/admin/src/views/admin/Articles.vue` | 文章管理表格 |
| 28 | `web/admin/src/views/admin/ArticleEdit.vue` | 文章编辑 |
| 29 | `web/admin/src/views/admin/ArticleView.vue` | 文章查看 |
| 30 | `web/admin/src/views/admin/Works.vue` | 作品管理 |
| 31 | `web/admin/src/views/admin/Categories.vue` | 分类管理 |
| 32 | `web/admin/src/views/admin/Tags.vue` | 标签管理 |
| 33 | `web/admin/src/views/admin/Comments.vue` | 评论管理 |
| 34 | `web/admin/src/views/admin/Links.vue` | 友链管理 |
| 35 | `web/admin/src/views/admin/Settings.vue` | 系统设置 |
| 36 | `web/admin/src/views/admin/Users.vue` | 用户管理 |
| 37 | `web/admin/src/views/admin/Ads.vue` | 广告管理 |
| 38 | `web/admin/src/components/ImageCropUpload.vue` | 上传区域 |
| 39 | `web/admin/src/components/VditorEditor.vue` | 编辑器 |
| 40 | `web/admin/src/components/MarkdownEditor.vue` | 编辑器 |

#### 第六批：Blog 遗留 Admin 页面（低优先级）

| 序号 | 文件 | 说明 |
|------|------|------|
| 41-50 | `web/blog/src/views/admin/*.vue` | 已弃用，仅清理硬编码颜色即可，不做大改 |

**每个文件的操作规范**：

1. **读取**当前文件的 `<template>` 和 `<style>` 部分
2. **分析**哪些样式需要修改（硬编码颜色、不一致间距、缺失过渡等）
3. **仅修改** `<style scoped>` / `<style>` 块和 `<template>` 中的 `class` / `style` 属性
4. **自检**：
   - [ ] 所有硬编码颜色值替换为 CSS 变量
   - [ ] 间距使用设计令牌（`var(--spacing-*)`）
   - [ ] 可交互元素有 `cursor: pointer`
   - [ ] hover/focus 状态有平滑过渡（`transition` 150-300ms）
   - [ ] 卡片/容器使用统一的 `border-radius` 和 `box-shadow`
   - [ ] 文字层级通过 `font-weight` 和 `color` 区分，非仅靠 `font-size`
   - [ ] 图片有 `object-fit` 和 `alt`
   - [ ] Blog 页面兼容 4 主题
   - [ ] 响应式布局未破坏
   - [ ] **未修改任何 `<script>` 内容**
   - [ ] **未修改任何事件绑定**
   - [ ] **未修改任何条件/列表渲染逻辑**
   - [ ] **未修改任何组件 props/emits**

---

### Phase 6: 主题系统对齐

**目标**：确保 Blog 的 4 主题在新令牌体系下正确工作。

**步骤**：

1. 更新 `web/blog/src/utils/theme.js` 中的运行时变量映射
2. 逐一验证 4 主题的视觉效果：
   - `day`：明亮的编辑/阅读体验
   - `night`：护眼的深色模式，确保对比度 ≥ 4.5:1
   - `holiday`：节日氛围（保留但优化配色）
   - `mourning`：灰度模式
3. 确保 `--theme-*` 变量正确覆盖 `--color-*` 基础令牌
4. 确保代码高亮主题（`codeTheme.js`）与新风格协调

---

### Phase 7: 验证与收尾

**目标**：全面验证，确保无回归。

**步骤**：

1. **运行 lint**：
   ```bash
   cd web/blog && pnpm lint
   cd web/admin && pnpm lint
   ```

2. **构建验证**：
   ```bash
   cd web/blog && pnpm build
   cd web/admin && pnpm build
   ```

3. **视觉回归检查清单**：

| 检查项 | Blog | Admin |
|--------|------|-------|
| 首页/仪表盘正常渲染 | | |
| 文章列表和详情正常 | | |
| 表单输入和提交正常 | | |
| 4 主题切换正常（Blog） | | - |
| 响应式布局正常 | | |
| 所有 hover/focus 状态有反馈 | | |
| 无硬编码颜色残留 | | |
| Element Plus 组件样式统一 | | |

4. **清理**：移除不再需要的 `!important` 覆盖、冗余 CSS

---

## 3. 质量守则

### 必须遵守

- **零逻辑变更**：每个文件修改后，`<script>` 内容必须与修改前逐字一致
- **渐进式修改**：一次只改一个文件，改完验证再继续
- **变量优先**：所有颜色、间距、圆角、阴影必须使用 CSS 变量，禁止硬编码
- **过渡平滑**：所有交互状态变化必须有 `transition`（150-300ms）
- **对比度合规**：文字与背景对比度 ≥ 4.5:1（WCAG AA）
- **cursor-pointer**：所有可点击元素必须设置 `cursor: pointer`
- **主题兼容**：Blog 的 4 主题必须全部正常工作
- **无 emoji 图标**：图标使用 Element Plus Icons 或 SVG，不用 emoji

### 禁止事项

- 不得引入新的 npm 依赖
- 不得修改 `package.json` 或 `pnpm-lock.yaml`
- 不得将 `<style scoped>` 改为 `<style lang="scss">`
- 不得添加 Tailwind 或其他 CSS 框架
- 不得修改 `vite.config.js`、`index.html`
- 不得修改 `router/`、`stores/`、`utils/api.js` 中的任何文件
- 不得删除或新增 DOM 元素（仅可修改现有元素的 `class` 和 `style` 属性）
- 不得修改 `v-for`、`v-if`、`v-show` 的条件或结构
- 不得修改 `:class`、`v-bind:style` 的绑定表达式中的逻辑（仅可修改静态 class 名和静态 style 值）

---

## 4. 设计方向参考

根据 InkSpace 作为博客/内容平台的定位，推荐的设计方向：

- **Blog 前台**：Editorial Grid / Magazine 风格 + Swiss Modernism — 排版优先，阅读体验为核心，克制的配色，清晰的层级
- **Admin 后台**：Minimalism & Swiss Style — 信息密度高，操作效率优先，扁平化，少装饰
- **配色**：以中性色（zinc/slate）为主色调，一个主色作为点缀（避免默认蓝 `#409eff`）
- **字体**：正文使用系统字体栈，标题可考虑引入一款衬线或特色字体提升编辑感
- **圆角**：统一使用 `--radius-md`（6-8px），不过圆不过方
- **阴影**：极少使用重阴影，用 `border` + `shadow-sm` 做层级区分

---

## 5. 文件修改记录模板

每修改一个文件，记录如下信息：

```
### [文件路径]
- 修改类型：样式重构
- 修改内容：[简述修改了什么样式]
- 逻辑验证：[确认 <script> 未修改]
- 主题验证：[确认 4 主题正常 / N/A]
- 残留问题：[如有]
```
