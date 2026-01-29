<template>
  <div class="user-container">
    <!-- 用户信息卡片 -->
    <div class="user-profile-card">
      <div class="profile-header">
        <div class="avatar-section">
          <el-avatar
            :size="80"
            :src="userStore.avatar"
            @click="handleAvatarClick"
          >
            {{ userStore.nickname?.charAt(0) || 'U' }}
          </el-avatar>
          <el-button
            type="text"
            size="small"
            @click="handleAvatarClick"
          >
            更换头像
          </el-button>
        </div>

        <div class="user-info">
          <h2>{{ userStore.nickname || '未设置昵称' }}</h2>
          <p class="user-id">ID: {{ userStore.id }}</p>
          <div class="user-stats">
            <div class="stat-item">
              <span class="stat-value">{{ stats.completedTasks }}</span>
              <span class="stat-label">完成任务</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ stats.totalEarnings }}</span>
              <span class="stat-label">总收益</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ stats.credit }}</span>
              <span class="stat-label">信用分</span>
            </div>
          </div>
        </div>

        <div class="profile-actions">
          <el-button type="primary" @click="editProfile">
            <el-icon><Edit /></el-icon>
            编辑资料
          </el-button>
          <el-button @click="shareProfile">
            <el-icon><Share /></el-icon>
            分享主页
          </el-button>
        </div>
      </div>

      <!-- 认证状态 -->
      <div class="verification-section">
        <div class="verification-item">
          <el-icon class="verified-icon" :color="userStore.phone_verified ? '#67c23a' : '#909399'">
            <CircleCheck />
          </el-icon>
          <span>手机认证</span>
        </div>
        <div class="verification-item">
          <el-icon class="verified-icon" :color="userStore.wechat_verified ? '#67c23a' : '#909399'">
            <CircleCheck />
          </el-icon>
          <span>微信认证</span>
        </div>
        <div class="verification-item">
          <el-icon class="verified-icon" :color="userStore.alipay_verified ? '#67c23a' : '#909399'">
            <CircleCheck />
          </el-icon>
          <span>支付宝认证</span>
        </div>
      </div>
    </div>

    <!-- 功能菜单 -->
    <div class="menu-section">
      <el-row :gutter="20">
        <el-col :xs="12" :sm="6">
          <div class="menu-item" @click="goToWallet">
            <div class="menu-icon wallet">
              <el-icon><Wallet /></el-icon>
            </div>
            <div class="menu-content">
              <h4>我的钱包</h4>
              <p>余额 ¥{{ formatMoney(wallet.balance) }}</p>
            </div>
          </div>
        </el-col>

        <el-col :xs="12" :sm="6">
          <div class="menu-item" @click="goToTasks">
            <div class="menu-icon tasks">
              <el-icon><Document /></el-icon>
            </div>
            <div class="menu-content">
              <h4>我的任务</h4>
              <p>{{ taskStats.total }} 个任务</p>
            </div>
          </div>
        </el-col>

        <el-col :xs="12" :sm="6">
          <div class="menu-item" @click="goToReviews">
            <div class="menu-icon reviews">
              <el-icon><Star /></el-icon>
            </div>
            <div class="menu-content">
              <h4>我的评价</h4>
              <p>{{ reviewStats.total }} 条评价</p>
            </div>
          </div>
        </el-col>

        <el-col :xs="12" :sm="6">
          <div class="menu-item" @click="goToSettings">
            <div class="menu-icon settings">
              <el-icon><Setting /></el-icon>
            </div>
            <div class="menu-content">
              <h4>账户设置</h4>
              <p>管理账户信息</p>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 详细信息 -->
    <div class="detail-section">
      <el-tabs v-model="activeTab" class="user-tabs">
        <!-- 任务记录 -->
        <el-tab-pane label="任务记录" name="tasks">
          <div class="tab-content">
            <div class="task-filters">
              <el-radio-group v-model="taskFilter" @change="loadTasks">
                <el-radio-button label="all">全部</el-radio-button>
                <el-radio-button label="published">发布</el-radio-button>
                <el-radio-button label="taken">接取</el-radio-button>
                <el-radio-button label="completed">已完成</el-radio-button>
              </el-radio-group>
            </div>

            <div v-loading="tasksLoading" class="task-list">
              <div
                v-for="task in tasks"
                :key="task.id"
                class="task-item"
                @click="goToTaskDetail(task.id)"
              >
                <div class="task-info">
                  <h4>{{ task.title }}</h4>
                  <p class="task-desc">{{ task.description }}</p>
                  <div class="task-meta">
                    <el-tag :type="getTaskStatusType(task.status)" size="small">
                      {{ getTaskStatusText(task.status) }}
                    </el-tag>
                    <span class="task-time">
                      {{ formatTime(task.createdAt) }}
                    </span>
                  </div>
                </div>
                <div class="task-price">
                  <span class="price">¥{{ task.price }}</span>
                </div>
              </div>

              <el-empty
                v-if="tasks.length === 0"
                description="暂无任务记录"
              />
            </div>
          </div>
        </el-tab-pane>

        <!-- 收支明细 -->
        <el-tab-pane label="收支明细" name="transactions">
          <div class="tab-content">
            <div class="transaction-filters">
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                @change="loadTransactions"
              />
              <el-select v-model="transactionType" @change="loadTransactions">
                <el-option label="全部" value="all" />
                <el-option label="收入" value="income" />
                <el-option label="支出" value="expense" />
              </el-select>
            </div>

            <div v-loading="transactionsLoading" class="transaction-list">
              <div
                v-for="transaction in transactions"
                :key="transaction.id"
                class="transaction-item"
              >
                <div class="transaction-icon">
                  <el-icon :color="transaction.tradeType === 'income' ? '#67c23a' : '#f56c6c'">
                    <component :is="getTransactionIcon(transaction.tradeType)" />
                  </el-icon>
                </div>
                <div class="transaction-info">
                  <h4>{{ transaction.remark || transaction.tradeType }}</h4>
                  <p class="transaction-time">{{ formatTime(transaction.createdAt) }}</p>
                </div>
                <div class="transaction-amount">
                  <span :class="transaction.tradeType">
                    {{ transaction.tradeType === 'income' ? '+' : '-' }}¥{{ formatMoney(transaction.amount) }}
                  </span>
                </div>
              </div>

              <el-empty
                v-if="transactions.length === 0"
                description="暂无交易记录"
              />
            </div>
          </div>
        </el-tab-pane>

        <!-- 评价管理 -->
        <el-tab-pane label="评价管理" name="reviews">
          <div class="tab-content">
            <div class="review-filters">
              <el-radio-group v-model="reviewFilter" @change="loadReviews">
                <el-radio-button label="received">收到的评价</el-radio-button>
                <el-radio-button label="given">给出的评价</el-radio-button>
              </el-radio-group>
            </div>

            <div v-loading="reviewsLoading" class="review-list">
              <div
                v-for="review in reviews"
                :key="review.id"
                class="review-item"
              >
                <div class="review-header">
                  <div class="reviewer-info">
                    <el-avatar :size="40" :src="review.reviewerAvatar">
                      {{ review.reviewerName?.charAt(0) }}
                    </el-avatar>
                    <div class="reviewer-details">
                      <h4>{{ review.reviewerName }}</h4>
                      <div class="rating">
                        <el-rate
                          v-model="review.rating"
                          disabled
                          show-score
                        />
                      </div>
                    </div>
                  </div>
                  <span class="review-time">{{ formatTime(review.createTime) }}</span>
                </div>
                <div class="review-content">
                  <p>{{ review.content }}</p>
                </div>
                <div class="review-footer">
                  <span class="task-link">任务：{{ review.taskId }}</span>
                </div>
              </div>

              <el-empty
                v-if="reviews.length === 0"
                description="暂无评价"
              />
            </div>
          </div>
        </el-tab-pane>

        <!-- 账户安全 -->
        <el-tab-pane label="账户安全" name="security">
          <div class="tab-content">
            <div class="security-list">
              <div class="security-item">
                <div class="security-info">
                  <h4>登录密码</h4>
                  <p>保护账户安全的重要凭证</p>
                </div>
                <el-button @click="changePassword">修改密码</el-button>
              </div>

              <div class="security-item">
                <div class="security-info">
                  <h4>手机号</h4>
                  <p>{{ userStore.phone ? `${userStore.phone.slice(0, 3)}****${userStore.phone.slice(-4)}` : '未绑定' }}</p>
                </div>
                <el-button>{{ userStore.phone ? '更换' : '绑定' }}手机</el-button>
              </div>

              <div class="security-item">
                <div class="security-info">
                  <h4>微信</h4>
                  <p>{{ userStore.userInfo?.wechat_verified ? '已认证' : '未认证' }}</p>
                </div>
                <el-button>{{ userStore.userInfo?.wechat_verified ? '解除' : '绑定' }}微信</el-button>
              </div>

              <div class="security-item">
                <div class="security-info">
                  <h4>支付宝</h4>
                  <p>{{ userStore.userInfo?.alipay_verified ? '已认证' : '未认证' }}</p>
                </div>
                <el-button>{{ userStore.userInfo?.alipay_verified ? '解除' : '绑定' }}支付宝</el-button>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 头像上传 -->
    <input
      ref="avatarInput"
      type="file"
      accept="image/*"
      style="display: none"
      @change="handleAvatarUpload"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Edit,
  Share,
  CircleCheck,
  Wallet,
  Document,
  Star,
  Setting,
  ArrowUp,
  ArrowDown,
  Money
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { userApi, taskApi } from '@/api'
import { formatMoney, formatTime } from '@/utils/format'
import type { Task, Transaction, Review, Wallet as WalletType } from '@/types'

