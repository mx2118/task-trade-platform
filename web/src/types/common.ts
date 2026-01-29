// API 响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  success: boolean
}

// 分页数据类型
export interface PageData<T = any> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// 请求参数类型
export interface PaginationParams {
  page?: number
  pageSize?: number
  keyword?: string
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}