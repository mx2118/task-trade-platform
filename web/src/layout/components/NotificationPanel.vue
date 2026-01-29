<template>
  <div class="notification-panel">
    <!-- 头部 -->
    <div class="panel-header">
      <h4>消息通知</h4>
      <el-button text size="small" @click="markAllAsRead">
        全部已读
      </el-button>
    </div>
    
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" class="notification-tabs">
      <el-tab-pane label="全部" name="all">
        <NotificationList
          :notifications="allNotifications"
          :loading="loading"
          @read="handleRead"
          @delete="handleDelete"
        />
      </el-tab-pane>
      
      <el-tab-pane label="未读" name="unread">
        <NotificationList
          :notifications="unreadNotifications"
          :loading="loading"
          @read="handleRead"
          @delete="handleDelete"
        />
      </el-tab-pane>
      
      <el-tab-pane label="系统" name="system">
        <NotificationList
          :notifications="systemNotifications"
          :loading="loading"
          @read="handleRead"
          @delete="handleDelete"
        />
      </el-tab-pane>
    </el-tabs>
    
    <!-- 底部 -->
    <div class="panel-footer">
      <el-button text size="small" @click="viewAll">
        查看全部消息
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface Notification {
  id: string
  type: 'system' | 'task' | 'payment' | 'message'
  title: string
  content: string
  isRead: boolean
  createdAt: string
  data?: any
}

const emit = defineEmits<{
  read: [id: string]
}>()

const router = useRouter()

const activeTab = ref('all')
const loading = ref(false)

// 模拟数据
const notifications = ref<Notification[]>([
  {
    id: '1',
    type: 'task',
    title: '任务申请通过',
    content: '您申请的"网站设计任务"已通过审核，请尽快开始工作。',
    isRead: false,
    createdAt: '2024-01-15 10:30:00'
  },
  {
    id: '2',
    type: 'payment',
    title: '收款通知',
    content: '您已收到任务报酬 ¥500.00，请查看钱包余额。',
    isRead: false,
    createdAt: '2024-01-15 09:15:00'
  },
  {
    id: '3',
    type: 'system',
    title: '系统维护通知',
    content: '系统将于今晚23:00-24:00进行维护，期间可能影响正常使用。',
    isRead: true,
    createdAt: '2024-01-14 16:20:00'
  },
  {
    id: '4',
    type: 'message',
    title: '新消息',
    content: '张三给您发送了一条消息，请及时查看。',
    isRead: true,
    createdAt: '2024-01-14 14:30:00'
  }
])

// 计算属性
const allNotifications = computed(() => notifications.value)

const unreadNotifications = computed(() => 
  notifications.value.filter(n => !n.isRead)
)

const systemNotifications = computed(() => 
  notifications.value.filter(n => n.type === 'system')
)

// 方法
const handleRead = (id: string) => {
  const notification = notifications.value.find(n => n.id === id)
  if (notification) {
    notification.isRead = true
    emit('read', id)
  }
}

const handleDelete = (id: string) => {
  const index = notifications.value.findIndex(n => n.id === id)
  if (index > -1) {
    notifications.value.splice(index, 1)
  }
}

const markAllAsRead = () => {
  notifications.value.forEach(n => {
    n.isRead = true
  })
  emit('read', 'all')
}

const viewAll = () => {
  router.push('/messages')
}

// 获取通知数据
const fetchNotifications = async () => {
  loading.value = true
  try {
    // 这里调用API获取真实数据
    // const response = await api.getNotifications()
    // notifications.value = response.data
  } catch (error) {
    console.error('获取通知失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchNotifications()
})
</script>

<style lang="scss" scoped>
.notification-panel {
  width: 360px;
  max-height: 480px;
  display: flex;
  flex-direction: column;
  background: $bg-white;
  border-radius: $border-radius-base;
  box-shadow: $box-shadow-dark;
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-md;
  border-bottom: 1px solid $border-lighter;
  
  h4 {
    margin: 0;
    font-size: $font-size-medium;
    font-weight: $font-weight-primary;
    color: $text-primary;
  }
}

.notification-tabs {
  flex: 1;
  overflow: hidden;
  
  :deep(.el-tabs__header) {
    margin: 0;
    padding: 0 $spacing-sm;
    border-bottom: 1px solid $border-lighter;
  }
  
  :deep(.el-tabs__nav-wrap) {
    &::after {
      display: none;
    }
  }
  
  :deep(.el-tabs__item) {
    padding: 0 $spacing-md;
    font-size: $font-size-sm;
    height: 40px;
    line-height: 40px;
    
    &.is-active {
      color: $primary-color;
      font-weight: $font-weight-primary;
    }
  }
  
  :deep(.el-tabs__active-bar) {
    background-color: $primary-color;
  }
  
  :deep(.el-tab-pane) {
    height: 320px;
    overflow-y: auto;
    padding: 0;
  }
}

.panel-footer {
  padding: $spacing-sm $spacing-md;
  border-top: 1px solid $border-lighter;
  text-align: center;
  
  .el-button {
    color: $primary-color;
    
    &:hover {
      color: lighten($primary-color, 10%);
    }
  }
}

// 响应式
@media (max-width: $breakpoint-sm) {
  .notification-panel {
    width: 100vw;
    max-width: 100vw;
    height: 100vh;
    max-height: 100vh;
    border-radius: 0;
  }
}

// 暗色主题
.dark .notification-panel {
  background: var(--el-bg-color-overlay);
  
  .panel-header {
    border-bottom-color: var(--el-border-color);
    
    h4 {
      color: var(--el-text-color-primary);
    }
  }
  
  .notification-tabs {
    :deep(.el-tabs__header) {
      border-bottom-color: var(--el-border-color);
    }
    
    :deep(.el-tabs__item) {
      color: var(--el-text-color-secondary);
      
      &.is-active {
        color: $primary-color;
      }
    }
  }
  
  .panel-footer {
    border-top-color: var(--el-border-color);
  }
}
</style>