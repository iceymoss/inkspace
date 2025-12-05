-- 快速修复：允许 article_id 为 NULL
USE mysite;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 删除外键约束
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS fk_comments_article;

-- 修改列允许 NULL
ALTER TABLE comments MODIFY COLUMN article_id BIGINT UNSIGNED NULL DEFAULT NULL;

-- 添加 work_id 列（如果不存在）
ALTER TABLE comments ADD COLUMN IF NOT EXISTS work_id BIGINT UNSIGNED NULL DEFAULT NULL AFTER article_id;

-- 重新添加外键
ALTER TABLE comments ADD CONSTRAINT fk_comments_article FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE;
ALTER TABLE comments ADD CONSTRAINT fk_comments_work FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 验证
DESCRIBE comments;

SELECT '✅ 修复完成！article_id 现在允许 NULL 了' AS result;

