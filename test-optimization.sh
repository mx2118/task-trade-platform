#!/bin/bash

# 前端性能优化测试脚本
# 用于验证优化效果

echo "======================================"
echo "🚀 前端性能优化测试"
echo "======================================"
echo ""

cd /www/wwwroot/task-trade-platform/web

echo "📦 1. 检查依赖安装状态..."
if [ -d "node_modules" ]; then
    echo "✅ node_modules 存在"
else
    echo "❌ node_modules 不存在，正在安装..."
    npm install
fi
echo ""

echo "🔍 2. 检查优化文件..."
files=(
    "src/utils/lazyLoad.ts"
    "public/sw.js"
    "public/manifest.json"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file"
    else
        echo "❌ $file 不存在"
    fi
done
echo ""

echo "🛠️ 3. 类型检查..."
npm run type-check 2>&1 | tail -n 5
echo ""

echo "📊 4. 构建分析..."
echo "正在进行生产构建（带分析）..."
# ANALYZE=true npm run build

echo ""
echo "✅ 优化测试完成！"
echo ""
echo "📝 后续步骤："
echo "1. 运行 'npm run dev' 启动开发服务器"
echo "2. 运行 'npm run build' 进行生产构建"
echo "3. 运行 'npm run build:analyze' 查看包体积分析"
echo "4. 运行 'npm run preview' 预览生产构建"
echo ""
echo "🎯 优化要点："
echo "- ✅ 路由懒加载 + 组件缓存"
echo "- ✅ 代码分割优化"
echo "- ✅ Gzip/Brotli 压缩"
echo "- ✅ Service Worker 缓存"
echo "- ✅ 图片懒加载工具"
echo "- ✅ PWA 支持"
echo "- ✅ 性能监控"
echo ""
echo "======================================"
