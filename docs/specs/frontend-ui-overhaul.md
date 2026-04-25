# InkSpace 前端 UI 全面美化规格

## Status
- Created: 2026-04-25
- Status: Draft

## Goal

在不改变任何交互逻辑的前提下，将 InkSpace 博客前台和管理后台的 UI 从当前的"功能可用"状态，全面提升为**精致拟物（Refined Skeuomorphic）**风格的专业级视觉体验，同时引入 Tailwind CSS + shadcn-vue 替换当前的 Element Plus + 纯 CSS 方案。

## Non-goals

- **不改变任何交互逻辑**：所有 `<script>` 内容、事件处理、API 调用、路由、状态管理、表单验证规则保持不变
- **不修改后端代码**：不涉及任何 Go 代码变更
- **不改变信息架构**：页面路由结构、导航层级、数据流向不变
- **不清理遗留代码**：blog/views/admin/ 下 11 个遗留文件和 admin/views/user/ 下 3 个遗留文件不做处理
- **不做 SEO 优化**：本次仅关注视觉和样式
- **不做国际化**：不涉及多语言

## User Story

作为 InkSpace 的**博客访客**，我希望看到精致、有质感、像高端杂志一样的阅读界面，让我感受到内容的品质和平台的专业性。

作为 InkSpace 的**管理员**，我希望后台界面精致拟物、层次分明、操作反馈清晰，让我在长时间使用时感到舒适和高效。

## 绝对红线（Read-Only Zones）

**以下内容严禁修改：**

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
| 表单验证规则 | `rules` 对象、`required`、`trigger`、`validator` 等 |
| DOM 结构层级 | 元素的父子嵌套关系、`v-for` 包裹的元素层级不得改变（但可替换为等效的 shadcn-vue 组件标签） |

**允许修改的是：**
- `<template>` 中的组件标签（Element Plus → shadcn-vue 等效替换）
- `<template>` 中的 `class` 和 `style` 属性
- `<style>` / `<style scoped>` 块
- CSS 变量定义
- `main.js` 中的全局注册（移除 Element Plus，添加 shadcn-vue 初始化）

---

## 设计方向

### 视觉风格：精致拟物（Refined Skeuomorphic）

结合 **Dimensional Layering** + **Swiss Modernism 2.0** + 博客 Editorial 特质：

- **层级感**：4 级阴影系统，卡片悬浮感明显，层次分明
- **质感**：细腻的渐变、微妙的光泽、精致的边框处理
- **排版**：编辑/阅读优先，衬线标题 + 无衬线正文，黄金比例间距
- **配色**：Blog 以中性暖灰为底，玫瑰粉为点缀；Admin 以深紫为底，琥珀橙为点缀
- **动效**：所有交互有 150-300ms 平滑过渡，卡片 hover 微浮起

### 技术选型决策

| 当前 | 替换为 | 理由 |
|------|--------|------|
| Element Plus | **shadcn-vue** | 高度可定制、Tailwind 原生、精致拟物风格天然适配、源码可控 |
| 纯 CSS + CSS 变量 | **Tailwind CSS v3** | 原子化一致性强、与 shadcn-vue 配合最佳、开发效率高 |
| 全量引入 Element Plus | **按需引入 shadcn-vue 组件** | 大幅减小 bundle 体积 |
| Element Plus Icons | **Lucide Icons** | 更精致、线条更细腻、与 shadcn-vue 默认搭配 |
| `!important` 覆盖 EP 样式 | **Tailwind + CSS 变量** | 无需 hack，样式可控 |
| `sass` devDependency | 移除 | 不再需要 |

### shadcn-vue 组件映射表

Element Plus → shadcn-vue 等效替换：

