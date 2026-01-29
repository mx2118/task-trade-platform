/**
 * 表单验证工具函数
 */

/**
 * 验证手机号
 * @param rule 规则对象
 * @param value 输入值
 * @param callback 回调函数
 */
export const validatePhone = (rule: any, value: string, callback: Function) => {
  if (!value) {
    callback(new Error('请输入手机号'))
    return
  }
  
  const phoneRegex = /^1[3-9]\d{9}$/
  if (!phoneRegex.test(value)) {
    callback(new Error('请输入正确的手机号'))
    return
  }
  
  callback()
}

/**
 * 验证验证码
 * @param rule 规则对象
 * @param value 输入值
 * @param callback 回调函数
 */
export const validateCode = (rule: any, value: string, callback: Function) => {
  if (!value) {
    callback(new Error('请输入验证码'))
    return
  }
  
  if (value.length !== 6) {
    callback(new Error('验证码为6位数字'))
    return
  }
  
  if (!/^\d{6}$/.test(value)) {
    callback(new Error('验证码只能是数字'))
    return
  }
  
  callback()
}

/**
 * 验证密码
 * @param rule 规则对象
 * @param value 输入值
 * @param callback 回调函数
 */
export const validatePassword = (rule: any, value: string, callback: Function) => {
  if (!value) {
    callback(new Error('请输入密码'))
    return
  }
  
  if (value.length < 6) {
    callback(new Error('密码长度不能少于6位'))
    return
  }
  
  if (value.length > 20) {
    callback(new Error('密码长度不能超过20位'))
    return
  }
  
  // 至少包含字母和数字
  if (!/(?=.*[a-zA-Z])(?=.*\d)/.test(value)) {
    callback(new Error('密码必须包含字母和数字'))
    return
  }
  
  callback()
}

/**
 * 验证确认密码
 * @param password 原密码
 * @returns 验证函数
 */
export const validateConfirmPassword = (password: string) => {
  return (rule: any, value: string, callback: Function) => {
    if (!value) {
      callback(new Error('请确认密码'))
      return
    }
    
    if (value !== password) {
      callback(new Error('两次输入的密码不一致'))
      return
    }
    
    callback()
  }
}

/**
 * 验证昵称
 * @param rule 规则对象
 * @param value 输入值
 * @param callback 回调函数
 */
export const validateNickname = (rule: any, value: string, callback: Function) => {
  if (!value) {
    callback(new Error('请输入昵称'))
    return
  }
  
  if (value.length < 2) {
    callback(new Error('昵称长度不能少于2位'))
    return
  }
  
  if (value.length > 20) {
    callback(new Error('昵称长度不能超过20位'))
    return
  }
  
  // 不允许特殊字符
  if (!/^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/.test(value)) {
    callback(new Error('昵称只能包含中文、字母、数字、下划线和横线'))
    return
  }
  
  callback()
}

/**
 * 验证邮箱
 * @param rule 规则对象
 * @param value 输入值
 * @param callback 回调函数
 */
export const validateEmail = (rule: any, value: string, callback: Function) => {
  if (!value) {
    callback(new Error('请输入邮箱'))
    return
  }
  
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  if (!emailRegex.test(value)) {
    callback(new Error('请输入正确的邮箱地址'))
    return
  }
  
  callback()
}

/**
 * 验证身份证号
 * @param rule 规则对象
 * @param value 输入值
 * @param callback 回调函数
 */
export const validateIdCard = (rule: any, value: string, callback: Function) => {
  if (!value) {
    callback(new Error('请输入身份证号'))
    return
  }
  
  // 18位身份证号正则
  const idCardRegex = /^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dX]$/
  if (!idCardRegex.test(value)) {
    callback(new Error('请输入正确的身份证号'))
    return
  }
  
  // 验证最后一位校验码
  const weights = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]
  const checkCodes = ['1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2']
  
  let sum = 0
  for (let i = 0; i < 17; i++) {
    sum += parseInt(value[i]) * weights[i]
  }
  
  const checkCode = checkCodes[sum % 11]
  if (value[17].toUpperCase() !== checkCode) {
    callback(new Error('请输入正确的身份证号'))
    return
  }
  
  callback()
}

/**
 * 验证银行卡号
 * @param rule 规则对象
 * @param value 输入值
 * @param callback 回调函数
 */
export const validateBankCard = (rule: any, value: string, callback: Function) => {
  if (!value) {
    callback(new Error('请输入银行卡号'))
    return
  }
  
  // 去除空格
  const cardNumber = value.replace(/\s/g, '')
  
  // 银行卡号长度验证（通常16-19位）
  if (cardNumber.length < 16 || cardNumber.length > 19) {
    callback(new Error('银行卡号长度不正确'))
    return
  }
  
  // Luhn算法验证
  let sum = 0
  let isEven = false
  
  for (let i = cardNumber.length - 1; i >= 0; i--) {
    let digit = parseInt(cardNumber[i])
    
    if (isEven) {
      digit *= 2
      if (digit > 9) {
        digit = digit % 10 + 1
      }
    }
    
    sum += digit
    isEven = !isEven
  }
  
  if (sum % 10 !== 0) {
    callback(new Error('请输入正确的银行卡号'))
    return
  }
  
  callback()
}

