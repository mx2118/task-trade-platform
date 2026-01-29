<template>
  <div class="minimal-home">
    <div class="container">
      <h1>🚀 任务交易平台</h1>
      <p>欢迎来到任务交易平台</p>
      <p>当前时间: {{ currentTime }}</p>
      
      <div class="buttons">
        <button @click="handleClick">点击测试</button>
      </div>
      
      <div v-if="message" class="message">{{ message }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const currentTime = ref('')
const message = ref('')

const handleClick = () => {
  message.value = '按钮点击成功！'
  console.log('按钮被点击了')
}

onMounted(() => {
  currentTime.value = new Date().toLocaleString('zh-CN')
  console.log('MinimalHome 组件已挂载')
})
</script>

<style scoped>
.minimal-home {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.container {
  text-align: center;
  color: white;
  padding: 40px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  max-width: 600px;
}

h1 {
  font-size: 48px;
  margin: 0 0 20px 0;
}

p {
  font-size: 18px;
  margin: 10px 0;
}

.buttons {
  margin: 30px 0;
}

button {
  padding: 12px 32px;
  font-size: 16px;
  background: white;
  color: #667eea;
  border: none;
  border-radius: 24px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s;
}

button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 255, 255, 0.3);
}

.message {
  margin-top: 20px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  font-size: 16px;
}
</style>
