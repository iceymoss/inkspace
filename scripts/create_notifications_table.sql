-- 创建通知表
USE mysite;

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

SELECT 'Notifications table created successfully!' as message;

