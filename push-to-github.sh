#!/bin/bash

echo "=========================================="
echo "任务交易平台 - GitHub 推送脚本"
echo "=========================================="
echo ""

# 检查是否在项目根目录
if [ ! -f "go.mod" ]; then
    echo "❌ 错误：请在项目根目录运行此脚本"
    exit 1
fi

echo "✅ 项目目录验证通过"
echo ""

# 检查git状态
echo "📋 检查Git状态..."
git status --porcelain
echo ""

# 检查是否有未提交的更改
if [ -n "$(git status --porcelain)" ]; then
    echo "⚠️  发现未提交的更改，正在提交..."
    
    # 添加所有文件
    git add .
    
    # 提交更改
    git commit -m "$(cat <<'EOF'
feat: 清理敏感文件并优化项目结构

- 删除包含敏感信息的脚本文件
- 清理英文技术文档
- 添加.gitignore排除敏感配置
- 保留中文文档和README.md
- 优化项目安全性

安全改进:
- 移除明文密码和SSH密钥
- 清理服务器IP地址信息
- 规范化.gitignore配置
EOF
)"
    
    echo "✅ 更改已提交"
else
    echo "✅ 没有未提交的更改"
fi

echo ""
echo "🚀 准备推送到GitHub..."
echo ""

# 提供推送说明
echo "请选择推送方式："
echo "1. 使用Personal Access Token (推荐)"
echo "2. 使用SSH密钥"
echo "3. 手动推送"
echo ""

read -p "请输入选择 (1-3): " choice

case $choice in
    1)
        echo ""
        echo "📝 使用Personal Access Token推送"
        echo "请按以下步骤操作："
        echo "1. 访问 https://github.com/settings/tokens"
        echo "2. 点击 'Generate new token (classic)'"
        echo "3. 选择 'repo' 权限"
        echo "4. 复制生成的token"
        echo ""
        read -s -p "请输入您的GitHub Personal Access Token: " token
        echo ""
        
        # 使用token推送
        git remote set-url origin https://$token@github.com/mx2118/task-trade-platform.git
        git push -u origin main
        ;;
        
    2)
        echo ""
        echo "🔑 使用SSH密钥推送"
        echo "请按以下步骤操作："
        echo "1. 生成SSH密钥：ssh-keygen -t ed25519 -C 'your_email@example.com'"
        echo "2. 添加到ssh-agent：eval \$(ssh-agent -s)"
        echo "3. 添加私钥：ssh-add ~/.ssh/id_ed25519"
        echo "4. 复制公钥到GitHub：cat ~/.ssh/id_ed25519.pub"
        echo "5. 在GitHub设置中添加SSH密钥"
        echo ""
        read -p "完成后按回车继续..."
        
        # 切换到SSH URL
        git remote set-url origin git@github.com:mx2118/task-trade-platform.git
        git push -u origin main
        ;;
        
    3)
        echo ""
        echo "📖 手动推送说明："
        echo ""
        echo "方法1 - 使用token："
        echo "  git remote set-url origin https://YOUR_TOKEN@github.com/mx2118/task-trade-platform.git"
        echo "  git push -u origin main"
        echo ""
        echo "方法2 - 使用SSH："
        echo "  git remote set-url origin git@github.com:mx2118/task-trade-platform.git"
        echo "  git push -u origin main"
        echo ""
        echo "方法3 - 使用git credential helper："
        echo "  git config --global credential.helper store"
        echo "  git push -u origin main"
        echo ""
        echo "仓库地址：https://github.com/mx2118/task-trade-platform.git"
        ;;
        
    *)
        echo "❌ 无效选择"
        exit 1
        ;;
esac

echo ""
echo "✅ 推送完成！"
echo "🌐 项目地址：https://github.com/mx2118/task-trade-platform"