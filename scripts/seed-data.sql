-- 任务平台测试数据脚本
USE task_platform;

-- 清空现有数据（保留表结构）
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE user_sessions;
TRUNCATE TABLE task_deliveries;
TRUNCATE TABLE task_applications;
TRUNCATE TABLE task_stages;
TRUNCATE TABLE tasks;
TRUNCATE TABLE user_credits;
TRUNCATE TABLE wallet_transactions;
TRUNCATE TABLE wallets;
TRUNCATE TABLE users;
TRUNCATE TABLE task_categories;
SET FOREIGN_KEY_CHECKS = 1;

-- 插入任务分类数据
INSERT INTO task_categories (name, icon, description, sort_order, status) VALUES
('软件开发', '💻', '网站开发、APP开发、小程序开发等软件类任务', 1, 1),
('设计美工', '🎨', 'UI设计、平面设计、视频制作等设计类任务', 2, 1),
('文案写作', '✍️', '文章撰写、文案策划、翻译等写作类任务', 3, 1),
('市场推广', '📱', '社交媒体推广、广告投放等营销类任务', 4, 1),
('数据录入', '📊', '数据整理、录入、分析等数据处理任务', 5, 1),
('其他服务', '🔧', '其他各类服务型任务', 6, 1);

-- 插入测试用户数据
INSERT INTO users (openid, unionid, auth_type, nickname, avatar, phone, email, credit_score, level, status) VALUES
('wx_test_publisher_001', 'union_test_001', 'wechat', '任务发布者小王', 'https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ1.png', '13800138001', 'wang@example.com', 8.5, 3, 1),
('wx_test_taker_001', 'union_test_002', 'wechat', '接单达人小李', 'https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ2.png', '13800138002', 'li@example.com', 9.2, 4, 1),
('alipay_test_publisher_001', 'alipay_union_001', 'alipay', '企业用户张总', 'https://tfs.alipayobjects.com/images/partner/T1.jpg', '13800138003', 'zhang@company.com', 7.8, 2, 1),
('wx_test_taker_002', 'union_test_003', 'wechat', '兼职小刘', 'https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ3.png', '13800138004', 'liu@example.com', 8.0, 3, 1),
('alipay_test_taker_001', 'alipay_union_002', 'alipay', '自由职业者陈小姐', 'https://tfs.alipayobjects.com/images/partner/T2.jpg', '13800138005', 'chen@freelance.com', 9.5, 5, 1);

-- 获取用户ID（用于后续插入）
SET @publisher1 = (SELECT user_id FROM users WHERE openid = 'wx_test_publisher_001');
SET @taker1 = (SELECT user_id FROM users WHERE openid = 'wx_test_taker_001');
SET @publisher2 = (SELECT user_id FROM users WHERE openid = 'alipay_test_publisher_001');
SET @taker2 = (SELECT user_id FROM users WHERE openid = 'wx_test_taker_002');
SET @taker3 = (SELECT user_id FROM users WHERE openid = 'alipay_test_taker_001');

-- 插入用户信誉数据
INSERT INTO user_credits (user_id, score, level, complete_rate, accept_rate, violate_count) VALUES
(@publisher1, 8.5, 3, 0.92, 0.95, 0),
(@taker1, 9.2, 4, 0.98, 0.96, 0),
(@publisher2, 7.8, 2, 0.85, 0.90, 1),
(@taker2, 8.0, 3, 0.88, 0.92, 0),
(@taker3, 9.5, 5, 0.99, 0.98, 0);

-- 插入钱包数据
INSERT INTO wallets (user_id, balance, frozen_balance, total_income, total_withdraw) VALUES
(@publisher1, 5000.00, 200.00, 10000.00, 5200.00),
(@taker1, 3200.50, 0.00, 8500.00, 5299.50),
(@publisher2, 8500.00, 500.00, 20000.00, 12000.00),
(@taker2, 1500.00, 100.00, 4500.00, 3100.00),
(@taker3, 6800.00, 0.00, 15000.00, 8200.00);

