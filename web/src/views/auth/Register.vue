<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <div class="logo">
          <img src="@/assets/logo.png" alt="Task Platform" />
        </div>
        <h2>注册新用户</h2>
        <p>加入任务交易平台，开启您的任务之旅</p>
      </div>

      <el-steps :active="currentStep" align-center class="register-steps">
        <el-step title="手机验证" />
        <el-step title="基本信息" />
        <el-step title="完成注册" />
      </el-steps>

      <!-- 步骤1：手机验证 -->
      <div v-if="currentStep === 0" class="step-content">
        <el-form
          ref="step1FormRef"
          :model="registerForm"
          :rules="step1Rules"
          class="register-form"
        >
          <el-form-item prop="phone">
            <el-input
              v-model="registerForm.phone"
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
                v-model="registerForm.code"
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
              :loading="step1Loading"
              class="next-btn"
              @click="handleStep1Next"
            >
              下一步
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤2：基本信息 -->
      <div v-if="currentStep === 1" class="step-content">
        <el-form
          ref="step2FormRef"
          :model="registerForm"
          :rules="step2Rules"
          class="register-form"
        >
          <el-form-item prop="nickname">
            <el-input
              v-model="registerForm.nickname"
              placeholder="请输入昵称"
              size="large"
              :prefix-icon="User"
              clearable
              maxlength="20"
              show-word-limit
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="请设置密码"
              size="large"
              :prefix-icon="Lock"
              show-password
              clearable
            />
          </el-form-item>

          <el-form-item prop="confirmPassword">
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="请确认密码"
              size="large"
              :prefix-icon="Lock"
              show-password
              clearable
            />
          </el-form-item>

          <el-form-item prop="agreeTerms">
            <el-checkbox v-model="registerForm.agreeTerms">
              我已阅读并同意
              <el-link type="primary" @click="showTerms">《用户协议》</el-link>
              和
              <el-link type="primary" @click="showPrivacy">《隐私政策》</el-link>
            </el-checkbox>
          </el-form-item>

          <div class="step-buttons">
            <el-button
              size="large"
              class="prev-btn"
              @click="currentStep = 0"
            >
              上一步
            </el-button>
            <el-button
              type="primary"
              size="large"
              :loading="step2Loading"
              class="next-btn"
              @click="handleStep2Next"
            >
              完成注册
            </el-button>
          </div>
        </el-form>
      </div>

      <!-- 步骤3：完成注册 -->
      <div v-if="currentStep === 2" class="step-content">
        <div class="success-content">
          <el-icon class="success-icon"><CircleCheck /></el-icon>
          <h3>注册成功！</h3>
          <p>欢迎加入任务交易平台</p>
          <div class="success-actions">
            <el-button
              type="primary"
              size="large"
              @click="goToHome"
            >
              开始使用
            </el-button>
            <el-button
              size="large"
              @click="goToProfile"
            >
              完善资料
            </el-button>
          </div>
        </div>
      </div>

      <div class="register-footer">
        <p>已有账号？
          <el-link type="primary" @click="goToLogin">立即登录</el-link>
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Phone, Key, User, Lock, CircleCheck } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { validatePhone, validateCode, validateNickname, validatePassword } from '@/utils/validators'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

// 表单数据
const registerForm = reactive({
  phone: '',
  code: '',
  nickname: '',
  password: '',
  confirmPassword: '',
  agreeTerms: false
})

// 表单引用
const step1FormRef = ref<FormInstance>()
const step2FormRef = ref<FormInstance>()

// 表单验证规则
const step1Rules: FormRules = {
  phone: [
    { validator: validatePhone, trigger: 'blur' }
  ],
  code: [
    { validator: validateCode, trigger: 'blur' }
  ]
}

