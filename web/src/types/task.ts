// 任务相关类型定义

export interface Task {
  id: number
  title: string
  description: string
  content: string
  category: string
  tags: string[]
  amount: number
  serviceFeeRate: number
  depositRate: number
  deadline: string
  location?: string
  workMode: 'online' | 'offline' | 'both'
  skillRequirements: string[]
  attachments?: TaskAttachment[]
  publisherId: number
  publisher: TaskUser
  publisher_avatar?: string
  publisher_name?: string
  takerId?: number
  taker?: TaskUser
  status: TaskStatus
  priority: 'low' | 'normal' | 'high' | 'urgent'
  is_urgent: boolean
  is_remote: boolean
  viewCount: number
  view_count: number
  applyCount: number
  maxApplicants?: number
  amount: number
  price: number
  stages?: TaskStage[]
  delivery?: TaskDelivery
  settlement?: TaskSettlement
  createdAt: string
  updatedAt: string
  publishedAt?: string
  takenAt?: string
  completedAt?: string
  expiredAt?: string
}

export interface TaskUser {
  id: number
  username: string
  nickname: string
  avatar?: string
  role: string
  creditScore: number
  successRate: number
  completedTasks: number
  isVerified: boolean
  avgRating: number
  responseTime: number
}

export interface TaskAttachment {
  id: number
  taskId: number
  filename: string
  originalName: string
  fileSize: number
  mimeType: string
  url: string
  uploadedBy: number
  createdAt: string
}

export interface TaskStage {
  id: number
  taskId: number
  stageName: string
  description: string
  amountRatio: number
  deadline?: string
  status: 'pending' | 'in_progress' | 'completed' | 'rejected'
  deliverables?: TaskDeliverable[]
  createdAt: string
  updatedAt: string
  completedAt?: string
}

export interface TaskDeliverable {
  id: number
  stageId: number
  name: string
  description: string
  type: 'file' | 'text' | 'image' | 'video' | 'other'
  required: boolean
  format?: string
  maxSize?: number
  submittedValue?: any
  submittedAt?: string
}

export interface TaskDelivery {
  id: number
  taskId: number
  takerId: number
  content: string
  attachments: TaskAttachment[]
  submittedAt: string
  status: 'pending' | 'accepted' | 'rejected'
  reviewedAt?: string
  reviewComment?: string
  reviewedBy?: number
}

export interface TaskSettlement {
  id: number
  taskId: number
  publisherId: number
  takerId: number
  publisherAmount: number
  takerAmount: number
  platformFee: number
  penalty?: number
  status: 'pending' | 'processing' | 'completed' | 'failed'
  transferNo?: string
  settleTime?: string
  createdAt: string
  updatedAt: string
}

export type TaskStatus = 
  | 'draft'           // 草稿
  | 'pending_audit'   // 待审核
  | 'auditing'        // 审核中
  | 'approved'        // 已通过
  | 'rejected'        // 已拒绝
  | 'open'            // 开放接取
  | 'closed'          // 已关闭
  | 'in_progress'     // 进行中
  | 'pending_delivery' // 待交付
  | 'pending_accept'  // 待验收
  | 'completed'       // 已完成
  | 'cancelled'       // 已取消
  | 'expired'         // 已过期
  | 'disputed'        // 争议中

// 任务筛选参数
export interface TaskFilter {
  keyword?: string
  category?: string
  tags?: string[]
  minAmount?: number
  maxAmount?: number
  workMode?: Task['workMode']
  status?: TaskStatus[]
  priority?: Task['priority']
  location?: string
  skillRequirements?: string[]
  publisherId?: number
  dateRange?: [string, string]
  sortBy?: 'createdAt' | 'amount' | 'deadline' | 'viewCount' | 'applyCount'
  sortOrder?: 'asc' | 'desc'
}

// 任务查询参数
export interface TaskQuery extends TaskFilter {
  page?: number
  pageSize?: number
}

