/**
 * 图片优化工具
 */

// 图片懒加载指令
export const lazyImageDirective = {
  mounted(el: HTMLImageElement) {
    // 添加 loading="lazy" 属性
    el.loading = 'lazy'
    
    // 如果浏览器不支持原生懒加载，使用 Intersection Observer
    if (!('loading' in HTMLImageElement.prototype)) {
      const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            const img = entry.target as HTMLImageElement
            const src = img.dataset.src
            if (src) {
              img.src = src
              observer.unobserve(img)
            }
          }
        })
      }, { rootMargin: '50px' })
      
      observer.observe(el)
    }
  }
}

// 图片预加载
export function preloadImage(url: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => resolve()
    img.onerror = reject
    img.src = url
  })
}

// 批量预加载图片
export function preloadImages(urls: string[]): Promise<void[]> {
  return Promise.all(urls.map(preloadImage))
}
