# 🔍 通知功能调试指南

## 问题描述
用户A评论/点赞/收藏用户B的作品，用户B没有收到通知。

---

## 🚀 快速修复步骤

### 步骤1：确保数据库表已创建

```bash
cd /home/jeff/icey/open-source/inkspace

# 执行完整部署脚本
mysql -h localhost -u root -proot mysite < scripts/complete_deployment.sql
```

### 步骤2：重启后端服务

```bash
# 停止现有服务（Ctrl+C）
# 然后重新启动
make dev
```

### 步骤3：测试通知功能

#### 测试评论通知
1. 用户A登录
2. 访问用户B的作品：http://127.0.0.1:3001/works/1
3. 发表评论
4. 查看后端日志，应该看到：
   ```
   ✅ 成功创建作品评论通知: 用户A -> 用户B, 作品1
   ```
5. 用户B登录
6. 访问通知中心：http://127.0.0.1:3001/dashboard/notifications
7. 应该看到评论通知

#### 测试点赞通知
1. 用户A登录
2. 访问用户B的作品
3. 点击"点赞"按钮
4. 查看后端日志，应该看到：
   ```
   ✅ 成功创建作品点赞通知: 用户A -> 用户B, 作品1
   ```
5. 用户B登录查看通知

#### 测试收藏通知
1. 用户A登录
2. 访问用户B的作品
3. 点击"收藏"按钮
4. 查看后端日志，应该看到：
   ```
   ✅ 成功创建作品收藏通知: 用户A -> 用户B, 作品1
   ```
5. 用户B登录查看通知

---

## 🔍 详细排查步骤

### 1. 检查数据库表

```sql
USE mysite;

-- 检查通知表是否存在
SHOW TABLES LIKE 'notifications';

-- 检查表结构
DESC notifications;

-- 查看所有通知
SELECT * FROM notifications ORDER BY created_at DESC LIMIT 10;

-- 统计通知数量
SELECT 
    type,
    COUNT(*) as count,
    SUM(CASE WHEN is_read = 0 THEN 1 ELSE 0 END) as unread_count
FROM notifications
GROUP BY type;
```

**预期结果：**
- 表应该存在
- 应该有通知记录
- 通知类型应该包括 'comment', 'like', 'favorite'

### 2. 检查后端日志

启动服务后，观察日志输出：

```bash
make dev
```

当用户A评论用户B的作品时，应该看到类似的日志：
```
✅ 成功创建作品评论通知: 用户2 -> 用户1, 作品1
```

**如果看到错误日志：**
```
❌ 创建作品评论通知失败: ...
❌ 获取作品信息失败: ...
```

说明有问题，根据错误信息进一步排查。

### 3. 测试API

```bash
# 替换 YOUR_TOKEN 为实际的JWT token
TOKEN="eyJhbGc..."

# 获取通知列表
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/notifications

# 获取未读数量
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/notifications/unread-count

# 创建测试评论（触发通知）
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"work_id": 1, "content": "测试评论"}' \
  http://localhost:8081/api/comments
```

### 4. 检查前端

打开浏览器开发者工具（F12）：

#### Network标签
- 查看 `/api/notifications` 请求
- 检查响应数据
- 确认status code是200

#### Console标签
- 查看是否有JavaScript错误
- 检查API调用日志

---

## 🐛 常见问题和解决方案

### 问题1：通知表不存在

**症状：**
```
Error: Table 'mysite.notifications' doesn't exist
```

**解决：**
```bash
mysql -h localhost -u root -proot mysite < scripts/create_notifications_table.sql
```

### 问题2：外键约束失败

**症状：**
```
Error 1452: Cannot add or update a child row: a foreign key constraint fails
```

**解决：**
```sql
USE mysite;

-- 检查相关用户和作品是否存在
SELECT id FROM users WHERE id IN (1, 2);
SELECT id FROM works WHERE id = 1;

-- 如果不存在，需要创建测试数据
```

### 问题3：通知服务未初始化

**症状：** 后端日志没有任何通知相关输出

**检查代码：**
```go
// CommentService, LikeService, FavoriteService 
// 都应该有 notificationService 成员变量

type CommentService struct{
    notificationService *NotificationService
}

func NewCommentService() *CommentService {
    return &CommentService{
        notificationService: NewNotificationService(),
    }
}
```

### 问题4：goroutine 错误被吞噬

**症状：** 通知创建失败但没有日志

**解决：** 已添加详细日志记录，重启服务后可以看到详细信息

---

## 📊 验证清单

- [ ] 数据库表已创建
- [ ] 服务已重启
- [ ] 可以看到日志输出
- [ ] API返回正确数据
- [ ] 前端显示通知
- [ ] 未读数量正确
- [ ] 点击通知可跳转

---

## 🎯 测试场景

### 场景1：评论作品
1. 用户A（id=2）登录
2. 访问用户B（id=1）的作品
3. 发表评论
4. **后端日志应显示：** `✅ 成功创建作品评论通知: 用户2 -> 用户1, 作品X`
5. 用户B登录
6. 导航栏应显示未读数量
7. 访问通知中心，应该看到通知

### 场景2：点赞作品
1. 用户A登录
2. 访问用户B的作品
3. 点击"点赞"
4. **后端日志应显示：** `✅ 成功创建作品点赞通知: 用户2 -> 用户1, 作品X`
5. 用户B查看通知

### 场景3：收藏作品
1. 用户A登录
2. 访问用户B的作品
3. 点击"收藏"
4. **后端日志应显示：** `✅ 成功创建作品收藏通知: 用户2 -> 用户1, 作品X`
5. 用户B查看通知

---

## 💡 调试技巧

### 1. 查看实时日志

```bash
# 启动服务并实时查看日志
make dev | grep "通知"
```

### 2. 查询数据库

```sql
-- 查看最新的通知
SELECT 
    n.id,
    n.type,
    n.content,
    u1.username as '接收者',
    u2.username as '发送者',
    n.work_id,
    n.created_at
FROM notifications n
LEFT JOIN users u1 ON n.user_id = u1.id
LEFT JOIN users u2 ON n.from_user_id = u2.id
ORDER BY n.created_at DESC
LIMIT 10;
```

### 3. 使用 Postman 测试API

导入以下请求：
- POST /api/comments（创建评论）
- POST /api/works/:id/like（点赞）
- POST /api/works/:id/favorite（收藏）
- GET /api/notifications（获取通知）

---

## 📝 预期行为

### 文章通知（已正常）
- ✅ 评论文章 → 文章作者收到通知
- ✅ 点赞文章 → 文章作者收到通知
- ✅ 收藏文章 → 文章作者收到通知

### 作品通知（应该正常）
- ✅ 评论作品 → 作品作者收到通知
- ✅ 点赞作品 → 作品作者收到通知
- ✅ 收藏作品 → 作品作者收到通知

---

## 🎊 如果问题持续

请提供以下信息：
1. 后端日志输出
2. 数据库通知表的内容
3. API响应数据
4. 浏览器控制台错误

这样可以更准确地定位问题！

