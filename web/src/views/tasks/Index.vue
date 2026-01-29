<template>
  <div class="tasks-container">
    <div class="tasks-header">
      <div class="header-left">
        <h1>任务大厅</h1>
        <p>发现适合您的任务，开始赚钱之旅</p>
      </div>
      
      <div class="header-right">
        <el-button type="primary" size="large" @click="publishTask">
          <el-icon><Plus /></el-icon>
          发布任务
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="tasks-filters">
      <el-card shadow="never">
        <el-form
          ref="filterFormRef"
          :model="filterForm"
          inline
          class="filter-form"
        >
          <el-form-item label="关键词">
            <el-input
              v-model="filterForm.keyword"
              placeholder="搜索任务标题或描述"
              clearable
              style="width: 240px"
              @keyup.enter="handleSearch"
            />
          </el-form-item>

          <el-form-item label="分类">
            <el-select
              v-model="filterForm.category_id"
              placeholder="选择分类"
              clearable
              style="width: 160px"
            >
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="价格区间">
            <el-input-number
              v-model="filterForm.min_price"
              placeholder="最低"
              :min="0"
              style="width: 100px"
            />
            <span style="margin: 0 8px">-</span>
            <el-input-number
              v-model="filterForm.max_price"
              placeholder="最高"
              :min="0"
              style="width: 100px"
            />
          </el-form-item>

          <el-form-item label="排序">
            <el-select
              v-model="filterForm.sort"
              placeholder="排序方式"
              style="width: 140px"
            >
              <el-option label="最新发布" value="created_at:desc" />
              <el-option label="价格从低到高" value="price:asc" />
              <el-option label="价格从高到低" value="price:desc" />
              <el-option label="截止时间" value="deadline:asc" />
            </el-select>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="handleReset">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 快捷筛选 -->
        <div class="quick-filters">
          <el-radio-group v-model="quickFilter" @change="handleQuickFilter">
            <el-radio-button label="all">全部</el-radio-button>
            <el-radio-button label="urgent">加急</el-radio-button>
            <el-radio-button label="high_price">高价</el-radio-button>
            <el-radio-button label="nearby">附近</el-radio-button>
            <el-radio-button label="recommend">推荐</el-radio-button>
          </el-radio-group>
        </div>
      </el-card>
    </div>

    <!-- 任务列表 -->
    <div class="tasks-content">
      <el-row :gutter="20">
        <el-col :xs="24" :lg="16">
          <div class="task-list">
            <div v-loading="loading" class="list-container">
              <div
                v-for="task in tasks"
                :key="task.id"
                class="task-item"
                @click="goToTaskDetail(task.id)"
              >
                <div class="task-header">
                  <div class="task-title">
                    <h3>{{ task.title }}</h3>
                    <div class="task-badges">
                      <el-tag v-if="task.is_urgent" type="danger" size="small">
                        加急
                      </el-tag>
                      <el-tag v-if="task.is_remote" type="success" size="small">
                        远程
                      </el-tag>
                      <el-tag
                        :type="getTaskStatusType(task.status)"
                        size="small"
                      >
                        {{ getTaskStatusText(task.status) }}
                      </el-tag>
                    </div>
                  </div>
                  <div class="task-price">
                    <span class="price">¥{{ task.price }}</span>
                    <span class="price-label">赏金</span>
                  </div>
                </div>

                <div class="task-description">
                  {{ task.description }}
                </div>

                <div class="task-meta">
                  <div class="meta-left">
                    <el-tag size="small">{{ task.category_name }}</el-tag>
                    <span class="location">
                      <el-icon><Location /></el-icon>
                      {{ task.location || '不限地区' }}
                    </span>
                    <span class="deadline">
                      <el-icon><Clock /></el-icon>
                      截止：{{ formatTime(task.deadline) }}
                    </span>
                  </div>
                  
                  <div class="meta-right">
                    <el-avatar
                      :size="32"
                      :src="task.publisher_avatar"
                    >
                      {{ task.publisher_name?.charAt(0) }}
                    </el-avatar>
                    <span class="publisher-name">{{ task.publisher_name }}</span>
                  </div>
                </div>

                <div class="task-footer">
                  <div class="stats">
                    <span>
                      <el-icon><View /></el-icon>
                      {{ task.view_count }}
                    </span>
                    <span>
                      <el-icon><User /></el-icon>
                      {{ task.apply_count }}
                    </span>
                    <span>{{ formatTime(task.created_at) }}</span>
                  </div>
                  
                  <div class="actions">
                    <el-button
                      v-if="task.status === 'pending'"
                      type="primary"
                      size="small"
                      @click.stop="handleApply(task)"
                    >
                      申请接取
                    </el-button>
                    <el-button
                      size="small"
                      @click.stop="handleShare(task)"
                    >
                      分享
                    </el-button>
                  </div>
                </div>
              </div>

              <!-- 空状态 -->
              <el-empty
                v-if="tasks.length === 0 && !loading"
                description="暂无任务"
              >
                <el-button type="primary" @click="publishTask">
                  发布第一个任务
                </el-button>
              </el-empty>

              <!-- 加载更多 -->
              <div v-if="hasMore && !loading" class="load-more">
                <el-button @click="loadMore">
                  加载更多
                </el-button>
              </div>
            </div>
          </div>
        </el-col>

        <el-col :xs="24" :lg="8">
          <!-- 任务统计 -->
          <el-card class="sidebar-card">
            <template #header>
              <h3>任务统计</h3>
            </template>
            <div class="stats-grid">
              <div class="stat-item">
                <div class="stat-value">{{ stats.total }}</div>
                <div class="stat-label">总任务</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ stats.pending }}</div>
                <div class="stat-label">待接取</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ stats.in_progress }}</div>
                <div class="stat-label">进行中</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ stats.completed }}</div>
                <div class="stat-label">已完成</div>
              </div>
            </div>
          </el-card>

          <!-- 热门分类 -->
          <el-card class="sidebar-card">
            <template #header>
              <h3>热门分类</h3>
            </template>
            <div class="category-list">
              <div
                v-for="category in popularCategories"
                :key="category.id"
                class="category-item"
                @click="filterByCategory(category.id)"
              >
                <span class="category-name">{{ category.name }}</span>
                <span class="category-count">{{ category.task_count }}</span>
              </div>
            </div>
          </el-card>

          <!-- 推荐用户 -->
          <el-card class="sidebar-card">
            <template #header>
              <h3>推荐用户</h3>
            </template>
            <div class="user-list">
              <div
                v-for="user in recommendUsers"
                :key="user.id"
                class="user-item"
                @click="goToUserProfile(user.id)"
              >
                <el-avatar :size="40" :src="user.avatar">
                  {{ user.nickname?.charAt(0) }}
                </el-avatar>
                <div class="user-info">
                  <div class="user-name">{{ user.nickname }}</div>
                  <div class="user-stats">
                    完成 {{ user.completed_tasks }} 个任务
                  </div>
                </div>
                <el-button
                  size="small"
                  @click.stop="followUser(user.id)"
                >
                  关注
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 申请任务弹窗 -->
    <el-dialog
      v-model="applyDialogVisible"
      title="申请接取任务"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="applyFormRef"
        :model="applyForm"
        :rules="applyRules"
        label-width="80px"
      >
        <el-form-item label="申请说明" prop="message">
          <el-input
            v-model="applyForm.message"
            type="textarea"
            :rows="4"
            placeholder="请简要说明您的相关经验和完成此任务的计划"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="预计时间" prop="estimated_time">
          <el-input
            v-model="applyForm.estimated_time"
            placeholder="例如：3天、1周等"
          />
        </el-form-item>

        <el-form-item label="联系方式" prop="contact">
          <el-input
            v-model="applyForm.contact"
            placeholder="您的联系方式"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="applyDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="applyLoading"
          @click="handleApplySubmit"
        >
          提交申请
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Location,
  Clock,
  View,
  User
} from '@element-plus/icons-vue'
import { taskApi, categoryApi, userApi } from '@/api'
import { formatTime } from '@/utils/format'
import type { Task, Category } from '@/types'

