/**
 * 性能测试脚本
 */

import http from 'k6/http'
import { check, sleep } from 'k6'
import { Rate } from 'k6/metrics'

// 自定义指标
const errorRate = new Rate('errors')

// 测试配置
export const options = {
  stages: [
    { duration: '2m', target: 100 }, // 2分钟内逐渐增加到100用户
    { duration: '5m', target: 100 }, // 保持100用户5分钟
    { duration: '2m', target: 200 }, // 2分钟内逐渐增加到200用户
    { duration: '5m', target: 200 }, // 保持200用户5分钟
    { duration: '2m', target: 0 },   // 2分钟内逐渐减少到0用户
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95%的请求在500ms内完成
    http_req_failed: ['rate<0.1'],     // 错误率低于10%
    errors: ['rate<0.1'],             // 自定义错误率低于10%
  },
}

const BASE_URL = 'http://localhost:8080'

// 测试数据
const testUsers = [
  { phone: '13800138001', code: '123456' },
  { phone: '13800138002', code: '123456' },
  { phone: '13800138003', code: '123456' },
  { phone: '13800138004', code: '123456' },
  { phone: '13800138005', code: '123456' },
]

const testTasks = [
  {
    title: '性能测试任务1',
    description: '这是一个性能测试任务',
    category_id: 1,
    price: 100,
    deadline: '2024-12-31 23:59:59',
  },
  {
    title: '性能测试任务2',
    description: '这是另一个性能测试任务',
    category_id: 2,
    price: 200,
    deadline: '2024-12-31 23:59:59',
  },
]

export function setup() {
  // 预热：创建一些测试数据
  console.log('开始预热阶段...')
  
  // 登录测试用户
  testUsers.forEach(user => {
    const loginResponse = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify(user), {
      headers: { 'Content-Type': 'application/json' },
    })
    
    if (loginResponse.status === 200) {
      console.log(`用户 ${user.phone} 登录成功`)
      user.token = loginResponse.json('data.token')
    }
  })
  
  console.log('预热完成')
}

export default function () {
  // 随机选择测试场景
  const scenarios = [
    testHomePage,
    testUserLogin,
    testTaskList,
    testTaskDetail,
    testCreateTask,
    testUserDashboard,
    testPaymentFlow,
    testAPIPerformance,
  ]
  
  const scenario = scenarios[Math.floor(Math.random() * scenarios.length)]
  scenario()
}

// 测试场景函数

function testHomePage() {
  const response = http.get(`${BASE_URL}/api/v1/system/info`)
  
  check(response, {
    '首页加载状态为200': (r) => r.status === 200,
    '首页响应时间<200ms': (r) => r.timings.duration < 200,
  })
}

function testUserLogin() {
  const user = testUsers[Math.floor(Math.random() * testUsers.length)]
  
  const response = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify(user), {
    headers: { 'Content-Type': 'application/json' },
  })
  
  const success = check(response, {
    '登录状态为200': (r) => r.status === 200,
    '登录响应时间<500ms': (r) => r.timings.duration < 500,
    '返回token': (r) => r.json('data.token') !== undefined,
  })
  
  errorRate.add(!success)
}

function testTaskList() {
  const params = {
    page: Math.floor(Math.random() * 10) + 1,
    limit: Math.floor(Math.random() * 20) + 10,
    keyword: ['设计', '开发', '写作', '翻译'][Math.floor(Math.random() * 4)],
  }
  
  const response = http.get(`${BASE_URL}/api/v1/tasks?${Object.entries(params).map(([k, v]) => `${k}=${v}`).join('&')}`)
  
  check(response, {
    '任务列表状态为200': (r) => r.status === 200,
    '任务列表响应时间<300ms': (r) => r.timings.duration < 300,
    '返回数据格式正确': (r) => {
      const data = r.json()
      return data.data && Array.isArray(data.data.list)
    },
  })
}

function testTaskDetail() {
  // 先获取任务列表
  const listResponse = http.get(`${BASE_URL}/api/v1/tasks?limit=10`)
  
  if (listResponse.status === 200) {
    const tasks = listResponse.json('data.list')
    if (tasks && tasks.length > 0) {
      const taskId = tasks[Math.floor(Math.random() * tasks.length)].id
      
      const response = http.get(`${BASE_URL}/api/v1/tasks/${taskId}`)
      
      check(response, {
        '任务详情状态为200': (r) => r.status === 200,
        '任务详情响应时间<200ms': (r) => r.timings.duration < 200,
      })
    }
  }
}

function testCreateTask() {
  const user = testUsers.find(u => u.token)
  if (!user) return
  
  const task = testTasks[Math.floor(Math.random() * testTasks.length)]
  
  const response = http.post(`${BASE_URL}/api/v1/tasks`, JSON.stringify(task), {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${user.token}`,
    },
  })
  
  check(response, {
    '创建任务状态为200或400': (r) => [200, 400].includes(r.status),
    '创建任务响应时间<1000ms': (r) => r.timings.duration < 1000,
  })
}

function testUserDashboard() {
  const user = testUsers.find(u => u.token)
  if (!user) return
  
  const response = http.get(`${BASE_URL}/api/v1/user/stats`, {
    headers: {
      'Authorization': `Bearer ${user.token}`,
    },
  })
  
  check(response, {
    '用户统计状态为200': (r) => r.status === 200,
    '用户统计响应时间<300ms': (r) => r.timings.duration < 300,
  })
}

function testPaymentFlow() {
  const user = testUsers.find(u => u.token)
  if (!user) return
  
  // 创建预支付订单
  const paymentData = {
    order_type: 'recharge',
    amount: 100,
    payment_method: 'alipay',
  }
  
  const response = http.post(`${BASE_URL}/api/v1/pay/prepay`, JSON.stringify(paymentData), {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${user.token}`,
    },
  })
  
  check(response, {
    '预支付状态为200或400': (r) => [200, 400].includes(r.status),
    '预支付响应时间<1500ms': (r) => r.timings.duration < 1500,
  })
}

function testAPIPerformance() {
  // 测试API响应性能
  const apis = [
    { path: '/api/v1/categories', name: '分类列表' },
    { path: '/api/v1/system/announce', name: '系统公告' },
    { path: '/api/v1/system/info', name: '系统信息' },
  ]
  
  apis.forEach(api => {
    const response = http.get(`${BASE_URL}${api.path}`)
    
    check(response, {
      [`${api.name}状态为200`]: (r) => r.status === 200,
      [`${api.name}响应时间<200ms`]: (r) => r.timings.duration < 200,
    })
  })
}

export function teardown() {
  console.log('测试完成，清理数据...')
  
  // 这里可以添加清理逻辑
  testUsers.forEach(user => {
    if (user.token) {
      http.post(`${BASE_URL}/api/v1/auth/logout`, '', {
        headers: {
          'Authorization': `Bearer ${user.token}`,
        },
      })
    }
  })
}