/**
 * 格式化工具函数
 */

/**
 * 格式化金额
 * @param amount 金额
 * @param decimals 小数位数
 * @param separator 千分位分隔符
 * @returns 格式化后的金额字符串
 */
export const formatMoney = (
  amount: number | string,
  decimals: number = 2,
  separator: string = ','
): string => {
  if (amount === null || amount === undefined || amount === '') {
    return '0.00'
  }

  const num = Number(amount)
  if (isNaN(num)) {
    return '0.00'
  }

  const fixed = num.toFixed(decimals)
  const parts = fixed.split('.')
  
  // 格式化整数部分
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, separator)
  
  return parts.join('.')
}

/**
 * 格式化时间
 * @param time 时间字符串、时间戳或Date对象
 * @param format 格式化字符串
 * @returns 格式化后的时间字符串
 */
export const formatTime = (
  time: string | number | Date,
  format: string = 'YYYY-MM-DD HH:mm:ss'
): string => {
  if (!time) return ''

  let date: Date
  
  if (typeof time === 'string') {
    // 处理相对时间格式
    if (time.includes('ago') || time.includes('前')) {
      return time
    }
    date = new Date(time)
  } else if (typeof time === 'number') {
    // 判断是秒还是毫秒
    date = time.toString().length === 10 ? new Date(time * 1000) : new Date(time)
  } else {
    date = time
  }

  if (isNaN(date.getTime())) {
    return ''
  }

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return format
    .replace('YYYY', year.toString())
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 格式化相对时间
 * @param time 时间字符串、时间戳或Date对象
 * @returns 相对时间字符串
 */
export const formatRelativeTime = (
  time: string | number | Date
): string => {
  if (!time) return ''

  const now = new Date()
  let date: Date

  if (typeof time === 'string') {
    date = new Date(time)
  } else if (typeof time === 'number') {
    date = time.toString().length === 10 ? new Date(time * 1000) : new Date(time)
  } else {
    date = time
  }

  if (isNaN(date.getTime())) {
    return ''
  }

  const diff = now.getTime() - date.getTime()
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  const months = Math.floor(days / 30)
  const years = Math.floor(months / 12)

  if (seconds < 60) {
    return '刚刚'
  } else if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else if (days < 30) {
    return `${days}天前`
  } else if (months < 12) {
    return `${months}个月前`
  } else {
    return `${years}年前`
  }
}

/**
 * 格式化文件大小
 * @param bytes 字节数
 * @param decimals 小数位数
 * @returns 格式化后的文件大小字符串
 */
export const formatFileSize = (
  bytes: number,
  decimals: number = 2
): string => {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(decimals)) + ' ' + sizes[i]
}

/**
 * 格式化数字
 * @param num 数字
 * @param decimals 小数位数
 * @param separator 千分位分隔符
 * @returns 格式化后的数字字符串
 */
export const formatNumber = (
  num: number,
  decimals: number = 0,
  separator: string = ','
): string => {
  if (isNaN(num)) return '0'

  const fixed = num.toFixed(decimals)
  const parts = fixed.split('.')
  
  // 格式化整数部分
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, separator)
  
  return parts.join('.')
}

/**
 * 格式化百分比
 * @param value 数值
 * @param total 总数
 * @param decimals 小数位数
 * @returns 百分比字符串
 */
export const formatPercentage = (
  value: number,
  total: number,
  decimals: number = 1
): string => {
  if (total === 0) return '0%'

  const percentage = (value / total) * 100
  return `${percentage.toFixed(decimals)}%`
}

/**
 * 格式化手机号
 * @param phone 手机号
 * @returns 格式化后的手机号
 */
export const formatPhone = (phone: string): string => {
  if (!phone || phone.length !== 11) return phone

  return phone.replace(/(\d{3})(\d{4})(\d{4})/, '$1****$3')
}

/**
 * 格式化身份证号
 * @param idCard 身份证号
 * @returns 格式化后的身份证号
 */
