/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module '*.svg' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent
  export default component
}

declare module '*.png' {
  const src: string
  export default src
}

declare module '*.jpg' {
  const src: string
  export default src
}

declare module '*.jpeg' {
  const src: string
  export default src
}

declare module '*.gif' {
  const src: string
  export default src
}

declare module '*.webp' {
  const src: string
  export default src
}

declare module '*.ico' {
  const src: string
  export default src
}

declare module '*.bmp' {
  const src: string
  export default src
}

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_API_BASE_URL: string
  readonly VITE_UPLOAD_URL: string
  readonly VITE_WS_URL: string
  readonly VITE_USE_MOCK: string
  readonly VITE_SHOW_DEV_TOOLS: string
  readonly VITE_ENABLE_ERROR_MONITOR: string
  readonly VITE_ENABLE_PERFORMANCE_MONITOR: string
  readonly VITE_SHOUQIANBA_APP_ID: string
  readonly VITE_SHOUQIANBA_MERCHANT_NO: string
  readonly VITE_SHOUQIANBA_SANDBOX: string
  readonly VITE_WECHAT_APP_ID: string
  readonly VITE_ALIPAY_APP_ID: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}