/**
 * 验证金额
 * @param min 最小金额
 * @param max 最大金额
 * @returns 验证函数
 */
export const validateAmount = (min: number = 0, max: number = 1000000) => {
  return (rule: any, value: string, callback: Function) => {
    if (!value) {
      callback(new Error('请输入金额'))
      return
    }
    
    const amount = parseFloat(value)
    
    if (isNaN(amount)) {
      callback(new Error('请输入有效的金额'))
      return
    }
    
    if (amount < min) {
      callback(new Error(`金额不能少于${min}元`))
      return
    }
    
    if (amount > max) {
      callback(new Error(`金额不能超过${max}元`))
      return
    }
    
    // 验证小数位数（最多2位）
    if (!/^\d+(\.\d{1,2})?$/.test(value)) {
      callback(new Error('金额最多保留2位小数'))
      return
    }
    
    callback()
  }
}

/**
 * 验证URL
 * @param rule 规则对象
 * @param value 输入值
 * @param callback 回调函数
 */
export const validateUrl = (rule: any, value: string, callback: Function) => {
  if (!value) {
    callback(new Error('请输入URL'))
    return
  }
  
  try {
    new URL(value)
    callback()
  } catch (error) {
    callback(new Error('请输入正确的URL'))
  }
}

/**
 * 验证中文姓名
 * @param rule 规则对象
 * @param value 输入值
 * @param callback 回调函数
 */
export const validateChineseName = (rule: any, value: string, callback: Function) => {
  if (!value) {
    callback(new Error('请输入姓名'))
    return
  }
  
  const nameRegex = /^[\u4e00-\u9fa5]{2,10}$/
  if (!nameRegex.test(value)) {
    callback(new Error('请输入2-10位中文姓名'))
    return
  }
  
  callback()
}

/**
 * 验证年龄
 * @param min 最小年龄
 * @param max 最大年龄
 * @returns 验证函数
 */
export const validateAge = (min: number = 18, max: number = 100) => {
  return (rule: any, value: string, callback: Function) => {
    if (!value) {
      callback(new Error('请输入年龄'))
      return
    }
    
    const age = parseInt(value)
    
    if (isNaN(age)) {
      callback(new Error('请输入有效的年龄'))
      return
    }
    
    if (age < min) {
      callback(new Error(`年龄不能小于${min}岁`))
      return
    }
    
    if (age > max) {
      callback(new Error(`年龄不能大于${max}岁`))
      return
    }
    
    callback()
  }
}

/**
 * 验证文件大小
 * @param maxSize 最大大小（字节）
 * @returns 验证函数
 */
export const validateFileSize = (maxSize: number) => {
  return (file: File) => {
    if (file.size > maxSize) {
      const sizeMB = (maxSize / 1024 / 1024).toFixed(2)
      return `文件大小不能超过${sizeMB}MB`
    }
    return ''
  }
}

/**
 * 验证文件类型
 * @param allowedTypes 允许的文件类型
 * @returns 验证函数
 */
export const validateFileType = (allowedTypes: string[]) => {
  return (file: File) => {
    const fileExtension = file.name.split('.').pop()?.toLowerCase()
    if (!fileExtension || !allowedTypes.includes(fileExtension)) {
      return `只支持${allowedTypes.join('、')}格式的文件`
    }
    return ''
  }
}

/**
 * 自定义验证器工厂
 * @param validator 验证函数
 * @param message 错误消息
 * @returns 验证规则
 */
export const createValidator = (validator: Function, message: string) => {
  return {
    validator: (rule: any, value: any, callback: Function) => {
      if (validator(value)) {
        callback()
      } else {
        callback(new Error(message))
      }
    }
  }
}

/**
 * 组合验证器
 * @param validators 验证器数组
 * @returns 验证函数
 */
export const combineValidators = (...validators: Function[]) => {
  return (rule: any, value: any, callback: Function) => {
    for (const validator of validators) {
      try {
        validator(rule, value, callback)
        return
      } catch (error) {
        callback(error)
        return
      }
    }
    callback()
  }
}

/**
 * 异步验证器
 * @param asyncValidator 异步验证函数
 * @param timeout 超时时间
 * @returns 验证函数
 */
export const createAsyncValidator = (asyncValidator: Function, timeout: number = 5000) => {
  return (rule: any, value: any, callback: Function) => {
    const timer = setTimeout(() => {
      callback(new Error('验证超时'))
    }, timeout)
    
    Promise.resolve(asyncValidator(value))
      .then(() => {
        clearTimeout(timer)
        callback()
      })
      .catch((error) => {
        clearTimeout(timer)
        callback(new Error(error.message || '验证失败'))
      })
  }
}