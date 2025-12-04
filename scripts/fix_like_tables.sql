-- ============================================
-- 修复点赞表的外键约束问题
-- 将 user_id 改为允许 NULL，以支持游客点赞
-- ============================================

USE `mysite`;

-- 修改 article_likes 表
-- 1. 删除外键约束
ALTER TABLE `article_likes` 
DROP FOREIGN KEY IF EXISTS `fk_article_likes_user`;

-- 2. 修改 user_id 字段，允许 NULL
ALTER TABLE `article_likes` 
MODIFY COLUMN `user_id` BIGINT UNSIGNED NULL;

-- 3. 重新添加外键约束（允许 NULL）
ALTER TABLE `article_likes`
ADD CONSTRAINT `fk_article_likes_user` 
FOREIGN KEY (`user_id`) 
REFERENCES `users` (`id`) 
ON DELETE CASCADE;

-- 修改 comment_likes 表
-- 1. 删除外键约束
ALTER TABLE `comment_likes` 
DROP FOREIGN KEY IF EXISTS `fk_comment_likes_user`;

-- 2. 修改 user_id 字段，允许 NULL
ALTER TABLE `comment_likes` 
MODIFY COLUMN `user_id` BIGINT UNSIGNED NULL;

-- 3. 重新添加外键约束（允许 NULL）
ALTER TABLE `comment_likes`
ADD CONSTRAINT `fk_comment_likes_user` 
FOREIGN KEY (`user_id`) 
REFERENCES `users` (`id`) 
ON DELETE CASCADE;

-- 更新现有数据：将 user_id = 0 的记录改为 NULL
UPDATE `article_likes` SET `user_id` = NULL WHERE `user_id` = 0;
UPDATE `comment_likes` SET `user_id` = NULL WHERE `user_id` = 0;

-- 验证修改
SELECT 
    TABLE_NAME,
    COLUMN_NAME,
    IS_NULLABLE,
    COLUMN_DEFAULT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'mysite' 
  AND TABLE_NAME IN ('article_likes', 'comment_likes')
  AND COLUMN_NAME = 'user_id';

