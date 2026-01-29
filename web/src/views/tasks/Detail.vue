<template>
  <div class="task-detail-container">
    <div v-loading="loading" class="task-detail">
      <!-- 任务头部 -->
      <div class="task-header">
        <div class="task-badges">
          <el-tag v-if="task.is_urgent" type="danger">加急</el-tag>
          <el-tag v-if="task.is_remote" type="success">远程</el-tag>
          <el-tag :type="getTaskStatusType(task.status)">
            {{ getTaskStatusText(task.status) }}
          </el-tag>
        </div>
        
        <h1 class="task-title">{{ task.title }}</h1>
        
        <div class="task-meta">
          <div class="meta-item">
            <el-icon><Money /></el-icon>
            <span class="price">¥{{ task.price }}</span>
            <span class="price-label">赏金</span>
          </div>
          
          <div class="meta-item">
            <el-icon><Clock /></el-icon>
            <span>{{ formatTime(task.deadline) }}</span>
            <span class="meta-label">截止时间</span>
          </div>
          
          <div class="meta-item">
            <el-icon><Location /></el-icon>
            <span>{{ task.location || '不限地区' }}</span>
            <span class="meta-label">地点</span>
          </div>
          
          <div class="meta-item">
            <el-icon><User /></el-icon>
            <span>{{ task.category_name }}</span>
            <span class="meta-label">分类</span>
          </div>
        </div>
      </div>

      <!-- 发布者信息 -->
      <div class="publisher-section">
        <div class="publisher-header">
          <h3>发布者信息</h3>
          <el-button
            v-if="task.publisher_id !== userStore.id"
            size="small"
            @click="followPublisher"
          >
            {{ isFollowing ? '已关注' : '关注' }}
          </el-button>
        </div>
        
        <div class="publisher-info">
          <el-avatar :size="60" :src="task.publisher_avatar">
            {{ task.publisher_name?.charAt(0) }}
          </el-avatar>
          
          <div class="publisher-details">
            <h4>{{ task.publisher_name }}</h4>
            <p>发布 {{ task.publisher_published_count }} 个任务</p>
            <div class="publisher-stats">
              <span>信用分: {{ task.publisher_credit }}</span>
              <span>好评率: {{ task.publisher_rating }}%</span>
            </div>
          </div>
          
          <el-button
            v-if="task.publisher_id !== userStore.id"
            type="primary"
            @click="contactPublisher"
          >
            <el-icon><ChatDotRound /></el-icon>
            私信
          </el-button>
        </div>
      </div>

      <!-- 任务详情 -->
      <div class="task-content">
        <div class="content-section">
          <h3>任务详情</h3>
          <div class="task-description" v-html="task.description"></div>
          
          <!-- 任务图片 -->
          <div v-if="task.images && task.images.length > 0" class="task-images">
            <el-image
              v-for="(image, index) in task.images"
              :key="index"
              :src="image"
              :preview-src-list="task.images"
              fit="cover"
              class="task-image"
            />
          </div>
          
          <!-- 任务要求 -->
          <div v-if="task.requirements" class="task-requirements">
            <h4>任务要求</h4>
            <div v-html="task.requirements"></div>
          </div>
        </div>

        <!-- 任务统计 -->
        <div class="content-section">
          <h3>任务统计</h3>
          <div class="stats-grid">
            <div class="stat-item">
              <span class="stat-value">{{ task.view_count }}</span>
              <span class="stat-label">浏览量</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ task.apply_count }}</span>
              <span class="stat-label">申请人数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ task.favorite_count }}</span>
              <span class="stat-label">收藏数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ task.share_count }}</span>
              <span class="stat-label">分享数</span>
            </div>
          </div>
        </div>

        <!-- 申请列表 -->
        <div v-if="task.publisher_id === userStore.id && task.applications" class="content-section">
          <h3>申请列表</h3>
          <div class="application-list">
            <div
              v-for="application in task.applications"
              :key="application.id"
              class="application-item"
            >
              <div class="applicant-info">
                <el-avatar :size="40" :src="application.applicant_avatar">
                  {{ application.applicant_name?.charAt(0) }}
                </el-avatar>
                <div class="applicant-details">
                  <h4>{{ application.applicant_name }}</h4>
                  <p>信用分: {{ application.applicant_credit }}</p>
                  <p>完成任务: {{ application.applicant_completed_tasks }}</p>
                </div>
              </div>
              
              <div class="application-content">
                <p>{{ application.message }}</p>
                <p class="application-time">
                  预计完成时间: {{ application.estimated_time }}
                </p>
              </div>
              
              <div class="application-actions">
                <el-button
                  v-if="application.status === 'pending'"
                  type="primary"
                  size="small"
                  @click="acceptApplication(application.id)"
                >
                  接受申请
                </el-button>
                <el-button
                  v-if="application.status === 'pending'"
                  type="danger"
                  size="small"
                  @click="rejectApplication(application.id)"
                >
                  拒绝申请
                </el-button>
                <el-tag
                  v-if="application.status !== 'pending'"
                  :type="getApplicationStatusType(application.status)"
                >
                  {{ getApplicationStatusText(application.status) }}
                </el-tag>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 侧边栏 -->
    <div class="task-sidebar">
      <!-- 操作按钮 -->
      <div class="action-card">
        <div v-if="task.status === 'pending' && task.publisher_id !== userStore.id" class="action-buttons">
          <el-button
            type="primary"
            size="large"
            :loading="applyLoading"
            @click="showApplyDialog"
          >
            申请接取
          </el-button>
          
          <el-button
            :type="isFavorited ? 'danger' : 'default'"
            size="large"
            @click="toggleFavorite"
          >
            <el-icon><Star /></el-icon>
            {{ isFavorited ? '已收藏' : '收藏' }}
          </el-button>
        </div>
        
        <div v-if="task.publisher_id === userStore.id" class="owner-actions">
          <el-button
            v-if="task.status === 'pending'"
            type="warning"
            size="large"
            @click="editTask"
          >
            <el-icon><Edit /></el-icon>
            编辑任务
          </el-button>
          
          <el-button
            v-if="task.status === 'in_progress'"
            type="success"
            size="large"
            @click="confirmDelivery"
          >
            <el-icon><CircleCheck /></el-icon>
            确认完成
          </el-button>
          
          <el-button
            v-if="['pending', 'in_progress'].includes(task.status)"
            type="danger"
            size="large"
            @click="cancelTask"
          >
            <el-icon><Close /></el-icon>
            取消任务
          </el-button>
        </div>

        <!-- 分享按钮 -->
        <div class="share-section">
          <p>分享任务</p>
          <div class="share-buttons">
            <el-button size="small" @click="shareToWechat">
              <el-icon><WechatFilled /></el-icon>
              微信
            </el-button>
            <el-button size="small" @click="shareToAlipay">
              <el-icon><CreditCard /></el-icon>
              支付宝
            </el-button>
            <el-button size="small" @click="copyLink">
              <el-icon><CopyDocument /></el-icon>
              复制链接
            </el-button>
          </div>
        </div>
      </div>

      <!-- 安全提醒 -->
      <div class="safety-card">
        <h4>
          <el-icon><Warning /></el-icon>
          安全提醒
        </h4>
        <ul>
          <li>请在线下见面前先确认对方身份</li>
          <li>不要轻易泄露个人隐私信息</li>
          <li>建议选择公共场所进行交易</li>
          <li>如遇到问题请及时联系平台客服</li>
        </ul>
      </div>

      <!-- 相关任务 -->
      <div class="related-tasks-card">
        <h4>相关任务</h4>
        <div v-loading="relatedLoading" class="related-tasks">
          <div
            v-for="relatedTask in relatedTasks"
            :key="relatedTask.id"
            class="related-task-item"
            @click="goToTask(relatedTask.id)"
          >
            <h5>{{ relatedTask.title }}</h5>
            <div class="related-task-meta">
              <span class="price">¥{{ relatedTask.price }}</span>
              <span>{{ formatTime(relatedTask.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 申请弹窗 -->
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
          @click="submitApply"
        >
          提交申请
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Money,
  Clock,
  Location,
  User,
  ChatDotRound,
  Star,
  Edit,
  CircleCheck,
  Close,
  Warning,
  WechatFilled,
  CreditCard,
  CopyDocument
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { taskApi } from '@/api'
import { formatTime } from '@/utils/format'
import type { Task } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 状态变量
const loading = ref(false)
const applyLoading = ref(false)
const relatedLoading = ref(false)
const applyDialogVisible = ref(false)
const isFavorited = ref(false)
const isFollowing = ref(false)

// 数据
const task = ref<Task>({} as Task)
const relatedTasks = ref<Task[]>([])

// 申请表单
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

// 加载任务详情
const loadTaskDetail = async () => {
  try {
    loading.value = true
    const taskId = Number(route.params.id)
    const response = await taskApi.getTaskDetail(taskId)
    task.value = response.data
  } catch (error: any) {
    ElMessage.error('加载任务详情失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 加载相关任务
const loadRelatedTasks = async () => {
  try {
    relatedLoading.value = true
    const response = await taskApi.getRelatedTasks(task.value.id, {
      limit: 5
    })
    relatedTasks.value = response.data
  } catch (error: any) {
    console.error('加载相关任务失败:', error)
  } finally {
    relatedLoading.value = false
  }
}

// 获取任务状态类型
const getTaskStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: 'info',
    pending_payment: 'warning',
    in_progress: 'primary',
    delivered: 'success',
    completed: 'success',
    cancelled: 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取任务状态文本
const getTaskStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待接取',
    pending_payment: '待支付',
    in_progress: '进行中',
    delivered: '已交付',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || '未知'
}

// 获取申请状态类型
const getApplicationStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: 'warning',
    accepted: 'success',
    rejected: 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取申请状态文本
const getApplicationStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待处理',
    accepted: '已接受',
    rejected: '已拒绝'
  }
  return statusMap[status] || '未知'
}

// 显示申请弹窗
const showApplyDialog = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  applyDialogVisible.value = true
}

// 提交申请
const submitApply = async () => {
  if (!applyFormRef.value) return

  try {
    await applyFormRef.value.validate()
    
    applyLoading.value = true
    
    await taskApi.applyTask(task.value.id, applyForm)
    
    ElMessage.success('申请已提交，请等待任务发布者确认')
    applyDialogVisible.value = false
    
    // 重新加载任务详情
    await loadTaskDetail()
    
  } catch (error: any) {
    if (error.message) {
      ElMessage.error(error.message)
    }
  } finally {
    applyLoading.value = false
  }
}

// 收藏任务
const toggleFavorite = async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }

  try {
    if (isFavorited.value) {
      await taskApi.unfavoriteTask(task.value.id)
      ElMessage.success('已取消收藏')
    } else {
      await taskApi.favoriteTask(task.value.id)
      ElMessage.success('收藏成功')
    }
    isFavorited.value = !isFavorited.value
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

// 关注发布者
const followPublisher = async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }

  try {
    if (isFollowing.value) {
      await taskApi.unfollowUser(task.value.publisher_id)
      ElMessage.success('已取消关注')
    } else {
      await taskApi.followUser(task.value.publisher_id)
      ElMessage.success('关注成功')
    }
    isFollowing.value = !isFollowing.value
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

// 联系发布者
const contactPublisher = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  // 跳转到聊天页面
  router.push(`/chat/${task.value.publisher_id}`)
}

// 编辑任务
const editTask = () => {
  router.push(`/tasks/${task.value.id}/edit`)
}

// 确认完成
const confirmDelivery = async () => {
  try {
    await ElMessageBox.confirm('确认该任务已完成？', '确认完成', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await taskApi.confirmDelivery(task.value.id)
    ElMessage.success('任务已标记为完成')
    
    // 重新加载任务详情
    await loadTaskDetail()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '操作失败')
    }
  }
}

