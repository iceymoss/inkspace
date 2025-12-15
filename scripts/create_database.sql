-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS inkspace CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（如果不存在）
CREATE USER IF NOT EXISTS 'inkspace'@'%' IDENTIFIED BY 'inkspace123';

-- 授权
GRANT ALL PRIVILEGES ON inkspace.* TO 'inkspace'@'%';

-- 刷新权限
FLUSH PRIVILEGES;

