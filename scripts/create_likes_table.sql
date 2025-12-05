-- 创建点赞表
USE mysite;

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

-- 为 works 表添加 like_count 和 favorite_count 字段（如果不存在）
ALTER TABLE `works` 
ADD COLUMN IF NOT EXISTS `like_count` int DEFAULT '0',
ADD COLUMN IF NOT EXISTS `favorite_count` int DEFAULT '0';

SELECT 'Likes table created successfully!' as message;

