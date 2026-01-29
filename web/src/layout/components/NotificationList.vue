<template>
  <div class="notification-list">
    <el-scrollbar height="320px">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="3" animated />
      </div>
      
      <div v-else-if="notifications.length === 0" class="empty-container">
        <el-empty description="暂无通知" :image-size="60">
          <template #image>
            <el-icon :size="40" class="empty-icon">
              <Bell />
            </el-icon>
          </template>
        </el-empty>
      </div>
      
      <div v-else class="notification-items">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          :class="[
            'notification-item',
            { 'notification-unread': !notification.isRead }
          ]"
          @click="handleClick(notification)"
        >
          <div class="notification-content">
            <div class="notification-header">
              <span class="notification-title">{{ notification.title }}</span>
              <span class="notification-time">{{ formatTime(notification.createdAt) }}</span>
            </div>
            
            <p class="notification-text">{{ notification.content }}</p>
            
            <div class="notification-footer">
              <el-tag
                :type="getTagType(notification.type)"
                size="small"
                class="notification-tag"
              >
                {{ getTypeLabel(notification.type) }}
              </el-tag>
              
              <div class="notification-actions">
                <el-button
                  v-if="!notification.isRead"
                  text
                  size="small"
                  @click.stop="markAsRead(notification.id)"
                >
                  标记已读
                </el-button>
                
                <el-dropdown trigger="click" @command="(cmd) => handleAction(cmd, notification)">
                  <el-button text size="small">
                    <el-icon>
                      <MoreFilled />
                    </el-icon>
                  </el-button>
                  
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="read" v-if="!notification.isRead">
                        <el-icon><View /></el-icon>
                        标记已读
                      </el-dropdown-item>
                      <el-dropdown-item command="unread" v-if="notification.isRead">
                        <el-icon><Hide /></el-icon>
                        标记未读
                      </el-dropdown-item>
                      <el-dropdown-item command="delete" divided>
                        <el-icon><Delete /></el-icon>
                        删除通知
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import {
  Bell,
  View,
  Hide,
  Delete,
  MoreFilled
} from '@element-plus/icons-vue'

interface Notification {
  id: string
  type: 'system' | 'task' | 'payment' | 'message'
  title: string
  content: string
  isRead: boolean
  createdAt: string
  data?: any
}

interface Props {
  notifications: Notification[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  read: [id: string]
  delete: [id: string]
}>()

const router = useRouter()

// 配置dayjs
dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

// 格式化时间
const formatTime = (time: string) => {
  const date = dayjs(time)
  const now = dayjs()
  
  // 如果是今天，显示相对时间
  if (date.isSame(now, 'day')) {
    return date.fromNow()
  }
  
  // 如果是昨天，显示"昨天 HH:mm"
  if (date.isSame(now.subtract(1, 'day'), 'day')) {
    return `昨天 ${date.format('HH:mm')}`
  }
  
  // 如果是今年，显示"MM-DD HH:mm"
  if (date.isSame(now, 'year')) {
    return date.format('MM-DD HH:mm')
  }
  
  // 否则显示"YYYY-MM-DD HH:mm"
  return date.format('YYYY-MM-DD HH:mm')
}

// 获取标签类型
const getTagType = (type: string) => {
  const typeMap = {
    system: 'info',
    task: 'success',
    payment: 'warning',
    message: 'primary'
  }
  return typeMap[type] || 'info'
}

// 获取类型标签
const getTypeLabel = (type: string) => {
  const labelMap = {
    system: '系统',
    task: '任务',
    payment: '支付',
    message: '消息'
  }
  return labelMap[type] || '通知'
}

// 处理点击
const handleClick = (notification: Notification) => {
  // 标记为已读
  if (!notification.isRead) {
    markAsRead(notification.id)
  }
  
  // 根据通知类型跳转
  const routes: Record<string, string> = {
    task: '/my-tasks',
    payment: '/wallet',
    message: '/messages'
  }
  
  const route = routes[notification.type]
  if (route) {
    router.push(route)
  }
}

// 标记已读
const markAsRead = (id: string) => {
  emit('read', id)
}

// 标记未读
const markAsUnread = (id: string) => {
  // 这里可以实现标记未读的逻辑
  console.log('标记未读:', id)
}

// 处理操作
const handleAction = (command: string, notification: Notification) => {
  switch (command) {
    case 'read':
      markAsRead(notification.id)
      break
    case 'unread':
      markAsUnread(notification.id)
      break
    case 'delete':
      emit('delete', notification.id)
      break
  }
}
</script>

<style lang="scss" scoped>
.notification-list {
  height: 320px;
}

.loading-container {
  padding: $spacing-lg;
}

.empty-container {
  height: 320px;
  display: flex;
  align-items: center;
  justify-content: center;
  
  .empty-icon {
    color: $text-placeholder;
  }
}

.notification-items {
  .notification-item {
    padding: $spacing-md;
    border-bottom: 1px solid $border-extra-light;
    cursor: pointer;
    transition: background-color 0.2s;
    
    &:hover {
      background-color: $bg-page;
    }
    
    &:last-child {
      border-bottom: none;
    }
    
    &.notification-unread {
      background-color: rgba($primary-color, 0.03);
      border-left: 3px solid $primary-color;
      padding-left: calc(#{$spacing-md} - 3px);
    }
  }
  
  .notification-content {
    .notification-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: $spacing-xs;
    }
    
    .notification-title {
      font-weight: $font-weight-primary;
      color: $text-primary;
      font-size: $font-size-sm;
      flex: 1;
      @extend .text-ellipsis;
    }
    
    .notification-time {
      font-size: $font-size-extra-small;
      color: $text-placeholder;
      flex-shrink: 0;
      margin-left: $spacing-sm;
    }
    
    .notification-text {
      color: $text-regular;
      font-size: $font-size-sm;
      line-height: 1.4;
      margin-bottom: $spacing-sm;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .notification-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    
    .notification-tag {
      flex-shrink: 0;
    }
    
    .notification-actions {
      display: flex;
      gap: $spacing-xs;
      opacity: 0;
      transition: opacity 0.2s;
    }
  }
  
  .notification-item:hover .notification-actions {
    opacity: 1;
  }
}

// 暗色主题
.dark .notification-items {
  .notification-item {
    border-bottom-color: var(--el-border-color);
    
    &:hover {
      background-color: var(--el-bg-color);
    }
    
    &.notification-unread {
      background-color: rgba($primary-color, 0.08);
    }
  }
  
  .notification-title {
    color: var(--el-text-color-primary);
  }
  
  .notification-time {
    color: var(--el-text-color-placeholder);
  }
  
  .notification-text {
    color: var(--el-text-color-regular);
  }
}
</style>