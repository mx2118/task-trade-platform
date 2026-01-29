<template>
  <aside
    class="layout-sidebar"
    :class="{ 'sidebar-collapsed': collapsed }"
    @mouseleave="handleMouseLeave"
  >
    <div class="sidebar-content">
      <!-- 菜单 -->
      <el-menu
        :default-active="activeMenu"
        :collapse="collapsed"
        :unique-opened="true"
        router
        class="sidebar-menu"
      >
        <!-- 任务中心 -->
        <el-sub-menu index="tasks">
          <template #title>
            <el-icon>
              <List />
            </el-icon>
            <span>任务中心</span>
          </template>
          
          <el-menu-item index="/tasks">
            <el-icon>
              <Search />
            </el-icon>
            <template #title>任务大厅</template>
          </el-menu-item>
          
          <el-menu-item 
            v-for="category in categories" 
            :key="category.id"
            :index="`/tasks?category_id=${category.id}`" 
            class="sub-item"
          >
            <el-icon>
              <Folder />
            </el-icon>
            <template #title>{{ category.name }}</template>
          </el-menu-item>
          
          <el-menu-item index="/publish" v-if="isLoggedIn">
            <el-icon>
              <Plus />
            </el-icon>
            <template #title>发布任务</template>
          </el-menu-item>
          
          <el-menu-item index="/my-tasks" v-if="isLoggedIn">
            <el-icon>
              <User />
            </el-icon>
            <template #title>我的任务</template>
          </el-menu-item>
        </el-sub-menu>
        
        <!-- 钱包 -->
        <el-menu-item index="/wallet" v-if="isLoggedIn">
          <el-icon>
            <Wallet />
          </el-icon>
          <template #title>我的钱包</template>
        </el-menu-item>
        
        <!-- 个人中心 -->
        <el-menu-item index="/profile" v-if="isLoggedIn">
          <el-icon>
            <UserFilled />
          </el-icon>
          <template #title>个人中心</template>
        </el-menu-item>
        
        <!-- 消息中心 -->
        <el-menu-item index="/messages" v-if="isLoggedIn">
          <el-icon>
            <Message />
          </el-icon>
          <template #title>消息中心</template>
        </el-menu-item>
        
        <!-- 帮助中心 -->
        <el-menu-item index="/help">
          <el-icon>
            <QuestionFilled />
          </el-icon>
          <template #title>帮助中心</template>
        </el-menu-item>
        
        <!-- 关于我们 -->
        <el-menu-item index="/about">
          <el-icon>
            <InfoFilled />
          </el-icon>
          <template #title>关于我们</template>
        </el-menu-item>
      </el-menu>
      
      <!-- 用户信息（收起时显示） -->
      <div v-if="collapsed && isLoggedIn" class="sidebar-user-collapsed">
        <el-tooltip content="个人中心" placement="right">
          <el-avatar :size="32" :src="userAvatar" class="user-avatar">
            <el-icon>
              <User />
            </el-icon>
          </el-avatar>
        </el-tooltip>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { categoryApi } from '@/api'
import type { Category } from '@/types'
import {
  List,
  Search,
  Plus,
  User,
  Wallet,
  UserFilled,
  Message,
  QuestionFilled,
  InfoFilled,
  Clock,
  Loading,
  CircleCheck,
  Document,
  View,
  CircleClose,
  Folder
} from '@element-plus/icons-vue'

interface Props {
  collapsed: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
}>()

const route = useRoute()
const appStore = useAppStore()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const isLoggedIn = computed(() => userStore.isLoggedIn)
const userAvatar = computed(() => userStore.avatar || '')

// 分类列表
const categories = ref<Category[]>([])

// 加载分类
const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategories()
    const responseData = response.data?.data || response.data || {}
    categories.value = responseData.categories || responseData.list || []
  } catch (error: any) {
    console.error('加载分类失败:', error)
    categories.value = []
  }
}

// 鼠标离开处理（仅桌面端）
const handleMouseLeave = () => {
  if (!appStore.isMobile && !props.collapsed) {
    // 可以添加自动收起逻辑
  }
}

// 处理菜单项点击
const handleMenuClick = () => {
  if (appStore.isMobile) {
    emit('close')
  }
}

// 快捷键支持
const handleKeydown = (e: KeyboardEvent) => {
  // Ctrl/Cmd + K 快速搜索
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    // 触发搜索
  }
  
  // Ctrl/Cmd + / 显示快捷键帮助
  if ((e.ctrlKey || e.metaKey) && e.key === '/') {
    e.preventDefault()
    // 显示帮助
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
  loadCategories()
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style lang="scss" scoped>
.layout-sidebar {
  position: fixed;
  top: $header-height;
  left: 0;
  height: calc(100vh - #{$header-height});
  width: $sidebar-width;
  background: $bg-white;
  border-right: 1px solid $border-lighter;
  transition: width 0.3s ease;
  z-index: $z-index-fixed;
  overflow: hidden;
  
  &.sidebar-collapsed {
    width: $sidebar-collapsed-width;
  }
  
  @media (max-width: $breakpoint-md) {
    // 移动端：默认隐藏侧边栏
    transform: translateX(-100%);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    width: 280px; // 移动端固定宽度
    z-index: $z-index-modal;
    box-shadow: none;
    
    // 移动端：当未收起时显示侧边栏
    &:not(.sidebar-collapsed) {
      transform: translateX(0);
      box-shadow: 4px 0 24px rgba(0, 0, 0, 0.12);
    }
  }
}

.sidebar-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative;
}

.sidebar-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: $spacing-md;
  border-bottom: 1px solid $border-lighter;
  background: linear-gradient(135deg, rgba($primary-color, 0.03) 0%, rgba($primary-color, 0.01) 100%);
}

