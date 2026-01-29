import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { ElMessage } from 'element-plus'

// 权限控制
export const setupPermission = () => {
  // 监听路由变化
  const userStore = useUserStore()
  const appStore = useAppStore()
  
  // 检查用户权限
  const checkPermission = (permission: string): boolean => {
    if (!userStore.userInfo) {
      return false
    }
    
    // 管理员拥有所有权限
    if (userStore.userInfo.role === 'admin') {
      return true
    }
    
    // 检查用户权限列表
    return userStore.userInfo.permissions?.includes(permission) || false
  }
  
  // 权限指令
  const hasPermission = (el: HTMLElement, binding: { value: string | string[] }) => {
    const { value } = binding
    
    if (!value) {
      return
    }
    
    const permissions = Array.isArray(value) ? value : [value]
    const hasAuth = permissions.some(permission => checkPermission(permission))
    
    if (!hasAuth) {
      el.remove()
    }
  }
  
  // 将权限检查函数添加到全局
  window.$checkPermission = checkPermission
  window.$hasPermission = hasPermission
  
  return {
    checkPermission,
    hasPermission
  }
}

// 路由权限控制
export const setupRouterPermission = () => {
  // 页面访问权限检查
  const checkPagePermission = (permission?: string): boolean => {
    const userStore = useUserStore()
    
    if (!userStore.userInfo) {
      return false
    }
    
    if (!permission) {
      return true
    }
    
    return setupPermission().checkPermission(permission)
  }
  
  return {
    checkPagePermission
  }
}

// API权限控制
export const setupApiPermission = () => {
  // 请求拦截器中的权限检查
  const checkApiPermission = (url: string, method: string): boolean => {
    const userStore = useUserStore()
    
    if (!userStore.userInfo) {
      return false
    }
    
    // 管理员拥有所有API访问权限
    if (userStore.userInfo.role === 'admin') {
      return true
    }
    
    // 根据URL和method检查权限
    // 这里可以根据实际需求实现更复杂的权限逻辑
    const apiPermissions = userStore.userInfo.apiPermissions || []
    const requiredPermission = `${method.toUpperCase()}:${url}`
    
    return apiPermissions.includes(requiredPermission) || apiPermissions.includes('*')
  }
  
  return {
    checkApiPermission
  }
}

// 按钮权限控制
export const setupButtonPermission = () => {
  const hasButtonPermission = (permission: string): boolean => {
    const { checkPermission } = setupPermission()
    return checkPermission(permission)
  }
  
  return {
    hasButtonPermission
  }
}

// 数据权限控制
export const setupDataPermission = () => {
  const checkDataPermission = (dataType: string, dataId?: string | number): boolean => {
    const userStore = useUserStore()
    
    if (!userStore.userInfo) {
      return false
    }
    
    // 管理员拥有所有数据权限
    if (userStore.userInfo.role === 'admin') {
      return true
    }
    
    // 根据数据类型和ID检查权限
    switch (dataType) {
      case 'task':
        // 检查任务权限（创建者、接取者或有相关权限）
        return true // 这里需要根据实际业务逻辑实现
      case 'user':
        // 检查用户数据权限（只能查看自己的信息）
        return true
      case 'order':
        // 检查订单权限
        return true
      default:
        return false
    }
  }
  
  return {
    checkDataPermission
  }
}

// 功能权限控制
export const setupFeaturePermission = () => {
  const checkFeaturePermission = (feature: string): boolean => {
    const userStore = useUserStore()
    const appStore = useAppStore()
    
    if (!userStore.userInfo) {
      return false
    }
    
    // 检查功能是否开启
    if (!appStore.features[feature]) {
      return false
    }
    
    // 管理员拥有所有功能权限
    if (userStore.userInfo.role === 'admin') {
      return true
    }
    
    // 检查用户功能权限
    const userFeatures = userStore.userInfo.features || []
    return userFeatures.includes(feature) || userFeatures.includes('*')
  }
  
  return {
    checkFeaturePermission
  }
}

// 权限错误处理
export const handlePermissionError = (error: any, type: 'page' | 'api' | 'button' | 'data' | 'feature') => {
  const userStore = useUserStore()
  
  let message = '权限不足'
  
  switch (type) {
    case 'page':
      message = '您没有访问此页面的权限'
      break
    case 'api':
      message = '您没有执行此操作的权限'
      break
    case 'button':
      message = '您没有使用此功能的权限'
      break
    case 'data':
      message = '您没有访问此数据的权限'
      break
    case 'feature':
      message = '此功能暂未开启或您没有使用权限'
      break
  }
  
  ElMessage.error(message)
  
  // 如果是页面权限错误，可以跳转到404页面
  if (type === 'page') {
    // router.push('/404')
  }
}

// 初始化权限系统
export const initPermissionSystem = () => {
  setupPermission()
  setupRouterPermission()
  setupApiPermission()
  setupButtonPermission()
  setupDataPermission()
  setupFeaturePermission()
}