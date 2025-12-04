# My Site - Go个人网站

> 基于 Go + Gin + GORM + MySQL + Redis + Vue 3 的现代化个人网站系统

[![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.3-green.svg)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## ✨ 功能特性

### 前台功能
- 🏠 **首页展示** - 展示最新文章和精选作品
- 📝 **博客系统** - 文章列表、详情、分类、标签筛选
- 🎨 **作品展示** - 个人作品集展示和详情
- 💬 **评论系统** - 支持文章评论和回复
- 🔍 **搜索功能** - 文章关键词搜索
- 👤 **关于页面** - 个人信息展示

### 后台管理
- 📊 **数据统计** - 文章、作品、评论等数据概览
- ✍️ **文章管理** - 文章CRUD、Markdown编辑
- 🖼️ **作品管理** - 作品CRUD、图片展示
- 📁 **分类管理** - 分类标签管理
- 🏷️ **标签管理** - 标签CRUD
- 💭 **评论管理** - 评论审核和删除

### 技术特性
- ✅ **RESTful API** - 标准化的API接口设计
- 🔐 **JWT认证** - 基于Token的身份认证
- 🗄️ **Redis缓存** - 提升系统性能
- 📱 **响应式设计** - 完美适配移动端
- 🐳 **Docker部署** - 一键容器化部署
- 🔄 **前后端分离** - 独立开发和部署

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 7+
- Docker & Docker Compose (可选)

### 本地开发

#### 1. 克隆项目

```bash
git clone <repository-url>
cd mysite
```

#### 2. 启动后端

```bash
# 安装依赖
go mod download

# 配置数据库
# 编辑 config/config.yaml 文件，修改数据库配置

# 运行
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

#### 3. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

前端服务将在 `http://localhost:3000` 启动

#### 4. 初始化数据

执行 `scripts/init.sql` 文件初始化数据库和创建默认管理员账号：

- 用户名: `admin`
- 密码: `admin123`

### Docker 部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

访问地址：
- 前端：http://localhost
- 后端：http://localhost:8080
- 管理后台：http://localhost/admin

## 📁 项目结构

```
mysite/
├── cmd/                    # 命令行工具
├── internal/               # 内部代码
│   ├── config/            # 配置管理
│   ├── database/          # 数据库连接
│   ├── handler/           # HTTP处理器
│   ├── middleware/        # 中间件
│   ├── models/            # 数据模型
│   ├── router/            # 路由
│   ├── service/           # 业务逻辑
│   └── utils/             # 工具函数
├── frontend/              # 前端项目
│   ├── src/
│   │   ├── assets/       # 静态资源
│   │   ├── components/   # 组件
│   │   ├── layouts/      # 布局
│   │   ├── router/       # 路由
│   │   ├── stores/       # 状态管理
│   │   ├── utils/        # 工具函数
│   │   └── views/        # 页面
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── config/                # 配置文件
├── scripts/               # 脚本文件
├── uploads/               # 上传文件
├── docker-compose.yml     # Docker编排
├── Dockerfile            # 后端Docker文件
├── Makefile              # Make命令
├── go.mod                # Go依赖
└── main.go               # 入口文件
```

## 🔧 配置说明

### 后端配置 (config/config.yaml)

```yaml
server:
  port: 8080              # 服务端口
  mode: debug             # 运行模式: debug, release, test

database:
  host: localhost         # 数据库地址
  port: 3306             # 数据库端口
  username: root         # 用户名
  password: root         # 密码
  database: mysite       # 数据库名

redis:
  host: localhost        # Redis地址
  port: 6379            # Redis端口
  password: ""          # 密码
  db: 0                 # 数据库编号

jwt:
  secret: your-secret-key-change-this-in-production
  expireHours: 168      # Token过期时间（小时）
```

### 前端配置 (frontend/vite.config.js)

```javascript
export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

## 📖 API文档

### 认证相关

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/register` | POST | 用户注册 |
| `/api/login` | POST | 用户登录 |
| `/api/profile` | GET | 获取个人信息 |

### 文章相关

| 接口 | 方法 | 说明 | 权限 |
|------|------|------|------|
| `/api/articles` | GET | 文章列表 | 公开 |
| `/api/articles/:id` | GET | 文章详情 | 公开 |
| `/api/articles` | POST | 创建文章 | 登录 |
| `/api/articles/:id` | PUT | 更新文章 | 作者/管理员 |
| `/api/articles/:id` | DELETE | 删除文章 | 作者/管理员 |

### 作品相关

| 接口 | 方法 | 说明 | 权限 |
|------|------|------|------|
| `/api/works` | GET | 作品列表 | 公开 |
| `/api/works/:id` | GET | 作品详情 | 公开 |
| `/api/admin/works` | POST | 创建作品 | 管理员 |
| `/api/admin/works/:id` | PUT | 更新作品 | 管理员 |
| `/api/admin/works/:id` | DELETE | 删除作品 | 管理员 |

更多API文档请参考代码中的路由定义。

## 🛠️ 技术栈

### 后端

- **框架**: Gin - 高性能的HTTP Web框架
- **ORM**: GORM - Go语言ORM库
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **认证**: JWT (golang-jwt/jwt)
- **配置**: Viper - 配置管理
- **密码**: bcrypt - 密码加密

### 前端

- **框架**: Vue 3 - 渐进式JavaScript框架
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **UI框架**: Element Plus
- **HTTP客户端**: Axios
- **Markdown**: markdown-it
- **代码高亮**: highlight.js
- **构建工具**: Vite

### 部署

- **容器化**: Docker & Docker Compose
- **Web服务器**: Nginx
- **反向代理**: Nginx

## 📝 开发指南

### 添加新功能

1. 在 `internal/models/` 中定义数据模型
2. 在 `internal/service/` 中实现业务逻辑
3. 在 `internal/handler/` 中实现HTTP处理器
4. 在 `internal/router/` 中注册路由
5. 在前端 `src/views/` 中创建页面

### 数据库迁移

GORM会自动处理数据库迁移，启动时会自动创建或更新表结构。

### 代码规范

- 后端遵循Go语言官方代码规范
- 前端遵循Vue 3官方风格指南
- 使用有意义的变量和函数命名
- 添加必要的注释

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

[MIT License](LICENSE)

## 📮 联系方式

- Email: your.email@example.com
- GitHub: [@yourusername](https://github.com/yourusername)

## 🙏 致谢

本项目参考了 [WinterChenS/my-site](https://github.com/WinterChenS/my-site) 的设计思路。

---

⭐ 如果这个项目对你有帮助，请给个Star支持一下！

