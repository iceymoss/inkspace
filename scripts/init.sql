-- Initialize database

-- Create default admin user (password: admin123)
INSERT INTO `users` (`username`, `password`, `email`, `nickname`, `role`, `status`, `created_at`, `updated_at`)
VALUES ('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbdxKnB7dAt6TvJZW', 'admin@example.com', '管理员', 'admin', 1, NOW(), NOW());

-- Create default categories
INSERT INTO `categories` (`name`, `slug`, `description`, `sort`, `created_at`, `updated_at`)
VALUES 
('技术分享', 'tech', '技术相关文章', 1, NOW(), NOW()),
('生活随笔', 'life', '生活感悟和随笔', 2, NOW(), NOW()),
('学习笔记', 'notes', '学习过程中的笔记', 3, NOW(), NOW());

-- Create default tags
INSERT INTO `tags` (`name`, `slug`, `color`, `created_at`, `updated_at`)
VALUES 
('Go', 'go', '#00ADD8', NOW(), NOW()),
('Vue', 'vue', '#4FC08D', NOW(), NOW()),
('Docker', 'docker', '#2496ED', NOW(), NOW()),
('MySQL', 'mysql', '#4479A1', NOW(), NOW()),
('Redis', 'redis', '#DC382D', NOW(), NOW());

