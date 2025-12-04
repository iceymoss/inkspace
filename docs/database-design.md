# 数据库设计文档

## 概述

本文档详细描述了个人博客系统的数据库表结构设计，包括所有表的字段定义、索引、约束和表之间的关系。

## 数据库信息

- **数据库类型**: MySQL 8.0+
- **字符集**: utf8mb4
- **排序规则**: utf8mb4_unicode_ci
- **引擎**: InnoDB
- **时区**: UTC

## 表数量概览

| 分类 | 表名 | 数量 |
|------|------|------|
| **核心表** | users, articles, categories, tags, article_tags, comments, works | 7 |
| **扩展表** | links, settings, attachments, article_likes, comment_likes, **article_favorites**, **user_follows**, notifications, subscriptions | 9 |
| **日志表** | visit_logs, visit_log_summaries | 2 |
| **总计** | | **18张表** |

## 表结构设计

### 1. 用户表 (users)

存储系统用户信息，包括管理员和普通用户。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键，用户ID |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| username | VARCHAR | 50 | - | NOT NULL | 用户名（唯一） |
| password | VARCHAR | 255 | - | NOT NULL | 密码（bcrypt加密） |
| email | VARCHAR | 100 | - | NOT NULL | 邮箱（唯一） |
| nickname | VARCHAR | 50 | - | NULL | 昵称 |
| avatar | VARCHAR | 255 | - | NULL | 头像URL |
| bio | VARCHAR | 500 | - | NULL | 个人简介 |
| role | VARCHAR | 20 | 'user' | NOT NULL | 角色：admin/user |
| status | TINYINT | - | 1 | NOT NULL | 状态：1启用/0禁用 |
| last_login_at | DATETIME(3) | - | NULL | NULL | 最后登录时间 |
| last_login_ip | VARCHAR | 50 | - | NULL | 最后登录IP |
| article_count | INT | - | 0 | NOT NULL | 文章数量 |
| comment_count | INT | - | 0 | NOT NULL | 评论数量 |

**索引：**
- PRIMARY KEY: `id`
- UNIQUE INDEX: `idx_username` (username)
- UNIQUE INDEX: `idx_email` (email)
- INDEX: `idx_deleted_at` (deleted_at)
- INDEX: `idx_role_status` (role, status)

**约束：**
- username 长度: 3-50字符
- email 必须符合邮箱格式
- password 存储bcrypt加密后的值

---

### 2. 文章表 (articles)

存储博客文章的核心信息。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键，文章ID |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| title | VARCHAR | 200 | - | NOT NULL | 文章标题 |
| content | LONGTEXT | - | - | NOT NULL | 文章内容（Markdown） |
| content_html | LONGTEXT | - | - | NULL | 文章HTML（可选预渲染） |
| summary | VARCHAR | 500 | - | NULL | 文章摘要 |
| cover | VARCHAR | 255 | - | NULL | 封面图URL |
| category_id | BIGINT UNSIGNED | - | - | NULL | 分类ID |
| author_id | BIGINT UNSIGNED | - | - | NOT NULL | 作者ID |
| view_count | INT | - | 0 | NOT NULL | 浏览次数 |
| like_count | INT | - | 0 | NOT NULL | 点赞次数 |
| comment_count | INT | - | 0 | NOT NULL | 评论数量 |
| word_count | INT | - | 0 | NOT NULL | 字数统计 |
| reading_time | INT | - | 0 | NOT NULL | 预计阅读时间（分钟） |
| status | TINYINT | - | 1 | NOT NULL | 状态：1发布/0草稿 |
| is_top | TINYINT | - | 0 | NOT NULL | 是否置顶：1是/0否 |
| is_recommend | TINYINT | - | 0 | NOT NULL | 是否推荐：1是/0否 |
| is_original | TINYINT | - | 1 | NOT NULL | 是否原创：1是/0否 |
| source_url | VARCHAR | 255 | - | NULL | 来源URL（转载时） |
| publish_at | DATETIME(3) | - | NULL | NULL | 发布时间 |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- INDEX: `idx_category_id` (category_id)
- INDEX: `idx_author_id` (author_id)
- INDEX: `idx_status` (status)
- INDEX: `idx_top_status_created` (is_top DESC, status, created_at DESC) - 列表查询优化
- INDEX: `idx_created_at` (created_at DESC)
- FULLTEXT INDEX: `idx_title_content` (title, content) - 全文搜索

