<template>
  <div class="payment-container">
    <div class="payment-header">
      <h1>支付订单</h1>
      <p>请选择支付方式完成付款</p>
    </div>

    <div class="payment-content">
      <el-row :gutter="24">
        <!-- 订单信息 -->
        <el-col :xs="24" :lg="16">
          <el-card class="order-card">
            <template #header>
              <h3>订单信息</h3>
            </template>

            <div class="order-info">
              <div class="order-item">
                <span class="label">订单号：</span>
                <span class="value">{{ order.order_no }}</span>
                <el-button
                  type="text"
                  size="small"
                  @click="copyOrderNo"
                >
                  复制
                </el-button>
              </div>
              
              <div class="order-item">
                <span class="label">订单类型：</span>
                <span class="value">{{ getOrderTypeText(order.order_type) }}</span>
              </div>
              
              <div v-if="order.task_info" class="task-info">
                <h4>任务详情</h4>
                <p>{{ order.task_info.title }}</p>
                <span class="task-price">¥{{ order.task_info.price }}</span>
              </div>
              
              <div class="order-time">
                <span class="label">创建时间：</span>
                <span class="value">{{ formatTime(order.created_at) }}</span>
              </div>
            </div>
          </el-card>

          <!-- 支付方式选择 -->
          <el-card class="payment-methods-card">
            <template #header>
              <h3>选择支付方式</h3>
            </template>

            <div class="payment-methods">
              <div
                v-for="method in paymentMethods"
                :key="method.code"
                class="payment-method"
                :class="{ active: selectedMethod === method.code }"
                @click="selectPaymentMethod(method.code)"
              >
                <div class="method-icon">
                  <el-icon :size="32">
                    <component :is="method.icon" />
                  </el-icon>
                </div>
                <div class="method-info">
                  <h4>{{ method.name }}</h4>
                  <p>{{ method.description }}</p>
                </div>
                <div class="method-radio">
                  <el-radio
                    v-model="selectedMethod"
                    :label="method.code"
                    size="large"
                  />
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 支付侧边栏 -->
        <el-col :xs="24" :lg="8">
          <el-card class="payment-summary">
            <template #header>
              <h3>支付信息</h3>
            </template>

            <div class="summary-content">
              <div class="amount-section">
                <div class="amount-row">
                  <span>订单金额</span>
                  <span>¥{{ order.amount }}</span>
                </div>
                
                <div v-if="order.discount_amount > 0" class="amount-row">
                  <span>优惠金额</span>
                  <span class="discount">-¥{{ order.discount_amount }}</span>
                </div>
                
                <div class="amount-row total">
                  <span>实付金额</span>
                  <span class="total-amount">¥{{ order.final_amount || order.amount }}</span>
                </div>
              </div>

              <div class="payment-actions">
                <el-button
                  type="primary"
                  size="large"
                  :loading="paymentLoading"
                  :disabled="!selectedMethod"
                  class="pay-button"
                  @click="handlePay"
                >
                  立即支付 ¥{{ order.final_amount || order.amount }}
                </el-button>
                
                <el-button
                  size="large"
                  class="cancel-button"
                  @click="handleCancel"
                >
                  取消订单
                </el-button>
              </div>

              <!-- 安全提示 -->
              <div class="security-tips">
                <h4>
                  <el-icon><Lock /></el-icon>
                  安全支付
                </h4>
                <ul>
                  <li>平台担保交易，资金安全有保障</li>
                  <li>支持主流支付方式，方便快捷</li>
                  <li>7×24小时客服支持</li>
                </ul>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 支付二维码弹窗 -->
    <el-dialog
      v-model="qrCodeVisible"
      title="扫码支付"
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
    >
      <div class="qr-payment">
        <div class="qr-header">
          <h4>请使用{{ selectedMethodName }}扫码支付</h4>
          <p>支付金额：¥{{ order.final_amount || order.amount }}</p>
        </div>
        
        <div class="qr-code">
          <el-image
            v-if="qrCodeUrl"
            :src="qrCodeUrl"
            fit="contain"
            class="qr-image"
          />
          <div v-else class="qr-loading">
            <el-icon class="is-loading"><Loading /></el-icon>
            <p>正在生成支付二维码...</p>
          </div>
        </div>
        
        <div class="qr-tips">
          <p>请打开{{ selectedMethodName }}扫一扫，扫描二维码完成支付</p>
          <p>支付完成后页面将自动跳转，请勿关闭此页面</p>
        </div>
        
        <div class="qr-actions">
          <el-button @click="refreshQRCode">刷新二维码</el-button>
          <el-button type="primary" @click="checkPaymentStatus">我已完成支付</el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 支付结果弹窗 -->
    <el-dialog
      v-model="resultVisible"
      :title="paymentSuccess ? '支付成功' : '支付失败'"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="payment-result">
        <div class="result-icon">
          <el-icon
            :size="64"
            :color="paymentSuccess ? '#67c23a' : '#f56c6c'"
          >
            <CircleCheck v-if="paymentSuccess" />
            <CircleClose v-else />
          </el-icon>
        </div>
        
        <h3>{{ paymentSuccess ? '支付成功！' : '支付失败' }}</h3>
        
        <p v-if="paymentSuccess">
          您已成功支付 ¥{{ order.final_amount || order.amount }}，
          订单将在确认后开始执行。
        </p>
        
        <p v-else>
          支付过程中出现错误，请重试或选择其他支付方式。
        </p>
        
        <div class="result-actions">
          <el-button
            v-if="paymentSuccess"
            type="primary"
            @click="goToOrderDetail"
          >
            查看订单
          </el-button>
          <el-button
            v-else
            type="primary"
            @click="retryPayment"
          >
            重新支付
          </el-button>
          <el-button @click="goToHome">返回首页</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Lock,
  Loading,
  CircleCheck,
  CircleClose,
  Wallet,
  CreditCard,
  Iphone,
  Coin
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { paymentApi } from '@/api'
import { formatTime } from '@/utils/format'
import type { Order } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 状态变量
const loading = ref(false)
const paymentLoading = ref(false)
const qrCodeVisible = ref(false)
const resultVisible = ref(false)
const paymentSuccess = ref(false)
const statusPollingTimer = ref<NodeJS.Timeout | null>(null)