const router = useRouter()

// 筛选表单
const filterForm = reactive({
  keyword: '',
  category_id: undefined as number | undefined,
  min_price: undefined as number | undefined,
  max_price: undefined as number | undefined,
  sort: 'created_at:desc'
})

// 快捷筛选
const quickFilter = ref('all')

// 状态变量
const loading = ref(false)
const hasMore = ref(true)
const currentPage = ref(1)
const pageSize = 20

// 数据
const tasks = ref<Task[]>([])
const categories = ref<Category[]>([])
const stats = reactive({
  total: 0,
  pending: 0,
  in_progress: 0,
  completed: 0
})
const popularCategories = ref<Category[]>([])
const recommendUsers = ref<any[]>([])

// 申请任务相关
const applyDialogVisible = ref(false)
const applyLoading = ref(false)
const currentTask = ref<Task | null>(null)
const applyFormRef = ref()

const applyForm = reactive({
  message: '',
  estimated_time: '',
  contact: ''
})

const applyRules = {
  message: [
    { required: true, message: '请输入申请说明', trigger: 'blur' },
    { min: 10, message: '申请说明至少10个字符', trigger: 'blur' }
  ],
  estimated_time: [
    { required: true, message: '请输入预计完成时间', trigger: 'blur' }
  ],
  contact: [
    { required: true, message: '请输入联系方式', trigger: 'blur' }
  ]
}