**外键：**
- `fk_articles_category`: category_id REFERENCES categories(id) ON DELETE SET NULL
- `fk_articles_author`: author_id REFERENCES users(id) ON DELETE CASCADE

---

### 3. 分类表 (categories)

存储文章分类信息。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键，分类ID |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| name | VARCHAR | 50 | - | NOT NULL | 分类名称（唯一） |
| slug | VARCHAR | 50 | - | NULL | URL别名（唯一） |
| description | VARCHAR | 200 | - | NULL | 分类描述 |
| cover | VARCHAR | 255 | - | NULL | 分类封面图 |
| sort | INT | - | 0 | NOT NULL | 排序权重 |
| article_count | INT | - | 0 | NOT NULL | 文章数量 |
| parent_id | BIGINT UNSIGNED | - | NULL | NULL | 父分类ID（预留，支持二级分类） |

**索引：**
- PRIMARY KEY: `id`
- UNIQUE INDEX: `idx_name` (name)
- UNIQUE INDEX: `idx_slug` (slug)
- INDEX: `idx_deleted_at` (deleted_at)
- INDEX: `idx_sort` (sort DESC)
- INDEX: `idx_parent_id` (parent_id)

---

### 4. 标签表 (tags)

存储文章标签信息。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键，标签ID |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| name | VARCHAR | 50 | - | NOT NULL | 标签名称（唯一） |
| slug | VARCHAR | 50 | - | NULL | URL别名（唯一） |
| color | VARCHAR | 20 | '#409eff' | NULL | 标签颜色 |
| article_count | INT | - | 0 | NOT NULL | 文章数量 |

**索引：**
- PRIMARY KEY: `id`
- UNIQUE INDEX: `idx_name` (name)
- UNIQUE INDEX: `idx_slug` (slug)
- INDEX: `idx_deleted_at` (deleted_at)

---

### 5. 文章标签关联表 (article_tags)

多对多关系：文章和标签的关联。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| article_id | BIGINT UNSIGNED | - | - | NOT NULL | 文章ID |
| tag_id | BIGINT UNSIGNED | - | - | NOT NULL | 标签ID |

**索引：**
- PRIMARY KEY: `article_id, tag_id`
- INDEX: `idx_article_id` (article_id)
- INDEX: `idx_tag_id` (tag_id)

**外键：**
- `fk_article_tags_article`: article_id REFERENCES articles(id) ON DELETE CASCADE
- `fk_article_tags_tag`: tag_id REFERENCES tags(id) ON DELETE CASCADE

---

### 6. 评论表 (comments)

存储文章评论信息，支持树形结构（父子评论）。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键，评论ID |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| article_id | BIGINT UNSIGNED | - | - | NOT NULL | 文章ID |
| user_id | BIGINT UNSIGNED | - | - | NULL | 用户ID（登录用户） |
| parent_id | BIGINT UNSIGNED | - | NULL | NULL | 父评论ID |
| root_id | BIGINT UNSIGNED | - | NULL | NULL | 根评论ID（优化查询） |
| reply_to_id | BIGINT UNSIGNED | - | NULL | NULL | 回复的评论ID |
| content | TEXT | - | - | NOT NULL | 评论内容 |
| nickname | VARCHAR | 50 | - | NULL | 昵称（游客评论） |
| email | VARCHAR | 100 | - | NULL | 邮箱（游客评论） |
| website | VARCHAR | 200 | - | NULL | 网站（游客评论） |
| ip | VARCHAR | 50 | - | NULL | IP地址 |
| user_agent | VARCHAR | 255 | - | NULL | User Agent |
| status | TINYINT | - | 1 | NOT NULL | 状态：1通过/0待审核/-1拒绝 |
| like_count | INT | - | 0 | NOT NULL | 点赞数 |
| reply_count | INT | - | 0 | NOT NULL | 回复数 |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- INDEX: `idx_article_id` (article_id)
- INDEX: `idx_user_id` (user_id)
- INDEX: `idx_parent_id` (parent_id)
- INDEX: `idx_root_id` (root_id)
- INDEX: `idx_status_created` (status, created_at DESC)

**外键：**
- `fk_comments_article`: article_id REFERENCES articles(id) ON DELETE CASCADE
- `fk_comments_user`: user_id REFERENCES users(id) ON DELETE SET NULL
- `fk_comments_parent`: parent_id REFERENCES comments(id) ON DELETE CASCADE

---

### 7. 作品表 (works)

