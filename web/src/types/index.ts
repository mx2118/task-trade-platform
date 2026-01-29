// 类型定义统一导出

// 任务相关类型
export * from './task'

// 用户相关类型
export * from './user'

// 支付相关类型
export * from './payment'

// 分类相关类型
export * from './category'

// 认证相关类型
export * from './auth'

// 评价相关类型
export * from './review'

// 通知相关类型
export * from './notification'

// 通用类型
export * from './common'
export interface Pagination {
  page: number
  pageSize: number
  total: number
  totalPages: number
}

export interface PaginatedResponse<T> {
  list: T[]
  pagination: Pagination
}

export interface ApiError {
  code: number
  message: string
  details?: any
  timestamp: number
}

export interface SelectOption {
  label: string
  value: any
  disabled?: boolean
}

export interface FileUpload {
  file: File
  url?: string
  progress?: number
  status?: 'pending' | 'uploading' | 'success' | 'error'
}

export interface Address {
  id: number
  user_id: number
  type: string
  name: string
  phone: string
  province: string
  city: string
  district: string
  address: string
  postal_code?: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface BankCard {
  id: number
  user_id: number
  bank_name: string
  card_number: string
  cardholder_name: string
  is_default: boolean
  status: string
  created_at: string
  updated_at: string
}

export interface NotificationSettings {
  task_notifications: boolean
  payment_notifications: boolean
  system_notifications: boolean
  email_notifications: boolean
  sms_notifications: boolean
  push_notifications: boolean
}

export interface UserStats {
  total_tasks: number
  completed_tasks: number
  total_earnings: number
  success_rate: number
  average_rating: number
  response_time: number
  monthly_stats: MonthlyStats[]
}

export interface MonthlyStats {
  month: string
  year: number
  tasks_completed: number
  earnings: number
  rating: number
}

export interface SystemStats {
  total_users: number
  active_users: number
  total_tasks: number
  completed_tasks: number
  total_amount: number
  today_stats: {
    new_users: number
    new_tasks: number
    completed_tasks: number
    total_amount: number
  }
  category_stats: CategoryStats[]
}