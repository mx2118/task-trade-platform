package performance

import (
	"context"
	"fmt"
	"time"
	
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	
	"task-platform-api/internal/models"
)

// DatabaseOptimizer 数据库性能优化器
type DatabaseOptimizer struct {
	db *gorm.DB
}

// NewDatabaseOptimizer 创建数据库优化器
func NewDatabaseOptimizer(db *gorm.DB) *DatabaseOptimizer {
	return &DatabaseOptimizer{db: db}
}

// CreateOptimizedIndexes 创建优化索引
func (d *DatabaseOptimizer) CreateOptimizedIndexes() error {
	indexes := []string{
		// 用户表索引
		"CREATE INDEX IF NOT EXISTS idx_users_status_created ON users(status, create_time)",
		"CREATE INDEX IF NOT EXISTS idx_users_auth_status ON users(auth_type, status)",
		"CREATE INDEX IF NOT EXISTS idx_users_credit_level ON users(credit_score, level)",
		
		// 会话表索引
		"CREATE INDEX IF NOT EXISTS idx_sessions_expire ON user_sessions(expire_time)",
		"CREATE INDEX IF NOT EXISTS idx_sessions_user_expire ON user_sessions(user_id, expire_time)",
		
		// 任务表核心索引
		"CREATE INDEX IF NOT EXISTS idx_tasks_status_created ON tasks(status, create_time DESC)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_publisher_status ON tasks(publisher_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_taker_status ON tasks(taker_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_category_status ON tasks(category_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_deadline_status ON tasks(deadline, status)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_amount_desc ON tasks(amount DESC)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_status_amount ON tasks(status, amount DESC)",
		
		// 任务阶段表索引
		"CREATE INDEX IF NOT EXISTS idx_stages_task_status ON task_stages(task_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_stages_deadline ON task_stages(deadline)",
		"CREATE INDEX IF NOT EXISTS idx_stages_task_order ON task_stages(task_id, sort_order)",
		
		// 任务交付表索引
		"CREATE INDEX IF NOT EXISTS idx_deliveries_task_status ON task_deliveries(task_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_deliveries_taker_status ON task_deliveries(taker_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_deliveries_created ON task_deliveries(create_time DESC)",
		
		// 任务申请表索引
		"CREATE INDEX IF NOT EXISTS idx_applications_task_status ON task_applications(task_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_applications_applicant ON task_applications(applicant_id, create_time DESC)",
		"CREATE INDEX IF NOT EXISTS idx_applications_price ON task_applications(quoted_price DESC)",
	}

	for _, idx := range indexes {
		if err := d.db.Exec(idx).Error; err != nil {
			return fmt.Errorf("创建索引失败: %s, 错误: %w", idx, err)
		}
	}

	return nil
}

// OptimizedTaskListQuery 优化的任务列表查询
type TaskListQuery struct {
	Status    *int8
	CategoryID *uint64
	PublisherID *uint64
	TakerID   *uint64
	MinAmount *float64
	MaxAmount *float64
	Keyword   string
	Offset    int
	Limit     int
	OrderBy   string // "created_desc", "amount_desc", "deadline_asc"
}

// GetTaskListWithOptimization 优化的任务列表查询
func (d *DatabaseOptimizer) GetTaskListWithOptimization(ctx context.Context, query TaskListQuery) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	// 构建基础查询
	db := d.db.WithContext(ctx).Model(&models.Task{})

	// 条件过滤
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.CategoryID != nil {
		db = db.Where("category_id = ?", *query.CategoryID)
	}
	if query.PublisherID != nil {
		db = db.Where("publisher_id = ?", *query.PublisherID)
	}
	if query.TakerID != nil {
		db = db.Where("taker_id = ?", *query.TakerID)
	}
	if query.MinAmount != nil {
		db = db.Where("amount >= ?", *query.MinAmount)
	}
	if query.MaxAmount != nil {
		db = db.Where("amount <= ?", *query.MaxAmount)
	}
	if query.Keyword != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取任务总数失败: %w", err)
	}

	// 排序优化
	switch query.OrderBy {
	case "created_desc":
		db = db.Order("create_time DESC")
	case "amount_desc":
		db = db.Order("amount DESC")
	case "deadline_asc":
		db = db.Order("deadline ASC")
	default:
		db = db.Order("create_time DESC")
	}

	// 分页查询（使用覆盖索引优化）
	if err := db.Preload("Publisher").
		Preload("Taker").
		Offset(query.Offset).
		Limit(query.Limit).
		Find(&tasks).Error; err != nil {
		return nil, 0, fmt.Errorf("查询任务列表失败: %w", err)
	}

	return tasks, total, nil
}

