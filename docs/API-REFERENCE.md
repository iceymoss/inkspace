# API 接口文档

> 基于 RESTful 风格的API设计

## 📖 目录

- [认证相关](#认证相关)
- [用户系统](#用户系统)
- [关注功能](#关注功能-🆕)
- [文章系统](#文章系统)
- [收藏功能](#收藏功能-🆕)
- [评论系统](#评论系统)
- [分类标签](#分类标签)
- [作品展示](#作品展示)
- [管理后台](#管理后台)

---

## 🔐 认证说明

### JWT Token

所有需要认证的接口需要在请求头中携带Token：

```http
Authorization: Bearer <your_jwt_token>
```

### 权限级别

- 🔓 **Public**: 无需认证
- 🔒 **Auth**: 需要登录
- 👑 **Admin**: 需要管理员权限

---

## 认证相关

### 用户注册
```http
POST /api/register
```

**权限**: 🔓 Public

**请求体**:
```json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com",
  "nickname": "测试用户"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 用户登录
```http
POST /api/login
```

**权限**: 🔓 Public

**请求体**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  }
}
```

---

## 用户系统

### 获取个人信息
```http
GET /api/profile
```

**权限**: 🔒 Auth

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "nickname": "管理员",
    "avatar": "",
    "bio": "",
    "article_count": 5,
    "comment_count": 10,
    "following_count": 3,
    "follower_count": 8,
    "favorite_count": 12
  }
}
```

---

### 更新个人信息
```http
PUT /api/profile
```

**权限**: 🔒 Auth

**请求体**:
```json
{
  "nickname": "新昵称",
  "email": "new@example.com",
  "bio": "个人简介",
  "avatar": "https://example.com/avatar.jpg"
}
```

---

### 获取用户主页
```http
GET /api/users/:id
```

**权限**: 🔓 Public

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 2,
    "username": "user1",
    "nickname": "用户1",
    "avatar": "",
    "bio": "这是我的个人简介",
    "article_count": 15,
    "follower_count": 100,
    "following_count": 50
  }
}
```

---

### 获取用户文章列表
```http
GET /api/users/:id/articles?page=1&page_size=10
```

**权限**: 🔓 Public

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

---

## 关注功能 🆕

### 关注用户
```http
POST /api/users/:id/follow
```

**权限**: 🔒 Auth

**响应**:
```json
{
  "code": 0,
  "message": "关注成功"
}
```

---

### 取消关注
```http
DELETE /api/users/:id/follow
```

**权限**: 🔒 Auth

**响应**:
```json
{
  "code": 0,
  "message": "取消关注成功"
}
```

---

### 获取关注统计
```http
GET /api/users/:id/follow-stats
```

**权限**: 🔓 Public

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "following_count": 50,
    "follower_count": 100,
    "is_following": true,
    "is_follower": false
  }
}
```

**说明**:
- `is_following`: 当前用户是否已关注该用户
- `is_follower`: 该用户是否关注了当前用户（互关）

---

### 获取关注列表
```http
GET /api/users/:id/following?page=1&page_size=20
```

**权限**: 🔓 Public

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "user": {
          "id": 3,
          "username": "user3",
          "nickname": "用户3",
          "avatar": ""
        },
        "created_at": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 获取粉丝列表
```http
GET /api/users/:id/followers?page=1&page_size=20
```

**权限**: 🔓 Public

**响应**: 同关注列表格式

---

## 文章系统

### 获取文章列表
```http
GET /api/articles?page=1&page_size=10&category_id=1&tag_id=1&keyword=Go
```

**权限**: 🔓 Public

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）
- `category_id`: 分类ID
- `tag_id`: 标签ID
- `keyword`: 搜索关键词
- `status`: 状态（默认1，仅显示已发布）

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Go语言入门教程",
        "summary": "这是一篇Go语言入门教程",
        "cover": "",
        "category": {
          "id": 1,
          "name": "技术分享"
        },
        "tags": [
          {"id": 1, "name": "Go", "color": "#00ADD8"}
        ],
        "author": {
          "id": 1,
          "nickname": "管理员"
        },
        "view_count": 100,
        "like_count": 10,
        "comment_count": 5,
        "favorite_count": 8,
        "is_top": false,
        "created_at": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 获取文章详情
```http
GET /api/articles/:id
```

**权限**: 🔓 Public

**说明**: 会自动增加浏览计数

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "Go语言入门教程",
    "content": "# Go语言入门\n\n这是文章内容...",
    "summary": "摘要",
    "cover": "",
    "category": {...},
    "tags": [...],
    "author": {...},
    "view_count": 101,
    "like_count": 10,
    "comment_count": 5,
    "favorite_count": 8,
    "word_count": 1500,
    "reading_time": 5,
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 创建文章
```http
POST /api/articles
```

**权限**: 🔒 Auth

**请求体**:
```json
{
  "title": "文章标题",
  "content": "文章内容（Markdown）",
  "summary": "文章摘要",
  "cover": "https://example.com/cover.jpg",
  "category_id": 1,
  "tag_ids": [1, 2, 3],
  "status": 1,
  "is_top": false,
  "is_recommend": false
}
```

---

### 更新文章
```http
PUT /api/articles/:id
```

**权限**: 🔒 Auth (作者或管理员)

---

### 删除文章
```http
DELETE /api/articles/:id
```

**权限**: 🔒 Auth (作者或管理员)

---

### 点赞文章
```http
POST /api/articles/:id/like
```

**权限**: 🔓 Public

---

## 收藏功能 🆕

### 收藏文章
```http
POST /api/articles/:id/favorite
```

**权限**: 🔒 Auth

**响应**:
```json
{
  "code": 0,
  "message": "收藏成功"
}
```

---

### 取消收藏
```http
DELETE /api/articles/:id/favorite
```

**权限**: 🔒 Auth

---

### 检查收藏状态
```http
GET /api/articles/:id/is-favorited
```

**权限**: 🔒 Auth

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_favorited": true
  }
}
```

---

### 获取我的收藏
```http
GET /api/favorites?page=1&page_size=20
```

**权限**: 🔒 Auth

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "article_id": 10,
        "article": {
          "id": 10,
          "title": "收藏的文章",
          "summary": "摘要",
          "cover": "",
          "view_count": 500
        },
        "created_at": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 12,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 获取用户收藏列表
```http
GET /api/users/:id/favorites?page=1&page_size=20
```

**权限**: 🔓 Public

---

## 评论系统

### 获取评论列表
```http
GET /api/comments?article_id=1&page=1&page_size=10
```

**权限**: 🔓 Public

**查询参数**:
- `article_id`: 文章ID（必填）
- `page`: 页码
- `page_size`: 每页数量

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "article_id": 1,
        "user": {
          "id": 2,
          "nickname": "用户1"
        },
        "content": "这是一条评论",
        "parent_id": null,
        "like_count": 5,
        "reply_count": 3,
        "created_at": "2024-01-01T12:00:00Z",
        "replies": [
          {
            "id": 2,
            "content": "这是回复",
            "created_at": "2024-01-01T13:00:00Z"
          }
        ]
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 发表评论
```http
POST /api/comments
```

**权限**: 🔓 Public (登录用户或游客)

**请求体**:
```json
{
  "article_id": 1,
  "content": "这是我的评论",
  "parent_id": null,
  "nickname": "游客昵称",
  "email": "guest@example.com",
  "website": "https://example.com"
}
```

**说明**:
- 登录用户：不需要填写 `nickname`, `email`, `website`
- 游客：需要填写 `nickname` 和 `email`

---

### 删除评论
```http
DELETE /api/comments/:id
```

**权限**: 🔒 Auth (作者或管理员)

---

## 分类标签

### 获取分类列表
```http
GET /api/categories
```

**权限**: 🔓 Public

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "技术分享",
      "slug": "tech",
      "description": "技术相关文章",
      "article_count": 50,
      "sort": 10
    }
  ]
}
```

---

### 获取标签列表
```http
GET /api/tags
```

**权限**: 🔓 Public

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Go",
      "slug": "go",
      "color": "#00ADD8",
      "article_count": 30
    }
  ]
}
```

---

## 作品展示

### 获取作品列表
```http
GET /api/works?page=1&page_size=12&status=1
```

**权限**: 🔓 Public

---

### 获取作品详情
```http
GET /api/works/:id
```

**权限**: 🔓 Public

**说明**: 会自动增加浏览计数

---

## 管理后台

### 用户管理

#### 获取用户列表
```http
GET /api/admin/users?page=1&page_size=10
```

**权限**: 👑 Admin

---

### 文章管理

#### 创建文章
```http
POST /api/admin/articles
```

**权限**: 👑 Admin

**请求体**: 同创建文章接口

---

#### 更新文章
```http
PUT /api/admin/articles/:id
```

**权限**: 👑 Admin

---

#### 删除文章
```http
DELETE /api/admin/articles/:id
```

**权限**: 👑 Admin

---

### 分类管理

#### 创建分类
```http
POST /api/admin/categories
```

**权限**: 👑 Admin

**请求体**:
```json
{
  "name": "新分类",
  "slug": "new-category",
  "description": "分类描述",
  "sort": 10
}
```

---

#### 更新分类
```http
PUT /api/admin/categories/:id
```

**权限**: 👑 Admin

---

#### 删除分类
```http
DELETE /api/admin/categories/:id
```

**权限**: 👑 Admin

**说明**: 如果分类下有文章，无法删除

---

### 标签管理

#### 创建标签
```http
POST /api/admin/tags
```

**权限**: 👑 Admin

**请求体**:
```json
{
  "name": "新标签",
  "slug": "new-tag",
  "color": "#409eff"
}
```

---

### 评论管理

#### 审核评论
```http
PUT /api/admin/comments/:id/status
```

**权限**: 👑 Admin

**请求体**:
```json
{
  "status": 1
}
```

**状态说明**:
- `1`: 通过
- `0`: 待审核
- `-1`: 拒绝

---

### 作品管理

#### 创建作品
```http
POST /api/admin/works
```

**权限**: 👑 Admin

**请求体**:
```json
{
  "title": "作品标题",
  "description": "作品描述",
  "cover": "https://example.com/cover.jpg",
  "images": ["url1", "url2"],
  "link": "https://project.com",
  "github_url": "https://github.com/user/repo",
  "demo_url": "https://demo.com",
  "tech_stack": "Go,Vue,Docker",
  "sort": 10,
  "status": 1
}
```

---

## 📊 响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {...}
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误信息"
}
```

### 分页响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 🔢 状态码

| 状态码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未登录 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

---

## 📝 数据验证规则

### 用户注册
- username: 3-50字符，字母数字下划线
- password: 6-50字符
- email: 有效的邮箱格式

### 文章创建
- title: 必填，最多200字符
- content: 必填，支持Markdown
- summary: 最多500字符

### 评论发表
- content: 必填，最多500字符
- nickname: 最多50字符
- email: 有效的邮箱格式

---

## 🔧 开发测试

### 使用 curl 测试

```bash
# 登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 获取文章列表
curl http://localhost:8080/api/articles

# 关注用户（需要token）
curl -X POST http://localhost:8080/api/users/2/follow \
  -H "Authorization: Bearer <your_token>"

# 收藏文章（需要token）
curl -X POST http://localhost:8080/api/articles/1/favorite \
  -H "Authorization: Bearer <your_token>"
```

### 使用 Postman

1. 导入环境变量
2. 设置 BASE_URL: `http://localhost:8080`
3. 登录获取Token
4. 在Headers中添加Authorization

---

## 📚 相关文档

- [数据库设计文档](database-design.md) - 详细的表结构
- [开发计划](DEVELOPMENT-PLAN.md) - 功能开发规划
- [快速开始](../QUICKSTART.md) - 项目启动指南

---

**API版本**: v1.0  
**最后更新**: 2024-01-01  
**接口总数**: 28+个

