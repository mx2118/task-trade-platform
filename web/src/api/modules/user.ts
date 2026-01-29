/**
 * 用户相关API
 */
import { request, type ApiResponse, type PageData } from '../index'
import type { User, Wallet, Transaction, Review, Notification } from '@/types'

export const userApi = {
  /**
   * 获取用户信息
   */
  getUserInfo(): Promise<ApiResponse<User>> {
    return request({
      url: '/user/profile',
      method: 'GET'
    })
  },

  /**
   * 更新用户信息
   */
  updateUserInfo(data: Partial<User>): Promise<ApiResponse<User>> {
    return request({
      url: '/user/profile',
      method: 'PUT',
      data
    })
  },

  /**
   * 更新头像
   */
  updateAvatar(formData: FormData): Promise<ApiResponse<{ avatar: string }>> {
    return request({
      url: '/user/avatar',
      method: 'POST',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 获取用户统计信息
   */
  getStats(): Promise<ApiResponse<{
    completed_tasks: number
    total_earnings: number
    credit: number
    published_tasks: number
    taken_tasks: number
    following_count: number
    followers_count: number
  }>> {
    return request({
      url: '/user/stats',
      method: 'GET'
    })
  },

  /**
   * 获取钱包信息
   */
  getWallet(): Promise<ApiResponse<Wallet>> {
    return request({
      url: '/user/wallet',
      method: 'GET'
    })
  },

  /**
   * 获取交易记录
   */
  getTransactions(params?: {
    type?: string // 'income' | 'expense' | 'freeze' | 'unfreeze'
    start_date?: string
    end_date?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<Transaction>>> {
    return request({
      url: '/user/wallet/transactions',
      method: 'GET',
      params
    })
  },

  /**
   * 获取充值记录
   */
  getRechargeHistory(params?: {
    status?: string
    start_date?: string
    end_date?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<Transaction>>> {
    return request({
      url: '/user/wallet/recharge',
      method: 'GET',
      params
    })
  },

  /**
   * 获取提现记录
   */
  getWithdrawHistory(params?: {
    status?: string
    start_date?: string
    end_date?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<Transaction>>> {
    return request({
      url: '/user/wallet/withdraw',
      method: 'GET',
      params
    })
  },

  /**
   * 充值
   */
  recharge(data: {
    amount: number
    payment_method: string
  }): Promise<ApiResponse<{
    order_no: string
    payment_url: string
    qr_code: string
  }>> {
    return request({
      url: '/user/wallet/recharge',
      method: 'POST',
      data
    })
  },

  /**
   * 提现
   */
  withdraw(data: {
    amount: number
    account_type: string // 'bank' | 'alipay' | 'wechat'
    account_info: {
      bank_name?: string
      account_number?: string
      account_name?: string
      alipay_account?: string
      wechat_account?: string
    }
  }): Promise<ApiResponse<{
    order_no: string
    status: string
  }>> {
    return request({
      url: '/user/wallet/withdraw',
      method: 'POST',
      data
    })
  },

  /**
   * 获取评价列表
   */
  getReviews(params?: {
    type?: string // 'received' | 'given'
    rating?: number
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<Review>>> {
    return request({
      url: '/user/reviews',
      method: 'GET',
      params
    })
  },

  /**
   * 创建评价
   */
  createReview(data: {
    target_user_id: number
    task_id: number
    rating: number
    comment: string
    images?: string[]
  }): Promise<ApiResponse<Review>> {
    return request({
      url: '/user/reviews',
      method: 'POST',
      data
    })
  },

  /**
   * 更新评价
   */
  updateReview(reviewId: number, data: {
    rating?: number
    comment?: string
  }): Promise<ApiResponse<Review>> {
    return request({
      url: `/user/reviews/${reviewId}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除评价
   */
  deleteReview(reviewId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/user/reviews/${reviewId}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取通知列表
   */
  getNotifications(params?: {
    type?: string
    is_read?: boolean
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<Notification>>> {
    return request({
      url: '/user/notifications',
      method: 'GET',
      params
    })
  },

  /**
   * 标记通知为已读
   */
  markNotificationRead(notificationId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/user/notifications/${notificationId}/read`,
      method: 'PUT'
    })
  },

  /**
   * 标记所有通知为已读
   */
  markAllNotificationsRead(): Promise<ApiResponse<void>> {
    return request({
      url: '/user/notifications/read-all',
      method: 'PUT'
    })
  },

  /**
   * 删除通知
   */
  deleteNotification(notificationId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/user/notifications/${notificationId}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取系统公告
   */
  getAnnouncements(params?: {
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<Notification>>> {
    return request({
      url: '/announcements',
      method: 'GET',
      params
    })
  },

  /**
   * 获取公告详情
   */
  getAnnouncementDetail(id: number): Promise<ApiResponse<Notification>> {
    return request({
      url: `/announcements/${id}`,
      method: 'GET'
    })
  },

  /**
   * 获取关注列表
   */
  getFollowing(params?: {
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<User>>> {
    return request({
      url: '/user/following',
      method: 'GET',
      params
    })
  },

  /**
   * 获取粉丝列表
   */
  getFollowers(params?: {
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<User>>> {
    return request({
      url: '/user/followers',
      method: 'GET',
      params
    })
  },

  /**
   * 关注用户
   */
  followUser(userId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/user/follow/${userId}`,
      method: 'POST'
    })
  },

  /**
   * 取消关注用户
   */
  unfollowUser(userId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/user/unfollow/${userId}`,
      method: 'POST'
    })
  },

  /**
   * 获取推荐用户
   */
  getRecommendUsers(params?: {
    limit?: number
    category?: string
  }): Promise<ApiResponse<User[]>> {
    return request({
      url: '/users/recommend',
      method: 'GET',
      params
    })
  },

  /**
   * 搜索用户
   */
  searchUsers(params: {
    keyword: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<User>>> {
    return request({
      url: '/users/search',
      method: 'GET',
      params
    })
  },

  /**
   * 获取用户详情
   */
  getUserDetail(userId: number): Promise<ApiResponse<User>> {
    return request({
      url: `/users/${userId}`,
      method: 'GET'
    })
  },

  /**
   * 用户举报
   */
  reportUser(userId: number, data: {
    reason: string
    description?: string
    images?: string[]
  }): Promise<ApiResponse<void>> {
    return request({
      url: `/users/${userId}/report`,
      method: 'POST',
      data
    })
  },

  /**
   * 修改密码
   */
  changePassword(data: {
    old_password: string
    new_password: string
    confirm_password: string
  }): Promise<ApiResponse<void>> {
    return request({
      url: '/user/password',
      method: 'PUT',
      data
    })
  },

  /**
   * 重置密码
   */
  resetPassword(data: {
    phone: string
    code: string
    new_password: string
  }): Promise<ApiResponse<void>> {
    return request({
      url: '/user/password/reset',
      method: 'POST',
      data
    })
  },

  /**
   * 绑定手机号
   */
  bindPhone(data: {
    phone: string
    code: string
  }): Promise<ApiResponse<void>> {
    return request({
      url: '/user/bind-phone',
      method: 'POST',
      data
    })
  },

  /**
   * 解绑手机号
   */
  unbindPhone(data: {
    password: string
  }): Promise<ApiResponse<void>> {
    return request({
      url: '/user/unbind-phone',
      method: 'POST',
      data
    })
  },

  /**
   * 绑定微信
   */
  bindWechat(data: {
    code: string
    state?: string
  }): Promise<ApiResponse<void>> {
    return request({
      url: '/user/bind-wechat',
      method: 'POST',
      data
    })
  },

  /**
   * 解绑微信
   */
  unbindWechat(): Promise<ApiResponse<void>> {
    return request({
      url: '/user/unbind-wechat',
      method: 'POST'
    })
  },

  /**
   * 绑定支付宝
   */
  bindAlipay(data: {
    auth_code: string
    state?: string
  }): Promise<ApiResponse<void>> {
    return request({
      url: '/user/bind-alipay',
      method: 'POST',
      data
    })
  },

  /**
   * 解绑支付宝
   */
  unbindAlipay(): Promise<ApiResponse<void>> {
    return request({
      url: '/user/unbind-alipay',
      method: 'POST'
    })
  },

  /**
   * 实名认证
   */
  realNameAuth(data: {
    real_name: string
    id_card: string
    id_card_front: string
    id_card_back: string
  }): Promise<ApiResponse<void>> {
    return request({
      url: '/user/real-name-auth',
      method: 'POST',
      data
    })
  },

  /**
   * 获取实名认证状态
   */
  getRealNameAuthStatus(): Promise<ApiResponse<{
    status: string // 'pending' | 'verified' | 'rejected'
    real_name?: string
    id_card?: string
    reject_reason?: string
  }>> {
    return request({
      url: '/user/real-name-auth/status',
      method: 'GET'
    })
  },

  /**
   * 获取用户设置
   */
  getSettings(): Promise<ApiResponse<{
    email_notifications: boolean
    push_notifications: boolean
    message_notifications: boolean
    task_reminders: boolean
    privacy_settings: {
      show_phone: boolean
      show_email: boolean
      show_location: boolean
    }
  }>> {
    return request({
      url: '/user/settings',
      method: 'GET'
    })
  },

  /**
   * 更新用户设置
   */
  updateSettings(data: {
    email_notifications?: boolean
    push_notifications?: boolean
    message_notifications?: boolean
    task_reminders?: boolean
    privacy_settings?: {
      show_phone?: boolean
      show_email?: boolean
      show_location?: boolean
    }
  }): Promise<ApiResponse<void>> {
    return request({
      url: '/user/settings',
      method: 'PUT',
      data
    })
  },

  /**
   * 获取用户标签
   */
  getUserTags(): Promise<ApiResponse<string[]>> {
    return request({
      url: '/user/tags',
      method: 'GET'
    })
  },

  /**
   * 更新用户标签
   */
  updateUserTags(tags: string[]): Promise<ApiResponse<void>> {
    return request({
      url: '/user/tags',
      method: 'PUT',
      data: { tags }
    })
  },

  /**
   * 获取用户证书
   */
  getUserCertificates(): Promise<ApiResponse<{
    id: number
    type: string
    name: string
    image: string
    issued_date: string
    expiry_date?: string
    status: string
  }[]>> {
    return request({
      url: '/user/certificates',
      method: 'GET'
    })
  },

  /**
   * 上传用户证书
   */
  uploadCertificate(data: {
    type: string
    name: string
    image: string
    issued_date: string
    expiry_date?: string
  }): Promise<ApiResponse<void>> {
    return request({
      url: '/user/certificates',
      method: 'POST',
      data
    })
  },

  /**
   * 删除用户证书
   */
  deleteCertificate(certificateId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/user/certificates/${certificateId}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取用户贡献值
   */
  getUserContribution(): Promise<ApiResponse<{
    total_contribution: number
    monthly_contribution: number
    rank: number
    badges: {
      id: number
      name: string
      icon: string
      description: string
      earned_date: string
    }[]
  }>> {
    return request({
      url: '/user/contribution',
      method: 'GET'
    })
  }
}