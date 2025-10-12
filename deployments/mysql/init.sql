-- MySQL 初始化脚本
-- 创建数据库和用户

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `go_demo` 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE `go_demo`;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `mobile` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `role` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'user',
  `last_login` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认管理员用户
-- 密码: admin123 (bcrypt hash)
INSERT IGNORE INTO `users` (`username`, `email`, `password`, `mobile`, `status`, `role`, `created_at`, `updated_at`) 
VALUES 
('admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '13800138000', 1, 'admin', NOW(), NOW()),
('demo', 'demo@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '13812345678', 1, 'user', NOW(), NOW());

-- 创建索引优化查询性能
ALTER TABLE `users` ADD INDEX `idx_users_role` (`role`);
ALTER TABLE `users` ADD INDEX `idx_users_created_at` (`created_at`);

-- 设置自增起始值
ALTER TABLE `users` AUTO_INCREMENT = 1000;

SET FOREIGN_KEY_CHECKS = 1;

-- 显示初始化完成信息
SELECT 'Database initialization completed successfully!' as status;
SELECT COUNT(*) as total_users FROM users;