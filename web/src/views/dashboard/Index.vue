<template>
  <div class="dashboard-container">
    <div class="dashboard-header">
      <div class="welcome-section">
        <h1>欢迎回来，{{ userStore.nickname || '用户' }}！</h1>
        <p>今天是美好的一天，让我们开始新的任务吧</p>
      </div>
      
      <div class="quick-actions">
        <el-button type="primary" size="large" @click="publishTask">
          <el-icon><Plus /></el-icon>
          发布任务
        </el-button>
        <el-button size="large" @click="browseTasks">
          <el-icon><Search /></el-icon>
          浏览任务
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="6">
          <div class="stats-card">
            <div class="stats-icon published">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stats-content">
              <h3>{{ stats.publishedTasks }}</h3>
              <p>发布任务</p>
            </div>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="6">
          <div class="stats-card">
            <div class="stats-icon taken">
              <el-icon><TakeawayBox /></el-icon>
            </div>
            <div class="stats-content">
              <h3>{{ stats.takenTasks }}</h3>
              <p>接取任务</p>
            </div>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="6">
          <div class="stats-card">
            <div class="stats-icon completed">
              <el-icon><CircleCheck /></el-icon>
            </div>
            <div class="stats-content">
              <h3>{{ stats.completedTasks }}</h3>
              <p>完成任务</p>
            </div>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="6">
          <div class="stats-card">
            <div class="stats-icon earnings">
              <el-icon><Wallet /></el-icon>
            </div>
            <div class="stats-content">
              <h3>¥{{ formatMoney(stats.totalEarnings) }}</h3>
              <p>总收益</p>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <div class="dashboard-content">
      <el-row :gutter="20">
        <!-- 我的任务 -->
        <el-col :xs="24" :lg="16">
          <div class="content-card">
            <div class="card-header">
              <h3>我的任务</h3>
              <el-tabs v-model="activeTaskTab" class="task-tabs">
                <el-tab-pane label="进行中" name="ongoing">
                  <div v-loading="tasksLoading" class="task-list">
                    <div
                      v-for="task in ongoingTasks"
                      :key="task.id"
                      class="task-item"
                      @click="goToTaskDetail(task.id)"
                    >
                      <div class="task-info">
                        <h4>{{ task.title }}</h4>
                        <p class="task-desc">{{ task.description }}</p>
                        <div class="task-meta">
                          <el-tag :type="getTaskStatusType(task.status)">
                            {{ getTaskStatusText(task.status) }}
                          </el-tag>
                          <span class="task-time">
                            {{ formatTime(task.created_at) }}
                          </span>
                        </div>
                      </div>
                      <div class="task-price">
                        <span class="price">¥{{ task.price }}</span>
                        <el-button
                          v-if="task.status === 'pending_payment'"
                          type="primary"
                          size="small"
                          @click.stop="handlePay(task)"
                        >
                          立即支付
                        </el-button>
                      </div>
                    </div>
                    
                    <el-empty
                      v-if="ongoingTasks.length === 0"
                      description="暂无进行中的任务"
                    />
                  </div>
                </el-tab-pane>
                
                <el-tab-pane label="已完成" name="completed">
                  <div v-loading="tasksLoading" class="task-list">
                    <div
                      v-for="task in completedTasks"
                      :key="task.id"
                      class="task-item"
                      @click="goToTaskDetail(task.id)"
                    >
                      <div class="task-info">
                        <h4>{{ task.title }}</h4>
                        <p class="task-desc">{{ task.description }}</p>
                        <div class="task-meta">
                          <el-tag type="success">已完成</el-tag>
                          <span class="task-time">
                            {{ formatTime(task.completed_at) }}
                          </span>
                        </div>
                      </div>
                      <div class="task-price">
                        <span class="price">¥{{ task.price }}</span>
                      </div>
                    </div>
                    
                    <el-empty
                      v-if="completedTasks.length === 0"
                      description="暂无已完成的任务"
                    />
                  </div>
                </el-tab-pane>
              </el-tabs>
            </div>
          </div>
        </el-col>

        <!-- 右侧信息 -->
        <el-col :xs="24" :lg="8">
          <!-- 余额信息 -->
          <div class="content-card">
            <div class="card-header">
              <h3>账户余额</h3>
              <el-link type="primary" @click="goToWallet">详情</el-link>
            </div>
            <div class="wallet-info">
              <div class="balance-amount">
                <span class="currency">¥</span>
                <span class="amount">{{ formatMoney(wallet.balance) }}</span>
              </div>
              <div class="wallet-actions">
                <el-button type="primary" size="small" @click="recharge">
                  充值
                </el-button>
                <el-button size="small" @click="withdraw">
                  提现
                </el-button>
              </div>
            </div>
          </div>

          <!-- 推荐任务 -->
          <div class="content-card">
            <div class="card-header">
              <h3>推荐任务</h3>
              <el-link type="primary" @click="browseTasks">更多</el-link>
            </div>
            <div v-loading="recommendLoading" class="recommend-list">
              <div
                v-for="task in recommendTasks"
                :key="task.id"
                class="recommend-item"
                @click="goToTaskDetail(task.id)"
              >
                <h4>{{ task.title }}</h4>
                <div class="recommend-meta">
                  <span class="price">¥{{ task.price }}</span>
                  <span class="category">{{ task.category_name }}</span>
                </div>
              </div>
              
              <el-empty
                v-if="recommendTasks.length === 0"
                description="暂无推荐任务"
                :image-size="100"
              />
            </div>
          </div>

          <!-- 系统公告 -->
          <div class="content-card">
            <div class="card-header">
              <h3>系统公告</h3>
              <el-link type="primary" @click="goToAnnouncements">全部</el-link>
            </div>
            <div v-loading="announcementsLoading" class="announcement-list">
              <div
                v-for="announcement in announcements"
                :key="announcement.id"
                class="announcement-item"
                @click="goToAnnouncement(announcement.id)"
              >
                <div class="announcement-title">
                  <el-icon v-if="!announcement.read" class="unread-dot"><Circle /></el-icon>
                  {{ announcement.title }}
                </div>
                <div class="announcement-time">
                  {{ formatTime(announcement.created_at) }}
                </div>
              </div>
              
              <el-empty
                v-if="announcements.length === 0"
                description="暂无公告"
                :image-size="100"
              />
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Plus,
  Search,
  Document,
  TakeawayBox,
  CircleCheck,
  Wallet,
  Circle
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { taskApi, paymentApi, userApi } from '@/api'
import { formatMoney, formatTime } from '@/utils/format'
import type { Task, Wallet as WalletType, Announcement } from '@/types'

