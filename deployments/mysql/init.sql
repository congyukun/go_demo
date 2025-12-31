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
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 1:正常 0:禁用',
  `role` tinyint(1) NOT NULL DEFAULT '1' COMMENT '角色 1:用户 2:管理员',
  `name` char(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '名字',
  `is_activated` tinyint(1) unsigned zerofill NOT NULL DEFAULT '1' COMMENT '1:正常 2：封禁',
  `last_login` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
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