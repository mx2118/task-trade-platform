-- 任务平台数据库初始化脚本

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS task_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE task_platform;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    user_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    openid VARCHAR(128) UNIQUE NOT NULL COMMENT '微信/支付宝用户标识',
    unionid VARCHAR(128) DEFAULT NULL COMMENT '跨平台用户标识',
    auth_type ENUM('wechat', 'alipay') NOT NULL COMMENT '授权类型',
    nickname VARCHAR(100) DEFAULT NULL COMMENT '用户昵称',
    avatar VARCHAR(500) DEFAULT NULL COMMENT '用户头像',
    phone VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    email VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
    credit_score DECIMAL(3,1) DEFAULT 5.0 COMMENT '信用评分(0-10)',
    level INT DEFAULT 1 COMMENT '用户等级',
    status TINYINT DEFAULT 1 COMMENT '状态:0-禁用,1-正常,2-待审核',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_openid (openid),
    INDEX idx_unionid (unionid),
    INDEX idx_auth_type (auth_type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
    session_id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    token VARCHAR(512) NOT NULL COMMENT 'JWT令牌',
    expire_time TIMESTAMP NOT NULL COMMENT '过期时间',
    device_info VARCHAR(255) DEFAULT NULL COMMENT '设备信息',
    ip_address VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    user_agent VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_token (token),
    INDEX idx_expire_time (expire_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户会话表';

-- 用户信誉表
CREATE TABLE IF NOT EXISTS user_credits (
    credit_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNIQUE NOT NULL COMMENT '用户ID',
    score DECIMAL(3,1) DEFAULT 5.0 COMMENT '信用评分',
    level INT DEFAULT 1 COMMENT '信用等级',
    complete_rate DECIMAL(3,2) DEFAULT 0.00 COMMENT '任务完成率',
    accept_rate DECIMAL(3,2) DEFAULT 0.00 COMMENT '任务接取率',
    violate_count INT DEFAULT 0 COMMENT '违规次数',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='信誉表';

-- 任务分类表
CREATE TABLE IF NOT EXISTS task_categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '分类名称',
    icon VARCHAR(255) DEFAULT NULL COMMENT '图标URL',
    description TEXT DEFAULT NULL COMMENT '分类描述',
    sort_order INT DEFAULT 0 COMMENT '排序权重',
    status TINYINT DEFAULT 1 COMMENT '状态:0-禁用,1-启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_sort_order (sort_order),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务分类表';

-- 任务表
CREATE TABLE IF NOT EXISTS tasks (
    task_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    publisher_id BIGINT NOT NULL COMMENT '发布者ID',
    taker_id BIGINT DEFAULT NULL COMMENT '接取者ID',
    title VARCHAR(200) NOT NULL COMMENT '任务标题',
    content TEXT NOT NULL COMMENT '任务内容',
    amount DECIMAL(10,2) NOT NULL COMMENT '任务金额',
    service_fee_ratio DECIMAL(3,2) DEFAULT 0.06 COMMENT '服务费比例',
    deposit_ratio DECIMAL(3,2) DEFAULT 0.10 COMMENT '保证金比例',
    deadline TIMESTAMP NOT NULL COMMENT '截止时间',
    status TINYINT DEFAULT 0 COMMENT '状态:0-草稿,1-待接取,2-进行中,3-待验收,4-已完成,5-已取消',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    apply_count INT DEFAULT 0 COMMENT '申请次数',
    category_id BIGINT DEFAULT NULL COMMENT '分类ID',
    tags JSON DEFAULT NULL COMMENT '标签',
    attachments JSON DEFAULT NULL COMMENT '附件',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_publisher_id (publisher_id),
    INDEX idx_taker_id (taker_id),
    INDEX idx_status (status),
    INDEX idx_deadline (deadline),
    INDEX idx_category_id (category_id),
    INDEX idx_create_time (create_time),
    FOREIGN KEY (publisher_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (taker_id) REFERENCES users(user_id) ON DELETE SET NULL,
    FOREIGN KEY (category_id) REFERENCES task_categories(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务表';

-- 任务阶段表
CREATE TABLE IF NOT EXISTS task_stages (
    stage_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL COMMENT '任务ID',
    stage_name VARCHAR(100) NOT NULL COMMENT '阶段名称',
    amount_ratio DECIMAL(3,2) NOT NULL COMMENT '阶段金额比例',
    amount DECIMAL(10,2) DEFAULT NULL COMMENT '阶段金额',
    deadline TIMESTAMP DEFAULT NULL COMMENT '阶段截止时间',
    status TINYINT DEFAULT 0 COMMENT '状态:0-待开始,1-进行中,2-已完成',
    description TEXT DEFAULT NULL COMMENT '阶段描述',
    sort_order INT DEFAULT 0 COMMENT '排序',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_task_id (task_id),
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务阶段表';

-- 任务申请表
CREATE TABLE IF NOT EXISTS task_applications (
    application_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL COMMENT '任务ID',
    applicant_id BIGINT NOT NULL COMMENT '申请者ID',
    message TEXT DEFAULT NULL COMMENT '申请留言',
    quoted_price DECIMAL(10,2) DEFAULT NULL COMMENT '报价',
    attachments JSON DEFAULT NULL COMMENT '附件',
    status TINYINT DEFAULT 0 COMMENT '状态:0-待审核,1-已接受,2-已拒绝',
    review_note TEXT DEFAULT NULL COMMENT '审核备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_task_id (task_id),
    INDEX idx_applicant_id (applicant_id),
    INDEX idx_status (status),
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
    FOREIGN KEY (applicant_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务申请表';

-- 交付凭证表
CREATE TABLE IF NOT EXISTS task_deliveries (
    delivery_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL COMMENT '任务ID',
    taker_id BIGINT NOT NULL COMMENT '接取者ID',
    stage_id BIGINT DEFAULT NULL COMMENT '阶段ID',
    file_url VARCHAR(500) DEFAULT NULL COMMENT '文件URL',
    content TEXT DEFAULT NULL COMMENT '交付说明',
    status TINYINT DEFAULT 0 COMMENT '状态:0-待验收,1-已验收,2-需整改',
    feedback TEXT DEFAULT NULL COMMENT '验收反馈',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_task_id (task_id),
    INDEX idx_taker_id (taker_id),
    INDEX idx_stage_id (stage_id),
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
    FOREIGN KEY (taker_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (stage_id) REFERENCES task_stages(stage_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交付凭证表';

-- 钱包表
CREATE TABLE IF NOT EXISTS wallets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNIQUE NOT NULL COMMENT '用户ID',
    balance DECIMAL(10,2) DEFAULT 0.00 COMMENT '可用余额',
    frozen_balance DECIMAL(10,2) DEFAULT 0.00 COMMENT '冻结余额',
    total_income DECIMAL(10,2) DEFAULT 0.00 COMMENT '总收入',
    total_withdraw DECIMAL(10,2) DEFAULT 0.00 COMMENT '总提现',
    version INT DEFAULT 0 COMMENT '乐观锁版本号',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='钱包表';

-- 交易表
CREATE TABLE IF NOT EXISTS trades (
    trade_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    task_id BIGINT DEFAULT NULL COMMENT '关联任务ID',
    trade_type ENUM('prepay','settle','refund','penalty') NOT NULL COMMENT '交易类型',
    amount DECIMAL(10,2) NOT NULL COMMENT '交易金额',
    third_party_no VARCHAR(64) DEFAULT NULL COMMENT '第三方交易号',
    internal_no VARCHAR(64) UNIQUE NOT NULL COMMENT '内部交易号',
    status TINYINT DEFAULT 0 COMMENT '状态:0-待支付,1-已支付,2-已失败,3-已退款',
    payment_method VARCHAR(20) DEFAULT NULL COMMENT '支付方式',
    description VARCHAR(500) DEFAULT NULL COMMENT '交易描述',
    pay_time TIMESTAMP DEFAULT NULL COMMENT '支付时间',
    expire_time TIMESTAMP DEFAULT NULL COMMENT '过期时间',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_task_id (task_id),
    INDEX idx_third_party_no (third_party_no),
    INDEX idx_internal_no (internal_no),
    INDEX idx_trade_type (trade_type),
    INDEX idx_status (status),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易表';

-- 结算表
CREATE TABLE IF NOT EXISTS settlements (
    settle_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL COMMENT '任务ID',
    publisher_amount DECIMAL(10,2) NOT NULL COMMENT '发布方收入',
    taker_amount DECIMAL(10,2) NOT NULL COMMENT '接取方收入',
    platform_fee DECIMAL(10,2) NOT NULL COMMENT '平台费用',
    penalty DECIMAL(10,2) DEFAULT 0 COMMENT '违约金',
    settle_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '结算时间',
    status TINYINT DEFAULT 0 COMMENT '状态:0-待结算,1-已结算,2-结算失败',
    remark TEXT DEFAULT NULL COMMENT '结算备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_task_id (task_id),
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='结算表';

-- 退款表
CREATE TABLE IF NOT EXISTS refunds (
    refund_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    trade_id BIGINT NOT NULL COMMENT '原交易ID',
    refund_amount DECIMAL(10,2) NOT NULL COMMENT '退款金额',
    reason VARCHAR(500) DEFAULT NULL COMMENT '退款原因',
    status TINYINT DEFAULT 0 COMMENT '状态:0-处理中,1-已成功,2-已失败',
    refund_no VARCHAR(64) UNIQUE NOT NULL COMMENT '退款单号',
    refund_time TIMESTAMP DEFAULT NULL COMMENT '退款时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_trade_id (trade_id),
    INDEX idx_refund_no (refund_no),
    FOREIGN KEY (trade_id) REFERENCES trades(trade_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款表';

-- 钱包交易记录表
CREATE TABLE IF NOT EXISTS wallet_transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    trade_id BIGINT DEFAULT NULL COMMENT '关联交易ID',
    type ENUM('income','expense','freeze','unfreeze') NOT NULL COMMENT '交易类型',
    amount DECIMAL(10,2) NOT NULL COMMENT '交易金额',
    balance_before DECIMAL(10,2) NOT NULL COMMENT '交易前余额',
    balance_after DECIMAL(10,2) NOT NULL COMMENT '交易后余额',
    description VARCHAR(500) DEFAULT NULL COMMENT '交易描述',
    related_id BIGINT DEFAULT NULL COMMENT '关联业务ID',
    related_type VARCHAR(20) DEFAULT NULL COMMENT '关联业务类型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_trade_id (trade_id),
    INDEX idx_related (related_id, related_type),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (trade_id) REFERENCES trades(trade_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='钱包交易记录表';

-- 提现申请表
CREATE TABLE IF NOT EXISTS withdraw_requests (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    amount DECIMAL(10,2) NOT NULL COMMENT '提现金额',
    withdraw_method ENUM('alipay','wechat','bank') NOT NULL COMMENT '提现方式',
    account_info JSON DEFAULT NULL COMMENT '账户信息',
    status TINYINT DEFAULT 0 COMMENT '状态:0-待处理,1-处理中,2-已完成,3-已拒绝',
    request_no VARCHAR(64) UNIQUE NOT NULL COMMENT '提现申请单号',
    process_time TIMESTAMP DEFAULT NULL COMMENT '处理时间',
    reject_reason VARCHAR(500) DEFAULT NULL COMMENT '拒绝原因',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_request_no (request_no),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='提现申请表';

-- 违规表
CREATE TABLE IF NOT EXISTS violations (
    violate_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    task_id BIGINT DEFAULT NULL COMMENT '关联任务ID',
    violate_type ENUM('fraud','delay','quality','other') NOT NULL COMMENT '违规类型',
    penalty DECIMAL(10,2) DEFAULT 0 COMMENT '处罚金额',
    description TEXT DEFAULT NULL COMMENT '违规描述',
    evidence JSON DEFAULT NULL COMMENT '违规证据',
    status TINYINT DEFAULT 0 COMMENT '状态:0-待处理,1-已处理,2-已申诉',
    handle_time TIMESTAMP DEFAULT NULL COMMENT '处理时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_task_id (task_id),
    INDEX idx_violate_type (violate_type),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='违规表';

-- 申诉表
CREATE TABLE IF NOT EXISTS complaints (
    complaint_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    task_id BIGINT DEFAULT NULL COMMENT '关联任务ID',
    violation_id BIGINT DEFAULT NULL COMMENT '关联违规ID',
    type ENUM('quality','delay','payment','other') NOT NULL COMMENT '申诉类型',
    content TEXT NOT NULL COMMENT '申诉内容',
    evidence JSON DEFAULT NULL COMMENT '申诉证据',
    status TINYINT DEFAULT 0 COMMENT '状态:0-待处理,1-处理中,2-已解决,3-已驳回',
    result TEXT DEFAULT NULL COMMENT '处理结果',
    handle_time TIMESTAMP DEFAULT NULL COMMENT '处理时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_task_id (task_id),
    INDEX idx_violation_id (violation_id),
    INDEX idx_type (type),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE SET NULL,
    FOREIGN KEY (violation_id) REFERENCES violations(violate_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='申诉表';

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
    notify_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    title VARCHAR(200) NOT NULL COMMENT '通知标题',
    content TEXT NOT NULL COMMENT '通知内容',
    type ENUM('task','payment','complaint','system') NOT NULL COMMENT '通知类型',
    is_read TINYINT DEFAULT 0 COMMENT '是否已读:0-未读,1-已读',
    related_id BIGINT DEFAULT NULL COMMENT '关联业务ID',
    related_type VARCHAR(20) DEFAULT NULL COMMENT '关联业务类型',
    data JSON DEFAULT NULL COMMENT '扩展数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_is_read (is_read),
    INDEX idx_related (related_id, related_type),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知表';

-- 风控日志表
CREATE TABLE IF NOT EXISTS risk_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT DEFAULT NULL COMMENT '用户ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型',
    risk_level TINYINT DEFAULT 0 COMMENT '风险等级:0-低,1-中,2-高',
    description TEXT DEFAULT NULL COMMENT '风险描述',
    ip_address VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    device_info JSON DEFAULT NULL COMMENT '设备信息',
    user_agent VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
    status TINYINT DEFAULT 0 COMMENT '处理状态:0-待处理,1-已通过,2-已拒绝',
    handle_note TEXT DEFAULT NULL COMMENT '处理备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_risk_level (risk_level),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='风控日志表';

-- 设备指纹表
CREATE TABLE IF NOT EXISTS device_fingerprints (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    fingerprint VARCHAR(255) UNIQUE NOT NULL COMMENT '设备指纹',
    ip_address VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    user_agent VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
    screen_info VARCHAR(200) DEFAULT NULL COMMENT '屏幕信息',
    platform VARCHAR(20) DEFAULT NULL COMMENT '平台类型',
    first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '首次发现时间',
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后发现时间',
    visit_count INT DEFAULT 1 COMMENT '访问次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_fingerprint (fingerprint),
    INDEX idx_ip_address (ip_address),
    INDEX idx_last_seen (last_seen)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备指纹表';

-- 插入初始数据
INSERT INTO task_categories (name, description, sort_order) VALUES
('技术开发', '软件开发、网站建设、APP开发等', 1),
('设计创意', 'UI设计、平面设计、logo设计等', 2),
('文案写作', '文案策划、内容创作、翻译等', 3),
('营销推广', '网络营销、SEO优化、社媒运营等', 4),
('生活服务', '家政服务、跑腿代办、维修等', 5),
('教育培训', '在线辅导、技能培训、咨询服务等', 6);

SET FOREIGN_KEY_CHECKS = 1;