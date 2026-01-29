// 用户相关类型定义

export interface UserInfo {
  id: number
  username: string
  nickname?: string
  avatar?: string
  email?: string
  phone?: string
  role: string
  permissions?: string[]
  roles?: string[]
  creditScore?: number
  balance?: number
  isVerified?: boolean
  realName?: string
  idCard?: string
  alipayAccount?: string
  wechatAccount?: string
  createdAt?: string
  updatedAt?: string
}

export interface LoginParams {
  username?: string
  password?: string
  type?: 'account' | 'wechat' | 'alipay'
  code?: string
}

export interface LoginResponse {
  token: string
  userInfo: UserInfo
  expiresIn?: number
  refreshToken?: string
}

export interface RegisterParams {
  username: string
  password: string
  confirmPassword: string
  email?: string
  phone?: string
  code?: string
  agree: boolean
}

export interface UpdateProfileParams {
  nickname?: string
  avatar?: string
  email?: string
  phone?: string
  realName?: string
  idCard?: string
  alipayAccount?: string
  wechatAccount?: string
}

export interface ChangePasswordParams {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

export interface UserBalance {
  userId: number
  totalIncome: number
  totalExpense: number
  availableBalance: number
  frozenBalance: number
}

export interface UserTransaction {
  id: number
  orderNo: string
  tradeNo?: string
  userId: number
  taskId?: number
  tradeType: string
  amount: number
  status: string
  payTime?: string
  payMethod?: string
  transactionId?: string
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface UserStatistics {
  publishedTasks: number
  takenTasks: number
  completedTasks: number
  totalEarnings: number
  successRate: number
  avgRating: number
  responseTime: number
}

export interface NotificationSettings {
  enable: boolean
  sound: boolean
  desktop: boolean
  email: boolean
  taskReminders: boolean
  paymentNotifications: boolean
  systemAnnouncements: boolean
  marketingMessages: boolean
}

export interface UserPreference {
  theme: 'light' | 'dark'
  language: string
  timezone: string
  currency: string
  autoAccept: boolean
  quickPay: boolean
  showAmount: boolean
  showRealName: boolean
}

export interface UserSession {
  id: string
  userId: number
  device: string
  browser: string
  os: string
  ip: string
  location: string
  loginTime: string
  lastActiveTime: string
  isActive: boolean
}

export interface UserActivity {
  id: number
  userId: number
  action: string
  target: string
  targetId?: number
  description: string
  ip: string
  userAgent: string
  createdAt: string
}

// 用户权限相关
export interface Permission {
  id: number
  name: string
  code: string
  description: string
  type: string
  module: string
  level: number
  parentId?: number
  children?: Permission[]
}

export interface Role {
  id: number
  name: string
  code: string
  description: string
  permissions: string[]
  isDefault: boolean
  createdAt: string
  updatedAt: string
}

// 用户认证相关
export interface CaptchaParams {
  type: 'image' | 'sms' | 'email'
  target: string
}

export interface CaptchaResponse {
  captchaId: string
  captchaImage?: string
  expiresIn: number
}

export interface VerifyCaptchaParams {
  captchaId: string
  code: string
}

// 第三方登录
export interface WechatAuthParams {
  appId: string
  redirectUri: string
  scope?: string
  state?: string
}

export interface WechatAuthResponse {
  code: string
  state?: string
}

export interface AlipayAuthParams {
  appId: string
  redirectUri: string
  scope?: string
  state?: string
}

export interface AlipayAuthResponse {
  auth_code: string
  state?: string
}