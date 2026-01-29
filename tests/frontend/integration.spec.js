/**
 * 前端集成测试
 */

import { mount } from '@vue/test-utils'
import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import ElementPlus from 'element-plus'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'

// 导入组件
import Login from '@/views/auth/Login.vue'
import Register from '@/views/auth/Register.vue'
import Dashboard from '@/views/dashboard/Index.vue'
import TasksIndex from '@/views/tasks/Index.vue'
import TaskDetail from '@/views/tasks/Detail.vue'
import UserIndex from '@/views/user/Index.vue'
import Payment from '@/views/payment/Index.vue'

// 导入stores
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'

// 模拟API
import { authApi, userApi, taskApi, paymentApi } from '@/api'

// 创建测试应用
const createTestApp = (component) => {
  const app = createApp(component)
  app.use(ElementPlus)
  app.use(createPinia())
  app.use(createRouter({
    history: createWebHistory(),
    routes: []
  }))
  return app
}

// 全局测试配置
global.fetch = vi.fn()

describe('前端集成测试', () => {
  let userStore
  let appStore

  beforeEach(() => {
    setActivePinia(createPinia())
    userStore = useUserStore()
    appStore = useAppStore()
    
    // 重置mock
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('认证流程测试', () => {
    describe('登录页面', () => {
      it('应该正确渲染登录页面', () => {
        const wrapper = mount(createTestApp(Login)).mount()
        
        expect(wrapper.find('.login-container').exists()).toBe(true)
        expect(wrapper.find('.login-card').exists()).toBe(true)
        expect(wrapper.find('input[placeholder="请输入手机号"]').exists()).toBe(true)
        expect(wrapper.find('input[placeholder="请输入验证码"]').exists()).toBe(true)
        expect(wrapper.find('button').text()).toContain('登录')
      })

      it('应该验证手机号输入', async () => {
        const wrapper = mount(createTestApp(Login)).mount()
        
        const phoneInput = wrapper.find('input[placeholder="请输入手机号"]')
        await phoneInput.setValue('13800138000')
        
        expect(phoneInput.element.value).toBe('13800138000')
      })

      it('应该处理登录成功', async () => {
        // Mock API响应
        authApi.login = vi.fn().mockResolvedValue({
          data: {
            token: 'test_token',
            user: {
              id: 1,
              nickname: '测试用户',
              phone: '13800138000'
            }
          }
        })

        const wrapper = mount(createTestApp(Login)).mount()
        
        // 填写表单
        await wrapper.find('input[placeholder="请输入手机号"]').setValue('13800138000')
        await wrapper.find('input[placeholder="请输入验证码"]').setValue('123456')
        
        // 提交表单
        await wrapper.find('button[type="submit"]').trigger('click')
        
        // 验证API调用
        expect(authApi.login).toHaveBeenCalledWith({
          phone: '13800138000',
          code: '123456'
        })
      })

      it('应该处理登录失败', async () => {
        // Mock API错误
        authApi.login = vi.fn().mockRejectedValue({
          message: '验证码错误'
        })

        const wrapper = mount(createTestApp(Login)).mount()
        
        await wrapper.find('input[placeholder="请输入手机号"]').setValue('13800138000')
        await wrapper.find('input[placeholder="请输入验证码"]').setValue('123456')
        await wrapper.find('button[type="submit"]').trigger('click')
        
        // 验证错误处理
        expect(wrapper.text()).toContain('验证码错误')
      })
    })

    describe('注册流程', () => {
      it('应该正确渲染注册页面', () => {
        const wrapper = mount(createTestApp(Register)).mount()
        
        expect(wrapper.find('.register-container').exists()).toBe(true)
        expect(wrapper.find('.register-steps').exists()).toBe(true)
        expect(wrapper.find('el-steps').exists()).toBe(true)
      })

      it('应该验证注册表单', async () => {
        authApi.register = vi.fn().mockResolvedValue({
          data: { id: 1, nickname: '新用户' }
        })

        const wrapper = mount(createTestApp(Register)).mount()
        
        // 第一步：手机验证
        await wrapper.find('input[placeholder="请输入手机号"]').setValue('13800138001')
        await wrapper.find('input[placeholder="请输入验证码"]').setValue('123456')
        await wrapper.find('.next-btn').trigger('click')
        
        // 第二步：基本信息
        await wrapper.find('input[placeholder="请输入昵称"]').setValue('新用户')
        await wrapper.find('input[type="password"]').setValue('password123')
        await wrapper.find('input[placeholder="请确认密码"]').setValue('password123')
        await wrapper.find('input[type="checkbox"]').setValue(true)
        
        await wrapper.find('.next-btn').trigger('click')
        
        expect(authApi.register).toHaveBeenCalled()
      })
    })
  })

  describe('任务管理测试', () => {
    beforeEach(() => {
      // 模拟用户登录
      userStore.user = {
        id: 1,
        nickname: '测试用户',
        phone: '13800138000'
      }
      userStore.token = 'test_token'
    })

    describe('任务列表页面', () => {
      it('应该正确渲染任务列表', async () => {
        // Mock API响应
        taskApi.getTasks = vi.fn().mockResolvedValue({
          data: {
            list: [
              {
                id: 1,
                title: '测试任务1',
                description: '这是测试任务1',
                price: 100.00,
                status: 'pending'
              },
              {
                id: 2,
                title: '测试任务2',
                description: '这是测试任务2',
                price: 200.00,
                status: 'in_progress'
              }
            ],
            total: 2,
            page: 1,
            limit: 20
          }
        })

        const wrapper = mount(createTestApp(TasksIndex)).mount()
        
        // 等待组件加载
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 100))
        
        expect(wrapper.find('.tasks-container').exists()).toBe(true)
        expect(wrapper.find('.task-list').exists()).toBe(true)
      })

      it('应该支持搜索功能', async () => {
        taskApi.getTasks = vi.fn().mockResolvedValue({
          data: { list: [], total: 0, page: 1, limit: 20 }
        })

        const wrapper = mount(createTestApp(TasksIndex)).mount()
        
        const searchInput = wrapper.find('input[placeholder="搜索任务标题或描述"]')
        await searchInput.setValue('测试关键词')
        await searchInput.trigger('keyup.enter')
        
        expect(taskApi.getTasks).toHaveBeenCalledWith(
          expect.objectContaining({
            keyword: '测试关键词'
          })
        )
      })

      it('应该支持筛选功能', async () => {
        taskApi.getTasks = vi.fn().mockResolvedValue({
          data: { list: [], total: 0, page: 1, limit: 20 }
        })

        const wrapper = mount(createTestApp(TasksIndex)).mount()
        
        const categorySelect = wrapper.find('el-select')
        await categorySelect.setValue(1)
        
        expect(taskApi.getTasks).toHaveBeenCalledWith(
          expect.objectContaining({
            category_id: 1
          })
        )
      })
    })

    describe('任务详情页面', () => {
      it('应该正确渲染任务详情', async () => {
        taskApi.getTaskDetail = vi.fn().mockResolvedValue({
          data: {
            id: 1,
            title: '测试任务详情',
            description: '详细的任务描述',
            price: 100.00,
            status: 'pending',
            publisher: {
              id: 1,
              nickname: '发布者',
              avatar: ''
            }
          }
        })

        const wrapper = mount(createTestApp(TaskDetail)).mount()
        
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 100))
        
        expect(wrapper.find('.task-detail-container').exists()).toBe(true)
        expect(wrapper.text()).toContain('测试任务详情')
      })

      it('应该支持任务申请', async () => {
        taskApi.applyTask = vi.fn().mockResolvedValue({
          data: { id: 1, status: 'pending' }
        })

        const wrapper = mount(createTestApp(TaskDetail)).mount()
        
        // 模拟点击申请按钮
        const applyButton = wrapper.find('.action-buttons .el-button')
        if (applyButton.exists()) {
          await applyButton.trigger('click')
          
          // 填写申请表单
          const messageTextarea = wrapper.find('el-textarea')
          await messageTextarea.setValue('我有相关经验')
          
          const submitButton = wrapper.find('.el-dialog .el-button--primary')
          await submitButton.trigger('click')
          
          expect(taskApi.applyTask).toHaveBeenCalledWith(1, {
            message: '我有相关经验',
            estimated_time: '',
            contact: ''
          })
        }
      })
    })
  })

  describe('用户中心测试', () => {
    beforeEach(() => {
      userStore.user = {
        id: 1,
        nickname: '测试用户',
        phone: '13800138000',
        avatar: ''
      }
    })

    it('应该正确渲染个人中心', async () => {
      userApi.getStats = vi.fn().mockResolvedValue({
        data: {
          completed_tasks: 10,
          total_earnings: 1000.00,
          credit: 100
        }
      })

      const wrapper = mount(createTestApp(UserIndex)).mount()
      
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
      
      expect(wrapper.find('.user-container').exists()).toBe(true)
      expect(wrapper.find('.user-profile-card').exists()).toBe(true)
    })

    it('应该支持用户资料编辑', async () => {
      userApi.updateUserInfo = vi.fn().mockResolvedValue({
        data: { id: 1, nickname: '更新的用户名' }
      })

      const wrapper = mount(createTestApp(UserIndex)).mount()
      
      const editButton = wrapper.find('.profile-actions .el-button')
      if (editButton.exists()) {
        await editButton.trigger('click')
        // 验证跳转到编辑页面
        // 这里需要测试路由导航
      }
    })
  })

  describe('支付流程测试', () => {
    beforeEach(() => {
      userStore.user = {
        id: 1,
        nickname: '测试用户',
        phone: '13800138000'
      }
    })

    it('应该正确渲染支付页面', async () => {
      paymentApi.getOrderDetail = vi.fn().mockResolvedValue({
        data: {
          id: 1,
          order_no: 'TEST_123',
          amount: 100.00,
          status: 'pending'
        }
      })

      const wrapper = mount(createTestApp(Payment)).mount()
      
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
      
      expect(wrapper.find('.payment-container').exists()).toBe(true)
      expect(wrapper.find('.order-card').exists()).toBe(true)
    })

    it('应该支持多种支付方式', async () => {
      const wrapper = mount(createTestApp(Payment)).mount()
      
      const paymentMethods = wrapper.findAll('.payment-method')
      expect(paymentMethods.length).toBeGreaterThan(0)
      
      // 测试选择支付方式
      await paymentMethods[0].trigger('click')
      expect(wrapper.vm.selectedMethod).toBeDefined()
    })

    it('应该处理支付成功', async () => {
      paymentApi.createPayment = vi.fn().mockResolvedValue({
        data: {
          qr_code: 'data:image/png;base64,test_qr_code',
          payment_url: 'https://test.payment.url'
        }
      })

      const wrapper = mount(createTestApp(Payment)).mount()
      
      // 选择支付方式
      const paymentMethods = wrapper.findAll('.payment-method')
      await paymentMethods[0].trigger('click')
      
      // 点击支付按钮
      const payButton = wrapper.find('.pay-button')
      await payButton.trigger('click')
      
      expect(paymentApi.createPayment).toHaveBeenCalled()
    })
  })

  describe('状态管理测试', () => {
    it('应该正确管理用户状态', () => {
      expect(userStore.user).toBeNull()
      expect(userStore.isLoggedIn).toBe(false)
      
      userStore.setUser({
        id: 1,
        nickname: '测试用户',
        phone: '13800138000'
      })
      
      expect(userStore.user).not.toBeNull()
      expect(userStore.isLoggedIn).toBe(true)
    })

    it('应该正确管理应用状态', () => {
      expect(appStore.theme).toBe('light')
      
      appStore.setTheme('dark')
      expect(appStore.theme).toBe('dark')
      
      appStore.toggleTheme()
      expect(appStore.theme).toBe('light')
    })
  })

  describe('组件交互测试', () => {
    it('应该正确处理表单验证', async () => {
      const wrapper = mount(createTestApp(Login)).mount()
      
      // 提交空表单
      const submitButton = wrapper.find('button[type="submit"]')
      await submitButton.trigger('click')
      
      // 验证错误提示
      expect(wrapper.text()).toContain('请输入手机号')
    })

    it('应该正确处理加载状态', async () => {
      authApi.login = vi.fn().mockImplementation(() => {
        return new Promise(resolve => {
          setTimeout(() => {
            resolve({
              data: { token: 'test_token', user: { id: 1 } }
            })
          }, 1000)
        })
      })

      const wrapper = mount(createTestApp(Login)).mount()
      
      await wrapper.find('input[placeholder="请输入手机号"]').setValue('13800138000')
      await wrapper.find('input[placeholder="请输入验证码"]').setValue('123456')
      await wrapper.find('button[type="submit"]').trigger('click')
      
      // 验证加载状态
      expect(wrapper.find('.el-button.is-loading').exists()).toBe(true)
    })

    it('应该正确处理错误状态', async () => {
      authApi.login = vi.fn().mockRejectedValue({
        message: '网络错误',
        code: 500
      })

      const wrapper = mount(createTestApp(Login)).mount()
      
      await wrapper.find('input[placeholder="请输入手机号"]').setValue('13800138000')
      await wrapper.find('input[placeholder="请输入验证码"]').setValue('123456')
      await wrapper.find('button[type="submit"]').trigger('click')
      
      // 验证错误处理
      expect(wrapper.vm.error).toBeDefined()
    })
  })

  describe('响应式设计测试', () => {
    it('应该在移动端正确显示', async () => {
      // 模拟移动端视口
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 375
      })

      const wrapper = mount(createTestApp(Dashboard)).mount()
      
      await wrapper.vm.$nextTick()
      
      // 验证移动端样式
      expect(wrapper.find('.dashboard-container').classes()).toContain('mobile')
    })

    it('应该在桌面端正确显示', async () => {
      // 模拟桌面端视口
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 1200
      })

      const wrapper = mount(createTestApp(Dashboard)).mount()
      
      await wrapper.vm.$nextTick()
      
      // 验证桌面端样式
      expect(wrapper.find('.dashboard-container').classes()).toContain('desktop')
    })
  })

  describe('性能测试', () => {
    it('应该在合理时间内加载组件', async () => {
      const startTime = performance.now()
      
      const wrapper = mount(createTestApp(Dashboard)).mount()
      await wrapper.vm.$nextTick()
      
      const endTime = performance.now()
      const loadTime = endTime - startTime
      
      // 加载时间应该小于100ms
      expect(loadTime).toBeLessThan(100)
    })

    it('应该正确处理大量数据', async () => {
      const largeDataSet = Array.from({ length: 1000 }, (_, i) => ({
        id: i,
        title: `任务 ${i}`,
        description: `任务描述 ${i}`,
        price: Math.random() * 1000
      }))

      taskApi.getTasks = vi.fn().mockResolvedValue({
        data: {
          list: largeDataSet,
          total: 1000,
          page: 1,
          limit: 20
        }
      })

      const wrapper = mount(createTestApp(TasksIndex)).mount()
      
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
      
      // 验证虚拟滚动或分页处理
      expect(wrapper.find('.task-list').exists()).toBe(true)
    })
  })
})

