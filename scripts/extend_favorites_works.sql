-- 扩展 favorites 表以支持作品收藏
USE mysite;

-- 添加 work_id 字段（如果不存在）
ALTER TABLE `favorites` 
ADD COLUMN IF NOT EXISTS `work_id` bigint unsigned DEFAULT NULL AFTER `article_id`;

-- 添加索引
ALTER TABLE `favorites` 
ADD KEY IF NOT EXISTS `idx_work` (`work_id`);

-- 添加外键约束
ALTER TABLE `favorites` 
ADD CONSTRAINT IF NOT EXISTS `fk_favorites_work` 
FOREIGN KEY (`work_id`) REFERENCES `works` (`id`) ON DELETE CASCADE;

-- 为 works 表添加收藏和点赞计数字段（如果不存在）
ALTER TABLE `works` 
ADD COLUMN IF NOT EXISTS `like_count` int DEFAULT '0',
ADD COLUMN IF NOT EXISTS `favorite_count` int DEFAULT '0';

SELECT 'Favorites table extended for works successfully!' as message;

