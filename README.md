# Inkspace - 现代化个人博客系统

<div align="center">

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Vue](https://img.shields.io/badge/Vue-3.3-4FC08D?style=for-the-badge&logo=vue.js)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Production_Ready-success?style=for-the-badge)

**基于 Go + Vue 3 构建的现代化多用户博客系统**

[功能特性](#-功能特性) • [快速开始](#-快速开始) • [技术栈](#-技术栈) • [文档](#-文档) • [许可证](#-许可证)

</div>

---

## ✨ 功能特性

### 核心功能
- ✅ **用户系统** - 注册登录、个人主页、用户关注/粉丝系统、个人资料管理
- ✅ **内容管理** - Markdown 编辑器、文章发布编辑、分类标签管理、作品展示（开源项目/摄影作品）
- ✅ **社交互动** - 评论系统（支持回复）、点赞、收藏、实时通知、用户关注
- ✅ **内容发现** - 热门文章排名、推荐文章/作品、分类浏览、标签筛选、搜索功能
- ✅ **作品展示** - 支持开源项目和摄影作品两种类型，摄影作品支持相册管理和EXIF信息
- ✅ **扩展功能** - 友情链接管理、文件上传/附件管理、访问统计
- ✅ **管理后台** - 完整的后台管理系统，包括：
  - 内容管理：文章、作品、分类、标签、评论审核
  - 用户管理：用户列表、权限管理、状态控制
  - 系统配置：首页轮播图、系统参数设置、主题风格
  - 广告管理：广告位管理、广告内容管理、广告投放
  - 友链管理：友情链接的增删改查
- ✅ **定时任务** - 独立的调度器服务，自动处理热门文章统计、数据更新等后台任务

### 为什么选择 Inkspace
- 🎯 **开箱即用** - 完整的博客系统，无需从零开始搭建
- 🚀 **快速部署** - Docker Compose 一键启动，几分钟即可上线
- 🎨 **现代化 UI** - 基于 Vue 3 和 Element Plus，界面美观易用
- 🔧 **易于扩展** - 清晰的代码结构，方便二次开发和定制
- 📱 **功能完整** - 从内容管理到社交互动，满足个人博客的所有需求

---

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+ (推荐使用 pnpm)
- Docker & Docker Compose (用于数据库)
- MySQL 8.0+ 和 Redis 7+ (或使用 Docker)

### 开发环境启动

```bash
# 1. 克隆项目
git clone https://github.com/your-username/inkspace.git
cd inkspace

# 2. 启动数据库服务
docker-compose up -d mysql redis

# 3. 配置环境变量（可选）
cp env.example .env
# 编辑 .env 文件修改数据库配置

# 4. 初始化数据库（可选，包含默认管理员账号）
mysql -h localhost -u inkspace -pinkspace123 inkspace < scripts/init.sql

# 5. 启动后端服务
go run cmd/server/main.go    # 用户服务 :8081
go run cmd/admin/main.go     # 管理服务 :8083
go run cmd/scheduler/main.go # 定时任务调度器（可选）

# 6. 启动前端（新终端）
cd web/blog && pnpm install && pnpm dev   # 博客前端 :3001
cd web/admin && pnpm install && pnpm dev  # 管理前端 :3002
```

**访问地址：**
- 博客前端: http://localhost:3001
- 管理后台: http://localhost:3002/login (admin / admin123)

### Docker Compose 一键部署

使用 Docker Compose 可以快速启动所有服务，适合生产环境部署：

```bash
# 方式一：完整部署（包含 MySQL 和 Redis）
docker-compose up -d

# 方式二：使用外部数据库服务
docker-compose -f docker-compose.external-db.yml up -d
```

**部署后访问：**
- 博客前端: http://is.iceymoss.com（需配置 DNS）
- 管理后台: http://admin.is.iceymoss.com（需配置 DNS）

**默认账号：**
- 管理后台: admin / admin123

> 💡 详细部署步骤、DNS 配置、HTTPS 配置等请查看 [部署文档](docs/DEPLOYMENT.md)

---

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin (HTTP 路由)
- **ORM**: GORM (数据库操作)
- **数据库**: MySQL 8.0+
- **缓存**: Redis 7+
- **认证**: JWT

### 前端
- **框架**: Vue 3 (Composition API)
- **UI 库**: Element Plus
- **状态管理**: Pinia
- **构建工具**: Vite

### 部署
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **负载均衡**: Nginx Upstream

---

## 📊 项目规模

| 类型 | 数量  | 说明 |
|------|-----|------|
| 数据库表 | 21张 | 完整的关系型数据库设计 |
| API 接口 | 53个 | RESTful 风格 API |
| 前端页面 | 21个 | Vue 3 + Element Plus |
| 服务模块 | 12个 | Service 层业务逻辑 |

---

## 📚 文档

- 📖 [快速开始指南](docs/QUICKSTART.md) - 详细的开发环境搭建步骤
- 🚀 [部署文档](docs/DEPLOYMENT.md) - 生产环境部署指南（Docker Compose）
- 🗄️ [数据库设计](docs/database-design.md) - 数据库表结构设计说明
- 🔌 [API 参考](docs/API-REFERENCE.md) - 完整的 API 接口文档
- ⏰ [定时任务](docs/SCHEDULER.md) - 调度器服务说明

---

## ⚙️ 配置说明

项目支持多种配置方式，优先级从高到低：

1. **环境变量** - 系统环境变量
2. **.env 文件** - 项目根目录下的 `.env` 文件
3. **YAML 配置文件** - `config/config.yaml` 和 `config/admin.yaml`

主要配置项：
- 数据库连接（MySQL）
- Redis 连接
- JWT 密钥
- 服务端口

完整配置说明请参考 `env.example` 文件。

---

## 🏗️ 项目结构

```
inkspace/
├── cmd/                    # 服务入口
│   ├── server/            # 用户服务 (8081)
│   ├── admin/             # 管理服务 (8083)
│   └── scheduler/         # 定时任务调度器
├── internal/              # 内部代码
│   ├── config/            # 配置管理
│   ├── database/          # 数据库连接和迁移
│   ├── handler/           # HTTP 处理器
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── router/            # 路由定义
│   └── service/           # 业务逻辑层
├── web/                   # 前端项目
│   ├── blog/             # 博客前端
│   └── admin/            # 管理后台前端
├── config/                # 配置文件
├── scripts/               # 脚本文件
├── nginx/                 # Nginx 配置
└── docs/                  # 项目文档
```

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

---


<div align="center">

**如果这个项目对你有帮助，请给一个 ⭐ Star！**

Made with ❤️ by Inkspace Team

</div>