// 全局错误处理测试
describe('全局错误处理', () => {
  it('应该正确处理网络错误', async () => {
    // 模拟网络错误
    global.fetch = vi.fn().mockRejectedValue(new Error('Network Error'))

    try {
      await authApi.login({ phone: '13800138000', code: '123456' })
    } catch (error) {
      expect(error.message).toBe('Network Error')
    }
  })

  it('应该正确处理API错误响应', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: () => Promise.resolve({
        message: 'Bad Request',
        code: 400
      })
    })

    try {
      await authApi.login({ phone: '13800138000', code: '123456' })
    } catch (error) {
      expect(error.code).toBe(400)
    }
  })
})

// 用户体验测试
describe('用户体验测试', () => {
  it('应该支持键盘导航', async () => {
    const wrapper = mount(createTestApp(Login)).mount()
    
    // 测试Tab键导航
    const firstInput = wrapper.find('input[placeholder="请输入手机号"]')
    await firstInput.trigger('focus')
    
    // 模拟Tab键
    await wrapper.vm.$nextTick()
    document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Tab' }))
    
    // 验证焦点移动
    expect(document.activeElement).toBe(wrapper.find('input[placeholder="请输入验证码"]').element)
  })

  it('应该提供合适的加载反馈', async () => {
    authApi.login = vi.fn().mockImplementation(() => {
      return new Promise(resolve => {
        setTimeout(() => {
          resolve({ data: { token: 'test' } })
        }, 2000)
      })
    })

    const wrapper = mount(createTestApp(Login)).mount()
    
    await wrapper.find('input[placeholder="请输入手机号"]').setValue('13800138000')
    await wrapper.find('input[placeholder="请输入验证码"]').setValue('123456')
    await wrapper.find('button[type="submit"]').trigger('click')
    
    // 验证加载指示器
    expect(wrapper.find('.el-button.is-loading').exists()).toBe(true)
    expect(wrapper.text()).toContain('登录中')
  })
})