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

### 步骤0：配置环境变量（可选）

项目支持使用 `.env` 文件或环境变量来配置，环境变量会覆盖 YAML 配置文件中的值。

```bash
# 复制配置模板
cp .env.example .env

# 编辑 .env 文件，修改数据库、Redis等配置
# 如果不创建 .env 文件，将使用 config/config.yaml 中的默认配置
```

**配置优先级**：环境变量 > .env 文件 > YAML 配置文件

### 步骤1：启动数据库
```bash
docker-compose up -d mysql redis
```

### 步骤2：初始化数据库
```bash
# 数据库迁移会在服务启动时自动执行（通过 GORM AutoMigrate）
# 初始化基础数据（可选，包含默认管理员账号等）
mysql -h localhost -u inkspace -pinkspace123 inkspace < scripts/init.sql
```

### 步骤3：启动后端服务

**终端1 - 用户服务**:
```bash
go run cmd/server/main.go
```

**终端2 - 管理服务**:
```bash
go run cmd/admin/main.go
```

**终端3 - 定时任务调度器（可选但推荐）**:
```bash
go run cmd/scheduler/main.go
```

### 步骤4：启动前端

**终端4 - 博客前端**:
```bash
cd web/blog
pnpm install  # 首次
pnpm dev
```

**终端5 - 管理前端**:
```bash
cd web/admin
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

## 📝 常用命令

```bash
# 启动服务
go run cmd/server/main.go    # 启动用户服务 (8081)
go run cmd/admin/main.go     # 启动管理服务 (8083)
go run cmd/scheduler/main.go # 启动定时任务调度器

# 编译服务
go build -o bin/server cmd/server/main.go      # 编译用户服务
go build -o bin/admin cmd/admin/main.go        # 编译管理服务
go build -o bin/scheduler cmd/scheduler/main.go # 编译定时任务调度器

# 数据库
# 数据库迁移会在服务启动时自动执行
mysql -h localhost -u inkspace -pinkspace123 inkspace < scripts/init.sql  # 初始化数据
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
- `web/blog/vite.config.js` - server.port (博客前端)
- `web/admin/vite.config.js` - server.port (管理前端)

### pnpm未安装
```bash
npm install -g pnpm
```

### 数据库连接失败
检查 `config/config.yaml` 中的数据库配置