-- 插入任务数据
INSERT INTO tasks (publisher_id, taker_id, title, content, amount, service_fee_ratio, deposit_ratio, deadline, status, view_count, apply_count, category_id, tags, attachments) VALUES
-- 待接取的任务
(@publisher1, NULL, '开发企业官网', '需要开发一个响应式企业官网，包含首页、产品展示、新闻中心、关于我们等栏目。要求使用Vue3+Element Plus，页面美观大方，兼容移动端。', 3800.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL 15 DAY), 1, 156, 8, 1, '["网站开发", "Vue3", "响应式"]', '["requirement.pdf", "design.fig"]'),

(@publisher2, NULL, '设计公司LOGO', '为新创科技公司设计LOGO，要求简洁现代，体现科技感和创新精神。需要提供多个设计方案，包含彩色和黑白版本，以及不同尺寸的应用效果图。', 1500.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL 7 DAY), 1, 89, 5, 2, '["LOGO设计", "品牌设计", "VI设计"]', '["brand_brief.docx"]'),

(@publisher1, NULL, '撰写产品推广文案', '为新上市的智能手表撰写一系列推广文案，包括产品介绍、功能特点、使用场景等。文案需要简洁有力，突出产品卖点，适合在社交媒体和电商平台使用。', 800.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL 5 DAY), 1, 67, 3, 3, '["文案撰写", "产品推广", "营销文案"]', '[]'),

(@publisher2, NULL, '小程序开发', '开发一个在线预约服务的微信小程序，包含用户注册登录、服务浏览、在线预约、订单管理等功能。UI设计简洁清晰，操作流程顺畅。', 5800.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL 20 DAY), 1, 234, 12, 1, '["小程序", "微信开发", "预约系统"]', '["requirement.pdf", "prototype.rp"]'),

(@publisher1, NULL, '短视频剪辑', '将拍摄好的产品宣传素材剪辑成3-5分钟的短视频，需要添加字幕、配乐、转场效果等。要求画面流畅，节奏明快，符合抖音、快手等平台风格。', 1200.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL 10 DAY), 1, 98, 6, 2, '["视频剪辑", "短视频", "后期制作"]', '["raw_footage.zip"]'),

-- 进行中的任务
(@publisher2, @taker1, 'CRM系统开发', '开发一套客户关系管理系统，包含客户管理、销售跟进、数据统计等模块。前端使用Vue3，后端使用Go，数据库MySQL。', 12000.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL 30 DAY), 2, 178, 4, 1, '["系统开发", "CRM", "全栈开发"]', '["requirement.pdf", "database.sql"]'),

(@publisher1, @taker3, '品牌VI设计', '设计完整的品牌VI系统，包括LOGO、名片、信纸、PPT模板等应用设计。需要提供设计规范手册和所有源文件。', 4500.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL 25 DAY), 2, 145, 7, 2, '["VI设计", "品牌设计", "视觉识别"]', '["brand_strategy.pdf"]'),

-- 待验收的任务
(@publisher1, @taker2, '数据录入整理', '将纸质文档扫描件中的客户信息录入到Excel表格中，约2000条数据。要求准确无误，格式统一规范。', 600.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL -2 DAY), 3, 45, 1, 5, '["数据录入", "Excel", "文档整理"]', '["scanned_docs.zip"]'),

-- 已完成的任务
(@publisher2, @taker1, 'APP界面设计', '设计电商类APP的UI界面，包含首页、分类、购物车、个人中心等10个主要页面。设计风格现代简约，色彩搭配和谐。', 2800.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL -5 DAY), 4, 123, 2, 2, '["UI设计", "APP设计", "移动端"]', '["design.xd"]'),

