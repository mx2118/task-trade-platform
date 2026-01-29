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
  
  // 始终添加首页
  items.push({
    title: '首页',
    path: '/',
    icon: HomeFilled
  })
  
  // 解析当前路径
  const pathSegments = route.path.split('/').filter(Boolean)
  
  // 构建路径
  let currentPath = ''
  pathSegments.forEach((segment, index) => {
    currentPath += `/${segment}`
    
    // 处理动态路由（如任务详情）
    if (/^\d+$/.test(segment)) {
      const parentPath = currentPath.substring(0, currentPath.lastIndexOf('/'))
      const parentRoute = routeMap[parentPath]
      
      if (parentRoute && parentPath === '/tasks') {
        items.push({
          title: '任务详情',
          path: currentPath
        })
      } else if (parentRoute && parentPath === '/my-tasks') {
        items.push({
          title: '任务详情',
          path: currentPath
        })
      }
    } else {
      const routeInfo = routeMap[currentPath]
      if (routeInfo) {
        items.push({
          title: routeInfo.title,
          path: currentPath,
          icon: routeInfo.icon
        })
      }
    }
  })
  
  // 移除重复的首页（如果当前路径就是首页）
  if (items.length > 1 && items[0].path === route.path) {
    items.shift()
  }
  
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
      font-size: $font-size-sm;
      
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
        font-size: $font-size-xs;
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