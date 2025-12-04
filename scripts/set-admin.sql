-- 将指定用户设为管理员

-- 方式1：通过用户名设置
UPDATE users SET role = 'admin' WHERE username = 'iceymoss';

-- 方式2：通过ID设置
-- UPDATE users SET role = 'admin' WHERE id = 1;

-- 查看所有管理员
SELECT id, username, email, nickname, role FROM users WHERE role = 'admin';

