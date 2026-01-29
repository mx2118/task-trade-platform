# Vue3.5前端项目结构设计

## 1. 项目目录结构

```
task-platform-frontend/
├── public/                       # 静态资源
│   ├── favicon.ico
│   ├── index.html
│   └── manifest.json
├── src/                          # 源代码
│   ├── api/                      # API接口
│   │   ├── index.js             # API配置
│   │   ├── request.js           # 请求封装
│   │   ├── auth.js              # 认证相关
│   │   ├── user.js              # 用户相关
│   │   ├── task.js              # 任务相关
│   │   ├── order.js             # 订单相关
│   │   ├── payment.js           # 支付相关
│   │   └── upload.js            # 文件上传
│   ├── assets/                   # 静态资源
│   │   ├── images/              # 图片
│   │   ├── icons/               # 图标
│   │   └── styles/              # 样式文件
│   │       ├── index.scss       # 全局样式
│   │       ├── variables.scss   # 变量定义
│   │       ├── mixins.scss      # 混入
│   │       └── reset.scss        # 样式重置
│   ├── components/               # 公共组件
│   │   ├── common/              # 通用组件
│   │   │   ├── AppHeader.vue    # 头部导航
│   │   │   ├── AppFooter.vue    # 底部
│   │   │   ├── AppSidebar.vue   # 侧边栏
│   │   │   ├── AppLoading.vue   # 加载组件
│   │   │   ├── AppPagination.vue # 分页组件
│   │   │   └── AppEmpty.vue     # 空状态
│   │   ├── forms/               # 表单组件
│   │   │   ├── TaskForm.vue     # 任务表单
│   │   │   ├── LoginForm.vue    # 登录表单
│   │   │   ├── RegisterForm.vue # 注册表单
│   │   │   └── UserForm.vue     # 用户表单
│   │   ├── cards/               # 卡片组件
│   │   │   ├── TaskCard.vue     # 任务卡片
│   │   │   ├── OrderCard.vue    # 订单卡片
│   │   │   └── UserCard.vue     # 用户卡片
│   │   └── business/            # 业务组件
│   │       ├── TaskStatus.vue   # 任务状态
│   │       ├── UserAvatar.vue   # 用户头像
│   │       ├── FileUpload.vue   # 文件上传
│   │       └── RichText.vue     # 富文本编辑器
│   ├── composables/              # 组合式函数
│   │   ├── useAuth.js           # 认证逻辑
│   │   ├── useRequest.js        # 请求逻辑
│   │   ├── usePagination.js     # 分页逻辑
│   │   ├── useUpload.js         # 上传逻辑
│   │   ├── useWebSocket.js      # WebSocket
│   │   └── useLocalStorage.js   # 本地存储
│   ├── layouts/                  # 布局组件
│   │   ├── DefaultLayout.vue    # 默认布局
│   │   ├── AuthLayout.vue       # 认证布局
│   │   ├── DashboardLayout.vue  # 仪表板布局
│   │   └── MobileLayout.vue     # 移动端布局
│   ├── pages/                    # 页面组件
│   │   ├── auth/                # 认证页面
│   │   │   ├── Login.vue        # 登录页
│   │   │   ├── Register.vue     # 注册页
│   │   │   └── ForgotPassword.vue # 忘记密码
│   │   ├── home/                # 首页
│   │   │   └── Index.vue        # 首页
│   │   ├── tasks/               # 任务相关
│   │   │   ├── List.vue         # 任务列表
│   │   │   ├── Detail.vue       # 任务详情
│   │   │   ├── Create.vue       # 创建任务
│   │   │   └── Edit.vue         # 编辑任务
│   │   ├── user/                # 用户相关
│   │   │   ├── Profile.vue      # 个人资料
│   │   │   ├── Settings.vue     # 设置
│   │   │   ├── Skills.vue       # 技能管理
│   │   │   └── Wallet.vue       # 钱包
│   │   ├── orders/              # 订单相关
│   │   │   ├── List.vue         # 订单列表
│   │   │   ├── Detail.vue       # 订单详情
│   │   │   └── Create.vue       # 创建订单
│   │   ├── payment/             # 支付相关
│   │   │   ├── Pay.vue          # 支付页面
│   │   │   └── Withdraw.vue     # 提现页面
│   │   ├── reviews/             # 评价相关
│   │   │   ├── List.vue         # 评价列表
│   │   │   └── Create.vue       # 创建评价
│   │   ├── messages/            # 消息相关
│   │   │   ├── List.vue         # 消息列表
│   │   │   └── Detail.vue       # 消息详情
│   │   └── admin/               # 管理员页面
│   │       ├── Dashboard.vue    # 管理仪表板
│   │       ├── Users.vue        # 用户管理
│   │       ├── Tasks.vue        # 任务管理
│   │       └── Orders.vue       # 订单管理
│   ├── router/                   # 路由配置
│   │   ├── index.js             # 路由主文件
│   │   ├── guards.js            # 路由守卫
│   │   └── routes.js            # 路由定义
│   ├── stores/                   # 状态管理
│   │   ├── index.js             # Store入口
│   │   ├── auth.js              # 认证状态
│   │   ├── user.js              # 用户状态
│   │   ├── task.js              # 任务状态
│   │   ├── order.js             # 订单状态
│   │   ├── message.js           # 消息状态
│   │   └── app.js               # 应用状态
│   ├── utils/                    # 工具函数
│   │   ├── auth.js              # 认证工具
│   │   ├── storage.js           # 存储工具
│   │   ├── validator.js         # 验证工具
│   │   ├── format.js            # 格式化工具
│   │   ├── constants.js         # 常量定义
│   │   └── helpers.js           # 辅助函数
│   ├── directives/               # 自定义指令
│   │   ├── lazy.js              # 懒加载指令
│   │   ├── permission.js        # 权限指令
│   │   └── loading.js           # 加载指令
│   ├── plugins/                  # 插件配置
│   │   ├── element-plus.js      # Element Plus
│   │   ├── axios.js             # Axios
│   │   └── app.js               # 应用插件
│   ├── types/                    # TypeScript类型定义
│   │   ├── api.d.ts             # API类型
│   │   ├── user.d.ts            # 用户类型
│   │   ├── task.d.ts            # 任务类型
│   │   └── common.d.ts          # 通用类型
│   ├── App.vue                   # 根组件
│   └── main.js                   # 入口文件
├── tests/                        # 测试文件
│   ├── unit/                    # 单元测试
│   ├── e2e/                     # 端到端测试
│   └── fixtures/                # 测试数据
├── .env                          # 环境变量
├── .env.development              # 开发环境变量
├── .env.production               # 生产环境变量
├── .eslintrc.js                  # ESLint配置
├── .prettierrc                   # Prettier配置
├── tsconfig.json                 # TypeScript配置
├── vite.config.js               # Vite配置
├── package.json                  # 项目配置
└── README.md                     # 项目说明
```

