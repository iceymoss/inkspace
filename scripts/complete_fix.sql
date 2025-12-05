-- 完整修复脚本 - 支持作品评论功能
-- 一次性运行所有修复

USE mysite;

SELECT '========================================' AS '';
SELECT '开始修复数据库结构...' AS '';
SELECT '========================================' AS '';

-- ==========================================
-- Part 1: 修复 Comments 表
-- ==========================================
SELECT '步骤 1/3: 修复 Comments 表...' AS '';

SET FOREIGN_KEY_CHECKS = 0;

-- 删除旧外键
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_article;
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_work;

-- 修改 article_id 允许 NULL
ALTER TABLE comments MODIFY COLUMN article_id BIGINT UNSIGNED NULL DEFAULT NULL;

-- 添加 work_id 字段
SELECT COUNT(*) INTO @work_id_exists 
FROM information_schema.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'comments'
  AND COLUMN_NAME = 'work_id';

SET @add_work_id = IF(@work_id_exists = 0,
    'ALTER TABLE comments ADD COLUMN work_id BIGINT UNSIGNED NULL DEFAULT NULL AFTER article_id',
    'SELECT "work_id 已存在" AS skip');
PREPARE stmt FROM @add_work_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 添加索引
CREATE INDEX IF NOT EXISTS idx_work_id ON comments(work_id);

-- 重新添加外键
ALTER TABLE comments 
ADD CONSTRAINT fk_comments_article 
FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE;

ALTER TABLE comments 
ADD CONSTRAINT fk_comments_work 
FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE;

SELECT '✅ Comments 表修复完成' AS result;

-- ==========================================
-- Part 2: 修复 Works 表
-- ==========================================
SELECT '步骤 2/3: 修复 Works 表...' AS '';

-- 添加 author_id 字段
SELECT COUNT(*) INTO @author_id_exists 
FROM information_schema.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'works'
  AND COLUMN_NAME = 'author_id';

SET @add_author_id = IF(@author_id_exists = 0,
    'ALTER TABLE works ADD COLUMN author_id BIGINT UNSIGNED NOT NULL DEFAULT 1 AFTER tech_stack',
    'SELECT "author_id 已存在" AS skip');
PREPARE stmt FROM @add_author_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 添加 comment_count 字段
SELECT COUNT(*) INTO @comment_count_exists 
FROM information_schema.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'works'
  AND COLUMN_NAME = 'comment_count';

SET @add_comment_count = IF(@comment_count_exists = 0,
    'ALTER TABLE works ADD COLUMN comment_count INT NOT NULL DEFAULT 0 AFTER view_count',
    'SELECT "comment_count 已存在" AS skip');
PREPARE stmt FROM @add_comment_count;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 获取有效用户ID
SET @valid_user = (SELECT id FROM users ORDER BY id LIMIT 1);

-- 更新无效的 author_id
UPDATE works 
SET author_id = @valid_user 
WHERE author_id = 0 
   OR author_id IS NULL 
   OR author_id NOT IN (SELECT id FROM users);

-- 添加索引
CREATE INDEX IF NOT EXISTS idx_author_id ON works(author_id);

-- 删除旧外键（如果存在）
ALTER TABLE works DROP FOREIGN KEY IF EXISTS fk_works_author;

-- 添加外键（注意：不使用 ON DELETE CASCADE 避免误删）
ALTER TABLE works 
ADD CONSTRAINT fk_works_author 
FOREIGN KEY (author_id) REFERENCES users(id);

SELECT '✅ Works 表修复完成' AS result;

-- ==========================================
-- Part 3: 同步计数器
-- ==========================================
SELECT '步骤 3/3: 同步计数器...' AS '';

-- 同步作品评论数
UPDATE works w 
SET comment_count = (
    SELECT COUNT(*) FROM comments c 
    WHERE c.work_id = w.id AND c.deleted_at IS NULL
);

SELECT '✅ 计数器同步完成' AS result;

SET FOREIGN_KEY_CHECKS = 1;

-- ==========================================
-- 验证最终结果
-- ==========================================
SELECT '========================================' AS '';
SELECT '✅✅✅ 所有修复完成！✅✅✅' AS '';
SELECT '========================================' AS '';

SELECT '最终验证:' AS '';

SELECT CONCAT('Comments 总数: ', COUNT(*)) AS stat FROM comments;
SELECT CONCAT('Works 总数: ', COUNT(*)) AS stat FROM works;
SELECT CONCAT('Users 总数: ', COUNT(*)) AS stat FROM users;

SELECT '现在可以重启服务了！' AS next_step;
SELECT 'make dev && make dev-admin && make dev-scheduler' AS command;

