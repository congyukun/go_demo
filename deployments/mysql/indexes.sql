-- 数据库索引优化
-- 用于提升查询性能

-- 用户表索引
-- 用户名索引（唯一索引，用于登录查询）
ALTER TABLE users ADD UNIQUE INDEX idx_username (username);

-- 邮箱索引（唯一索引，用于邮箱查询和验证）
ALTER TABLE users ADD UNIQUE INDEX idx_email (email);

-- 手机号索引（用于手机号查询）
ALTER TABLE users ADD INDEX idx_mobile (mobile);

-- 状态和创建时间复合索引（用于分页查询活跃用户）
ALTER TABLE users ADD INDEX idx_status_created (status, created_at DESC);

-- 软删除索引（用于过滤已删除记录）
ALTER TABLE users ADD INDEX idx_deleted_at (deleted_at);

-- 角色索引（用于按角色筛选用户）
ALTER TABLE users ADD INDEX idx_role (role);

-- 最后登录时间索引（用于统计活跃用户）
ALTER TABLE users ADD INDEX idx_last_login (last_login_at DESC);

-- 创建时间索引（用于按时间范围查询）
ALTER TABLE users ADD INDEX idx_created_at (created_at DESC);

-- 如果有其他表，继续添加索引
-- 例如：订单表、日志表等

-- 示例：如果有订单表
-- ALTER TABLE orders ADD INDEX idx_user_id (user_id);
-- ALTER TABLE orders ADD INDEX idx_status (status);
-- ALTER TABLE orders ADD INDEX idx_created_at (created_at DESC);
-- ALTER TABLE orders ADD INDEX idx_user_status (user_id, status);

-- 示例：如果有登录日志表
-- ALTER TABLE login_logs ADD INDEX idx_user_id (user_id);
-- ALTER TABLE login_logs ADD INDEX idx_login_time (login_time DESC);
-- ALTER TABLE login_logs ADD INDEX idx_ip_address (ip_address);

-- 查看索引使用情况的查询
-- SHOW INDEX FROM users;

-- 分析表以更新索引统计信息
ANALYZE TABLE users;

-- 优化表以重建索引和回收空间
-- OPTIMIZE TABLE users;