| Element Plus | shadcn-vue | 说明 |
|---|---|---|
| `el-button` | `Button` | shadcn-vue 的 Button 组件 |
| `el-input` | `Input` | 文本输入框 |
| `el-textarea` | `Textarea` | 多行输入 |
| `el-select` + `el-option` | `Select` + `SelectTrigger` + `SelectContent` + `SelectItem` | 下拉选择 |
| `el-form` + `el-form-item` | 保留 `<form>` + 自定义样式 | 表单验证逻辑不变，仅替换 UI |
| `el-table` + `el-table-column` | `Table` + `TableHeader` + `TableBody` + `TableRow` + `TableCell` | 数据表格 |
| `el-card` | `Card` + `CardHeader` + `CardContent` + `CardFooter` | 卡片容器 |
| `el-dialog` | `Dialog` + `DialogContent` + `DialogHeader` | 对话框 |
| `el-dropdown` | `DropdownMenu` + `DropdownMenuTrigger` + `DropdownMenuContent` + `DropdownMenuItem` | 下拉菜单 |
| `el-menu` + `el-menu-item` | 自定义侧边栏 + `NavigationMenu` | 侧边导航 |
| `el-tabs` + `el-tab-pane` | `Tabs` + `TabsList` + `TabsTrigger` + `TabsContent` | 选项卡 |
| `el-tag` | `Badge` | 标签/徽标 |
| `el-avatar` | `Avatar` | 头像 |
| `el-pagination` | `Pagination` | 分页 |
| `el-switch` | `Switch` | 开关 |
| `el-badge` | `Badge` | 角标 |
| `el-alert` | `Alert` | 提示 |
| `el-empty` | 自定义空状态组件 | 空状态 |
| `el-divider` | `Separator` | 分隔线 |
| `el-breadcrumb` | `Breadcrumb` + `BreadcrumbItem` | 面包屑 |
| `el-carousel` | 自定义轮播组件 | 轮播（shadcn-vue 无此组件） |
| `el-upload` | 自定义上传区域 | 文件上传 |
| `el-radio-group` + `el-radio` / `el-radio-button` | `RadioGroup` + `RadioGroupItem` | 单选 |
| `el-image` | `<img>` + 自定义样式 | 图片展示 |
| `el-link` | `<a>` / `<router-link>` + 自定义样式 | 链接 |
| `el-row` + `el-col` | Tailwind Grid/Flex | 布局 |
| `el-space` | Tailwind `gap-*` | 间距 |
| `el-collapse` | `Accordion` + `AccordionItem` | 折叠面板 |
| `el-date-picker` | `DatePicker` | 日期选择 |
| `el-input-number` | `Input` + 自定义步进按钮 | 数字输入 |
| `ElMessage` | `toast()` (shadcn-vue Sonner) | 消息提示 |
| `ElMessageBox` | `AlertDialog` | 确认对话框 |
| `ElLoading` | 自定义 Skeleton / Spinner | 加载状态 |

### Blog 前台配色方案

```
:root {
  /* === 主色 === */
  --color-primary: #18181B;        /* 近黑 zinc-900 */
  --color-primary-hover: #27272A;  /* zinc-800 */
  --color-accent: #EC4899;         /* 玫瑰粉 pink-500 */
  --color-accent-hover: #DB2777;   /* pink-600 */

  /* === 功能色 === */
  --color-success: #22C55E;        /* green-500 */
  --color-warning: #F59E0B;        /* amber-500 */
  --color-danger: #EF4444;         /* red-500 */
  --color-info: #6B7280;           /* gray-500 */

  /* === 文本色阶（5级）=== */
  --color-text-primary: #09090B;   /* zinc-950 */
  --color-text-regular: #18181B;   /* zinc-900 */
  --color-text-secondary: #52525B; /* zinc-600 */
  --color-text-tertiary: #71717A;  /* zinc-500 */
  --color-text-placeholder: #A1A1AA; /* zinc-400 */
  --color-text-inverse: #FAFAFA;   /* zinc-50 */

  /* === 背景色阶（5级）=== */
  --color-bg-primary: #FAFAFA;     /* zinc-50 */
  --color-bg-secondary: #F4F4F5;   /* zinc-100 */
  --color-bg-tertiary: #E4E4E7;   /* zinc-200 */
  --color-bg-card: #FFFFFF;        /* white */
  --color-bg-hover: #F4F4F5;       /* zinc-100 */
  --color-bg-active: #E4E4E7;     /* zinc-200 */

  /* === 边框（3级）=== */
  --color-border: #D4D4D8;         /* zinc-300 */
  --color-border-light: #E4E4E7;   /* zinc-200 */
  --color-border-lighter: #F4F4F5; /* zinc-100 */

  /* === 精致拟物阴影（4级）=== */
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.04), 0 1px 3px rgba(0,0,0,0.06);
  --shadow-md: 0 4px 6px -1px rgba(0,0,0,0.06), 0 6px 16px -1px rgba(0,0,0,0.08);
  --shadow-lg: 0 10px 15px -3px rgba(0,0,0,0.08), 0 12px 28px -2px rgba(0,0,0,0.10);
  --shadow-xl: 0 20px 25px -5px rgba(0,0,0,0.10), 0 25px 50px -3px rgba(0,0,0,0.12);

  /* === 圆角 === */
  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-full: 9999px;

  /* === 间距 === */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;
  --spacing-2xl: 48px;

  /* === 字体 === */
  --font-serif: 'Libre Bodoni', Georgia, serif;
  --font-sans: 'Public Sans', system-ui, -apple-system, sans-serif;
  --font-mono: 'Fira Code', ui-monospace, monospace;

  /* === 字号（7级）=== */
  --font-size-xs: 0.75rem;
  --font-size-sm: 0.875rem;
  --font-size-base: 1rem;
  --font-size-lg: 1.125rem;
  --font-size-xl: 1.25rem;
  --font-size-2xl: 1.5rem;
  --font-size-3xl: 2rem;
  --font-size-4xl: 2.5rem;

  /* === 行高 === */
  --line-height-tight: 1.25;
  --line-height-base: 1.6;
  --line-height-relaxed: 1.8;

  /* === 过渡 === */
  --transition-fast: 150ms cubic-bezier(0.4, 0, 0.2, 1);
  --transition-base: 200ms cubic-bezier(0.4, 0, 0.2, 1);
  --transition-slow: 300ms cubic-bezier(0.4, 0, 0.2, 1);
  --transition-spring: 500ms cubic-bezier(0.34, 1.56, 0.64, 1);

  /* === 精致拟物效果 === */
  --card-inset-shadow: inset 0 1px 0 rgba(255,255,255,0.5);
  --card-border: 1px solid var(--color-border-light);
  --hover-lift: -2px;
  --hover-lift-shadow: var(--shadow-lg);
}
```

