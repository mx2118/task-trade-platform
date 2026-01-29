<template>
  <div ref="containerRef" class="lazy-image-wrapper" :style="wrapperStyle">
    <!-- 占位符 -->
    <div v-if="!isLoaded && !hasError" class="lazy-image-placeholder">
      <div class="placeholder-shimmer"></div>
    </div>
    
    <!-- 实际图片 -->
    <img
      v-show="isLoaded"
      :src="currentSrc"
      :alt="alt"
      class="lazy-image"
      :class="{ 'fade-in': isLoaded }"
      @load="handleLoad"
      @error="handleError"
    />
    
    <!-- 错误状态 -->
    <div v-if="hasError" class="lazy-image-error">
      <el-icon><Picture /></el-icon>
      <span>加载失败</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Picture } from '@element-plus/icons-vue'

interface Props {
  src: string
  alt?: string
  width?: string | number
  height?: string | number
  lazy?: boolean
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  alt: '',
  lazy: true,
  placeholder: ''
})

const containerRef = ref<HTMLElement | null>(null)
const isLoaded = ref(false)
const hasError = ref(false)
const isIntersecting = ref(false)

const currentSrc = computed(() => {
  if (!props.lazy || isIntersecting.value) {
    return props.src
  }
  return props.placeholder || ''
})

const wrapperStyle = computed(() => {
  const style: Record<string, string> = {}
  if (props.width) {
    style.width = typeof props.width === 'number' ? `${props.width}px` : props.width
  }
  if (props.height) {
    style.height = typeof props.height === 'number' ? `${props.height}px` : props.height
  }
  return style
})

const handleLoad = () => {
  isLoaded.value = true
}

const handleError = () => {
  hasError.value = true
}

let observer: IntersectionObserver | null = null

onMounted(() => {
  if (!props.lazy) {
    isIntersecting.value = true
    return
  }

  if (!containerRef.value) return

  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          isIntersecting.value = true
          // 加载后停止观察
          if (observer) {
            observer.disconnect()
          }
        }
      })
    },
    {
      rootMargin: '100px',
      threshold: 0.01
    }
  )

  observer.observe(containerRef.value)
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect()
  }
})
</script>

<style lang="scss" scoped>
.lazy-image-wrapper {
  position: relative;
  overflow: hidden;
  background-color: #f5f7fa;
  display: inline-block;
  
  .lazy-image-placeholder {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    
    .placeholder-shimmer {
      width: 100%;
      height: 100%;
      background: linear-gradient(
        90deg,
        #f5f7fa 25%,
        #ebeef5 50%,
        #f5f7fa 75%
      );
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }
  }
  
  .lazy-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    
    &.fade-in {
      animation: fadeIn 0.3s ease-in;
    }
  }
  
  .lazy-image-error {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #c0c4cc;
    font-size: 14px;
    
    .el-icon {
      font-size: 32px;
      margin-bottom: 8px;
    }
  }
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
