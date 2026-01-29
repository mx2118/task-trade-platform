import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

// 配置NProgress
NProgress.configure({
  showSpinner: false,
  minimum: 0.1,
  speed: 200,
  trickleSpeed: 200
})

// 基础路由（无需权限）
export const constantRoutes: RouteRecordRaw[] = [
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
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '页面不存在',
      hidden: true,
      requiresAuth: false
    }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layout/index.vue'),
    redirect: '/home',
    children: [
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/views/home/index.vue'),
        meta: {
          title: '首页',
          icon: 'HomeFilled',
          requiresAuth: false
        }
      },
      {
        path: 'tasks',
        name: 'Tasks',
        component: () => import('@/views/tasks/index.vue'),
        meta: {
          title: '任务大厅',
          icon: 'List',
          requiresAuth: false
        }
      },
      {
        path: 'tasks/:id',
        name: 'TaskDetail',
        component: () => import('@/views/tasks/detail.vue'),
        meta: {
          title: '任务详情',
          hidden: true,
          requiresAuth: false
        },
        props: true
      },
      {
        path: 'publish',
        name: 'PublishTask',
        component: () => import('@/views/tasks/publish.vue'),
        meta: {
          title: '发布任务',
          icon: 'Plus',
          requiresAuth: true
        }
      },
      {
        path: 'my-tasks',
        name: 'MyTasks',
        component: () => import('@/views/tasks/my-tasks.vue'),
        meta: {
          title: '我的任务',
          icon: 'User',
          requiresAuth: true
        }
      },
      {
        path: 'wallet',
        name: 'Wallet',
        component: () => import('@/views/wallet/index.vue'),
        meta: {
          title: '我的钱包',
          icon: 'Wallet',
          requiresAuth: true
        }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/user/profile.vue'),
        meta: {
          title: '个人中心',
          icon: 'UserFilled',
          requiresAuth: true
        }
      },
      {
        path: 'messages',
        name: 'Messages',
        component: () => import('@/views/messages/index.vue'),
        meta: {
          title: '消息中心',
          icon: 'Message',
          requiresAuth: true
        }
      },
      {
        path: 'help',
        name: 'Help',
        component: () => import('@/views/help/index.vue'),
        meta: {
          title: '帮助中心',
          icon: 'QuestionFilled',
          requiresAuth: false
        }
      },
      {
        path: 'about',
        name: 'About',
        component: () => import('@/views/about/index.vue'),
        meta: {
          title: '关于我们',
          icon: 'InfoFilled',
          requiresAuth: false
        }
      }
    ]
  }
]

// 错误页面路由
export const errorRoutes: RouteRecordRaw[] = [
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404',
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
  // 开始进度条
  NProgress.start()
  
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 任务交易平台`
  }
  
  // 获取用户信息
  const userStore = useUserStore()
  const hasToken = userStore.token
  
  if (to.meta.requiresAuth) {
    // 需要登录的页面
    if (hasToken) {
      if (userStore.userInfo) {
        // 已有用户信息，直接放行
        next()
      } else {
        try {
          // 获取用户信息
          await userStore.getUserInfo()
          next()
        } catch (error) {
          // 获取用户信息失败，清除token并跳转到登录页
          userStore.logout()
          next(`/login?redirect=${to.fullPath}`)
        }
      }
    } else {
      // 未登录，跳转到登录页
      next(`/login?redirect=${to.fullPath}`)
    }
  } else {
    // 不需要登录的页面
    if (to.path === '/login' && hasToken) {
      // 已登录用户访问登录页，跳转到首页
      next('/')
    } else {
      next()
    }
  }
})

router.afterEach(() => {
  // 结束进度条
  NProgress.done()
})

export default router