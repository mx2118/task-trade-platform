import request from '../request'
import type { ApiResponse, PageData, Review, CreateReviewParams, ReviewListParams, ReviewStats } from '@/types'

// 获取评价列表
export const getReviews = (params: ReviewListParams): Promise<ApiResponse<PageData<Review>>> => {
  return request.get('/api/reviews', { params })
}

// 创建评价
export const createReview = (data: CreateReviewParams): Promise<ApiResponse<Review>> => {
  return request.post('/api/reviews', data)
}

// 获取评价详情
export const getReview = (id: number): Promise<ApiResponse<Review>> => {
  return request.get(`/api/reviews/${id}`)
}

// 更新评价
export const updateReview = (id: number, data: Partial<Review>): Promise<ApiResponse<Review>> => {
  return request.put(`/api/reviews/${id}`, data)
}

// 删除评价
export const deleteReview = (id: number): Promise<ApiResponse<null>> => {
  return request.delete(`/api/reviews/${id}`)
}

// 获取用户收到的评价
export const getUserReceivedReviews = (userId: number, params?: ReviewListParams): Promise<ApiResponse<PageData<Review>>> => {
  return request.get(`/api/reviews/received/${userId}`, { params })
}

// 获取用户发出的评价
export const getUserGivenReviews = (userId: number, params?: ReviewListParams): Promise<ApiResponse<PageData<Review>>> => {
  return request.get(`/api/reviews/given/${userId}`, { params })
}

// 获取任务相关的评价
export const getTaskReviews = (taskId: number, params?: ReviewListParams): Promise<ApiResponse<PageData<Review>>> => {
  return request.get(`/api/reviews/task/${taskId}`, { params })
}

// 获取评价统计
export const getReviewStats = (userId?: number): Promise<ApiResponse<ReviewStats>> => {
  const url = userId ? `/api/reviews/stats/${userId}` : '/api/reviews/stats/me'
  return request.get(url)
}

// 标记评价为已读/未读
export const markReviewAsRead = (id: number, isRead: boolean): Promise<ApiResponse<null>> => {
  return request.patch(`/api/reviews/${id}/read`, { isRead })
}

// 导出 API 对象
export const reviewApi = {
  getReviews,
  createReview,
  getReview,
  updateReview,
  deleteReview,
  getUserReceivedReviews,
  getUserGivenReviews,
  getTaskReviews,
  getReviewStats,
  markReviewAsRead
}