// 任务创建参数
export interface CreateTaskParams {
  title: string
  description: string
  content: string
  category: string
  tags: string[]
  amount: number
  serviceFeeRate?: number
  depositRate?: number
  deadline: string
  location?: string
  workMode: Task['workMode']
  skillRequirements: string[]
  priority?: Task['priority']
  maxApplicants?: number
  stages?: CreateStageParams[]
  attachments?: File[]
}

export interface CreateStageParams {
  stageName: string
  description: string
  amountRatio: number
  deadline?: string
  deliverables?: CreateDeliverableParams[]
}

export interface CreateDeliverableParams {
  name: string
  description: string
  type: TaskDeliverable['type']
  required: boolean
  format?: string
  maxSize?: number
}

// 任务更新参数
export interface UpdateTaskParams extends Partial<CreateTaskParams> {
  id: number
}

// 任务申请参数
export interface ApplyTaskParams {
  taskId: number
  proposal: string
  attachments?: File[]
  estimatedTime?: number
  bidAmount?: number
  message?: string
}

// 任务接取参数
export interface TakeTaskParams {
  taskId: number
  proposal?: string
  estimatedTime?: number
  message?: string
}

// 任务交付参数
export interface DeliverTaskParams {
  taskId: number
  content: string
  attachments?: File[]
  stageId?: number
}

// 任务验收参数
export interface AcceptTaskParams {
  taskId: number
  status: 'accepted' | 'rejected'
  comment?: string
  stageId?: number
  attachments?: File[]
}

// 任务评价参数
export interface ReviewTaskParams {
  taskId: number
  rating: number
  comment: string
  tags?: string[]
  isAnonymous?: boolean
}

export interface TaskReview {
  id: number
  taskId: number
  reviewerId: number
  reviewer: TaskUser
  targetId: number
  target: TaskUser
  rating: number
  comment: string
  tags: string[]
  isAnonymous: boolean
  isPositive: boolean
  helpfulCount: number
  createdAt: string
  updatedAt: string
}

// 任务统计
export interface TaskStatistics {
  totalTasks: number
  publishedTasks: number
  takenTasks: number
  completedTasks: number
  cancelledTasks: number
  disputedTasks: number
  totalAmount: number
  averageAmount: number
  successRate: number
  averageTime: number
  categoryStats: CategoryStats[]
  monthlyStats: MonthlyStats[]
}

export interface CategoryStats {
  category: string
  count: number
  amount: number
  percentage: number
}

export interface MonthlyStats {
  month: string
  published: number
  completed: number
  amount: number
  growth: number
}

// 任务标签
export interface TaskTag {
  id: number
  name: string
  color: string
  description?: string
  usageCount: number
  createdAt: string
}

// 任务分类
export interface TaskCategory {
  id: number
  name: string
  icon: string
  description?: string
  parentId?: number
  children?: TaskCategory[]
  taskCount: number
  isActive: boolean
  createdAt: string
  updatedAt: string
}

// 任务模板
export interface TaskTemplate {
  id: number
  name: string
  description: string
  category: string
  tags: string[]
  defaultAmount: number
  contentTemplate: string
  stagesTemplate: CreateStageParams[]
  isActive: boolean
  usageCount: number
  createdBy: number
  createdAt: string
  updatedAt: string
}

// 任务日志
export interface TaskLog {
  id: number
  taskId: number
  userId: number
  user: TaskUser
  action: string
  content: string
  oldStatus?: string
  newStatus?: string
  metadata?: Record<string, any>
  createdAt: string
}

// 任务申请
export interface TaskApplication {
  id: number
  taskId: number
  applicantId: number
  applicant: TaskUser
  proposal: string
  attachments: TaskAttachment[]
  estimatedTime?: number
  bidAmount?: number
  message?: string
  status: 'pending' | 'accepted' | 'rejected'
  reviewedAt?: string
  reviewComment?: string
  reviewedBy?: number
  createdAt: string
  updatedAt: string
}