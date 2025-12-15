# Docker Compose 使用指南

本项目提供两个版本的 docker-compose 配置，适用于不同的部署场景。

## 📋 版本说明

### 1. 完整版本 (`docker-compose.yml`)
包含所有服务：MySQL、Redis、后端服务、前端服务
- **适用场景**：全新部署、开发环境、测试环境
- **特点**：一键启动所有服务，包括数据库

### 2. 简化版本 (`docker-compose.external-db.yml`)
只包含业务服务，使用外部已有的 MySQL 和 Redis
- **适用场景**：生产环境、已有数据库服务
- **特点**：只启动业务服务，复用现有数据库

---

## 🚀 使用方法

### 完整版本（包含 MySQL/Redis）

```bash
# 启动所有服务（包括 MySQL、Redis）
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止所有服务
docker-compose down
```

**服务列表：**
- `mysql` - MySQL 数据库 (端口 3306)
- `redis` - Redis 缓存 (端口 6379)
- `backend` - 博客用户 API (端口 8081)
- `admin-backend` - 管理后台 API (端口 8083)
- `scheduler` - 定时任务调度器
- `blog-frontend` - 博客前端 (端口 3001)
- `admin-frontend` - 管理前端 (端口 3002)

---

### 简化版本（使用外部数据库）

#### 前置条件

1. **确保外部 MySQL 和 Redis 容器已运行**
   ```bash
   docker ps | grep -E "(mysql|redis)"
   ```

2. **创建网络并连接外部容器**
   ```bash
   # 创建网络（如果不存在）
   docker network create inkspace-network
   
   # 将外部 MySQL 容器加入网络
   docker network connect inkspace-network mysql-inkspace
   
   # 将外部 Redis 容器加入网络
   docker network connect inkspace-network redis-inkspace
   ```

3. **配置数据库连接**
   
   编辑 `config/config.yaml` 和 `config/admin.yaml`，确保数据库配置指向外部容器：
   
   ```yaml
   database:
     host: mysql-inkspace  # 外部 MySQL 容器名
     port: 3306
     username: root        # 根据你的外部容器配置调整
     password: root         # 根据你的外部容器配置调整
     database: mysite       # 根据你的外部容器数据库名调整
   
   redis:
     host: redis-inkspace  # 外部 Redis 容器名
     port: 6379
   ```

#### 启动服务

```bash
# 使用简化版本启动服务
docker-compose -f docker-compose.external-db.yml up -d --build

# 查看服务状态
docker-compose -f docker-compose.external-db.yml ps

# 查看日志
docker-compose -f docker-compose.external-db.yml logs -f

# 停止服务
docker-compose -f docker-compose.external-db.yml down
```

**服务列表：**
- `backend` - 博客用户 API (端口 8081)
- `admin-backend` - 管理后台 API (端口 8083)
- `scheduler` - 定时任务调度器
- `blog-frontend` - 博客前端 (端口 3001)
- `admin-frontend` - 管理前端 (端口 3002)

---

## 🌐 访问地址

启动成功后，可以通过以下地址访问：

- **博客前端**: http://localhost:3001
- **管理前端**: http://localhost:3002
- **用户 API**: http://localhost:8081/api
- **管理 API**: http://localhost:8083/api

---

## 🔧 常见问题

### 端口冲突

如果遇到端口冲突（如 3306、6379 已被占用），可以使用简化版本：

```bash
docker-compose -f docker-compose.external-db.yml up -d --build
```

### 网络连接问题

如果使用简化版本时，服务无法连接到外部数据库：

1. 确认外部容器已加入 `inkspace-network` 网络
2. 检查配置文件中的数据库主机名是否正确
3. 确认外部容器正在运行：`docker ps | grep -E "(mysql|redis)"`

### 配置文件路径

配置文件通过 volume 挂载：
- `./config` → `/app/config` (容器内)
- `./uploads` → `/app/uploads` (容器内)

修改配置文件后，需要重启服务：
```bash
docker-compose restart backend admin-backend scheduler
```

---

## 📝 注意事项

1. **数据持久化**：完整版本使用 Docker volumes 持久化 MySQL 和 Redis 数据
2. **配置文件**：两种版本都使用相同的配置文件 (`config/config.yaml`, `config/admin.yaml`)
3. **网络隔离**：所有服务在 `inkspace-network` 网络中，可以通过服务名互相访问
4. **端口映射**：确保宿主机端口未被占用

