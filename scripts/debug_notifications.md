# 通知功能调试指南

## 问题描述
用户A评论用户B的作品，用户B没有收到通知。

## 排查步骤

### 1. 检查数据库表是否存在

```bash
mysql -h localhost -u root -proot mysite < scripts/test_notifications.sql
```

### 2. 手动测试通知创建

```bash
# 进入MySQL
mysql -h localhost -u root -proot mysite

# 手动插入一条测试通知
INSERT INTO notifications (user_id, from_user_id, type, content, work_id, is_read, created_at, updated_at)
VALUES (1, 2, 'comment', '评论了你的作品', 1, 0, NOW(), NOW());

# 查询通知
SELECT * FROM notifications WHERE user_id = 1;
```

### 3. 检查后端日志

观察服务器日志，看是否有通知创建的错误信息。

### 4. API测试

```bash
# 获取通知列表（替换YOUR_TOKEN为实际token）
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8081/api/notifications

# 获取未读数量
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8081/api/notifications/unread-count
```

### 5. 前端测试

1. 打开浏览器开发者工具（F12）
2. 访问通知中心：http://127.0.0.1:3001/dashboard/notifications
3. 查看Network标签，检查API请求和响应
4. 查看Console标签，查看是否有错误

---

## 可能的问题

### 问题1：通知表未创建
**检查：**
```sql
SHOW TABLES LIKE 'notifications';
```

**解决：**
```bash
mysql -h localhost -u root -proot mysite < scripts/create_notifications_table.sql
```

### 问题2：外键约束问题
**检查：**
```sql
SHOW CREATE TABLE notifications;
```

**解决：**
可能需要暂时禁用外键检查：
```sql
SET FOREIGN_KEY_CHECKS = 0;
-- 执行操作
SET FOREIGN_KEY_CHECKS = 1;
```

### 问题3：通知服务未正确初始化
**检查代码：**
- CommentService 是否正确初始化 notificationService
- LikeService 是否正确初始化 notificationService
- FavoriteService 是否正确初始化 notificationService

### 问题4：异步通知失败但未记录
**问题：** 通知创建在 goroutine 中，错误可能被忽略

**解决：** 添加日志记录

---

## 建议的修复

### 1. 添加日志记录

在通知创建时添加日志：

```go
// 发送通知给作品作者
var work models.Work
if err := database.DB.First(&work, workID).Error; err == nil {
    if work.AuthorID != userID {
        go func() {
            err := s.notificationService.CreateCommentNotification(userID, work.AuthorID, nil, workID, comment.ID)
            if err != nil {
                log.Printf("创建评论通知失败: %v", err)
            } else {
                log.Printf("成功创建评论通知: 用户%d -> 用户%d, 作品%d", userID, work.AuthorID, workID)
            }
        }()
    }
}
```

### 2. 检查数据库迁移

确保执行了所有必要的SQL脚本：
```bash
mysql -h localhost -u root -proot mysite < scripts/complete_deployment.sql
```

### 3. 重启服务

```bash
# 重启后端服务
make dev
make dev-admin
```

