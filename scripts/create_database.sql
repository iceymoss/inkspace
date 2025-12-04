-- ============================================
-- 创建数据库 SQL 脚本
-- 数据库类型: MySQL 8.0+
-- 字符集: utf8mb4
-- 排序规则: utf8mb4_unicode_ci
-- 引擎: InnoDB
-- 时区: UTC
-- ============================================

-- 设置时区为 UTC
SET time_zone = '+00:00';

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `mysite`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE `mysite`;

-- 设置会话字符集
SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 显示数据库信息
SELECT 
    SCHEMA_NAME AS '数据库名称',
    DEFAULT_CHARACTER_SET_NAME AS '字符集',
    DEFAULT_COLLATION_NAME AS '排序规则'
FROM 
    INFORMATION_SCHEMA.SCHEMATA
WHERE 
    SCHEMA_NAME = 'mysite';

-- 显示当前时区
SELECT @@global.time_zone AS '全局时区', @@session.time_zone AS '会话时区';

SELECT '==================================' AS '';
SELECT '数据库创建完成！' AS '状态';
SELECT '==================================' AS '';
SELECT '数据库名称: mysite' AS '信息';
SELECT '字符集: utf8mb4' AS '信息';
SELECT '排序规则: utf8mb4_unicode_ci' AS '信息';
SELECT '时区: UTC' AS '信息';
SELECT '==================================' AS '';