export const formatIdCard = (idCard: string): string => {
  if (!idCard || idCard.length !== 18) return idCard

  return idCard.replace(/(\d{6})(\d{8})(\d{4})/, '$1********$3')
}

/**
 * 格式化银行卡号
 * @param bankCard 银行卡号
 * @returns 格式化后的银行卡号
 */
export const formatBankCard = (bankCard: string): string => {
  if (!bankCard) return bankCard

  const cleaned = bankCard.replace(/\s/g, '')
  if (cleaned.length < 16) return bankCard

  return cleaned.replace(/(\d{4})(?=\d)/g, '$1 ')
}

/**
 * 格式化地址
 * @param province 省
 * @param city 市
 * @param district 区县
 * @param detail 详细地址
 * @returns 完整地址
 */
export const formatAddress = (
  province?: string,
  city?: string,
  district?: string,
  detail?: string
): string => {
  const parts = [province, city, district, detail].filter(Boolean)
  return parts.join('')
}

/**
 * 格式化任务状态
 * @param status 状态码
 * @returns 状态文本
 */
export const formatTaskStatus = (status: string): string => {
  const statusMap: Record<string, string> = {
    pending: '待接取',
    pending_payment: '待支付',
    in_progress: '进行中',
    delivered: '已交付',
    completed: '已完成',
    cancelled: '已取消',
    rejected: '已拒绝'
  }

  return statusMap[status] || '未知'
}

/**
 * 格式化支付状态
 * @param status 状态码
 * @returns 状态文本
 */
export const formatPaymentStatus = (status: string): string => {
  const statusMap: Record<string, string> = {
    pending: '待支付',
    processing: '处理中',
    success: '支付成功',
    failed: '支付失败',
    cancelled: '已取消',
    refunded: '已退款'
  }

  return statusMap[status] || '未知'
}

/**
 * 格式化订单类型
 * @param type 类型码
 * @returns 类型文本
 */
export const formatOrderType = (type: string): string => {
  const typeMap: Record<string, string> = {
    task: '任务支付',
    recharge: '账户充值',
    withdrawal: '账户提现',
    settlement: '任务结算',
    refund: '退款'
  }

  return typeMap[type] || '其他'
}

/**
 * 格式化交易类型
 * @param type 类型码
 * @returns 类型文本
 */
export const formatTransactionType = (type: string): string => {
  const typeMap: Record<string, string> = {
    income: '收入',
    expense: '支出',
    freeze: '冻结',
    unfreeze: '解冻',
    refund: '退款'
  }

  return typeMap[type] || '其他'
}

/**
 * 格式化星级评分
 * @param rating 评分
 * @param maxRating 最高评分
 * @returns 星级字符串
 */
export const formatStarRating = (rating: number, maxRating: number = 5): string => {
  const fullStars = Math.floor(rating)
  const hasHalfStar = rating - fullStars >= 0.5
  const emptyStars = maxRating - fullStars - (hasHalfStar ? 1 : 0)

  return '★'.repeat(fullStars) + (hasHalfStar ? '☆' : '') + '☆'.repeat(emptyStars)
}

/**
 * 格式化标签列表
 * @param tags 标签数组
 * @param separator 分隔符
 * @returns 标签字符串
 */
export const formatTags = (tags: string[], separator: string = ', '): string => {
  if (!Array.isArray(tags) || tags.length === 0) {
    return ''
  }

  return tags.filter(Boolean).join(separator)
}

/**
 * 格式化URL参数
 * @param params 参数对象
 * @returns URL参数字符串
 */
export const formatUrlParams = (params: Record<string, any>): string => {
  const searchParams = new URLSearchParams()

  Object.keys(params).forEach(key => {
    const value = params[key]
    if (value !== null && value !== undefined && value !== '') {
      searchParams.append(key, String(value))
    }
  })

  return searchParams.toString()
}

