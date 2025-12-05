-- 测试通知功能
USE mysite;

-- 1. 检查通知表是否存在
SHOW TABLES LIKE 'notifications';

-- 2. 查看通知表结构
DESC notifications;

-- 3. 查看所有通知记录
SELECT 
    n.id,
    n.type,
    n.content,
    n.user_id as '接收者ID',
    n.from_user_id as '发送者ID',
    n.article_id,
    n.work_id,
    n.is_read as '已读',
    n.created_at as '创建时间'
FROM notifications n
ORDER BY n.created_at DESC
LIMIT 20;

-- 4. 统计各类型通知数量
SELECT 
    type as '通知类型',
    COUNT(*) as '数量',
    SUM(CASE WHEN is_read = 0 THEN 1 ELSE 0 END) as '未读数量'
FROM notifications
GROUP BY type;

-- 5. 查看最近的作品评论
SELECT 
    c.id as '评论ID',
    c.work_id as '作品ID',
    c.user_id as '评论者ID',
    c.content as '评论内容',
    w.author_id as '作品作者ID',
    c.created_at as '评论时间'
FROM comments c
LEFT JOIN works w ON c.work_id = w.id
WHERE c.work_id IS NOT NULL
ORDER BY c.created_at DESC
LIMIT 10;

-- 6. 查看用户信息
SELECT id, username, nickname FROM users LIMIT 5;

