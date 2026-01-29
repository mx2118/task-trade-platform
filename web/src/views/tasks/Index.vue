<template>
  <div class="tasks-container">
    <!-- 当前分类提示 -->
    <div v-if="currentCategoryName" class="category-breadcrumb">
      <el-icon><Folder /></el-icon>
      <span>当前分类：{{ currentCategoryName }}</span>
    </div>

    <!-- 任务统计 -->
    <div class="tasks-stats-bar">
      <el-button 
        class="stat-btn stat-total" 
        :class="{ active: statusFilter === 'all' }"
        @click="handleStatusFilter('all')"
      >
        总任务 <span class="stat-value">{{ stats.total }}</span>
      </el-button>
      <el-button 
        class="stat-btn stat-pending"
        :class="{ active: statusFilter === 'pending' }"
        @click="handleStatusFilter('pending')"
      >
        待接取 <span class="stat-value">{{ stats.pending }}</span>
      </el-button>
      <el-button 
        class="stat-btn stat-progress"
        :class="{ active: statusFilter === 'in_progress' }"
        @click="handleStatusFilter('in_progress')"
      >
        进行中 <span class="stat-value">{{ stats.in_progress }}</span>
      </el-button>
      <el-button 
        class="stat-btn stat-reviewing"
        :class="{ active: statusFilter === 'reviewing' }"
        @click="handleStatusFilter('reviewing')"
      >
        待验收 <span class="stat-value">{{ stats.reviewing }}</span>
      </el-button>
      <el-button 
        class="stat-btn stat-completed"
        :class="{ active: statusFilter === 'completed' }"
        @click="handleStatusFilter('completed')"
      >
        已完成 <span class="stat-value">{{ stats.completed }}</span>
      </el-button>
      <el-button 
        class="stat-btn stat-cancelled"
        :class="{ active: statusFilter === 'cancelled' }"
        @click="handleStatusFilter('cancelled')"
      >
        已取消 <span class="stat-value">{{ stats.cancelled }}</span>
      </el-button>
    </div>

    <!-- 搜索筛选模态窗口 -->
    <el-dialog
      v-model="searchDialogVisible"
      title="任务检索"
      width="90%"
      :close-on-click-modal="false"
      destroy-on-close
      style="max-width: 800px"
      class="search-dialog"
      :modal-class="'search-modal-mask'"
    >
      <el-form
        ref="filterFormRef"
        :model="filterForm"
        label-width="80px"
        class="filter-form"
      >
        <el-form-item label="关键词">
          <el-input
            v-model="filterForm.keyword"
            placeholder="搜索任务标题或描述"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>

        <el-form-item label="分类">
          <el-select
            v-model="filterForm.category_id"
            placeholder="选择分类"
            clearable
            style="width: 100%"
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
          <div style="display: flex; align-items: center; gap: 8px;">
            <el-input-number
              v-model="filterForm.min_price"
              placeholder="最低"
              :min="0"
              style="flex: 1"
            />
            <span>-</span>
            <el-input-number
              v-model="filterForm.max_price"
              placeholder="最高"
              :min="0"
              style="flex: 1"
            />
          </div>
        </el-form-item>

        <el-form-item label="排序方式">
          <el-select
            v-model="filterForm.sort"
            placeholder="排序方式"
            style="width: 100%"
          >
            <el-option label="最新发布" value="created_at:desc" />
            <el-option label="价格从低到高" value="price:asc" />
            <el-option label="价格从高到低" value="price:desc" />
            <el-option label="截止时间" value="deadline:asc" />
          </el-select>
        </el-form-item>

        <!-- 快捷筛选 -->
        <el-form-item label="快捷筛选">
          <el-radio-group v-model="quickFilter" @change="handleQuickFilter">
            <el-radio-button value="all">全部</el-radio-button>
            <el-radio-button value="urgent">加急</el-radio-button>
            <el-radio-button value="high_price">高价</el-radio-button>
            <el-radio-button value="nearby">附近</el-radio-button>
            <el-radio-button value="recommend">推荐</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="handleReset">
          <el-icon><Refresh /></el-icon>
          重置
        </el-button>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
      </template>
    </el-dialog>

    <!-- 任务列表 -->
    <div class="tasks-content">
      <el-row :gutter="20">
        <el-col :xs="24" :lg="24">
          <div class="task-list">
            <div v-loading="loading" class="list-container">
              <!-- 测试：直接显示任务数量 -->
              <div v-if="!loading && tasks.length === 0" style="padding: 40px; text-align: center; background: #fff; border-radius: 8px;">
                <p style="font-size: 18px; color: #666;">暂无任务数据</p>
              </div>

              <!-- 简化的任务卡片 -->
              <div
                v-for="task in tasks"
                :key="task.id"
                class="task-item"
                style="background: white; padding: 20px; margin-bottom: 16px; border-radius: 8px; border: 1px solid #e0e0e0;"
              >
                <h3 style="margin: 0 0 10px 0; font-size: 18px;">{{ task.title }}</h3>
                <p style="margin: 0; color: #666;">{{ task.description }}</p>
                <div style="margin-top: 10px; display: flex; justify-content: space-between; align-items: center;">
                  <span style="color: #ff4d4f; font-size: 20px; font-weight: bold;">¥{{ task.price }}</span>
                  <el-tag :class="'task-status-' + task.status" size="small">
                    {{ getTaskStatusText(task.status) }}
                  </el-tag>
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

              <!-- 移动端加载更多 -->
              <div v-if="total > 0 && isMobile" class="mobile-load-more">
                <div class="task-count-info">
                  已加载 {{ tasks.length }} / {{ total }} 个任务
                </div>
                <el-button
                  v-if="hasMore"
                  type="primary"
                  size="large"
                  :loading="loading"
                  class="load-more-btn"
                  @click="loadMore"
                >
                  <el-icon v-if="!loading"><ArrowDown /></el-icon>
                  {{ loading ? '加载中...' : '加载更多' }}
                </el-button>
                <div v-else class="no-more-tip">
                  已加载全部任务
                  <el-button text @click="showAllData">查看分页版本</el-button>
                </div>
              </div>

              <!-- 桌面端分页 -->
              <div v-if="total > 0 && !isMobile" class="pagination-wrapper" style="padding: 20px; text-align: center;">
                <el-pagination
                  v-model:current-page="currentPage"
                  v-model:page-size="pageSize"
                  :page-sizes="[10, 20, 50, 100]"
                  :total="total"
                  layout="total, sizes, prev, pager, next, jumper"
                  background
                  @size-change="handleSizeChange"
                  @current-change="handlePageChange"
                />
              </div>
            </div>
          </div>
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
import { ref, reactive, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Location,
  Clock,
  View,
  User,
  ArrowDown,
  Folder,
  Close
} from '@element-plus/icons-vue'
import { taskApi, categoryApi, userApi } from '@/api'
import { formatTime } from '@/utils/format'
import type { Task, Category } from '@/types'

