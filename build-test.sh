#!/bin/bash

# 快速构建测试脚本

echo "🚀 开始构建测试..."
echo ""

cd /www/wwwroot/task-trade-platform/web

# 清理旧构建
echo "📦 清理旧构建..."
rm -rf dist/
echo "✅ 清理完成"
echo ""

# 执行构建
echo "🔨 开始生产构建..."
npm run build

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 构建成功！"
    echo ""
    
    # 显示构建产物大小
    echo "📊 构建产物分析："
    echo "----------------------------------------"
    
    if [ -d "dist" ]; then
        echo "总大小："
        du -sh dist/
        echo ""
        
        echo "JS 文件："
        du -sh dist/js/ 2>/dev/null || echo "  未找到 js 目录"
        
        echo "CSS 文件："
        du -sh dist/css/ 2>/dev/null || echo "  未找到 css 目录"
        
        echo "图片文件："
        du -sh dist/images/ 2>/dev/null || echo "  未找到 images 目录"
        
        echo ""
        echo "文件列表："
        ls -lh dist/ | head -n 20
    fi
    
    echo "----------------------------------------"
    echo ""
    echo "🎉 构建测试完成！"
    echo ""
    echo "💡 下一步："
    echo "   1. 运行 'npm run preview' 预览构建结果"
    echo "   2. 使用 Chrome DevTools 测试性能"
    echo "   3. 使用 Lighthouse 进行评分"
else
    echo ""
    echo "❌ 构建失败，请检查错误信息"
    exit 1
fi
