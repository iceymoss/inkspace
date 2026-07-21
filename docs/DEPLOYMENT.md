# 部署指南

## 📋 环境要求

- Docker 20.10+
- Docker Compose 2.0+
- MySQL 8.0+（如果使用外部数据库）
- Redis 7+（如果使用外部 Redis）

---

## 🚀 部署方式

本项目提供两种部署方式，根据你的实际情况选择：

### 方式一：完整部署（包含 MySQL 和 Redis）

适用于：**全新部署、开发环境、测试环境**

使用 `docker-compose.yml`，一键启动所有服务，包括 MySQL 和 Redis。

### 方式二：使用外部数据库

适用于：**生产环境、已有数据库服务**

使用 `docker-compose.external-db.yml`，只启动业务服务，连接外部已有的 MySQL 和 Redis。

---

## 方式一：完整部署（包含 MySQL 和 Redis）

### 1. 配置 DNS（必需）

在 DNS 服务商处添加以下 A 记录，指向服务器 IP：

```
is.iceymoss.com        A    <your-server-ip>
admin.is.iceymoss.com  A    <your-server-ip>
```

**注意：** 确保 DNS 解析生效后再继续部署。

### 2. 克隆项目

```bash
git clone <repository-url>
cd inkspace
```

### 3. 配置环境变量（可选）

```bash
# 复制配置模板
cp env.example .env

# 编辑 .env 文件，修改数据库、Redis等配置
# 如果不创建 .env 文件，将使用 config/config.yaml 中的默认配置
```

**配置优先级**：环境变量 > .env 文件 > YAML 配置文件

**默认配置：**
- MySQL 端口：3306
- Redis 端口：6379
- 数据库名：inkspace
- 数据库用户：inkspace
- 数据库密码：inkspace123

### 4. 启动所有服务

```bash
# 构建并启动所有服务（前端会构建并嵌入对应 Go 服务）
docker-compose up -d --build
```

**启动的服务包括：**
- `mysql` - MySQL 数据库 (端口 3306)
- `redis` - Redis 缓存 (端口 6379)
- `backend-1/2/3` - 用户服务（3个实例，负载均衡）
- `admin-backend` - 管理后台服务 (端口 8083)
- `scheduler` - 定时任务调度器
- `nginx-proxy` - Nginx 反向代理 (端口 80/443)

### 5. 查看服务状态

```bash
# 查看所有服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend-1
docker-compose logs -f admin-backend
```

### 6. 访问服务

启动成功后，通过子域名访问：

- **博客前端**: http://is.iceymoss.com
- **管理前端**: http://admin.is.iceymoss.com
- **管理 API**: http://admin.is.iceymoss.com/api （或直接访问 http://<server-ip>:8083/api）

**注意：**
- 博客和管理前端分别嵌入 `server` 与 `admin` 二进制，由 Nginx 按域名代理
- 确保 DNS 解析已生效
- 如果使用 HTTPS，需要配置 SSL 证书（参考 nginx/nginx.conf 中的 HTTPS 配置）

### 8. 默认账号

- **管理后台**: admin / admin123
- **博客系统**: 可注册新账号

### 9. 停止服务

```bash
# 停止所有服务（保留数据）
docker-compose stop

# 停止并删除容器（保留数据卷）
docker-compose down

# 停止并删除容器和数据卷（⚠️ 危险操作）
docker-compose down -v
```

---

## 方式二：使用外部数据库

### 前置条件

1. **配置 DNS（必需）**

在 DNS 服务商处添加以下 A 记录，指向服务器 IP：

```
is.iceymoss.com        A    <your-server-ip>
admin.is.iceymoss.com  A    <your-server-ip>
```

**注意：** 确保 DNS 解析生效后再继续部署。

2. **已有 MySQL 和 Redis 服务**
   - MySQL 8.0+ 已运行
   - Redis 7+ 已运行

3. **创建数据库和用户**

```bash
# 连接到 MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE IF NOT EXISTS inkspace CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 创建用户（如果不存在）
CREATE USER IF NOT EXISTS 'inkspace'@'%' IDENTIFIED BY 'inkspace123';

# 授权
GRANT ALL PRIVILEGES ON inkspace.* TO 'inkspace'@'%';
FLUSH PRIVILEGES;
```