存储个人作品展示信息。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键，作品ID |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| title | VARCHAR | 200 | - | NOT NULL | 作品标题 |
| description | TEXT | - | - | NULL | 作品描述 |
| cover | VARCHAR | 255 | - | NULL | 封面图URL |
| images | TEXT | - | - | NULL | 图片列表（JSON数组） |
| link | VARCHAR | 255 | - | NULL | 项目链接 |
| github_url | VARCHAR | 255 | - | NULL | GitHub链接 |
| demo_url | VARCHAR | 255 | - | NULL | 演示链接 |
| tech_stack | VARCHAR | 500 | - | NULL | 技术栈（逗号分隔） |
| sort | INT | - | 0 | NOT NULL | 排序权重 |
| view_count | INT | - | 0 | NOT NULL | 浏览次数 |
| status | TINYINT | - | 1 | NOT NULL | 状态：1发布/0草稿 |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- INDEX: `idx_status_sort` (status, sort DESC)

---

### 8. 友情链接表 (links)

存储友情链接信息。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键，链接ID |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| name | VARCHAR | 100 | - | NOT NULL | 链接名称 |
| url | VARCHAR | 255 | - | NOT NULL | 链接地址 |
| logo | VARCHAR | 255 | - | NULL | 网站Logo |
| description | VARCHAR | 200 | - | NULL | 描述 |
| email | VARCHAR | 100 | - | NULL | 联系邮箱 |
| sort | INT | - | 0 | NOT NULL | 排序权重 |
| status | TINYINT | - | 1 | NOT NULL | 状态：1显示/0隐藏 |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- INDEX: `idx_status_sort` (status, sort DESC)

---

### 9. 系统配置表 (settings)

存储系统配置信息。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键 |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| key | VARCHAR | 100 | - | NOT NULL | 配置键（唯一） |
| value | TEXT | - | - | NULL | 配置值 |
| type | VARCHAR | 20 | 'string' | NOT NULL | 类型：string/int/bool/json |
| description | VARCHAR | 200 | - | NULL | 描述 |
| group | VARCHAR | 50 | 'general' | NOT NULL | 分组 |

**索引：**
- PRIMARY KEY: `id`
- UNIQUE INDEX: `idx_key` (key)
- INDEX: `idx_group` (group)

**预定义配置：**
- `site_name` - 网站名称
- `site_description` - 网站描述
- `site_keywords` - 网站关键词
- `site_icp` - 备案号
- `comment_audit` - 评论审核开关
- `register_enabled` - 注册开关

---

### 10. 附件表 (attachments)

存储上传的文件附件信息。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键，附件ID |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| user_id | BIGINT UNSIGNED | - | - | NOT NULL | 上传者ID |
| file_name | VARCHAR | 255 | - | NOT NULL | 原始文件名 |
| file_path | VARCHAR | 500 | - | NOT NULL | 存储路径 |
| file_size | BIGINT | - | - | NOT NULL | 文件大小（字节） |
| file_type | VARCHAR | 50 | - | NOT NULL | 文件类型 |
| mime_type | VARCHAR | 100 | - | NOT NULL | MIME类型 |
| extension | VARCHAR | 20 | - | NOT NULL | 文件扩展名 |
| width | INT | - | - | NULL | 图片宽度 |
| height | INT | - | - | NULL | 图片高度 |
| storage_type | VARCHAR | 20 | 'local' | NOT NULL | 存储类型 |
| url | VARCHAR | 500 | - | NULL | 访问URL |
| usage_count | INT | - | 0 | NOT NULL | 使用次数 |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- INDEX: `idx_user_id` (user_id)
- INDEX: `idx_file_type` (file_type)

**外键：**
- `fk_attachments_user`: user_id REFERENCES users(id) ON DELETE SET NULL

---

### 11. 文章点赞记录表 (article_likes)

记录文章点赞，防止重复点赞。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键 |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| article_id | BIGINT UNSIGNED | - | - | NOT NULL | 文章ID |
| user_id | BIGINT UNSIGNED | - | - | NOT NULL | 用户ID（0表示游客） |
| ip | VARCHAR | 50 | - | NULL | IP地址（游客点赞） |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- UNIQUE INDEX: `idx_article_user` (article_id, user_id) - 防止重复点赞
- INDEX: `idx_user_id` (user_id)
- INDEX: `idx_ip` (ip)

**外键：**
- `fk_article_likes_article`: article_id REFERENCES articles(id) ON DELETE CASCADE
- `fk_article_likes_user`: user_id REFERENCES users(id) ON DELETE CASCADE

---

### 12. 评论点赞记录表 (comment_likes)