// 取消任务
const cancelTask = async () => {
  try {
    await ElMessageBox.confirm('确定要取消该任务？', '取消任务', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await taskApi.cancelTask(task.value.id)
    ElMessage.success('任务已取消')
    
    // 重新加载任务详情
    await loadTaskDetail()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '操作失败')
    }
  }
}

// 接受申请
const acceptApplication = async (applicationId: number) => {
  try {
    await taskApi.acceptApplication(applicationId)
    ElMessage.success('申请已接受')
    
    // 重新加载任务详情
    await loadTaskDetail()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

// 拒绝申请
const rejectApplication = async (applicationId: number) => {
  try {
    await taskApi.rejectApplication(applicationId)
    ElMessage.success('申请已拒绝')
    
    // 重新加载任务详情
    await loadTaskDetail()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

// 分享功能
const shareToWechat = () => {
  // 微信分享功能
  ElMessage.info('请在微信中打开进行分享')
}

const shareToAlipay = () => {
  // 支付宝分享功能
  ElMessage.info('请在支付宝中打开进行分享')
}

const copyLink = async () => {
  try {
    await navigator.clipboard.writeText(window.location.href)
    ElMessage.success('链接已复制到剪贴板')
  } catch (error) {
    ElMessage.info('请手动复制链接')
  }
}

// 跳转到任务
const goToTask = (id: number) => {
  router.push(`/tasks/${id}`)
}

// 页面初始化
onMounted(async () => {
  await loadTaskDetail()
  await loadRelatedTasks()
})
</script>

<style lang="scss" scoped>
.task-detail-container {
  display: flex;
  gap: 24px;
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;

  .task-detail {
    flex: 1;
    min-width: 0;

    .task-header {
      background: white;
      border-radius: 12px;
      padding: 24px;
      margin-bottom: 20px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

      .task-badges {
        display: flex;
        gap: 8px;
        margin-bottom: 16px;
        flex-wrap: wrap;
      }

      .task-title {
        margin: 0 0 20px 0;
        font-size: 28px;
        font-weight: 600;
        line-height: 1.4;
        color: #1a1a1a;
      }

      .task-meta {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 16px;

        .meta-item {
          display: flex;
          align-items: center;
          gap: 8px;

          .el-icon {
            font-size: 16px;
            color: #666;
          }

          .price {
            font-size: 20px;
            font-weight: 600;
            color: #ff4d4f;
          }

          .price-label,
          .meta-label {
            font-size: 12px;
            color: #999;
          }
        }
      }
    }

    .publisher-section {
      background: white;
      border-radius: 12px;
      padding: 24px;
      margin-bottom: 20px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

      .publisher-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;

        h3 {
          margin: 0;
          font-size: 18px;
          font-weight: 600;
          color: #1a1a1a;
        }
      }

      .publisher-info {
        display: flex;
        align-items: center;
        gap: 16px;

        .publisher-details {
          flex: 1;

          h4 {
            margin: 0 0 4px 0;
            font-size: 16px;
            font-weight: 500;
            color: #1a1a1a;
          }

          p {
            margin: 0 0 4px 0;
            font-size: 14px;
            color: #666;
          }

          .publisher-stats {
            display: flex;
            gap: 16px;
            font-size: 12px;
            color: #999;
          }
        }
      }
    }

    .task-content {
      .content-section {
        background: white;
        border-radius: 12px;
        padding: 24px;
        margin-bottom: 20px;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

        h3 {
          margin: 0 0 16px 0;
          font-size: 18px;
          font-weight: 600;
          color: #1a1a1a;
        }

        .task-description {
          font-size: 16px;
          line-height: 1.8;
          color: #333;
          margin-bottom: 20px;

          :deep(img) {
            max-width: 100%;
            border-radius: 8px;
          }
        }

        .task-images {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
          gap: 12px;
          margin-bottom: 20px;

          .task-image {
            width: 100%;
            height: 200px;
            border-radius: 8px;
            cursor: pointer;
          }
        }

        .task-requirements {
          h4 {
            margin: 0 0 12px 0;
            font-size: 16px;
            font-weight: 500;
            color: #1a1a1a;
          }

          div {
            font-size: 14px;
            line-height: 1.6;
            color: #666;
          }
        }

        .stats-grid {
          display: grid;
          grid-template-columns: repeat(4, 1fr);
          gap: 16px;

          .stat-item {
            text-align: center;
            padding: 16px;
            background: #fafafa;
            border-radius: 8px;

            .stat-value {
              display: block;
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

        .application-list {
          .application-item {
            border: 1px solid #f0f0f0;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 16px;

            .applicant-info {
              display: flex;
              align-items: center;
              gap: 12px;
              margin-bottom: 16px;

              .applicant-details {
                flex: 1;

                h4 {
                  margin: 0 0 4px 0;
                  font-size: 16px;
                  font-weight: 500;
                  color: #1a1a1a;
                }

                p {
                  margin: 0;
                  font-size: 14px;
                  color: #666;
                }
              }
            }

            .application-content {
              margin-bottom: 16px;

              p {
                margin: 0 0 4px 0;
                font-size: 14px;
                color: #333;
              }

              .application-time {
                font-size: 12px;
                color: #999;
              }
            }

            .application-actions {
              text-align: right;
            }
          }
        }
      }
    }
  }

  .task-sidebar {
    width: 360px;
    flex-shrink: 0;

    .action-card {
      background: white;
      border-radius: 12px;
      padding: 24px;
      margin-bottom: 20px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

      .action-buttons,
      .owner-actions {
        display: flex;
        flex-direction: column;
        gap: 12px;
        margin-bottom: 24px;

        .el-button {
          height: 44px;
          font-size: 16px;

          .el-icon {
            margin-right: 8px;
          }
        }
      }

      .share-section {
        p {
          margin: 0 0 12px 0;
          font-size: 14px;
          color: #666;
          text-align: center;
        }

        .share-buttons {
          display: flex;
          gap: 8px;

          .el-button {
            flex: 1;
          }
        }
      }
    }

    .safety-card {
      background: white;
      border-radius: 12px;
      padding: 24px;
      margin-bottom: 20px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

      h4 {
        display: flex;
        align-items: center;
        gap: 8px;
        margin: 0 0 16px 0;
        font-size: 16px;
        font-weight: 600;
        color: #1a1a1a;

        .el-icon {
          color: #e6a23c;
        }
      }

      ul {
        margin: 0;
        padding-left: 20px;

        li {
          font-size: 14px;
          color: #666;
          line-height: 1.6;
          margin-bottom: 8px;
        }
      }
    }

    .related-tasks-card {
      background: white;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

      h4 {
        margin: 0 0 16px 0;
        font-size: 16px;
        font-weight: 600;
        color: #1a1a1a;
      }

      .related-tasks {
        .related-task-item {
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

          h5 {
            margin: 0 0 8px 0;
            font-size: 14px;
            font-weight: 500;
            color: #1a1a1a;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .related-task-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 12px;
            color: #999;

            .price {
              font-size: 14px;
              font-weight: 600;
              color: #ff4d4f;
            }
          }
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 1024px) {
  .task-detail-container {
    flex-direction: column;

    .task-sidebar {
      width: 100%;
    }
  }
}

@media (max-width: 768px) {
  .task-detail-container {
    padding: 16px;

    .task-detail {
      .task-header {
        padding: 16px;

        .task-title {
          font-size: 24px;
        }

        .task-meta {
          grid-template-columns: 1fr;
        }
      }

      .publisher-section,
      .task-content .content-section {
        padding: 16px;

        .stats-grid {
          grid-template-columns: repeat(2, 1fr);
        }
      }
    }

    .task-sidebar {
      .action-card,
      .safety-card,
      .related-tasks-card {
        padding: 16px;
      }
    }
  }
}
</style>