### 1. 克隆项目

```bash
git clone <repository-url>
cd inkspace
```

### 2. 配置数据库连接

编辑 `config/config.yaml` 和 `config/admin.yaml`，修改数据库和 Redis 连接信息：

```yaml
# config/config.yaml
database:
  host: <your-mysql-host>      # 例如: localhost 或 192.168.1.100
  port: 3306
  username: inkspace
  password: inkspace123
  database: inkspace

redis:
  host: <your-redis-host>       # 例如: localhost 或 192.168.1.100
  port: 6379
  password: ""                  # 如果有密码，填写密码
```

### 3. 配置上传文件目录（可选）

默认情况下，上传文件挂载到 `/var/www/inkspace/uploads`。如果需要修改路径：

1. **创建上传目录**（如果使用默认路径）：
```bash
sudo mkdir -p /var/www/inkspace/uploads
sudo chown -R $USER:$USER /var/www/inkspace/uploads
```

2. **修改挂载路径**（如果需要）：
编辑 `docker-compose.external-db.yml`，将所有服务中的：
```yaml
volumes:
  - /var/www/inkspace/uploads:/app/uploads
```
修改为你想要的路径，例如：
```yaml
volumes:
  - /data/inkspace/uploads:/app/uploads  # 或其他路径
```

**注意：**
- 确保目录存在且有正确的读写权限
- 所有后端服务（backend-1/2/3、admin-backend、scheduler）必须使用相同的挂载路径
- 建议使用绝对路径，避免相对路径带来的问题

### 4. 创建 Docker 网络（如果使用外部容器）

如果你的 MySQL 和 Redis 也是 Docker 容器，需要将它们加入同一个网络：

```bash
# 创建网络（如果不存在）
docker network create inkspace-network

# 将外部 MySQL 容器加入网络
docker network connect inkspace-network <your-mysql-container-name>

# 将外部 Redis 容器加入网络
docker network connect inkspace-network <your-redis-container-name>
```

**注意：** 如果 MySQL/Redis 容器在不同的网络，需要确保：
- 在配置文件中使用容器名作为 host（如果在同一网络）
- 或使用宿主机 IP 地址（如果不在同一网络）

### 5. 启动业务服务

```bash
# 使用外部数据库配置启动服务
docker-compose -f docker-compose.external-db.yml up -d --build
```

**启动的服务包括：**
- `backend-1/2/3` - 用户服务（3个实例，负载均衡）
- `admin-backend` - 管理后台服务 (端口 8083)
- `scheduler` - 定时任务调度器
- `nginx-proxy` - Nginx 反向代理 (端口 80/443)

### 6. 查看服务状态

```bash
# 查看服务状态
docker-compose -f docker-compose.external-db.yml ps

# 查看服务日志
docker-compose -f docker-compose.external-db.yml logs -f
```

### 7. 访问服务

通过子域名访问：

- **博客前端**: http://is.iceymoss.com
- **管理前端**: http://admin.is.iceymoss.com
- **管理 API**: http://admin.is.iceymoss.com/api （或直接访问 http://<server-ip>:8083/api）

**注意：**
- 博客和管理前端分别嵌入 `server` 与 `admin` 二进制，由 Nginx 按域名代理
- 确保 DNS 解析已生效
- 用户服务有 3 个实例（backend-1/2/3），通过外层 Nginx 进行负载均衡
- 如果使用 HTTPS，需要配置 SSL 证书（参考 nginx/nginx.conf 中的 HTTPS 配置）

### 8. 停止服务

```bash
# 停止服务
docker-compose -f docker-compose.external-db.yml stop

# 停止并删除容器
docker-compose -f docker-compose.external-db.yml down
```

---

## 🔧 配置说明

### 配置文件位置

- `config/config.yaml` - 用户服务配置
- `config/admin.yaml` - 管理服务配置
- `.env` - 环境变量配置（可选，优先级最高）

### 配置优先级

1. **环境变量** - 系统环境变量
2. **.env 文件** - 项目根目录下的 `.env` 文件
3. **YAML 配置文件** - `config/config.yaml` 和 `config/admin.yaml`

### 主要配置项

