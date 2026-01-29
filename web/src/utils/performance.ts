/**
 * 前端性能优化工具集
 */

// 性能监控器
class PerformanceMonitor {
  private static instance: PerformanceMonitor
  private metrics: Map<string, any> = new Map()
  
  static getInstance(): PerformanceMonitor {
    if (!PerformanceMonitor.instance) {
      PerformanceMonitor.instance = new PerformanceMonitor()
    }
    return PerformanceMonitor.instance
  }
  
  // 记录页面加载性能
  recordPageLoad() {
    if (typeof window !== 'undefined' && 'performance' in window) {
      const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming
      
      this.metrics.set('pageLoad', {
        domContentLoaded: navigation.domContentLoadedEventEnd - navigation.domContentLoadedEventStart,
        loadComplete: navigation.loadEventEnd - navigation.loadEventStart,
        firstPaint: performance.getEntriesByType('paint').find(entry => entry.name === 'first-paint')?.startTime,
        firstContentfulPaint: performance.getEntriesByType('paint').find(entry => entry.name === 'first-contentful-paint')?.startTime,
        domInteractive: navigation.domInteractive - navigation.fetchStart,
        domComplete: navigation.domComplete - navigation.fetchStart
      })
    }
  }
  
  // 记录API请求性能
  recordAPICall(url: string, duration: number, success: boolean) {
    const key = `api_${url}`
    const existing = this.metrics.get(key) || { count: 0, totalDuration: 0, errors: 0 }
    
    this.metrics.set(key, {
      count: existing.count + 1,
      totalDuration: existing.totalDuration + duration,
      averageDuration: (existing.totalDuration + duration) / (existing.count + 1),
      errors: success ? existing.errors : existing.errors + 1
    })
  }
  
  // 获取性能指标
  getMetrics() {
    return Object.fromEntries(this.metrics)
  }
  
  // 发送性能数据到服务器
  async sendMetrics() {
    try {
      const metrics = this.getMetrics()
      // 发送到性能监控端点
      await fetch('/api/performance/metrics', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          timestamp: Date.now(),
          userAgent: navigator.userAgent,
          url: window.location.href,
          metrics
        })
      })
    } catch (error) {
      console.error('Failed to send performance metrics:', error)
    }
  }
}

// 图片懒加载管理器
class LazyImageLoader {
  private static instance: LazyImageLoader
  private observer: IntersectionObserver | null = null
  private imageMap: WeakMap<HTMLImageElement, string> = new WeakMap()
  
  static getInstance(): LazyImageLoader {
    if (!LazyImageLoader.instance) {
      LazyImageLoader.instance = new LazyImageLoader()
    }
    return LazyImageLoader.instance
  }
  
  init() {
    if ('IntersectionObserver' in window) {
      this.observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            const img = entry.target as HTMLImageElement
            const src = this.imageMap.get(img)
            if (src) {
              img.src = src
              img.classList.remove('lazy')
              this.observer?.unobserve(img)
              this.imageMap.delete(img)
            }
          }
        })
      }, {
        rootMargin: '50px 0px',
        threshold: 0.01
      })
    }
  }
  
  observe(img: HTMLImageElement, src: string) {
    if (this.observer) {
      this.imageMap.set(img, src)
      img.classList.add('lazy')
      this.observer.observe(img)
    } else {
      // 降级处理
      img.src = src
    }
  }
  
  destroy() {
    this.observer?.disconnect()
    this.observer = null
  }
}

// 虚拟滚动组件
export class VirtualScroller {
  private container: HTMLElement
  private itemHeight: number
  private visibleCount: number
  private startIndex = 0
  private items: any[] = []
  private renderedItems: HTMLElement[] = []
  
  constructor(container: HTMLElement, itemHeight: number) {
    this.container = container
    this.itemHeight = itemHeight
    this.visibleCount = Math.ceil(container.clientHeight / itemHeight) + 2 // 缓冲区
    
    this.init()
  }
  
  private init() {
    this.container.style.overflowY = 'auto'
    this.container.style.height = `${this.container.clientHeight}px`
    
    this.container.addEventListener('scroll', this.handleScroll.bind(this))
  }
  
  setItems(items: any[]) {
    this.items = items
    this.render()
  }
  
  private handleScroll() {
    const scrollTop = this.container.scrollTop
    const newStartIndex = Math.floor(scrollTop / this.itemHeight)
    
    if (newStartIndex !== this.startIndex) {
      this.startIndex = newStartIndex
      this.render()
    }
  }
  