### Admin 后台配色方案

```
:root {
  --color-primary: #7C3AED;        /* violet-600 */
  --color-primary-hover: #6D28D9;  /* violet-700 */
  --color-accent: #F97316;         /* orange-500 */
  --color-accent-hover: #EA580C;   /* orange-600 */

  --color-text-primary: #1E1B4B;   /* indigo-950 */
  --color-text-regular: #312E81;   /* indigo-800 */
  --color-text-secondary: #6366F1; /* indigo-500 */
  --color-text-tertiary: #818CF8;  /* indigo-400 */
  --color-text-placeholder: #A5B4FC; /* indigo-300 */
  --color-text-inverse: #FAFAFA;

  --color-bg-primary: #FAF5FF;     /* violet-50 */
  --color-bg-secondary: #F5F3FF;   /* violet-100 */
  --color-bg-tertiary: #EDE9FE;   /* violet-200 */
  --color-bg-card: #FFFFFF;
  --color-bg-hover: #F5F3FF;
  --color-bg-active: #EDE9FE;

  --color-border: #DDD6FE;         /* violet-300 */
  --color-border-light: #EDE9FE;   /* violet-200 */
  --color-border-lighter: #F5F3FF; /* violet-100 */

  /* 侧边栏专用 */
  --color-sidebar-bg: #1E1B4B;     /* indigo-950 */
  --color-sidebar-text: #C4B5FD;   /* violet-300 */
  --color-sidebar-text-active: #FFFFFF;
  --color-sidebar-hover: #312E81;  /* indigo-800 */
  --color-sidebar-border: #3730A3; /* indigo-700 */

  /* 阴影、圆角、间距同 Blog */
}
```

### 4 主题运行时变量（Blog 专属）

`--theme-*` 变量在运行时由 `theme.js` 动态设置，覆盖 `--color-*` 静态令牌：