**数据库配置：**
```yaml
database:
  host: localhost
  port: 3306
  username: inkspace
  password: inkspace123
  database: inkspace
  charset: utf8mb4
```

**Redis 配置：**
```yaml
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

**服务端口：**
```yaml
server:
  port: 8081  # 用户服务端口

admin:
  port: 8083  # 管理服务端口
```

修改配置文件后，需要重启服务：

```bash
# 方式一：重启所有服务
docker-compose restart

# 方式二：重启特定服务
docker-compose restart backend-1 admin-backend scheduler
```

---

## 🌐 子域名配置

### DNS 设置

无论使用哪种部署方式，都需要配置 DNS 记录：

```
is.iceymoss.com        A    <your-server-ip>
admin.is.iceymoss.com  A    <your-server-ip>
```

**DNS 配置步骤：**

1. 登录你的 DNS 服务商（如 Cloudflare、阿里云 DNS、腾讯云 DNS 等）
2. 添加两条 A 记录：
   - 主机记录：`is`，记录值：服务器 IP 地址
   - 主机记录：`admin`，记录值：服务器 IP 地址
3. 等待 DNS 解析生效（通常几分钟到几小时）

**验证 DNS 解析：**

```bash
# 检查 DNS 解析
nslookup is.iceymoss.com
nslookup admin.is.iceymoss.com

# 或使用 dig
dig is.iceymoss.com
dig admin.is.iceymoss.com
```

### Nginx 配置

Nginx 配置文件位于 `nginx/nginx.conf`，已配置两个子域名：

- `is.iceymoss.com` → 博客前端
- `admin.is.iceymoss.com` → 管理前端

如需修改子域名，编辑 `nginx/nginx.conf` 中的 `server_name` 配置。

### HTTPS 配置

项目提供 `docker-compose.https.yml` 和 `nginx/nginx.https.conf`。以下步骤会申请一张同时覆盖 `is.iceymoss.com` 和 `admin.is.iceymoss.com` 的 Let's Encrypt 证书。

1. 先使用基础 Compose 启动 HTTP 服务。基础 Nginx 已开放 `/.well-known/acme-challenge/`：

```bash
# 内置 MySQL/Redis
docker compose up -d --build

# 或使用外部 MySQL/Redis
docker compose -f docker-compose.external-db.yml up -d --build
```

2. 在服务器仓库根目录签发证书，将邮箱替换为实际运维邮箱：

```bash
docker run --rm \
  -v "$PWD/certbot/www:/var/www/certbot" \
  -v "$PWD/certbot/conf:/etc/letsencrypt" \
  certbot/certbot certonly --webroot \
  --webroot-path /var/www/certbot \
  --email your-email@example.com \
  --agree-tos --no-eff-email \
  --cert-name is.iceymoss.com \
  -d is.iceymoss.com \
  -d admin.is.iceymoss.com
```

3. 证书签发成功后叠加 HTTPS 配置启动。HTTP 请求会自动 301 跳转到 HTTPS：

```bash
# 内置 MySQL/Redis
docker compose -f docker-compose.yml -f docker-compose.https.yml up -d

# 外部 MySQL/Redis
docker compose -f docker-compose.external-db.yml -f docker-compose.https.yml up -d
```

4. 验证两个域名：

```bash
curl -I https://is.iceymoss.com/health
curl -I https://admin.is.iceymoss.com/health
```

5. 手动续期并热加载 Nginx：

```bash
docker run --rm \
  -v "$PWD/certbot/www:/var/www/certbot" \
  -v "$PWD/certbot/conf:/etc/letsencrypt" \
  certbot/certbot renew --webroot --webroot-path /var/www/certbot

