import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import './styles/index.scss'

import App from './App.vue'
import router from './router'
import { useAppStore } from './stores/app'

// 创建应用实例
const app = createApp(App)

// 创建 Pinia 实例
const pinia = createPinia()

// 使用插件
app.use(pinia)
app.use(router)
app.use(ElementPlus)

// 注册所有图标组件
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 初始化应用设置
const appStore = useAppStore()
appStore.initAppSettings()

// 挂载应用
app.mount('#app')

// 隐藏加载动画
if (typeof window !== 'undefined' && window.__HIDE_LOADING__) {
  setTimeout(() => {
    window.__HIDE_LOADING__()
  }, 100)
}
