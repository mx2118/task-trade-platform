// 全局类型扩展
declare global {
  interface Window {
    __HIDE_LOADING__?: () => void
    __LOW_BANDWIDTH__?: boolean
    $checkPermission?: (permission: string) => boolean
    $hasPermission?: (permission: string) => boolean
  }

  // requestIdleCallback 类型声明
  interface IdleDeadline {
    readonly didTimeout: boolean
    timeRemaining: () => number
  }

  interface IdleRequestOptions {
    timeout?: number
  }

  function requestIdleCallback(
    callback: (deadline: IdleDeadline) => void,
    options?: IdleRequestOptions
  ): number
  
  function cancelIdleCallback(handle: number): void
}

export {}
