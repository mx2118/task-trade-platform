/**
 * 端到端用户旅程测试
 */

const { test, expect } = require('@playwright/test')

test.describe('完整用户旅程测试', () => {
  test.beforeEach(async ({ page }) => {
    // 设置测试环境
    await page.goto('http://localhost:3000')
    
    // 清除存储
    await page.evaluate(() => {
      localStorage.clear()
      sessionStorage.clear()
    })
  })

  test('新用户注册到完成任务的完整流程', async ({ page }) => {
    // 1. 访问首页
    await page.goto('http://localhost:3000')
    await expect(page).toHaveTitle(/任务交易平台/)
    
    // 2. 点击注册
    await page.click('text=注册新用户')
    await expect(page.locator('.register-container')).toBeVisible()
    
    // 3. 第一步：手机验证
    await page.fill('input[placeholder="请输入手机号"]', '13800138000')
    await page.click('text=获取验证码')
    await expect(page.locator('.el-button:has-text("60s后重试")')).toBeVisible()
    
    await page.fill('input[placeholder="请输入验证码"]', '123456')
    await page.click('text=下一步')
    
    // 4. 第二步：基本信息
    await page.fill('input[placeholder="请输入昵称"]', '测试用户')
    await page.fill('input[type="password"]', 'test123456')
    await page.fill('input[placeholder="请确认密码"]', 'test123456')
    await page.click('input[type="checkbox"]')
    await page.click('text=完成注册')
    
    // 5. 验证注册成功
    await expect(page.locator('.success-content')).toBeVisible()
    await expect(page.locator('text=注册成功')).toBeVisible()
    
    // 6. 点击开始使用
    await page.click('text=开始使用')
    await expect(page).toHaveURL('http://localhost:3000/dashboard')
    
    // 7. 发布第一个任务
    await page.click('text=发布任务')
    await expect(page).toHaveURL('http://localhost:3000/tasks/publish')
    
    await page.fill('input[name="title"]', '测试任务：设计一个Logo')
    await page.fill('textarea[name="description"]', '需要设计一个现代化的公司Logo，要求简洁、专业、易识别')
    await page.fill('textarea[name="requirements"]', '1. 提供至少3个初稿\n2. 支持修改2次\n3. 提供源文件')
    await page.selectOption('select[name="category_id"]', '设计')
    await page.fill('input[name="price"]', '500')
    await page.fill('input[name="deadline"]', '2024-12-31 23:59:59')
    await page.fill('input[name="location"]', '北京市朝阳区')
    await page.check('input[name="is_urgent"]')
    
    await page.click('text=发布任务')
    await expect(page.locator('.el-message:has-text("发布成功")')).toBeVisible()
    
    // 8. 查看发布的任务
    await page.goto('http://localhost:3000/user')
    await expect(page.locator('.user-container')).toBeVisible()
    
    await page.click('text=我的任务')
    await expect(page.locator('.task-item')).toHaveCount(1)
    await expect(page.locator('.task-item')).toContainText('测试任务：设计一个Logo')
    
    // 9. 充值账户
    await page.click('text=我的钱包')
    await page.click('text=充值')
    
    await page.fill('input[name="amount"]', '1000')
    await page.selectOption('select[name="payment_method"]', 'alipay')
    await page.click('text=立即充值')
    
    // 10. 查看其他任务
    await page.goto('http://localhost:3000/tasks')
    await expect(page.locator('.tasks-container')).toBeVisible()
    
    // 等待任务列表加载
    await page.waitForSelector('.task-item')
    const taskItems = await page.locator('.task-item').count()
    expect(taskItems).toBeGreaterThan(0)
    
    // 11. 申请一个任务
    const firstTask = page.locator('.task-item').first()
    await firstTask.click()
    
    await expect(page).toHaveURL(/\/tasks\/\d+/)
    await page.click('text=申请接取')
    
    // 填写申请信息
    await page.fill('textarea[name="message"]', '我有丰富的设计经验，可以高质量完成这个任务')
    await page.fill('input[name="estimated_time"]', '3天')
    await page.fill('input[name="contact"]', '13800138001')
    
    await page.click('text=提交申请')
    await expect(page.locator('.el-message:has-text("申请已提交")')).toBeVisible()
    
    // 12. 查看申请记录
    await page.goto('http://localhost:3000/user')
    await page.click('text=我的任务')
    await page.click('text=已申请')
    
    await expect(page.locator('.task-item')).toHaveCount(1)
    
    // 13. 查看个人统计
    await page.click('text=返回首页')
    await expect(page.locator('.stats-cards')).toBeVisible()
    
    const completedTasks = await page.locator('.stats-card:has-text("完成任务") .stat-value').textContent()
    expect(completedTasks).toBe('0')
    
    const publishedTasks = await page.locator('.stats-card:has-text("发布任务") .stat-value').textContent()
    expect(publishedTasks).toBe('1')
  })

  test('用户登录到支付任务的完整流程', async ({ page }) => {
    // 1. 登录
    await page.goto('http://localhost:3000/login')
    
    await page.fill('input[placeholder="请输入手机号"]', '13800138001')
    await page.fill('input[placeholder="请输入验证码"]', '123456')
    await page.click('text=登录')
    
    await expect(page).toHaveURL('http://localhost:3000/dashboard')
    
    // 2. 查看钱包余额
    await page.goto('http://localhost:3000/user')
    await page.click('text=我的钱包')
    
    await expect(page.locator('.wallet-info')).toBeVisible()
    
    // 3. 充值
    await page.click('text=充值')
    await page.fill('input[name="amount"]', '100')
    await page.selectOption('select[name="payment_method"]', 'alipay')
    await page.click('text=立即充值')
    
    // 4. 完成支付流程
    await expect(page).toHaveURL(/\/payment/)
    await expect(page.locator('.payment-container')).toBeVisible()
    
    await page.click('.payment-method')
    await page.click('text=立即支付')
    
    // 模拟支付成功
    await page.waitForTimeout(2000)
    await expect(page.locator('.payment-result')).toBeVisible()
    await expect(page.locator('text=支付成功')).toBeVisible()
    
    await page.click('text=查看订单')
    
    // 5. 查看交易记录
    await page.goto('http://localhost:3000/user')
    await page.click('text=收支明细')
    
    await expect(page.locator('.transaction-list')).toBeVisible()
    await expect(page.locator('.transaction-item')).toHaveCount.greaterThan(0)
  })

  test('任务发布到完成结算的完整流程', async ({ page }) => {
    // 1. 登录并发布任务
    await page.goto('http://localhost:3000/login')
    await page.fill('input[placeholder="请输入手机号"]', '13800138002')
    await page.fill('input[placeholder="请输入验证码"]', '123456')
    await page.click('text=登录')
    
    await page.click('text=发布任务')
    
    await page.fill('input[name="title"]', '文案撰写任务')
    await page.fill('textarea[name="description"]', '需要为新产品撰写宣传文案')
    await page.fill('textarea[name="requirements"]', '1. 突出产品特点\n2. 吸引目标用户\n3. 500字左右')
    await page.selectOption('select[name="category_id"]', '写作')
    await page.fill('input[name="price"]', '200')
    await page.fill('input[name="deadline"]', '2024-12-15 23:59:59')
    
    await page.click('text=发布任务')
    
    // 2. 支付任务费用
    await expect(page).toHaveURL(/\/payment/)
    await page.click('.payment-method')
    await page.click('text=立即支付')
    
    await page.waitForTimeout(2000)
    await page.click('text=查看订单')
    
    // 3. 切换到接取者账号
    await page.goto('http://localhost:3000/login')
    await page.fill('input[placeholder="请输入手机号"]', '13800138003')
    await page.fill('input[placeholder="请输入验证码"]', '123456')
    await page.click('text=登录')
    
    // 4. 查找并申请任务
    await page.goto('http://localhost:3000/tasks')
    await page.fill('input[placeholder="搜索任务标题或描述"]', '文案撰写')
    await page.click('text=搜索')
    
    const taskItem = page.locator('.task-item:has-text("文案撰写任务")').first()
    await taskItem.click()
    
    await page.click('text=申请接取')
    await page.fill('textarea[name="message"]', '我是专业文案撰写者，有5年经验')
    await page.fill('input[name="estimated_time"]', '2天')
    await page.fill('input[name="contact"]', '13800138003')
    await page.click('text=提交申请')
    
    // 5. 切换回发布者账号接受申请
    await page.goto('http://localhost:3000/login')
    await page.fill('input[placeholder="请输入手机号"]', '13800138002')
    await page.fill('input[placeholder="请输入验证码"]', '123456')
    await page.click('text=登录')
    
    await page.goto('http://localhost:3000/tasks')
    await page.click('a[href="/tasks"]:has-text("我的任务")')
    await page.click('text=文案撰写任务')
    
    await page.click('text=接受申请')
    
    // 6. 接取者交付任务
    await page.goto('http://localhost:3000/login')
    await page.fill('input[placeholder="请输入手机号"]', '13800138003')
    await page.fill('input[placeholder="请输入验证码"]', '123456')
    await page.click('text=登录')
    
    await page.goto('http://localhost:3000/tasks')
    await page.click('text=进行中')
    
    const taskToDeliver = page.locator('.task-item:has-text("文案撰写任务")')
    await taskToDeliver.click()
    
    await page.click('text=交付任务')
    await page.fill('textarea[name="delivery_content"]', '文案已完成，请查看附件')
    await page.click('text=提交交付')
    
    // 7. 发布者确认完成
    await page.goto('http://localhost:3000/login')
    await page.fill('input[placeholder="请输入手机号"]', '13800138002')
    await page.fill('input[placeholder="请输入验证码"]', '123456')
    await page.click('text=登录')
    
    await page.goto('http://localhost:3000/tasks')
    await page.click('text=已交付')
    
    const taskToConfirm = page.locator('.task-item:has-text("文案撰写任务")')
    await taskToConfirm.click()
    
    await page.click('text=确认完成')
    
    // 8. 评价任务
    await page.fill('el-rate', '5')
    await page.fill('textarea[name="comment"]', '任务完成得很好，推荐合作')
    await page.click('text=提交评价')
    
    // 9. 查看收益
    await page.goto('http://localhost:3000/user')
    await page.click('text=我的钱包')
    
    await expect(page.locator('.balance-amount')).toContainText('180') // 200 - 10%平台费用
    await page.click('text=收支明细')
    
    await expect(page.locator('.transaction-item:has-text("收入")')).toBeVisible()
  })

  test('用户体验和性能测试', async ({ page }) => {
    // 1. 测试页面加载性能
    const startTime = Date.now()
    await page.goto('http://localhost:3000')
    await page.waitForLoadState('networkidle')
    const loadTime = Date.now() - startTime
    
    expect(loadTime).toBeLessThan(3000) // 页面加载时间应小于3秒
    
    // 2. 测试响应式设计
    await page.setViewportSize({ width: 375, height: 667 }) // 移动端
    await expect(page.locator('.dashboard-container')).toBeVisible()
    
    await page.setViewportSize({ width: 1200, height: 800 }) // 桌面端
    await expect(page.locator('.sidebar')).toBeVisible()
    
    // 3. 测试键盘导航
    await page.keyboard.press('Tab')
    await expect(page.locator(':focus')).toBeVisible()
    
    // 4. 测试搜索功能性能
    await page.goto('http://localhost:3000/tasks')
    await page.fill('input[placeholder="搜索任务标题或描述"]', '设计')
    
    const searchStartTime = Date.now()
    await page.waitForSelector('.task-list')
    const searchTime = Date.now() - searchStartTime
    
    expect(searchTime).toBeLessThan(1000) // 搜索响应时间应小于1秒
    
    // 5. 测试无限滚动
    await page.scrollTo('bottom')
    await page.waitForTimeout(1000)
    
    const initialTaskCount = await page.locator('.task-item').count()
    await page.scrollTo('bottom')
    await page.waitForTimeout(1000)
    
    const finalTaskCount = await page.locator('.task-item').count()
    expect(finalTaskCount).toBeGreaterThan(initialTaskCount)
    
    // 6. 测试错误处理
    await page.goto('http://localhost:3000/non-existent-page')
    await expect(page.locator('.error-page')).toBeVisible()
    await expect(page.locator('text=页面未找到')).toBeVisible()
    
    // 7. 测试网络错误处理
    await page.route('**/api/**', route => route.abort())
    await page.goto('http://localhost:3000/tasks')
    await expect(page.locator('.error-message')).toBeVisible()
    await page.unroute('**/api/**')
  })

  test('跨浏览器兼容性测试', async ({ page, browserName }) => {
    await page.goto('http://localhost:3000')
    
    // 基础功能测试
    await expect(page).toHaveTitle(/任务交易平台/)
    await expect(page.locator('.dashboard-container')).toBeVisible()
    
    // 测试现代JavaScript特性
    const modernJS = await page.evaluate(() => {
      return {
        fetch: typeof fetch !== 'undefined',
        promise: typeof Promise !== 'undefined',
        arrow: (() => { try { eval('() => {}'); return true; } catch { return false; } })(),
        classes: typeof class {} !== 'undefined'
      }
    })
    
    expect(modernJS.fetch).toBe(true)
    expect(modernJS.promise).toBe(true)
    expect(modernJS.arrow).toBe(true)
    expect(modernJS.classes).toBe(true)
    
    // 测试CSS特性
    const cssSupport = await page.evaluate(() => {
      const testDiv = document.createElement('div')
      testDiv.style.display = 'flex'
      return testDiv.style.display === 'flex'
    })
    
    expect(cssSupport).toBe(true)
    
    console.log(`Browser ${browserName} compatibility test passed`)
  })
})

