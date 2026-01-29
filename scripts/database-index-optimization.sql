-- 任务交易平台数据库索引优化脚本
-- 执行前请备份数据库

-- 1. 用户表索引优化
-- 用户表已有基础索引，添加复合索引优化查询
CREATE INDEX IF NOT EXISTS idx_users_status_created ON users(status, create_time);
CREATE INDEX IF NOT EXISTS idx_users_auth_status ON users(auth_type, status);
CREATE INDEX IF NOT EXISTS idx_users_credit_level ON users(credit_score, level);

-- 2. 用户会话表索引优化
-- 添加过期时间索引用于清理过期会话
CREATE INDEX IF NOT EXISTS idx_sessions_expire ON user_sessions(expire_time);
CREATE INDEX IF NOT EXISTS idx_sessions_user_expire ON user_sessions(user_id, expire_time);

-- 3. 用户信誉表索引优化
CREATE INDEX IF NOT EXISTS idx_credits_score_level ON user_credits(score, level);
CREATE INDEX IF NOT EXISTS idx_credits_complete_rate ON user_credits(complete_rate DESC);

-- 4. 任务表核心索引优化
-- 任务状态和创建时间复合索引（列表查询最常用）
CREATE INDEX IF NOT EXISTS idx_tasks_status_created ON tasks(status, create_time DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_publisher_status ON tasks(publisher_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_taker_status ON tasks(taker_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_category_status ON tasks(category_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_deadline_status ON tasks(deadline, status);

-- 金额相关索引（按金额排序查询）
CREATE INDEX IF NOT EXISTS idx_tasks_amount_desc ON tasks(amount DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_status_amount ON tasks(status, amount DESC);

-- 全文搜索索引（任务标题和内容）
CREATE FULLTEXT INDEX IF NOT EXISTS ft_tasks_title_content ON tasks(title, content);

-- 5. 任务阶段表索引优化
CREATE INDEX IF NOT EXISTS idx_stages_task_status ON task_stages(task_id, status);
CREATE INDEX IF NOT EXISTS idx_stages_deadline ON task_stages(deadline);
CREATE INDEX IF NOT EXISTS idx_stages_task_order ON task_stages(task_id, sort_order);

-- 6. 任务交付表索引优化
CREATE INDEX IF NOT EXISTS idx_deliveries_task_status ON task_deliveries(task_id, status);
CREATE INDEX IF NOT EXISTS idx_deliveries_taker_status ON task_deliveries(taker_id, status);
CREATE INDEX IF NOT EXISTS idx_deliveries_created ON task_deliveries(create_time DESC);

-- 7. 任务申请表索引优化
CREATE INDEX IF NOT EXISTS idx_applications_task_status ON task_applications(task_id, status);
CREATE INDEX IF NOT EXISTS idx_applications_applicant ON task_applications(applicant_id, create_time DESC);
CREATE INDEX IF NOT EXISTS idx_applications_price ON task_applications(quoted_price DESC);

-- 8. 交易表索引优化（需要在models中添加相应表结构）
-- 假设存在交易表 trades
-- CREATE INDEX IF NOT EXISTS idx_trades_task_status ON trades(task_id, status);
-- CREATE INDEX IF NOT EXISTS idx_trades_user_created ON trades(user_id, create_time DESC);
-- CREATE INDEX IF NOT EXISTS idx_trades_amount ON trades(amount DESC);

-- 9. 通知表索引优化（需要在models中添加相应表结构）
-- 假设存在通知表 notifications
-- CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read, create_time DESC);
-- CREATE INDEX IF NOT EXISTS idx_notifications_type_created ON notifications(type, create_time DESC);

-- 10. 分区表优化（针对大数据量表）
-- 对于日志表、通知表等大数据量表，考虑按时间分区
-- 示例：按月分区通知表（需要MySQL 8.0+）
-- ALTER TABLE notifications PARTITION BY RANGE (YEAR(create_time) * 100 + MONTH(create_time)) (
--     PARTITION p202401 VALUES LESS THAN (202402),
--     PARTITION p202402 VALUES LESS THAN (202403),
--     PARTITION p202403 VALUES LESS THAN (202404),
--     -- 继续添加分区...
--     PARTITION p_future VALUES LESS THAN MAXVALUE
-- );

-- 查询优化：覆盖索引设计
-- 任务列表查询覆盖索引
CREATE INDEX IF NOT EXISTS idx_tasks_list_cover ON tasks(status, create_time DESC, task_id, title, amount, publisher_id);

-- 用户统计覆盖索引
CREATE INDEX IF NOT EXISTS idx_users_stats_cover ON users(status, create_time, user_id, credit_score, level);

-- 定期清理和优化建议
-- 1. 定期清理过期会话（每天）
-- DELETE FROM user_sessions WHERE expire_time < NOW();

-- 2. 定期优化表结构（每周）
-- OPTIMIZE TABLE users, tasks, user_sessions, task_applications;

-- 3. 分析表统计信息（每天）
-- ANALYZE TABLE users, tasks, user_sessions, task_applications;

-- 性能监控查询
-- 查看索引使用情况
-- SELECT 
--     TABLE_NAME,
--     INDEX_NAME,
--     CARDINALITY,
--     SUB_PART,
--     PACKED,
--     NULLABLE,
--     INDEX_TYPE
-- FROM information_schema.STATISTICS 
-- WHERE TABLE_SCHEMA = 'task_platform'
-- ORDER BY TABLE_NAME, SEQ_IN_INDEX;

-- 查看慢查询日志状态
-- SHOW VARIABLES LIKE 'slow_query_log%';
-- SHOW VARIABLES LIKE 'long_query_time';