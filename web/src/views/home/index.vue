<template>
  <div class="home-container">
    <!-- 首页横幅 -->
    <div class="hero-banner">
      <div class="hero-content">
        <div class="hero-text">
          <h1 class="hero-title">🚀 任务交易平台</h1>
          <p class="hero-subtitle">发现优质任务，实现价值变现</p>
          <div class="hero-actions">
            <el-button type="primary" size="large" @click="goToTasks" class="action-btn">
              <el-icon><Search /></el-icon>
              浏览任务
            </el-button>
            <el-button type="success" size="large" @click="goToPublish" class="action-btn">
              <el-icon><Plus /></el-icon>
              发布任务
            </el-button>
          </div>
        </div>
        <div class="hero-stats">
          <div class="stat-item">
            <div class="stat-number">{{ stats.total }}</div>
            <div class="stat-label">总任务</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ stats.active_users }}</div>
            <div class="stat-label">活跃用户</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ stats.completed_today }}</div>
            <div class="stat-label">今日完成</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 热门任务预览 -->
    <div class="featured-tasks">
      <div class="section-header">
        <h2>🔥 热门任务</h2>
        <el-button text @click="goToTasks">
          查看更多
          <el-icon><ArrowRight /></el-icon>
        </el-button>
      </div>
      
      <div class="tasks-grid">
        <el-row :gutter="20">
          <el-col 
            v-for="task in featuredTasks" 
            :key="task.id"
            :xs="24" 
            :sm="12" 
            :lg="8" 
            :xl="6"
          >
            <TaskCard :task="task" @click="goToTaskDetail(task.id)" />
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 快速分类 -->
    <div class="quick-categories">
      <div class="section-header">
        <h2>📂 快速分类</h2>
      </div>
      
      <div class="categories-grid">
        <div 
          v-for="category in categories.slice(0, 8)" 
          :key="category.id"
          class="category-item"
          @click="goToTasksByCategory(category.id)"
        >
          <div class="category-icon">
            <el-icon size="24">
              <component :is="getCategoryIcon(category.id)" />
            </el-icon>
          </div>
          <div class="category-info">
            <div class="category-name">{{ category.name }}</div>
            <div class="category-count">{{ category.task_count }} 个任务</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 为什么选择我们 -->
    <div class="features-section">
      <div class="section-header">
        <h2>✨ 为什么选择我们</h2>
      </div>
      
      <el-row :gutter="24">
        <el-col :xs="24" :sm="12" :lg="6" v-for="feature in features" :key="feature.title">
          <div class="feature-card">
            <div class="feature-icon">
              <el-icon size="32" :color="feature.color">
                <component :is="feature.icon" />
              </el-icon>
            </div>
            <h3 class="feature-title">{{ feature.title }}</h3>
            <p class="feature-description">{{ feature.description }}</p>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 最新动态 -->
    <div class="latest-activities">
      <div class="section-header">
        <h2>📈 最新动态</h2>
      </div>
      
      <el-timeline>
        <el-timeline-item
          v-for="activity in activities"
          :key="activity.id"
          :timestamp="formatTime(activity.created_at)"
          placement="top"
        >
          <div class="activity-content">
            <div class="activity-title">{{ activity.title }}</div>
            <div class="activity-description">{{ activity.description }}</div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Search, 
  Plus, 
  ArrowRight,
  Trophy,
  Shield,
  Clock,
  Wallet,
  Service
} from '@element-plus/icons-vue'
import { taskApi, categoryApi } from '@/api'
import { formatTime } from '@/utils/format'
import TaskCard from '@/components/common/TaskCard.vue'
import type { Task, Category } from '@/types'

const router = useRouter()

// 数据状态
const featuredTasks = ref<Task[]>([])
const categories = ref<Category[]>([])
const stats = ref({
  total: 0,
  active_users: 0,
  completed_today: 0
})
const activities = ref<any[]>([])

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

// 获取分类图标
const getCategoryIcon = (categoryId: number) => {
  const iconMap: Record<number, any> = {
    1: Trophy,    // 设计创意
    2: Service,   // 程序开发
    3: Clock,     // 文案写作
    4: Wallet,    // 数据分析
    5: Shield     // 其他
  }
  return iconMap[categoryId] || Trophy
}

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

// 按分类查看任务
const goToTasksByCategory = (categoryId: number) => {
  router.push({
    path: '/tasks',
    query: { category_id: categoryId }
  })
}

// 加载热门任务
const loadFeaturedTasks = async () => {
  try {
    const response = await taskApi.getTasks({
      limit: 8,
      sort: 'view_count:desc'
    })
    featuredTasks.value = response.data.list
  } catch (error) {
    console.error('加载热门任务失败:', error)
  }
}

