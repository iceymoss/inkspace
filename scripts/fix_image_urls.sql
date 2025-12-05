-- 修复数据库中硬编码的图片URL
-- 将 http://localhost:8081/uploads/... 替换为 /uploads/...

USE mysite;

-- 修复用户头像URL
UPDATE users 
SET avatar = REPLACE(avatar, 'http://localhost:8081', '')
WHERE avatar LIKE 'http://localhost:8081%';

-- 修复作品封面URL
UPDATE works 
SET cover = REPLACE(cover, 'http://localhost:8081', '')
WHERE cover LIKE 'http://localhost:8081%';

-- 修复文章封面URL
UPDATE articles 
SET cover = REPLACE(cover, 'http://localhost:8081', '')
WHERE cover LIKE 'http://localhost:8081%';

-- 修复分类Logo URL
UPDATE categories 
SET logo = REPLACE(logo, 'http://localhost:8081', '')
WHERE logo LIKE 'http://localhost:8081%';

-- 显示修复结果
SELECT 'Users' as table_name, COUNT(*) as fixed_count 
FROM users 
WHERE avatar LIKE '/uploads%'
UNION ALL
SELECT 'Works', COUNT(*) 
FROM works 
WHERE cover LIKE '/uploads%'
UNION ALL
SELECT 'Articles', COUNT(*) 
FROM articles 
WHERE cover LIKE '/uploads%'
UNION ALL
SELECT 'Categories', COUNT(*) 
FROM categories 
WHERE logo LIKE '/uploads%';