/**
 * 格式化持续时间
 * @param seconds 秒数
 * @returns 格式化的持续时间字符串
 */
export const formatDuration = (seconds: number): string => {
  if (seconds < 60) {
    return `${seconds}秒`
  } else if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = seconds % 60
    return remainingSeconds > 0 ? `${minutes}分${remainingSeconds}秒` : `${minutes}分钟`
  } else if (seconds < 86400) {
    const hours = Math.floor(seconds / 3600)
    const remainingMinutes = Math.floor((seconds % 3600) / 60)
    return remainingMinutes > 0 ? `${hours}小时${remainingMinutes}分钟` : `${hours}小时`
  } else {
    const days = Math.floor(seconds / 86400)
    const remainingHours = Math.floor((seconds % 86400) / 3600)
    return remainingHours > 0 ? `${days}天${remainingHours}小时` : `${days}天`
  }
}

/**
 * 格式化倒计时
 * @param endTime 结束时间
 * @returns 倒计时字符串
 */
export const formatCountdown = (endTime: string | number | Date): string => {
  const now = new Date().getTime()
  let end: number

  if (typeof endTime === 'string') {
    end = new Date(endTime).getTime()
  } else if (typeof endTime === 'number') {
    end = endTime.toString().length === 10 ? endTime * 1000 : endTime
  } else {
    end = endTime.getTime()
  }

  if (end <= now) {
    return '已结束'
  }

  const diff = end - now
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)

  const parts: string[] = []

  if (days > 0) parts.push(`${days}天`)
  if (hours > 0) parts.push(`${hours}时`)
  if (minutes > 0) parts.push(`${minutes}分`)
  if (seconds > 0 && days === 0) parts.push(`${seconds}秒`)

  return parts.join('') || '即将结束'
}

/**
 * 格式化数字单位
 * @param num 数字
 * @param decimals 小数位数
 * @returns 格式化后的数字字符串
 */
export const formatUnit = (num: number, decimals: number = 1): string => {
  if (num < 1000) {
    return num.toString()
  } else if (num < 10000) {
    return `${(num / 1000).toFixed(decimals)}K`
  } else if (num < 100000000) {
    return `${(num / 10000).toFixed(decimals)}W`
  } else {
    return `${(num / 100000000).toFixed(decimals)}亿`
  }
}

/**
 * 格式化颜色值
 * @param color 颜色值
 * @returns 格式化后的颜色值
 */
export const formatColor = (color: string): string => {
  if (!color) return '#000000'

  // 如果已经是HEX格式，直接返回
  if (color.startsWith('#')) {
    return color
  }

  // 如果是RGB格式，转换为HEX
  const rgbMatch = color.match(/rgb\((\d+),\s*(\d+),\s*(\d+)\)/)
  if (rgbMatch) {
    const r = parseInt(rgbMatch[1])
    const g = parseInt(rgbMatch[2])
    const b = parseInt(rgbMatch[3])
    return `#${((r << 16) | (g << 8) | b).toString(16).padStart(6, '0')}`
  }

  return color
}

/**
 * 格式化文本截断
 * @param text 文本
 * @param maxLength 最大长度
 * @param suffix 后缀
 * @returns 截断后的文本
 */
export const formatTextTruncate = (
  text: string,
  maxLength: number,
  suffix: string = '...'
): string => {
  if (!text || text.length <= maxLength) {
    return text || ''
  }

  return text.substring(0, maxLength - suffix.length) + suffix
}

/**
 * 格式化关键词高亮
 * @param text 文本
 * @param keyword 关键词
 * @param tag 标签名
 * @returns 高亮后的HTML字符串
 */
export const formatKeywordHighlight = (
  text: string,
  keyword: string,
  tag: string = 'mark'
): string => {
  if (!text || !keyword) {
    return text || ''
  }

  const regex = new RegExp(`(${keyword})`, 'gi')
  return text.replace(regex, `<${tag}>$1</${tag}>`)
}