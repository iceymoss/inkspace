-- 修复 notifications 表的 from_user_id 字段，允许 NULL（用于系统通知）
USE mysite;

-- 1. 删除现有的外键约束
ALTER TABLE `notifications` DROP FOREIGN KEY IF EXISTS `fk_notifications_from_user`;

-- 2. 修改 from_user_id 字段，允许 NULL
ALTER TABLE `notifications` 
MODIFY COLUMN `from_user_id` bigint unsigned NULL COMMENT '触发通知的用户（NULL表示系统通知）';

-- 3. 将现有的 from_user_id = 0 的记录改为 NULL（系统通知）
UPDATE `notifications` 
SET `from_user_id` = NULL 
WHERE `from_user_id` = 0;

-- 4. 重新创建外键约束（允许 NULL）
ALTER TABLE `notifications`
ADD CONSTRAINT `fk_notifications_from_user` 
FOREIGN KEY (`from_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

SELECT 'Notifications table updated successfully! from_user_id now allows NULL for system notifications.' as message;

