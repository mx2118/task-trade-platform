import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ElNotification } from 'element-plus'

export const useAppStore = defineStore('app', () => {
  // 状态
  const sidebarCollapsed = ref(true) // 默认收起，在 initAppSettings 中根据设备类型调整
  const device = ref<'desktop' | 'mobile'>('desktop')
  const theme = ref<'light' | 'dark'>('light')
  const language = ref('zh-CN')
  const globalLoading = ref(false)
  const features = ref<Record<string, boolean>>({
    taskPublish: true,
    taskTake: true,
    payment: true,
    chat: false,
    live: false,
    vip: false
  })
  
  // 系统配置
  const config = ref({
    version: '1.0.0',
    buildTime: '',
    environment: import.meta.env.MODE,
    apiVersion: 'v1',
    wsUrl: import.meta.env.VITE_WS_URL,
    uploadUrl: import.meta.env.VITE_UPLOAD_URL
  })
  
  // 通知配置
  const notification = ref({
    enable: true,
    sound: true,
    desktop: true,
    email: false
  })
  
  // 计算属性
  const isMobile = computed(() => device.value === 'mobile')
  const isDesktop = computed(() => device.value === 'desktop')
  const isDarkTheme = computed(() => theme.value === 'dark')
  
  // Actions
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
    // 只在桌面端保存侧边栏状态
    if (device.value === 'desktop') {
      localStorage.setItem('sidebarCollapsed', JSON.stringify(sidebarCollapsed.value))
    }
  }
  
  const setSidebarCollapsed = (collapsed: boolean) => {
    sidebarCollapsed.value = collapsed
    // 只在桌面端保存侧边栏状态
    if (device.value === 'desktop') {
      localStorage.setItem('sidebarCollapsed', JSON.stringify(collapsed))
    }
  }
  
  const setDevice = (type: 'desktop' | 'mobile') => {
    device.value = type
    // 移动端自动收起侧边栏
    if (type === 'mobile') {
      sidebarCollapsed.value = true
    }
  }
  
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    // 应用主题到DOM
    document.documentElement.classList.toggle('dark', theme.value === 'dark')
    // 保存主题设置到localStorage
    localStorage.setItem('theme', theme.value)
  }
  
  const setTheme = (newTheme: 'light' | 'dark') => {
    theme.value = newTheme
    document.documentElement.classList.toggle('dark', newTheme === 'dark')
    localStorage.setItem('theme', newTheme)
  }
  
  const setLanguage = (lang: string) => {
    language.value = lang
    localStorage.setItem('language', lang)
  }
  
  const setGlobalLoading = (loading: boolean) => {
    globalLoading.value = loading
  }
  
  const updateFeatures = (newFeatures: Record<string, boolean>) => {
    features.value = { ...features.value, ...newFeatures }
  }
  
  const updateConfig = (newConfig: Partial<typeof config.value>) => {
    config.value = { ...config.value, ...newConfig }
  }
  
  const updateNotification = (newNotification: Partial<typeof notification.value>) => {
    notification.value = { ...notification.value, ...newNotification }
  }
  
  // 显示系统通知
  const showNotification = (title: string, message: string, type: 'success' | 'warning' | 'info' | 'error' = 'info') => {
    if (notification.value.enable) {
      ElNotification({
        title,
        message,
        type,
        duration: 3000
      })
    }
  }
  
  // 检查功能是否开启
  const isFeatureEnabled = (feature: string) => {
    return features.value[feature] || false
  }
  
  // 初始化应用设置
  const initAppSettings = () => {
    // 检测设备类型
    const checkDevice = () => {
      const width = window.innerWidth
      if (width < 768) {
        setDevice('mobile')
      } else {
        setDevice('desktop')
      }
    }
    
    // 先检测设备类型（这会自动设置移动端的侧边栏状态）
    checkDevice()
    
    // 只在桌面端恢复侧边栏状态
    if (device.value === 'desktop') {
      const savedSidebarCollapsed = localStorage.getItem('sidebarCollapsed')
      if (savedSidebarCollapsed !== null) {
        setSidebarCollapsed(JSON.parse(savedSidebarCollapsed))
      } else {
        // 桌面端默认展开
        setSidebarCollapsed(false)
      }
    }
    
    // 从localStorage恢复主题设置，如果没有则默认使用日间模式
    const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null
    if (savedTheme) {
      setTheme(savedTheme)
    } else {
      // 确保默认是日间模式
      setTheme('light')
    }
    
    // 从localStorage恢复语言设置
    const savedLanguage = localStorage.getItem('language')
    if (savedLanguage) {
      setLanguage(savedLanguage)
    }
    
    // 监听窗口大小变化
    window.addEventListener('resize', checkDevice)
  }
  
  // 重置应用状态
  const resetAppState = () => {
    sidebarCollapsed.value = false
    device.value = 'desktop'
    theme.value = 'light'
    language.value = 'zh-CN'
    globalLoading.value = false
    features.value = {
      taskPublish: true,
      taskTake: true,
      payment: true,
      chat: false,
      live: false,
      vip: false
    }
  }
  
  return {
    // 状态
    sidebarCollapsed,
    device,
    theme,
    language,
    globalLoading,
    features,
    config,
    notification,
    
    // 计算属性
    isMobile,
    isDesktop,
    isDarkTheme,
    
    // Actions
    toggleSidebar,
    setSidebarCollapsed,
    setDevice,
    toggleTheme,
    setTheme,
    setLanguage,
    setGlobalLoading,
    updateFeatures,
    updateConfig,
    updateNotification,
    showNotification,
    isFeatureEnabled,
    initAppSettings,
    resetAppState
  }
})