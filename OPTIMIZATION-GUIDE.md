# 🎯 前端性能优化使用指南

## 快速开始

### 1. 查看优化总结
```bash
./optimization-summary.sh
```

### 2. 运行测试
```bash
./test-optimization.sh
```

### 3. 构建测试
```bash
chmod +x build-test.sh
./build-test.sh
```

## 开发流程

### 启动开发服务器
```bash
cd web
npm run dev
```

### 构建生产版本
```bash
cd web
npm run build
```

### 查看构建分析
```bash
cd web
npm run build:analyze
# 会在浏览器中打开可视化分析报告
```

### 预览生产构建
```bash
cd web
npm run preview
```

## 性能测试

### 使用 Lighthouse
```bash
# 确保预览服务器正在运行
npm run preview

# 在另一个终端运行
lighthouse http://localhost:4173 --view --output html --output-path ./lighthouse-report.html
```

### 使用 Chrome DevTools
1. 打开预览服务器：`npm run preview`
2. 在 Chrome 中打开 http://localhost:4173
3. 按 F12 打开 DevTools
4. 切换到 Lighthouse 标签
5. 点击 "Generate report"

## 优化要点说明

### 1. 路由懒加载
所有路由组件都使用懒加载，并且带有缓存机制。如果需要添加新路由：

```typescript
{
  path: 'new-page',
  name: 'NewPage',
  component: lazyLoadView('path/to/NewPage'),
  meta: {
    title: '新页面',
    keepAlive: true, // 是否缓存该页面
    preload: true    // 是否预加载
  }
}
```

### 2. 使用图片懒加载
```vue
<script setup>
import { useLazyLoad } from '@/utils/lazyLoad'

const imageRef = ref(null)
const { isVisible } = useLazyLoad(imageRef)
</script>

<template>
  <div ref="imageRef">
    <img v-if="isVisible" :src="imageSrc" />
  </div>
</template>
```

### 3. 网络自适应加载
```typescript
import { shouldLoadHighQuality } from '@/utils/lazyLoad'

const imageUrl = shouldLoadHighQuality() 
  ? 'high-quality.jpg' 
  : 'low-quality.jpg'
```

### 4. Service Worker 控制
在生产环境中，Service Worker 会自动注册。如果需要手动清理缓存：

```javascript
// 在浏览器控制台中执行
navigator.serviceWorker.controller?.postMessage({
  type: 'CLEAR_CACHE'
})
```

## 性能监控

应用已内置性能监控，在生产环境会自动收集以下指标：

- **LCP** (Largest Contentful Paint)
- **FID** (First Input Delay)
- **页面加载时间**

数据会在控制台输出，可以接入你的监控系统。

## 故障排查

### 构建失败
```bash
# 清理依赖重新安装
rm -rf node_modules package-lock.json
npm install
```

### Service Worker 问题
```bash
# 在 Chrome DevTools 中
# Application > Service Workers > Unregister
```

### 缓存问题
```bash
# 硬刷新：Ctrl+Shift+R (Windows/Linux) 或 Cmd+Shift+R (Mac)
# 或清除浏览器缓存
```

## 文档索引

- [优化总结](./OPTIMIZATION-SUMMARY.md) - 优化概览和成果
- [详细报告](./PERFORMANCE-OPTIMIZATION-V2.md) - 深度技术分析
- [前端优化](./FRONTEND-OPTIMIZATION-FINAL.md) - 前台页面优化

## 技术支持

如有问题，请检查：
1. Node.js 版本 >= 18.0.0
2. npm 版本 >= 9.0.0
3. 浏览器为最新版 Chrome/Firefox/Safari

---

**最后更新**: 2026-01-29