const router = useRouter()
const route = useRoute()

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

// 状态筛选
const statusFilter = ref('all')

// 状态变量
const loading = ref(false)
const hasMore = ref(true)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const isMobile = ref(false)

// 数据
const tasks = ref<Task[]>([])
const categories = ref<Category[]>([])
const stats = reactive({
  total: 0,
  pending: 0,
  in_progress: 0,
  reviewing: 0,
  completed: 0,
  cancelled: 0
})
const popularCategories = ref<Category[]>([])
const recommendUsers = ref<any[]>([])

// 当前分类名称
const currentCategoryName = computed(() => {
  if (!filterForm.category_id) return ''
  const category = categories.value.find(c => c.id === filterForm.category_id)
  return category?.name || ''
})

// 清除分类筛选
const clearCategoryFilter = async () => {
  filterForm.category_id = undefined
  await router.push('/tasks')
}

// 申请任务相关
const applyDialogVisible = ref(false)
const applyLoading = ref(false)
const currentTask = ref<Task | null>(null)
const applyFormRef = ref()

// 搜索对话框
const searchDialogVisible = ref(false)

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
    
    // 将前端状态转换为后端状态码（字符串格式）
    const statusMap: Record<string, string> = {
      'pending': '1',
      'in_progress': '2',
      'reviewing': '3',
      'completed': '4',
      'cancelled': '5'
    }
    
    const params: any = {
      page,
      limit: pageSize.value,
      ...filterForm,
      ...(quickFilter.value !== 'all' && { [quickFilter.value]: true })
    }
    
    // 添加状态筛选参数（转换为后端字符串格式）
    if (statusFilter.value !== 'all') {
      params.status = statusMap[statusFilter.value]
    }

    // 设置10秒超时
    const timeoutPromise = new Promise((_, reject) => {
      setTimeout(() => reject(new Error('请求超时')), 10000)
    })
    
    const response = await Promise.race([
      taskApi.getTasks(params),
      timeoutPromise
    ]) as any
    
    // 后端返回的数据在 response.data.data.tasks
    const responseData = response.data?.data || response.data || {}
    
    let taskList = responseData.tasks || responseData.list || []
    
    // 更新总数
    total.value = responseData.total || 0
    
    // 映射后端字段到前端期望的字段
    taskList = taskList.map((task: any) => ({
      ...task,
      id: task.task_id || task.id,
      description: task.content || task.description,
      created_at: task.create_time || task.created_at,
      price: task.amount || task.price,
      status: task.status === 1 ? 'pending' : task.status === 2 ? 'in_progress' : task.status === 3 ? 'reviewing' : task.status === 4 ? 'completed' : 'cancelled'
    }))
    
    if (append) {
      tasks.value = [...tasks.value, ...taskList]
    } else {
      tasks.value = taskList
    }

    hasMore.value = taskList.length === pageSize.value
    currentPage.value = page
    
  } catch (error: any) {
    console.error('❌ 加载任务失败:', error)
    ElMessage.error('加载任务失败: ' + (error.message || '网络错误'))
    // 即使失败也显示空列表，不要一直卡在加载状态
    if (!append) {
      tasks.value = []
      total.value = 0
    }
  } finally {
    loading.value = false
  }
}

