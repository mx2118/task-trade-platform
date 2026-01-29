<template>
  <header class="layout-header">
    <div class="header-left">
      <!-- 侧边栏切换按钮 -->
      <el-button
        class="sidebar-toggle"
        text
        @click="$emit('toggleSidebar')"
      >
        <el-icon :size="20">
          <Menu />
        </el-icon>
      </el-button>
      
      <!-- Logo -->
      <router-link to="/" class="logo">
        <img src="/logo.svg" alt="任务交易平台" class="logo-img" />
        <span v-if="!isMobile" class="logo-text">任务交易平台</span>
      </router-link>
    </div>
    
    <div class="header-center">
      <!-- 搜索框 -->
      <el-input
        v-model="searchKeyword"
        placeholder="搜索任务、用户..."
        class="search-input"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <el-icon>
            <Search />
          </el-icon>
        </template>
      </el-input>
    </div>
    
    <div class="header-right">
      <!-- 主题切换 -->
      <el-button
        class="header-action"
        text
        @click="toggleTheme"
      >
        <el-icon :size="18">
          <Sunny v-if="isDark" />
          <Moon v-else />
        </el-icon>
      </el-button>
      
      <!-- 通知 -->
      <el-popover placement="bottom-end" :width="360" trigger="click">
        <template #reference>
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="header-action">
            <el-button text>
              <el-icon :size="18">
                <Bell />
              </el-icon>
            </el-button>
          </el-badge>
        </template>
        
        <NotificationPanel @read="markAsRead" />
      </el-popover>
      
      <!-- 用户菜单 -->
      <el-dropdown trigger="click" @command="handleUserCommand">
        <div class="user-info">
          <el-avatar :size="32" :src="userAvatar" class="user-avatar">
            <el-icon>
              <User />
            </el-icon>
          </el-avatar>
          <span v-if="!isMobile" class="user-name">{{ username }}</span>
          <el-icon class="dropdown-icon">
            <ArrowDown />
          </el-icon>
        </div>
        
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item command="wallet">
              <el-icon><Wallet /></el-icon>
              我的钱包
            </el-dropdown-item>
            <el-dropdown-item command="tasks">
              <el-icon><List /></el-icon>
              我的任务
            </el-dropdown-item>
            <el-dropdown-item command="messages">
              <el-icon><Message /></el-icon>
              消息中心
            </el-dropdown-item>
            <el-dropdown-item divided command="settings">
              <el-icon><Setting /></el-icon>
              设置
            </el-dropdown-item>
            <el-dropdown-item command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import {
  Menu,
  Search,
  Sunny,
  Moon,
  Bell,
  User,
  Wallet,
  List,
  Message,
  Setting,
  SwitchButton,
  ArrowDown
} from '@element-plus/icons-vue'
import NotificationPanel from './NotificationPanel.vue'

const emit = defineEmits<{
  toggleSidebar: []
}>()

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

const searchKeyword = ref('')
const unreadCount = ref(3) // 模拟未读消息数量

const isMobile = computed(() => appStore.isMobile)
const isDark = computed(() => appStore.isDarkTheme)
const username = computed(() => userStore.nickname || userStore.username)
const userAvatar = computed(() => userStore.avatar)

// 切换主题
const toggleTheme = () => {
  appStore.toggleTheme()
}

// 处理搜索
const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({
      path: '/tasks',
      query: { keyword: searchKeyword.value.trim() }
    })
  }
}

// 处理用户菜单命令
const handleUserCommand = (command: string) => {
  const commands: Record<string, () => void> = {
    profile: () => router.push('/profile'),
    wallet: () => router.push('/wallet'),
    tasks: () => router.push('/my-tasks'),
    messages: () => router.push('/messages'),
    settings: () => router.push('/settings'),
    logout: () => {
      userStore.logout()
      router.push('/login')
    }
  }
  
  const handler = commands[command]
  if (handler) {
    handler()
  }
}

// 标记消息为已读
const markAsRead = () => {
  unreadCount.value = 0
}
</script>

<style lang="scss" scoped>
.layout-header {
  height: $header-height;
  display: flex;
  align-items: center;
  padding: 0 $spacing-lg;
  background: $bg-white;
  border-bottom: 1px solid $border-lighter;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: $z-index-sticky;
}

.header-left {
  display: flex;
  align-items: center;
  flex: 1;
}

.sidebar-toggle {
  margin-right: $spacing-sm;
  
  @media (min-width: $breakpoint-md) {
    display: none;
  }
}

.logo {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: $text-primary;
  font-weight: $font-weight-primary;
  font-size: $font-size-medium;
}

.logo-img {
  height: 32px;
  width: 32px;
  margin-right: $spacing-sm;
}

.logo-text {
  font-size: $font-size-large;
  font-weight: $font-weight-primary;
  color: $primary-color;
}

.header-center {
  flex: 2;
  max-width: 400px;
  margin: 0 $spacing-lg;
  
  @media (max-width: $breakpoint-md) {
    display: none;
  }
}

.search-input {
  width: 100%;
}

.header-right {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  flex: 1;
  justify-content: flex-end;
}

.header-action {
  position: relative;
  color: $text-regular;
  
  &:hover {
    color: $primary-color;
  }
}

.user-info {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  cursor: pointer;
  padding: $spacing-xs $spacing-sm;
  border-radius: $border-radius-base;
  transition: background-color 0.2s;
  
  &:hover {
    background-color: $bg-page;
  }
}

.user-avatar {
  flex-shrink: 0;
}

.user-name {
  font-size: $font-size-sm;
  color: $text-primary;
  max-width: 100px;
  @extend .text-ellipsis;
}

.dropdown-icon {
  font-size: 12px;
  color: $text-secondary;
  transition: transform 0.2s;
}

.user-info:hover .dropdown-icon {
  transform: rotate(180deg);
}

// 响应式调整
@media (max-width: $breakpoint-sm) {
  .layout-header {
    padding: 0 $spacing-md;
  }
  
  .header-left {
    flex: 1;
  }
  
  .header-right {
    gap: $spacing-xs;
  }
}

// 暗色主题
.dark .layout-header {
  background: var(--el-bg-color-overlay);
  border-bottom-color: var(--el-border-color);
  
  .logo {
    color: var(--el-text-color-primary);
  }
  
  .logo-text {
    color: $primary-color;
  }
  
  .header-action {
    color: var(--el-text-color-secondary);
    
    &:hover {
      color: $primary-color;
    }
  }
  
  .user-info {
    &:hover {
      background-color: var(--el-bg-color);
    }
  }
  
  .user-name {
    color: var(--el-text-color-primary);
  }
}
</style>