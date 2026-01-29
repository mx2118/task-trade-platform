<template>
  <div class="simple-home">
    <!-- 首页横幅 -->
    <div class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">🚀 任务交易平台</h1>
        <p class="hero-subtitle">发现优质任务，实现价值变现</p>
        <div class="hero-actions">
          <el-button type="primary" size="large" @click="goToTasks">
            <el-icon><Search /></el-icon>
            浏览任务
          </el-button>
          <el-button type="success" size="large" @click="goToPublish">
            <el-icon><Plus /></el-icon>
            发布任务
          </el-button>
        </div>
      </div>
    </div>

    <!-- 任务列表预览 -->
    <div class="tasks-preview">
      <div class="section-header">
        <h2>🔥 热门任务</h2>
        <el-button text @click="goToTasks">
          查看更多
          <el-icon><ArrowRight /></el-icon>
        </el-button>
      </div>
      
      <div class="tasks-grid" v-loading="loading">
        <div
          v-for="task in tasks"
          :key="task.id"
          class="task-item"
          @click="goToTaskDetail(task.id)"
        >
          <div class="task-card">
            <div class="task-header">
              <h3>{{ task.title }}</h3>
              <div class="task-price">¥{{ task.amount || task.price }}</div>
            </div>
            
            <div class="task-description">
              {{ task.description || task.content }}
            </div>
            
            <div class="task-meta">
              <span class="deadline">
                <el-icon><Clock /></el-icon>
                截止：{{ formatRelativeTime(task.deadline) }}
              </span>
              <el-avatar :size="32" :src="task.publisher_avatar">
                {{ (task.publisher_name || 'U')?.charAt(0) }}
              </el-avatar>
            </div>
            
            <div class="task-footer">
              <el-button
                v-if="task.status === 1 || task.status === 'pending'"
                type="primary"
                size="small"
                @click.stop="handleApply(task)"
              >
                申请接取
              </el-button>
              <el-button size="small" @click.stop="handleShare(task)">
                分享
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 特性介绍 -->
    <div class="features-section">
      <div class="section-header">
        <h2>✨ 平台特色</h2>
      </div>
      
      <el-row :gutter="24">
        <el-col :xs="24" :sm="12" :lg="6" v-for="feature in features" :key="feature.title">
          <div class="feature-card">
            <div class="feature-icon" :style="{ color: feature.color }">
              <el-icon size="32">
                <component :is="feature.icon" />
              </el-icon>
            </div>
            <h3 class="feature-title">{{ feature.title }}</h3>
            <p class="feature-description">{{ feature.description }}</p>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Search, 
  Plus, 
  ArrowRight,
  Clock,
  Shield,
  Wallet,
  Trophy,
  Service
} from '@element-plus/icons-vue'
import { taskApi } from '@/api'
import { formatRelativeTime } from '@/utils/format'
import type { Task } from '@/types'

const router = useRouter()

// 状态变量
const tasks = ref<Task[]>([])
const loading = ref(false)

// 特性数据
const features = [
  {
    title: '真实可靠',
    description: '严格的实名认证体系，确保任务和用户的真实性',
    icon: Shield,
    color: '#67c23a'
  },
  {
    title: '快速结算',
    description: '任务完成后即时结算，资金安全有保障',
    icon: Wallet,
    color: '#e6a23c'
  },
  {
    title: '海量任务',
    description: '涵盖各类技能领域，总有一款适合你',
    icon: Trophy,
    color: '#f56c6c'
  },
  {
    title: '贴心服务',
    description: '7x24小时客服支持，随时解决您的问题',
    icon: Service,
    color: '#409eff'
  }
]

// 跳转到任务列表
const goToTasks = () => {
  router.push('/tasks')
}

// 跳转到发布任务
const goToPublish = () => {
  router.push('/tasks/publish')
}

// 跳转到任务详情
const goToTaskDetail = (id: number) => {
  router.push(`/tasks/${id}`)
}

// 申请接取任务
const handleApply = (task: Task) => {
  ElMessage.info('请先登录后再申请任务')
  goToTaskDetail(task.id)
}

// 分享任务
const handleShare = async (task: Task) => {
  const shareUrl = `${window.location.origin}/tasks/${task.id}`
  
  try {
    await navigator.clipboard.writeText(shareUrl)
    ElMessage.success('链接已复制到剪贴板')
  } catch (error) {
    ElMessage.info(`任务链接：${shareUrl}`)
  }
}