// 加载分类
const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategories()
    const responseData = response.data?.data || response.data || {}
    categories.value = responseData.categories || responseData.list || []
  } catch (error: any) {
    console.error('加载分类失败:', error)
  }
}

// 加载统计数据（尝试从API获取，失败则使用任务列表计算）
const loadStats = async () => {
  try {
    // 构建参数，包含当前分类筛选
    const params: any = {}
    if (filterForm.category_id) {
      params.category_id = filterForm.category_id
    }
    
    const response = await taskApi.getTaskStats(params)
    const responseData = response.data?.data || response.data || {}
    
    // 如果API返回了有效数据，则使用
    if (responseData.total !== undefined) {
      Object.assign(stats, {
        total: responseData.total || 0,
        pending: responseData.pending || 0,
        in_progress: responseData.in_progress || responseData.inProgress || 0,
        reviewing: responseData.reviewing || 0,
        completed: responseData.completed || 0,
        cancelled: responseData.cancelled || 0
      })
      return true
    }
    return false
  } catch (error: any) {
    return false
  }
}

// 从所有任务加载统计数据（根据当前分类筛选条件）
const loadStatsFromAllTasks = async () => {
  try {
    // 构建查询参数，包含当前分类筛选
    const params: any = {
      page: 1,
      limit: 9999
    }
    
    // 如果有分类筛选，添加到参数中
    if (filterForm.category_id) {
      params.category_id = filterForm.category_id
    }
    
    const response = await taskApi.getTasks(params) as any
    const responseData = response.data?.data || response.data || {}
    let taskList = responseData.tasks || responseData.list || []
    
    // 映射状态
    taskList = taskList.map((task: any) => ({
      ...task,
      status: task.status === 1 ? 'pending' : task.status === 2 ? 'in_progress' : task.status === 3 ? 'reviewing' : task.status === 4 ? 'completed' : 'cancelled'
    }))
    
    const pending = taskList.filter((t: any) => t.status === 'pending').length
    const inProgress = taskList.filter((t: any) => t.status === 'in_progress').length
    const reviewing = taskList.filter((t: any) => t.status === 'reviewing').length
    const completed = taskList.filter((t: any) => t.status === 'completed').length
    const cancelled = taskList.filter((t: any) => t.status === 'cancelled').length
    
    Object.assign(stats, {
      total: taskList.length,
      pending,
      in_progress: inProgress,
      reviewing,
      completed,
      cancelled
    })
  } catch (error: any) {
    console.error('❌ 加载统计数据失败:', error)
  }
}

