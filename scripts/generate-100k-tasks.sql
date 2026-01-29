-- 生成10万条任务数据
-- 使用UTF-8编码
SET NAMES utf8mb4;
USE task_platform;

-- 禁用外键检查和索引以加速插入
SET FOREIGN_KEY_CHECKS = 0;
SET UNIQUE_CHECKS = 0;
SET AUTOCOMMIT = 0;

-- 删除旧任务数据（保留用户和分类）
-- 只删除tasks表数据
TRUNCATE TABLE tasks;

-- 重置自增ID
ALTER TABLE tasks AUTO_INCREMENT = 1;

-- 创建存储过程生成任务数据
DELIMITER $$

DROP PROCEDURE IF EXISTS generate_tasks$$

CREATE PROCEDURE generate_tasks()
BEGIN
    DECLARE i INT DEFAULT 1;
    DECLARE batch_size INT DEFAULT 1000;
    DECLARE title_prefix VARCHAR(100);
    DECLARE content_text TEXT;
    DECLARE category_id_val INT;
    DECLARE publisher_id_val INT;
    DECLARE status_val INT;
    DECLARE amount_val DECIMAL(10,2);
    DECLARE deadline_days INT;
    DECLARE view_count_val INT;
    DECLARE apply_count_val INT;
    
    -- 任务标题模板
    DECLARE titles_array TEXT DEFAULT '开发企业官网|设计公司LOGO|编写技术文档|制作宣传视频|翻译英文资料|数据录入整理|撰写营销文案|APP开发项目|小程序开发|网站UI设计|平面海报设计|产品摄影拍摄|短视频剪辑|音频后期制作|PPT设计制作|Excel数据分析|问卷调查收集|市场调研报告|SEO优化服务|社交媒体运营|客服外包服务|电商店铺装修|产品详情页设计|品牌VI设计|包装设计制作|插画绘制|漫画创作|3D建模渲染|动画制作|游戏美术设计|程序Bug修复|代码重构优化|数据库设计|API接口开发|微信公众号开发|H5页面制作|响应式网站|后台管理系统|移动端适配|跨平台应用开发|云服务器配置|网络安全测试|软件测试服务|自动化脚本开发|爬虫程序开发|数据可视化|机器学习模型|AI算法开发|区块链开发|智能合约编写';
    
    -- 任务描述模板
    DECLARE contents_array TEXT DEFAULT '需要开发一个功能完善的系统，要求界面美观，用户体验良好，代码规范|寻找有经验的专业人士，项目周期灵活，价格可以协商，质量第一|急需完成此项目，时间紧迫，希望尽快交付，可以加价|长期合作项目，有稳定的工作量，适合兼职人员|简单任务，适合新手练手，完成质量好可以长期合作|要求有相关经验，能够独立完成，提供作品集优先|预算充足，追求高品质，欢迎资深人士投标|小型项目，工作量不大，适合快速完成|需要创意和专业技能，欢迎提供方案和报价|远程工作，时间自由，按质量付费';
    
    WHILE i <= 100000 DO
        -- 随机选择标题
        SET title_prefix = SUBSTRING_INDEX(SUBSTRING_INDEX(titles_array, '|', FLOOR(1 + RAND() * 50)), '|', -1);
        
        -- 随机选择描述
        SET content_text = CONCAT(
            SUBSTRING_INDEX(SUBSTRING_INDEX(contents_array, '|', FLOOR(1 + RAND() * 10)), '|', -1),
            '。要求：',
            CASE FLOOR(1 + RAND() * 5)
                WHEN 1 THEN '有相关经验，能独立完成，提供案例'
                WHEN 2 THEN '工作认真负责，按时交付，沟通顺畅'
                WHEN 3 THEN '技术过硬，代码规范，文档完善'
                WHEN 4 THEN '创意新颖，设计美观，符合品牌调性'
                ELSE '价格合理，质量保证，售后完善'
            END,
            '。联系方式：微信/QQ私聊详谈。'
        );
        
        -- 随机分类 (1-10)
        SET category_id_val = FLOOR(1 + RAND() * 10);
        
        -- 随机发布者 (1-5)
        SET publisher_id_val = FLOOR(1 + RAND() * 5);
        
        -- 随机状态：60%待接取，25%进行中，10%已完成，5%已取消
        SET status_val = CASE 
            WHEN RAND() < 0.60 THEN 1
            WHEN RAND() < 0.85 THEN 2
            WHEN RAND() < 0.95 THEN 4
            ELSE 5
        END;
        
        -- 随机金额 (50-10000)
        SET amount_val = ROUND(50 + RAND() * 9950, 2);
        
        -- 随机截止时间 (7-90天后)
        SET deadline_days = FLOOR(7 + RAND() * 83);
        
        -- 随机浏览量 (0-1000)
        SET view_count_val = FLOOR(RAND() * 1000);
        
        -- 随机申请人数 (0-50)
        SET apply_count_val = FLOOR(RAND() * 50);
        
        -- 插入数据
        INSERT INTO tasks (
            title, content, amount, status, view_count, apply_count,
            category_id, publisher_id, deadline, 
            create_time, update_time
        ) VALUES (
            CONCAT(title_prefix, ' #', i),
            content_text,
            amount_val,
            status_val,
            view_count_val,
            apply_count_val,
            category_id_val,
            publisher_id_val,
            DATE_ADD(NOW(), INTERVAL deadline_days DAY),
            DATE_SUB(NOW(), INTERVAL FLOOR(RAND() * 30) DAY),
            NOW()
        );
        
        -- 每1000条提交一次
        IF i % batch_size = 0 THEN
            COMMIT;
            SELECT CONCAT('已生成 ', i, ' 条任务数据...') AS progress;
        END IF;
        
        SET i = i + 1;
    END WHILE;
    
    COMMIT;
    
    -- 恢复设置
    SET FOREIGN_KEY_CHECKS = 1;
    SET UNIQUE_CHECKS = 1;
    SET AUTOCOMMIT = 1;
    
    SELECT '✅ 成功生成 100,000 条任务数据！' AS result;
    
END$$

DELIMITER ;

-- 执行存储过程
CALL generate_tasks();

-- 显示统计信息
SELECT 
    '=== 数据统计 ===' AS metric,
    '' AS count
UNION ALL
SELECT 
    '总任务数:' AS metric,
    COUNT(*) AS count
FROM tasks
UNION ALL
SELECT 
    '待接取任务:' AS metric,
    COUNT(*) AS count
FROM tasks WHERE status = 1
UNION ALL
SELECT 
    '进行中任务:' AS metric,
    COUNT(*) AS count
FROM tasks WHERE status = 2
UNION ALL
SELECT 
    '已完成任务:' AS metric,
    COUNT(*) AS count
FROM tasks WHERE status = 4
UNION ALL
SELECT 
    '平均金额:' AS metric,
    CONCAT('¥', ROUND(AVG(amount), 2)) AS count
FROM tasks
UNION ALL
SELECT 
    '总金额:' AS metric,
    CONCAT('¥', FORMAT(SUM(amount), 2)) AS count
FROM tasks;

-- 清理存储过程
DROP PROCEDURE IF EXISTS generate_tasks;
