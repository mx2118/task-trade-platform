<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="logo">
          <img src="@/assets/logo.png" alt="Task Platform" />
        </div>
        <h2>任务交易平台</h2>
        <p>安全便捷的任务交易服务</p>
      </div>

      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="phone">
          <el-input
            v-model="loginForm.phone"
            placeholder="请输入手机号"
            size="large"
            :prefix-icon="Phone"
            clearable
            maxlength="11"
            show-word-limit
          />
        </el-form-item>

        <el-form-item prop="code">
          <div class="code-input-group">
            <el-input
              v-model="loginForm.code"
              placeholder="请输入验证码"
              size="large"
              :prefix-icon="Key"
              maxlength="6"
              show-word-limit
            />
            <el-button
              :disabled="codeCooldown > 0"
              :loading="codeLoading"
              size="large"
              @click="sendCode"
            >
              {{ codeCooldown > 0 ? `${codeCooldown}s后重试` : '获取验证码' }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loginLoading"
            class="login-btn"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>

      <div class="divider">
        <span>或</span>
      </div>

      <div class="social-login">
        <el-button
          size="large"
          class="wechat-btn"
          @click="handleWechatLogin"
        >
          <el-icon><WechatFilled /></el-icon>
          微信登录
        </el-button>
        
        <el-button
          size="large"
          class="alipay-btn"
          @click="handleAlipayLogin"
        >
          <el-icon><CreditCard /></el-icon>
          支付宝登录
        </el-button>
      </div>

      <div class="login-footer">
        <p>登录即表示同意
          <el-link type="primary" @click="showTerms">《用户协议》</el-link>
          和
          <el-link type="primary" @click="showPrivacy">《隐私政策》</el-link>
        </p>
      </div>
    </div>

    <!-- 协议弹窗 -->
    <el-dialog
      v-model="termsVisible"
      title="用户协议"
      width="60%"
      destroy-on-close
    >
      <div class="terms-content">
        <h3>1. 服务条款</h3>
        <p>欢迎使用任务交易平台。本平台为用户提供任务发布、接取、支付等服务的在线交易平台。</p>
        
        <h3>2. 用户责任</h3>
        <p>用户应确保发布的信息真实有效，遵守相关法律法规，不得发布违法、违规内容。</p>
        
        <h3>3. 平台责任</h3>
        <p>平台将为用户提供安全、稳定的交易环境，保护用户信息安全，处理纠纷调解。</p>
        
        <h3>4. 交易规则</h3>
        <p>用户应按照平台规则进行交易，确保任务质量和支付安全。</p>
      </div>
      <template #footer>
        <el-button @click="termsVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 隐私政策弹窗 -->
    <el-dialog
      v-model="privacyVisible"
      title="隐私政策"
      width="60%"
      destroy-on-close
    >
      <div class="privacy-content">
        <h3>1. 信息收集</h3>
        <p>我们收集您的手机号、微信/支付宝信息用于账户认证和交易服务。</p>
        
        <h3>2. 信息使用</h3>
        <p>您的个人信息仅用于提供服务、改善用户体验和确保交易安全。</p>
        
        <h3>3. 信息保护</h3>
        <p>我们采用行业标准的加密技术保护您的个人信息安全。</p>
        
        <h3>4. 信息共享</h3>
        <p>未经您的同意，我们不会向第三方共享您的个人信息。</p>
      </div>
      <template #footer>
        <el-button @click="privacyVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElLoading } from 'element-plus'
import { Phone, Key, WechatFilled, CreditCard } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { validatePhone, validateCode } from '@/utils/validators'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 表单数据
const loginForm = reactive({
  phone: '',
  code: ''
})

// 表单引用
const loginFormRef = ref<FormInstance>()

// 表单验证规则
const loginRules: FormRules = {
  phone: [
    { validator: validatePhone, trigger: 'blur' }
  ],
  code: [
    { validator: validateCode, trigger: 'blur' }
  ]
}

// 状态变量
const loginLoading = ref(false)
const codeLoading = ref(false)
const codeCooldown = ref(0)
const termsVisible = ref(false)
const privacyVisible = ref(false)

// 发送验证码
const sendCode = async () => {
  if (!validatePhone(null, loginForm.phone, () => {})) {
    ElMessage.error('请输入正确的手机号')
    return
  }

  codeLoading.value = true
  
  try {
    await userStore.sendVerificationCode(loginForm.phone)
    ElMessage.success('验证码已发送')
    
    // 开始倒计时
    codeCooldown.value = 60
    const timer = setInterval(() => {
      codeCooldown.value--
      if (codeCooldown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error: any) {
    ElMessage.error(error.message || '发送验证码失败')
  } finally {
    codeLoading.value = false
  }
}

// 手机号登录
const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    
    loginLoading.value = true
    
    await userStore.phoneLogin({
      phone: loginForm.phone,
      code: loginForm.code
    })

    ElMessage.success('登录成功')
    
    // 跳转到目标页面
    const redirect = route.query.redirect as string || '/dashboard'
    router.replace(redirect)
    
  } catch (error: any) {
    if (error.message) {
      ElMessage.error(error.message)
    }
  } finally {
    loginLoading.value = false
  }
}

// 微信登录
const handleWechatLogin = async () => {
  try {
    loginLoading.value = true
    
    // 获取微信授权URL
    const authUrl = await userStore.getWechatAuthUrl()
    
    // 跳转到微信授权页面
    window.location.href = authUrl
    
  } catch (error: any) {
    ElMessage.error(error.message || '微信登录失败')
    loginLoading.value = false
  }
}

// 支付宝登录
const handleAlipayLogin = async () => {
  try {
    loginLoading.value = true
    
    // 获取支付宝授权URL
    const authUrl = await userStore.getAlipayAuthUrl()
    
    // 跳转到支付宝授权页面
    window.location.href = authUrl
    
  } catch (error: any) {
    ElMessage.error(error.message || '支付宝登录失败')
    loginLoading.value = false
  }
}

// 显示用户协议
const showTerms = () => {
  termsVisible.value = true
}

// 显示隐私政策
const showPrivacy = () => {
  privacyVisible.value = true
}

// 处理第三方登录回调
onMounted(async () => {
  const { code, state, platform } = route.query
  
  if (code && platform) {
    loginLoading.value = true
    
    try {
      if (platform === 'wechat') {
        await userStore.wechatLogin({
          code: code as string,
          state: state as string
        })
      } else if (platform === 'alipay') {
        await userStore.alipayLogin({
          auth_code: code as string,
          state: state as string
        })
      }
      
      ElMessage.success('登录成功')
      
      const redirect = sessionStorage.getItem('login_redirect') || '/dashboard'
      sessionStorage.removeItem('login_redirect')
      router.replace(redirect)
      
    } catch (error: any) {
      ElMessage.error(error.message || '登录失败')
      router.replace('/login')
    } finally {
      loginLoading.value = false
    }
  }
})
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;

  .login-card {
    width: 100%;
    max-width: 420px;
    background: white;
    border-radius: 16px;
    padding: 40px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);

    .login-header {
      text-align: center;
      margin-bottom: 40px;

      .logo {
        margin-bottom: 20px;

        img {
          width: 80px;
          height: 80px;
          border-radius: 12px;
        }
      }

      h2 {
        margin: 0 0 8px 0;
        font-size: 28px;
        font-weight: 600;
        color: #1a1a1a;
      }

      p {
        margin: 0;
        font-size: 14px;
        color: #666;
      }
    }

    .login-form {
      .code-input-group {
        display: flex;
        gap: 12px;

        .el-input {
          flex: 1;
        }

        .el-button {
          width: 120px;
          flex-shrink: 0;
        }
      }

      .login-btn {
        width: 100%;
        height: 48px;
        font-size: 16px;
        font-weight: 500;
        border-radius: 8px;
      }
    }

    .divider {
      text-align: center;
      margin: 32px 0;
      position: relative;

      &::before {
        content: '';
        position: absolute;
        top: 50%;
        left: 0;
        right: 0;
        height: 1px;
        background: #e0e0e0;
      }

      span {
        background: white;
        padding: 0 16px;
        color: #999;
        font-size: 14px;
        position: relative;
        z-index: 1;
      }
    }

    .social-login {
      display: flex;
      flex-direction: column;
      gap: 12px;

      .wechat-btn {
        background: #09bb07;
        color: white;
        border: none;
        height: 44px;
        font-size: 15px;

        &:hover {
          background: #08a006;
        }

        .el-icon {
          margin-right: 8px;
        }
      }

      .alipay-btn {
        background: #1677ff;
        color: white;
        border: none;
        height: 44px;
        font-size: 15px;

        &:hover {
          background: #1466e6;
        }

        .el-icon {
          margin-right: 8px;
        }
      }
    }

    .login-footer {
      margin-top: 24px;
      text-align: center;

      p {
        margin: 0;
        font-size: 12px;
        color: #999;
        line-height: 1.6;

        .el-link {
          font-size: 12px;
        }
      }
    }
  }
}

.terms-content,
.privacy-content {
  h3 {
    color: #333;
    margin-top: 20px;
    margin-bottom: 8px;

    &:first-child {
      margin-top: 0;
    }
  }

  p {
    color: #666;
    line-height: 1.6;
    margin-bottom: 12px;
  }
}

// 响应式设计
@media (max-width: 480px) {
  .login-container {
    padding: 12px;

    .login-card {
      padding: 24px;

      .login-header {
        margin-bottom: 24px;

        h2 {
          font-size: 24px;
        }
      }

      .code-input-group {
        flex-direction: column;

        .el-button {
          width: 100%;
        }
      }

      .social-login {
        .wechat-btn,
        .alipay-btn {
          height: 48px;
        }
      }
    }
  }
}
</style>