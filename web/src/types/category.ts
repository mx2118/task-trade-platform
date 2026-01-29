// 分类相关类型定义

export interface TaskCategory {
  id: number
  name: string
  icon: string
  description?: string
  parent_id?: number
  children?: TaskCategory[]
  task_count: number
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CategoryStats {
  category_id: number
  category_name: string
  task_count: number
  total_amount: number
  average_amount: number
  completed_count: number
  success_rate: number
}

export interface CreateCategoryParams {
  name: string
  description?: string
  icon?: string
  parent_id?: number
  sort_order?: number
  is_active?: boolean
}

export interface UpdateCategoryParams extends Partial<CreateCategoryParams> {
  id: number
}

export interface CategoryFilter {
  parent_id?: number
  is_active?: boolean
  with_stats?: boolean
}

// 分类树结构
export interface CategoryTree extends TaskCategory {
  children: CategoryTree[]
  level: number
  path: string
}