记录评论点赞，防止重复点赞。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键 |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| comment_id | BIGINT UNSIGNED | - | - | NOT NULL | 评论ID |
| user_id | BIGINT UNSIGNED | - | - | NOT NULL | 用户ID（0表示游客） |
| ip | VARCHAR | 50 | - | NULL | IP地址（游客点赞） |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- UNIQUE INDEX: `idx_comment_user` (comment_id, user_id) - 防止重复点赞
- INDEX: `idx_user_id` (user_id)
- INDEX: `idx_ip` (ip)

**外键：**
- `fk_comment_likes_comment`: comment_id REFERENCES comments(id) ON DELETE CASCADE
- `fk_comment_likes_user`: user_id REFERENCES users(id) ON DELETE CASCADE

---

### 13. 文章收藏表 (article_favorites) 🆕

记录用户收藏的文章，实现收藏功能。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键 |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| user_id | BIGINT UNSIGNED | - | - | NOT NULL | 用户ID |
| article_id | BIGINT UNSIGNED | - | - | NOT NULL | 文章ID |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- UNIQUE INDEX: `idx_user_article` (user_id, article_id) - 防止重复收藏
- INDEX: `idx_article_id` (article_id)

**外键：**
- `fk_article_favorites_user`: user_id REFERENCES users(id) ON DELETE CASCADE
- `fk_article_favorites_article`: article_id REFERENCES articles(id) ON DELETE CASCADE

**特点：**
- 组合唯一索引防止用户重复收藏同一篇文章
- 支持快速查询用户的收藏列表
- 支持快速查询文章被收藏的次数

---

### 14. 用户关注关系表 (user_follows) 🆕

存储用户之间的关注关系，实现关注/粉丝功能。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键 |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| follower_id | BIGINT UNSIGNED | - | - | NOT NULL | 关注者ID（粉丝） |
| following_id | BIGINT UNSIGNED | - | - | NOT NULL | 被关注者ID |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- UNIQUE INDEX: `idx_follower_following` (follower_id, following_id) - 防止重复关注
- INDEX: `idx_follower_id` (follower_id)
- INDEX: `idx_following_id` (following_id)

**外键：**
- `fk_user_follows_follower`: follower_id REFERENCES users(id) ON DELETE CASCADE
- `fk_user_follows_following`: following_id REFERENCES users(id) ON DELETE CASCADE

**特点：**
- 组合唯一索引防止重复关注
- 支持快速查询关注列表（我关注的人）
- 支持快速查询粉丝列表（关注我的人）
- 支持检查互相关注状态

**关系说明：**
```
用户A 关注 用户B
follower_id = A (关注者/粉丝)
following_id = B (被关注者)

查询A的关注列表: WHERE follower_id = A
查询B的粉丝列表: WHERE following_id = B
检查A是否关注B: WHERE follower_id = A AND following_id = B
```

---

### 17. 通知表 (notifications)