// 加载热门分类
const loadPopularCategories = async () => {
  try {
    const response = await categoryApi.getPopularCategories()
    const responseData = response.data?.data || response.data || {}
    popularCategories.value = responseData.categories || responseData.list || []
  } catch (error: any) {
    console.error('加载热门分类失败:', error)
    // 使用已加载的分类数据
    popularCategories.value = categories.value.filter((c: any) => c.task_count > 0)
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
  searchDialogVisible.value = false // 关闭对话框
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
  // 不关闭对话框，允许继续调整筛选条件
}

// 状态筛选
const handleStatusFilter = async (status: string) => {
  statusFilter.value = status
  currentPage.value = 1
  await loadTasks(1, false)
  
  // 状态筛选不影响统计数据（统计始终显示当前分类下的全部状态）
  // 所以不需要重新加载统计
}

// 加载更多
const loadMore = () => {
  if (hasMore.value && !loading.value) {
    loadTasks(currentPage.value + 1, true)
  }
}

// 分页处理
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadTasks(page, false)
  // 滚动到顶部
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadTasks(1, false)
  // 滚动到顶部
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 检测移动端
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}

// 滚动处理（移动端无限滚动）
const handleScroll = () => {
  if (!isMobile.value || loading.value || !hasMore.value) return
  
  const scrollHeight = document.documentElement.scrollHeight
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop
  const clientHeight = document.documentElement.clientHeight
  
  // 距离底部150px时开始加载
  if (scrollHeight - scrollTop - clientHeight < 150) {
    loadMore()
  }
}

// 查看全部数据（切换到PC分页模式）
const showAllData = () => {
  isMobile.value = false
  currentPage.value = Math.ceil(tasks.value.length / pageSize.value)
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

// 获取任务状态类型（返回空字符串，使用自定义颜色）
const getTaskStatusType = (status: string) => {
  return ''
}

// 获取任务状态文本
const getTaskStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待接取',
    in_progress: '进行中',
    reviewing: '待验收',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || '未知'
}

// 页面初始化
onMounted(async () => {
  // 检测移动端
  checkMobile()
  window.addEventListener('resize', checkMobile)
  
  // 添加滚动监听（用于移动端无限滚动）
  if (isMobile.value) {
    window.addEventListener('scroll', handleScroll)
  }
  
  // 监听全局搜索对话框事件
  const handleOpenSearch = () => {
    searchDialogVisible.value = true
  }
  window.addEventListener('open-search-dialog', handleOpenSearch)
  
  // 从路由参数读取初始状态筛选
  const routeStatus = route.query.status as string
  if (routeStatus && ['pending', 'in_progress', 'reviewing', 'completed', 'cancelled'].includes(routeStatus)) {
    statusFilter.value = routeStatus
  }
  
  // 从路由参数读取初始分类筛选
  const routeCategoryId = route.query.category_id as string
  if (routeCategoryId) {
    filterForm.category_id = parseInt(routeCategoryId)
  }
  
  // 加载关键数据
  await Promise.all([
    loadTasks(),
    loadCategories()
  ])
  
  // 加载统计数据（优先使用API，失败则从所有任务计算）
  const statsLoaded = await loadStats()
  if (!statsLoaded) {
    await loadStatsFromAllTasks()
  }
  
  // 加载非关键数据（失败不影响页面显示）
  Promise.allSettled([
    loadPopularCategories(),
    loadRecommendUsers()
  ]).catch(err => console.warn('加载辅助数据失败:', err))
})

// 监听路由参数变化
watch(() => route.query.status, (newStatus) => {
  if (newStatus && typeof newStatus === 'string' && ['pending', 'in_progress', 'reviewing', 'completed', 'cancelled'].includes(newStatus)) {
    statusFilter.value = newStatus
    loadTasks(1, false)
  } else if (!newStatus) {
    statusFilter.value = 'all'
    loadTasks(1, false)
  }
})

