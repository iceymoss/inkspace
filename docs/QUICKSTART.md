# 快速开始指南

本指南将帮助你快速搭建 Inkspace 开发环境。

## 📋 前置要求

- Go 1.21+
- Node.js 18+ (推荐使用 pnpm)
- Docker & Docker Compose
- MySQL 8.0+ 和 Redis 7+ (或使用 Docker 容器)

## ⚡ 快速启动

### 1. 克隆项目

```bash
git clone https://github.com/your-username/inkspace.git
cd inkspace
```

### 2. 启动数据库服务

```bash
docker-compose up -d mysql redis
```

等待数据库服务启动完成（约 30 秒）。

### 3. 配置环境变量（可选）

项目支持使用 `.env` 文件或环境变量来配置，环境变量会覆盖 YAML 配置文件中的值。

```bash
# 复制配置模板
cp env.example .env

# 编辑 .env 文件，修改数据库、Redis等配置
# 如果不创建 .env 文件，将使用 config/config.yaml 中的默认配置
```

**配置优先级**：环境变量 > .env 文件 > YAML 配置文件

### 4. 初始化数据库（可选）

数据库迁移会在服务启动时自动执行（通过 GORM AutoMigrate）。

如果需要初始化基础数据（包含默认管理员账号等），执行：

```bash
mysql -h localhost -u inkspace -pinkspace123 inkspace < scripts/init.sql
```

### 5. 启动后端服务

打开三个终端窗口：

**终端1 - 用户服务**:
```bash
go run cmd/server/main.go
# 服务运行在 :8081
```

**终端2 - 管理服务**:
```bash
go run cmd/admin/main.go
# 服务运行在 :8083
```

**终端3 - 定时任务调度器（可选但推荐）**:
```bash
go run cmd/scheduler/main.go
```

### 6. 启动前端服务

**终端4 - 博客前端**:
```bash
cd web/blog
pnpm install  # 首次运行需要安装依赖
pnpm dev
# 前端运行在 :3001
```

**终端5 - 管理前端**:
```bash
cd web/admin
pnpm install  # 首次运行需要安装依赖
pnpm dev
# 前端运行在 :3002
```

## 🌐 访问地址

启动成功后，访问以下地址：

- **博客前端**: http://localhost:3001
- **管理后台**: http://localhost:3002/login
- **用户API**: http://localhost:8081/api
- **管理API**: http://localhost:8083/api

## 🔐 默认账号

### 管理后台
```
地址: http://localhost:3002/login
账号: admin
密码: admin123
```

### 博客系统
```
地址: http://localhost:3001
可以注册新账号
```

## ✅ 验证启动成功

### 检查后端服务

```bash
# 检查用户服务
curl http://localhost:8081/health
# 预期输出: {"status":"ok"}

# 检查管理服务
curl http://localhost:8083/health
# 预期输出: {"status":"ok","service":"admin"}
```

### 检查前端

- 访问 http://localhost:3001 - 应该看到博客首页
- 访问 http://localhost:3002 - 应该跳转到管理登录页面

## 🐛 常见问题

### 端口冲突

如果端口被占用，修改对应的配置文件：

- `config/config.yaml` - `server.port` (用户服务，默认 8081)
- `config/admin.yaml` - `server.port` (管理服务，默认 8083)
- `web/blog/vite.config.js` - `server.port` (博客前端，默认 3001)
- `web/admin/vite.config.js` - `server.port` (管理前端，默认 3002)

### pnpm 未安装

```bash
npm install -g pnpm
```

### 数据库连接失败

1. 检查 Docker 容器是否正常运行：
   ```bash
   docker-compose ps
   ```

2. 检查 `config/config.yaml` 或 `.env` 文件中的数据库配置：
   - 数据库主机和端口
   - 用户名和密码
   - 数据库名称

3. 检查数据库是否已创建：
   ```bash
   docker exec -it inkspace-mysql mysql -u inkspace -pinkspace123 -e "SHOW DATABASES;"
   ```

### Go 模块下载失败

如果在中国大陆，可以配置 Go 代理：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

## 📝 下一步

- 查看 [部署文档](DEPLOYMENT.md) 了解生产环境部署
- 查看 [API 文档](API-REFERENCE.md) 了解接口详情
- 查看 [数据库设计](database-design.md) 了解数据结构

