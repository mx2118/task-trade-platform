package models

import (
    "time"
    "gorm.io/gorm"
)

// Trade 交易表
type Trade struct {
    ID            uint64    `json:"id" gorm:"primaryKey;column:trade_id"`
    UserID        uint64    `json:"user_id" gorm:"index;not null;comment:用户ID"`
    TaskID        *uint64   `json:"task_id" gorm:"index;comment:关联任务ID"`
    TradeType     string    `json:"trade_type" gorm:"type:enum('prepay','settle','refund','penalty');not null;comment:交易类型"`
    Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null;comment:交易金额"`
    ThirdPartyNo  string    `json:"third_party_no" gorm:"size:64;index;comment:第三方交易号"`
    InternalNo    string    `json:"internal_no" gorm:"size:64;uniqueIndex;comment:内部交易号"`
    Status        int8      `json:"status" gorm:"default:0;comment:状态:0-待支付,1-已支付,2-已失败,3-已退款"`
    PaymentMethod string    `json:"payment_method" gorm:"size:20;comment:支付方式"`
    Description   string    `json:"description" gorm:"size:500;comment:交易描述"`
    PayTime       *time.Time `json:"pay_time" gorm:"comment:支付时间"`
    ExpireTime    *time.Time `json:"expire_time" gorm:"comment:过期时间"`
    CreatedAt     time.Time `json:"created_at" gorm:"column:create_time"`
    UpdatedAt     time.Time `json:"updated_at" gorm:"column:update_time"`
    
    User User  `json:"user" gorm:"foreignKey:UserID"`
    Task *Task `json:"task" gorm:"foreignKey:TaskID"`
}

// TableName 设置表名
func (Trade) TableName() string {
    return "trades"
}

// Settlement 结算表
type Settlement struct {
    ID             uint64    `json:"id" gorm:"primaryKey;column:settle_id"`
    TaskID         uint64    `json:"task_id" gorm:"index;not null;comment:任务ID"`
    PublisherAmount float64  `json:"publisher_amount" gorm:"type:decimal(10,2);not null;comment:发布方收入"`
    TakerAmount    float64  `json:"taker_amount" gorm:"type:decimal(10,2);not null;comment:接取方收入"`
    PlatformFee    float64  `json:"platform_fee" gorm:"type:decimal(10,2);not null;comment:平台费用"`
    Penalty        float64  `json:"penalty" gorm:"type:decimal(10,2);default:0;comment:违约金"`
    SettleTime     time.Time `json:"settle_time" gorm:"column:settle_time"`
    Status         int8      `json:"status" gorm:"default:0;comment:状态:0-待结算,1-已结算,2-结算失败"`
    Remark         string    `json:"remark" gorm:"type:text;comment:结算备注"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    
    Task Task `json:"task" gorm:"foreignKey:TaskID"`
}

// TableName 设置表名
func (Settlement) TableName() string {
    return "settlements"
}

// Refund 退款表
type Refund struct {
    ID          uint64    `json:"id" gorm:"primaryKey;column:refund_id"`
    TradeID     uint64    `json:"trade_id" gorm:"index;not null;comment:原交易ID"`
    RefundAmount float64   `json:"refund_amount" gorm:"type:decimal(10,2);not null;comment:退款金额"`
    Reason      string    `json:"reason" gorm:"size:500;comment:退款原因"`
    Status      int8      `json:"status" gorm:"default:0;comment:状态:0-处理中,1-已成功,2-已失败"`
    RefundNo    string    `json:"refund_no" gorm:"size:64;uniqueIndex;comment:退款单号"`
    RefundTime  *time.Time `json:"refund_time" gorm:"comment:退款时间"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    Trade Trade `json:"trade" gorm:"foreignKey:TradeID"`
}

// TableName 设置表名
func (Refund) TableName() string {
    return "refunds"
}

