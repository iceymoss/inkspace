-- 添加作品评论功能支持
-- 运行时间：2025-12-05

-- 1. 删除旧的外键约束（如果存在）
SET @fk_exists = (
    SELECT COUNT(*) 
    FROM information_schema.TABLE_CONSTRAINTS 
    WHERE CONSTRAINT_SCHEMA = DATABASE()
    AND TABLE_NAME = 'comments'
    AND CONSTRAINT_NAME = 'fk_comments_article'
    AND CONSTRAINT_TYPE = 'FOREIGN KEY'
);

SET @sql = IF(@fk_exists > 0, 
    'ALTER TABLE comments DROP FOREIGN KEY fk_comments_article', 
    'SELECT "外键不存在，跳过删除" AS info');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 2. 修改 comments 表的 article_id 允许为 NULL
ALTER TABLE comments 
MODIFY COLUMN article_id BIGINT UNSIGNED NULL DEFAULT NULL;

-- 3. 为 comments 表添加 work_id 字段
ALTER TABLE comments 
ADD COLUMN IF NOT EXISTS work_id BIGINT UNSIGNED NULL DEFAULT NULL AFTER article_id,
ADD INDEX IF NOT EXISTS idx_work_id (work_id);

-- 4. 为 works 表添加 author_id 字段
ALTER TABLE works 
ADD COLUMN IF NOT EXISTS author_id BIGINT UNSIGNED NOT NULL DEFAULT 1 AFTER tech_stack,
ADD INDEX IF NOT EXISTS idx_author_id (author_id);

-- 5. 为 works 表添加 comment_count 字段
ALTER TABLE works 
ADD COLUMN IF NOT EXISTS comment_count INT NOT NULL DEFAULT 0 AFTER view_count;

-- 4. 同步现有作品的评论数（如果有）
UPDATE works w 
SET comment_count = (
    SELECT COUNT(*) FROM comments c 
    WHERE c.work_id = w.id AND c.deleted_at IS NULL
);

-- 5. 验证
SELECT 'Comments表结构更新' AS status;
DESCRIBE comments;

SELECT 'Works表结构更新' AS status;
DESCRIBE works;

SELECT '✅ 数据库迁移完成' AS result;

