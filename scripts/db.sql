-- ============================================
-- My Site 数据库表结构定义 (DDL)
-- ============================================
-- 
-- 说明：
-- 1. 本项目使用 GORM AutoMigrate 自动生成表结构
-- 2. 此文件仅作为参考文档，展示所有表的DDL定义
-- 3. 实际表结构由 internal/database/migrate.go 中的 AutoMigrate 自动创建
-- 4. 如需查看实际表结构，请运行: SHOW CREATE TABLE table_name;
--
-- 使用方法：
-- 1. 开发环境：直接运行程序，GORM会自动创建表结构
-- 2. 生产环境：确保数据库已创建，然后运行程序进行迁移
-- ============================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 核心表
-- ============================================

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(100) NOT NULL,
  `nickname` varchar(50) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `bio` varchar(500) DEFAULT NULL,
  `role` varchar(20) DEFAULT 'user',
  `status` int DEFAULT '1',
  `last_login_at` datetime(3) DEFAULT NULL,
  `last_login_ip` varchar(50) DEFAULT NULL,
  `article_count` int NOT NULL DEFAULT '0',
  `work_count` int NOT NULL DEFAULT '0',
  `comment_count` int NOT NULL DEFAULT '0',
  `following_count` int NOT NULL DEFAULT '0',
  `follower_count` int NOT NULL DEFAULT '0',
  `favorite_count` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_role_status` (`role`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文章表
CREATE TABLE IF NOT EXISTS `articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(200) NOT NULL,
  `content` longtext NOT NULL,
  `content_html` longtext,
  `summary` varchar(500) DEFAULT NULL,
  `cover` varchar(255) DEFAULT NULL,
  `category_id` bigint unsigned DEFAULT NULL,
  `author_id` bigint unsigned NOT NULL,
  `view_count` int NOT NULL DEFAULT '0',
  `like_count` int NOT NULL DEFAULT '0',
  `comment_count` int NOT NULL DEFAULT '0',
  `favorite_count` int NOT NULL DEFAULT '0',
  `word_count` int NOT NULL DEFAULT '0',
  `reading_time` int NOT NULL DEFAULT '0',
  `status` int DEFAULT '1',
  `is_top` tinyint(1) DEFAULT '0',
  `is_recommend` tinyint(1) DEFAULT '0',
  `is_original` tinyint(1) DEFAULT '1',
  `source_url` varchar(255) DEFAULT NULL,
  `publish_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_articles_deleted_at` (`deleted_at`),
  KEY `idx_title` (`title`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_author_id` (`author_id`),
  KEY `idx_status` (`status`),
  KEY `idx_top_status_created` (`is_top` DESC,`status`,`created_at` DESC),
  FULLTEXT KEY `idx_articles_title_content` (`title`,`content`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 分类表
CREATE TABLE IF NOT EXISTS `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(50) NOT NULL,
  `slug` varchar(50) DEFAULT NULL,
  `description` varchar(200) DEFAULT NULL,
  `logo` varchar(255) NOT NULL,
  `cover` varchar(255) DEFAULT NULL,
  `sort` int DEFAULT '0',
  `article_count` int NOT NULL DEFAULT '0',
  `parent_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_categories_name` (`name`),
  UNIQUE KEY `idx_categories_slug` (`slug`),
  KEY `idx_categories_deleted_at` (`deleted_at`),
  KEY `idx_sort` (`sort`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 标签表
CREATE TABLE IF NOT EXISTS `tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(50) NOT NULL,
  `slug` varchar(50) DEFAULT NULL,
  `color` varchar(20) DEFAULT '#409eff',
  `article_count` int NOT NULL DEFAULT '0',
  `user_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_tag_name` (`user_id`,`name`),
  KEY `idx_tags_deleted_at` (`deleted_at`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文章标签关联表
CREATE TABLE IF NOT EXISTS `article_tags` (
  `article_id` bigint unsigned NOT NULL,
  `tag_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`article_id`,`tag_id`),
  KEY `fk_article_tags_tag` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 评论表
CREATE TABLE IF NOT EXISTS `comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `article_id` bigint unsigned DEFAULT NULL,
  `work_id` bigint unsigned DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `content` text NOT NULL,
  `parent_id` bigint unsigned DEFAULT NULL,
  `root_id` bigint unsigned DEFAULT NULL,
  `reply_to_id` bigint unsigned DEFAULT NULL,
  `nickname` varchar(50) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `website` varchar(200) DEFAULT NULL,
  `ip` varchar(50) DEFAULT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `status` int DEFAULT '1',
  `like_count` int NOT NULL DEFAULT '0',
  `reply_count` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_comments_deleted_at` (`deleted_at`),
  KEY `idx_article_id` (`article_id`),
  KEY `idx_work_id` (`work_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_root_id` (`root_id`),
  KEY `idx_status_created` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 作品表
CREATE TABLE IF NOT EXISTS `works` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(200) NOT NULL,
  `type` varchar(50) NOT NULL DEFAULT 'project',
  `metadata` text,
  `daily_quota` tinyint(1) DEFAULT '0',
  `description` text,
  `cover` varchar(255) DEFAULT NULL,
  `images` text,
  `link` varchar(255) DEFAULT NULL,
  `github_url` varchar(255) DEFAULT NULL,
  `demo_url` varchar(255) DEFAULT NULL,
  `tech_stack` varchar(500) DEFAULT NULL,
  `author_id` bigint unsigned NOT NULL DEFAULT '1',
  `sort` int DEFAULT '0',
  `view_count` int NOT NULL DEFAULT '0',
  `comment_count` int NOT NULL DEFAULT '0',
  `like_count` int NOT NULL DEFAULT '0',
  `favorite_count` int NOT NULL DEFAULT '0',
  `status` int DEFAULT '1',
  `is_recommend` tinyint(1) DEFAULT '0',
  `audit_message` text,
  PRIMARY KEY (`id`),
  KEY `idx_works_deleted_at` (`deleted_at`),
  KEY `idx_type` (`type`),
  KEY `idx_author_id` (`author_id`),
  KEY `idx_status_sort` (`status`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 互动功能表
-- ============================================

-- 点赞表
CREATE TABLE IF NOT EXISTS `likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `article_id` bigint unsigned DEFAULT NULL,
  `work_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_likes_deleted_at` (`deleted_at`),
  KEY `idx_user_target` (`user_id`,`article_id`,`work_id`),
  KEY `idx_article` (`article_id`),
  KEY `idx_work` (`work_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 收藏表（统一收藏表，支持文章和作品）
CREATE TABLE IF NOT EXISTS `favorites` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `article_id` bigint unsigned DEFAULT NULL,
  `work_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_favorites_deleted_at` (`deleted_at`),
  KEY `idx_user_target` (`user_id`,`article_id`,`work_id`),
  KEY `idx_article` (`article_id`),
  KEY `idx_work` (`work_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文章收藏表（兼容旧版本）
CREATE TABLE IF NOT EXISTS `article_favorites` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `article_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_article_favorites_deleted_at` (`deleted_at`),
  KEY `idx_user_article` (`user_id`,`article_id`),
  KEY `idx_article_id` (`article_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户关注表
CREATE TABLE IF NOT EXISTS `user_follows` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `follower_id` bigint unsigned NOT NULL,
  `following_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_follows_deleted_at` (`deleted_at`),
  KEY `idx_follower_following` (`follower_id`,`following_id`),
  KEY `idx_follower_id` (`follower_id`),
  KEY `idx_following_id` (`following_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 通知表
CREATE TABLE IF NOT EXISTS `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `from_user_id` bigint unsigned DEFAULT NULL,
  `type` varchar(50) NOT NULL,
  `content` text,
  `article_id` bigint unsigned DEFAULT NULL,
  `work_id` bigint unsigned DEFAULT NULL,
  `comment_id` bigint unsigned DEFAULT NULL,
  `is_read` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_notifications_deleted_at` (`deleted_at`),
  KEY `idx_notifications_user_id` (`user_id`),
  KEY `idx_notifications_from_user_id` (`from_user_id`),
  KEY `idx_notifications_article_id` (`article_id`),
  KEY `idx_notifications_work_id` (`work_id`),
  KEY `idx_notifications_comment_id` (`comment_id`),
  KEY `idx_notifications_is_read` (`is_read`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 系统功能表
-- ============================================

-- 友情链接表
CREATE TABLE IF NOT EXISTS `links` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `url` varchar(255) NOT NULL,
  `logo` varchar(255) DEFAULT NULL,
  `description` varchar(200) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `sort` int DEFAULT '0',
  `status` int DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `idx_links_deleted_at` (`deleted_at`),
  KEY `idx_status_sort` (`status`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 系统配置表
CREATE TABLE IF NOT EXISTS `settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `key` varchar(100) NOT NULL,
  `value` text,
  `type` varchar(20) DEFAULT 'string',
  `description` varchar(200) DEFAULT NULL,
  `group` varchar(50) DEFAULT 'general',
  `is_public` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_settings_key` (`key`),
  KEY `idx_settings_deleted_at` (`deleted_at`),
  KEY `idx_settings_group` (`group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 附件表
CREATE TABLE IF NOT EXISTS `attachments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `file_name` varchar(255) NOT NULL,
  `file_path` varchar(500) NOT NULL,
  `file_size` bigint NOT NULL,
  `file_type` varchar(50) NOT NULL,
  `mime_type` varchar(100) NOT NULL,
  `extension` varchar(20) NOT NULL,
  `width` int DEFAULT NULL,
  `height` int DEFAULT NULL,
  `storage_type` varchar(20) DEFAULT 'local',
  `url` varchar(500) DEFAULT NULL,
  `usage_count` int DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_attachments_deleted_at` (`deleted_at`),
  KEY `idx_attachments_user_id` (`user_id`),
  KEY `idx_attachments_file_type` (`file_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 订阅表
CREATE TABLE IF NOT EXISTS `subscriptions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `email` varchar(100) NOT NULL,
  `status` int DEFAULT '1',
  `token` varchar(64) DEFAULT NULL,
  `ip` varchar(50) DEFAULT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `confirm_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_subscriptions_email` (`email`),
  UNIQUE KEY `idx_subscriptions_token` (`token`),
  KEY `idx_subscriptions_deleted_at` (`deleted_at`),
  KEY `idx_subscriptions_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 广告位置表
CREATE TABLE IF NOT EXISTS `ad_positions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `code` varchar(50) NOT NULL,
  `description` text,
  `max_count` int DEFAULT '4',
  `status` int DEFAULT '1',
  `sort` int DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_ad_positions_code` (`code`),
  KEY `idx_ad_positions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 广告表
CREATE TABLE IF NOT EXISTS `advertisements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(200) NOT NULL,
  `image` varchar(500) DEFAULT NULL,
  `link` varchar(500) DEFAULT NULL,
  `description` text,
  `status` int DEFAULT '1',
  `click_count` int DEFAULT '0',
  `view_count` int DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_advertisements_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 广告投放表
CREATE TABLE IF NOT EXISTS `ad_placements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `ad_position_id` bigint unsigned NOT NULL,
  `advertisement_id` bigint unsigned NOT NULL,
  `sort` int DEFAULT '0',
  `status` int DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `idx_ad_placements_deleted_at` (`deleted_at`),
  KEY `idx_ad_placements_position_id` (`ad_position_id`),
  KEY `idx_ad_placements_advertisement_id` (`advertisement_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 日志统计表
-- ============================================

-- 访问日志表
CREATE TABLE IF NOT EXISTS `visit_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `ip` varchar(50) DEFAULT NULL,
  `user_agent` varchar(500) DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `method` varchar(10) DEFAULT NULL,
  `referer` varchar(500) DEFAULT NULL,
  `duration` int DEFAULT NULL,
  `status_code` int DEFAULT NULL,
  `country` varchar(50) DEFAULT NULL,
  `province` varchar(50) DEFAULT NULL,
  `city` varchar(50) DEFAULT NULL,
  `browser` varchar(50) DEFAULT NULL,
  `os` varchar(50) DEFAULT NULL,
  `device` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_visit_logs_created_at` (`created_at`),
  KEY `idx_visit_logs_user_id` (`user_id`),
  KEY `idx_visit_logs_ip` (`ip`),
  KEY `idx_visit_logs_path` (`path`),
  KEY `idx_visit_logs_status_code` (`status_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 访问统计汇总表
CREATE TABLE IF NOT EXISTS `visit_log_summaries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `date` date NOT NULL,
  `pv` int NOT NULL DEFAULT '0',
  `uv` int NOT NULL DEFAULT '0',
  `ip` int NOT NULL DEFAULT '0',
  `new_users` int NOT NULL DEFAULT '0',
  `article_view` int NOT NULL DEFAULT '0',
  `avg_duration` int NOT NULL DEFAULT '0',
  `bounce_rate` decimal(5,2) DEFAULT '0.00',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 说明
-- ============================================
-- 
-- 1. 此文件仅作为参考文档，实际表结构由 GORM AutoMigrate 自动生成
-- 2. 表结构定义基于 internal/models 目录下的模型文件
-- 3. 索引和外键约束由 internal/database/migrate.go 中的 createIndexes() 和 createForeignKeys() 创建
-- 4. 如需查看实际表结构，请运行: SHOW CREATE TABLE table_name;
-- 5. 初始数据请参考 scripts/init.sql
-- 
-- ============================================

