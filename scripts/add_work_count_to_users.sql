-- ============================================
-- 添加 work_count 字段到 users 表
-- ============================================

-- 如果字段不存在，则添加
ALTER TABLE `users` 
ADD COLUMN IF NOT EXISTS `work_count` INT NOT NULL DEFAULT 0 
AFTER `article_count`;

-- 同步现有用户的作品数
UPDATE `users` u
SET `work_count` = (
    SELECT COUNT(*) 
    FROM `works` w 
    WHERE w.author_id = u.id 
    AND w.deleted_at IS NULL 
    AND w.status = 1  -- 只统计已发布的作品
);

