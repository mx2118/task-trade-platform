package models

import (
    "time"
    "gorm.io/gorm"
)

// User 用户表
type User struct {
    ID         uint64    `json:"id" gorm:"primaryKey;column:user_id"`
    OpenID     string    `json:"openid" gorm:"uniqueIndex;size:128;comment:微信/支付宝用户标识"`
    UnionID    string    `json:"unionid" gorm:"index;size:128;comment:跨平台用户标识"`
    AuthType   string    `json:"auth_type" gorm:"type:enum('wechat','alipay');not null;comment:授权类型"`
    Nickname   string    `json:"nickname" gorm:"size:100;comment:用户昵称"`
    Avatar     string    `json:"avatar" gorm:"size:500;comment:用户头像"`
    Phone      string    `json:"phone" gorm:"size:20;comment:手机号"`
    Email      string    `json:"email" gorm:"size:100;comment:邮箱"`
    CreditScore float32   `json:"credit_score" gorm:"type:decimal(3,1);default:5.0;comment:信用评分(0-10)"`
    Level      int       `json:"level" gorm:"default:1;comment:用户等级"`
    Status     int8      `json:"status" gorm:"default:1;comment:状态:0-禁用,1-正常,2-待审核"`
    CreatedAt  time.Time `json:"created_at" gorm:"column:create_time"`
    UpdatedAt  time.Time `json:"updated_at" gorm:"column:update_time"`
    DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 设置表名
func (User) TableName() string {
    return "users"
}

// UserSession 用户会话表
type UserSession struct {
    SessionID  string    `json:"session_id" gorm:"primaryKey;size:64"`
    UserID     uint64    `json:"user_id" gorm:"index;not null;comment:用户ID"`
    Token      string    `json:"token" gorm:"size:512;not null;comment:JWT令牌"`
    ExpireTime time.Time `json:"expire_time" gorm:"not null;comment:过期时间"`
    DeviceInfo string    `json:"device_info" gorm:"size:255;comment:设备信息"`
    IPAddress  string    `json:"ip_address" gorm:"size:45;comment:IP地址"`
    UserAgent  string    `json:"user_agent" gorm:"size:500;comment:用户代理"`
    CreatedAt  time.Time `json:"created_at" gorm:"column:create_time"`
    
    User User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 设置表名
func (UserSession) TableName() string {
    return "user_sessions"
}

// UserCredit 用户信誉表
type UserCredit struct {
    ID           uint64    `json:"id" gorm:"primaryKey;column:credit_id"`
    UserID       uint64    `json:"user_id" gorm:"uniqueIndex;not null;comment:用户ID"`
    Score        float32   `json:"score" gorm:"type:decimal(3,1);default:5.0;comment:信用评分"`
    Level        int       `json:"level" gorm:"default:1;comment:信用等级"`
    CompleteRate float32   `json:"complete_rate" gorm:"type:decimal(3,2);default:0.00;comment:任务完成率"`
    AcceptRate   float32   `json:"accept_rate" gorm:"type:decimal(3,2);default:0.00;comment:任务接取率"`
    ViolateCount int       `json:"violate_count" gorm:"default:0;comment:违规次数"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"column:update_time"`
    
    User User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 设置表名
func (UserCredit) TableName() string {
    return "user_credits"
}

// BeforeCreate GORM钩子：创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.CreatedAt.IsZero() {
        u.CreatedAt = time.Now()
    }
    if u.UpdatedAt.IsZero() {
        u.UpdatedAt = time.Now()
    }
    return nil
}

// BeforeUpdate GORM钩子：更新前
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    u.UpdatedAt = time.Now()
    return nil
}

// IsNormal 用户是否正常状态
func (u *User) IsNormal() bool {
    return u.Status == 1
}

// IsDisabled 用户是否被禁用
func (u *User) IsDisabled() bool {
    return u.Status == 0
}

// IsPending 用户是否待审核
func (u *User) IsPending() bool {
    return u.Status == 2
}