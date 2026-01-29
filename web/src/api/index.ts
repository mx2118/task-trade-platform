/**
 * API 接口统一导出
 */

import request from './request'
import { authApi } from './modules/auth'
import { userApi } from './modules/user'
import { taskApi } from './modules/task'
import { categoryApi } from './modules/category'
import { paymentApi } from './modules/payment'
import { reviewApi } from './modules/review'
import { notificationApi } from './modules/notification'
import { uploadApi } from './modules/upload'

export {
  // 请求实例
  request,
  
  // API 模块
  authApi,
  userApi,
  taskApi,
  categoryApi,
  paymentApi,
  reviewApi,
  notificationApi,
  uploadApi
}

// 导出所有API类型
export type { 
  LoginParams as AuthLoginParams,
  RegisterParams as AuthRegisterParams,
  UserInfo as AuthUserInfo
} from '../types/auth'
export type { 
  UserInfo as UserType,
  Wallet as UserWallet,
  LoginParams as UserLoginParams,
  RegisterParams as UserRegisterParams,
  UserBalance,
  UserTransaction as Transaction,
  Review,
  Notification as UserNotification
} from '../types/user'
export type * from '../types/task'
export type * from '../types/payment'
export type { 
  ApiResponse as CommonApiResponse,
  PageData as CommonPageData,
  PaginationParams as CommonPaginationParams
} from '../types/common'

// API 响应基础类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  timestamp: number
}

export interface PageData<T = any> {
  list: T[]
  total: number
  page: number
  limit: number
  has_more: boolean
}

export interface ErrorInfo {
  code: number
  message: string
  details?: any
  stack?: string
}

// 请求配置类型
export interface RequestConfig {
  url?: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  params?: any
  data?: any
  headers?: Record<string, string>
  timeout?: number
  showLoading?: boolean
  showError?: boolean
}

// 文件上传配置
export interface UploadConfig {
  file: File
  progress?: (progress: number) => void
  cancel?: () => void
}

export default {
  request,
  authApi,
  userApi,
  taskApi,
  categoryApi,
  paymentApi,
  reviewApi,
  notificationApi,
  uploadApi
}