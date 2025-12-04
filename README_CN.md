# My Site - Go个人网站

[English](README.md) | 简体中文

完整的项目文档请查看 [README.md](README.md)

## 快速开始

### 使用Docker部署（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd mysite

# 启动所有服务
docker-compose up -d

# 访问应用
# 前端: http://localhost
# 后端API: http://localhost:8080
# 管理后台: http://localhost/admin
```

默认管理员账号：
- 用户名: `admin`
- 密码: `admin123`

### 本地开发

```bash
# 后端
go mod download
go run main.go

# 前端（新终端）
cd frontend
npm install
npm run dev
```

## 主要功能

- ✅ 文章管理（发布、编辑、删除）
- ✅ 作品展示
- ✅ 评论系统
- ✅ 分类和标签
- ✅ 用户认证（JWT）
- ✅ 后台管理
- ✅ Redis缓存
- ✅ Docker部署

## 技术栈

**后端**: Go + Gin + GORM + MySQL + Redis

**前端**: Vue 3 + Element Plus + Pinia + Vite

## 项目结构

```
mysite/
├── internal/          # 后端代码
├── frontend/          # 前端代码
├── config/           # 配置文件
├── scripts/          # 数据库脚本
└── docker-compose.yml # Docker配置
```

## 许可证

MIT License

