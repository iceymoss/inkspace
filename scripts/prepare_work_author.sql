-- 准备作品表的作者字段
-- 在添加外键约束之前运行

USE mysite;

-- 1. 添加 author_id 字段（如果不存在）
ALTER TABLE works 
ADD COLUMN IF NOT EXISTS author_id BIGINT UNSIGNED NOT NULL DEFAULT 1 AFTER tech_stack;

-- 2. 获取第一个管理员用户ID
SET @admin_id = (SELECT id FROM users WHERE role = 'admin' ORDER BY id LIMIT 1);

-- 3. 如果没有管理员，使用第一个用户
SET @default_author = IFNULL(@admin_id, (SELECT id FROM users ORDER BY id LIMIT 1));

-- 4. 更新所有现有作品的 author_id 为有效用户
UPDATE works 
SET author_id = @default_author 
WHERE author_id = 0 OR author_id IS NULL OR author_id NOT IN (SELECT id FROM users);

-- 5. 验证所有 author_id 都有效
SELECT COUNT(*) as '需要修复的作品数量'
FROM works 
WHERE author_id NOT IN (SELECT id FROM users);

-- 6. 添加索引（如果不存在）
CREATE INDEX IF NOT EXISTS idx_author_id ON works(author_id);

-- 7. 添加 comment_count 字段（如果不存在）
ALTER TABLE works 
ADD COLUMN IF NOT EXISTS comment_count INT NOT NULL DEFAULT 0 AFTER view_count;

SELECT '✅ 作品表准备完成，可以安全添加外键了' AS result;
SELECT CONCAT('默认作者ID: ', @default_author) AS info;