// 加载任务列表
const loadTasks = async (page = 1, append = false) => {
  try {
    loading.value = true
    
    const params = {
      page,
      limit: pageSize,
      ...filterForm,
      ...(quickFilter.value !== 'all' && { [quickFilter.value]: true })
    }

    const response = await taskApi.getTasks(params)
    
    if (append) {
      tasks.value = [...tasks.value, ...response.data.list]
    } else {
      tasks.value = response.data.list
    }

    hasMore.value = response.data.list.length === pageSize
    currentPage.value = page
    
  } catch (error: any) {
    ElMessage.error('加载任务失败')
  } finally {
    loading.value = false
  }
}

// 加载分类
const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategories()
    categories.value = response.data
  } catch (error: any) {
    console.error('加载分类失败:', error)
  }
}

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await taskApi.getStats()
    Object.assign(stats, response.data)
  } catch (error: any) {
    console.error('加载统计失败:', error)
  }
}

// 加载热门分类
const loadPopularCategories = async () => {
  try {
    const response = await categoryApi.getPopularCategories()
    popularCategories.value = response.data
  } catch (error: any) {
    console.error('加载热门分类失败:', error)
  }
}

// 加载推荐用户
const loadRecommendUsers = async () => {
  try {
    const response = await userApi.getRecommendUsers()
    recommendUsers.value = response.data
  } catch (error: any) {
    console.error('加载推荐用户失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  loadTasks(1, false)
}

// 重置
const handleReset = () => {
  Object.assign(filterForm, {
    keyword: '',
    category_id: undefined,
    min_price: undefined,
    max_price: undefined,
    sort: 'created_at:desc'
  })
  quickFilter.value = 'all'
  handleSearch()
}

// 快捷筛选
const handleQuickFilter = () => {
  currentPage.value = 1
  loadTasks(1, false)
}

// 加载更多
const loadMore = () => {
  if (hasMore.value && !loading.value) {
    loadTasks(currentPage.value + 1, true)
  }
}

// 按分类筛选
const filterByCategory = (categoryId: number) => {
  filterForm.category_id = categoryId
  handleSearch()
}

// 跳转到任务详情
const goToTaskDetail = (id: number) => {
  router.push(`/tasks/${id}`)
}

// 发布任务
const publishTask = () => {
  router.push('/tasks/publish')
}

// 跳转到用户资料
const goToUserProfile = (id: number) => {
  router.push(`/user/${id}`)
}

// 关注用户
const followUser = async (userId: number) => {
  try {
    await userApi.followUser(userId)
    ElMessage.success('关注成功')
  } catch (error: any) {
    ElMessage.error(error.message || '关注失败')
  }
}

// 申请接取任务
const handleApply = (task: Task) => {
  currentTask.value = task
  applyDialogVisible.value = true
}

// 提交申请
const handleApplySubmit = async () => {
  if (!applyFormRef.value || !currentTask.value) return

  try {
    await applyFormRef.value.validate()
    
    applyLoading.value = true
    
    await taskApi.applyTask(currentTask.value.id, applyForm)
    
    ElMessage.success('申请已提交，请等待任务发布者确认')
    applyDialogVisible.value = false
    
    // 重新加载任务列表
    loadTasks(currentPage.value, false)
    
  } catch (error: any) {
    if (error.message) {
      ElMessage.error(error.message)
    }
  } finally {
    applyLoading.value = false
  }
}

// 分享任务
const handleShare = async (task: Task) => {
  const shareUrl = `${window.location.origin}/tasks/${task.id}`
  
  try {
    await navigator.clipboard.writeText(shareUrl)
    ElMessage.success('链接已复制到剪贴板')
  } catch (error) {
    ElMessageBox.alert(
      `任务链接：${shareUrl}`,
      '分享任务',
      {
        confirmButtonText: '确定',
        type: 'info'
      }
    )
  }
}

// 获取任务状态类型
const getTaskStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: 'info',
    in_progress: 'primary',
    completed: 'success',
    cancelled: 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取任务状态文本
const getTaskStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待接取',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || '未知'
}

// 页面初始化
onMounted(() => {
  loadTasks()
  loadCategories()
  loadStats()
  loadPopularCategories()
  loadRecommendUsers()
})
</script>

<style lang="scss" scoped>
.tasks-container {
  padding: 24px;

  .tasks-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 32px;

    .header-left {
      h1 {
        margin: 0 0 8px 0;
        font-size: 28px;
        font-weight: 600;
        color: #1a1a1a;
      }

      p {
        margin: 0;
        font-size: 16px;
        color: #666;
      }
    }

    .header-right {
      .el-button {
        height: 48px;
        padding: 0 24px;
        font-size: 16px;

        .el-icon {
          margin-right: 8px;
        }
      }
    }
  }

  .tasks-filters {
    margin-bottom: 24px;

    .filter-form {
      .el-form-item {
        margin-bottom: 16px;
      }
    }

    .quick-filters {
      margin-top: 16px;
    }
  }

  .tasks-content {
    .task-list {
      .list-container {
        .task-item {
          background: white;
          border: 1px solid #f0f0f0;
          border-radius: 12px;
          padding: 20px;
          margin-bottom: 16px;
          cursor: pointer;
          transition: all 0.3s ease;

          &:hover {
            border-color: #1890ff;
            box-shadow: 0 4px 12px rgba(24, 144, 255, 0.1);
            transform: translateY(-2px);
          }

          .task-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 12px;

            .task-title {
              flex: 1;

              h3 {
                margin: 0 0 8px 0;
                font-size: 18px;
                font-weight: 600;
                color: #1a1a1a;
                line-height: 1.4;
              }

              .task-badges {
                display: flex;
                gap: 8px;
                flex-wrap: wrap;
              }
            }

            .task-price {
              text-align: right;
              margin-left: 16px;

              .price {
                display: block;
                font-size: 24px;
                font-weight: 700;
                color: #ff4d4f;
                line-height: 1;
              }

              .price-label {
                font-size: 12px;
                color: #999;
                margin-top: 4px;
              }
            }
          }

          .task-description {
            color: #666;
            line-height: 1.6;
            margin-bottom: 16px;
            overflow: hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
          }

          .task-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 16px;

            .meta-left {
              display: flex;
              align-items: center;
              gap: 16px;
              flex-wrap: wrap;

              .location,
              .deadline {
                display: flex;
                align-items: center;
                gap: 4px;
                font-size: 14px;
                color: #666;

                .el-icon {
                  font-size: 16px;
                }
              }
            }

            .meta-right {
              display: flex;
              align-items: center;
              gap: 8px;

              .publisher-name {
                font-size: 14px;
                color: #666;
              }
            }
          }

          .task-footer {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding-top: 16px;
            border-top: 1px solid #f5f5f5;

            .stats {
              display: flex;
              gap: 16px;
              font-size: 12px;
              color: #999;

              span {
                display: flex;
                align-items: center;
                gap: 4px;

                .el-icon {
                  font-size: 14px;
                }
              }
            }

            .actions {
              display: flex;
              gap: 8px;
            }
          }
        }

        .load-more {
          text-align: center;
          padding: 20px 0;
        }
      }
    }

    .sidebar-card {
      margin-bottom: 20px;

      h3 {
        margin: 0;
        font-size: 16px;
        font-weight: 600;
        color: #1a1a1a;
      }

      .stats-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 16px;

        .stat-item {
          text-align: center;
          padding: 16px;
          background: #fafafa;
          border-radius: 8px;

          .stat-value {
            font-size: 24px;
            font-weight: 600;
            color: #1890ff;
            margin-bottom: 4px;
          }

          .stat-label {
            font-size: 12px;
            color: #666;
          }
        }
      }

      .category-list {
        .category-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 12px 0;
          border-bottom: 1px solid #f0f0f0;
          cursor: pointer;
          transition: all 0.3s ease;

          &:hover {
            background: #fafafa;
            margin: 0 -12px;
            padding: 12px;
          }

          &:last-child {
            border-bottom: none;
          }

          .category-name {
            font-size: 14px;
            color: #1a1a1a;
          }

          .category-count {
            font-size: 12px;
            color: #999;
            background: #f0f0f0;
            padding: 2px 8px;
            border-radius: 4px;
          }
        }
      }

      .user-list {
        .user-item {
          display: flex;
          align-items: center;
          padding: 12px 0;
          border-bottom: 1px solid #f0f0f0;
          cursor: pointer;
          transition: all 0.3s ease;

          &:hover {
            background: #fafafa;
            margin: 0 -12px;
            padding: 12px;
          }

          &:last-child {
            border-bottom: none;
          }

          .user-info {
            flex: 1;
            margin: 0 12px;

            .user-name {
              font-size: 14px;
              font-weight: 500;
              color: #1a1a1a;
              margin-bottom: 2px;
            }

            .user-stats {
              font-size: 12px;
              color: #999;
            }
          }
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .tasks-container {
    padding: 16px;

    .tasks-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 16px;

      .header-right {
        width: 100%;

        .el-button {
          width: 100%;
        }
      }
    }

    .tasks-filters {
      .filter-form {
        .el-form-item {
          width: 100%;
          margin-right: 0;
        }
      }
    }

    .tasks-content {
      .task-list {
        .task-item {
          .task-header {
            flex-direction: column;

            .task-price {
              margin-left: 0;
              margin-top: 12px;
              text-align: left;
            }
          }

          .task-meta {
            flex-direction: column;
            align-items: flex-start;
            gap: 12px;
          }

          .task-footer {
            flex-direction: column;
            gap: 12px;

            .stats {
              order: 2;
            }

            .actions {
              order: 1;
              width: 100%;
              justify-content: flex-end;
            }
          }
        }
      }
    }
  }
}
</style>