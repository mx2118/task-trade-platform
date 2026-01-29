<template>
  <teleport to="body">
    <transition-group name="message" tag="div" class="global-message-container">
      <div
        v-for="message in messages"
        :key="message.id"
        :class="[
          'global-message-item',
          `message-${message.type}`,
          { 'message-with-icon': message.showIcon }
        ]"
        @click="removeMessage(message.id)"
      >
        <el-icon v-if="message.showIcon" class="message-icon">
          <component :is="getIcon(message.type)" />
        </el-icon>
        <div class="message-content">
          <div v-if="message.title" class="message-title">{{ message.title }}</div>
          <div class="message-text">{{ message.text }}</div>
        </div>
        <el-icon v-if="message.closable" class="message-close" @click.stop="removeMessage(message.id)">
          <Close />
        </el-icon>
      </div>
    </transition-group>
  </teleport>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Close, SuccessFilled, WarningFilled, InfoFilled, CircleCloseFilled } from '@element-plus/icons-vue'

interface MessageItem {
  id: string
  type: 'success' | 'warning' | 'info' | 'error'
  title?: string
  text: string
  duration?: number
  closable?: boolean
  showIcon?: boolean
}

const messages = ref<MessageItem[]>([])

let messageIdCounter = 0

// 添加消息
const addMessage = (message: Omit<MessageItem, 'id'>) => {
  const id = `message-${++messageIdCounter}`
  const messageItem: MessageItem = {
    id,
    duration: 4500,
    closable: true,
    showIcon: true,
    ...message
  }
  
  messages.value.push(messageItem)
  
  // 自动移除消息
  if (messageItem.duration && messageItem.duration > 0) {
    setTimeout(() => {
      removeMessage(id)
    }, messageItem.duration)
  }
  
  return id
}

// 移除消息
const removeMessage = (id: string) => {
  const index = messages.value.findIndex(msg => msg.id === id)
  if (index > -1) {
    messages.value.splice(index, 1)
  }
}

// 清空所有消息
const clearMessages = () => {
  messages.value = []
}

// 获取图标
const getIcon = (type: string) => {
  const iconMap = {
    success: SuccessFilled,
    warning: WarningFilled,
    info: InfoFilled,
    error: CircleCloseFilled
  }
  return iconMap[type] || InfoFilled
}

// 扩展 window 类型
declare global {
  interface Window {
    $message: {
      success: (text: string, options?: Partial<MessageItem>) => void
      warning: (text: string, options?: Partial<MessageItem>) => void
      error: (text: string, options?: Partial<MessageItem>) => void
      info: (text: string, options?: Partial<MessageItem>) => void
      close: (id?: number) => void
      closeAll: () => void
      clear: () => void
    }
  }
}

// 全局消息方法
window.$message = {
  success: (text: string, options?: Partial<MessageItem>) => {
    return addMessage({ type: 'success', text, ...options })
  },
  warning: (text: string, options?: Partial<MessageItem>) => {
    return addMessage({ type: 'warning', text, ...options })
  },
  info: (text: string, options?: Partial<MessageItem>) => {
    return addMessage({ type: 'info', text, ...options })
  },
  error: (text: string, options?: Partial<MessageItem>) => {
    return addMessage({ type: 'error', text, ...options })
  },
  close: (id?: number) => {
    if (id !== undefined) {
      removeMessage(id)
    }
  },
  closeAll: clearMessages,
  clear: clearMessages
}

// 监听键盘事件，ESC键清空消息
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    clearMessages()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style lang="scss" scoped>
.global-message-container {
  position: fixed;
  top: 80px;
  right: 20px;
  z-index: $z-index-notification;
  pointer-events: none;
  max-width: 400px;
  width: 100%;
  
  @media (max-width: $breakpoint-sm) {
    top: 60px;
    left: 10px;
    right: 10px;
    max-width: none;
  }
}

.global-message-item {
  position: relative;
  display: flex;
  align-items: flex-start;
  padding: $spacing-md;
  margin-bottom: $spacing-sm;
  border-radius: $border-radius-base;
  box-shadow: $box-shadow-light;
  background: $bg-white;
  border-left: 4px solid;
  pointer-events: auto;
  cursor: pointer;
  transition: all $transition-base;
  
  &:hover {
    transform: translateX(-4px);
    box-shadow: $box-shadow-hover;
  }
  
  &.message-success {
    border-left-color: $success-color;
    background: linear-gradient(135deg, rgba($success-color, 0.1) 0%, rgba($success-color, 0.05) 100%);
    
    .message-title {
      color: $success-color;
    }
  }
  
  &.message-warning {
    border-left-color: $warning-color;
    background: linear-gradient(135deg, rgba($warning-color, 0.1) 0%, rgba($warning-color, 0.05) 100%);
    
    .message-title {
      color: $warning-color;
    }
  }
  
  &.message-info {
    border-left-color: $info-color;
    background: linear-gradient(135deg, rgba($info-color, 0.1) 0%, rgba($info-color, 0.05) 100%);
    
    .message-title {
      color: $info-color;
    }
  }
  
  &.message-error {
    border-left-color: $danger-color;
    background: linear-gradient(135deg, rgba($danger-color, 0.1) 0%, rgba($danger-color, 0.05) 100%);
    
    .message-title {
      color: $danger-color;
    }
  }
  
  .message-icon {
    margin-right: $spacing-sm;
    margin-top: 2px;
    flex-shrink: 0;
  }
  
  .message-content {
    flex: 1;
    min-width: 0;
  }
  
  .message-title {
    font-weight: $font-weight-primary;
    margin-bottom: $spacing-xs;
    font-size: $font-size-medium;
  }
  
  .message-text {
    color: $text-regular;
    line-height: 1.4;
    word-break: break-word;
  }
  
  .message-close {
    margin-left: $spacing-sm;
    flex-shrink: 0;
    color: $text-secondary;
    cursor: pointer;
    transition: color 0.2s;
    
    &:hover {
      color: $text-primary;
    }
  }
  
  &:not(.message-with-icon) {
    .message-content {
      margin-left: 0;
    }
  }
}

// 过渡动画
.message-enter-active {
  transition: all 0.3s ease-out;
}

.message-leave-active {
  transition: all 0.2s ease-in;
}

.message-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.message-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.message-move {
  transition: transform 0.3s ease;
}

// 暗色主题
.dark .global-message-item {
  background: var(--el-bg-color-overlay);
  
  &.message-success {
    background: linear-gradient(135deg, rgba($success-color, 0.2) 0%, rgba($success-color, 0.1) 100%);
  }
  
  &.message-warning {
    background: linear-gradient(135deg, rgba($warning-color, 0.2) 0%, rgba($warning-color, 0.1) 100%);
  }
  
  &.message-info {
    background: linear-gradient(135deg, rgba($info-color, 0.2) 0%, rgba($info-color, 0.1) 100%);
  }
  
  &.message-error {
    background: linear-gradient(135deg, rgba($danger-color, 0.2) 0%, rgba($danger-color, 0.1) 100%);
  }
  
  .message-text {
    color: var(--el-text-color-regular);
  }
  
  .message-close {
    color: var(--el-text-color-secondary);
    
    &:hover {
      color: var(--el-text-color-primary);
    }
  }
}
</style>