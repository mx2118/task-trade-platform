import request from '../request'
import type { ApiResponse } from '@/types'

// 上传文件
export const uploadFile = (file: File, options?: {
  type?: 'avatar' | 'task' | 'document' | 'image'
  progress?: (progress: number) => void
}): Promise<ApiResponse<{
  url: string
  filename: string
  size: number
  type: string
}>> => {
  const formData = new FormData()
  formData.append('file', file)
  if (options?.type) {
    formData.append('type', options.type)
  }

  return request.post('/api/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: (progressEvent) => {
      if (options?.progress && progressEvent.total) {
        const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        options.progress(progress)
      }
    }
  })
}

// 上传头像
export const uploadAvatar = (file: File): Promise<ApiResponse<{
  url: string
  filename: string
}>> => {
  const formData = new FormData()
  formData.append('avatar', file)
  
  return request.post('/api/upload/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 上传任务附件
export const uploadTaskFile = (file: File, taskId: number): Promise<ApiResponse<{
  url: string
  filename: string
  size: number
  type: string
}>> => {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('taskId', taskId.toString())
  
  return request.post('/api/upload/task', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取上传限制
export const getUploadLimits = (): Promise<ApiResponse<{
  maxFileSize: number
  allowedTypes: string[]
  maxFiles: number
}>> => {
  return request.get('/api/upload/limits')
}

// 删除文件
export const deleteFile = (filename: string): Promise<ApiResponse<null>> => {
  return request.delete(`/api/upload/${filename}`)
}

// 导出 API 对象
export const uploadApi = {
  uploadFile,
  uploadAvatar,
  uploadTaskFile,
  getUploadLimits,
  deleteFile
}