// 数据
const order = reactive<Order>({
  id: 0,
  order_no: '',
  order_type: '',
  amount: 0,
  discount_amount: 0,
  final_amount: 0,
  status: '',
  payment_method: '',
  task_info: null,
  created_at: ''
})

const selectedMethod = ref('alipay')
const qrCodeUrl = ref('')
const selectedMethodName = ref('')

// 支付方式
const paymentMethods = [
  {
    code: 'alipay',
    name: '支付宝',
    description: '使用支付宝余额、花呗、银行卡等支付',
    icon: CreditCard
  },
  {
    code: 'wechat',
    name: '微信支付',
    description: '使用微信余额、银行卡等支付',
    icon: Iphone
  },
  {
    code: 'balance',
    name: '余额支付',
    description: '使用账户余额支付',
    icon: Wallet
  },
  {
    code: 'shouqianba',
    name: '收钱吧',
    description: '支持多种支付方式',
    icon: Coin
  }
]

// 加载订单信息
const loadOrderInfo = async () => {
  try {
    loading.value = true
    const orderNo = route.query.order_no as string
    
    if (!orderNo) {
      ElMessage.error('订单号不存在')
      router.push('/')
      return
    }

    const response = await paymentApi.getOrderDetail(orderNo)
    Object.assign(order, response.data)

    // 检查订单状态
    if (order.status === 'paid') {
      showPaymentResult(true)
    } else if (order.status === 'cancelled') {
      ElMessage.warning('订单已取消')
      router.push('/user/orders')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载订单信息失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 选择支付方式
const selectPaymentMethod = (method: string) => {
  selectedMethod.value = method
  const methodInfo = paymentMethods.find(m => m.code === method)
  selectedMethodName.value = methodInfo?.name || ''
}

// 复制订单号
const copyOrderNo = async () => {
  try {
    await navigator.clipboard.writeText(order.order_no)
    ElMessage.success('订单号已复制')
  } catch (error) {
    ElMessage.info('请手动复制订单号')
  }
}

// 获取订单类型文本
const getOrderTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    task: '任务支付',
    recharge: '账户充值',
    settlement: '任务结算'
  }
  return typeMap[type] || '其他'
}

// 处理支付
const handlePay = async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }

  try {
    paymentLoading.value = true

    // 余额支付直接处理
    if (selectedMethod.value === 'balance') {
      await handleBalancePayment()
      return
    }

    // 其他支付方式生成二维码
    const response = await paymentApi.createPayment({
      order_no: order.order_no,
      payment_method: selectedMethod.value,
      return_url: `${window.location.origin}/payment/success`
    })

    qrCodeUrl.value = response.data.qr_code
    qrCodeVisible.value = true

    // 开始轮询支付状态
    startStatusPolling()

  } catch (error: any) {
    ElMessage.error(error.message || '发起支付失败')
  } finally {
    paymentLoading.value = false
  }
}

