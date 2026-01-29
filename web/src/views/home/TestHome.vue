<template>
  <div class="test-home">
    <div class="hero">
      <h1>🚀 任务交易平台</h1>
      <p>测试页面 - Vue组件正常渲染</p>
      <p>当前时间: {{ currentTime }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const currentTime = ref('')

onMounted(() => {
  currentTime.value = new Date().toLocaleString('zh-CN')
  console.log('TestHome组件已挂载')
})
</script>

<style lang="scss" scoped>
.test-home {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  
  .hero {
    text-align: center;
    color: white;
    padding: 40px;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    
    h1 {
      font-size: 48px;
      margin: 0 0 20px 0;
    }
    
    p {
      font-size: 20px;
      margin: 10px 0;
    }
  }
}
</style>
