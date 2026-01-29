import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { visualizer } from 'rollup-plugin-visualizer'

// 性能优化配置
export default defineConfig({
  plugins: [
    vue(),
    
    // 打包体积分析
    visualizer({
      filename: 'dist/stats.html',
      open: false,
      gzipSize: true,
      brotliSize: true
    })
  ],
  
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@components': resolve(__dirname, 'src/components'),
      '@views': resolve(__dirname, 'src/views'),
      '@utils': resolve(__dirname, 'src/utils'),
      '@api': resolve(__dirname, 'src/api'),
      '@stores': resolve(__dirname, 'src/stores'),
      '@assets': resolve(__dirname, 'src/assets')
    }
  },
  
  build: {
    // 代码分割优化
    rollupOptions: {
      output: {
        // 手动代码分割
        manualChunks: {
          // Vue相关
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          // UI组件库
          'ui-vendor': ['element-plus', '@element-plus/icons-vue'],
          // 工具库
          'utils-vendor': ['axios', 'dayjs', 'lodash-es', 'crypto-js'],
          // 图表库
          'chart-vendor': ['echarts', 'vue-echarts'],
          // 其他第三方库
          'other-vendor': ['qrcode', 'js-cookie', 'nprogress']
        },
        
        // 优化文件命名
        chunkFileNames: 'assets/js/[name]-[hash].js',
        entryFileNames: 'assets/js/[name]-[hash].js',
        assetFileNames: 'assets/[ext]/[name]-[hash].[ext]'
      }
    },
    
    // 压缩配置
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // 移除console
        drop_debugger: true, // 移除debugger
        pure_funcs: ['console.log', 'console.info', 'console.debug']
      },
      mangle: {
        safari10: true
      }
    },
    
    // 目标浏览器
    target: 'es2015',
    
    // 启用CSS代码分割
    cssCodeSplit: true,
    
    // 构建报告
    reportCompressedSize: true,
    
    // 分包大小警告限制
    chunkSizeWarningLimit: 1000
  },
  
  // 开发服务器配置
  server: {
    port: 3000,
    host: '0.0.0.0',
    cors: true,
    proxy: {
      '/api': {
        target: 'http://49.234.39.189:8080',
        changeOrigin: true,
        rewrite: (path) => path
      }
    }
  },
  
  // 预览服务器配置
  preview: {
    port: 3001,
    host: '0.0.0.0'
  },
  
  // CSS配置
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "@/assets/styles/variables.scss";`
      }
    },
    
    // PostCSS配置
    postcss: {
      plugins: [
        require('autoprefixer'),
        require('cssnano')({
          preset: 'default'
        })
      ]
    }
  },
  
  // 依赖优化
  optimizeDeps: {
    include: [
      'vue',
      'vue-router',
      'pinia',
      'element-plus',
      '@element-plus/icons-vue',
      'axios',
      'dayjs',
      'lodash-es',
      'echarts',
      'vue-echarts'
    ],
    
    exclude: [
      // 排除不需要预构建的包
    ]
  },
  
  // 环境变量
  define: {
    __VUE_PROD_DEVTOOLS__: false,  // 生产环境关闭Vue DevTools
    __VUE_OPTIONS_API__: true,
    __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false
  },
  
  // esbuild配置
  esbuild: {
    // 生产环境移除console和debugger
    drop: process.env.NODE_ENV === 'production' ? ['console', 'debugger'] : []
  }
})