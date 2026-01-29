#!/bin/bash

# 数据库备份脚本
# 使用方法: ./backup.sh [类型] [保留天数]
# 类型: mysql|redis|all
# 保留天数: 默认7天

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
RETENTION_DAYS=${2:-7}

# 创建备份目录
mkdir -p "$BACKUP_DIR/mysql"
mkdir -p "$BACKUP_DIR/redis"

# 生成时间戳
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# MySQL备份函数
backup_mysql() {
    log_info "开始备份MySQL数据库..."
    
    local backup_file="$BACKUP_DIR/mysql/mysql_backup_$TIMESTAMP.sql.gz"
    
    # 执行备份
    docker exec "$MYSQL_CONTAINER" mysqldump \
        -u root \
        -p"$MYSQL_ROOT_PASSWORD" \
        --single-transaction \
        --routines \
        --triggers \
        "$MYSQL_DATABASE" | gzip > "$backup_file"
    
    if [ $? -eq 0 ]; then
        log_success "MySQL备份完成: $backup_file"
        
        # 获取备份文件大小
        local file_size=$(du -h "$backup_file" | cut -f1)
        log_info "备份文件大小: $file_size"
    else
        log_error "MySQL备份失败"
        exit 1
    fi
}

# Redis备份函数
backup_redis() {
    log_info "开始备份Redis数据..."
    
    local backup_file="$BACKUP_DIR/redis/redis_backup_$TIMESTAMP.rdb.gz"
    
    # 在Redis容器中执行BGSAVE
    docker exec "$REDIS_CONTAINER" redis-cli BGSAVE
    
    # 等待备份完成
    while [ "$(docker exec \"$REDIS_CONTAINER\" redis-cli LASTSAVE)" = "$(docker exec \"$REDIS_CONTAINER\" redis-cli LASTSAVE)" ]; do
        sleep 1
    done
    
    # 复制RDB文件并压缩
    docker cp "$REDIS_CONTAINER:/data/dump.rdb" - | gzip > "$backup_file"
    
    if [ $? -eq 0 ]; then
        log_success "Redis备份完成: $backup_file"
        
        # 获取备份文件大小
        local file_size=$(du -h "$backup_file" | cut -f1)
        log_info "备份文件大小: $file_size"
    else
        log_error "Redis备份失败"
        exit 1
    fi
}

# 清理旧备份函数
cleanup_old_backups() {
    log_info "清理超过 $RETENTION_DAYS 天的旧备份..."
    
    # 清理MySQL备份
    find "$BACKUP_DIR/mysql" -name "mysql_backup_*.sql.gz" -mtime +$RETENTION_DAYS -delete
    
    # 清理Redis备份
    find "$BACKUP_DIR/redis" -name "redis_backup_*.rdb.gz" -mtime +$RETENTION_DAYS -delete
    
    log_success "旧备份清理完成"
}

# 验证备份函数
verify_backups() {
    log_info "验证备份文件完整性..."
    
    local mysql_files=$(find "$BACKUP_DIR/mysql" -name "mysql_backup_$TIMESTAMP.sql.gz" | wc -l)
    local redis_files=$(find "$BACKUP_DIR/redis" -name "redis_backup_$TIMESTAMP.rdb.gz" | wc -l)
    
    if [ "$BACKUP_TYPE" = "mysql" ] || [ "$BACKUP_TYPE" = "all" ]; then
        if [ "$mysql_files" -eq 1 ]; then
            log_success "MySQL备份文件验证通过"
        else
            log_error "MySQL备份文件验证失败"
            exit 1
        fi
    fi
    
    if [ "$BACKUP_TYPE" = "redis" ] || [ "$BACKUP_TYPE" = "all" ]; then
        if [ "$redis_files" -eq 1 ]; then
            log_success "Redis备份文件验证通过"
        else
            log_error "Redis备份文件验证失败"
            exit 1
        fi
    fi
}

# 备份统计函数
backup_stats() {
    log_info "备份统计信息:"
    
    echo "=== 备份目录: $BACKUP_DIR ==="
    echo "MySQL备份数量: $(find "$BACKUP_DIR/mysql" -name "*.sql.gz" | wc -l)"
    echo "Redis备份数量: $(find "$BACKUP_DIR/redis" -name "*.rdb.gz" | wc -l)"
    echo "MySQL备份总大小: $(du -sh "$BACKUP_DIR/mysql" 2>/dev/null | cut -f1 || echo "0B")"
    echo "Redis备份总大小: $(du -sh "$BACKUP_DIR/redis" 2>/dev/null | cut -f1 || echo "0B")"
    echo
}

# 主函数
main() {
    BACKUP_TYPE=${1:-all}
    
    # 验证备份类型
    if [[ ! "$BACKUP_TYPE" =~ ^(mysql|redis|all)$ ]]; then
        log_error "无效的备份类型: $BACKUP_TYPE，支持的类型: mysql, redis, all"
        exit 1
    fi
    
    log_info "开始备份任务 - 类型: $BACKUP_TYPE, 保留天数: $RETENTION_DAYS"
    
    # 执行备份
    case "$BACKUP_TYPE" in
        mysql)
            backup_mysql
            ;;
        redis)
            backup_redis
            ;;
        all)
            backup_mysql
            backup_redis
            ;;
    esac
    
    # 验证备份
    verify_backups
    
    # 清理旧备份
    cleanup_old_backups
    
    # 显示统计信息
    backup_stats
    
    log_success "备份任务完成！"
}

# 读取环境变量
if [ -f "/app/.env.prod" ]; then
    source /app/.env.prod
fi

# 执行主函数
main "$@"