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
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
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
INSERT IGNORE INTO `users` (`username`, `email`, `password_hash`, `phone`, `status`, `role`, `created_at`, `updated_at`) 
VALUES 
('admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '13800138000', 1, 'admin', NOW(), NOW()),
('demo', 'demo@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '13812345678', 1, 'user', NOW(), NOW());

-- 创建索引优化查询性能
ALTER TABLE `users` ADD INDEX `idx_users_status` (`status`);
ALTER TABLE `users` ADD INDEX `idx_users_role` (`role`);
ALTER TABLE `users` ADD INDEX `idx_users_created_at` (`created_at`);

-- 设置自增起始值
ALTER TABLE `users` AUTO_INCREMENT = 1000;

-- 创建视图：活跃用户
CREATE OR REPLACE VIEW `active_users` AS
SELECT 
    `id`,
    `username`,
    `email`,
    `phone`,
    `role`,
    `last_login`,
    `created_at`,
    `updated_at`
FROM `users` 
WHERE `deleted_at` IS NULL AND `status` = 1;

-- 创建视图：用户统计
CREATE OR REPLACE VIEW `user_stats` AS
SELECT 
    COUNT(*) as total_users,
    SUM(CASE WHEN `status` = 1 THEN 1 ELSE 0 END) as active_users,
    SUM(CASE WHEN `status` = 0 THEN 1 ELSE 0 END) as inactive_users,
    SUM(CASE WHEN DATE(`created_at`) = CURDATE() THEN 1 ELSE 0 END) as today_registered,
    SUM(CASE WHEN YEAR(`created_at`) = YEAR(NOW()) AND MONTH(`created_at`) = MONTH(NOW()) THEN 1 ELSE 0 END) as this_month_registered
FROM `users` 
WHERE `deleted_at` IS NULL;

-- 创建存储过程：清理软删除的用户（可选）
DELIMITER $$
DROP PROCEDURE IF EXISTS `CleanupDeletedUsers`;
CREATE PROCEDURE `CleanupDeletedUsers`(IN days_old INT)
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE user_count INT DEFAULT 0;
    
    -- 计算要删除的用户数量
    SELECT COUNT(*) INTO user_count 
    FROM `users` 
    WHERE `deleted_at` IS NOT NULL 
    AND `deleted_at` < DATE_SUB(NOW(), INTERVAL days_old DAY);
    
    -- 物理删除超过指定天数的软删除用户
    DELETE FROM `users` 
    WHERE `deleted_at` IS NOT NULL 
    AND `deleted_at` < DATE_SUB(NOW(), INTERVAL days_old DAY);
    
    -- 返回删除的用户数量
    SELECT CONCAT('Cleaned up ', user_count, ' deleted users older than ', days_old, ' days') as result;
END$$
DELIMITER ;

-- 创建触发器：更新用户时自动更新 updated_at
DELIMITER $$
DROP TRIGGER IF EXISTS `users_updated_at`;
CREATE TRIGGER `users_updated_at`
BEFORE UPDATE ON `users`
FOR EACH ROW 
BEGIN
    SET NEW.updated_at = NOW();
END$$
DELIMITER ;

-- 设置权限（如果需要）
-- GRANT SELECT, INSERT, UPDATE, DELETE ON go_test.* TO 'demo_user'@'%';
-- FLUSH PRIVILEGES;

SET FOREIGN_KEY_CHECKS = 1;

-- 显示初始化完成信息
SELECT 'Database initialization completed successfully!' as status;
SELECT COUNT(*) as total_users FROM users;