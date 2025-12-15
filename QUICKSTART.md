# 快速启动指南

## 🏗️ 系统架构

完全独立的**四服务架构**：

```
博客前端 (:3001) → 用户服务 (:8081)
管理前端 (:3002) → 管理服务 (:8083)
             ↓
        共享数据库
```

---

## ⚡ 快速启动（推荐）

### 步骤1：启动数据库
```bash
docker-compose up -d mysql redis
```

### 步骤2：初始化数据库
```bash
make db-migrate && make db-init
```

### 步骤3：启动后端服务

**终端1 - 用户服务**:
```bash
make dev
# 或: go run cmd/server/main.go
```

**终端2 - 管理服务**:
```bash
make dev-admin
# 或: go run cmd/admin/main.go
```

**终端3 - 定时任务调度器（可选但推荐）**:
```bash
make dev-scheduler
# 或: go run cmd/scheduler/main.go
```

### 步骤4：启动前端

**终端4 - 博客前端**:
```bash
cd frontend/blog
pnpm install  # 首次
pnpm dev
```

**终端5 - 管理前端**:
```bash
cd frontend/admin
pnpm install  # 首次
pnpm dev
```

---

## 🌐 访问地址

- **博客前端**: http://localhost:3001
- **管理前端**: http://localhost:3002/login
- **用户API**: http://localhost:8081/api
- **管理API**: http://localhost:8083/api

---

## 🔐 默认账号

### 管理后台
```
地址: http://localhost:3002/login
账号: admin
密码: admin123
```

### 博客系统
```
地址: http://localhost:3001/login
可以注册新账号
```

---

## 📝 Makefile命令

```bash
make dev              # 启动用户服务 (8081)
make dev-admin        # 启动管理服务 (8083)
make dev-scheduler    # 启动定时任务调度器
make build            # 编译用户服务
make build-admin      # 编译管理服务
make build-scheduler  # 编译定时任务调度器
make build-all        # 编译所有服务
make db-migrate       # 数据库迁移
make db-init          # 初始化数据
```

---

## 🎯 服务说明

| 服务 | 端口 | 目录 | 用途 |
|------|------|------|------|
| 博客前端 | 3001 | frontend/blog | 用户界面 |
| 管理前端 | 3002 | frontend/admin | 管理界面 |
| 用户服务 | 8081 | cmd/server | 用户API |
| 管理服务 | 8083 | cmd/admin | 管理API |
| 定时任务 | - | cmd/scheduler | 后台统计（可选） |

---

## ✅ 验证启动成功

### 检查后端
```bash
curl http://localhost:8081/health  # {"status":"ok"}
curl http://localhost:8083/health  # {"status":"ok","service":"admin"}
```

### 检查前端
- 访问 http://localhost:3001 - 应该看到博客首页
- 访问 http://localhost:3002 - 应该跳转到管理登录

---

## 🐛 常见问题

### 端口冲突
修改对应的配置文件：
- `config/config.yaml` - server.port (用户服务)
- `config/admin.yaml` - server.port (管理服务)
- `frontend/blog/vite.config.js` - server.port (博客前端)
- `frontend/admin/vite.config.js` - server.port (管理前端)

### pnpm未安装
```bash
npm install -g pnpm
```

### 数据库连接失败
检查 `config/config.yaml` 中的数据库配置

