// 支付相关类型定义

export interface Trade {
  id: number
  orderNo: string
  tradeNo?: string
  userId: number
  taskId?: number
  tradeType: OrderType
  amount: number
  status: PaymentStatus
  payTime?: string
  payMethod?: string
  transactionId?: string
  remark?: string
  expiredAt?: string
  createdAt: string
  updatedAt: string
}

export interface Refund {
  id: number
  refundNo: string
  tradeId: number
  orderNo: string
  userId: number
  taskId?: number
  amount: number
  reason: string
  status: RefundStatus
  refundTime?: string
  createdAt: string
  updatedAt: string
}

export interface Settlement {
  id: number
  taskId: number
  publisherId: number
  takerId: number
  publisherAmount: number
  takerAmount: number
  platformFee: number
  penalty?: number
  status: SettlementStatus
  transferNo?: string
  settleTime?: string
  createdAt: string
  updatedAt: string
}

export interface Wallet {
  userId: number
  balance: number
  frozenBalance: number
  totalIncome: number
  totalExpense: number
  creditScore: number
  level: number
  withdrawLimit: number
  dailyWithdrawLimit: number
  monthlyWithdrawLimit: number
  lastUpdated: string
}

export interface Withdraw {
  id: number
  withdrawNo: string
  userId: number
  amount: number
  fee: number
  actualAmount: number
  accountType: 'alipay' | 'wechat' | 'bank'
  accountNo: string
  accountName: string
  bankCode?: string
  status: WithdrawStatus
  reason?: string
  processedAt?: string
  createdAt: string
  updatedAt: string
}

export interface PaymentMethod {
  id: number
  userId: number
  type: 'alipay' | 'wechat' | 'bank_card'
  name: string
  accountNo: string
  bankCode?: string
  isDefault: boolean
  isVerified: boolean
  createdAt: string
  updatedAt: string
}

// 订单类型
export type OrderType = 
  | 'task_publish'  // 任务发布
  | 'task_take'     // 任务接取
  | 'deposit'       // 保证金
  | 'service_fee'   // 服务费
  | 'penalty'       // 违约金
  | 'withdraw_fee'  // 提现手续费
  | 'refund'        // 退款

// 支付状态
export type PaymentStatus = 
  | 'pending'    // 待支付
  | 'paid'       // 已支付
  | 'failed'     // 支付失败
  | 'cancelled'  // 已取消
  | 'expired'    // 已过期
  | 'refunded'   // 已退款
  | 'settled'    // 已结算

// 退款状态
export type RefundStatus = 
  | 'pending'    // 退款中
  | 'processing' // 处理中
  | 'success'    // 退款成功
  | 'failed'     // 退款失败
  | 'cancelled'  // 退款取消

// 结算状态
export type SettlementStatus = 
  | 'pending'    // 待结算
  | 'processing' // 结算中
  | 'completed'  // 已完成
  | 'failed'     // 失败

// 提现状态
export type WithdrawStatus = 
  | 'pending'    // 待审核
  | 'processing' // 处理中
  | 'completed'  // 已完成
  | 'failed'     // 失败
  | 'rejected'   // 已拒绝

// 支付方式
export type PayMethod = 
  | 'alipay'     // 支付宝
  | 'wechat'     // 微信支付
  | 'bank_card'  // 银行卡
  | 'balance'    // 余额支付
  | 'credit'     // 信用支付

// 预支付请求参数
export interface PrePayParams {
  taskId: number
  orderType: OrderType
  amount: number
  returnURL?: string
  clientIP?: string
  remark?: string
}

// 预支付响应
export interface PrePayResponse {
  orderNo: string
  tradeNo: string
  payURL: string
  qrcode?: string
  amount: number
  expireTime: number
  paymentMethods?: PaymentMethodInfo[]
}

export interface PaymentMethodInfo {
  method: PayMethod
  name: string
  icon: string
  enabled: boolean
  discount?: number
}

// 支付查询参数
export interface PaymentQueryParams {
  orderNo: string
}

// 支付回调数据
export interface PaymentCallbackData {
  orderNo: string
  tradeNo: string
  status: PaymentStatus
  amount: number
  payTime: string
  payMethod: PayMethod
  transactionId: string
  signature: string
  extra?: string
}

// 退款请求参数
export interface RefundParams {
  orderNo: string
  amount: number
  reason: string
}

// 提现请求参数
export interface WithdrawParams {
  amount: number
  accountType: Withdraw['accountType']
  accountNo: string
  accountName: string
  bankCode?: string
  password?: string
}

// 结算请求参数
export interface SettlementParams {
  taskId: number
  publisherAmount: number
  takerAmount: number
}

// 交易查询参数
export interface TransactionQueryParams {
  page?: number
  pageSize?: number
  tradeType?: OrderType
  status?: PaymentStatus
  startTime?: string
  endTime?: string
  minAmount?: number
  maxAmount?: number
}

// 钱包查询参数
export interface WalletQueryParams {
  userId?: number
  includeFrozen?: boolean
  includeStatistics?: boolean
}

// 支付统计
export interface PaymentStatistics {
  totalTrades: number
  totalAmount: number
  successRate: number
  averageAmount: number
  methodStats: PaymentMethodStats[]
  monthlyStats: PaymentMonthlyStats[]
}

export interface PaymentMethodStats {
  method: PayMethod
  count: number
  amount: number
  percentage: number
}

export interface PaymentMonthlyStats {
  month: string
  trades: number
  amount: number
  refunds: number
  refundAmount: number
  growth: number
}

// 账单记录
export interface Bill {
  id: number
  userId: number
  billNo: string
  type: 'income' | 'expense'
  category: string
  amount: number
  balance: number
  description: string
  relatedId?: number
  relatedType?: string
  createdAt: string
}

export interface BillQueryParams {
  page?: number
  pageSize?: number
  type?: Bill['type']
  category?: string
  startTime?: string
  endTime?: string
  minAmount?: number
  maxAmount?: number
}

// 余额变动记录
export interface BalanceChange {
  id: number
  userId: number
  changeType: string
  amount: number
  balanceBefore: number
  balanceAfter: number
  description: string
  relatedId?: number
  relatedType?: string
  createdAt: string
}

// 支付配置
export interface PaymentConfig {
  enableAlipay: boolean
  enableWechat: boolean
  enableBankCard: boolean
  enableBalance: boolean
  minAmount: number
  maxAmount: number
  serviceFeeRate: number
  depositRate: number
  withdrawFee: number
  withdrawMinAmount: number
  withdrawMaxAmount: number
  dailyWithdrawLimit: number
  monthlyWithdrawLimit: number
}

// 收钱吧配置
export interface ShouqianbaConfig {
  appId: string
  merchantNo: string
  sandbox: boolean
  supportedMethods: PayMethod[]
}

// 支付日志
export interface PaymentLog {
  id: number
  type: 'payment' | 'refund' | 'withdraw' | 'settlement'
  action: string
  userId?: number
  orderNo?: string
  amount?: number
  status?: string
  message: string
  request?: any
  response?: any
  ip?: string
  userAgent?: string
  createdAt: string
}