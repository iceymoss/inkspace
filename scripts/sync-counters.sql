-- ============================================
-- 同步所有计数器字段
-- ============================================

-- 1. 同步用户的文章数
UPDATE users u 
SET article_count = (
    SELECT COUNT(*) 
    FROM articles a 
    WHERE a.author_id = u.id 
    AND a.deleted_at IS NULL
);

-- 2. 同步用户的评论数
UPDATE users u 
SET comment_count = (
    SELECT COUNT(*) 
    FROM comments c 
    WHERE c.user_id = u.id 
    AND c.deleted_at IS NULL
);

-- 3. 同步文章的评论数
UPDATE articles a 
SET comment_count = (
    SELECT COUNT(*) 
    FROM comments c 
    WHERE c.article_id = a.id 
    AND c.deleted_at IS NULL
);

-- 4. 同步分类的文章数
UPDATE categories cat 
SET article_count = (
    SELECT COUNT(*) 
    FROM articles a 
    WHERE a.category_id = cat.id 
    AND a.deleted_at IS NULL
);

-- 5. 同步标签的文章数
UPDATE tags t 
SET article_count = (
    SELECT COUNT(*) 
    FROM article_tags at 
    INNER JOIN articles a ON at.article_id = a.id 
    WHERE at.tag_id = t.id 
    AND a.deleted_at IS NULL
);

-- 6. 同步评论的回复数
UPDATE comments c 
SET reply_count = (
    SELECT COUNT(*) 
    FROM comments r 
    WHERE r.parent_id = c.id 
    AND r.deleted_at IS NULL
);

-- 7. 同步用户的关注数 🆕
UPDATE users u
SET following_count = (
    SELECT COUNT(*)
    FROM user_follows uf
    WHERE uf.follower_id = u.id
    AND uf.deleted_at IS NULL
);

-- 8. 同步用户的粉丝数 🆕
UPDATE users u
SET follower_count = (
    SELECT COUNT(*)
    FROM user_follows uf
    WHERE uf.following_id = u.id
    AND uf.deleted_at IS NULL
);

-- 9. 同步用户的收藏数 🆕
UPDATE users u
SET favorite_count = (
    SELECT COUNT(*)
    FROM article_favorites af
    WHERE af.user_id = u.id
    AND af.deleted_at IS NULL
);

-- 10. 同步文章的收藏数 🆕
UPDATE articles a
SET favorite_count = (
    SELECT COUNT(*)
    FROM article_favorites af
    WHERE af.article_id = a.id
    AND af.deleted_at IS NULL
);

-- 11. 同步作品的评论数 🆕
UPDATE works w 
SET comment_count = (
    SELECT COUNT(*) 
    FROM comments c 
    WHERE c.work_id = w.id 
    AND c.deleted_at IS NULL
);

-- 显示同步结果
SELECT '========================================' AS '';
SELECT '✅ 计数器同步完成！' AS '状态';
SELECT '========================================' AS '';
SELECT CONCAT('用户总数: ', COUNT(*)) AS '统计' FROM users WHERE deleted_at IS NULL;
SELECT CONCAT('文章总数: ', COUNT(*)) AS '统计' FROM articles WHERE deleted_at IS NULL;
SELECT CONCAT('作品总数: ', COUNT(*)) AS '统计' FROM works WHERE deleted_at IS NULL;
SELECT CONCAT('评论总数: ', COUNT(*)) AS '统计' FROM comments WHERE deleted_at IS NULL;
SELECT CONCAT('分类总数: ', COUNT(*)) AS '统计' FROM categories WHERE deleted_at IS NULL;
SELECT CONCAT('标签总数: ', COUNT(*)) AS '统计' FROM tags WHERE deleted_at IS NULL;
SELECT '========================================' AS '';
SELECT '统计字段已同步：' AS '提示';
SELECT '  - 用户文章数、评论数、关注数、粉丝数、收藏数' AS '提示';
SELECT '  - 文章评论数、收藏数' AS '提示';
SELECT '  - 作品评论数' AS '提示';
SELECT '  - 分类文章数、标签文章数' AS '提示';
SELECT '  - 评论回复数' AS '提示';
SELECT '========================================' AS '';

