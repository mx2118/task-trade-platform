<template>
  <div class="layout-container">
    <!-- 头部导航 -->
    <Header @toggle-sidebar="toggleSidebar" />
    
    <!-- 主体内容 -->
    <div class="layout-main">
      <!-- 侧边栏 -->
      <Sidebar :collapsed="sidebarCollapsed" @close="closeSidebar" />
      
      <!-- 内容区域 -->
      <main class="layout-content" :class="{ 'content-collapsed': sidebarCollapsed }">
        <!-- 面包屑导航 -->
        <Breadcrumb class="breadcrumb" />
        
        <!-- 路由视图 -->
        <div class="page-container">
          <router-view v-slot="{ Component, route }">
            <transition name="page" mode="out-in">
              <component :is="Component" :key="route.path" />
            </transition>
          </router-view>
        </div>
      </main>
    </div>
    
    <!-- 移动端遮罩 -->
    <div
      v-if="isMobile && !sidebarCollapsed"
      class="mobile-overlay"
      @click="closeSidebar"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import Header from './components/Header.vue'
import Sidebar from './components/Sidebar.vue'
import Breadcrumb from './components/Breadcrumb.vue'

const appStore = useAppStore()

const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const isMobile = computed(() => appStore.isMobile)

const toggleSidebar = () => {
  appStore.toggleSidebar()
}

const closeSidebar = () => {
  if (isMobile.value) {
    appStore.setSidebarCollapsed(true)
  }
}
</script>

<style lang="scss" scoped>
.layout-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.layout-main {
  flex: 1;
  display: flex;
  position: relative;
  overflow: hidden;
}

.layout-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: margin-left 0.3s ease;
  margin-left: $sidebar-width;
  
  &.content-collapsed {
    margin-left: $sidebar-collapsed-width;
  }
}

.breadcrumb {
  padding: $spacing-md $spacing-lg;
  background: $bg-white;
  border-bottom: 1px solid $border-lighter;
  flex-shrink: 0;
}

.page-container {
  flex: 1;
  overflow-y: auto;
  padding: $spacing-lg;
}

.mobile-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: $z-index-modal - 1;
  
  @media (min-width: $breakpoint-md) {
    display: none;
  }
}

// 响应式布局
@media (max-width: $breakpoint-md) {
  .layout-content {
    margin-left: 0;
    
    &.content-collapsed {
      margin-left: 0;
    }
  }
  
  .page-container {
    padding: $spacing-md;
  }
  
  .breadcrumb {
    padding: $spacing-sm $spacing-md;
  }
}

// 暗色主题
.dark {
  .breadcrumb {
    background: var(--el-bg-color-overlay);
    border-bottom-color: var(--el-border-color);
  }
}

// 页面过渡动画
.page-enter-active {
  transition: all 0.4s ease;
}

.page-leave-active {
  transition: all 0.3s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>