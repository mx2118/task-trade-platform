<template>
  <div class="task-card">
    <div class="task-header">
      <div class="task-title">
        <h3>{{ task.title }}</h3>
        <div class="task-badges">
          <el-tag v-if="task.is_urgent" type="danger" size="small">
            加急
          </el-tag>
          <el-tag v-if="task.is_remote" type="success" size="small">
            远程
          </el-tag>
          <el-tag
            :type="getTaskStatusType(task.status)"
            size="small"
          >
            {{ getTaskStatusText(task.status) }}
          </el-tag>
        </div>
      </div>
      <div class="task-price">
        <span class="price">¥{{ task.amount || task.price }}</span>
        <span class="price-label">赏金</span>
      </div>
    </div>

    <div class="task-description">
      {{ task.description || task.content }}
    </div>

    <div class="task-meta">
      <div class="meta-left">
        <span class="deadline">
          <el-icon><Clock /></el-icon>
          截止：{{ formatTime(task.deadline) }}
        </span>
        <span class="location" v-if="task.location">
          <el-icon><Location /></el-icon>
          {{ task.location }}
        </span>
      </div>
      
      <div class="meta-right">
        <el-avatar
          :size="32"
          :src="task.publisher_avatar"
        >
          {{ task.publisher_name?.charAt(0) || 'U' }}
        </el-avatar>
      </div>
    </div>

    <div class="task-footer">
      <div class="stats">
        <span>
          <el-icon><View /></el-icon>
          {{ task.view_count || 0 }}
        </span>
        <span>
          <el-icon><User /></el-icon>
          {{ task.applyCount || 0 }}
        </span>
        <span>{{ formatTime(task.createdAt) }}</span>
      </div>
      
      <div class="actions">
        <el-button
          v-if="task.status === 1"
          type="primary"
          size="small"
          @click.stop="handleApply(task)"
        >
          申请接取
        </el-button>
        <el-button
          size="small"
          @click.stop="handleShare(task)"
        >
          分享
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus'
import { Clock, Location, View, User } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/format'
import type { Task } from '@/types'

interface Props {
  task: Task
}

const props = defineProps<Props>()
const emit = defineEmits<{
  click: [task: Task]
}>()

// 获取任务状态类型
const getTaskStatusType = (status: number | string): 'success' | 'warning' | 'info' | 'primary' | 'danger' => {
  const statusMap: Record<number | string, 'success' | 'warning' | 'info' | 'primary' | 'danger'> = {
    0: 'info',    // 草稿
    1: 'primary', // 待接取
    2: 'warning', // 进行中
    3: 'info',    // 待验收
    4: 'success', // 已完成
    5: 'danger'   // 已取消
  }
  return statusMap[status] || 'info'
}

// 获取任务状态文本
const getTaskStatusText = (status: number | string) => {
  const statusMap: Record<number | string, string> = {
    0: '草稿',
    1: '待接取',
    2: '进行中',
    3: '待验收',
    4: '已完成',
    5: '已取消'
  }
  return statusMap[status] || '未知'
}

// 申请接取任务
const handleApply = (task: Task) => {
  ElMessageBox.confirm(
    '确定要申请接取这个任务吗？',
    '申请任务',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  ).then(() => {
    ElMessage.success('申请已提交，请等待任务发布者确认')
  }).catch(() => {
    // 用户取消
  })
}

// 分享任务
const handleShare = async (task: Task) => {
  const shareUrl = `${window.location.origin}/tasks/${task.id}`
  
  try {
    await navigator.clipboard.writeText(shareUrl)
    ElMessage.success('链接已复制到剪贴板')
  } catch (error) {
    ElMessageBox.alert(
      `任务链接：${shareUrl}`,
      '分享任务',
      {
        confirmButtonText: '确定',
        type: 'info'
      }
    )
  }
}

// 处理点击
const handleClick = () => {
  emit('click', props.task)
}
</script>

<style lang="scss" scoped>
.task-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  
  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    transform: translateY(-2px);
  }

  .task-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 12px;

    .task-title {
      flex: 1;

      h3 {
        margin: 0 0 8px 0;
        font-size: 16px;
        font-weight: 600;
        color: #1a1a1a;
        line-height: 1.4;
      }

      .task-badges {
        display: flex;
        gap: 6px;
        flex-wrap: wrap;
      }
    }

    .task-price {
      text-align: right;
      margin-left: 16px;

      .price {
        display: block;
        font-size: 20px;
        font-weight: 700;
        color: #ff4d4f;
        line-height: 1;
      }

      .price-label {
        font-size: 12px;
        color: #999;
        margin-top: 4px;
      }
    }
  }

  .task-description {
    color: #666;
    font-size: 14px;
    line-height: 1.5;
    margin-bottom: 16px;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .task-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    .meta-left {
      display: flex;
      align-items: center;
      gap: 12px;
      flex-wrap: wrap;

      .deadline,
      .location {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 12px;
        color: #666;

        .el-icon {
          font-size: 14px;
        }
      }
    }

    .meta-right {
      display: flex;
      align-items: center;
    }
  }

  .task-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 12px;
    border-top: 1px solid #f0f0f0;

    .stats {
      display: flex;
      gap: 12px;
      font-size: 12px;
      color: #999;

      span {
        display: flex;
        align-items: center;
        gap: 4px;

        .el-icon {
          font-size: 14px;
        }
      }
    }

    .actions {
      display: flex;
      gap: 8px;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .task-card {
    padding: 16px;

    .task-header {
      flex-direction: column;

      .task-price {
        margin-left: 0;
        margin-top: 12px;
        text-align: left;

        .price {
          font-size: 18px;
        }
      }
    }

    .task-meta {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;
    }

    .task-footer {
      flex-direction: column;
      gap: 12px;

      .actions {
        width: 100%;
        justify-content: flex-end;
      }
    }
  }
}
</style>