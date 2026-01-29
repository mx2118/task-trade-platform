package models

import (
    "time"
    "gorm.io/gorm"
)

// Violation 违规表
type Violation struct {
    ID        uint64    `json:"id" gorm:"primaryKey;column:violate_id"`
    UserID    uint64    `json:"user_id" gorm:"index;not null;comment:用户ID"`
    TaskID    *uint64   `json:"task_id" gorm:"index;comment:关联任务ID"`
    ViolateType string   `json:"violate_type" gorm:"type:enum('fraud','delay','quality','other');not null;comment:违规类型"`
    Penalty   float64   `json:"penalty" gorm:"type:decimal(10,2);default:0;comment:处罚金额"`
    Description string   `json:"description" gorm:"type:text;comment:违规描述"`
    Evidence   string    `json:"evidence" gorm:"type:json;comment:违规证据"`
    Status     int8      `json:"status" gorm:"default:0;comment:状态:0-待处理,1-已处理,2-已申诉"`
    HandleTime *time.Time `json:"handle_time" gorm:"comment:处理时间"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    
    User User  `json:"user" gorm:"foreignKey:UserID"`
    Task *Task `json:"task" gorm:"foreignKey:TaskID"`
}

// TableName 设置表名
func (Violation) TableName() string {
    return "violations"
}

// Complaint 申诉表
type Complaint struct {
    ID         uint64    `json:"id" gorm:"primaryKey;column:complaint_id"`
    UserID     uint64    `json:"user_id" gorm:"index;not null;comment:用户ID"`
    TaskID     *uint64   `json:"task_id" gorm:"index;comment:关联任务ID"`
    ViolationID *uint64  `json:"violation_id" gorm:"index;comment:关联违规ID"`
    Type       string    `json:"type" gorm:"type:enum('quality','delay','payment','other');not null;comment:申诉类型"`
    Content    string    `json:"content" gorm:"type:text;not null;comment:申诉内容"`
    Evidence   string    `json:"evidence" gorm:"type:json;comment:申诉证据"`
    Status     int8      `json:"status" gorm:"default:0;comment:状态:0-待处理,1-处理中,2-已解决,3-已驳回"`
    Result     string    `json:"result" gorm:"type:text;comment:处理结果"`
    HandleTime *time.Time `json:"handle_time" gorm:"comment:处理时间"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    
    User      User      `json:"user" gorm:"foreignKey:UserID"`
    Task      *Task     `json:"task" gorm:"foreignKey:TaskID"`
    Violation *Violation `json:"violation" gorm:"foreignKey:ViolationID"`
}

// TableName 设置表名
func (Complaint) TableName() string {
    return "complaints"
}