存储用户通知信息。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键 |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| user_id | BIGINT UNSIGNED | - | - | NOT NULL | 接收者ID |
| from_user_id | BIGINT UNSIGNED | - | - | NULL | 发送者ID |
| type | VARCHAR | 50 | - | NOT NULL | 通知类型 |
| title | VARCHAR | 200 | - | NOT NULL | 标题 |
| content | TEXT | - | - | NULL | 内容 |
| target_type | VARCHAR | 50 | - | NULL | 目标类型 |
| target_id | BIGINT UNSIGNED | - | - | NULL | 目标ID |
| link | VARCHAR | 255 | - | NULL | 跳转链接 |
| is_read | TINYINT | - | 0 | NOT NULL | 是否已读 |
| read_at | DATETIME(3) | - | NULL | NULL | 阅读时间 |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_deleted_at` (deleted_at)
- INDEX: `idx_user_status` (user_id, is_read)
- INDEX: `idx_from_user_id` (from_user_id)
- INDEX: `idx_type` (type)
- INDEX: `idx_target_id` (target_id)

**通知类型：**
- `comment` - 评论通知
- `reply` - 回复通知
- `like` - 点赞通知
- `system` - 系统通知
- `mention` - @提及通知

**外键：**
- `fk_notifications_user`: user_id REFERENCES users(id) ON DELETE CASCADE
- `fk_notifications_from_user`: from_user_id REFERENCES users(id) ON DELETE CASCADE

---

### 18. 订阅表 (subscriptions)

存储邮件订阅信息。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键 |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |
| deleted_at | DATETIME(3) | - | NULL | NULL | 软删除时间 |
| email | VARCHAR | 100 | - | NOT NULL | 订阅邮箱（唯一） |
| status | TINYINT | - | 1 | NOT NULL | 状态：1已确认/0未确认/-1已取消 |
| token | VARCHAR | 64 | - | NULL | 确认/取消令牌（唯一） |
| ip | VARCHAR | 50 | - | NULL | 订阅IP |
| user_agent | VARCHAR | 255 | - | NULL | User Agent |
| confirm_at | DATETIME(3) | - | NULL | NULL | 确认时间 |

**索引：**
- PRIMARY KEY: `id`
- UNIQUE INDEX: `idx_email` (email)
- UNIQUE INDEX: `idx_token` (token)
- INDEX: `idx_deleted_at` (deleted_at)
- INDEX: `idx_status` (status)

---

### 19. 访问日志表 (visit_logs)

记录网站访问日志，用于统计分析。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键 |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 访问时间 |
| user_id | BIGINT UNSIGNED | - | - | NULL | 用户ID（0表示游客） |
| ip | VARCHAR | 50 | - | NULL | IP地址 |
| user_agent | VARCHAR | 500 | - | NULL | User Agent |
| path | VARCHAR | 255 | - | NULL | 访问路径 |
| method | VARCHAR | 10 | - | NULL | 请求方法 |
| referer | VARCHAR | 500 | - | NULL | 来源URL |
| duration | INT | - | - | NULL | 请求耗时（毫秒） |
| status_code | INT | - | - | NULL | HTTP状态码 |
| country | VARCHAR | 50 | - | NULL | 国家 |
| province | VARCHAR | 50 | - | NULL | 省份 |
| city | VARCHAR | 50 | - | NULL | 城市 |
| browser | VARCHAR | 50 | - | NULL | 浏览器 |
| os | VARCHAR | 50 | - | NULL | 操作系统 |
| device | VARCHAR | 50 | - | NULL | 设备类型 |

**索引：**
- PRIMARY KEY: `id`
- INDEX: `idx_created_at` (created_at)
- INDEX: `idx_user_id` (user_id)
- INDEX: `idx_ip` (ip)
- INDEX: `idx_path` (path)
- INDEX: `idx_status_code` (status_code)

**说明：**
- 数据量大时建议按月分表
- 定期归档到统计表

---

### 20. 访问统计汇总表 (visit_log_summaries)

按天汇总的访问统计数据。

| 字段名 | 类型 | 长度 | 默认值 | 允许空 | 说明 |
|--------|------|------|--------|--------|------|
| id | BIGINT UNSIGNED | - | AUTO_INCREMENT | NOT NULL | 主键 |
| date | DATE | - | - | NOT NULL | 日期（唯一） |
| pv | INT | - | 0 | NOT NULL | 页面浏览量 |
| uv | INT | - | 0 | NOT NULL | 独立访客数 |
| ip | INT | - | 0 | NOT NULL | 独立IP数 |
| new_users | INT | - | 0 | NOT NULL | 新增用户数 |
| article_view | INT | - | 0 | NOT NULL | 文章浏览量 |
| avg_duration | INT | - | 0 | NOT NULL | 平均访问时长（秒） |
| bounce_rate | DECIMAL | 5,2 | 0.00 | NOT NULL | 跳出率 |
| created_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) | NOT NULL | 创建时间 |
| updated_at | DATETIME(3) | - | CURRENT_TIMESTAMP(3) ON UPDATE | NOT NULL | 更新时间 |

**索引：**
- PRIMARY KEY: `id`
- UNIQUE INDEX: `idx_date` (date)

**说明：**
- 由定时任务每天凌晨汇总生成
- 用于统计面板展示

---

## 表关系图

### 核心表关系

```
users (用户表)
  ├── 1:N → articles (作者)
  ├── 1:N → comments (评论者)
  ├── 1:N → attachments (上传者)
  ├── 1:N → article_likes (点赞记录)
  ├── 1:N → comment_likes (点赞记录)
  ├── 1:N → article_favorites (收藏记录) 🆕
  ├── 1:N → user_follows (关注记录 - 作为关注者) 🆕
  ├── 1:N → user_follows (被关注记录 - 作为被关注者) 🆕
  ├── 1:N → notifications (接收通知)
  └── 1:N → notifications (发送通知)

articles (文章表)
  ├── N:1 → users (作者)
  ├── N:1 → categories (分类)
  ├── N:M → tags (标签，通过 article_tags)
  ├── 1:N → comments (评论)
  ├── 1:N → article_likes (点赞记录)
  └── 1:N → article_favorites (收藏记录) 🆕

