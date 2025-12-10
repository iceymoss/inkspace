-- ============================================
-- My Site 数据库初始化脚本
-- ============================================

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 1. 创建默认管理员账号
-- ============================================
-- 密码: admin123 (bcrypt加密)
INSERT INTO `users` (
    `username`, `password`, `email`, `nickname`, `role`, `status`, 
    `article_count`, `comment_count`, `created_at`, `updated_at`
)
VALUES (
    'admin', 
    '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbdxKnB7dAt6TvJZW', 
    'admin@example.com', 
    '管理员', 
    'admin', 
    1, 
    0, 
    0, 
    NOW(), 
    NOW()
) ON DUPLICATE KEY UPDATE updated_at = NOW();

-- ============================================
-- 2. 创建默认分类
-- ============================================
INSERT INTO `categories` (`name`, `slug`, `description`, `sort`, `article_count`, `created_at`, `updated_at`)
VALUES 
('技术分享', 'tech', '技术相关的文章和教程', 10, 0, NOW(), NOW()),
('开发笔记', 'dev-notes', '开发过程中的笔记和心得', 9, 0, NOW(), NOW()),
('生活随笔', 'life', '生活感悟和随笔', 8, 0, NOW(), NOW()),
('项目经验', 'projects', '项目开发经验分享', 7, 0, NOW(), NOW()),
('学习资源', 'resources', '学习资源和教程推荐', 6, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- ============================================
-- 3. 创建默认标签
-- ============================================
INSERT INTO `tags` (`name`, `slug`, `color`, `article_count`, `created_at`, `updated_at`)
VALUES 
-- 编程语言
('Go', 'go', '#00ADD8', 0, NOW(), NOW()),
('Python', 'python', '#3776AB', 0, NOW(), NOW()),
('JavaScript', 'javascript', '#F7DF1E', 0, NOW(), NOW()),
('TypeScript', 'typescript', '#3178C6', 0, NOW(), NOW()),

-- 前端框架
('Vue', 'vue', '#4FC08D', 0, NOW(), NOW()),
('React', 'react', '#61DAFB', 0, NOW(), NOW()),

-- 后端框架
('Gin', 'gin', '#00ADD8', 0, NOW(), NOW()),
('Spring Boot', 'spring-boot', '#6DB33F', 0, NOW(), NOW()),

-- 数据库
('MySQL', 'mysql', '#4479A1', 0, NOW(), NOW()),
('PostgreSQL', 'postgresql', '#336791', 0, NOW(), NOW()),
('MongoDB', 'mongodb', '#47A248', 0, NOW(), NOW()),
('Redis', 'redis', '#DC382D', 0, NOW(), NOW()),

-- 开发工具
('Docker', 'docker', '#2496ED', 0, NOW(), NOW()),
('Kubernetes', 'k8s', '#326CE5', 0, NOW(), NOW()),
('Git', 'git', '#F05032', 0, NOW(), NOW()),

-- 其他
('微服务', 'microservices', '#409EFF', 0, NOW(), NOW()),
('算法', 'algorithm', '#67C23A', 0, NOW(), NOW()),
('设计模式', 'design-patterns', '#E6A23C', 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- ============================================
-- 4. 创建示例文章（可选）
-- ============================================
-- 获取管理员ID和分类ID
SET @admin_id = (SELECT id FROM users WHERE username = 'admin' LIMIT 1);
SET @category_id = (SELECT id FROM categories WHERE slug = 'tech' LIMIT 1);

INSERT INTO `articles` (
    `title`, `content`, `summary`, `cover`, 
    `category_id`, `author_id`, 
    `view_count`, `like_count`, `comment_count`, 
    `word_count`, `reading_time`, 
    `status`, `is_top`, `is_recommend`, `is_original`,
    `publish_at`, `created_at`, `updated_at`
)
VALUES (
    '欢迎使用 My Site 个人博客系统',
    '# 欢迎使用 My Site

这是一个基于 Go + Gin + Vue 3 开发的现代化个人博客系统。

## 主要特性

- **前后端分离**: 使用 Go + Gin 构建后端API，Vue 3 + Element Plus 构建前端
- **Markdown编辑**: 支持Markdown格式撰写文章
- **分类和标签**: 灵活的分类和标签管理
- **评论系统**: 支持文章评论和回复
- **作品展示**: 展示个人作品和项目
- **响应式设计**: 完美适配PC和移动端
- **Docker部署**: 一键容器化部署

## 技术栈

### 后端
- Go 1.21+
- Gin Web框架
- GORM ORM框架
- MySQL 8.0
- Redis缓存
- JWT认证

### 前端
- Vue 3
- Element Plus
- Pinia状态管理
- Vue Router
- Axios

## 快速开始

使用 Docker Compose 一键启动：

```bash
docker-compose up -d
```

访问地址：http://localhost

默认管理员账号：
- 用户名：admin
- 密码：admin123

## 后续计划

- [ ] 文章搜索功能优化
- [ ] 图片上传功能
- [ ] 文章导入导出
- [ ] 主题切换
- [ ] 数据统计面板

祝你使用愉快！',
    '欢迎使用 My Site 个人博客系统！这是一个功能完善、开箱即用的个人网站解决方案。',
    '',
    @category_id,
    @admin_id,
    0, 0, 0,
    350, 2,
    1, 1, 1, 1,
    NOW(),
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 为示例文章添加标签
SET @article_id = LAST_INSERT_ID();
SET @tag_go = (SELECT id FROM tags WHERE slug = 'go' LIMIT 1);
SET @tag_vue = (SELECT id FROM tags WHERE slug = 'vue' LIMIT 1);
SET @tag_docker = (SELECT id FROM tags WHERE slug = 'docker' LIMIT 1);

INSERT INTO `article_tags` (`article_id`, `tag_id`)
VALUES 
(@article_id, @tag_go),
(@article_id, @tag_vue),
(@article_id, @tag_docker)
ON DUPLICATE KEY UPDATE article_id = article_id;

-- ============================================
-- 5. 更新统计数据
-- ============================================
-- 更新管理员的文章数
UPDATE users SET article_count = (
    SELECT COUNT(*) FROM articles WHERE author_id = users.id
) WHERE username = 'admin';

-- 更新分类的文章数
UPDATE categories SET article_count = (
    SELECT COUNT(*) FROM articles WHERE category_id = categories.id
);

-- 更新标签的文章数
UPDATE tags t SET article_count = (
    SELECT COUNT(*) FROM article_tags at
    INNER JOIN articles a ON at.article_id = a.id
    WHERE at.tag_id = t.id
);

-- ============================================
-- 6. 创建系统配置
-- ============================================
INSERT INTO `settings` (`key`, `value`, `type`, `description`, `group`, `is_public`, `created_at`, `updated_at`)
VALUES 
-- 网站基础信息
('site_name', 'My Site', 'string', '网站名称', 'site', TRUE, NOW(), NOW()),
('site_description', '基于 Go + Gin + Vue 3 的现代化个人博客系统', 'string', '网站描述', 'site', TRUE, NOW(), NOW()),
('site_keywords', 'Go,Gin,Vue,博客,个人网站', 'string', '网站关键词', 'site', TRUE, NOW(), NOW()),
('site_icp', '', 'string', '备案号', 'site', TRUE, NOW(), NOW()),
('site_copyright', 'Copyright © 2024 My Site. All rights reserved.', 'string', '版权信息', 'site', TRUE, NOW(), NOW()),
('site_logo', '', 'string', '网站Logo', 'site', TRUE, NOW(), NOW()),
('site_favicon', '', 'string', '网站图标', 'site', TRUE, NOW(), NOW()),

-- 功能开关
('comment_audit', '0', 'bool', '评论是否需要审核（0=否，1=是）', 'feature', FALSE, NOW(), NOW()),
('work_audit', '0', 'bool', '作品是否需要审核（0=否，1=是）', 'feature', FALSE, NOW(), NOW()),
('register_enabled', '1', 'bool', '是否开放注册（0=否，1=是）', 'feature', TRUE, NOW(), NOW()),
('article_comment_enabled', '1', 'bool', '是否开放文章评论（0=否，1=是）', 'feature', TRUE, NOW(), NOW()),
('work_comment_enabled', '1', 'bool', '是否开放作品评论（0=否，1=是）', 'feature', TRUE, NOW(), NOW()),

-- 上传设置
('upload_max_size', '10485760', 'int', '上传文件最大大小（字节）', 'upload', FALSE, NOW(), NOW()),
('upload_allowed_types', 'image/jpeg,image/jpg,image/png,image/gif', 'string', '允许上传的文件类型', 'upload', FALSE, NOW(), NOW()),

-- SEO设置
('seo_title', 'My Site - 个人博客', 'string', 'SEO标题', 'seo', TRUE, NOW(), NOW()),
('seo_description', '分享技术、记录生活、展示作品', 'string', 'SEO描述', 'seo', TRUE, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- ============================================
-- 7. 创建示例友情链接
-- ============================================
INSERT INTO `links` (`name`, `url`, `logo`, `description`, `sort`, `status`, `created_at`, `updated_at`)
VALUES 
('GitHub', 'https://github.com', '', '全球最大的开源社区', 10, 1, NOW(), NOW()),
('Go语言官网', 'https://golang.org', '', 'Go编程语言官网', 9, 1, NOW(), NOW()),
('Vue.js', 'https://vuejs.org', '', '渐进式JavaScript框架', 8, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- ============================================
-- 8. 优化设置
-- ============================================
-- 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 显示初始化结果
SELECT '==================================' AS '';
SELECT '数据库初始化完成！' AS '状态';
SELECT '==================================' AS '';
SELECT CONCAT('管理员账号: admin') AS '提示';
SELECT CONCAT('管理员密码: admin123') AS '提示';
SELECT CONCAT('邮箱: admin@example.com') AS '提示';
SELECT '==================================' AS '';
SELECT CONCAT('用户数量: ', COUNT(*)) AS '统计' FROM users;
SELECT CONCAT('分类数量: ', COUNT(*)) AS '统计' FROM categories;
SELECT CONCAT('标签数量: ', COUNT(*)) AS '统计' FROM tags;
SELECT CONCAT('文章数量: ', COUNT(*)) AS '统计' FROM articles;
SELECT CONCAT('友链数量: ', COUNT(*)) AS '统计' FROM links;
SELECT CONCAT('配置数量: ', COUNT(*)) AS '统计' FROM settings;
SELECT '==================================' AS '';
SELECT '✅ 所有表初始化完成！' AS '提示';
SELECT '==================================' AS '';