## 2. 核心配置文件

### 2.1 Vite配置

#### vite.config.js
```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
      imports: ['vue', 'vue-router', 'pinia'],
      dts: true
    }),
    Components({
      resolvers: [ElementPlusResolver()],
      dts: true
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@components': resolve(__dirname, 'src/components'),
      '@pages': resolve(__dirname, 'src/pages'),
      '@api': resolve(__dirname, 'src/api'),
      '@utils': resolve(__dirname, 'src/utils'),
      '@stores': resolve(__dirname, 'src/stores'),
      '@assets': resolve(__dirname, 'src/assets')
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: '@import "@/assets/styles/variables.scss";'
      }
    }
  },
  server: {
    port: 3000,
    open: true,
    cors: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api/v1')
      }
    }
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          elementPlus: ['element-plus'],
          utils: ['axios', 'dayjs']
        }
      }
    }
  }
})
```

### 2.2 包管理配置

#### package.json
```json
{
  "name": "task-platform-frontend",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "preview": "vite preview",
    "lint": "eslint . --ext .vue,.js,.jsx,.cjs,.mjs,.ts,.tsx,.cts,.mts --fix --ignore-path .gitignore",
    "format": "prettier --write src/",
    "test": "vitest",
    "test:e2e": "cypress run",
    "test:e2e:dev": "cypress open"
  },
  "dependencies": {
    "vue": "^3.5.0",
    "vue-router": "^4.2.0",
    "pinia": "^2.1.0",
    "element-plus": "^2.4.0",
    "axios": "^1.6.0",
    "dayjs": "^1.11.0",
    "lodash-es": "^4.17.21",
    "nprogress": "^0.2.0",
    "js-cookie": "^3.0.5",
    "@element-plus/icons-vue": "^2.3.1",
    "echarts": "^5.4.0",
    "vue-echarts": "^6.6.0",
    "swiper": "^11.0.0",
    "qrcode": "^1.5.3"
  },
  "devDependencies": {
    "@types/node": "^20.8.0",
    "@types/lodash-es": "^4.17.0",
    "@types/js-cookie": "^3.0.0",
    "@types/qrcode": "^1.5.0",
    "@typescript-eslint/eslint-plugin": "^6.0.0",
    "@typescript-eslint/parser": "^6.0.0",
    "@vitejs/plugin-vue": "^4.4.0",
    "@vue/eslint-config-prettier": "^8.0.0",
    "@vue/eslint-config-typescript": "^12.0.0",
    "cypress": "^13.3.0",
    "eslint": "^8.45.0",
    "eslint-plugin-vue": "^9.17.0",
    "prettier": "^3.0.0",
    "sass": "^1.69.0",
    "typescript": "^5.2.0",
    "unplugin-auto-import": "^0.16.0",
    "unplugin-vue-components": "^0.25.0",
    "vite": "^4.4.0",
    "vitest": "^0.34.0",
    "vue-tsc": "^1.8.0"
  }
}
```

