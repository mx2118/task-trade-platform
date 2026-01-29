import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './styles/index.scss'

import App from './App.vue'
import router from './router'

// 导入权限控制
import './permission'

// 创建应用实例
const app = createApp(App)

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 使用插件
app.use(createPinia())
app.use(router)
app.use(ElementPlus, {
  size: 'default',
  zIndex: 3000
})

// 全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('全局错误:', err)
  console.error('错误信息:', info)
  
  // 在生产环境中，可以发送错误信息到监控服务
  if (import.meta.env.PROD && import.meta.env.VITE_ENABLE_ERROR_MONITOR === 'true') {
    // 这里可以集成错误监控服务
  }
}

// 挂载应用
app.mount('#app')

// 开发环境下的调试工具
if (import.meta.env.DEV && import.meta.env.VITE_SHOW_DEV_TOOLS === 'true') {
  // 开发工具配置
  window.__VUE_DEVTOOLS_GLOBAL_HOOK__ = {
    on: () => {},
    off: () => {},
    emit: () => {}
  }
}