// 余额支付
const handleBalancePayment = async () => {
  try {
    const response = await paymentApi.balancePayment(order.order_no)
    
    if (response.data.success) {
      showPaymentResult(true)
    } else {
      ElMessage.error(response.data.message || '余额不足')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '余额支付失败')
  }
}

// 刷新二维码
const refreshQRCode = async () => {
  try {
    const response = await paymentApi.refreshQRCode(order.order_no, selectedMethod.value)
    qrCodeUrl.value = response.data.qr_code
  } catch (error: any) {
    ElMessage.error('刷新二维码失败')
  }
}

// 检查支付状态
const checkPaymentStatus = async () => {
  try {
    const response = await paymentApi.checkPaymentStatus(order.order_no)
    
    if (response.data.status === 'paid') {
      stopStatusPolling()
      showPaymentResult(true)
    } else {
      ElMessage.info('支付尚未完成，请继续扫码支付')
    }
  } catch (error: any) {
    ElMessage.error('检查支付状态失败')
  }
}

// 开始状态轮询
const startStatusPolling = () => {
  statusPollingTimer.value = setInterval(async () => {
    try {
      const response = await paymentApi.checkPaymentStatus(order.order_no)
      
      if (response.data.status === 'paid') {
        stopStatusPolling()
        showPaymentResult(true)
      }
    } catch (error) {
      // 轮询失败时静默处理
    }
  }, 3000)
}

// 停止状态轮询
const stopStatusPolling = () => {
  if (statusPollingTimer.value) {
    clearInterval(statusPollingTimer.value)
    statusPollingTimer.value = null
  }
}

// 显示支付结果
const showPaymentResult = (success: boolean) => {
  paymentSuccess.value = success
  qrCodeVisible.value = false
  resultVisible.value = true
  
  if (success) {
    // 支付成功，清除轮询
    stopStatusPolling()
  }
}

// 重新支付
const retryPayment = () => {
  resultVisible.value = false
  handlePay()
}

// 取消订单
const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消该订单？', '取消订单', {
      confirmButtonText: '确定',
      cancelButtonText: '继续支付',
      type: 'warning'
    })

    await paymentApi.cancelOrder(order.order_no)
    
    ElMessage.success('订单已取消')
    router.push('/user/orders')
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '取消订单失败')
    }
  }
}

// 页面跳转
const goToOrderDetail = () => {
  router.push(`/user/orders/${order.id}`)
}

const goToHome = () => {
  router.push('/')
}

// 组件销毁时清理定时器
onUnmounted(() => {
  stopStatusPolling()
})

// 页面初始化
onMounted(() => {
  loadOrderInfo()
})
</script>