docker exec inkspace-nginx-proxy nginx -s reload
```

可将上述续期命令加入服务器 cron，每天执行一次。Certbot 只会在证书接近过期时实际续期。`certbot/conf/` 中的证书和私钥不会提交到 Git。

---

## 📊 服务架构

### 方式一（完整部署）

```
┌─────────────────────────────────────┐
│  Nginx Proxy (80/443)                │
│  ├── is.iceymoss.com                │
│  └── admin.is.iceymoss.com          │
├─────────────────────────────────────┤
│  Application Services               │
│  ├── Backend-1/2/3 + Blog SPA       │
│  ├── Admin Backend + Admin SPA      │
│  └── Scheduler                      │
├─────────────────────────────────────┤
│  Data Services                      │
│  ├── MySQL (3306)                   │
│  └── Redis (6379)                   │
└─────────────────────────────────────┘
```

### 方式二（外部数据库）

```
┌─────────────────────────────────────┐
│  Nginx Proxy (80/443)                │
│  ├── is.iceymoss.com                │
│  └── admin.is.iceymoss.com          │
├─────────────────────────────────────┤
│  Application Services               │
│  ├── Backend-1/2/3 + Blog SPA       │
│  ├── Admin Backend + Admin SPA      │
│  └── Scheduler                      │
├─────────────────────────────────────┤
│  External Services                  │
│  ├── MySQL (外部)                   │
│  └── Redis (外部)                   │
└─────────────────────────────────────┘
```

---

## ✅ 验证部署

### 检查服务健康状态

**方式一（完整部署）：**
```bash
# 检查用户服务
curl http://localhost:8081/health
# 预期响应: {"status":"ok"}

# 检查管理服务
curl http://localhost:8083/health
# 预期响应: {"status":"ok","service":"admin"}

# 检查前端
curl http://localhost:3001
curl http://localhost:3002
```

**方式二（外部数据库）：**
```bash
# 检查后端服务（通过容器内部）
docker-compose -f docker-compose.external-db.yml exec backend-1 curl http://localhost:8081/health
# 或检查所有3个实例
docker-compose -f docker-compose.external-db.yml exec backend-1 curl http://localhost:8081/health
docker-compose -f docker-compose.external-db.yml exec backend-2 curl http://localhost:8081/health
docker-compose -f docker-compose.external-db.yml exec backend-3 curl http://localhost:8081/health

# 检查管理服务
curl http://localhost:8083/health
# 预期响应: {"status":"ok","service":"admin"}

# 检查前端（通过子域名访问）
curl http://is.iceymoss.com
curl http://admin.is.iceymoss.com
```

### 检查数据库连接

```bash
# 方式一：进入 MySQL 容器
docker-compose exec mysql mysql -u inkspace -pinkspace123 inkspace -e "SHOW TABLES;"

# 方式二：从宿主机连接
mysql -h localhost -u inkspace -pinkspace123 inkspace -e "SHOW TABLES;"
```

### 检查 Redis 连接

```bash
# 方式一：进入 Redis 容器
docker-compose exec redis redis-cli ping
# 预期响应: PONG

# 方式二：从宿主机连接
redis-cli -h localhost ping
```

---

## 🐛 常见问题

### 端口冲突

如果遇到端口冲突，可以修改：

1. **修改 docker-compose.yml 中的端口映射**
2. **修改配置文件中的端口设置**
   - `config/config.yaml` - 用户服务端口
   - `config/admin.yaml` - 管理服务端口

### 数据库连接失败

**方式一（完整部署）：**
```bash
# 检查 MySQL 容器是否运行
docker-compose ps mysql

# 查看 MySQL 日志
docker-compose logs mysql

# 检查网络连接
docker-compose exec backend-1 ping mysql
```

**方式二（外部数据库）：**
- 确认 MySQL 服务已启动
- 检查配置文件中的 host、port、username、password
- 确认防火墙允许连接
- 如果使用容器，确认网络配置正确

### Redis 连接失败

```bash
# 检查 Redis 容器是否运行
docker-compose ps redis

# 查看 Redis 日志
docker-compose logs redis

# 测试连接
docker-compose exec backend-1 redis-cli -h redis ping
```

### 服务启动失败

**如果遇到 ContainerConfig 错误：**
```
ERROR: for backend-1  'ContainerConfig'
docker.errors.ImageNotFound: 404 Client Error: Not Found
KeyError: 'ContainerConfig'
```

这是旧容器引用损坏镜像导致的，**请使用 [更新部署](#-更新部署) 部分的正确流程**：
1. 先执行 `docker-compose down`（或 `docker-compose -f docker-compose.external-db.yml down`）
2. 再执行 `docker-compose up -d --build --force-recreate`

**普通服务启动失败：**
```bash
# 查看详细日志
# 方式一
docker-compose logs -f <service-name>

