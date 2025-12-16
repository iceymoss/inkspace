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
# 数据库迁移会在服务启动时自动执行
# 初始化基础数据（可选）
mysql -h localhost -u inkspace -pinkspace123 inkspace < scripts/init.sql

# 3. 启动后端（2个终端）
go run cmd/server/main.go    # 用户服务 :8081
go run cmd/admin/main.go     # 管理服务 :8083

# 4. 启动前端（2个终端）
cd web/blog && pnpm dev    # 博客前端 :3001
cd web/admin && pnpm dev   # 管理前端 :3002
```

**开发环境访问：**
- **博客**: http://localhost:3001  
- **管理**: http://localhost:3002/login (admin / admin123)

**生产环境访问（Docker 部署）：**
- **博客**: http://is.iceymoss.com  
- **管理**: http://admin.is.iceymoss.com/login (admin / admin123)

**注意：** 生产环境需要先配置 DNS 记录指向服务器 IP，详见 [部署指南](docs/DEPLOYMENT.md)

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

## ⚙️ 配置说明

项目支持多种配置方式，优先级从高到低：

1. **环境变量** - 系统环境变量
2. **.env 文件** - 项目根目录下的 `.env` 文件
3. **YAML 配置文件** - `config/config.yaml` 和 `config/admin.yaml`

### 使用 .env 文件（推荐）

```bash
# 复制配置模板
cp env.example .env

# 编辑 .env 文件，修改数据库、Redis等配置
# 注意：.env 文件不会被提交到 Git（已在 .gitignore 中）
```

### 环境变量列表

主要环境变量：
- `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USERNAME`, `DATABASE_PASSWORD`, `DATABASE_NAME`
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- `JWT_SECRET`, `JWT_ADMIN_SECRET`
- `SERVER_PORT`, `ADMIN_PORT`

完整列表请参考 `env.example` 文件。

---

## 📝 命令

```bash
# 数据库
# 数据库迁移会在服务启动时自动执行（通过 GORM AutoMigrate）
# 初始化基础数据（可选）
mysql -h localhost -u inkspace -pinkspace123 inkspace < scripts/init.sql

# 开发
go run cmd/server/main.go    # 启动用户服务
go run cmd/admin/main.go     # 启动管理服务
go run cmd/scheduler/main.go # 启动定时任务调度器
cd web/blog && pnpm dev       # 启动博客前端
cd web/admin && pnpm dev     # 启动管理前端

# Docker
docker-compose up -d          # 启动
docker-compose down           # 停止
```

---

## 📚 文档

- [QUICKSTART.md](QUICKSTART.md) - 快速启动指南
- [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) - 部署指南
- [docs/database-design.md](docs/database-design.md) - 数据库设计（18张表详解）
- [docs/API-REFERENCE.md](docs/API-REFERENCE.md) - API文档（53个接口）
- [docs/SCHEDULER.md](docs/SCHEDULER.md) - 定时任务调度器文档

---

## 📄 许可证

MIT License

---

参考项目：[WinterChenS/my-site](https://github.com/WinterChenS/my-site)
