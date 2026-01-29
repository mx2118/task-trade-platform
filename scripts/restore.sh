#!/bin/bash

# 数据库恢复脚本
# 使用方法: ./restore.sh [类型] [备份文件]
# 类型: mysql|redis
# 备份文件: 备份文件路径（不提供则列出可用备份）

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 配置
BACKUP_DIR="/app/backups"
MYSQL_CONTAINER="task-trade-mysql-prod"
REDIS_CONTAINER="task-trade-redis-prod"
MYSQL_DATABASE="task_trade"

# 列出可用备份
list_backups() {
    local type=$1
    local backup_path="$BACKUP_DIR/$type"
    
    log_info "可用的 $type 备份文件:"
    
    if [ ! -d "$backup_path" ]; then
        log_warning "备份目录不存在: $backup_path"
        return 1
    fi
    
    case "$type" in
        mysql)
            ls -la "$backup_path"/mysql_backup_*.sql.gz 2>/dev/null | awk '{print $9, $5}' | while read -r file size; do
                if [ -n "$file" ]; then
                    filename=$(basename "$file")
                    echo "  $filename (大小: $size)"
                fi
            done
            ;;
        redis)
            ls -la "$backup_path"/redis_backup_*.rdb.gz 2>/dev/null | awk '{print $9, $5}' | while read -r file size; do
                if [ -n "$file" ]; then
                    filename=$(basename "$file")
                    echo "  $filename (大小: $size)"
                fi
            done
            ;;
    esac
}

# MySQL恢复函数
restore_mysql() {
    local backup_file=$1
    
    if [ -z "$backup_file" ]; then
        list_backups "mysql"
        log_info "请使用: ./restore.sh mysql [备份文件名]"
        exit 1
    fi
    
    local full_path="$BACKUP_DIR/mysql/$backup_file"
    
    if [ ! -f "$full_path" ]; then
        log_error "备份文件不存在: $full_path"
        exit 1
    fi
    
    log_warning "警告：此操作将覆盖现有的MySQL数据库"
    read -p "确认继续？(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "操作已取消"
        exit 0
    fi
    
    log_info "开始恢复MySQL数据库..."
    
    # 停止应用服务以避免数据冲突
    log_info "停止应用服务..."
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml stop backend
    
    # 恢复数据库
    gunzip -c "$full_path" | docker exec -i "$MYSQL_CONTAINER" mysql \
        -u root \
        -p"$MYSQL_ROOT_PASSWORD" \
        "$MYSQL_DATABASE"
    
    if [ $? -eq 0 ]; then
        log_success "MySQL数据库恢复完成"
        
        # 重启应用服务
        log_info "重启应用服务..."
        docker-compose -f docker-compose.yml -f docker-compose.prod.yml start backend
        
        log_success "服务重启完成"
    else
        log_error "MySQL数据库恢复失败"
        exit 1
    fi
}

# Redis恢复函数
restore_redis() {
    local backup_file=$1
    
    if [ -z "$backup_file" ]; then
        list_backups "redis"
        log_info "请使用: ./restore.sh redis [备份文件名]"
        exit 1
    fi
    
    local full_path="$BACKUP_DIR/redis/$backup_file"
    
    if [ ! -f "$full_path" ]; then
        log_error "备份文件不存在: $full_path"
        exit 1
    fi
    
    log_warning "警告：此操作将覆盖现有的Redis数据"
    read -p "确认继续？(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "操作已取消"
        exit 0
    fi
    
    log_info "开始恢复Redis数据..."
    
    # 停止Redis服务
    log_info "停止Redis服务..."
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml stop redis
    
    # 清除现有数据
    docker exec "$REDIS_CONTAINER" sh -c "rm -f /data/dump.rdb"
    
    # 解压并复制备份文件
    gunzip -c "$full_path" | docker cp - "$REDIS_CONTAINER:/data/dump.rdb"
    
    # 重启Redis服务
    log_info "重启Redis服务..."
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml start redis
    
    # 等待Redis启动
    sleep 5
    
    # 验证恢复
    if docker exec "$REDIS_CONTAINER" redis-cli ping > /dev/null 2>&1; then
        log_success "Redis数据恢复完成"
    else
        log_error "Redis数据恢复失败"
        exit 1
    fi
}

# 验证恢复函数
verify_restore() {
    local type=$1
    
    log_info "验证 $type 数据恢复..."
    
    case "$type" in
        mysql)
            if docker exec "$MYSQL_CONTAINER" mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "USE $MYSQL_DATABASE; SELECT COUNT(*) FROM users;" > /dev/null 2>&1; then
                log_success "MySQL数据验证通过"
            else
                log_error "MySQL数据验证失败"
                exit 1
            fi
            ;;
        redis)
            if docker exec "$REDIS_CONTAINER" redis-cli ping > /dev/null 2>&1; then
                local key_count=$(docker exec "$REDIS_CONTAINER" redis-cli DBSIZE)
                log_success "Redis数据验证通过，共有 $key_count 个键"
            else
                log_error "Redis数据验证失败"
                exit 1
            fi
            ;;
    esac
}

# 主函数
main() {
    if [ $# -lt 1 ]; then
        echo "使用方法: ./restore.sh [类型] [备份文件]"
        echo "类型: mysql, redis"
        echo
        echo "示例:"
        echo "  ./restore.sh mysql mysql_backup_20231201_120000.sql.gz"
        echo "  ./restore.sh redis redis_backup_20231201_120000.rdb.gz"
        exit 1
    fi
    
    RESTORE_TYPE=$1
    BACKUP_FILE=$2
    
    # 验证恢复类型
    if [[ ! "$RESTORE_TYPE" =~ ^(mysql|redis)$ ]]; then
        log_error "无效的恢复类型: $RESTORE_TYPE，支持的类型: mysql, redis"
        exit 1
    fi
    
    log_info "开始数据恢复任务 - 类型: $RESTORE_TYPE"
    
    # 执行恢复
    case "$RESTORE_TYPE" in
        mysql)
            restore_mysql "$BACKUP_FILE"
            ;;
        redis)
            restore_redis "$BACKUP_FILE"
            ;;
    esac
    
    # 验证恢复
    verify_restore "$RESTORE_TYPE"
    
    log_success "数据恢复任务完成！"
}

# 读取环境变量
if [ -f "/app/.env.prod" ]; then
    source /app/.env.prod
fi

# 执行主函数
main "$@"