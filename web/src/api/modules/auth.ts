import request from '../request'
import type { ApiResponse, LoginParams, RegisterParams, LoginResult, UserInfo, CaptchaResult, PasswordResetParams, EmailVerifyParams } from '@/types'

// 登录
export const login = (data: LoginParams): Promise<ApiResponse<LoginResult>> => {
  return request.post('/api/auth/login', data)
}

// 注册
export const register = (data: RegisterParams): Promise<ApiResponse<UserInfo>> => {
  return request.post('/api/auth/register', data)
}

// 登出
export const logout = (): Promise<ApiResponse<null>> => {
  return request.post('/api/auth/logout')
}

// 刷新Token
export const refreshToken = (refreshToken: string): Promise<ApiResponse<{ token: string; expiresIn: number }>> => {
  return request.post('/api/auth/refresh', { refreshToken })
}

// 获取当前用户信息
export const getCurrentUser = (): Promise<ApiResponse<UserInfo>> => {
  return request.get('/api/auth/me')
}

// 更新用户信息
export const updateUserInfo = (data: Partial<UserInfo>): Promise<ApiResponse<UserInfo>> => {
  return request.put('/api/auth/profile', data)
}

// 修改密码
export const changePassword = (data: { oldPassword: string; newPassword: string }): Promise<ApiResponse<null>> => {
  return request.post('/api/auth/change-password', data)
}

// 获取验证码
export const getCaptcha = (): Promise<ApiResponse<CaptchaResult>> => {
  return request.get('/api/auth/captcha')
}

// 忘记密码 - 发送邮件
export const forgotPassword = (email: string): Promise<ApiResponse<null>> => {
  return request.post('/api/auth/forgot-password', { email })
}

// 重置密码
export const resetPassword = (data: PasswordResetParams): Promise<ApiResponse<null>> => {
  return request.post('/api/auth/reset-password', data)
}

// 验证邮箱
export const verifyEmail = (data: EmailVerifyParams): Promise<ApiResponse<null>> => {
  return request.post('/api/auth/verify-email', data)
}

// 发送邮箱验证码
export const sendEmailVerification = (email: string): Promise<ApiResponse<null>> => {
  return request.post('/api/auth/send-email-verification', { email })
}

// 绑定微信
export const bindWechat = (code: string): Promise<ApiResponse<{ openid: string }>> => {
  return request.post('/api/auth/bind-wechat', { code })
}

// 绑定支付宝
export const bindAlipay = (authCode: string): Promise<ApiResponse<{ userId: string; verified: boolean }>> => {
  return request.post('/api/auth/bind-alipay', { authCode })
}

// 导出 API 对象
export const authApi = {
  login,
  register,
  logout,
  refreshToken,
  getCurrentUser,
  updateUserInfo,
  changePassword,
  getCaptcha,
  forgotPassword,
  resetPassword,
  verifyEmail,
  sendEmailVerification,
  bindWechat,
  bindAlipay
}