(@publisher1, @taker3, 'SEO优化服务', '对企业网站进行SEO优化，包括关键词优化、内容优化、外链建设等。目标是在3个月内将主要关键词排名提升到百度首页。', 3500.00, 0.06, 0.10, DATE_ADD(NOW(), INTERVAL -10 DAY), 4, 89, 1, 4, '["SEO", "网站优化", "搜索引擎"]', '["seo_report.pdf"]');

-- 插入任务申请记录
INSERT INTO task_applications (task_id, applicant_id, message, quoted_price, status) VALUES
((SELECT task_id FROM tasks WHERE title = '开发企业官网' LIMIT 1), @taker1, '有3年Vue开发经验，做过多个企业官网项目，可以提供作品案例查看。', 3800.00, 0),
((SELECT task_id FROM tasks WHERE title = '开发企业官网' LIMIT 1), @taker3, '专业前端开发工程师，精通Vue3和响应式布局，保证按时高质量交付。', 3800.00, 0),
((SELECT task_id FROM tasks WHERE title = '设计公司LOGO' LIMIT 1), @taker3, '资深品牌设计师，擅长科技类LOGO设计，可以提供多套方案供选择。', 1500.00, 0),
((SELECT task_id FROM tasks WHERE title = '小程序开发' LIMIT 1), @taker1, '开发过多个微信小程序项目，熟悉小程序开发规范和审核流程。', 5800.00, 0);

-- 插入任务交付记录
INSERT INTO task_deliveries (task_id, taker_id, content, file_url, status) VALUES
((SELECT task_id FROM tasks WHERE title = '数据录入整理' LIMIT 1), @taker2, '已完成所有2000条数据的录入工作，数据已按要求格式整理完毕。请查收附件中的Excel文件。', 'https://example.com/files/customer_data.xlsx', 0);

-- 插入钱包交易记录
INSERT INTO wallet_transactions (user_id, amount, type, balance_before, balance_after, description, related_id, related_type) VALUES
(@publisher1, 3800.00, 'freeze', 8800.00, 5000.00, '发布任务：开发企业官网', (SELECT task_id FROM tasks WHERE title = '开发企业官网' LIMIT 1), 'task'),
(@publisher2, 1500.00, 'freeze', 10000.00, 8500.00, '发布任务：设计公司LOGO', (SELECT task_id FROM tasks WHERE title = '设计公司LOGO' LIMIT 1), 'task'),
(@taker1, 2632.00, 'income', 568.50, 3200.50, '完成任务：APP界面设计（扣除手续费168元）', (SELECT task_id FROM tasks WHERE title = 'APP界面设计' LIMIT 1), 'task'),
(@publisher2, 2800.00, 'expense', 11300.00, 8500.00, '支付任务：APP界面设计', (SELECT task_id FROM tasks WHERE title = 'APP界面设计' LIMIT 1), 'task'),
(@taker3, 3290.00, 'income', 3510.00, 6800.00, '完成任务：SEO优化服务（扣除手续费210元）', (SELECT task_id FROM tasks WHERE title = 'SEO优化服务' LIMIT 1), 'task'),
(@publisher1, 3500.00, 'expense', 8500.00, 5000.00, '支付任务：SEO优化服务', (SELECT task_id FROM tasks WHERE title = 'SEO优化服务' LIMIT 1), 'task');

-- 提交更改
COMMIT;

-- 验证数据
SELECT '=== 数据统计 ===' as '';
SELECT '用户总数:' as metric, COUNT(*) as count FROM users;
SELECT '任务总数:' as metric, COUNT(*) as count FROM tasks;
SELECT '待接取任务:' as metric, COUNT(*) as count FROM tasks WHERE status = 1;
SELECT '进行中任务:' as metric, COUNT(*) as count FROM tasks WHERE status = 2;
SELECT '已完成任务:' as metric, COUNT(*) as count FROM tasks WHERE status = 4;
SELECT '任务分类数:' as metric, COUNT(*) as count FROM task_categories;
SELECT '任务申请数:' as metric, COUNT(*) as count FROM task_applications;
