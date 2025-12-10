-- 修复通知表和作品表的字段问题
USE mysite;

-- 1. 检查并修复 notifications 表
-- 删除不需要的 title 字段（如果存在）
SET @exist := (SELECT COUNT(*) FROM information_schema.columns 
               WHERE table_schema = 'mysite' 
               AND table_name = 'notifications' 
               AND column_name = 'title');

SET @sqlstmt := IF(@exist > 0, 
                   'ALTER TABLE notifications DROP COLUMN title', 
                   'SELECT ''title column does not exist'' as message');

PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 2. 为 works 表添加 like_count 和 favorite_count（如果不存在）
SET @exist := (SELECT COUNT(*) FROM information_schema.columns 
               WHERE table_schema = 'mysite' 
               AND table_name = 'works' 
               AND column_name = 'like_count');

SET @sqlstmt := IF(@exist = 0, 
                   'ALTER TABLE works ADD COLUMN like_count int DEFAULT 0 AFTER comment_count', 
                   'SELECT ''like_count already exists'' as message');

PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @exist := (SELECT COUNT(*) FROM information_schema.columns 
               WHERE table_schema = 'mysite' 
               AND table_name = 'works' 
               AND column_name = 'favorite_count');

SET @sqlstmt := IF(@exist = 0, 
                   'ALTER TABLE works ADD COLUMN favorite_count int DEFAULT 0 AFTER like_count', 
                   'SELECT ''favorite_count already exists'' as message');

PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 3. 显示修复结果
SELECT '✅ 修复完成！' as message;

-- 4. 验证表结构
SELECT '--- notifications 表结构 ---' as '';
DESC notifications;

SELECT '--- works 表计数字段 ---' as '';
SHOW COLUMNS FROM works LIKE '%count%';

-- 5. 测试插入通知
INSERT INTO notifications (user_id, from_user_id, type, content, work_id, is_read, created_at, updated_at)
VALUES (1, 2, 'test', '测试通知', 2, 0, NOW(), NOW());

SELECT '--- 测试通知已插入 ---' as '';
SELECT * FROM notifications WHERE type = 'test' ORDER BY id DESC LIMIT 1;

-- 6. 清理测试数据
DELETE FROM notifications WHERE type = 'test';

SELECT '✅ 所有修复已完成，可以重新测试通知功能！' as message;

