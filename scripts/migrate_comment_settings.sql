-- 迁移评论配置：将 comment_enabled 拆分为 article_comment_enabled 和 work_comment_enabled
-- 执行方式：mysql -h localhost -u root -proot mysite < scripts/migrate_comment_settings.sql

-- 1. 如果 article_comment_enabled 不存在，从 comment_enabled 创建
INSERT INTO `settings` (`key`, `value`, `type`, `description`, `group`, `is_public`, `created_at`, `updated_at`)
SELECT 
  'article_comment_enabled',
  COALESCE((SELECT `value` FROM `settings` WHERE `key` = 'comment_enabled' LIMIT 1), '1'),
  'bool',
  '是否开放文章评论（0=否，1=是）',
  'feature',
  TRUE,
  NOW(),
  NOW()
WHERE NOT EXISTS (SELECT 1 FROM `settings` WHERE `key` = 'article_comment_enabled');

-- 2. 如果 work_comment_enabled 不存在，从 comment_enabled 创建
INSERT INTO `settings` (`key`, `value`, `type`, `description`, `group`, `is_public`, `created_at`, `updated_at`)
SELECT 
  'work_comment_enabled',
  COALESCE((SELECT `value` FROM `settings` WHERE `key` = 'comment_enabled' LIMIT 1), '1'),
  'bool',
  '是否开放作品评论（0=否，1=是）',
  'feature',
  TRUE,
  NOW(),
  NOW()
WHERE NOT EXISTS (SELECT 1 FROM `settings` WHERE `key` = 'work_comment_enabled');

-- 3. 更新已存在的配置（如果comment_enabled存在，则同步其值）
UPDATE `settings` 
SET `value` = (SELECT `value` FROM (SELECT `value` FROM `settings` WHERE `key` = 'comment_enabled' LIMIT 1) AS tmp)
WHERE `key` = 'article_comment_enabled' 
  AND EXISTS (SELECT 1 FROM `settings` WHERE `key` = 'comment_enabled');

UPDATE `settings` 
SET `value` = (SELECT `value` FROM (SELECT `value` FROM `settings` WHERE `key` = 'comment_enabled' LIMIT 1) AS tmp)
WHERE `key` = 'work_comment_enabled' 
  AND EXISTS (SELECT 1 FROM `settings` WHERE `key` = 'comment_enabled');

-- 注意：不删除旧的 comment_enabled 配置，以保持向后兼容
-- 如果需要删除，可以手动执行：
-- DELETE FROM `settings` WHERE `key` = 'comment_enabled';

