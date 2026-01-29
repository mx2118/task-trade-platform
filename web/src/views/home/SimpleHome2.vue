<template>
  <div class="simple-home">
    <!-- 首页横幅 -->
    <div class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">🚀 任务交易平台</h1>
        <p class="hero-subtitle">发现优质任务，实现价值变现</p>
        <div class="hero-actions">
          <button class="btn btn-primary" @click="goToTasks">
            浏览任务
          </button>
          <button class="btn btn-success" @click="goToPublish">
            发布任务
          </button>
        </div>
      </div>
    </div>

    <!-- 任务列表预览 -->
    <div class="tasks-preview">
      <div class="section-header">
        <h2>🔥 热门任务</h2>
        <button class="link-btn" @click="goToTasks">
          查看更多 →
        </button>
      </div>
      
      <div class="tasks-grid">
        <div
          v-for="task in tasks"
          :key="task.id"
          class="task-card"
          @click="goToTaskDetail(task.id)"
        >
          <div class="task-header">
            <h3>{{ task.title }}</h3>
            <div class="task-price">¥{{ task.amount }}</div>
          </div>
          
          <div class="task-description">
            {{ task.description }}
          </div>
          
          <div class="task-footer">
            <span class="deadline">截止：{{ formatDate(task.deadline) }}</span>
            <span class="publisher">{{ task.publisher_name }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 特性介绍 -->
    <div class="features-section">
      <div class="section-header">
        <h2>✨ 平台特色</h2>
      </div>
      
      <div class="features-grid">
        <div class="feature-card">
          <div class="feature-icon">🛡️</div>
          <h3>真实可靠</h3>
          <p>严格的实名认证体系，确保任务和用户的真实性</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">💰</div>
          <h3>快速结算</h3>
          <p>任务完成后即时结算，资金安全有保障</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">🏆</div>
          <h3>海量任务</h3>
          <p>涵盖各类技能领域，总有一款适合你</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">🎯</div>
          <h3>贴心服务</h3>
          <p>7x24小时客服支持，随时解决您的问题</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 状态变量
const tasks = ref([
  {
    id: 1,
    title: '网站前端开发',
    description: '需要开发一个企业官网，要求响应式设计，兼容主流浏览器',
    amount: 5000,
    deadline: new Date(Date.now() + 7 * 24 * 3600 * 1000).toISOString(),
    publisher_name: '张先生',
    status: 1
  },
  {
    id: 2,
    title: 'Logo设计',
    description: '为新创公司设计一个现代简约风格的Logo，提供多个方案',
    amount: 2000,
    deadline: new Date(Date.now() + 5 * 24 * 3600 * 1000).toISOString(),
    publisher_name: '李女士',
    status: 1
  },
  {
    id: 3,
    title: '小程序开发',
    description: '开发一个电商类小程序，包含商品展示、购物车、订单管理等功能',
    amount: 8000,
    deadline: new Date(Date.now() + 10 * 24 * 3600 * 1000).toISOString(),
    publisher_name: '王总',
    status: 1
  },
  {
    id: 4,
    title: 'UI界面设计',
    description: '设计一套移动App的UI界面，风格简约现代',
    amount: 3500,
    deadline: new Date(Date.now() + 6 * 24 * 3600 * 1000).toISOString(),
    publisher_name: '刘经理',
    status: 1
  },
  {
    id: 5,
    title: '数据分析报告',
    description: '分析电商平台用户数据，提供详细的数据分析报告',
    amount: 4000,
    deadline: new Date(Date.now() + 8 * 24 * 3600 * 1000).toISOString(),
    publisher_name: '陈总监',
    status: 1
  },
  {
    id: 6,
    title: '文案撰写',
    description: '为产品撰写营销文案，包括产品介绍、宣传语等',
    amount: 1500,
    deadline: new Date(Date.now() + 3 * 24 * 3600 * 1000).toISOString(),
    publisher_name: '周小姐',
    status: 1
  }
])

// 格式化日期
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = date.getTime() - now.getTime()
  const days = Math.ceil(diff / (1000 * 3600 * 24))
  
  if (days < 0) return '已截止'
  if (days === 0) return '今天'
  if (days === 1) return '明天'
  return `${days}天后`
}

// 跳转到任务列表
const goToTasks = () => {
  router.push('/layout/tasks')
}

// 跳转到发布任务
const goToPublish = () => {
  router.push('/login')
}

// 跳转到任务详情
const goToTaskDetail = (id: number) => {
  router.push(`/layout/tasks/${id}`)
}

onMounted(() => {
  console.log('SimpleHome2 组件已挂载')
})
</script>

<style scoped>
.simple-home {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.hero-section {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 80px 24px;
  text-align: center;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
}

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
}

.btn {
  height: 48px;
  padding: 0 32px;
  font-size: 16px;
  border: none;
  border-radius: 24px;
  cursor: pointer;
  transition: all 0.3s;
  font-weight: 500;
}

.btn-primary {
  background: #409eff;
  color: white;
}

.btn-primary:hover {
  background: #66b1ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
}

.btn-success {
  background: #67c23a;
  color: white;
}

.btn-success:hover {
  background: #85ce61;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.4);
}

.tasks-preview,
.features-section {
  max-width: 1200px;
  margin: 0 auto;
  padding: 48px 24px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32px;
}

.section-header h2 {
  font-size: 28px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.link-btn {
  background: none;
  border: none;
  color: #fff;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s;
}

.link-btn:hover {
  opacity: 0.8;
}

.tasks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 24px;
}

.task-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
}

.task-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  transform: translateY(-4px);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.task-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  flex: 1;
}

.task-price {
  color: #ff4d4f;
  font-size: 20px;
  font-weight: 700;
  margin-left: 16px;
  white-space: nowrap;
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

.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
  font-size: 13px;
  color: #999;
}

.features-section {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 20px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
}

.feature-card {
  text-align: center;
  padding: 32px 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.feature-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.feature-card h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 12px 0;
}

.feature-card p {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .hero-title {
    font-size: 32px;
  }
  
  .hero-subtitle {
    font-size: 16px;
  }
  
  .tasks-grid {
    grid-template-columns: 1fr;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .btn {
    width: 100%;
  }
}
</style>