<style lang="scss" scoped>
.payment-container {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;

  .payment-header {
    text-align: center;
    margin-bottom: 40px;

    h1 {
      margin: 0 0 8px 0;
      font-size: 32px;
      font-weight: 600;
      color: #1a1a1a;
    }

    p {
      margin: 0;
      font-size: 16px;
      color: #666;
    }
  }

  .payment-content {
    .order-card {
      margin-bottom: 24px;

      .order-info {
        .order-item {
          display: flex;
          align-items: center;
          margin-bottom: 16px;

          .label {
            font-size: 14px;
            color: #666;
            width: 80px;
          }

          .value {
            flex: 1;
            font-size: 14px;
            color: #1a1a1a;
            font-weight: 500;
          }
        }

        .task-info {
          background: #f8f9fa;
          border-radius: 8px;
          padding: 16px;
          margin-bottom: 16px;

          h4 {
            margin: 0 0 8px 0;
            font-size: 14px;
            color: #666;
          }

          p {
            margin: 0 0 8px 0;
            font-size: 16px;
            color: #1a1a1a;
            font-weight: 500;
          }

          .task-price {
            font-size: 18px;
            font-weight: 600;
            color: #ff4d4f;
          }
        }

        .order-time {
          display: flex;
          align-items: center;

          .label {
            font-size: 14px;
            color: #666;
            width: 80px;
          }

          .value {
            font-size: 14px;
            color: #999;
          }
        }
      }
    }

    .payment-methods-card {
      margin-bottom: 24px;

      .payment-methods {
        .payment-method {
          display: flex;
          align-items: center;
          padding: 20px;
          border: 2px solid #f0f0f0;
          border-radius: 8px;
          margin-bottom: 16px;
          cursor: pointer;
          transition: all 0.3s ease;

          &:hover {
            border-color: #1890ff;
            background: #f6f8ff;
          }

          &.active {
            border-color: #1890ff;
            background: #f6f8ff;
          }

          .method-icon {
            width: 60px;
            height: 60px;
            border-radius: 8px;
            background: #f5f5f5;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 16px;

            .el-icon {
              color: #1890ff;
            }
          }

          .method-info {
            flex: 1;

            h4 {
              margin: 0 0 4px 0;
              font-size: 16px;
              font-weight: 600;
              color: #1a1a1a;
            }

            p {
              margin: 0;
              font-size: 14px;
              color: #666;
            }
          }

          .method-radio {
            margin-left: 16px;
          }
        }
      }
    }

    .payment-summary {
      .summary-content {
        .amount-section {
          margin-bottom: 32px;

          .amount-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 12px;

            span {
              font-size: 14px;
            }

            .discount {
              color: #67c23a;
            }

            &.total {
              padding-top: 16px;
              border-top: 1px solid #f0f0f0;

              span {
                font-size: 16px;
                font-weight: 600;
                color: #1a1a1a;
              }

              .total-amount {
                font-size: 24px;
                font-weight: 700;
                color: #ff4d4f;
              }
            }
          }
        }

        .payment-actions {
          display: flex;
          flex-direction: column;
          gap: 12px;
          margin-bottom: 32px;

          .pay-button {
            height: 48px;
            font-size: 18px;
            font-weight: 600;
          }

          .cancel-button {
            height: 44px;
          }
        }

        .security-tips {
          h4 {
            display: flex;
            align-items: center;
            gap: 8px;
            margin: 0 0 16px 0;
            font-size: 16px;
            font-weight: 600;
            color: #1a1a1a;

            .el-icon {
              color: #67c23a;
            }
          }

          ul {
            margin: 0;
            padding-left: 20px;

            li {
              font-size: 14px;
              color: #666;
              line-height: 1.6;
              margin-bottom: 8px;
            }
          }
        }
      }
    }
  }
}

.qr-payment {
  text-align: center;

  .qr-header {
    margin-bottom: 24px;

    h4 {
      margin: 0 0 8px 0;
      font-size: 18px;
      font-weight: 600;
      color: #1a1a1a;
    }

    p {
      margin: 0;
      font-size: 16px;
      color: #ff4d4f;
      font-weight: 600;
    }
  }

  .qr-code {
    width: 200px;
    height: 200px;
    margin: 0 auto 24px;
    border: 1px solid #f0f0f0;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;

    .qr-image {
      width: 180px;
      height: 180px;
    }

    .qr-loading {
      text-align: center;

      .el-icon {
        font-size: 32px;
        color: #1890ff;
        margin-bottom: 8px;
      }

      p {
        margin: 0;
        font-size: 14px;
        color: #666;
      }
    }
  }

  .qr-tips {
    margin-bottom: 24px;

    p {
      margin: 0 0 8px 0;
      font-size: 14px;
      color: #666;
      line-height: 1.6;
    }
  }

  .qr-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
  }
}

.payment-result {
  text-align: center;

  .result-icon {
    margin-bottom: 16px;
  }

  h3 {
    margin: 0 0 16px 0;
    font-size: 20px;
    font-weight: 600;
    color: #1a1a1a;
  }

  p {
    margin: 0 0 32px 0;
    font-size: 14px;
    color: #666;
    line-height: 1.6;
  }

  .result-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .payment-container {
    padding: 16px;

    .payment-header {
      margin-bottom: 24px;

      h1 {
        font-size: 24px;
      }
    }

    .payment-content {
      .payment-methods {
        .payment-method {
          padding: 16px;

          .method-icon {
            width: 48px;
            height: 48px;
            margin-right: 12px;
          }

          .method-info {
            h4 {
              font-size: 14px;
            }

            p {
              font-size: 12px;
            }
          }
        }
      }

      .payment-summary {
        .summary-content {
          .payment-actions {
            .pay-button {
              font-size: 16px;
            }
          }
        }
      }
    }
  }

  .qr-payment {
    .qr-header {
      h4 {
        font-size: 16px;
      }
    }

    .qr-code {
      width: 160px;
      height: 160px;

      .qr-image {
        width: 140px;
        height: 140px;
      }
    }

    .qr-actions {
      flex-direction: column;

      .el-button {
        width: 100%;
      }
    }
  }

  .payment-result {
    .result-actions {
      flex-direction: column;

      .el-button {
        width: 100%;
      }
    }
  }
}
</style>