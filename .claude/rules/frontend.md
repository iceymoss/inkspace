---
description: Vue 3 前端代码规则
globs: ["web/**/*.{vue,js,jsx}"]
---

- 两个独立应用：`web/blog`（博客，:3001）与 `web/admin`（后台，:3002）。改动前确认改的是哪个应用
- 包管理用 **pnpm**，不要用 npm/yarn/bun
- UI 组件用 **Element Plus**（`element-plus` + `@element-plus/icons-vue`），不要引入其他 UI 库
- Vue 3 Composition API（`<script setup>`）；状态管理用 **Pinia**（`src/stores`），路由用 vue-router（`src/router`）
- HTTP 请求用 **axios**，复用 `src/utils` 里封装的实例/拦截器，不要在组件里裸调 axios
- Markdown 编辑/渲染用 **vditor**，不要另引 Markdown 库
- 界面为**中文单语**，文案直接写中文，不要引入 vue-i18n 或搞多语言 key
- 提交前跑 `pnpm lint`（ESLint --fix）
