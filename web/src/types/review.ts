// 评价相关类型
export interface Review {
  id: number
  taskId: number
  reviewerId: number
  revieweeId: number
  rating: number // 1-5 星评分
  content: string
  tags?: string[]
  isPublic: boolean
  createTime: string
  updateTime: string
  reviewerName: string
  revieweeName: string
  reviewerAvatar?: string
  revieweeAvatar?: string
}

export interface CreateReviewParams {
  taskId: number
  revieweeId: number
  rating: number
  content: string
  tags?: string[]
  isPublic?: boolean
}

export interface ReviewListParams {
  page?: number
  pageSize?: number
  taskId?: number
  reviewerId?: number
  revieweeId?: number
  rating?: number
  keyword?: string
}

export interface ReviewStats {
  averageRating: number
  totalReviews: number
  ratingDistribution: {
    1: number
    2: number
    3: number
    4: number
    5: number
  }
}