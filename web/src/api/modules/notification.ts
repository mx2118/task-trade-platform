import request from '../request'
import type { ApiResponse, PageData, Notification, NotificationListParams, NotificationStats } from '@/types'

// 获取通知列表
export const getNotifications = (params: NotificationListParams): Promise<ApiResponse<PageData<Notification>>> => {
  return request.get('/api/notifications', { params })
}

// 获取通知详情
export const getNotification = (id: number): Promise<ApiResponse<Notification>> => {
  return request.get(`/api/notifications/${id}`)
}

// 标记通知为已读
export const markAsRead = (ids: number | number[]): Promise<ApiResponse<null>> => {
  return request.post('/api/notifications/read', { ids })
}

// 标记所有通知为已读
export const markAllAsRead = (): Promise<ApiResponse<null>> => {
  return request.post('/api/notifications/read-all')
}

// 删除通知
export const deleteNotifications = (ids: number | number[]): Promise<ApiResponse<null>> => {
  return request.delete('/api/notifications', { data: { ids } })
}

// 获取未读通知数量
export const getUnreadCount = (): Promise<ApiResponse<{ count: number }>> => {
  return request.get('/api/notifications/unread-count')
}

// 获取通知统计
export const getNotificationStats = (): Promise<ApiResponse<NotificationStats>> => {
  return request.get('/api/notifications/stats')
}

// 创建系统通知（管理员）
export const createNotification = (data: {
  title: string
  content: string
  type: string
  priority?: string
  targetUsers?: number[]
  data?: any
}): Promise<ApiResponse<Notification>> => {
  return request.post('/api/notifications', data)
}

// 导出 API 对象
export const notificationApi = {
  getNotifications,
  getNotification,
  markAsRead,
  markAllAsRead,
  deleteNotifications,
  getUnreadCount,
  getNotificationStats,
  createNotification
}