const router = useRouter()
const userStore = useUserStore()

// 状态变量
const activeTaskTab = ref('ongoing')
const tasksLoading = ref(false)
const recommendLoading = ref(false)
const announcementsLoading = ref(false)

// 数据
const stats = reactive({
  publishedTasks: 0,
  takenTasks: 0,
  completedTasks: 0,
  totalEarnings: 0
})

const ongoingTasks = ref<Task[]>([])
const completedTasks = ref<Task[]>([])
const recommendTasks = ref<Task[]>([])
const announcements = ref<Announcement[]>([])
const wallet = ref<WalletType>({
  balance: 0,
  frozen_amount: 0,
  total_income: 0,
  total_expense: 0
})

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await userApi.getStats()
    Object.assign(stats, response.data)
  } catch (error: any) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载我的任务
const loadMyTasks = async () => {
  try {
    tasksLoading.value = true
    
    const [ongoingResponse, completedResponse] = await Promise.all([
      taskApi.getMyTasks({ status: 'ongoing', page: 1, limit: 5 }),
      taskApi.getMyTasks({ status: 'completed', page: 1, limit: 5 })
    ])
    
    ongoingTasks.value = ongoingResponse.data.list
    completedTasks.value = completedResponse.data.list
  } catch (error: any) {
    ElMessage.error('加载任务失败')
  } finally {
    tasksLoading.value = false
  }
}

// 加载推荐任务
const loadRecommendTasks = async () => {
  try {
    recommendLoading.value = true
    const response = await taskApi.getTasks({
      recommend: true,
      page: 1,
      limit: 5
    })
    recommendTasks.value = response.data.list
  } catch (error: any) {
    console.error('加载推荐任务失败:', error)
  } finally {
    recommendLoading.value = false
  }
}

