-- 强制修复 comments 表结构
USE mysite;

-- 1. 完全禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;
SET SQL_SAFE_UPDATES = 0;

-- 2. 显示当前状态
SELECT '=== 当前 comments 表结构 ===' AS '';
SHOW COLUMNS FROM comments WHERE Field IN ('article_id', 'work_id');

-- 3. 删除所有相关外键约束
SELECT CONCAT('删除外键约束...') AS '';

-- 查找并删除所有 comments 表的外键
SELECT CONCAT('ALTER TABLE comments DROP FOREIGN KEY ', CONSTRAINT_NAME, ';') AS sql_statements
FROM information_schema.TABLE_CONSTRAINTS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'comments'
  AND CONSTRAINT_TYPE = 'FOREIGN KEY';

-- 实际删除
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_article;
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_work;
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_user;
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_parent;

SELECT '✅ 外键约束已删除' AS '';

-- 4. 强制修改 article_id 列
SELECT '修改 article_id 列...' AS '';
ALTER TABLE comments MODIFY COLUMN article_id BIGINT UNSIGNED NULL DEFAULT NULL;
SELECT '✅ article_id 修改完成' AS '';

-- 5. 检查 work_id 是否存在，不存在则添加
SELECT COUNT(*) INTO @work_id_exists
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'comments'
  AND COLUMN_NAME = 'work_id';

SELECT IF(@work_id_exists = 0, '添加 work_id 列...', 'work_id 已存在') AS '';

-- 如果不存在则添加
SET @add_work_id = IF(@work_id_exists = 0,
    'ALTER TABLE comments ADD COLUMN work_id BIGINT UNSIGNED NULL DEFAULT NULL AFTER article_id',
    'SELECT 1');
PREPARE stmt FROM @add_work_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT '✅ work_id 处理完成' AS '';

-- 6. 重新添加外键约束
SELECT '重新添加外键约束...' AS '';

ALTER TABLE comments 
ADD CONSTRAINT fk_comments_article 
FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE;

ALTER TABLE comments 
ADD CONSTRAINT fk_comments_work 
FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE;

ALTER TABLE comments 
ADD CONSTRAINT fk_comments_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE comments 
ADD CONSTRAINT fk_comments_parent 
FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE;

SELECT '✅ 外键约束已重新添加' AS '';

-- 7. 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;
SET SQL_SAFE_UPDATES = 1;

-- 8. 显示最终结果
SELECT '========================================' AS '';
SELECT '✅✅✅ 修复完成 ✅✅✅' AS '';
SELECT '========================================' AS '';

SELECT '最终 comments 表结构:' AS '';
SHOW COLUMNS FROM comments WHERE Field IN ('article_id', 'work_id');

-- 验证约束
SELECT '外键约束:' AS '';
SELECT 
    CONSTRAINT_NAME as 约束名,
    COLUMN_NAME as 列名,
    REFERENCED_TABLE_NAME as 引用表
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'comments'
  AND REFERENCED_TABLE_NAME IS NOT NULL;

SELECT '现在可以重启服务测试了！' AS '';

