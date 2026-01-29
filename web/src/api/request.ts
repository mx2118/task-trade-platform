import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getToken } from '@/utils/auth'
import router from '@/router'

// 响应数据类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  timestamp?: number
}

// 请求配置
export interface RequestConfig extends AxiosRequestConfig {
  showLoading?: boolean
  showError?: boolean
  skipAuth?: boolean
}

// 创建axios实例
const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8'
  }
})

// 请求拦截器
service.interceptors.request.use(
  (config: any) => {
    // 获取用户token
    const userStore = useUserStore()
    const token = getToken()
    
    // 添加token到请求头
    if (token && !config.skipAuth) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    // 添加请求ID用于追踪
    config.headers['X-Request-ID'] = generateRequestId()
    
    // 添加设备信息
    config.headers['X-User-Agent'] = navigator.userAgent
    config.headers['X-Platform'] = getPlatform()
    
    // 处理请求数据
    if (config.data && config.method === 'get') {
      config.params = config.data
      config.data = undefined
    }
    
    return config
  },
  (error) => {
    console.error('请求拦截器错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { data, config } = response
    
    // 请求成功
    if (data.code === 200) {
      return data
    }
    
    // 业务错误处理
    handleBusinessError(data, config)
    return Promise.reject(new Error(data.message || '请求失败'))
  },
  (error) => {
    const { response, config } = error
    
    // 网络错误处理
    if (!response) {
      handleNetworkError(error, config)
      return Promise.reject(error)
    }
    
    // HTTP状态码错误处理
    handleHttpError(response, config)
    return Promise.reject(error)
  }
)

// 业务错误处理
const handleBusinessError = (data: ApiResponse, config: RequestConfig) => {
  const message = data.message || '业务错误'
  
  // 显示错误消息
  if (config.showError !== false) {
    ElMessage.error(message)
  }
  
  // 特殊业务错误码处理
  switch (data.code) {
    case 401:
      // 未授权，跳转到登录页
      const userStore = useUserStore()
      userStore.logout()
      router.push(`/login?redirect=${router.currentRoute.value.fullPath}`)
      break
    case 403:
      // 权限不足
      ElMessage.error('权限不足，请联系管理员')
      break
    case 404:
      // 资源不存在
      ElMessage.error('请求的资源不存在')
      break
    case 429:
      // 请求过于频繁
      ElMessage.error('请求过于频繁，请稍后再试')
      break
    case 500:
      // 服务器内部错误
      ElMessage.error('服务器内部错误，请稍后再试')
      break
  }
}

// 网络错误处理
const handleNetworkError = (error: any, config: RequestConfig) => {
  let message = '网络错误，请检查网络连接'
  
  if (error.code === 'ECONNABORTED') {
    message = '请求超时，请稍后再试'
  } else if (error.message.includes('Network Error')) {
    message = '网络连接失败，请检查网络设置'
  }
  
  if (config.showError !== false) {
    ElMessage.error(message)
  }
}

// HTTP状态码错误处理
const handleHttpError = (response: AxiosResponse, config: RequestConfig) => {
  const { status, statusText } = response
  let message = `请求失败 (${status})`
  
  switch (status) {
    case 400:
      message = '请求参数错误'
      break
    case 401:
      message = '未授权，请重新登录'
      break
    case 403:
      message = '权限不足'
      break
    case 404:
      message = '请求的资源不存在'
      break
    case 405:
      message = '请求方法不被允许'
      break
    case 408:
      message = '请求超时'
      break
    case 429:
      message = '请求过于频繁'
      break
    case 500:
      message = '服务器内部错误'
      break
    case 502:
      message = '网关错误'
      break
    case 503:
      message = '服务暂时不可用'
      break
    case 504:
      message = '网关超时'
      break
    default:
      message = `${statusText || '请求失败'}`
  }
  
  if (config.showError !== false) {
    ElMessage.error(message)
  }
  
  // 401错误特殊处理
  if (status === 401) {
    const userStore = useUserStore()
    userStore.logout()
    router.push(`/login?redirect=${router.currentRoute.value.fullPath}`)
  }
}

// 生成请求ID
const generateRequestId = (): string => {
  return Date.now().toString(36) + Math.random().toString(36).substr(2)
}

// 获取平台信息
const getPlatform = (): string => {
  const userAgent = navigator.userAgent
  if (userAgent.includes('Android')) {
    return 'Android'
  } else if (userAgent.includes('iPhone') || userAgent.includes('iPad')) {
    return 'iOS'
  } else if (userAgent.includes('Windows')) {
    return 'Windows'
  } else if (userAgent.includes('Mac')) {
    return 'macOS'
  } else if (userAgent.includes('Linux')) {
    return 'Linux'
  }
  return 'Unknown'
}

// 通用请求方法
export const request = <T = any>(config: RequestConfig): Promise<ApiResponse<T>> => {
  return service.request(config) as Promise<ApiResponse<T>>
}

// GET请求
export const get = <T = any>(
  url: string,
  params?: any,
  config: RequestConfig = {}
): Promise<ApiResponse<T>> => {
  return request<T>({
    method: 'get',
    url,
    params,
    ...config
  })
}

// POST请求
export const post = <T = any>(
  url: string,
  data?: any,
  config: RequestConfig = {}
): Promise<ApiResponse<T>> => {
  return request<T>({
    method: 'post',
    url,
    data,
    ...config
  })
}

// PUT请求
export const put = <T = any>(
  url: string,
  data?: any,
  config: RequestConfig = {}
): Promise<ApiResponse<T>> => {
  return request<T>({
    method: 'put',
    url,
    data,
    ...config
  })
}

// DELETE请求
export const del = <T = any>(
  url: string,
  config: RequestConfig = {}
): Promise<ApiResponse<T>> => {
  return request<T>({
    method: 'delete',
    url,
    ...config
  })
}

// PATCH请求
export const patch = <T = any>(
  url: string,
  data?: any,
  config: RequestConfig = {}
): Promise<ApiResponse<T>> => {
  return request<T>({
    method: 'patch',
    url,
    data,
    ...config
  })
}

// 文件上传
export const upload = <T = any>(
  url: string,
  file: File,
  config: RequestConfig = {}
): Promise<ApiResponse<T>> => {
  const formData = new FormData()
  formData.append('file', file)
  
  return request<T>({
    method: 'post',
    url,
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    ...config
  })
}

// 下载文件
export const download = (
  url: string,
  params?: any,
  filename?: string
): Promise<void> => {
  return request({
    method: 'get',
    url,
    params,
    responseType: 'blob'
  }).then((response) => {
    const blob = new Blob([response.data])
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename || 'download'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)
  })
}

export default service