// BatchUpdateTaskStatus 批量更新任务状态（优化版）
func (d *DatabaseOptimizer) BatchUpdateTaskStatus(ctx context.Context, taskIDs []uint64, status int8) error {
	if len(taskIDs) == 0 {
		return nil
	}

	// 使用单个批量更新语句，避免多次数据库往返
	result := d.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("task_id IN ?", taskIDs).
		Updates(map[string]interface{}{
			"status":     status,
			"update_time": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("批量更新任务状态失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("没有找到需要更新的任务")
	}

	return nil
}

// GetUserTaskStatistics 获取用户任务统计（优化版）
func (d *DatabaseOptimizer) GetUserTaskStatistics(ctx context.Context, userID uint64) (map[string]int64, error) {
	var stats []struct {
		Status int8  `json:"status"`
		Count  int64 `json:"count"`
	}

	// 使用单个查询获取所有状态统计
	err := d.db.WithContext(ctx).
		Model(&models.Task{}).
		Select("status, COUNT(*) as count").
		Where("publisher_id = ?", userID).
		Group("status").
		Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("获取用户任务统计失败: %w", err)
	}

	// 转换为map格式
	result := make(map[string]int64)
	for _, stat := range stats {
		var statusName string
		switch stat.Status {
		case 0:
			statusName = "draft"
		case 1:
			statusName = "available"
		case 2:
			statusName = "in_progress"
		case 3:
			statusName = "pending_accept"
		case 4:
			statusName = "completed"
		case 5:
			statusName = "cancelled"
		default:
			statusName = "unknown"
		}
		result[statusName] = stat.Count
	}

	return result, nil
}

// OptimizedUserSearch 优化的用户搜索
func (d *DatabaseOptimizer) OptimizedUserSearch(ctx context.Context, keyword string, limit int) ([]models.User, error) {
	var users []models.User

	// 使用索引优化的搜索查询
	err := d.db.WithContext(ctx).
		Model(&models.User{}).
		Where("nickname LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Where("status = ?", 1). // 只搜索正常状态用户
		Order("credit_score DESC, create_time DESC").
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %w", err)
	}

	return users, nil
}

// CleanupExpiredSessions 清理过期会话（优化版）
func (d *DatabaseOptimizer) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	result := d.db.WithContext(ctx).
		Where("expire_time < ?", time.Now()).
		Delete(&models.UserSession{})

	if result.Error != nil {
		return 0, fmt.Errorf("清理过期会话失败: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// GetPopularTasks 获取热门任务（基于浏览量和完成数）
func (d *DatabaseOptimizer) GetPopularTasks(ctx context.Context, limit int) ([]models.Task, error) {
	var tasks []models.Task

	// 复合查询：综合考虑浏览量、完成数和最近时间
	err := d.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("status IN ?", []int8{1, 2, 4}). // 可接取、进行中、已完成
		Order("view_count DESC, create_time DESC").
		Limit(limit).
		Preload("Publisher").
		Find(&tasks).Error

	if err != nil {
		return nil, fmt.Errorf("获取热门任务失败: %w", err)
	}

	return tasks, nil
}

// UpsertUserCredit 使用 UPSERT 操作更新用户信用分
func (d *DatabaseOptimizer) UpsertUserCredit(ctx context.Context, credit *models.UserCredit) error {
	// 使用 ON DUPLICATE KEY UPDATE 避免查询-更新模式
	return d.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"score", "level", "complete_rate", "accept_rate", "violate_count", "update_time"}),
		}).
		Create(credit).Error
}

// AnalyzeTableStats 分析表统计信息（优化查询计划）
func (d *DatabaseOptimizer) AnalyzeTableStats(ctx context.Context, tableName string) error {
	return d.db.WithContext(ctx).Exec(fmt.Sprintf("ANALYZE TABLE %s", tableName)).Error
}

// GetQueryExecutionPlan 获取查询执行计划
func (d *DatabaseOptimizer) GetQueryExecutionPlan(ctx context.Context, query string) ([]map[string]interface{}, error) {
	var plans []map[string]interface{}
	err := d.db.WithContext(ctx).Raw("EXPLAIN " + query).Scan(&plans).Error
	return plans, err
}