const step2Rules: FormRules = {
  nickname: [
    { validator: validateNickname, trigger: 'blur' }
  ],
  password: [
    { validator: validatePassword, trigger: 'blur' }
  ],
  confirmPassword: [
    {
      validator: (rule, value, callback) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  agreeTerms: [
    {
      validator: (rule, value, callback) => {
        if (!value) {
          callback(new Error('请同意用户协议和隐私政策'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
}

// 状态变量
const currentStep = ref(0)
const step1Loading = ref(false)
const step2Loading = ref(false)
const codeLoading = ref(false)
const codeCooldown = ref(0)
const termsVisible = ref(false)
const privacyVisible = ref(false)

// 发送验证码
const sendCode = async () => {
  if (!validatePhone(null, registerForm.phone, () => {})) {
    ElMessage.error('请输入正确的手机号')
    return
  }

  codeLoading.value = true
  
  try {
    await userStore.sendVerificationCode(registerForm.phone, 'register')
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

// 步骤1：验证手机号
const handleStep1Next = async () => {
  if (!step1FormRef.value) return

  try {
    await step1FormRef.value.validate()
    
    step1Loading.value = true
    
    // 验证手机号和验证码
    await userStore.verifyPhoneCode({
      phone: registerForm.phone,
      code: registerForm.code,
      type: 'register'
    })
    
    ElMessage.success('验证成功')
    currentStep.value = 1
    
  } catch (error: any) {
    ElMessage.error(error.message || '验证失败')
  } finally {
    step1Loading.value = false
  }
}

// 步骤2：提交注册
const handleStep2Next = async () => {
  if (!step2FormRef.value) return

  try {
    await step2FormRef.value.validate()
    
    step2Loading.value = true
    
    await userStore.register({
      phone: registerForm.phone,
      nickname: registerForm.nickname,
      password: registerForm.password
    })
    
    ElMessage.success('注册成功')
    currentStep.value = 2
    
  } catch (error: any) {
    ElMessage.error(error.message || '注册失败')
  } finally {
    step2Loading.value = false
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

// 跳转到首页
const goToHome = () => {
  router.push('/dashboard')
}

// 跳转到个人资料
const goToProfile = () => {
  router.push('/user/profile')
}

// 跳转到登录页
const goToLogin = () => {
  router.push('/login')
}
</script>

<style lang="scss" scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;

  .register-card {
    width: 100%;
    max-width: 480px;
    background: white;
    border-radius: 16px;
    padding: 40px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);

    .register-header {
      text-align: center;
      margin-bottom: 32px;

      .logo {
        margin-bottom: 16px;

        img {
          width: 60px;
          height: 60px;
          border-radius: 12px;
        }
      }

      h2 {
        margin: 0 0 8px 0;
        font-size: 24px;
        font-weight: 600;
        color: #1a1a1a;
      }

      p {
        margin: 0;
        font-size: 14px;
        color: #666;
      }
    }

    .register-steps {
      margin-bottom: 32px;
    }

    .step-content {
      min-height: 300px;

      .register-form {
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

        .step-buttons {
          display: flex;
          gap: 12px;
          margin-top: 24px;

          .prev-btn {
            flex: 1;
            height: 44px;
          }

          .next-btn {
            flex: 1;
            height: 44px;
          }
        }

        .next-btn {
          width: 100%;
          height: 44px;
          font-size: 16px;
          font-weight: 500;
          border-radius: 8px;
        }
      }

      .success-content {
        text-align: center;
        padding: 40px 0;

        .success-icon {
          font-size: 64px;
          color: #67c23a;
          margin-bottom: 16px;
        }

        h3 {
          margin: 0 0 8px 0;
          font-size: 20px;
          font-weight: 600;
          color: #1a1a1a;
        }

        p {
          margin: 0 0 32px 0;
          font-size: 14px;
          color: #666;
        }

        .success-actions {
          display: flex;
          gap: 12px;
          justify-content: center;

          .el-button {
            width: 120px;
            height: 40px;
          }
        }
      }
    }

    .register-footer {
      text-align: center;
      margin-top: 24px;
      padding-top: 24px;
      border-top: 1px solid #eee;

      p {
        margin: 0;
        font-size: 14px;
        color: #666;

        .el-link {
          font-weight: 500;
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
  .register-container {
    padding: 12px;

    .register-card {
      padding: 24px;

      .register-header {
        margin-bottom: 24px;

        h2 {
          font-size: 20px;
        }
      }

      .step-content {
        min-height: 280px;

        .register-form {
          .code-input-group {
            flex-direction: column;

            .el-button {
              width: 100%;
            }
          }

          .step-buttons {
            flex-direction: column;

            .prev-btn,
            .next-btn {
              width: 100%;
            }
          }
        }

        .success-content {
          padding: 32px 0;

          .success-icon {
            font-size: 48px;
          }

          .success-actions {
            flex-direction: column;

            .el-button {
              width: 100%;
            }
          }
        }
      }
    }
  }
}
</style>