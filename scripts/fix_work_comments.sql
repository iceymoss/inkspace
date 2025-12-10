-- 修复作品评论功能的数据库结构
-- 运行时间：2025-12-05

USE mysite;

-- 1. 临时禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 2. 删除 comments 表的旧外键约束
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_article;

-- 3. 修改 article_id 允许为 NULL
ALTER TABLE comments MODIFY COLUMN article_id BIGINT UNSIGNED NULL;

-- 4. 添加 work_id 字段（如果不存在）
ALTER TABLE comments 
ADD COLUMN IF NOT EXISTS work_id BIGINT UNSIGNED NULL AFTER article_id;

-- 5. 添加索引（如果不存在）
CREATE INDEX IF NOT EXISTS idx_work_id ON comments(work_id);

-- 6. 添加新的外键约束（article_id 可以为 NULL）
ALTER TABLE comments 
ADD CONSTRAINT fk_comments_article 
FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE;

-- 7. 添加 work_id 外键约束
ALTER TABLE comments 
ADD CONSTRAINT fk_comments_work 
FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE;

-- 8. 为 works 表添加 author_id 字段（如果不存在）
ALTER TABLE works 
ADD COLUMN IF NOT EXISTS author_id BIGINT UNSIGNED NOT NULL DEFAULT 1 AFTER tech_stack;

-- 9. 添加索引
CREATE INDEX IF NOT EXISTS idx_author_id ON works(author_id);

-- 10. 添加外键约束
ALTER TABLE works 
ADD CONSTRAINT fk_works_author 
FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE;

-- 11. 为 works 表添加 comment_count 字段（如果不存在）
ALTER TABLE works 
ADD COLUMN IF NOT EXISTS comment_count INT NOT NULL DEFAULT 0 AFTER view_count;

-- 12. 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 13. 同步现有数据
UPDATE works w 
SET comment_count = (
    SELECT COUNT(*) FROM comments c 
    WHERE c.work_id = w.id AND c.deleted_at IS NULL
) WHERE w.id > 0;

-- 14. 验证结果
SELECT 'Comments 表结构:' AS info;
SHOW COLUMNS FROM comments LIKE '%_id';

SELECT 'Works 表结构:' AS info;
SHOW COLUMNS FROM works LIKE '%_id';
SHOW COLUMNS FROM works LIKE 'comment_count';

SELECT '✅ 数据库修复完成' AS result;

