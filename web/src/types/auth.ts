// 认证相关类型
export interface LoginParams {
  username: string
  password: string
  captcha?: string
}

export interface RegisterParams {
  username: string
  email: string
  password: string
  confirmPassword: string
  phone?: string
  captcha?: string
}

export interface UserInfo {
  id: number
  username: string
  nickname?: string
  avatar?: string
  email?: string
  phone?: string
  role: string
  status: 'active' | 'inactive' | 'banned'
  createTime: string
  updateTime: string
  lastLoginTime?: string
  isEmailVerified: boolean
  isPhoneVerified: boolean
  profile?: {
    bio?: string
    location?: string
    website?: string
    skills?: string[]
  }
  wechat_openid?: string
  alipay_userid?: string
  alipay_verified?: boolean
}

export interface LoginResult {
  token: string
  refreshToken: string
  userInfo: UserInfo
  expiresIn: number
}

export interface CaptchaResult {
  captchaId: string
  captchaImage: string
}

export interface PasswordResetParams {
  email: string
  captcha: string
  newPassword: string
}

export interface EmailVerifyParams {
  email: string
  code: string
}