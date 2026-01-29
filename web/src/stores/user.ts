import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, logout, getCurrentUser as getUserInfo } from '@/api/modules/auth'
import { getToken, setToken, removeToken } from '@/utils/auth'
import type { LoginParams, UserInfo } from '@/types/user'
import { ElMessage } from 'element-plus'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string>(getToken() || '')
  const userInfo = ref<UserInfo | null>(null)
  const permissions = ref<string[]>([])
  const roles = ref<string[]>([])
  
  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || '')
  const userId = computed(() => userInfo.value?.id || 0)
  const userRole = computed(() => userInfo.value?.role || 'user')
  const avatar = computed(() => userInfo.value?.avatar || '')
  const email = computed(() => userInfo.value?.email || '')
  const phone = computed(() => userInfo.value?.phone || '')
  
  // Actions
  const loginAction = async (loginData: LoginParams) => {
    try {
      const response = await login(loginData)
      const { token: newToken } = response.data
      
      // 保存token
      token.value = newToken
      setToken(newToken)
      
      // 获取用户信息
      await getUserInfoAction()
      
      ElMessage.success('登录成功')
      return response
    } catch (error) {
      console.error('登录失败:', error)
      ElMessage.error('登录失败，请检查用户名和密码')
      throw error
    }
  }
  
  const logoutAction = async () => {
    try {
      await logout()
    } catch (error) {
      console.error('退出登录失败:', error)
    } finally {
      // 清除本地状态
      token.value = ''
      userInfo.value = null
      permissions.value = []
      roles.value = []
      removeToken()
      
      ElMessage.success('已退出登录')
    }
  }
  
  const getUserInfoAction = async () => {
    try {
      const response = await getUserInfo()
      const info = response.data
      
      userInfo.value = info
      permissions.value = info.permissions || []
      roles.value = info.roles || [info.role || 'user']
      
      return info
    } catch (error) {
      console.error('获取用户信息失败:', error)
      throw error
    }
  }
  
  const updateUserInfo = (newInfo: Partial<UserInfo>) => {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...newInfo }
    }
  }
  
  const resetUserState = () => {
    token.value = ''
    userInfo.value = null
    permissions.value = []
    roles.value = []
    removeToken()
  }
  
  // 检查权限
  const hasPermission = (permission: string) => {
    if (roles.value.includes('admin')) {
      return true
    }
    return permissions.value.includes(permission)
  }
  
  const hasRole = (role: string) => {
    return roles.value.includes(role)
  }
  
  const hasAnyPermission = (permissionList: string[]) => {
    if (roles.value.includes('admin')) {
      return true
    }
    return permissionList.some(permission => permissions.value.includes(permission))
  }
  
  const hasAllPermissions = (permissionList: string[]) => {
    if (roles.value.includes('admin')) {
      return true
    }
    return permissionList.every(permission => permissions.value.includes(permission))
  }
  
  // 微信登录
  const wechatLogin = async (code: string) => {
    try {
      // 这里应该调用微信登录API
      // const response = await wechatLoginApi(code)
      // 暂时使用模拟数据
      const mockResponse = {
        code: 200,
        data: {
          token: 'wechat-mock-token',
          userInfo: {
            id: 1,
            username: 'wechat_user',
            nickname: '微信用户',
            avatar: '',
            email: '',
            phone: '',
            role: 'user',
            permissions: [],
            roles: ['user'],
            creditScore: 100,
            balance: 0,
            isVerified: false,
            createdAt: new Date().toISOString()
          }
        }
      }
      
      const { token: newToken, userInfo: info } = mockResponse.data
      
      token.value = newToken
      setToken(newToken)
      userInfo.value = info
      permissions.value = info.permissions || []
      roles.value = info.roles
      
      ElMessage.success('微信登录成功')
      return mockResponse
    } catch (error) {
      console.error('微信登录失败:', error)
      ElMessage.error('微信登录失败')
      throw error
    }
  }
  
  // 支付宝登录
  const alipayLogin = async (authCode: string) => {
    try {
      // 这里应该调用支付宝登录API
      // const response = await alipayLoginApi(authCode)
      // 暂时使用模拟数据
      const mockResponse = {
        code: 200,
        data: {
          token: 'alipay-mock-token',
          userInfo: {
            id: 2,
            username: 'alipay_user',
            nickname: '支付宝用户',
            avatar: '',
            email: '',
            phone: '',
            role: 'user',
            permissions: [],
            roles: ['user'],
            creditScore: 100,
            balance: 0,
            isVerified: false,
            createdAt: new Date().toISOString()
          }
        }
      }
      
      const { token: newToken, userInfo: info } = mockResponse.data
      
      token.value = newToken
      setToken(newToken)
      userInfo.value = info
      permissions.value = info.permissions || []
      roles.value = info.roles
      
      ElMessage.success('支付宝登录成功')
      return mockResponse
    } catch (error) {
      console.error('支付宝登录失败:', error)
      ElMessage.error('支付宝登录失败')
      throw error
    }
  }
  
  return {
    // 状态
    token,
    userInfo,
    permissions,
    roles,
    id: userId,
    
    // 计算属性
    isLoggedIn,
    username,
    userId,
    userRole,
    avatar,
    email,
    phone,
    
    // Actions (别名兼容)
    login: loginAction,
    logout: logoutAction,
    getUserInfo: getUserInfoAction,
    
    // Actions
    loginAction,
    logoutAction,
    getUserInfoAction,
    updateUserInfo,
    resetUserState,
    hasPermission,
    hasRole,
    hasAnyPermission,
    hasAllPermissions,
    wechatLogin,
    alipayLogin
  }
})