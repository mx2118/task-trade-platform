/**
 * 格式化工具函数
 */

import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

// 扩展dayjs
dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

/**
 * 格式化时间
 */
export function formatTime(time: string | Date | null | undefined, format: string = 'YYYY-MM-DD HH:mm'): string {
  if (!time) return ''
  
  const date = typeof time === 'string' ? dayjs(time) : dayjs(time)
  return date.format(format)
}

/**
 * 格式化相对时间
 */
export function formatRelativeTime(time: string | Date | null | undefined): string {
  if (!time) return ''
  
  const date = typeof time === 'string' ? dayjs(time) : dayjs(time)
  return date.fromNow()
}

/**
 * 格式化金额
 */
export function formatAmount(amount: number | string, precision: number = 2): string {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  if (isNaN(num)) return '¥0.00'
  
  return `¥${num.toFixed(precision)}`
}

/**
 * 格式化金额 (别名)
 */
export function formatMoney(amount: number | string, precision: number = 2): string {
  return formatAmount(amount, precision)
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

/**
 * 格式化数字（千分位）
 */
export function formatNumber(num: number | string): string {
  const n = typeof num === 'string' ? parseFloat(num) : num
  return n.toLocaleString('zh-CN')
}

/**
 * 格式化百分比
 */
export function formatPercentage(value: number, precision: number = 1): string {
  return `${(value * 100).toFixed(precision)}%`
}

/**
 * 格式化电话号码
 */
export function formatPhone(phone: string): string {
  if (!phone) return ''
  
  // 简单的手机号格式化（隐藏中间4位）
  if (phone.length === 11) {
    return `${phone.slice(0, 3)}****${phone.slice(7)}`
  }
  
  return phone
}

/**
 * 格式化银行卡号
 */
export function formatBankCard(cardNumber: string): string {
  if (!cardNumber) return ''
  
  // 隐藏部分数字，只显示前后4位
  if (cardNumber.length >= 8) {
    return `${cardNumber.slice(0, 4)} **** **** ${cardNumber.slice(-4)}`
  }
  
  return cardNumber
}

/**
 * 格式化身份证号
 */
export function formatIdCard(idCard: string): string {
  if (!idCard) return ''
  
  // 隐藏中间部分，只显示前后6位
  if (idCard.length >= 8) {
    return `${idCard.slice(0, 6)}********${idCard.slice(-4)}`
  }
  
  return idCard
}

/**
 * 获取时间段描述
 */
export function getTimePeriod(hour: number): string {
  if (hour >= 6 && hour < 12) return '上午'
  if (hour >= 12 && hour < 18) return '下午'
  if (hour >= 18 && hour < 24) return '晚上'
  return '凌晨'
}

/**
 * 格式化地址
 */
export function formatAddress(address: {
  province?: string
  city?: string
  district?: string
  detail?: string
}): string {
  const parts = [
    address.province,
    address.city,
    address.district,
    address.detail
  ].filter(Boolean)
  
  return parts.join('')
}

/**
 * 格式化URL（添加协议）
 */
export function formatUrl(url: string): string {
  if (!url) return ''
  
  if (!/^https?:\/\//.test(url)) {
    return `https://${url}`
  }
  
  return url
}

/**
 * 截断文本
 */
export function truncateText(text: string, maxLength: number = 50): string {
  if (!text) return ''
  
  if (text.length <= maxLength) return text
  
  return `${text.slice(0, maxLength)}...`
}

/**
 * 获取文件扩展名
 */
export function getFileExtension(filename: string): string {
  if (!filename) return ''
  
  const parts = filename.split('.')
  return parts.length > 1 ? parts[parts.length - 1].toLowerCase() : ''
}

/**
 * 判断是否为图片文件
 */
export function isImageFile(filename: string): boolean {
  const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp']
  const extension = getFileExtension(filename)
  return imageExtensions.includes(extension)
}

/**
 * 获取文件类型图标
 */
export function getFileIcon(filename: string): string {
  const extension = getFileExtension(filename)
  
  const iconMap: Record<string, string> = {
    // 图片
    jpg: 'Picture',
    jpeg: 'Picture',
    png: 'Picture',
    gif: 'Picture',
    webp: 'Picture',
    bmp: 'Picture',
    
    // 文档
    pdf: 'Document',
    doc: 'Document',
    docx: 'Document',
    xls: 'Document',
    xlsx: 'Document',
    ppt: 'Document',
    pptx: 'Document',
    txt: 'Document',
    
    // 压缩文件
    zip: 'FolderOpened',
    rar: 'FolderOpened',
    '7z': 'FolderOpened',
    tar: 'FolderOpened',
    
    // 音频
    mp3: 'Headphones',
    wav: 'Headphones',
    flac: 'Headphones',
    
    // 视频
    mp4: 'VideoPlay',
    avi: 'VideoPlay',
    mov: 'VideoPlay',
    wmv: 'VideoPlay',
    
    // 代码文件
    js: 'Code',
    ts: 'Code',
    html: 'Code',
    css: 'Code',
    java: 'Code',
    py: 'Code',
    cpp: 'Code',
    c: 'Code',
    
    // 其他
    default: 'Document'
  }
  
  return iconMap[extension] || iconMap.default
}

/**
 * 生成颜色（基于字符串）
 */
export function generateColorFromString(str: string): string {
  if (!str) return '#667eea'
  
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  
  const hue = hash % 360
  return `hsl(${hue}, 70%, 50%)`
}

/**
 * 验证手机号
 */
export function validatePhone(phone: string): boolean {
  const phoneRegex = /^1[3-9]\d{9}$/
  return phoneRegex.test(phone)
}

/**
 * 验证邮箱
 */
export function validateEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * 验证身份证号
 */
export function validateIdCard(idCard: string): boolean {
  if (!idCard) return false
  
  // 简单验证：长度为15或18位
  return /^\d{15}$|^\d{17}[\dXx]$/.test(idCard)
}

/**
 * 获取随机字符串
 */
export function getRandomString(length: number = 8): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  
  return result
}

/**
 * 防抖函数
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number,
  immediate?: boolean
): (...args: Parameters<T>) => void {
  let timeout: ReturnType<typeof setTimeout> | null = null
  
  return function(this: any, ...args: Parameters<T>) {
    const later = () => {
      timeout = null
      if (!immediate) func.apply(this, args)
    }
    
    const callNow = immediate && !timeout
    
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(later, wait)
    
    if (callNow) func.apply(this, args)
  }
}

/**
 * 节流函数
 */
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  wait: number,
  options?: { leading?: boolean; trailing?: boolean }
): (...args: Parameters<T>) => void {
  let timeout: ReturnType<typeof setTimeout> | null = null
  let previous = 0
  
  const { leading = true, trailing = true } = options || {}
  
  return function(this: any, ...args: Parameters<T>) {
    const now = Date.now()
    const remaining = wait - (now - previous)
    
    if (remaining <= 0 || remaining > wait) {
      if (timeout) {
        clearTimeout(timeout)
        timeout = null
      }
      
      if (leading) {
        previous = now
        func.apply(this, args)
      }
    } else if (!timeout && trailing) {
      timeout = setTimeout(() => {
        previous = Date.now()
        timeout = null
        func.apply(this, args)
      }, remaining)
    }
  }
}