// 监听路由分类参数变化
watch(() => route.query.category_id, async (newCategoryId) => {
  if (newCategoryId) {
    filterForm.category_id = parseInt(newCategoryId as string)
  } else {
    filterForm.category_id = undefined
  }
  
  // 重新加载任务和统计数据
  await loadTasks(1, false)
  
  // 重新加载统计数据（基于当前分类）
  const statsLoaded = await loadStats()
  if (!statsLoaded) {
    await loadStatsFromAllTasks()
  }
})

// 页面卸载时移除事件监听
onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  window.removeEventListener('scroll', handleScroll)
  window.removeEventListener('open-search-dialog', () => {})
})
</script>

<style lang="scss" scoped>
.tasks-container {
  padding: 24px;

  // 分类面包屑
  .category-breadcrumb {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 20px;
    margin-bottom: 16px;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.08) 0%, rgba(118, 75, 162, 0.08) 100%);
    border-radius: 8px;
    border: 1px solid rgba(102, 126, 234, 0.2);
    font-size: 14px;
    color: #374151;

    .el-icon {
      font-size: 16px;
      color: #667eea;
    }

    span {
      font-weight: 500;
      flex: 1;
    }

    .el-button {
      font-size: 13px;
    }
  }

  // 顶部任务统计栏
  .tasks-stats-bar {
    display: flex;
    gap: 8px;
    margin-bottom: 20px;
    flex-wrap: wrap;

    .stat-btn {
      flex: 0 1 auto;
      min-width: 90px;
      height: 42px;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 6px;
      font-size: 13px;
      font-weight: 500;
      border: none;
      cursor: pointer;
      transition: all 0.3s ease;
      opacity: 0.7;
      padding: 0 12px;

      &:hover {
        opacity: 0.85;
        transform: translateY(-2px);
      }

      &.active {
        opacity: 1;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }

      .stat-value {
        font-size: 18px;
        font-weight: 600;
      }

      // 总任务 - 深蓝色
      &.stat-total {
        background-color: #3b82f6;
        color: #ffffff;
        border-color: #3b82f6;
      }

      // 待接取 - 酒红色
      &.stat-pending {
        background-color: #be123c;
        color: #ffffff;
        border-color: #be123c;
      }

      // 进行中 - 橙色
      &.stat-progress {
        background-color: #f59e0b;
        color: #ffffff;
        border-color: #f59e0b;
      }

      // 待验收 - 紫色
      &.stat-reviewing {
        background-color: #8b5cf6;
        color: #ffffff;
        border-color: #8b5cf6;
      }

      // 已完成 - 绿色
      &.stat-completed {
        background-color: #10b981;
        color: #ffffff;
        border-color: #10b981;
      }

      // 已取消 - 深灰色
      &.stat-cancelled {
        background-color: #4b5563;
        color: #ffffff;
        border-color: #4b5563;
      }
    }
  }

  // 搜索对话框内的表单
  .filter-form {
    .el-form-item {
      margin-bottom: 24px;
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

        .mobile-load-more {
          padding: 32px 16px;
          text-align: center;
          border-top: 1px solid #eee;
          margin-top: 24px;

          .task-count-info {
            font-size: 14px;
            color: #666;
            margin-bottom: 16px;
            padding: 12px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border-radius: 8px;
            font-weight: 500;
          }

          .load-more-btn {
            width: 100%;
            max-width: 400px;
            height: 48px;
            font-size: 16px;
            font-weight: 600;
            border-radius: 24px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border: none;
            
            &:hover {
              background: linear-gradient(135deg, #5568d3 0%, #6a3f8f 100%);
            }
          }

          .no-more-tip {
            font-size: 13px;
            color: #999;
            margin-top: 16px;
            padding: 12px;
            
            .el-button--text {
              color: #1890ff;
              font-size: 13px;
              padding: 0 4px;
            }
          }
        }

        .pagination-wrapper {
          display: flex;
          justify-content: center;
          align-items: center;
          padding: 32px 0;
          margin-top: 24px;
          border-top: 1px solid #eee;

          :deep(.el-pagination) {
            .btn-prev,
            .btn-next,
            .el-pager li {
              border-radius: 6px;
              margin: 0 4px;
            }

            .el-pager li.is-active {
              background-color: #1890ff;
              font-weight: 600;
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

    // 移动端统计栏
    .tasks-stats-bar {
      gap: 6px;
      flex-wrap: wrap;

      .stat-btn {
        flex: 0 0 calc((100% - 12px) / 3) !important;
        min-width: 0 !important;
        max-width: calc((100% - 12px) / 3) !important;
        height: 42px !important;
        font-size: 10px !important;
        padding: 0 4px !important;
        flex-direction: column !important;
        gap: 1px !important;
        margin: 0 !important;

        .stat-value {
          font-size: 14px;
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

// 任务列表状态标签颜色
.tasks-container {
  :deep(.task-status-pending) {
    background-color: #be123c;
    color: #ffffff;
    border-color: #be123c;
  }

  :deep(.task-status-in_progress) {
    background-color: #f59e0b;
    color: #ffffff;
    border-color: #f59e0b;
  }

  :deep(.task-status-reviewing) {
    background-color: #8b5cf6;
    color: #ffffff;
    border-color: #8b5cf6;
  }

  :deep(.task-status-completed) {
    background-color: #10b981;
    color: #ffffff;
    border-color: #10b981;
  }

  :deep(.task-status-cancelled) {
    background-color: #4b5563;
    color: #ffffff;
    border-color: #4b5563;
  }
}
</style>

<style lang="scss">
// 搜索对话框遮罩层样式（必须使用全局样式）
.search-modal-mask {
  background-color: rgba(255, 255, 255, 0.8) !important;
  backdrop-filter: blur(8px);
}

// 强制对话框使用白天模式（全局样式）
.search-dialog.el-dialog {
  background: #ffffff;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.15);
  border-radius: 12px;
  overflow: hidden;
  
  .el-dialog__header {
    background: #ffffff;
    border-bottom: 1px solid #e5e7eb;
    padding: 20px 24px;
    
    .el-dialog__title {
      color: #1a1a1a;
      font-weight: 600;
      font-size: 18px;
    }
    
    .el-dialog__headerbtn {
      .el-dialog__close {
        color: #6b7280;
        font-size: 20px;
        
        &:hover {
          color: #1a1a1a;
        }
      }
    }
  }
  
  .el-dialog__body {
    background: #ffffff;
    padding: 24px;
    
    .el-form {
      .el-form-item__label {
        color: #374151;
        font-weight: 500;
      }
      
      .el-input__wrapper {
        background-color: #f9fafb;
        border: 1px solid #e5e7eb;
        box-shadow: none;
        
        &:hover,
        &.is-focus {
          border-color: #667eea;
          background-color: #ffffff;
        }
        
        .el-input__inner {
          color: #1a1a1a;
          
          &::placeholder {
            color: #9ca3af;
          }
        }
      }
      
      .el-select {
        .el-input__wrapper {
          background-color: #f9fafb;
          border: 1px solid #e5e7eb;
        }
      }
      
      .el-input-number {
        .el-input__wrapper {
          background-color: #f9fafb;
          border: 1px solid #e5e7eb;
        }
      }
      
      .el-radio-button {
        .el-radio-button__inner {
          background-color: #f9fafb;
          border: 1px solid #e5e7eb;
          color: #374151;
          
          &:hover {
            color: #667eea;
          }
        }
        
        &.is-active {
          .el-radio-button__inner {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-color: #667eea;
            color: #ffffff;
          }
        }
      }
    }
  }
  
  .el-dialog__footer {
    background: #f9fafb;
    border-top: 1px solid #e5e7eb;
    padding: 16px 24px;
    
    .el-button {
      &.el-button--default {
        background-color: #ffffff;
        border: 1px solid #e5e7eb;
        color: #374151;
        
        &:hover {
          background-color: #f9fafb;
          border-color: #d1d5db;
          color: #1a1a1a;
        }
      }
      
      &.el-button--primary {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border: none;
        color: #ffffff;
        
        &:hover {
          background: linear-gradient(135deg, #5568d3 0%, #6a3f8f 100%);
        }
      }
    }
  }
}
</style>