.logo-img {
  height: 32px;
  width: 32px;
  flex-shrink: 0;
}

.logo-text {
  margin-left: $spacing-sm;
  font-size: $font-size-large;
  font-weight: $font-weight-primary;
  color: $primary-color;
  white-space: nowrap;
}

.sidebar-menu {
  flex: 1;
  border-right: none;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px 0;
  
  // 自定义滚动条
  &::-webkit-scrollbar {
    width: 4px;
  }
  
  &::-webkit-scrollbar-track {
    background: transparent;
  }
  
  &::-webkit-scrollbar-thumb {
    background: $border-lighter;
    border-radius: 2px;
    
    &:hover {
      background: $border-light;
    }
  }
  
  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    height: 48px;
    line-height: 48px;
    margin: 2px 8px;
    border-radius: 8px;
    transition: all 0.2s ease;
    color: #374151;
    font-size: 15px;
    
    &:hover {
      background-color: rgba($primary-color, 0.08);
      color: $primary-color;
    }
    
    &.is-active {
      background: linear-gradient(135deg, $primary-color 0%, lighten($primary-color, 10%) 100%);
      color: white;
      font-weight: 500;
      box-shadow: 0 2px 8px rgba($primary-color, 0.25);
    }
    
    .el-icon {
      font-size: 18px;
      color: inherit;
    }
  }
  
  :deep(.el-sub-menu) {
    .el-sub-menu__title {
      color: #374151;
      font-weight: 500;
      
      &:hover {
        background-color: rgba($primary-color, 0.08);
        color: $primary-color;
      }
    }
    
    .el-menu-item {
      padding-left: 52px !important;
      color: #4b5563;
      
      &.sub-item {
        padding-left: 64px !important;
        height: 42px;
        line-height: 42px;
        font-size: 14px;
        
        .el-icon {
          font-size: 16px;
        }
      }
      
      &.is-active {
        background-color: rgba($primary-color, 0.12);
        color: $primary-color;
        font-weight: 500;
      }
      
      &:hover {
        background-color: rgba($primary-color, 0.08);
        color: $primary-color;
      }
    }
  }
}

.sidebar-user-collapsed {
  padding: $spacing-md;
  border-top: 1px solid $border-lighter;
  display: flex;
  justify-content: center;
}

.user-avatar {
  cursor: pointer;
  transition: transform 0.2s;
  
  &:hover {
    transform: scale(1.1);
  }
}

// Logo文字动画
.logo-text-enter-active,
.logo-text-leave-active {
  transition: opacity 0.2s, transform 0.2s;
}

.logo-text-enter-from,
.logo-text-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}

// 响应式
@media (max-width: $breakpoint-md) {
  .layout-sidebar {
    width: 280px;
    
    &.sidebar-collapsed {
      width: 280px;
    }
  }
  
  .sidebar-menu {
    .el-menu-item,
    .el-sub-menu__title {
      height: 52px;
      line-height: 52px;
    }
  }
}

// 暗色主题
.dark .layout-sidebar {
  background: var(--el-bg-color-overlay);
  border-right-color: var(--el-border-color);
  
  .sidebar-logo {
    border-bottom-color: var(--el-border-color);
    background: linear-gradient(135deg, rgba($primary-color, 0.08) 0%, rgba($primary-color, 0.03) 100%);
  }
  
  .logo-text {
    color: $primary-color;
  }
  
  .sidebar-menu {
    :deep(.el-menu-item),
    :deep(.el-sub-menu__title) {
      color: #d1d5db;
      
      &:hover {
        background-color: rgba($primary-color, 0.12);
        color: $primary-color;
      }
      
      &.is-active {
        background: linear-gradient(135deg, $primary-color 0%, lighten($primary-color, 10%) 100%);
        color: white;
        box-shadow: 0 2px 12px rgba($primary-color, 0.3);
      }
      
      .el-icon {
        color: inherit;
      }
    }
    
    :deep(.el-sub-menu) {
      .el-sub-menu__title {
        color: #d1d5db;
        
        &:hover {
          background-color: rgba($primary-color, 0.12);
          color: $primary-color;
        }
      }
      
      .el-menu-item {
        color: #9ca3af;
        
        &.is-active {
          background-color: rgba($primary-color, 0.15);
          color: $primary-color;
        }
        
        &:hover {
          background-color: rgba($primary-color, 0.12);
          color: $primary-color;
        }
      }
    }
  }
  
  .sidebar-user-collapsed {
    border-top-color: var(--el-border-color);
  }
}
</style>