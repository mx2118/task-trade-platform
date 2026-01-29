/**
 * 任务相关API
 */
import { request, type ApiResponse, type PageData } from '../index'
import type { Task, TaskApplication, TaskCategory } from '@/types'

export const taskApi = {
  /**
   * 获取任务列表
   */
  getTasks(params?: {
    keyword?: string
    category_id?: number
    min_price?: number
    max_price?: number
    sort?: string
    urgent?: boolean
    remote?: boolean
    nearby?: boolean
    recommend?: boolean
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<Task>>> {
    return request({
      url: '/tasks',
      method: 'GET',
      params
    })
  },

  /**
   * 获取任务详情
   */
  getTaskDetail(id: number): Promise<ApiResponse<Task>> {
    return request({
      url: `/tasks/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建任务
   */
  createTask(data: {
    title: string
    description: string
    requirements?: string
    category_id: number
    price: number
    deadline: string
    location?: string
    is_urgent?: boolean
    is_remote?: boolean
    images?: string[]
  }): Promise<ApiResponse<Task>> {
    return request({
      url: '/tasks',
      method: 'POST',
      data
    })
  },

  /**
   * 更新任务
   */
  updateTask(id: number, data: Partial<Task>): Promise<ApiResponse<Task>> {
    return request({
      url: `/tasks/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除任务
   */
  deleteTask(id: number): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 申请接取任务
   */
  applyTask(taskId: number, data: {
    message: string
    estimated_time: string
    contact: string
  }): Promise<ApiResponse<TaskApplication>> {
    return request({
      url: `/tasks/${taskId}/apply`,
      method: 'POST',
      data
    })
  },

  /**
   * 接取任务（直接接取，无需申请）
   */
  takeTask(taskId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/${taskId}/take`,
      method: 'POST'
    })
  },

  /**
   * 交付任务
   */
  deliverTask(taskId: number, data: {
    delivery_content: string
    delivery_images?: string[]
  }): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/${taskId}/deliver`,
      method: 'POST',
      data
    })
  },

  /**
   * 确认任务完成
   */
  confirmTask(taskId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/${taskId}/confirm`,
      method: 'POST'
    })
  },

  /**
   * 取消任务
   */
  cancelTask(taskId: number, reason?: string): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/${taskId}/cancel`,
      method: 'POST',
      data: { reason }
    })
  },

  /**
   * 获取我的任务
   */
  getMyTasks(params?: {
    status?: string
    role?: string // 'publisher' | 'taker'
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<Task>>> {
    return request({
      url: '/user/tasks',
      method: 'GET',
      params
    })
  },

  /**
   * 获取任务统计
   */
  getTaskStats(): Promise<ApiResponse<{
    total: number
    pending: number
    in_progress: number
    completed: number
    cancelled: number
  }>> {
    return request({
      url: '/tasks/stats',
      method: 'GET'
    })
  },

  /**
   * 获取我的任务统计
   */
  getMyTaskStats(): Promise<ApiResponse<{
    total: number
    published: number
    taken: number
    completed: number
    cancelled: number
  }>> {
    return request({
      url: '/user/tasks/stats',
      method: 'GET'
    })
  },

  /**
   * 获取相关任务
   */
  getRelatedTasks(taskId: number, params?: {
    limit?: number
    category_id?: number
  }): Promise<ApiResponse<Task[]>> {
    return request({
      url: `/tasks/${taskId}/related`,
      method: 'GET',
      params
    })
  },

  /**
   * 收藏任务
   */
  favoriteTask(taskId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/${taskId}/favorite`,
      method: 'POST'
    })
  },

  /**
   * 取消收藏任务
   */
  unfavoriteTask(taskId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/${taskId}/unfavorite`,
      method: 'POST'
    })
  },

  /**
   * 分享任务
   */
  shareTask(taskId: number): Promise<ApiResponse<{
    share_url: string
    share_code: string
  }>> {
    return request({
      url: `/tasks/${taskId}/share`,
      method: 'POST'
    })
  },

  /**
   * 获取任务申请列表（任务发布者）
   */
  getTaskApplications(taskId: number, params?: {
    status?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<TaskApplication>>> {
    return request({
      url: `/tasks/${taskId}/applications`,
      method: 'GET',
      params
    })
  },

  /**
   * 接受申请
   */
  acceptApplication(applicationId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/applications/${applicationId}/accept`,
      method: 'POST'
    })
  },

  /**
   * 拒绝申请
   */
  rejectApplication(applicationId: number, reason?: string): Promise<ApiResponse<void>> {
    return request({
      url: `/applications/${applicationId}/reject`,
      method: 'POST',
      data: { reason }
    })
  },

  /**
   * 获取我的申请列表
   */
  getMyApplications(params?: {
    status?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<TaskApplication>>> {
    return request({
      url: '/user/applications',
      method: 'GET',
      params
    })
  },

  /**
   * 搜索任务
   */
  searchTasks(params: {
    keyword: string
    category_id?: number
    min_price?: number
    max_price?: number
    sort?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<PageData<Task>>> {
    return request({
      url: '/tasks/search',
      method: 'GET',
      params
    })
  },

  /**
   * 获取热门任务
   */
  getHotTasks(params?: {
    limit?: number
    category_id?: number
  }): Promise<ApiResponse<Task[]>> {
    return request({
      url: '/tasks/hot',
      method: 'GET',
      params
    })
  },

  /**
   * 获取推荐任务
   */
  getRecommendTasks(params?: {
    limit?: number
    category_id?: number
    exclude_completed?: boolean
  }): Promise<ApiResponse<Task[]>> {
    return request({
      url: '/tasks/recommend',
      method: 'GET',
      params
    })
  },

  /**
   * 获取加急任务
   */
  getUrgentTasks(params?: {
    limit?: number
    category_id?: number
  }): Promise<ApiResponse<Task[]>> {
    return request({
      url: '/tasks/urgent',
      method: 'GET',
      params
    })
  },

  /**
   * 获取附近任务
   */
  getNearbyTasks(params?: {
    latitude?: number
    longitude?: number
    radius?: number
    limit?: number
  }): Promise<ApiResponse<Task[]>> {
    return request({
      url: '/tasks/nearby',
      method: 'GET',
      params
    })
  },

  /**
   * 任务举报
   */
  reportTask(taskId: number, data: {
    reason: string
    description?: string
    images?: string[]
  }): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/${taskId}/report`,
      method: 'POST',
      data
    })
  },

  /**
   * 任务浏览记录
   */
  viewTask(taskId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/${taskId}/view`,
      method: 'POST'
    })
  },

  /**
   * 关注用户（从任务页面）
   */
  followUser(userId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/users/${userId}/follow`,
      method: 'POST'
    })
  },

  /**
   * 取消关注用户
   */
  unfollowUser(userId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/users/${userId}/unfollow`,
      method: 'POST'
    })
  },

  /**
   * 获取任务模板
   */
  getTaskTemplates(): Promise<ApiResponse<TaskTemplate[]>> {
    return request({
      url: '/tasks/templates',
      method: 'GET'
    })
  },

  /**
   * 创建任务草稿
   */
  createDraft(data: Partial<Task>): Promise<ApiResponse<Task>> {
    return request({
      url: '/tasks/draft',
      method: 'POST',
      data
    })
  },

  /**
   * 获取草稿列表
   */
  getDrafts(): Promise<ApiResponse<Task[]>> {
    return request({
      url: '/tasks/drafts',
      method: 'GET'
    })
  },

  /**
   * 删除草稿
   */
  deleteDraft(taskId: number): Promise<ApiResponse<void>> {
    return request({
      url: `/tasks/draft/${taskId}`,
      method: 'DELETE'
    })
  }
}

// 任务模板类型
export interface TaskTemplate {
  id: number
  title: string
  description: string
  requirements: string
  category_id: number
  price: number
  is_system: boolean
  usage_count: number
}