// Wallet 钱包表
type Wallet struct {
    ID              uint64    `json:"id" gorm:"primaryKey"`
    UserID          uint64    `json:"user_id" gorm:"uniqueIndex;not null;comment:用户ID"`
    Balance         float64   `json:"balance" gorm:"type:decimal(10,2);default:0.00;comment:可用余额"`
    FrozenBalance   float64   `json:"frozen_balance" gorm:"type:decimal(10,2);default:0.00;comment:冻结余额"`
    TotalIncome     float64   `json:"total_income" gorm:"type:decimal(10,2);default:0.00;comment:总收入"`
    TotalWithdraw   float64   `json:"total_withdraw" gorm:"type:decimal(10,2);default:0.00;comment:总提现"`
    Version         int       `json:"version" gorm:"default:0;comment:乐观锁版本号"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    
    User User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 设置表名
func (Wallet) TableName() string {
    return "wallets"
}

// WalletTransaction 钱包交易记录表
type WalletTransaction struct {
    ID             uint64    `json:"id" gorm:"primaryKey"`
    UserID         uint64    `json:"user_id" gorm:"index;not null;comment:用户ID"`
    TradeID        *uint64   `json:"trade_id" gorm:"index;comment:关联交易ID"`
    Type           string    `json:"type" gorm:"type:enum('income','expense','freeze','unfreeze');not null;comment:交易类型"`
    Amount         float64   `json:"amount" gorm:"type:decimal(10,2);not null;comment:交易金额"`
    BalanceBefore  float64   `json:"balance_before" gorm:"type:decimal(10,2);not null;comment:交易前余额"`
    BalanceAfter   float64   `json:"balance_after" gorm:"type:decimal(10,2);not null;comment:交易后余额"`
    Description    string    `json:"description" gorm:"size:500;comment:交易描述"`
    RelatedID      uint64    `json:"related_id" gorm:"index;comment:关联业务ID"`
    RelatedType    string    `json:"related_type" gorm:"size:20;comment:关联业务类型"`
    CreatedAt      time.Time `json:"created_at"`
    
    User  User  `json:"user" gorm:"foreignKey:UserID"`
    Trade *Trade `json:"trade" gorm:"foreignKey:TradeID"`
}

// TableName 设置表名
func (WalletTransaction) TableName() string {
    return "wallet_transactions"
}

// WithdrawRequest 提现申请表
type WithdrawRequest struct {
    ID             uint64    `json:"id" gorm:"primaryKey"`
    UserID         uint64    `json:"user_id" gorm:"index;not null;comment:用户ID"`
    Amount         float64   `json:"amount" gorm:"type:decimal(10,2);not null;comment:提现金额"`
    WithdrawMethod string    `json:"withdraw_method" gorm:"type:enum('alipay','wechat','bank');not null;comment:提现方式"`
    AccountInfo    string    `json:"account_info" gorm:"type:json;comment:账户信息"`
    Status         int8      `json:"status" gorm:"default:0;comment:状态:0-待处理,1-处理中,2-已完成,3-已拒绝"`
    RequestNo      string    `json:"request_no" gorm:"size:64;uniqueIndex;comment:提现申请单号"`
    ProcessTime    *time.Time `json:"process_time" gorm:"comment:处理时间"`
    RejectReason   string    `json:"reject_reason" gorm:"size:500;comment:拒绝原因"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    
    User User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 设置表名
func (WithdrawRequest) TableName() string {
    return "withdraw_requests"
}

// BeforeCreate GORM钩子：创建前
func (t *Trade) BeforeCreate(tx *gorm.DB) error {
    if t.CreatedAt.IsZero() {
        t.CreatedAt = time.Now()
    }
    if t.UpdatedAt.IsZero() {
        t.UpdatedAt = time.Now()
    }
    return nil
}

// BeforeUpdate GORM钩子：更新前
func (t *Trade) BeforeUpdate(tx *gorm.DB) error {
    t.UpdatedAt = time.Now()
    return nil
}

// IsPending 交易是否待支付
func (t *Trade) IsPending() bool {
    return t.Status == 0
}

// IsPaid 交易是否已支付
func (t *Trade) IsPaid() bool {
    return t.Status == 1
}

// IsFailed 交易是否失败
func (t *Trade) IsFailed() bool {
    return t.Status == 2
}

// IsRefunded 交易是否已退款
func (t *Trade) IsRefunded() bool {
    return t.Status == 3
}

// IsExpired 交易是否已过期
func (t *Trade) IsExpired() bool {
    if t.ExpireTime == nil {
        return false
    }
    return time.Now().After(*t.ExpireTime)
}

// GetAvailableBalance 获取可用余额
func (w *Wallet) GetAvailableBalance() float64 {
    return w.Balance
}

// GetFrozenBalance 获取冻结余额
func (w *Wallet) GetFrozenBalance() float64 {
    return w.FrozenBalance
}

// GetTotalBalance 获取总余额
func (w *Wallet) GetTotalBalance() float64 {
    return w.Balance + w.FrozenBalance
}