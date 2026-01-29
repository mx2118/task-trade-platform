/**
 * 分类相关API
 */
import { request, type ApiResponse, type PageData } from '../index'
import type { TaskCategory } from '@/types'

export const categoryApi = {
  /**
   * 获取分类列表
   */
  getCategories(): Promise<ApiResponse<TaskCategory[]>> {
    return request({
      url: '/categories',
      method: 'GET'
    })
  },

  /**
   * 获取分类详情
   */
  getCategoryDetail(id: number): Promise<ApiResponse<TaskCategory>> {
    return request({
      url: `/categories/${id}`,
      method: 'GET'
    })
  },

  /**
   * 获取热门分类
   */
  getPopularCategories(params?: {
    limit?: number
  }): Promise<ApiResponse<TaskCategory[]>> {
    return request({
      url: '/categories/popular',
      method: 'GET',
      params
    })
  },

  /**
   * 获取分类下的任务数量
   */
  getCategoryTaskCount(): Promise<ApiResponse<TaskCategory[]>> {
    return request({
      url: '/categories/stats',
      method: 'GET'
    })
  },

  /**
   * 创建分类（管理员）
   */
  createCategory(data: {
    name: string
    description?: string
    icon?: string
    sort_order?: number
  }): Promise<ApiResponse<TaskCategory>> {
    return request({
      url: '/categories',
      method: 'POST',
      data
    })
  },

  /**
   * 更新分类（管理员）
   */
  updateCategory(id: number, data: Partial<TaskCategory>): Promise<ApiResponse<TaskCategory>> {
    return request({
      url: `/categories/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除分类（管理员）
   */
  deleteCategory(id: number): Promise<ApiResponse<void>> {
    return request({
      url: `/categories/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取分类排序
   */
  sortCategories(categoryIds: number[]): Promise<ApiResponse<void>> {
    return request({
      url: '/categories/sort',
      method: 'POST',
      data: { category_ids: categoryIds }
    })
  }
}