  private render() {
    const endIndex = Math.min(this.startIndex + this.visibleCount, this.items.length)
    const fragment = document.createDocumentFragment()
    
    // 清空现有内容
    this.container.innerHTML = ''
    
    // 添加顶部占位
    if (this.startIndex > 0) {
      const topSpacer = document.createElement('div')
      topSpacer.style.height = `${this.startIndex * this.itemHeight}px`
      fragment.appendChild(topSpacer)
    }
    
    // 渲染可见项
    for (let i = this.startIndex; i < endIndex; i++) {
      const item = this.items[i]
      const itemElement = this.createItemElement(item, i)
      fragment.appendChild(itemElement)
    }
    
    // 添加底部占位
    const remainingHeight = (this.items.length - endIndex) * this.itemHeight
    if (remainingHeight > 0) {
      const bottomSpacer = document.createElement('div')
      bottomSpacer.style.height = `${remainingHeight}px`
      fragment.appendChild(bottomSpacer)
    }
    
    this.container.appendChild(fragment)
  }
  
  private createItemElement(item: any, index: number) {
    const div = document.createElement('div')
    div.style.height = `${this.itemHeight}px`
    div.style.borderBottom = '1px solid #eee'
    div.style.padding = '10px'
    div.textContent = `Item ${index}: ${JSON.stringify(item)}`
    return div
  }
}

// 防抖函数
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number,
  immediate = false
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

// 节流函数
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean = false
  
  return function(this: any, ...args: Parameters<T>) {
    if (!inThrottle) {
      func.apply(this, args)
      inThrottle = true
      setTimeout(() => inThobble = false, wait)
    }
  }
}

// 资源预加载器
export class ResourcePreloader {
  private static preloadedResources: Set<string> = new Set()
  
  // 预加载图片
  static preloadImage(src: string): Promise<HTMLImageElement> {
    return new Promise((resolve, reject) => {
      if (this.preloadedResources.has(src)) {
        const img = new Image()
        img.src = src
        resolve(img)
        return
      }
      
      const img = new Image()
      img.onload = () => {
        this.preloadedResources.add(src)
        resolve(img)
      }
      img.onerror = reject
      img.src = src
    })
  }
  
  // 预加载脚本
  static preloadScript(src: string): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.preloadedResources.has(src)) {
        resolve()
        return
      }
      
      const link = document.createElement('link')
      link.rel = 'preload'
      link.as = 'script'
      link.href = src
      link.onload = () => {
        this.preloadedResources.add(src)
        resolve()
      }
      link.onerror = reject
      document.head.appendChild(link)
    })
  }
  
  // 批量预加载资源
  static async preloadResources(resources: Array<{type: 'image' | 'script', src: string}>) {
    const promises = resources.map(resource => {
      if (resource.type === 'image') {
        return this.preloadImage(resource.src)
      } else {
        return this.preloadScript(resource.src)
      }
    })
    
    return Promise.allSettled(promises)
  }
}

// 缓存管理器
export class CacheManager {
  private static cache: Map<string, {data: any, timestamp: number, ttl: number}> = new Map()
  
  // 设置缓存
  static set(key: string, data: any, ttl: number = 300000) { // 默认5分钟
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl
    })
  }
  
  // 获取缓存
  static get(key: string): any | null {
    const cached = this.cache.get(key)
    if (!cached) return null
    
    if (Date.now() - cached.timestamp > cached.ttl) {
      this.cache.delete(key)
      return null
    }
    
    return cached.data
  }
  
  // 删除缓存
  static delete(key: string) {
    this.cache.delete(key)
  }
  
  // 清空过期缓存
  static clearExpired() {
    const now = Date.now()
    for (const [key, value] of this.cache.entries()) {
      if (now - value.timestamp > value.ttl) {
        this.cache.delete(key)
      }
    }
  }
}

// 错误边界处理
export class ErrorBoundary {
  private static errors: Array<{error: Error, timestamp: number, context: any}> = []
  
  static captureError(error: Error, context?: any) {
    this.errors.push({
      error,
      timestamp: Date.now(),
      context
    })
    
    // 发送错误到监控服务器
    this.reportError(error, context)
  }
  
  private static async reportError(error: Error, context?: any) {
    try {
      await fetch('/api/error/report', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          message: error.message,
          stack: error.stack,
          timestamp: Date.now(),
          userAgent: navigator.userAgent,
          url: window.location.href,
          context
        })
      })
    } catch (e) {
      console.error('Failed to report error:', e)
    }
  }
  
  static getErrors() {
    return this.errors
  }
}

// 导出单例实例
export const performanceMonitor = PerformanceMonitor.getInstance()
export const lazyImageLoader = LazyImageLoader.getInstance()

// 初始化性能监控
if (typeof window !== 'undefined') {
  // 页面加载完成后记录性能指标
  window.addEventListener('load', () => {
    performanceMonitor.recordPageLoad()
    
    // 定期发送性能数据
    setInterval(() => {
      performanceMonitor.sendMetrics()
    }, 60000) // 每分钟发送一次
  })
  
  // 初始化图片懒加载
  lazyImageLoader.init()
  
  // 定期清理缓存
  setInterval(() => {
    CacheManager.clearExpired()
  }, 300000) // 每5分钟清理一次
}