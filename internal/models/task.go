package models

import (
    "time"
    "gorm.io/gorm"
)

// Task 任务表
type Task struct {
    ID              uint64    `json:"id" gorm:"primaryKey;column:task_id"`
    PublisherID     uint64    `json:"publisher_id" gorm:"index;not null;comment:发布者ID"`
    TakerID         uint64    `json:"taker_id" gorm:"index;comment:接取者ID"`
    Title           string    `json:"title" gorm:"size:200;not null;comment:任务标题"`
    Content         string    `json:"content" gorm:"type:text;not null;comment:任务内容"`
    Amount          float64   `json:"amount" gorm:"type:decimal(10,2);not null;comment:任务金额"`
    ServiceFeeRatio float64   `json:"service_fee_ratio" gorm:"type:decimal(3,2);default:0.06;comment:服务费比例"`
    DepositRatio    float64   `json:"deposit_ratio" gorm:"type:decimal(3,2);default:0.10;comment:保证金比例"`
    Deadline        time.Time `json:"deadline" gorm:"not null;comment:截止时间"`
    Status          int8      `json:"status" gorm:"default:0;comment:状态:0-草稿,1-待接取,2-进行中,3-待验收,4-已完成,5-已取消"`
    ViewCount       int       `json:"view_count" gorm:"default:0;comment:浏览次数"`
    ApplyCount      int       `json:"apply_count" gorm:"default:0;comment:申请次数"`
    CategoryID      uint64    `json:"category_id" gorm:"index;comment:分类ID"`
    Tags           string    `json:"tags" gorm:"type:json;comment:标签"`
    Attachments    string    `json:"attachments" gorm:"type:json;comment:附件"`
    CreatedAt      time.Time `json:"created_at" gorm:"column:create_time"`
    UpdatedAt      time.Time `json:"updated_at" gorm:"column:update_time"`
    DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
    
    Publisher  User         `json:"publisher" gorm:"foreignKey:PublisherID"`
    Taker      *User        `json:"taker" gorm:"foreignKey:TakerID"`
    Stages     []TaskStage  `json:"stages" gorm:"foreignKey:TaskID"`
    Deliveries []TaskDelivery `json:"deliveries" gorm:"foreignKey:TaskID"`
    Trades     []Trade       `json:"trades" gorm:"foreignKey:TaskID"`
}

// TableName 设置表名
func (Task) TableName() string {
    return "tasks"
}

// TaskStage 任务阶段表
type TaskStage struct {
    ID          uint64    `json:"id" gorm:"primaryKey;column:stage_id"`
    TaskID      uint64    `json:"task_id" gorm:"index;not null;comment:任务ID"`
    StageName   string    `json:"stage_name" gorm:"size:100;not null;comment:阶段名称"`
    AmountRatio float64   `json:"amount_ratio" gorm:"type:decimal(3,2);not null;comment:阶段金额比例"`
    Amount      float64   `json:"amount" gorm:"type:decimal(10,2);comment:阶段金额"`
    Deadline    *time.Time `json:"deadline" gorm:"comment:阶段截止时间"`
    Status      int8      `json:"status" gorm:"default:0;comment:状态:0-待开始,1-进行中,2-已完成"`
    Description string    `json:"description" gorm:"type:text;comment:阶段描述"`
    SortOrder   int       `json:"sort_order" gorm:"default:0;comment:排序"`
    CreatedAt   time.Time `json:"created_at" gorm:"column:create_time"`
    
    Task Task `json:"task" gorm:"foreignKey:TaskID"`
}

// TableName 设置表名
func (TaskStage) TableName() string {
    return "task_stages"
}

// TaskDelivery 交付凭证表
type TaskDelivery struct {
    ID         uint64    `json:"id" gorm:"primaryKey;column:delivery_id"`
    TaskID     uint64    `json:"task_id" gorm:"index;not null;comment:任务ID"`
    TakerID    uint64    `json:"taker_id" gorm:"index;not null;comment:接取者ID"`
    StageID    *uint64   `json:"stage_id" gorm:"index;comment:阶段ID"`
    FileURL    string    `json:"file_url" gorm:"size:500;comment:文件URL"`
    Content    string    `json:"content" gorm:"type:text;comment:交付说明"`
    Status     int8      `json:"status" gorm:"default:0;comment:状态:0-待验收,1-已验收,2-需整改"`
    Feedback   string    `json:"feedback" gorm:"type:text;comment:验收反馈"`
    CreatedAt  time.Time `json:"created_at" gorm:"column:create_time"`
    UpdatedAt  time.Time `json:"updated_at" gorm:"column:update_time"`
    
    Task  Task        `json:"task" gorm:"foreignKey:TaskID"`
    Taker User        `json:"taker" gorm:"foreignKey:TakerID"`
    Stage *TaskStage  `json:"stage" gorm:"foreignKey:StageID"`
}

