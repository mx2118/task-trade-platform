import request from '../request'
import type { ApiResponse, PageData } from '@/types'

// 支付相关接口
export interface PaymentMethod {
  id: number
  name: string
  type: 'alipay' | 'wechat' | 'bank'
  account: string
  isDefault: boolean
  createTime: string
}

export interface PaymentOrder {
  id: number
  orderNo: string
  amount: number
  status: 'pending' | 'paid' | 'failed' | 'refunded'
  paymentMethod: string
  description: string
  createTime: string
  payTime?: string
  expireTime?: string
}

export interface WithdrawalRequest {
  id: number
  amount: number
  status: 'pending' | 'processing' | 'completed' | 'failed'
  method: string
  account: string
  fee: number
  actualAmount: number
  description?: string
  createTime: string
  processTime?: string
}

// 获取支付方式列表
export const getPaymentMethods = (): Promise<ApiResponse<PaymentMethod[]>> => {
  return request.get('/api/payment/methods')
}

// 创建支付订单
export const createPayment = (data: {
  taskId: number
  amount: number
  method: string
  description?: string
}): Promise<ApiResponse<{ orderNo: string; paymentUrl?: string; qrCode?: string }>> => {
  return request.post('/api/payment/create', data)
}

// 查询支付状态
export const getPaymentStatus = (orderNo: string): Promise<ApiResponse<PaymentOrder>> => {
  return request.get(`/api/payment/status/${orderNo}`)
}

// 获取支付记录
export const getPaymentHistory = (params: {
  page?: number
  pageSize?: number
  status?: string
  startTime?: string
  endTime?: string
}): Promise<ApiResponse<PageData<PaymentOrder>>> => {
  return request.get('/api/payment/history', { params })
}

// 提现申请
export const createWithdrawal = (data: {
  amount: number
  method: string
  account: string
  description?: string
}): Promise<ApiResponse<{ id: number }>> => {
  return request.post('/api/payment/withdrawal', data)
}

// 获取提现记录
export const getWithdrawalHistory = (params: {
  page?: number
  pageSize?: number
  status?: string
}): Promise<ApiResponse<PageData<WithdrawalRequest>>> => {
  return request.get('/api/payment/withdrawal-history', { params })
}

// 取消提现申请
export const cancelWithdrawal = (id: number): Promise<ApiResponse<null>> => {
  return request.post(`/api/payment/withdrawal/${id}/cancel`)
}

// 获取账户余额
export const getBalance = (): Promise<ApiResponse<{
  available: number
  frozen: number
  total: number
}>> => {
  return request.get('/api/payment/balance')
}

// 导出 API 对象
export const paymentApi = {
  getPaymentMethods,
  createPayment,
  getPaymentStatus,
  getPaymentHistory,
  createWithdrawal,
  getWithdrawalHistory,
  cancelWithdrawal,
  getBalance
}