const router = useRouter()
const userStore = useUserStore()

// 状态变量
const activeTab = ref('tasks')
const taskFilter = ref('all')
const reviewFilter = ref('received')
const transactionType = ref('all')
const dateRange = ref<[Date, Date] | null>(null)

const tasksLoading = ref(false)
const transactionsLoading = ref(false)
const reviewsLoading = ref(false)

// 数据
const stats = reactive({
  completedTasks: 0,
  totalEarnings: 0,
  credit: 100
})

const taskStats = reactive({
  total: 0,
  published: 0,
  taken: 0,
  completed: 0
})

const reviewStats = reactive({
  total: 0,
  averageRating: 0
})

const wallet = ref<WalletType>({
  balance: 0,
  frozenAmount: 0,
  frozen_amount: 0,
  totalIncome: 0,
  total_income: 0,
  totalExpense: 0,
  createdAt: '',
  updatedAt: '',
  id: 0,
  userId: 0
})

const tasks = ref<Task[]>([])
const transactions = ref<Transaction[]>([])
const reviews = ref<Review[]>([])

// 文件上传
const avatarInput = ref<HTMLInputElement>()

// 加载用户统计
const loadStats = async () => {
  try {
    const response = await userApi.getStats()
    Object.assign(stats, response.data)
  } catch (error: any) {
    console.error('加载用户统计失败:', error)
  }
}

