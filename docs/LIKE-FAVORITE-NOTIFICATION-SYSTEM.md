# 点赞、收藏和消息通知系统

## 🎯 功能概述

### 1. 点赞功能
- ✅ 文章点赞/取消点赞
- ✅ 作品点赞/取消点赞
- ✅ 评论点赞/取消点赞
- ✅ 实时更新点赞数量
- ✅ 点赞状态持久化

### 2. 收藏功能
- ✅ 文章收藏/取消收藏
- ✅ 作品收藏/取消收藏
- ✅ 实时更新收藏数量
- ✅ 收藏状态持久化

### 3. 消息通知
- ✅ 评论通知（文章/作品作者）
- ✅ 点赞通知（文章/作品作者）
- ✅ 收藏通知（文章/作品作者）
- ✅ 关注通知
- ✅ 回复通知

---

## 📊 数据库设计

### Likes 表

```sql
CREATE TABLE `likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `article_id` bigint unsigned DEFAULT NULL,
  `work_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_target` (`user_id`,`article_id`,`work_id`),
  KEY `idx_article` (`article_id`),
  KEY `idx_work` (`work_id`),
  KEY `idx_likes_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_likes_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_likes_work` FOREIGN KEY (`work_id`) REFERENCES `works` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);
```

### Favorites 表（已存在，需扩展）

```sql
ALTER TABLE `favorites` ADD COLUMN `work_id` bigint unsigned DEFAULT NULL;
ALTER TABLE `favorites` ADD KEY `idx_work` (`work_id`);
ALTER TABLE `favorites` ADD CONSTRAINT `fk_favorites_work` FOREIGN KEY (`work_id`) REFERENCES `works` (`id`) ON DELETE CASCADE;
```

### Notifications 表

```sql
CREATE TABLE `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL COMMENT '接收通知的用户',
  `from_user_id` bigint unsigned NOT NULL COMMENT '触发通知的用户',
  `type` varchar(50) NOT NULL COMMENT 'comment/like/favorite/follow/reply',
  `content` text COMMENT '通知内容',
  `article_id` bigint unsigned DEFAULT NULL,
  `work_id` bigint unsigned DEFAULT NULL,
  `comment_id` bigint unsigned DEFAULT NULL,
  `is_read` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_notifications_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_notifications_from_user` FOREIGN KEY (`from_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);
```

---

## 🔌 API 端点

### 点赞 API

```
POST   /api/works/:id/like          # 点赞/取消点赞作品
POST   /api/articles/:id/like       # 点赞/取消点赞文章
GET    /api/works/:id/liked         # 检查是否已点赞作品
GET    /api/articles/:id/liked      # 检查是否已点赞文章
```

### 收藏 API

```
POST   /api/works/:id/favorite      # 收藏作品
DELETE /api/works/:id/favorite      # 取消收藏作品
GET    /api/works/:id/favorited     # 检查是否已收藏作品
POST   /api/articles/:id/favorite   # 收藏文章
DELETE /api/articles/:id/favorite   # 取消收藏文章
```

### 通知 API

```
GET    /api/notifications           # 获取通知列表
GET    /api/notifications/unread    # 获取未读通知数量
PUT    /api/notifications/:id/read  # 标记通知为已读
PUT    /api/notifications/read-all  # 标记所有通知为已读
DELETE /api/notifications/:id       # 删除通知
```

---

## 🎨 前端实现

### WorkDetail.vue - 点赞和收藏按钮

```vue
<template>
  <div class="action-buttons">
    <el-button 
      :type="isLiked ? 'primary' : 'default'"
      @click="handleLike"
      :loading="liking"
    >
      <el-icon><Star /></el-icon>
      {{ work.like_count }} {{ isLiked ? '已点赞' : '点赞' }}
    </el-button>
    
    <el-button 
      :type="isFavorited ? 'warning' : 'default'"
      @click="handleFavorite"
      :loading="favoriting"
    >
      <el-icon><Star /></el-icon>
      {{ work.favorite_count }} {{ isFavorited ? '已收藏' : '收藏' }}
    </el-button>
  </div>
</template>

<script setup>
const isLiked = ref(false)
const isFavorited = ref(false)
const liking = ref(false)
const favoriting = ref(false)

const checkLikedStatus = async () => {
  if (!userStore.isLoggedIn) return
  try {
    const response = await api.get(`/works/${route.params.id}/liked`)
    isLiked.value = response.data.liked
  } catch (error) {
    console.error('Failed to check liked status:', error)
  }
}

const checkFavoritedStatus = async () => {
  if (!userStore.isLoggedIn) return
  try {
    const response = await api.get(`/works/${route.params.id}/favorited`)
    isFavorited.value = response.data.favorited
  } catch (error) {
    console.error('Failed to check favorited status:', error)
  }
}

const handleLike = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  liking.value = true
  try {
    await api.post(`/works/${route.params.id}/like`)
    isLiked.value = !isLiked.value
    work.value.like_count += isLiked.value ? 1 : -1
    ElMessage.success(isLiked.value ? '点赞成功' : '取消点赞')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    liking.value = false
  }
}

const handleFavorite = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  favoriting.value = true
  try {
    if (isFavorited.value) {
      await api.delete(`/works/${route.params.id}/favorite`)
    } else {
      await api.post(`/works/${route.params.id}/favorite`)
    }
    isFavorited.value = !isFavorited.value
    work.value.favorite_count += isFavorited.value ? 1 : -1
    ElMessage.success(isFavorited.value ? '收藏成功' : '取消收藏')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    favoriting.value = false
  }
}

onMounted(() => {
  loadWork()
  checkLikedStatus()
  checkFavoritedStatus()
})
</script>
```

---

## 🔔 消息通知触发时机

### 1. 评论通知
```go
// 在 comment_service.go 的 Create 方法中
func (s *CommentService) Create(req *CommentRequest) (*Comment, error) {
    // ... 创建评论 ...
    
    // 发送通知给作品/文章作者
    if req.WorkID != nil {
        work, _ := workService.GetByID(*req.WorkID)
        if work.AuthorID != req.UserID {
            notificationService.CreateNotification(
                work.AuthorID,
                req.UserID,
                "comment",
                "评论了你的作品",
                nil,
                req.WorkID,
                &comment.ID,
            )
        }
    }
    
    return comment, nil
}
```

### 2. 点赞通知
```go
// 在 like_service.go 的 LikeWork 方法中
func (s *LikeService) LikeWork(userID, workID uint) error {
    // ... 点赞逻辑 ...
    
    // 发送通知给作品作者
    work, _ := workService.GetByID(workID)
    if work.AuthorID != userID {
        notificationService.CreateNotification(
            work.AuthorID,
            userID,
            "like",
            "点赞了你的作品",
            nil,
            &workID,
            nil,
        )
    }
    
    return nil
}
```

### 3. 收藏通知
```go
// 在 favorite_service.go 的 AddWorkFavorite 方法中
func (s *FavoriteService) AddWorkFavorite(userID, workID uint) error {
    // ... 收藏逻辑 ...
    
    // 发送通知给作品作者
    work, _ := workService.GetByID(workID)
    if work.AuthorID != userID {
        notificationService.CreateNotification(
            work.AuthorID,
            userID,
            "favorite",
            "收藏了你的作品",
            nil,
            &workID,
            nil,
        )
    }
    
    return nil
}
```

---

## 📝 实现步骤

### 1. 数据库迁移

```bash
mysql -h localhost -u root -proot mysite < scripts/create_likes_table.sql
mysql -h localhost -u root -proot mysite < scripts/extend_favorites_table.sql
mysql -h localhost -u root -proot mysite < scripts/create_notifications_table.sql
```

### 2. 后端实现

- [x] `internal/models/like.go` - 点赞模型
- [x] `internal/service/like_service.go` - 点赞服务
- [x] `internal/handler/like_handler.go` - 点赞处理器
- [ ] `internal/models/notification.go` - 通知模型
- [ ] `internal/service/notification_service.go` - 通知服务
- [ ] `internal/handler/notification_handler.go` - 通知处理器
- [ ] 扩展 `favorite_service.go` - 添加作品收藏

### 3. 前端实现

- [ ] 更新 `WorkDetail.vue` - 添加点赞/收藏按钮
- [ ] 更新 `BlogDetail.vue` - 添加点赞/收藏按钮
- [ ] 创建 `Notifications.vue` - 通知中心
- [ ] 更新导航栏 - 添加通知图标和未读数量

---

## ✅ 完成清单

- [x] 评论区标注作者
- [x] 点赞功能（后端模型+服务+路由）
- [ ] 点赞功能（前端实现）
- [ ] 收藏功能（后端扩展）
- [ ] 收藏功能（前端实现）
- [ ] 消息通知系统（后端）
- [ ] 消息通知系统（前端）
- [ ] 优化摄影详情页样式

---

这是一个完整的实现方案，需要继续完成剩余部分。

