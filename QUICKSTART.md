# 快速开始指南

## 方式一：使用Docker（推荐）

### 1. 前置条件
- 安装 Docker 和 Docker Compose

### 2. 启动项目

```bash
# 克隆项目
git clone <your-repository-url>
cd mysite

# 启动所有服务（MySQL, Redis, 后端, 前端）
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 3. 访问应用

- **前端网站**: http://localhost
- **管理后台**: http://localhost/admin
- **后端API**: http://localhost:8080

### 4. 默认账号

- 用户名: `admin`
- 密码: `admin123`

### 5. 停止服务

```bash
docker-compose down
```

---

## 方式二：本地开发

### 1. 前置条件

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 7+

### 2. 数据库准备

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE mysite CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 导入初始数据
mysql -u root -p mysite < scripts/init.sql
```

### 3. 后端启动

```bash
# 安装依赖
go mod download

# 修改配置文件
# 编辑 config/config.yaml，配置数据库和Redis连接信息

# 启动后端
go run main.go

# 或使用 Make
make dev
```

后端将在 http://localhost:8080 运行

### 4. 前端启动

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端将在 http://localhost:3000 运行

---

## 使用 Makefile

项目提供了 Makefile 来简化常用操作：

```bash
# 查看所有可用命令
make help

# 开发模式运行后端
make dev

# 构建后端
make build

# 运行测试
make test

# Docker相关
make docker-build    # 构建镜像
make docker-up       # 启动容器
make docker-down     # 停止容器
make docker-logs     # 查看日志

# 前端相关
make frontend-install  # 安装依赖
make frontend-dev      # 开发模式
make frontend-build    # 构建生产版本
```

---

## 常见问题

### 1. 数据库连接失败

检查 `config/config.yaml` 中的数据库配置是否正确：

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: root
  database: mysite
```

### 2. Redis连接失败

检查Redis是否启动：

```bash
redis-cli ping
# 应该返回 PONG
```

### 3. 前端无法访问后端API

确保后端已启动，并检查前端 `vite.config.js` 中的代理配置：

```javascript
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true
  }
}
```

### 4. Docker容器启动失败

```bash
# 查看日志
docker-compose logs

# 重新构建
docker-compose build --no-cache
docker-compose up -d
```

### 5. 端口冲突

如果默认端口被占用，可以修改 `docker-compose.yml` 中的端口映射：

```yaml
services:
  frontend:
    ports:
      - "8888:80"  # 将前端改为8888端口
  
  backend:
    ports:
      - "9090:8080"  # 将后端改为9090端口
```

---

## 下一步

1. 修改个人信息
   - 编辑 `frontend/src/views/About.vue`
   - 更新个人简介和联系方式

2. 自定义配置
   - 修改 `config/config.yaml` 调整系统配置
   - 修改 JWT密钥以提高安全性

3. 添加内容
   - 登录管理后台
   - 创建分类和标签
   - 发布第一篇文章
   - 上传作品展示

4. 部署到生产环境
   - 修改数据库密码
   - 更新JWT密钥
   - 配置域名和HTTPS
   - 设置防火墙规则

---

## 技术支持

如有问题，请提交 Issue 或查看完整文档 [README.md](README.md)

