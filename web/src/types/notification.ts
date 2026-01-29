// 通知相关类型
export interface Notification {
  id: number
  userId: number
  title: string
  content: string
  type: 'system' | 'task' | 'payment' | 'review' | 'message'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  isRead: boolean
  data?: any // 额外数据
  createTime: string
  updateTime: string
  readTime?: string
}

export interface NotificationListParams {
  page?: number
  pageSize?: number
  type?: string
  isRead?: boolean
  keyword?: string
}

export interface CreateNotificationParams {
  userId: number
  title: string
  content: string
  type: string
  priority?: string
  data?: any
}

export interface NotificationStats {
  total: number
  unread: number
  read: number
  typeDistribution: {
    system: number
    task: number
    payment: number
    review: number
    message: number
  }
}