categories (分类表)
  ├── 1:N → articles (文章)
  └── 1:N → categories (父子分类，自关联)

tags (标签表)
  └── N:M → articles (通过 article_tags 关联)

comments (评论表)
  ├── N:1 → articles (所属文章)
  ├── N:1 → users (评论者，可选)
  ├── 1:N → comments (父子评论，自关联)
  └── 1:N → comment_likes (点赞记录)

works (作品表)
  └── 独立表，无外键关系

links (友情链接表)
  └── 独立表，无外键关系

settings (系统配置表)
  └── 独立表，无外键关系

attachments (附件表)
  └── N:1 → users (上传者)

article_likes (文章点赞表)
  ├── N:1 → articles (文章)
  └── N:1 → users (用户)

comment_likes (评论点赞表)
  ├── N:1 → comments (评论)
  └── N:1 → users (用户)

article_favorites (文章收藏表) 🆕
  ├── N:1 → users (用户)
  └── N:1 → articles (文章)

user_follows (用户关注表) 🆕
  ├── N:1 → users (关注者)
  └── N:1 → users (被关注者，自关联)

notifications (通知表)
  ├── N:1 → users (接收者)
  └── N:1 → users (发送者)

subscriptions (订阅表)
  └── 独立表，无外键关系

visit_logs (访问日志表)
  └── 独立表，记录访问信息

visit_log_summaries (访问统计汇总表)
  └── 独立表，汇总统计数据
```

### 表分类

**1. 核心业务表（7张）**
- users, articles, categories, tags, article_tags, comments, works

**2. 扩展功能表（9张）**
- links, settings, attachments, article_likes, comment_likes, **article_favorites**, **user_follows**, notifications, subscriptions

**3. 日志统计表（2张）**
- visit_logs, visit_log_summaries

**总计：18张表**

## 数据统计字段说明

为了提高查询性能，部分表使用了冗余的统计字段：

1. **users.article_count**: 用户发布的文章总数
2. **users.comment_count**: 用户的评论总数
3. **articles.comment_count**: 文章的评论总数
4. **categories.article_count**: 分类下的文章总数
5. **tags.article_count**: 标签下的文章总数
6. **comments.reply_count**: 评论的回复数

**维护方式：**
- 通过应用层代码维护（推荐）
- 使用数据库触发器维护
- 定期通过定时任务校准

## 性能优化建议

### 1. 索引优化
- ✅ 为常用查询条件添加索引
- ✅ 组合索引遵循最左前缀原则
- ✅ 避免过多索引影响写入性能

### 2. 查询优化
- 文章列表查询使用 `idx_top_status_created` 组合索引
- 分页查询使用 LIMIT + OFFSET
- 大数据量时考虑使用游标分页

### 3. 缓存策略
- 文章详情：缓存1小时
- 文章列表：缓存15分钟
- 分类标签：缓存30分钟
- 热门文章：缓存6小时

### 4. 分表策略（可选）
当数据量超过百万级时考虑：
- 文章表按年份分表
- 评论表按文章ID哈希分表

## 数据迁移和版本控制

### 版本管理
建议使用数据库迁移工具管理表结构变更：
- golang-migrate
- goose
- 自定义迁移脚本

### 备份策略
1. 每日全量备份
2. 增量备份（binlog）
3. 定期测试恢复流程

## 安全考虑

1. **密码安全**: 使用 bcrypt 加密，cost=10
2. **SQL注入**: 使用 GORM 参数化查询
3. **敏感信息**: 不在日志中输出密码、token等
4. **软删除**: 重要数据使用软删除而非物理删除
5. **数据加密**: 敏感字段考虑加密存储

## 数据字典导出

可以使用以下SQL导出完整的数据字典：

```sql
SELECT 
  TABLE_NAME as '表名',
  COLUMN_NAME as '字段名',
  COLUMN_TYPE as '数据类型',
  IS_NULLABLE as '允许空',
  COLUMN_DEFAULT as '默认值',
  COLUMN_COMMENT as '备注'
FROM 
  INFORMATION_SCHEMA.COLUMNS
WHERE 
  TABLE_SCHEMA = 'mysite'
ORDER BY 
  TABLE_NAME, ORDINAL_POSITION;
```

## 更新日志

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0.0 | 2024-01-01 | 初始版本 |

---

**文档维护**: 数据库结构变更时需同步更新此文档

