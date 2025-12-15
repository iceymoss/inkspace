# My Site - Go个人博客系统

![Go](https://img.shields.io/badge/Go-1.21-blue.svg)
![Vue](https://img.shields.io/badge/Vue-3.3-green.svg)
![Status](https://img.shields.io/badge/Status-Production_Ready-success.svg)

基于 **Go + Gin + GORM + MySQL + Redis + Vue 3** 的多用户博客系统

---

## ⚡ 快速启动

### 四服务架构

```bash
# 1. 启动数据库
docker-compose up -d mysql redis

# 2. 初始化数据库
make db-migrate && make db-init

# 3. 启动后端（2个终端）
make dev          # 用户服务 :8081
make dev-admin    # 管理服务 :8083

# 4. 启动前端（2个终端）
cd frontend/blog && pnpm dev    # 博客前端 :3001
cd frontend/admin && pnpm dev   # 管理前端 :3002
```

**博客**: http://localhost:3001  
**管理**: http://localhost:3002/login (admin / admin123)

详细步骤：[QUICKSTART.md](QUICKSTART.md)

---

## ✨ 功能特性

- ✅ 用户注册登录、个人主页
- ✅ 文章管理（Markdown编辑）
- ✅ 评论系统（树形结构）
- ✅ **用户关注/粉丝系统** 🔥
- ✅ **文章收藏功能** 🔥
- ✅ **点赞系统**（文章+评论）🔥
- ✅ **实时通知** 🔥
- ✅ **热门文章排名**（多维度计算）🔥
- ✅ **推荐文章/作品**（管理后台设置）🔥
- ✅ 分类标签、作品展示、友情链接
- ✅ 完整的管理后台
- ✅ 独立的定时任务调度器

---

## 🚀 技术亮点

### 数据库设计
- **18张表**完整设计
- **13个冗余计数字段**，自动维护，减少80% JOIN查询
- **55+个索引**，优化查询性能
- **事务保证**数据一致性

### 技术栈
**后端**：Go 1.21 + Gin + GORM + MySQL 8.0 + Redis 7  
**前端**：Vue 3 + Element Plus + Pinia + Vite  
**部署**：Docker + Docker Compose

### 项目规模
- 数据库表：18张
- API接口：53个
- 前端页面：21个
- Service层：12个

---

## 📝 命令

```bash
# 数据库
make db-migrate    # 创建18张表
make db-init       # 导入初始数据
make db-sync       # 同步计数器

# 开发（WSL终端）
make dev           # 启动后端
cd frontend && npm run dev  # 启动前端

# Docker
make docker-up     # 启动
make docker-down   # 停止
```

---

## 📚 文档

- [QUICKSTART.md](QUICKSTART.md) - 快速启动指南
- [DEPLOYMENT.md](DEPLOYMENT.md) - 部署指南
- [docs/database-design.md](docs/database-design.md) - 数据库设计（18张表详解）
- [docs/API-REFERENCE.md](docs/API-REFERENCE.md) - API文档（53个接口）
- [docs/SCHEDULER.md](docs/SCHEDULER.md) - 定时任务调度器文档

---

## 📄 许可证

MIT License

---

参考项目：[WinterChenS/my-site](https://github.com/WinterChenS/my-site)