| 变量名 | Day (默认) | Night | Holiday | Mourning |
|--------|-----------|-------|---------|----------|
| `--theme-bg-primary` | #FAFAFA | #0D1117 | #FFF5F5 | #2D2D2D |
| `--theme-bg-secondary` | #F4F4F5 | #161B22 | #FFE4E6 | #3A3A3A |
| `--theme-bg-card` | #FFFFFF | #1C2128 | #FFFFFF | #404040 |
| `--theme-bg-hover` | #F4F4F5 | #21262D | #FFE4E6 | #4A4A4A |
| `--theme-text-primary` | #09090B | #C9D1D9 | #8B1A1A | #D4D4D4 |
| `--theme-text-secondary` | #52525B | #8B949E | #B91C1C | #A3A3A3 |
| `--theme-text-tertiary` | #71717A | #6E7681 | #DC2626 | #737373 |
| `--theme-border` | #D4D4D8 | #30363D | #FECACA | #525252 |
| `--theme-border-light` | #E4E4E7 | #21262D | #FFE4E6 | #404040 |
| `--theme-primary` | #18181B | #58A6FF | #FF3333 | #999999 |
| `--theme-primary-hover` | #27272A | #79C0FF | #FF6666 | #AAAAAA |
| `--theme-accent` | #EC4899 | #58A6FF | #FF3333 | #999999 |
| `--theme-shadow` | rgba(0,0,0,0.08) | rgba(0,0,0,0.3) | rgba(139,26,26,0.1) | rgba(0,0,0,0.2) |
| `--theme-content-bg` | #FFFFFF | #0D1117 | #FFF1F2 | #2D2D2D |
| `--theme-hero-gradient` | linear-gradient(135deg, #18181B, #EC4899) | linear-gradient(135deg, #0D1117, #58A6FF) | linear-gradient(135deg, #8B1A1A, #FF3333) | linear-gradient(135deg, #1A1A1A, #666666) |

---

## Core Flow

### 迁移执行流程

1. **Phase 1**：安装 Tailwind CSS + shadcn-vue 依赖，配置构建工具
2. **Phase 2**：重建设计令牌体系，迁移 CSS 变量到 Tailwind 主题配置
3. **Phase 3**：创建共享 UI 基础组件（AppButton、AppCard 等，封装 shadcn-vue + 项目定制）
4. **Phase 4**：逐页面替换 Element Plus 组件为 shadcn-vue，重写样式
5. **Phase 5**：主题系统适配（确保 4 主题在新体系下正常）
6. **Phase 6**：动效与微交互增强
7. **Phase 7**：暗色模式完善 + 移动端适配
8. **Phase 8**：验证与收尾

---

## Technical Design

### 依赖变更

**新增依赖（web/blog/ 和 web/admin/）：**

```json
{
  "devDependencies": {
    "tailwindcss": "^3.4.0",
    "postcss": "^8.4.0",
    "autoprefixer": "^10.4.0",
    "@tailwindcss/typography": "^0.5.10"
  },
  "dependencies": {
    "radix-vue": "^1.9.0",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.1.0",
    "tailwind-merge": "^2.2.0",
    "lucide-vue-next": "^0.344.0",
    "shadcn-vue": "^0.10.0",
    "vue-sonner": "^1.0.0"
  }
}
```

**移除依赖：**

```json
{
  "dependencies": {
    "element-plus": "REMOVE",
    "@element-plus/icons-vue": "REMOVE"
  },
  "devDependencies": {
    "sass": "REMOVE"
  }
}
```

### 文件结构变更

```
web/blog/src/
├── assets/
│   ├── main.css                    # → 重写：Tailwind 指令 + 自定义样式
│   └── fonts/                      # → 新增：本地字体文件（可选）
├── components/
│   ├── ui/                         # → 新增：shadcn-vue 组件目录
│   │   ├── button.vue
│   │   ├── card.vue
│   │   ├── input.vue
│   │   ├── dialog.vue
│   │   ├── select.vue
│   │   ├── table.vue
│   │   ├── badge.vue
│   │   ├── avatar.vue
│   │   ├── tabs.vue
│   │   ├── dropdown-menu.vue
│   │   ├── pagination.vue
│   │   ├── switch.vue
│   │   ├── alert.vue
│   │   ├── separator.vue
│   │   ├── breadcrumb.vue
│   │   ├── accordion.vue
│   │   ├── skeleton.vue
│   │   ├── toast.vue
│   │   └── ... (按需添加)
│   ├── NotificationDropdown.vue    # → 样式重写 + 图标替换
│   ├── ImageCropUpload.vue         # → 样式重写
│   ├── VditorEditor.vue            # → 样式重写
│   └── MarkdownEditor.vue          # → 样式重写
├── lib/
│   └── utils.ts                    # → 新增：cn() 工具函数（clsx + tailwind-merge）
├── main.js                         # → 修改：移除 Element Plus，添加 Tailwind
├── tailwind.config.js              # → 新增
├── postcss.config.js               # → 新增
└── components.json                 # → 新增：shadcn-vue 配置

web/admin/src/
├── (同上结构)
```

### Tailwind 配置（web/blog/tailwind.config.js）

```js
/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        border: 'var(--color-border)',
        input: 'var(--color-border)',
        ring: 'var(--color-primary)',
        background: 'var(--color-bg-primary)',
        foreground: 'var(--color-text-primary)',
        primary: {
          DEFAULT: 'var(--color-primary)',
          foreground: 'var(--color-text-inverse)',
        },
        accent: {
          DEFAULT: 'var(--color-accent)',
          foreground: 'var(--color-text-inverse)',
        },
        card: {
          DEFAULT: 'var(--color-bg-card)',
          foreground: 'var(--color-text-primary)',
        },
        muted: {
          DEFAULT: 'var(--color-bg-secondary)',
          foreground: 'var(--color-text-secondary)',
        },
        destructive: {
          DEFAULT: 'var(--color-danger)',
          foreground: 'var(--color-text-inverse)',
        },
        sidebar: {
          DEFAULT: 'var(--color-sidebar-bg)',
          foreground: 'var(--color-sidebar-text)',
        },
      },
      fontFamily: {
        serif: ['var(--font-serif)'],
        sans: ['var(--font-sans)'],
        mono: ['var(--font-mono)'],
      },
      borderRadius: {
        sm: 'var(--radius-sm)',
        md: 'var(--radius-md)',
        lg: 'var(--radius-lg)',
        xl: 'var(--radius-xl)',
      },
      boxShadow: {
        sm: 'var(--shadow-sm)',
        md: 'var(--shadow-md)',
        lg: 'var(--shadow-lg)',
        xl: 'var(--shadow-xl)',
      },
      keyframes: {
        'accordion-down': {
          from: { height: '0' },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: '0' },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
      },
    },
  },
  plugins: [require('@tailwindcss/typography')],
}
```

### main.css 新结构（web/blog/src/assets/main.css）

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@import url('https://fonts.googleapis.com/css2?family=Libre+Bodoni:wght@400;500;600;700&family=Public+Sans:wght@300;400;500;600;700&family=Fira+Code:wght@400;500;600;700&display=swap');

@layer base {
  :root {
    /* 所有 CSS 自定义属性（如上面配色方案所示） */
  }

  /* 4 主题运行时覆盖 */
  body.theme-night { ... }
  body.theme-holiday { ... }
  body.theme-mourning { ... }

  * {
    @apply border-border;
  }

  body {
    @apply bg-background text-foreground font-sans;
    font-feature-settings: "rlig" 1, "calt" 1;
  }
}

@layer components {
  /* 精致拟物卡片 */
  .card-skeuomorphic {
    @apply bg-card text-card-foreground rounded-lg border border-border;
    box-shadow: var(--shadow-sm), var(--card-inset-shadow);
    transition: transform var(--transition-base), box-shadow var(--transition-base);
  }

  .card-skeuomorphic:hover {
    transform: translateY(var(--hover-lift));
    box-shadow: var(--hover-lift-shadow), var(--card-inset-shadow);
  }

  /* 容器 */
  .container-blog {
    @apply mx-auto w-full max-w-[1200px] px-[var(--spacing-md)];
  }
}
```

### Element Plus → shadcn-vue 迁移策略

#### 按钮迁移示例

**Before (Element Plus):**
```html
<el-button type="primary" @click="handleSubmit">提交</el-button>
```

**After (shadcn-vue):**
```html
<Button @click="handleSubmit">提交</Button>
```

> 注意：`@click="handleSubmit"` 保持完全不变，仅替换组件标签。

#### 表格迁移示例

**Before (Element Plus):**
```html
<el-table :data="articles">
  <el-table-column prop="title" label="标题" />
  <el-table-column prop="status" label="状态" />
</el-table>
```

**After (shadcn-vue):**
```html
<Table>
  <TableHeader>
    <TableRow>
      <TableHead>标题</TableHead>
      <TableHead>状态</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    <TableRow v-for="article in articles" :key="article.id">
      <TableCell>{{ article.title }}</TableCell>
      <TableCell>{{ article.status }}</TableCell>
    </TableRow>
  </TableBody>
</Table>
```

> 注意：`v-for="article in articles" :key="article.id"` 保持完全不变。

#### 消息提示迁移

**Before:**
```js
import { ElMessage } from 'element-plus'
ElMessage.success('操作成功')
```

**After:**
```js
import { toast } from 'vue-sonner'
toast.success('操作成功')
```

> 注意：这是 `<script>` 块中唯一的例外修改——替换消息提示 API 导入，因为 Element Plus 被移除。功能行为完全等效。

#### 确认对话框迁移

**Before:**
```js
import { ElMessageBox } from 'element-plus'
ElMessageBox.confirm('确定删除？', '提示', { type: 'warning' })
  .then(() => { /* 删除逻辑 */ })
  .catch(() => {})
```

**After:**
```js
// 使用自定义确认函数包装 shadcn-vue AlertDialog
confirmDialog('确定删除？', '提示').then(() => { /* 删除逻辑 */ })
```

> 注意：Promise 链逻辑完全保留，仅替换 UI 呈现方式。

### 图标迁移

Element Plus Icons → Lucide Icons 映射：

| Element Plus | Lucide | 用途 |
|---|---|---|
| `User` | `User` | 用户 |
| `Edit` | `Pencil` | 编辑 |
| `Delete` | `Trash2` | 删除 |
| `Search` | `Search` | 搜索 |
| `Star` | `Star` | 收藏 |
| `View` | `Eye` | 浏览 |
| `Clock` | `Clock` | 时间 |
| `ChatDotRound` | `MessageCircle` | 评论 |
| `Collection` | `Bookmark` | 收藏 |
| `Bell` | `Bell` | 通知 |
| `ArrowDown` | `ChevronDown` | 下拉 |
| `ArrowRight` | `ArrowRight` | 箭头 |
| `Setting` | `Settings` | 设置 |
| `Plus` | `Plus` | 新增 |
| `Lock` | `Lock` | 锁 |
| `Message` | `Mail` | 邮件 |
| `Picture` | `Image` | 图片 |
| `Document` | `FileText` | 文档 |
| `Link` | `Link` | 链接 |
| `Odometer` | `Gauge` | 仪表 |
| `Reading` | `BookOpen` | 阅读 |
| `PriceTag` | `Tag` | 标签 |
| `DataAnalysis` | `BarChart3` | 统计 |
| `RefreshLeft` | `RotateCcw` | 逆时针旋转 |
| `RefreshRight` | `RotateCw` | 顺时针旋转 |
| `Refresh` | `RefreshCw` | 刷新 |
| `FullScreen` | `Maximize` | 全屏 |
| `ZoomIn` | `ZoomIn` | 放大 |
| `ZoomOut` | `ZoomOut` | 缩小 |
| `Close` | `X` | 关闭 |
| `Check` | `Check` | 确认 |
| `SwitchButton` | `LogOut` | 退出 |
| `Calendar` | `Calendar` | 日期 |
| `Location` | `MapPin` | 位置 |
| `Folder` | `Folder` | 文件夹 |
| `ChatLineRound` | `MessageSquare` | 消息 |

---

## Implementation Steps（每步可独立提交）

### Phase 1: 基础设施搭建

1. [ ] 安装 Tailwind CSS + PostCSS + Autoprefixer
2. [ ] 安装 shadcn-vue 及其依赖（radix-vue, class-variance-authority, clsx, tailwind-merge）
3. [ ] 安装 lucide-vue-next 图标库
4. [ ] 安装 vue-sonner 消息提示库
5. [ ] 创建 `tailwind.config.js`（blog + admin 各一份）
6. [ ] 创建 `postcss.config.js`（blog + admin 各一份）
7. [ ] 创建 `src/lib/utils.ts`（cn() 工具函数）
8. [ ] 初始化 shadcn-vue（`components.json`）
9. [ ] 重写 `main.css`（Tailwind 指令 + 设计令牌 + 全局样式）
10. [ ] 更新 `main.js`：移除 Element Plus 注册，添加 Tailwind 指令
11. [ ] 移除 `element-plus` 和 `@element-plus/icons-vue` 依赖
12. [ ] 构建验证：确保基础框架可运行

### Phase 2: 设计令牌与全局样式

13. [ ] 在 `:root` 中定义完整设计令牌（颜色、间距、圆角、阴影、字体）
14. [ ] 定义 4 主题运行时覆盖变量（body.theme-night / .theme-holiday / .theme-mourning）
15. [ ] 定义 Element Plus CSS 变量映射（过渡期保留，逐步移除）
16. [ ] 更新 `theme.js`：适配新的变量命名体系
17. [ ] 创建全局组件样式（card-skeuomorphic、container-blog 等）
18. [ ] 定义 Tailwind 工具类扩展

### Phase 3: shadcn-vue 组件安装

19. [ ] 安装 Button 组件
20. [ ] 安装 Input 组件
21. [ ] 安装 Card 组件
22. [ ] 安装 Dialog 组件
23. [ ] 安装 Table 组件
24. [ ] 安装 Select 组件
25. [ ] 安装 Badge 组件
26. [ ] 安装 Avatar 组件
27. [ ] 安装 Tabs 组件
28. [ ] 安装 DropdownMenu 组件
29. [ ] 安装 Pagination 组件
30. [ ] 安装 Switch 组件
31. [ ] 安装 Alert 组件
32. [ ] 安装 Separator 组件
33. [ ] 安装 Breadcrumb 组件
34. [ ] 安装 Accordion 组件
35. [ ] 安装 Skeleton 组件
36. [ ] 安装 Toast / Sonner 组件
37. [ ] 安装 RadioGroup 组件
38. [ ] 安装 Textarea 组件
39. [ ] 安装 DatePicker 组件
40. [ ] 创建自定义轮播组件
41. [ ] 创建自定义上传组件
42. [ ] 创建自定义空状态组件

### Phase 4: Blog 布局组件迁移

43. [ ] 迁移 `MainLayout.vue`：Header + 导航 + Footer + 主体布局
44. [ ] 迁移 `UserCenterLayout.vue`：侧边栏 + 内容区域

### Phase 5: Blog 高频页面迁移（第一批）

45. [ ] 迁移 `Home.vue`：Hero + 文章卡片 + 统计区
46. [ ] 迁移 `Blog.vue`：文章列表 + 筛选 + 分页
47. [ ] 迁移 `BlogDetail.vue`：文章排版 + 评论区
48. [ ] 迁移 `Login.vue`：登录表单

### Phase 6: Blog 内容页面迁移（第二批）

49. [ ] 迁移 `Works.vue`：作品网格
50. [ ] 迁移 `WorkDetail.vue`：作品详情
51. [ ] 迁移 `About.vue`：关于页
52. [ ] 迁移 `Links.vue`：友链卡片
53. [ ] 迁移 `UserProfile.vue`：用户主页
54. [ ] 迁移 `UserSearch.vue`：搜索结果

### Phase 7: Blog 用户中心迁移（第三批）

55. [ ] 迁移 `user/Dashboard.vue`
56. [ ] 迁移 `user/MyArticles.vue`
57. [ ] 迁移 `user/MyWorks.vue`
58. [ ] 迁移 `user/MyComments.vue`
59. [ ] 迁移 `user/Notifications.vue`
60. [ ] 迁移 `user/ArticleEdit.vue`
61. [ ] 迁移 `user/WorkEdit.vue`
62. [ ] 迁移 `ProfileEdit.vue`
63. [ ] 迁移 `Favorites.vue`
64. [ ] 迁移 `FollowList.vue`

### Phase 8: Blog 公共组件迁移（第四批）

65. [ ] 迁移 `NotificationDropdown.vue`
66. [ ] 迁移 `ImageCropUpload.vue`
67. [ ] 迁移 `VditorEditor.vue`
68. [ ] 迁移 `MarkdownEditor.vue`

### Phase 9: Admin 后台迁移（第五批）

69. [ ] 迁移 `admin/Login.vue`
70. [ ] 迁移 `admin/Dashboard.vue`
71. [ ] 迁移 `admin/Articles.vue`
72. [ ] 迁移 `admin/ArticleEdit.vue`
73. [ ] 迁移 `admin/ArticleView.vue`
74. [ ] 迁移 `admin/Works.vue`
75. [ ] 迁移 `admin/Categories.vue`
76. [ ] 迁移 `admin/Tags.vue`
77. [ ] 迁移 `admin/Comments.vue`
78. [ ] 迁移 `admin/Links.vue`
79. [ ] 迁移 `admin/Settings.vue`
80. [ ] 迁移 `admin/Users.vue`
81. [ ] 迁移 `admin/Ads.vue`
82. [ ] 迁移 Admin `AdminLayout.vue`
83. [ ] 迁移 Admin 公共组件（ImageCropUpload、VditorEditor、MarkdownEditor）

### Phase 10: 主题系统对齐与完善

84. [ ] 更新 `theme.js` 运行时变量映射
85. [ ] 验证 Day 主题
86. [ ] 验证 Night 主题（对比度 ≥ 4.5:1）
87. [ ] 验证 Holiday 主题
88. [ ] 验证 Mourning 主题
89. [ ] 更新 `codeTheme.js` 协调代码高亮风格

### Phase 11: 动效与微交互增强

90. [ ] 为所有可交互元素添加 hover/focus 过渡效果
91. [ ] 添加卡片 hover 浮起动效
92. [ ] 添加页面切换过渡动画
93. [ ] 添加 Skeleton 加载状态
94. [ ] 尊重 `prefers-reduced-motion` 媒体查询

### Phase 12: 移动端适配

95. [ ] Blog 首页响应式优化（375px / 768px / 1024px / 1440px）
96. [ ] Blog 文章列表响应式优化
97. [ ] Blog 文章详情响应式优化
98. [ ] UserCenterLayout 移动端侧边栏适配
99. [ ] Admin 后台响应式优化

### Phase 13: 验证与收尾

100. [ ] 运行 `pnpm lint`（blog + admin）
101. [ ] 运行 `pnpm build`（blog + admin）
102. [ ] 全页面视觉回归验证
103. [ ] 移除残留的 Element Plus 引用和 `!important` 覆盖
104. [ ] 移除 `sass` devDependency
105. [ ] 清理冗余 CSS

---

## Exception Handling

| 场景 | 处理方式 |
|------|----------|
| shadcn-vue 缺少某个 Element Plus 组件 | 自定义实现，基于 radix-vue 原语或原生 HTML + Tailwind |
| 迁移某页面后功能异常 | 回滚该文件，检查是否误改了 `<script>` 或事件绑定 |
| Night 主题对比度不足 | 调整 `--theme-*` 变量值，确保 ≥ 4.5:1 |
| Vditor 编辑器样式与 Tailwind 冲突 | 使用 `:where()` 或前缀隔离，Vditor 自带样式不受 Tailwind preflight 影响 |
| `ElMessage` / `ElMessageBox` 替换后 Promise 行为不一致 | 封装 `confirmDialog()` 函数，确保返回相同接口的 Promise |
| 某组件的 `:class` 动态绑定使用了 Element Plus 特定类名 | 替换为等效的 Tailwind 类名，保持绑定逻辑不变 |
| 构建体积增大 | Tailwind PurgeCSS 自动移除未使用样式；shadcn-vue 按需引入 |

---

## Existing Patterns Referenced

- `web/blog/src/assets/main.css` — 当前设计令牌体系（`--color-*`、`--theme-*`、`--spacing-*` 等）
- `web/blog/src/utils/theme.js` — 4 主题运行时切换机制
- `web/blog/src/utils/codeTheme.js` — 代码高亮主题管理
- `web/blog/src/main.js` — 全量引入 Element Plus 的模式
- `web/blog/src/layouts/MainLayout.vue` — Header/Footer 布局结构
- `web/blog/src/layouts/UserCenterLayout.vue` — 侧边栏布局结构
- `web/admin/src/layouts/AdminLayout.vue` — Admin 侧边栏布局

---

## Test Plan

- [ ] `pnpm build` 成功（blog + admin）
- [ ] `pnpm lint` 无新增错误
- [ ] Blog 所有 20 个路由页面可正常渲染
- [ ] Admin 所有 14 个路由页面可正常渲染
- [ ] Blog 4 主题切换正常
- [ ] 登录/登出流程正常
- [ ] 文章 CRUD 流程正常（创建、编辑、删除）
- [ ] 作品 CRUD 流程正常
- [ ] 评论提交和展示正常
- [ ] 用户资料编辑正常
- [ ] 通知系统正常
- [ ] 收藏/取消收藏正常
- [ ] 图片上传和裁剪正常
- [ ] Markdown 编辑器正常（Vditor + markdown-it）
- [ ] 响应式布局在 375px / 768px / 1024px / 1440px 下正常
- [ ] 无硬编码颜色残留（搜索 `#[0-9a-fA-F]{3,8}` 和 `rgb` 排除变量引用）
- [ ] 所有可交互元素有 `cursor-pointer`
- [ ] 所有 hover/focus 状态有平滑过渡
- [ ] `prefers-reduced-motion` 被尊重

---

## Open Questions

1. **Vditor 兼容性**：Vditor 编辑器自带完整 CSS，与 Tailwind preflight 可能存在冲突，需要在实际迁移时验证并确定隔离策略。
2. **cropperjs 兼容性**：ImageCropUpload 使用 cropperjs，需确认其样式不受 Tailwind reset 影响。
3. **Admin 主题**：当前 Admin 无主题切换，是否需要新增暗色模式？（建议：本次不做，后续迭代）
4. **博客后台遗留页面**：`web/blog/src/views/admin/` 下 11 个遗留文件不在路由中，是否迁移？（根据用户回答：跳过）
5. **ElLoading 替换**：`utils/workNavigation.js` 使用 `ElLoading.service()`，需确认替换方案（建议：自定义 Spinner 或 Skeleton）。

---

## MVP Scope

### MVP（必须完成）
- Phase 1: 基础设施搭建
- Phase 2: 设计令牌与全局样式
- Phase 3: shadcn-vue 组件安装
- Phase 4: Blog 布局组件迁移
- Phase 5: Blog 高频页面迁移（Home、Blog、BlogDetail、Login）
- Phase 10: 主题系统对齐（至少 Day + Night 主题）
- Phase 13: 构建验证

### V2（第二批）
- Phase 6: Blog 内容页面
- Phase 7: Blog 用户中心
- Phase 8: Blog 公共组件

### V3（第三批）
- Phase 9: Admin 后台迁移
- Phase 11: 动效增强
- Phase 12: 移动端适配

---

## 文件修改记录模板

每修改一个文件，记录如下信息：

```
### [文件路径]
- 修改类型：[组件迁移 / 样式重写 / 配置变更]
- 修改内容：[简述修改了什么]
- 逻辑验证：[确认 <script> 未修改 / 仅替换消息提示 API]
- 主题验证：[确认 4 主题正常 / N/A]
- 残留问题：[如有]
```