# 方式二
docker-compose -f docker-compose.external-db.yml logs -f <service-name>

# 检查容器状态
# 方式一
docker-compose ps

# 方式二
docker-compose -f docker-compose.external-db.yml ps

# 重启服务
# 方式一
docker-compose restart <service-name>

# 方式二
docker-compose -f docker-compose.external-db.yml restart <service-name>
```

### 数据持久化

**方式一（完整部署）：**
- MySQL 数据存储在 Docker volume `mysql_data`
- Redis 数据存储在 Docker volume `redis_data`
- 上传文件存储在 `./uploads` 目录（项目根目录）

**方式二（外部数据库）：**
- 上传文件存储在 `/var/www/inkspace/uploads`（默认路径，可在 docker-compose.external-db.yml 中修改）

**查看 volumes：**
```bash
docker volume ls | grep inkspace
```

**备份数据：**
```bash
# 备份 MySQL
docker-compose exec mysql mysqldump -u inkspace -pinkspace123 inkspace > backup.sql

# 备份 Redis
docker-compose exec redis redis-cli SAVE
docker cp inkspace-redis:/data/dump.rdb ./backup.rdb
```

---


## 🔄 更新部署

### 正确更新部署流程（避免 ContainerConfig 错误）

**重要提示：** 更新部署时，必须先停止并删除旧容器，再重新构建和启动，这样可以避免 ContainerConfig 错误。

#### 方式一（完整部署）

```bash
# 1. 拉取最新代码
git pull

# 2. 停止并删除当前项目的所有容器和网络
# 注意：docker-compose down 只影响当前 compose 文件中的容器，不会影响其他服务的容器
docker-compose down

# 3. 重新构建镜像并启动服务（强制重新创建容器）
docker-compose up -d --build --force-recreate

# 4. 查看更新日志
docker-compose logs -f
```

#### 方式二（外部数据库）

```bash
# 1. 拉取最新代码
git pull

# 2. 停止并删除当前项目的所有容器和网络
# 注意：docker-compose down 只影响当前 compose 文件中的容器，不会影响其他服务的容器
docker-compose -f docker-compose.external-db.yml down

# 3. 重新构建镜像并启动服务（强制重新创建容器）
docker-compose -f docker-compose.external-db.yml up -d --build --force-recreate

# 4. 查看更新日志
docker-compose -f docker-compose.external-db.yml logs -f
```

**命令说明：**
- `docker-compose down` - 停止并删除当前 compose 文件中定义的所有容器和网络，**安全操作，只影响本项目容器**
- `--build` - 重新构建镜像
- `--force-recreate` - 强制重新创建所有容器，即使配置没有变化（**关键：避免 ContainerConfig 错误**）

### 只更新特定服务（不推荐，可能出错）

如果只修改了某个服务的代码，可以尝试只更新该服务，但**不推荐**，因为可能出现 ContainerConfig 错误：

```bash
# 方式一
docker-compose up -d --build <service-name>

# 方式二
docker-compose -f docker-compose.external-db.yml up -d --build <service-name>
```

**如果遇到 ContainerConfig 错误，必须使用上述完整更新流程。**

### 重新部署（完全重置）

如果需要完全重新部署（例如配置更改、环境问题等）：

**方式一：**
```bash
# 1. 停止并删除所有容器、网络、卷
# ⚠️ 注意：-v 参数会删除数据卷，包括 MySQL 和 Redis 数据
docker-compose down -v

# 2. 重新构建并启动
docker-compose up -d --build --force-recreate


**方式二：**
```bash
# 1. 停止并删除所有容器和网络
# 注意：不会删除数据卷，因为方式二使用外部数据库
docker-compose -f docker-compose.external-db.yml down

# 2. 重新构建并启动
docker-compose -f docker-compose.external-db.yml up -d --build --force-recreate
```

### 数据库迁移

数据库迁移会在服务启动时自动执行（通过 GORM AutoMigrate），无需手动操作。

如果遇到迁移问题，可以查看日志：
```bash
# 方式一
docker-compose logs backend-1 | grep -i migrate

# 方式二
docker-compose -f docker-compose.external-db.yml logs backend-1 | grep -i migrate
```