// Notification 通知表
type Notification struct {
    ID        uint64    `json:"id" gorm:"primaryKey;column:notify_id"`
    UserID    uint64    `json:"user_id" gorm:"index;not null;comment:用户ID"`
    Title     string    `json:"title" gorm:"size:200;not null;comment:通知标题"`
    Content   string    `json:"content" gorm:"type:text;not null;comment:通知内容"`
    Type      string    `json:"type" gorm:"type:enum('task','payment','complaint','system');not null;comment:通知类型"`
    IsRead    int8      `json:"is_read" gorm:"default:0;index;comment:是否已读:0-未读,1-已读"`
    RelatedID *uint64   `json:"related_id" gorm:"index;comment:关联业务ID"`
    RelatedType string   `json:"related_type" gorm:"size:20;comment:关联业务类型"`
    Data      string    `json:"data" gorm:"type:json;comment:扩展数据"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    User User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 设置表名
func (Notification) TableName() string {
    return "notifications"
}

// RiskLog 风控日志表
type RiskLog struct {
    ID          uint64    `json:"id" gorm:"primaryKey"`
    UserID      uint64    `json:"user_id" gorm:"index;comment:用户ID"`
    Action      string    `json:"action" gorm:"size:50;not null;comment:操作类型"`
    RiskLevel   int8      `json:"risk_level" gorm:"default:0;comment:风险等级:0-低,1-中,2-高"`
    Description string    `json:"description" gorm:"type:text;comment:风险描述"`
    IPAddress   string    `json:"ip_address" gorm:"size:45;comment:IP地址"`
    DeviceInfo  string    `json:"device_info" gorm:"type:json;comment:设备信息"`
    UserAgent   string    `json:"user_agent" gorm:"size:500;comment:用户代理"`
    Status      int8      `json:"status" gorm:"default:0;comment:处理状态:0-待处理,1-已通过,2-已拒绝"`
    HandleNote  string    `json:"handle_note" gorm:"type:text;comment:处理备注"`
    CreatedAt   time.Time `json:"created_at"`
    
    User User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 设置表名
func (RiskLog) TableName() string {
    return "risk_logs"
}

// DeviceFingerprint 设备指纹表
type DeviceFingerprint struct {
    ID         uint64    `json:"id" gorm:"primaryKey"`
    Fingerprint string   `json:"fingerprint" gorm:"uniqueIndex;size:255;not null;comment:设备指纹"`
    IPAddress  string    `json:"ip_address" gorm:"size:45;index;comment:IP地址"`
    UserAgent  string    `json:"user_agent" gorm:"size:500;comment:用户代理"`
    ScreenInfo string    `json:"screen_info" gorm:"size:200;comment:屏幕信息"`
    Platform   string    `json:"platform" gorm:"size:20;comment:平台类型"`
    FirstSeen  time.Time `json:"first_seen" gorm:"comment:首次发现时间"`
    LastSeen   time.Time `json:"last_seen" gorm:"comment:最后发现时间"`
    VisitCount int       `json:"visit_count" gorm:"default:1;comment:访问次数"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 设置表名
func (DeviceFingerprint) TableName() string {
    return "device_fingerprints"
}

// BeforeCreate GORM钩子：创建前
func (v *Violation) BeforeCreate(tx *gorm.DB) error {
    if v.CreatedAt.IsZero() {
        v.CreatedAt = time.Now()
    }
    if v.UpdatedAt.IsZero() {
        v.UpdatedAt = time.Now()
    }
    return nil
}

// BeforeUpdate GORM钩子：更新前
func (v *Violation) BeforeUpdate(tx *gorm.DB) error {
    v.UpdatedAt = time.Now()
    return nil
}

// BeforeCreate GORM钩子：创建前
func (c *Complaint) BeforeCreate(tx *gorm.DB) error {
    if c.CreatedAt.IsZero() {
        c.CreatedAt = time.Now()
    }
    if c.UpdatedAt.IsZero() {
        c.UpdatedAt = time.Now()
    }
    return nil
}

// BeforeUpdate GORM钩子：更新前
func (c *Complaint) BeforeUpdate(tx *gorm.DB) error {
    c.UpdatedAt = time.Now()
    return nil
}

// IsPending 申诉是否待处理
func (c *Complaint) IsPending() bool {
    return c.Status == 0
}

// IsProcessing 申诉是否处理中
func (c *Complaint) IsProcessing() bool {
    return c.Status == 1
}

// IsResolved 申诉是否已解决
func (c *Complaint) IsResolved() bool {
    return c.Status == 2
}

// IsRejected 申诉是否已驳回
func (c *Complaint) IsRejected() bool {
    return c.Status == 3
}

// IsUnread 通知是否未读
func (n *Notification) IsUnread() bool {
    return n.IsRead == 0
}

// IsReadStatus 通知是否已读
func (n *Notification) IsReadStatus() bool {
    return n.IsRead == 1
}

// MarkAsRead 标记为已读
func (n *Notification) MarkAsRead() {
    readValue := int8(1)
    n.IsRead = readValue
    n.UpdatedAt = time.Now()
}

// MarkAsUnread 标记为未读
func (n *Notification) MarkAsUnread() {
    readValue := int8(0)
    n.IsRead = readValue
    n.UpdatedAt = time.Now()
}

// IsLowRisk 风险等级是否为低
func (r *RiskLog) IsLowRisk() bool {
    return r.RiskLevel == 0
}

// IsMediumRisk 风险等级是否为中
func (r *RiskLog) IsMediumRisk() bool {
    return r.RiskLevel == 1
}

// IsHighRisk 风险等级是否为高
func (r *RiskLog) IsHighRisk() bool {
    return r.RiskLevel == 2
}