// 加载任务列表
const loadTasks = async () => {
  // 先设置一些模拟数据，让页面立即显示
  tasks.value = [
    {
      id: 1,
      title: '网站前端开发',
      description: '需要开发一个企业官网，要求响应式设计，兼容主流浏览器',
      amount: 5000,
      deadline: new Date(Date.now() + 7 * 24 * 3600 * 1000).toISOString(),
      publisher_name: '张先生',
      publisher_avatar: '',
      status: 1
    },
    {
      id: 2,
      title: 'Logo设计',
      description: '为新创公司设计一个现代简约风格的Logo，提供多个方案',
      amount: 2000,
      deadline: new Date(Date.now() + 5 * 24 * 3600 * 1000).toISOString(),
      publisher_name: '李女士',
      publisher_avatar: '',
      status: 1
    },
    {
      id: 3,
      title: '小程序开发',
      description: '开发一个电商类小程序，包含商品展示、购物车、订单管理等功能',
      amount: 8000,
      deadline: new Date(Date.now() + 10 * 24 * 3600 * 1000).toISOString(),
      publisher_name: '王总',
      publisher_avatar: '',
      status: 1
    }
  ]
  
  // 异步尝试加载真实数据
  try {
    loading.value = true
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 3000)
    
    const response = await taskApi.getTasks({
      limit: 6,
      sort: 'view_count:desc'
    })
    
    clearTimeout(timeoutId)
    
    if (response.data.list && response.data.list.length > 0) {
      tasks.value = response.data.list
    }
  } catch (error) {
    console.log('使用模拟数据:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadTasks()
})
</script>

<style lang="scss" scoped>
.simple-home {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .hero-section {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    padding: 80px 24px;
    text-align: center;
    
    .hero-content {
      max-width: 800px;
      margin: 0 auto;
      
      .hero-title {
        font-size: 48px;
        font-weight: 700;
        color: #fff;
        margin: 0 0 16px 0;
        text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
      }
      
      .hero-subtitle {
        font-size: 20px;
        color: rgba(255, 255, 255, 0.9);
        margin: 0 0 32px 0;
      }
      
      .hero-actions {
        display: flex;
        gap: 16px;
        justify-content: center;
        flex-wrap: wrap;
        
        .el-button {
          height: 48px;
          padding: 0 24px;
          font-size: 16px;
          border-radius: 24px;
          
          .el-icon {
            margin-right: 8px;
          }
        }
      }
    }
  }
  
  .tasks-preview,
  .features-section {
    max-width: 1200px;
    margin: 0 auto;
    padding: 48px 24px;
    
    .section-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 32px;
      
      h2 {
        font-size: 28px;
        font-weight: 600;
        color: #fff;
        margin: 0;
      }
      
      .el-button {
        color: #fff;
        
        .el-icon {
          margin-left: 4px;
        }
      }
    }
  }
  
  .tasks-preview {
    .tasks-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
      gap: 24px;
      
      .task-item {
        cursor: pointer;
        
        .task-card {
          background: white;
          border-radius: 12px;
          padding: 20px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          transition: all 0.3s ease;
          
          &:hover {
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
            transform: translateY(-4px);
          }
          
          .task-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 12px;
            
            h3 {
              margin: 0 0 8px 0;
              font-size: 16px;
              font-weight: 600;
              color: #1a1a1a;
              line-height: 1.4;
            }
            
            .task-price {
              text-align: right;
              margin-left: 16px;
              
              color: #ff4d4f;
              font-size: 20px;
              font-weight: 700;
            }
          }
          
          .task-description {
            color: #666;
            font-size: 14px;
            line-height: 1.5;
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
            
            .deadline {
              display: flex;
              align-items: center;
              gap: 4px;
              font-size: 12px;
              color: #666;
              
              .el-icon {
                font-size: 14px;
              }
            }
          }
          
          .task-footer {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding-top: 12px;
            border-top: 1px solid #f0f0f0;
            gap: 8px;
          }
        }
      }
    }
  }
  
  .features-section {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 20px;
    
    .feature-card {
      text-align: center;
      padding: 32px 20px;
      background: rgba(255, 255, 255, 0.95);
      border-radius: 12px;
      transition: all 0.3s ease;
      margin-bottom: 20px;
      
      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
      }
      
      .feature-icon {
        margin-bottom: 16px;
      }
      
      .feature-title {
        font-size: 18px;
        font-weight: 600;
        color: #1a1a1a;
        margin: 0 0 12px 0;
      }
      
      .feature-description {
        font-size: 14px;
        color: #666;
        line-height: 1.6;
        margin: 0;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .simple-home {
    .hero-section {
      padding: 60px 16px;
      
      .hero-content {
        .hero-title {
          font-size: 32px;
        }
        
        .hero-subtitle {
          font-size: 16px;
        }
        
        .hero-actions {
          .el-button {
            width: 100%;
            justify-content: center;
          }
        }
      }
    }
    
    .tasks-preview,
    .features-section {
      padding: 32px 16px;
      
      .section-header {
        flex-direction: column;
        gap: 16px;
        
        h2 {
          font-size: 24px;
        }
      }
    }
    
    .tasks-preview {
      .tasks-grid {
        grid-template-columns: 1fr;
        gap: 16px;
        
        .task-item {
          .task-card {
            padding: 16px;
            
            .task-header {
              flex-direction: column;
              
              .task-price {
                margin-left: 0;
                margin-top: 8px;
                text-align: left;
              }
            }
            
            .task-meta {
              flex-direction: column;
              align-items: flex-start;
              gap: 8px;
            }
            
            .task-footer {
              flex-direction: column;
              gap: 8px;
            }
          }
        }
      }
    }
    
    .features-section {
      .feature-card {
        padding: 24px 16px;
      }
    }
  }
}
</style>