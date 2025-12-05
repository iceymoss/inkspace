-- 修复 comments 表结构以支持作品评论
USE mysite;

-- 1. 临时禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 2. 删除旧的外键约束
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_article;
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_work;

-- 3. 修改 article_id 允许为 NULL
ALTER TABLE comments MODIFY COLUMN article_id BIGINT UNSIGNED NULL DEFAULT NULL;

-- 4. 添加 work_id 字段（如果不存在）
SELECT COUNT(*) INTO @col_exists 
FROM information_schema.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'comments'
  AND COLUMN_NAME = 'work_id';

SET @sql = IF(@col_exists = 0,
    'ALTER TABLE comments ADD COLUMN work_id BIGINT UNSIGNED NULL DEFAULT NULL AFTER article_id',
    'SELECT "work_id already exists" AS info');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 5. 添加索引
CREATE INDEX IF NOT EXISTS idx_work_id ON comments(work_id);

-- 6. 重新添加外键约束（允许 NULL）
ALTER TABLE comments 
ADD CONSTRAINT fk_comments_article 
FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE;

ALTER TABLE comments 
ADD CONSTRAINT fk_comments_work 
FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE;

-- 7. 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 8. 验证结果
SELECT '=== Comments 表结构 ===' AS info;
SELECT 
    COLUMN_NAME as 字段名,
    COLUMN_TYPE as 类型,
    IS_NULLABLE as 允许NULL,
    COLUMN_DEFAULT as 默认值
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'comments'
  AND COLUMN_NAME IN ('article_id', 'work_id');

SELECT '✅ Comments 表结构修复完成' AS result;

