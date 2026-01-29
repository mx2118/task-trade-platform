<template>
  <el-breadcrumb class="breadcrumb" separator="/">
    <el-breadcrumb-item
      v-for="item in breadcrumbs"
      :key="item.path"
      :to="item.path"
    >
      <el-icon v-if="item.icon" class="breadcrumb-icon">
        <component :is="item.icon" />
      </el-icon>
      {{ item.title }}
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  HomeFilled,
  List,
  UserFilled,
  Wallet,
  Message,
  QuestionFilled,
  InfoFilled,
  Plus
} from '@element-plus/icons-vue'

interface BreadcrumbItem {
  title: string
  path: string
  icon?: any
}

const route = useRoute()
const router = useRouter()

// 路由映射
const routeMap: Record<string, { title: string; icon?: any }> = {
  '/': { title: '首页', icon: HomeFilled },
  '/home': { title: '首页', icon: HomeFilled },
  '/layout': { title: '首页', icon: HomeFilled },
  '/layout/tasks': { title: '任务大厅', icon: List },
  '/layout/publish': { title: '发布任务', icon: Plus },
  '/layout/my-tasks': { title: '我的任务', icon: UserFilled },
  '/layout/wallet': { title: '我的钱包', icon: Wallet },
  '/layout/dashboard': { title: '控制台', icon: HomeFilled },
  '/layout/user': { title: '个人中心', icon: UserFilled },
  '/layout/payment': { title: '支付', icon: Wallet },
  '/layout/profile': { title: '个人中心', icon: UserFilled },
  '/layout/messages': { title: '消息中心', icon: Message },
  '/layout/help': { title: '帮助中心', icon: QuestionFilled },
  '/layout/about': { title: '关于我们', icon: InfoFilled },
  '/tasks': { title: '任务大厅', icon: List },
  '/publish': { title: '发布任务', icon: Plus },
  '/my-tasks': { title: '我的任务', icon: UserFilled },
  '/wallet': { title: '我的钱包', icon: Wallet },
  '/profile': { title: '个人中心', icon: UserFilled },
  '/messages': { title: '消息中心', icon: Message },
  '/help': { title: '帮助中心', icon: QuestionFilled },
  '/about': { title: '关于我们', icon: InfoFilled }
}

const breadcrumbs = computed(() => {
  const items: BreadcrumbItem[] = []
  
  // 如果当前路径是任务大厅，直接返回空数组（不显示面包屑）
  if (route.path === '/' || route.path === '/layout' || route.path === '/layout/tasks' || route.path === '/tasks') {
    return items
  }
  
  // 对于其他页面，始终添加首页
  items.push({
    title: '首页',
    path: '/layout/tasks',
    icon: HomeFilled
  })
  
  // 解析当前路径
  const pathSegments = route.path.split('/').filter(Boolean)
  
  // 构建路径
  let currentPath = ''
  pathSegments.forEach((segment, index) => {
    currentPath += `/${segment}`
    
    // 跳过 layout 段
    if (segment === 'layout') {
      return
    }
    
    // 处理动态路由（如任务详情）
    if (/^\d+$/.test(segment)) {
      const parentPath = currentPath.substring(0, currentPath.lastIndexOf('/'))
      
      if (parentPath.includes('/tasks')) {
        items.push({
          title: '任务详情',
          path: currentPath
        })
      } else if (parentPath.includes('/my-tasks')) {
        items.push({
          title: '任务详情',
          path: currentPath
        })
      }
    } else {
      const routeInfo = routeMap[currentPath]
      if (routeInfo && currentPath !== '/layout' && currentPath !== '/layout/tasks') {
        items.push({
          title: routeInfo.title,
          path: currentPath,
          icon: routeInfo.icon
        })
      }
    }
  })
  
  return items
})
</script>

<style lang="scss" scoped>
.breadcrumb {
  .el-breadcrumb__item {
    .el-breadcrumb__inner {
      display: flex;
      align-items: center;
      gap: $spacing-xs;
      color: $text-regular;
      font-size: $font-size-small;
      
      &:hover {
        color: $primary-color;
      }
      
      &.is-link {
        color: $text-regular;
        
        &:hover {
          color: $primary-color;
        }
      }
    }
    
    &:last-child {
      .el-breadcrumb__inner {
        color: $text-primary;
        font-weight: $font-weight-secondary;
      }
    }
  }
  
  .breadcrumb-icon {
    font-size: 14px;
  }
}

// 响应式
@media (max-width: $breakpoint-sm) {
  .breadcrumb {
    .el-breadcrumb__item {
      .el-breadcrumb__inner {
        font-size: $font-size-extra-small;
      }
    }
    
    .breadcrumb-icon {
      font-size: 12px;
    }
  }
}

// 暗色主题
.dark .breadcrumb {
  .el-breadcrumb__item {
    .el-breadcrumb__inner {
      color: var(--el-text-color-regular);
      
      &:hover {
        color: $primary-color;
      }
      
      &.is-link {
        color: var(--el-text-color-regular);
        
        &:hover {
          color: $primary-color;
        }
      }
    }
    
    &:last-child {
      .el-breadcrumb__inner {
        color: var(--el-text-color-primary);
      }
    }
  }
}
</style>