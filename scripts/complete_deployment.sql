-- 完整部署脚本 - 一次性执行所有数据库迁移
USE mysite;

-- ========================================
-- 1. 修复图片URL
-- ========================================
UPDATE users 
SET avatar = REPLACE(avatar, 'http://localhost:8081', '')
WHERE avatar LIKE 'http://localhost:8081%';

UPDATE works 
SET cover = REPLACE(cover, 'http://localhost:8081', '')
WHERE cover LIKE 'http://localhost:8081%';

UPDATE articles 
SET cover = REPLACE(cover, 'http://localhost:8081', '')
WHERE cover LIKE 'http://localhost:8081%';

UPDATE categories 
SET logo = REPLACE(logo, 'http://localhost:8081', '')
WHERE logo LIKE 'http://localhost:8081%';

-- ========================================
-- 2. 创建点赞表
-- ========================================
CREATE TABLE IF NOT EXISTS `likes` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ========================================
-- 3. 为 works 表添加计数字段
-- ========================================
ALTER TABLE `works` 
ADD COLUMN IF NOT EXISTS `like_count` int DEFAULT '0',
ADD COLUMN IF NOT EXISTS `favorite_count` int DEFAULT '0';

-- ========================================
-- 4. 扩展 favorites 表支持作品
-- ========================================
ALTER TABLE `favorites` 
ADD COLUMN IF NOT EXISTS `work_id` bigint unsigned DEFAULT NULL AFTER `article_id`;

-- 添加索引（如果不存在）
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics 
               WHERE table_schema = 'mysite' 
               AND table_name = 'favorites' 
               AND index_name = 'idx_work');
SET @sqlstmt := IF(@exist = 0, 
                   'ALTER TABLE favorites ADD KEY idx_work (work_id)', 
                   'SELECT ''Index already exists'' as message');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 添加外键约束（如果不存在）
SET @exist := (SELECT COUNT(*) FROM information_schema.table_constraints 
               WHERE constraint_schema = 'mysite' 
               AND table_name = 'favorites' 
               AND constraint_name = 'fk_favorites_work');
SET @sqlstmt := IF(@exist = 0, 
                   'ALTER TABLE favorites ADD CONSTRAINT fk_favorites_work FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE', 
                   'SELECT ''Foreign key already exists'' as message');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ========================================
-- 5. 创建通知表
-- ========================================
CREATE TABLE IF NOT EXISTS `notifications` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ========================================
-- 完成
-- ========================================
SELECT '✅ 所有数据库迁移已完成！' as message;
SELECT '📊 统计信息：' as '';
SELECT 'Users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Works', COUNT(*) FROM works
UNION ALL
SELECT 'Articles', COUNT(*) FROM articles
UNION ALL
SELECT 'Comments', COUNT(*) FROM comments
UNION ALL
SELECT 'Likes', COUNT(*) FROM likes
UNION ALL
SELECT 'Favorites', COUNT(*) FROM favorites
UNION ALL
SELECT 'Notifications', COUNT(*) FROM notifications;

