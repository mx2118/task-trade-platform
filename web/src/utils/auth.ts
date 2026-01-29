import Cookies from 'js-cookie'

const TOKEN_KEY = 'task-platform-token'
const REFRESH_TOKEN_KEY = 'task-platform-refresh-token'

// Token相关操作
export const getToken = (): string | null => {
  return Cookies.get(TOKEN_KEY) || localStorage.getItem(TOKEN_KEY)
}

export const setToken = (token: string, expiresIn?: number): void => {
  // 保存到localStorage
  localStorage.setItem(TOKEN_KEY, token)
  
  // 保存到cookie
  const expires = expiresIn ? expiresIn / 86400 : 7 // 默认7天
  Cookies.set(TOKEN_KEY, token, { 
    expires,
    secure: import.meta.env.PROD,
    sameSite: 'strict'
  })
}

export const removeToken = (): void => {
  localStorage.removeItem(TOKEN_KEY)
  Cookies.remove(TOKEN_KEY)
}

// 刷新Token
export const getRefreshToken = (): string | null => {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
}

export const setRefreshToken = (refreshToken: string): void => {
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
}

export const removeRefreshToken = (): void => {
  localStorage.removeItem(REFRESH_TOKEN_KEY)
}

// 检查Token是否过期
export const isTokenExpired = (token: string): boolean => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    const currentTime = Date.now() / 1000
    return payload.exp < currentTime
  } catch {
    return true
  }
}

// 解析Token
export const parseToken = (token: string): any => {
  try {
    return JSON.parse(atob(token.split('.')[1]))
  } catch {
    return null
  }
}

// 检查用户权限
export const hasPermission = (userPermissions: string[], requiredPermission: string): boolean => {
  return userPermissions.includes(requiredPermission) || userPermissions.includes('*')
}

// 检查用户角色
export const hasRole = (userRoles: string[], requiredRole: string): boolean => {
  return userRoles.includes(requiredRole) || userRoles.includes('admin')
}

// 格式化用户名显示
export const formatUsername = (userInfo: { username?: string; nickname?: string }): string => {
  return userInfo.nickname || userInfo.username || '未知用户'
}

// 隐藏手机号
export const hidePhone = (phone: string): string => {
  if (!phone || phone.length < 7) return phone
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 隐藏邮箱
export const hideEmail = (email: string): string => {
  if (!email) return email
  const [username, domain] = email.split('@')
  if (username.length <= 2) return email
  
  const visiblePart = username.substring(0, 2)
  const hiddenPart = '*'.repeat(username.length - 2)
  return `${visiblePart}${hiddenPart}@${domain}`
}

// 生成随机字符串
export const generateRandomString = (length: number): string => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

// 生成UUID
export const generateUUID = (): string => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

// 设备指纹生成
export const generateDeviceFingerprint = async (): Promise<string> {
  const data = {
    userAgent: navigator.userAgent,
    language: navigator.language,
    platform: navigator.platform,
    cookieEnabled: navigator.cookieEnabled,
    screen: {
      width: screen.width,
      height: screen.height,
      colorDepth: screen.colorDepth
    },
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
    localStorage: !!window.localStorage,
    sessionStorage: !!window.sessionStorage
  }

  // 生成简单的hash
  const jsonString = JSON.stringify(data)
  let hash = 0
  for (let i = 0; i < jsonString.length; i++) {
    const char = jsonString.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash // Convert to 32-bit integer
  }
  
  return Math.abs(hash).toString(36)
}

// 验证密码强度
export const validatePassword = (password: string): {
  isValid: boolean
  strength: 'weak' | 'medium' | 'strong'
  message: string
} => {
  const minLength = 8
  const hasUpperCase = /[A-Z]/.test(password)
  const hasLowerCase = /[a-z]/.test(password)
  const hasNumbers = /\d/.test(password)
  const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(password)
  
  const score = [hasUpperCase, hasLowerCase, hasNumbers, hasSpecialChar].filter(Boolean).length
  
  if (password.length < minLength) {
    return {
      isValid: false,
      strength: 'weak',
      message: '密码长度至少8位'
    }
  }
  
  if (score < 2) {
    return {
      isValid: false,
      strength: 'weak',
      message: '密码强度太弱，建议包含大小写字母、数字和特殊字符'
    }
  }
  
  if (score === 2) {
    return {
      isValid: true,
      strength: 'medium',
      message: '密码强度中等'
    }
  }
  
  return {
    isValid: true,
    strength: 'strong',
    message: '密码强度很好'
  }
}

// 格式化文件大小
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 下载文件
export const downloadFile = (url: string, filename?: string): void => {
  const link = document.createElement('a')
  link.href = url
  link.download = filename || 'download'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 复制到剪贴板
export const copyToClipboard = async (text: string): Promise<boolean> => {
  try {
    if (navigator.clipboard) {
      await navigator.clipboard.writeText(text)
      return true
    } else {
      // 降级方案
      const textArea = document.createElement('textarea')
      textArea.value = text
      document.body.appendChild(textArea)
      textArea.select()
      const result = document.execCommand('copy')
      document.body.removeChild(textArea)
      return result
    }
  } catch {
    return false
  }
}

// 防抖函数
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timeoutId: ReturnType<typeof setTimeout>
  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => func(...args), delay)
  }
}

// 节流函数
export const throttle = <T extends (...args: any[]) => any>(
  func: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let lastCall = 0
  return (...args: Parameters<T>) => {
    const now = Date.now()
    if (now - lastCall >= delay) {
      lastCall = now
      func(...args)
    }
  }
}