/**
 * 懒加载工具 - 简化版
 */

/**
 * 使用 requestIdleCallback 延迟执行非关键任务
 */
export function runWhenIdle(callback: () => void, timeout = 2000) {
  if ('requestIdleCallback' in window) {
    requestIdleCallback(callback, { timeout })
  } else {
    setTimeout(callback, 100)
  }
}

/**
 * 检测网络状况
 */
export function getNetworkStatus() {
  const connection = (navigator as any).connection || (navigator as any).mozConnection || (navigator as any).webkitConnection
  
  if (!connection) {
    return { effectiveType: '4g', saveData: false }
  }

  return {
    effectiveType: connection.effectiveType || '4g',
    saveData: connection.saveData || false,
    downlink: connection.downlink,
    rtt: connection.rtt
  }
}

/**
 * 根据网络状况判断是否加载高质量资源
 */
export function shouldLoadHighQuality(): boolean {
  const network = getNetworkStatus()
  
  // 如果用户开启了省流量模式，或者网络较慢，则加载低质量资源
  if (network.saveData || network.effectiveType === 'slow-2g' || network.effectiveType === '2g') {
    return false
  }
  
  return true
}