## 3. 核心代码实现

### 3.1 入口文件

#### src/main.js
```javascript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './assets/styles/index.scss'
import './plugins/element-plus'
import './plugins/axios'
import './directives'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
```

### 3.2 路由配置

#### src/router/index.js
```javascript
import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
import { useAuthStore } from '@/stores/auth'
import { getToken } from '@/utils/auth'

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const token = getToken()

  // 如果有token但用户信息为空，获取用户信息
  if (token && !authStore.userInfo) {
    try {
      await authStore.getUserInfo()
    } catch (error) {
      authStore.logout()
      next('/login')
      return
    }
  }

  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 任务平台`
  }

  // 权限检查
  if (to.meta.requiresAuth && !token) {
    next('/login')
    return
  }

  // 管理员权限检查
  if (to.meta.requiresAdmin && authStore.userInfo?.role !== 'admin') {
    next('/403')
    return
  }

  next()
})

export default router
```

#### src/router/routes.js
```javascript
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import DashboardLayout from '@/layouts/DashboardLayout.vue'

export default [
  {
    path: '/',
    component: DefaultLayout,
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/pages/home/Index.vue'),
        meta: { title: '首页' }
      },
      {
        path: '/tasks',
        name: 'TaskList',
        component: () => import('@/pages/tasks/List.vue'),
        meta: { title: '任务列表' }
      },
      {
        path: '/tasks/:id',
        name: 'TaskDetail',
        component: () => import('@/pages/tasks/Detail.vue'),
        meta: { title: '任务详情' }
      }
    ]
  },
  {
    path: '/auth',
    component: AuthLayout,
    children: [
      {
        path: '/login',
        name: 'Login',
        component: () => import('@/pages/auth/Login.vue'),
        meta: { title: '登录' }
      },
      {
        path: '/register',
        name: 'Register',
        component: () => import('@/pages/auth/Register.vue'),
        meta: { title: '注册' }
      }
    ]
  },
  {
    path: '/user',
    component: DashboardLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '/profile',
        name: 'UserProfile',
        component: () => import('@/pages/user/Profile.vue'),
        meta: { title: '个人资料' }
      },
      {
        path: '/tasks/published',
        name: 'PublishedTasks',
        component: () => import('@/pages/tasks/List.vue'),
        meta: { title: '我发布的任务' }
      },
      {
        path: '/tasks/applied',
        name: 'AppliedTasks',
        component: () => import('@/pages/tasks/List.vue'),
        meta: { title: '我申请的任务' }
      },
      {
        path: '/orders',
        name: 'Orders',
        component: () => import('@/pages/orders/List.vue'),
        meta: { title: '订单管理' }
      }
    ]
  },
  {
    path: '/admin',
    component: DashboardLayout,
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '/dashboard',
        name: 'AdminDashboard',
        component: () => import('@/pages/admin/Dashboard.vue'),
        meta: { title: '管理后台' }
      },
      {
        path: '/users',
        name: 'AdminUsers',
        component: () => import('@/pages/admin/Users.vue'),
        meta: { title: '用户管理' }
      }
    ]
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/pages/error/403.vue'),
    meta: { title: '禁止访问' }
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/pages/error/404.vue'),
    meta: { title: '页面不存在' }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]
```

### 3.3 状态管理

#### src/stores/auth.js
```javascript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, register, logout, getUserInfo } from '@/api/auth'
import { setToken, removeToken } from '@/utils/auth'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref(localStorage.getItem('token'))
  const userInfo = ref(null)
  const loading = ref(false)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')

  // 方法
  const loginAction = async (credentials) => {
    loading.value = true
    try {
      const response = await login(credentials)
      token.value = response.data.token
      userInfo.value = response.data.user
      setToken(response.data.token)
      return response
    } finally {
      loading.value = false
    }
  }

  const registerAction = async (userData) => {
    loading.value = true
    try {
      const response = await register(userData)
      return response
    } finally {
      loading.value = false
    }
  }

  const logoutAction = async () => {
    try {
      await logout()
    } finally {
      token.value = null
      userInfo.value = null
      removeToken()
    }
  }

  const getUserInfoAction = async () => {
    try {
      const response = await getUserInfo()
      userInfo.value = response.data
      return response
    } catch (error) {
      logoutAction()
      throw error
    }
  }

  return {
    token,
    userInfo,
    loading,
    isLoggedIn,
    isAdmin,
    loginAction,
    registerAction,
    logoutAction,
    getUserInfoAction
  }
})
```

### 3.4 API封装

#### src/api/request.js
```javascript
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, removeToken } from '@/utils/auth'
import router from '@/router'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const { data } = response
    
    // 如果返回的code不是200，说明有错误
    if (data.code !== 200) {
      ElMessage.error(data.message || '请求失败')
      return Promise.reject(new Error(data.message || '请求失败'))
    }
    
    return data
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          ElMessage.error('登录已过期，请重新登录')
          removeToken()
          router.push('/login')
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器错误')
          break
        default:
          ElMessage.error(data?.message || '网络错误')
      }
    } else {
      ElMessage.error('网络连接错误')
    }
    
    return Promise.reject(error)
  }
)