// 加载分类
const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategories()
    categories.value = response.data
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await taskApi.getStats()
    stats.value = {
      ...response.data,
      completed_today: response.data.completed_today || 0,
      active_users: response.data.active_users || 0
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载最新动态
const loadActivities = async () => {
  try {
    // 模拟动态数据
    activities.value = [
      {
        id: 1,
        title: '用户 张三 完成了任务 "网站UI设计"',
        description: '获得赏金 ¥500.00，评分 5.0',
        created_at: new Date().toISOString()
      },
      {
        id: 2,
        title: '新任务 "移动APP开发" 已发布',
        description: '赏金 ¥2000.00，预计完成时间 7天',
        created_at: new Date(Date.now() - 3600000).toISOString()
      },
      {
        id: 3,
        title: '系统升级通知',
        description: '新增任务推荐功能，提升用户体验',
        created_at: new Date(Date.now() - 7200000).toISOString()
      }
    ]
  } catch (error) {
    console.error('加载动态失败:', error)
  }
}

// 格式化时间
const formatTime = (time: string) => {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  return `${Math.floor(diff / 86400000)}天前`
}

// 页面初始化
onMounted(() => {
  loadFeaturedTasks()
  loadCategories()
  loadStats()
  loadActivities()
})
</script>

<style lang="scss" scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .hero-banner {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    padding: 60px 24px;
    
    .hero-content {
      max-width: 1200px;
      margin: 0 auto;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 40px;
      
      .hero-text {
        flex: 1;
        
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
          flex-wrap: wrap;
          
          .action-btn {
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
      
      .hero-stats {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 32px;
        
        .stat-item {
          text-align: center;
          padding: 24px;
          background: rgba(255, 255, 255, 0.1);
          border-radius: 16px;
          backdrop-filter: blur(10px);
          
          .stat-number {
            font-size: 36px;
            font-weight: 700;
            color: #fff;
            margin-bottom: 8px;
          }
          
          .stat-label {
            font-size: 14px;
            color: rgba(255, 255, 255, 0.8);
          }
        }
      }
    }
  }
  
  .featured-tasks,
  .quick-categories,
  .features-section,
  .latest-activities {
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
  
  .featured-tasks {
    .tasks-grid {
      .task-card {
        background: rgba(255, 255, 255, 0.95);
        border-radius: 12px;
        padding: 20px;
        transition: all 0.3s ease;
        cursor: pointer;
        
        &:hover {
          transform: translateY(-4px);
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
        }
      }
    }
  }
  
  .quick-categories {
    .categories-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 20px;
      
      .category-item {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 20px;
        background: rgba(255, 255, 255, 0.95);
        border-radius: 12px;
        cursor: pointer;
        transition: all 0.3s ease;
        
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }
        
        .category-icon {
          width: 48px;
          height: 48px;
          display: flex;
          align-items: center;
          justify-content: center;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          border-radius: 12px;
        }
        
        .category-info {
          .category-name {
            font-size: 16px;
            font-weight: 600;
            color: #1a1a1a;
            margin-bottom: 4px;
          }
          
          .category-count {
            font-size: 14px;
            color: #666;
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
  
  .latest-activities {
    .el-timeline {
      .el-timeline-item {
        .el-timeline-item__content {
          background: rgba(255, 255, 255, 0.95);
          border-radius: 12px;
          padding: 16px;
          
          .activity-content {
            .activity-title {
              font-size: 16px;
              font-weight: 600;
              color: #1a1a1a;
              margin-bottom: 8px;
            }
            
            .activity-description {
              font-size: 14px;
              color: #666;
              line-height: 1.5;
              margin: 0;
            }
          }
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .home-container {
    .hero-banner {
      padding: 40px 16px;
      
      .hero-content {
        flex-direction: column;
        gap: 32px;
        
        .hero-text {
          .hero-title {
            font-size: 32px;
          }
          
          .hero-subtitle {
            font-size: 16px;
          }
          
          .hero-actions {
            .action-btn {
              width: 100%;
              justify-content: center;
            }
          }
        }
        
        .hero-stats {
          grid-template-columns: 1fr;
          gap: 16px;
          
          .stat-item {
            padding: 16px;
            
            .stat-number {
              font-size: 28px;
            }
          }
        }
      }
    }
    
    .featured-tasks,
    .quick-categories,
    .features-section,
    .latest-activities {
      padding: 32px 16px;
      
      .section-header {
        flex-direction: column;
        gap: 16px;
        
        h2 {
          font-size: 24px;
        }
      }
    }
    
    .categories-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>