// 加载任务统计
const loadTaskStats = async () => {
  try {
    const response = await taskApi.getMyTaskStats()
    Object.assign(taskStats, response.data)
  } catch (error: any) {
    console.error('加载任务统计失败:', error)
  }
}

// 加载钱包信息
const loadWallet = async () => {
  try {
    const response = await userApi.getWallet()
    wallet.value = response.data
  } catch (error: any) {
    console.error('加载钱包信息失败:', error)
  }
}

// 加载任务记录
const loadTasks = async () => {
  try {
    tasksLoading.value = true
    
    const params: any = {
      page: 1,
      limit: 20
    }

    if (taskFilter.value !== 'all') {
      params.role = taskFilter.value
    }

    const response = await taskApi.getMyTasks(params)
    tasks.value = response.data.list
  } catch (error: any) {
    ElMessage.error('加载任务记录失败')
  } finally {
    tasksLoading.value = false
  }
}

// 加载交易记录
const loadTransactions = async () => {
  try {
    transactionsLoading.value = true
    
    const params: any = {
      page: 1,
      limit: 20
    }

    if (transactionType.value !== 'all') {
      params.type = transactionType.value
    }

    if (dateRange.value) {
      params.start_date = dateRange.value[0].toISOString()
      params.end_date = dateRange.value[1].toISOString()
    }

    const response = await userApi.getTransactions(params)
    transactions.value = response.data.list
  } catch (error: any) {
    ElMessage.error('加载交易记录失败')
  } finally {
    transactionsLoading.value = false
  }
}