export default request
```

### 3.5 组合式函数

#### src/composables/useRequest.js
```javascript
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'

export function useRequest(requestFn, options = {}) {
  const { immediate = false, onSuccess, onError } = options
  
  const loading = ref(false)
  const data = ref(null)
  const error = ref(null)
  
  const execute = async (...args) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await requestFn(...args)
      data.value = response.data
      
      if (onSuccess) {
        onSuccess(response)
      }
      
      return response
    } catch (err) {
      error.value = err
      
      if (onError) {
        onError(err)
      }
      
      throw err
    } finally {
      loading.value = false
    }
  }
  
  if (immediate) {
    execute()
  }
  
  return {
    loading,
    data,
    error,
    execute
  }
}

export function usePagination(requestFn, options = {}) {
  const { defaultPageSize = 10 } = options
  
  const pagination = reactive({
    page: 1,
    pageSize: defaultPageSize,
    total: 0
  })
  
  const { loading, data, error, execute } = useRequest(
    (params) => requestFn({ ...params, page: pagination.page, pageSize: pagination.pageSize }),
    { immediate: true }
  )
  
  const setPage = (page) => {
    pagination.page = page
    execute()
  }
  
  const setPageSize = (pageSize) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    execute()
  }
  
  const refresh = () => {
    execute()
  }
  
  return {
    loading,
    data,
    error,
    pagination,
    setPage,
    setPageSize,
    refresh
  }
}
```

## 4. 样式设计

### 4.1 SCSS变量

#### src/assets/styles/variables.scss
```scss
// 主题色
$primary-color: #409eff;
$success-color: #67c23a;
$warning-color: #e6a23c;
$danger-color: #f56c6c;
$info-color: #909399;

// 文字色
$text-primary: #303133;
$text-regular: #606266;
$text-secondary: #909399;
$text-placeholder: #c0c4cc;

// 边框色
$border-color: #dcdfe6;
$border-light: #e4e7ed;
$border-lighter: #ebeef5;

// 背景色
$bg-color: #f5f7fa;
$bg-white: #ffffff;

// 间距
$spacing-xs: 4px;
$spacing-sm: 8px;
$spacing-md: 16px;
$spacing-lg: 24px;
$spacing-xl: 32px;

// 断点
$breakpoint-xs: 480px;
$breakpoint-sm: 768px;
$breakpoint-md: 992px;
$breakpoint-lg: 1200px;
$breakpoint-xl: 1920px;

// 阴影
$box-shadow-base: 0 2px 4px rgba(0, 0, 0, 0.12), 0 0 6px rgba(0, 0, 0, 0.04);
$box-shadow-light: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
$box-shadow-dark: 0 4px 16px rgba(0, 0, 0, 0.12);

// 圆角
$border-radius-sm: 4px;
$border-radius-base: 6px;
$border-radius-lg: 8px;
```

### 4.2 全局样式

#### src/assets/styles/index.scss
```scss
@import './variables.scss';
@import './mixins.scss';
@import './reset.scss';

// 全局样式
body {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5;
  color: $text-regular;
  background-color: $bg-color;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

// 通用类
.clearfix::after {
  content: '';
  display: table;
  clear: both;
}

.text-center {
  text-align: center;
}

.text-left {
  text-align: left;
}

.text-right {
  text-align: right;
}

.flex {
  display: flex;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.flex-column {
  display: flex;
  flex-direction: column;
}

// 响应式工具类
@media (max-width: $breakpoint-sm) {
  .hidden-sm {
    display: none !important;
  }
}

@media (max-width: $breakpoint-md) {
  .hidden-md {
    display: none !important;
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
```

## 5. 移动端适配

### 5.1 响应式设计
- 使用rem单位和flex布局
- 移动端专用组件
- 触摸事件优化
- 横竖屏适配

### 5.2 PWA支持
- Service Worker缓存
- 离线功能
- 推送通知
- 应用安装

这个前端架构设计采用了现代化的Vue3.5技术栈，具备良好的可维护性和扩展性，支持PC端和移动端的适配。