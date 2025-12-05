-- 添加作品类型支持
USE mysite;

-- 1. 添加 type 字段
ALTER TABLE works 
ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'project' AFTER title;

-- 2. 添加 metadata 字段（存储类型专属数据）
ALTER TABLE works 
ADD COLUMN IF NOT EXISTS metadata TEXT NULL AFTER type;

-- 3. 添加 daily_quota 字段（是否受每日配额限制）
ALTER TABLE works 
ADD COLUMN IF NOT EXISTS daily_quota BOOLEAN DEFAULT FALSE AFTER metadata;

-- 4. 添加索引
CREATE INDEX IF NOT EXISTS idx_type ON works(type);
CREATE INDEX IF NOT EXISTS idx_author_date ON works(author_id, created_at);

-- 5. 验证结果
SELECT '=== Works 表新增字段 ===' AS '';
SHOW COLUMNS FROM works WHERE Field IN ('type', 'metadata', 'daily_quota');

SELECT '✅ 作品类型支持添加完成' AS result;
SELECT '支持的类型：project (开源项目), photography (摄影作品)' AS info;