// 加载系统公告
const loadAnnouncements = async () => {
  try {
    announcementsLoading.value = true
    const response = await userApi.getAnnouncements({ page: 1, limit: 5 })
    announcements.value = response.data.list
  } catch (error: any) {
    console.error('加载公告失败:', error)
  } finally {
    announcementsLoading.value = false
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

// 页面跳转
const publishTask = () => {
  router.push('/tasks/publish')
}

const browseTasks = () => {
  router.push('/tasks')
}

const goToTaskDetail = (id: number) => {
  router.push(`/tasks/${id}`)
}

const goToWallet = () => {
  router.push('/user/wallet')
}

const goToAnnouncements = () => {
  router.push('/announcements')
}

const goToAnnouncement = (id: number) => {
  router.push(`/announcements/${id}`)
}

// 支付处理
const handlePay = async (task: Task) => {
  try {
    const response = await paymentApi.prePay({
      order_type: 'task',
      order_id: task.id,
      amount: task.price,
      payment_method: 'alipay'
    })
    
    // 跳转到支付页面
    router.push({
      path: '/payment',
      query: {
        order_no: response.data.order_no,
        amount: task.price
      }
    })
  } catch (error: any) {
    ElMessage.error(error.message || '发起支付失败')
  }
}

// 充值
const recharge = () => {
  router.push('/user/wallet/recharge')
}

// 提现
const withdraw = () => {
  router.push('/user/wallet/withdraw')
}

// 页面初始化
onMounted(() => {
  loadStats()
  loadMyTasks()
  loadRecommendTasks()
  loadAnnouncements()
  loadWallet()
})
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 24px;

  .dashboard-header {
    margin-bottom: 32px;

    .welcome-section {
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

    .quick-actions {
      margin-top: 24px;
      display: flex;
      gap: 16px;

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

  .stats-cards {
    margin-bottom: 32px;

    .stats-card {
      background: white;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
      display: flex;
      align-items: center;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
      }

      .stats-icon {
        width: 56px;
        height: 56px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 16px;

        .el-icon {
          font-size: 24px;
          color: white;
        }

        &.published {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }

        &.taken {
          background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        }

        &.completed {
          background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
        }

        &.earnings {
          background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
        }
      }

      .stats-content {
        h3 {
          margin: 0 0 4px 0;
          font-size: 28px;
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

  .dashboard-content {
    .content-card {
      background: white;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
      margin-bottom: 20px;

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 20px;

        h3 {
          margin: 0;
          font-size: 18px;
          font-weight: 600;
          color: #1a1a1a;
        }
      }

      .task-tabs {
        .task-list {
          .task-item {
            padding: 16px;
            border: 1px solid #f0f0f0;
            border-radius: 8px;
            margin-bottom: 12px;
            cursor: pointer;
            transition: all 0.3s ease;
            display: flex;
            justify-content: space-between;
            align-items: center;

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

                .task-time {
                  font-size: 12px;
                  color: #999;
                }
              }
            }

            .task-price {
              text-align: right;

              .price {
                display: block;
                font-size: 18px;
                font-weight: 600;
                color: #ff4d4f;
                margin-bottom: 8px;
              }
            }
          }
        }
      }

      .wallet-info {
        text-align: center;

        .balance-amount {
          margin-bottom: 24px;

          .currency {
            font-size: 24px;
            color: #1890ff;
          }

          .amount {
            font-size: 36px;
            font-weight: 600;
            color: #1890ff;
          }
        }

        .wallet-actions {
          display: flex;
          gap: 12px;
          justify-content: center;

          .el-button {
            flex: 1;
          }
        }
      }

      .recommend-list {
        .recommend-item {
          padding: 12px 0;
          border-bottom: 1px solid #f0f0f0;
          cursor: pointer;
          transition: all 0.3s ease;

          &:hover {
            background: #fafafa;
            margin: 0 -8px;
            padding: 12px 8px;
          }

          &:last-child {
            border-bottom: none;
          }

          h4 {
            margin: 0 0 8px 0;
            font-size: 14px;
            font-weight: 500;
            color: #1a1a1a;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .recommend-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;

            .price {
              font-size: 16px;
              font-weight: 600;
              color: #ff4d4f;
            }

            .category {
              font-size: 12px;
              color: #999;
              background: #f5f5f5;
              padding: 2px 8px;
              border-radius: 4px;
            }
          }
        }
      }

      .announcement-list {
        .announcement-item {
          padding: 12px 0;
          border-bottom: 1px solid #f0f0f0;
          cursor: pointer;
          transition: all 0.3s ease;

          &:hover {
            background: #fafafa;
            margin: 0 -8px;
            padding: 12px 8px;
          }

          &:last-child {
            border-bottom: none;
          }

          .announcement-title {
            display: flex;
            align-items: center;
            font-size: 14px;
            font-weight: 500;
            color: #1a1a1a;
            margin-bottom: 4px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;

            .unread-dot {
              color: #ff4d4f;
              font-size: 8px;
              margin-right: 8px;
            }
          }

          .announcement-time {
            font-size: 12px;
            color: #999;
          }
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .dashboard-container {
    padding: 16px;

    .dashboard-header {
      .welcome-section {
        h1 {
          font-size: 24px;
        }

        p {
          font-size: 14px;
        }
      }

      .quick-actions {
        flex-direction: column;

        .el-button {
          width: 100%;
        }
      }
    }

    .dashboard-content {
      .content-card {
        padding: 16px;

        .task-list {
          .task-item {
            flex-direction: column;
            align-items: flex-start;

            .task-price {
              text-align: left;
              margin-top: 12px;

              .price {
                margin-bottom: 4px;
              }
            }
          }
        }
      }
    }
  }
}
</style>