// 压力测试
test.describe('系统压力测试', () => {
  test('并发用户访问测试', async ({ browser }) => {
    const concurrentUsers = 10
    const promises = []
    
    for (let i = 0; i < concurrentUsers; i++) {
      const context = await browser.newContext()
      const page = await context.newPage()
      
      promises.push(
        page.goto('http://localhost:3000')
          .then(() => page.waitForLoadState('networkidle'))
          .then(() => {
            return page.locator('.dashboard-container').isVisible()
          })
          .finally(() => context.close())
      )
    }
    
    const results = await Promise.all(promises)
    expect(results.every(visible => visible === true)).toBe(true)
    
    console.log(`Concurrent users test passed: ${concurrentUsers} users`)
  })

  test('大数据量处理测试', async ({ page }) => {
    // 模拟大量任务数据
    await page.goto('http://localhost:3000/tasks')
    
    // 测试虚拟滚动性能
    const scrollTest = async () => {
      const startTime = Date.now()
      
      for (let i = 0; i < 10; i++) {
        await page.scrollTo('bottom')
        await page.waitForTimeout(500)
      }
      
      const totalTime = Date.now() - startTime
      return totalTime
    }
    
    const scrollTime = await scrollTest()
    expect(scrollTime).toBeLessThan(10000) // 10次滚动应在10秒内完成
    
    console.log(`Large data handling test passed: ${scrollTime}ms for 10 scrolls`)
  })
})

