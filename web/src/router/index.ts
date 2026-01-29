import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

// 基础路由（无需权限）
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/layout/tasks',
    meta: {
      hidden: true
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: {
      title: '登录',
      hidden: true,
      requiresAuth: false
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: {
      title: '注册',
      hidden: true,
      requiresAuth: false
    }
  },
  {
    path: '/layout',
    name: 'Layout',
    component: () => import('@/layout/index.vue'),
    redirect: '/layout/tasks',
    children: [
      {
        path: 'tasks/:id',
        name: 'TaskDetail',
        component: () => import('@/views/tasks/Detail.vue'),
        meta: {
          title: '任务详情',
          hidden: true,
          requiresAuth: false
        },
        props: true
      },
      {
        path: 'tasks',
        name: 'Tasks',
        component: () => import('@/views/tasks/Index.vue'),
        meta: {
          title: '任务列表',
          icon: 'List',
          requiresAuth: false
        }
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Index.vue'),
        meta: {
          title: '控制台',
          icon: 'DataLine',
          requiresAuth: true
        }
      },
      {
        path: 'user',
        name: 'UserCenter',
        component: () => import('@/views/user/Index.vue'),
        meta: {
          title: '个人中心',
          icon: 'User',
          requiresAuth: true
        }
      },
      {
        path: 'payment',
        name: 'Payment',
        component: () => import('@/views/payment/Index.vue'),
        meta: {
          title: '支付',
          icon: 'Wallet',
          requiresAuth: true
        }
      }
    ]
  }
]

// 错误页面路由
export const errorRoutes: RouteRecordRaw[] = [
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
    meta: {
      hidden: true
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes: [...constantRoutes, ...errorRoutes],
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 任务交易平台`
  }
  
  // 直接放行，不检查权限
  next()
})

export default router