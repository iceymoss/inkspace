# 数据库维护指南

## 目录

1. [日常维护](#日常维护)
2. [数据库迁移](#数据库迁移)
3. [备份和恢复](#备份和恢复)
4. [性能优化](#性能优化)
5. [故障处理](#故障处理)

## 日常维护

### 1. 健康检查

定期检查数据库健康状态：

```bash
# 使用管理工具
go run cmd/migrate/main.go -action=health

# 或直接连接MySQL
mysql -u root -p -e "SELECT VERSION();"
```

### 2. 同步计数器

定期同步统计字段（建议每周一次）：

```bash
go run cmd/migrate/main.go -action=sync
```

或使用SQL直接执行：

```sql
-- 同步用户文章数
UPDATE users u 
SET article_count = (
    SELECT COUNT(*) FROM articles a 
    WHERE a.author_id = u.id AND a.deleted_at IS NULL
);

-- 同步文章评论数
UPDATE articles a 
SET comment_count = (
    SELECT COUNT(*) FROM comments c 
    WHERE c.article_id = a.id AND c.deleted_at IS NULL
);

-- 同步分类文章数
UPDATE categories cat 
SET article_count = (
    SELECT COUNT(*) FROM articles a 
    WHERE a.category_id = cat.id AND a.deleted_at IS NULL
);

-- 同步标签文章数
UPDATE tags t 
SET article_count = (
    SELECT COUNT(*) FROM article_tags at 
    INNER JOIN articles a ON at.article_id = a.id 
    WHERE at.tag_id = t.id AND a.deleted_at IS NULL
);
```

### 3. 清理软删除数据

定期清理软删除的旧数据（超过30天）：

```sql
-- 查看软删除数据
SELECT 
    'users' as table_name, COUNT(*) as count 
FROM users WHERE deleted_at IS NOT NULL
UNION ALL
SELECT 'articles', COUNT(*) FROM articles WHERE deleted_at IS NOT NULL
UNION ALL
SELECT 'comments', COUNT(*) FROM comments WHERE deleted_at IS NOT NULL;

-- 永久删除30天前的软删除数据（谨慎操作！）
DELETE FROM users 
WHERE deleted_at IS NOT NULL 
AND deleted_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

DELETE FROM articles 
WHERE deleted_at IS NOT NULL 
AND deleted_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

DELETE FROM comments 
WHERE deleted_at IS NOT NULL 
AND deleted_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
```

## 数据库迁移

### 执行迁移

```bash
# 执行所有迁移
go run cmd/migrate/main.go -action=migrate

# 查看所有表
go run cmd/migrate/main.go -action=tables

# 查看特定表信息
go run cmd/migrate/main.go -action=info articles
```

### 重置数据库（开发环境）

```bash
# 删除所有表
go run cmd/migrate/main.go -action=drop -force

# 重新执行迁移
go run cmd/migrate/main.go -action=migrate

# 导入初始数据
mysql -u root -p mysite < scripts/init.sql
```

## 备份和恢复

### 1. 完整备份

```bash
# 备份所有数据
mysqldump -u root -p mysite > backup_$(date +%Y%m%d_%H%M%S).sql

# 仅备份表结构
mysqldump -u root -p --no-data mysite > schema_$(date +%Y%m%d_%H%M%S).sql

# 备份特定表
mysqldump -u root -p mysite articles users > backup_important_$(date +%Y%m%d_%H%M%S).sql
```

### 2. 自动备份脚本

创建 `scripts/backup.sh`：

```bash
#!/bin/bash

# 配置
DB_USER="root"
DB_PASS="your_password"
DB_NAME="mysite"
BACKUP_DIR="/var/backups/mysql"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 执行备份
mysqldump -u$DB_USER -p$DB_PASS $DB_NAME | gzip > $BACKUP_DIR/mysite_$DATE.sql.gz

# 删除30天前的备份
find $BACKUP_DIR -name "mysite_*.sql.gz" -mtime +30 -delete

echo "备份完成: $BACKUP_DIR/mysite_$DATE.sql.gz"
```

设置定时任务（crontab）：

```bash
# 每天凌晨2点备份
0 2 * * * /path/to/scripts/backup.sh
```

### 3. 恢复数据

```bash
# 从备份恢复
mysql -u root -p mysite < backup_20240101_020000.sql

# 从压缩备份恢复
gunzip < backup_20240101_020000.sql.gz | mysql -u root -p mysite
```

## 性能优化

### 1. 分析慢查询

启用慢查询日志：

```sql
-- 查看慢查询配置
SHOW VARIABLES LIKE 'slow_query%';
SHOW VARIABLES LIKE 'long_query_time';

-- 启用慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2; -- 超过2秒的查询

-- 查看慢查询日志文件位置
SHOW VARIABLES LIKE 'slow_query_log_file';
```

### 2. 分析表和索引

```sql
-- 分析表
ANALYZE TABLE articles;

-- 检查表状态
SHOW TABLE STATUS LIKE 'articles';

-- 查看索引使用情况
SHOW INDEX FROM articles;

-- 检查未使用的索引
SELECT * FROM sys.schema_unused_indexes WHERE object_schema = 'mysite';
```

### 3. 优化查询

使用 EXPLAIN 分析查询：

```sql
EXPLAIN SELECT * FROM articles 
WHERE status = 1 
ORDER BY is_top DESC, created_at DESC 
LIMIT 10;
```

### 4. 表优化

```sql
-- 优化表（整理碎片）
OPTIMIZE TABLE articles;
OPTIMIZE TABLE comments;

-- 检查表
CHECK TABLE articles;

-- 修复表（如果有问题）
REPAIR TABLE articles;
```

## 故障处理

### 1. 连接问题

```sql
-- 查看当前连接数
SHOW STATUS LIKE 'Threads_connected';
SHOW STATUS LIKE 'Max_used_connections';

-- 查看最大连接数
SHOW VARIABLES LIKE 'max_connections';

-- 增加最大连接数
SET GLOBAL max_connections = 200;

-- 查看正在执行的查询
SHOW PROCESSLIST;

-- 杀死慢查询
KILL <process_id>;
```

### 2. 锁表问题

```sql
-- 查看表锁
SHOW OPEN TABLES WHERE In_use > 0;

-- 查看InnoDB锁等待
SELECT * FROM information_schema.innodb_locks;
SELECT * FROM information_schema.innodb_lock_waits;

-- 查看事务
SELECT * FROM information_schema.innodb_trx;
```

### 3. 磁盘空间

```sql
-- 查看表大小
SELECT 
    table_name AS 'Table',
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)'
FROM information_schema.TABLES 
WHERE table_schema = 'mysite'
ORDER BY (data_length + index_length) DESC;

-- 查看数据库总大小
SELECT 
    table_schema AS 'Database',
    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.TABLES 
WHERE table_schema = 'mysite'
GROUP BY table_schema;
```

### 4. 数据不一致

```bash
# 运行同步脚本
go run cmd/migrate/main.go -action=sync
```

## 监控指标

### 关键指标

1. **连接数**: 当前连接 / 最大连接
2. **查询速度**: 慢查询数量和时间
3. **表大小**: 各表数据量增长
4. **索引效率**: 索引使用率
5. **锁等待**: 锁等待时间和次数

### 监控查询

```sql
-- 数据库状态概览
SELECT 
    'Uptime' as metric, VARIABLE_VALUE as value FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Uptime'
UNION ALL
SELECT 'Threads_connected', VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Threads_connected'
UNION ALL
SELECT 'Questions', VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Questions'
UNION ALL
SELECT 'Slow_queries', VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Slow_queries';
```

## 最佳实践

1. **定期备份**: 每日自动备份，保留30天
2. **索引优化**: 定期分析查询性能，优化索引
3. **数据清理**: 定期清理软删除数据
4. **计数器同步**: 每周同步一次统计字段
5. **监控告警**: 设置关键指标告警
6. **版本控制**: 使用迁移工具管理表结构变更
7. **测试环境**: 重要操作先在测试环境验证

## 紧急联系

遇到严重问题时：
1. 立即停止写入操作
2. 检查错误日志
3. 评估数据完整性
4. 从最近的备份恢复
5. 记录问题和解决方案

---

**维护计划**:
- 每日: 自动备份
- 每周: 同步计数器、检查慢查询
- 每月: 优化表、清理旧数据、安全检查
- 每季度: 性能评估、容量规划

