<template>
  <div id="app" class="app-container">
    <!-- 路由视图 -->
    <router-view />
    
    <!-- 全局加载组件 -->
    <el-loading
      v-if="globalLoading"
      :lock="true"
      text="加载中..."
      background="rgba(0, 0, 0, 0.7)"
      :custom-class="'global-loading'"
    />
    
    <!-- 全局消息提示 -->
    <GlobalMessage />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import GlobalMessage from '@/components/GlobalMessage.vue'

// 获取应用状态
const appStore = useAppStore()

// 全局加载状态
const globalLoading = computed(() => appStore.globalLoading)
</script>

<style lang="scss">
.app-container {
  min-height: 100vh;
  width: 100vw;
  overflow-x: hidden;
  overflow-y: auto;
}

// 页面切换动画
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.2s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

// 全局加载样式
.global-loading {
  .el-loading-spinner {
    .circular {
      width: 50px;
      height: 50px;
    }
    
    .el-loading-text {
      font-size: 16px;
      color: #fff;
      margin-top: 16px;
    }
  }
}

// 滚动条样式
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

// 移动端适配
@media (max-width: 768px) {
  .app-container {
    font-size: 14px;
  }
}

// 暗色主题适配
.dark {
  .app-container {
    background-color: var(--el-bg-color-page);
    color: var(--el-text-color-primary);
  }
}
</style>