// 加载评价记录
const loadReviews = async () => {
  try {
    reviewsLoading.value = true
    
    const response = await userApi.getReviews({
      type: reviewFilter.value,
      page: 1,
      limit: 20
    })
    reviews.value = response.data.list
  } catch (error: any) {
    ElMessage.error('加载评价记录失败')
  } finally {
    reviewsLoading.value = false
  }
}

// 页面跳转
const goToWallet = () => {
  router.push('/user/wallet')
}

const goToTasks = () => {
  router.push('/user/tasks')
}

const goToReviews = () => {
  router.push('/user/reviews')
}

const goToSettings = () => {
  router.push('/user/settings')
}

const goToTaskDetail = (id: number) => {
  router.push(`/tasks/${id}`)
}

const editProfile = () => {
  router.push('/user/profile/edit')
}

const shareProfile = async () => {
  const shareUrl = `${window.location.origin}/user/${userStore.id}`
  
  try {
    await navigator.clipboard.writeText(shareUrl)
    ElMessage.success('链接已复制到剪贴板')
  } catch (error) {
    ElMessage.info(shareUrl)
  }
}

const changePassword = () => {
  router.push('/user/security/password')
}

// 头像相关
const handleAvatarClick = () => {
  avatarInput.value?.click()
}

const handleAvatarUpload = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  try {
    await userApi.uploadAvatar(file)
    
    ElMessage.success('头像上传成功')
    
    // 重新加载用户信息
    await userStore.getUserInfo()
  } catch (error: any) {
    ElMessage.error(error.message || '头像上传失败')
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

// 获取交易图标
const getTransactionIcon = (type: string) => {
  const iconMap: Record<string, any> = {
    income: ArrowUp,
    expense: ArrowDown,
    reward: Money
  }
  return iconMap[type] || ArrowUp
}

// 页面初始化
onMounted(() => {
  loadStats()
  loadTaskStats()
  loadWallet()
  loadTasks()
})
</script>

<style lang="scss" scoped>
.user-container {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;

  .user-profile-card {
    background: white;
    border-radius: 16px;
    padding: 32px;
    margin-bottom: 24px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

    .profile-header {
      display: flex;
      align-items: center;
      gap: 24px;
      margin-bottom: 24px;

      .avatar-section {
        text-align: center;

        .el-avatar {
          cursor: pointer;
          margin-bottom: 8px;
        }
      }

      .user-info {
        flex: 1;

        h2 {
          margin: 0 0 8px 0;
          font-size: 24px;
          font-weight: 600;
          color: #1a1a1a;
        }

        .user-id {
          margin: 0 0 16px 0;
          font-size: 14px;
          color: #999;
        }

        .user-stats {
          display: flex;
          gap: 32px;

          .stat-item {
            .stat-value {
              display: block;
              font-size: 20px;
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
      }

      .profile-actions {
        display: flex;
        flex-direction: column;
        gap: 8px;

        .el-button {
          .el-icon {
            margin-right: 8px;
          }
        }
      }
    }

    .verification-section {
      display: flex;
      gap: 32px;
      padding-top: 24px;
      border-top: 1px solid #f0f0f0;

      .verification-item {
        display: flex;
        align-items: center;
        gap: 8px;

        .verified-icon {
          font-size: 16px;
        }

        span {
          font-size: 14px;
          color: #666;
        }
      }
    }
  }

  .menu-section {
    margin-bottom: 24px;

    .menu-item {
      background: white;
      border-radius: 12px;
      padding: 24px;
      display: flex;
      align-items: center;
      gap: 16px;
      cursor: pointer;
      transition: all 0.3s ease;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
      }

      .menu-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;

        .el-icon {
          font-size: 24px;
          color: white;
        }

        &.wallet {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }

        &.tasks {
          background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        }

        &.reviews {
          background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
        }

        &.settings {
          background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
        }
      }

      .menu-content {
        h4 {
          margin: 0 0 4px 0;
          font-size: 16px;
          font-weight: 600;
          color: #1a1a1a;
        }

        p {
          margin: 0;
          font-size: 14px;
          color: #666;
        }
      }
    }
  }

  .detail-section {
    background: white;
    border-radius: 16px;
    padding: 24px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

    .user-tabs {
      .tab-content {
        margin-top: 20px;

        .task-filters,
        .transaction-filters,
        .review-filters {
          margin-bottom: 20px;
        }

        .task-list {
          .task-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 16px;
            border: 1px solid #f0f0f0;
            border-radius: 8px;
            margin-bottom: 12px;
            cursor: pointer;
            transition: all 0.3s ease;

            &:hover {
              border-color: #1890ff;
              box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
            }

            .task-info {
              flex: 1;

              h4 {
                margin: 0 0 8px 0;
                font-size: 16px;
                font-weight: 500;
                color: #1a1a1a;
              }

              .task-desc {
                margin: 0 0 8px 0;
                font-size: 14px;
                color: #666;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
              }

              .task-meta {
                display: flex;
                align-items: center;
                gap: 12px;

                .task-type,
                .task-time {
                  font-size: 12px;
                  color: #999;
                }
              }
            }

            .task-price {
              .price {
                font-size: 18px;
                font-weight: 600;
                color: #ff4d4f;
              }
            }
          }
        }

        .transaction-list {
          .transaction-item {
            display: flex;
            align-items: center;
            padding: 16px;
            border-bottom: 1px solid #f0f0f0;

            &:last-child {
              border-bottom: none;
            }

            .transaction-icon {
              width: 40px;
              height: 40px;
              border-radius: 8px;
              background: #f5f5f5;
              display: flex;
              align-items: center;
              justify-content: center;
              margin-right: 16px;
            }

            .transaction-info {
              flex: 1;

              h4 {
                margin: 0 0 4px 0;
                font-size: 14px;
                font-weight: 500;
                color: #1a1a1a;
              }

              .transaction-time {
                font-size: 12px;
                color: #999;
              }
            }

            .transaction-amount {
              span {
                font-size: 16px;
                font-weight: 600;

                &.income {
                  color: #67c23a;
                }

                &.expense {
                  color: #f56c6c;
                }
              }
            }
          }
        }

        .review-list {
          .review-item {
            padding: 20px;
            border: 1px solid #f0f0f0;
            border-radius: 8px;
            margin-bottom: 16px;

            .review-header {
              display: flex;
              justify-content: space-between;
              align-items: center;
              margin-bottom: 12px;

              .reviewer-info {
                display: flex;
                align-items: center;
                gap: 12px;

                .reviewer-details {
                  h4 {
                    margin: 0 0 4px 0;
                    font-size: 14px;
                    font-weight: 500;
                    color: #1a1a1a;
                  }
                }
              }

              .review-time {
                font-size: 12px;
                color: #999;
              }
            }

            .review-content {
              margin-bottom: 12px;

              p {
                margin: 0;
                font-size: 14px;
                line-height: 1.6;
                color: #666;
              }
            }

            .review-footer {
              .task-link {
                font-size: 12px;
                color: #1890ff;
                cursor: pointer;

                &:hover {
                  text-decoration: underline;
                }
              }
            }
          }
        }

        .security-list {
          .security-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 20px;
            border: 1px solid #f0f0f0;
            border-radius: 8px;
            margin-bottom: 16px;

            .security-info {
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
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .user-container {
    padding: 16px;

    .user-profile-card {
      padding: 20px;

      .profile-header {
        flex-direction: column;
        text-align: center;
        gap: 16px;

        .user-info {
          .user-stats {
            justify-content: center;
            gap: 24px;
          }
        }

        .profile-actions {
          flex-direction: row;
          width: 100%;
          justify-content: center;
        }
      }

      .verification-section {
        justify-content: center;
        gap: 24px;
      }
    }

    .menu-section {
      .menu-item {
        padding: 16px;

        .menu-content {
          h4 {
            font-size: 14px;
          }

          p {
            font-size: 12px;
          }
        }
      }
    }

    .detail-section {
      padding: 16px;
    }
  }
}
</style>