// 安全性测试
test.describe('安全性测试', () => {
  test('XSS防护测试', async ({ page }) => {
    await page.goto('http://localhost:3000/login')
    
    // 尝试输入恶意脚本
    await page.fill('input[placeholder="请输入手机号"]', '<script>alert("XSS")</script>')
    await page.click('text=登录')
    
    // 验证脚本没有执行
    await page.waitForTimeout(2000)
    await expect(page.locator('.login-container')).toBeVisible()
    // 如果XSS成功，会有alert弹窗，页面会失去焦点
    expect(await page.locator('body').isFocused()).toBe(true)
  })

  test('CSRF防护测试', async ({ page }) => {
    // 模拟跨站请求
    await page.goto('http://localhost:3000/login')
    await page.evaluate(() => {
      fetch('http://localhost:8080/api/v1/user/update', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          nickname: 'hacked_user'
        })
      })
    })
    
    await page.waitForTimeout(2000)
    
    // 验证请求被拒绝
    const response = await page.evaluate(async () => {
      try {
        const resp = await fetch('http://localhost:3000/api/v1/user/profile')
        return resp.status
      } catch {
        return 500
      }
    })
    
    expect(response).toBe(401) // 未授权
  })

  test('输入验证测试', async ({ page }) => {
    await page.goto('http://localhost:3000/tasks/publish')
    
    // 测试SQL注入
    await page.fill('input[name="title"]', "'; DROP TABLE tasks; --")
    await page.click('text=发布任务')
    
    await expect(page.locator('.el-message:has-text("标题不能包含特殊字符")')).toBeVisible()
    
    // 测试文件上传限制
    const [fileChooser] = await Promise.all([
      page.waitForEvent('filechooser'),
      page.click('input[type="file"]')
    ])
    
    await fileChooser.setFiles(['malicious.exe'])
    await expect(page.locator('.el-message:has-text("不支持的文件类型")')).toBeVisible()
  })
})