// TableName 设置表名
func (TaskDelivery) TableName() string {
    return "task_deliveries"
}

// TaskCategory 任务分类表
type TaskCategory struct {
    ID          uint64    `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"size:100;not null;comment:分类名称"`
    Icon        string    `json:"icon" gorm:"size:255;comment:图标URL"`
    Description string    `json:"description" gorm:"type:text;comment:分类描述"`
    SortOrder   int       `json:"sort_order" gorm:"default:0;comment:排序权重"`
    Status      int8      `json:"status" gorm:"default:1;comment:状态:0-禁用,1-启用"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    Tasks []Task `json:"tasks" gorm:"foreignKey:CategoryID"`
}

// TableName 设置表名
func (TaskCategory) TableName() string {
    return "task_categories"
}

// TaskApplication 任务申请表
type TaskApplication struct {
    ID           uint64    `json:"id" gorm:"primaryKey;column:application_id"`
    TaskID       uint64    `json:"task_id" gorm:"index;not null;comment:任务ID"`
    ApplicantID  uint64    `json:"applicant_id" gorm:"index;not null;comment:申请者ID"`
    Message      string    `json:"message" gorm:"type:text;comment:申请留言"`
    QuotedPrice  float64   `json:"quoted_price" gorm:"type:decimal(10,2);comment:报价"`
    Attachments  string    `json:"attachments" gorm:"type:json;comment:附件"`
    Status       int8      `json:"status" gorm:"default:0;comment:状态:0-待审核,1-已接受,2-已拒绝"`
    ReviewNote   string    `json:"review_note" gorm:"type:text;comment:审核备注"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    
    Task       Task `json:"task" gorm:"foreignKey:TaskID"`
    Applicant  User `json:"applicant" gorm:"foreignKey:ApplicantID"`
}

// TableName 设置表名
func (TaskApplication) TableName() string {
    return "task_applications"
}

// BeforeCreate GORM钩子：创建前
func (t *Task) BeforeCreate(tx *gorm.DB) error {
    if t.CreatedAt.IsZero() {
        t.CreatedAt = time.Now()
    }
    if t.UpdatedAt.IsZero() {
        t.UpdatedAt = time.Now()
    }
    return nil
}

// BeforeUpdate GORM钩子：更新前
func (t *Task) BeforeUpdate(tx *gorm.DB) error {
    t.UpdatedAt = time.Now()
    return nil
}

// IsDraft 任务是否为草稿状态
func (t *Task) IsDraft() bool {
    return t.Status == 0
}

// IsAvailable 任务是否可接取
func (t *Task) IsAvailable() bool {
    return t.Status == 1 && t.TakerID == 0
}

// IsInProgress 任务是否进行中
func (t *Task) IsInProgress() bool {
    return t.Status == 2
}

// IsPendingAccept 任务是否待验收
func (t *Task) IsPendingAccept() bool {
    return t.Status == 3
}

// IsCompleted 任务是否已完成
func (t *Task) IsCompleted() bool {
    return t.Status == 4
}

// IsCancelled 任务是否已取消
func (t *Task) IsCancelled() bool {
    return t.Status == 5
}

// HasTaker 任务是否已有人接取
func (t *Task) HasTaker() bool {
    return t.TakerID > 0
}

// IsExpired 任务是否已过期
func (t *Task) IsExpired() bool {
    return time.Now().After(t.Deadline)
}

// GetPublisherAmount 获取发布者应收金额（扣除服务费后）
func (t *Task) GetPublisherAmount() float64 {
    return t.Amount * (1 - t.ServiceFeeRatio)
}

// GetTakerAmount 获取接取者应收金额（扣除服务费后）
func (t *Task) GetTakerAmount() float64 {
    return t.Amount * (1 - t.ServiceFeeRatio)
}

// GetPlatformFee 获取平台服务费
func (t *Task) GetPlatformFee() float64 {
    return t.Amount * t.ServiceFeeRatio
}

// GetDepositAmount 获取保证金金额
func (t *Task) GetDepositAmount() float64 {
    return t.Amount * t.DepositRatio
}