# 🔧 API字段名修复

## 🐛 问题描述

### 症状
调用 `/api/articles/1/is-liked` 返回：
```json
{
  "code": 0,
  "data": {
    "liked": false  // ❌ 字段名不匹配
  }
}
```

前端读取 `response.data.is_liked` → `undefined`

---

## ✅ 修复方案

### 统一返回两种字段名

修改所有状态检查API，同时返回两种字段名以兼容前后端：

#### 1. 点赞状态检查（2个API）

**文件：** `internal/handler/like_handler.go`

```go
// CheckArticleLiked - 检查文章点赞状态
// GET /api/articles/:id/is-liked
utils.Success(c, gin.H{"liked": liked, "is_liked": liked})

// CheckWorkLiked - 检查作品点赞状态  
// GET /api/works/:id/liked
utils.Success(c, gin.H{"liked": liked, "is_liked": liked})
```

#### 2. 收藏状态检查（2个API）

**文件：** `internal/handler/favorite_handler.go`

```go
// CheckFavorited - 检查文章收藏状态
// GET /api/articles/:id/is-favorited
utils.Success(c, gin.H{"is_favorited": favorited, "favorited": favorited})

// CheckWorkFavorited - 检查作品收藏状态
// GET /api/works/:id/favorited
utils.Success(c, gin.H{"favorited": favorited, "is_favorited": favorited})
```

#### 3. 评论点赞状态检查

**文件：** `internal/handler/like_handler.go`

```go
// CheckCommentLiked - 检查评论点赞状态
// GET /api/comments/:id/is-liked
utils.Success(c, gin.H{"liked": false, "is_liked": false})
```

---

## 📊 修复后的API响应

### 点赞状态
```json
// GET /api/articles/1/is-liked
{
  "code": 0,
  "data": {
    "liked": true,      // ✅ 兼容旧版
    "is_liked": true    // ✅ 兼容新版
  }
}
```

### 收藏状态
```json
// GET /api/articles/1/is-favorited
{
  "code": 0,
  "data": {
    "is_favorited": true,  // ✅ 兼容旧版
    "favorited": true      // ✅ 兼容新版
  }
}
```

---

## 🎯 前端兼容性

前端代码已经兼容两种字段名：

```javascript
// 点赞状态
isLiked.value = response.data.is_liked || response.data.liked || false

// 收藏状态
isFavorited.value = response.data.is_favorited || response.data.favorited || false
```

---

## 🚀 重启服务生效

```bash
cd /home/jeff/icey/open-source/inkspace
make dev
```

---

## ✅ 完整修复清单

1. ✅ 缓存时间改为15秒
2. ✅ 点赞操作清除缓存（4处）
3. ✅ 收藏操作清除缓存（4处）
4. ✅ API返回字段兼容（5个API）
5. ✅ 前端字段读取兼容

---

## 🧪 测试验证

### 登录后测试

1. 访问：http://127.0.0.1:3001/login
2. 登录账号：iceymoss / 123456
3. 访问文章：http://127.0.0.1:3001/blog/1
4. 点击"点赞"按钮
5. 查看按钮状态变化
6. 刷新页面，状态应该保持

### API测试

```bash
# 获取token后测试
TOKEN="your_token_here"

# 检查点赞状态
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/articles/1/is-liked

# 应该返回：
# {"code":0,"data":{"liked":true,"is_liked":true}}
```

---

## 🎉 完成

所有API字段名问题已修复，前后端完全兼容！

**重启服务后立即生效！** 🚀

