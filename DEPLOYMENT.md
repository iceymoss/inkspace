# 部署指南

## 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 7+
- Docker & Docker Compose (可选)

## 快速部署

### 1. 克隆项目

```bash
git clone <repository-url>
cd inkspace
```

### 2. 配置环境

复制配置文件并修改：

```bash
cp config/config.yaml.example config/config.yaml
cp config/admin.yaml.example config/admin.yaml
```

编辑配置文件，设置数据库和Redis连接信息。

### 3. 启动数据库

使用 Docker Compose：

```bash
docker-compose up -d mysql redis
```

或手动启动 MySQL 和 Redis 服务。

### 4. 初始化数据库

```bash
# 创建数据库（如果不存在）
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS mysite CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行数据库迁移（GORM会自动创建表结构）
make db-migrate

# 初始化基础数据（可选）
make db-init
```

### 5. 启动后端服务

```bash
# 用户服务（端口8081）
make dev

# 管理服务（端口8083，新终端）
make dev-admin

# 定时任务调度器（可选，新终端）
make dev-scheduler
```

### 6. 启动前端

```bash
# 博客前端（端口3001，新终端）
cd frontend/blog
pnpm install
pnpm dev

# 管理前端（端口3002，新终端）
cd frontend/admin
pnpm install
pnpm dev
```

## 生产环境部署

### Docker 部署（推荐）

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 手动部署

#### 后端编译

```bash
# 编译所有服务
make build-all

# 或单独编译
make build          # 用户服务
make build-admin    # 管理服务
make build-scheduler # 调度器
```

#### 前端构建

```bash
# 构建博客前端
cd frontend/blog
pnpm build

# 构建管理前端
cd frontend/admin
pnpm build
```

#### 使用 Nginx 部署前端

参考 `frontend/blog/nginx.conf` 和 `frontend/admin/nginx.conf` 配置 Nginx。

## 验证部署

### 检查服务状态

```bash
# 检查后端健康状态
curl http://localhost:8081/health
curl http://localhost:8083/health

# 检查前端
curl http://localhost:3001
curl http://localhost:3002
```

### 默认账号

- **管理后台**: admin / admin123
- **博客系统**: 可注册新账号

## 常见问题

### 端口冲突

修改配置文件中的端口设置：
- `config/config.yaml` - 用户服务端口
- `config/admin.yaml` - 管理服务端口
- `frontend/blog/vite.config.js` - 博客前端端口
- `frontend/admin/vite.config.js` - 管理前端端口

### 数据库连接失败

检查配置文件中的数据库连接信息，确保：
- MySQL 服务已启动
- 数据库已创建
- 用户名和密码正确
- 防火墙允许连接

### 前端构建失败

```bash
# 清理并重新安装依赖
rm -rf node_modules pnpm-lock.yaml
pnpm install
pnpm build
```

## 更多信息

- [快速开始指南](QUICKSTART.md)
- [API文档](docs/API-REFERENCE.md)
- [数据库设计](docs/database-design.md)
