#!/bin/bash

# 任务交易平台性能优化部署脚本
# 使用方法: ./performance-deployment.sh [env: dev|prod]

set -e

ENV=${1:-dev}
PROJECT_ROOT="/www/wwwroot/task-trade-platform"
LOG_FILE="/var/log/task-platform/performance-deployment.log"

echo "==================================="
echo "🚀 任务交易平台性能优化部署"
echo "环境: $ENV"
echo "时间: $(date)"
echo "==================================="

# 创建日志目录
mkdir -p /var/log/task-platform

# 函数：记录日志
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# 函数：检查服务状态
check_service() {
    local service=$1
    local port=$2
    
    if netstat -tuln | grep -q ":$port "; then
        log "✅ $service 服务运行正常 (端口: $port)"
        return 0
    else
        log "❌ $service 服务未运行 (端口: $port)"
        return 1
    fi
}

# 函数：安装依赖
install_dependencies() {
    log "📦 安装系统依赖..."
    
    # 更新包管理器
    if command -v apt-get &> /dev/null; then
        apt-get update
        apt-get install -y nginx redis-server mysql-server
    elif command -v yum &> /dev/null; then
        yum update -y
        yum install -y nginx redis mysql-server
    fi
    
    # 安装性能监控工具
    if command -v apt-get &> /dev/null; then
        apt-get install -y htop iotop nethogs
    elif command -v yum &> /dev/null; then
        yum install -y htop iotop nethogs
    fi
}

# 函数：优化系统参数
optimize_system() {
    log "🔧 优化系统参数..."
    
    # 修改文件描述符限制
    echo "* soft nofile 65536" >> /etc/security/limits.conf
    echo "* hard nofile 65536" >> /etc/security/limits.conf
    
    # 优化内核参数
    cat >> /etc/sysctl.conf << EOF
# 网络优化
net.core.rmem_max = 134217728
net.core.wmem_max = 134217728
net.ipv4.tcp_rmem = 4096 65536 134217728
net.ipv4.tcp_wmem = 4096 65536 134217728
net.ipv4.tcp_congestion_control = bbr

# 内存优化
vm.swappiness = 10
vm.dirty_ratio = 15
vm.dirty_background_ratio = 5

# 文件系统优化
fs.file-max = 2097152
EOF
    
    sysctl -p
    
    log "✅ 系统参数优化完成"
}

# 函数：配置数据库优化
setup_database() {
    log "🗄️ 配置数据库优化..."
    
    # 复制优化脚本
    cp "$PROJECT_ROOT/scripts/database-index-optimization.sql" "/tmp/"
    
    # 应用数据库索引优化
    if command -v mysql &> /dev/null; then
        log "应用数据库索引优化..."
        mysql -u root -p task_platform < /tmp/database-index-optimization.sql
    fi
    
    # 优化MySQL配置
    cat > /etc/mysql/mysql.conf.d/mysqld.cnf << EOF
[mysqld]
# 内存优化
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
innodb_flush_log_at_trx_commit = 2
innodb_flush_method = O_DIRECT

# 连接优化
max_connections = 1000
max_connect_errors = 10000
wait_timeout = 28800
interactive_timeout = 28800

# 查询优化
query_cache_size = 64M
query_cache_type = 1
query_cache_limit = 2M

# 慢查询日志
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 1

# 二进制日志
log-bin = mysql-bin
binlog_format = ROW
expire_logs_days = 7
EOF
    
    # 重启MySQL服务
    systemctl restart mysql
    log "✅ 数据库优化配置完成"
}

# 函数：配置Redis优化
setup_redis() {
    log "⚡ 配置Redis优化..."
    
    # 备份原配置
    cp /etc/redis/redis.conf /etc/redis/redis.conf.backup
    
    # 优化Redis配置
    sed -i 's/^maxmemory .*/maxmemory 2gb/' /etc/redis/redis.conf
    sed -i 's/^maxmemory-policy .*/maxmemory-policy allkeys-lru/' /etc/redis/redis.conf
    sed -i 's/^save .*/save 900 1\nsave 300 10\nsave 60 10000/' /etc/redis/redis.conf
    
    # 性能优化
    cat >> /etc/redis/redis.conf << EOF

# 性能优化配置
tcp-keepalive 300
timeout 0
tcp-backlog 511

# 内存优化
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-size -2
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64

# 持久化优化
stop-writes-on-bgsave-error no
rdbcompression yes
rdbchecksum yes

# 日志优化
loglevel notice
syslog-enabled yes
EOF
    
    # 重启Redis服务
    systemctl restart redis
    log "✅ Redis优化配置完成"
}

# 函数：配置Nginx优化
setup_nginx() {
    log "🌐 配置Nginx优化..."
    
    # 复制性能配置
    cp "$PROJECT_ROOT/nginx-performance.conf" "/etc/nginx/sites-available/task-platform"
    
    # 启用站点
    if [ -d "/etc/nginx/sites-enabled" ]; then
        ln -sf /etc/nginx/sites-available/task-platform /etc/nginx/sites-enabled/
    fi
    
    # 备份原配置
    cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup
    
    # 优化Nginx主配置
    cat > /etc/nginx/nginx.conf << EOF
user www-data;
worker_processes auto;
worker_rlimit_nofile 65535;
pid /run/nginx.pid;

events {
    worker_connections 10240;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    # 日志格式
    log_format main '\$remote_addr - \$remote_user [\$time_local] "\$request" '
                    '\$status \$body_bytes_sent "\$http_referer" '
                    '"\$http_user_agent" "\$http_x_forwarded_for" '
                    '\$request_time \$upstream_response_time';
    
    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log warn;
    
    # 性能优化
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    keepalive_requests 1000;
    
    # 缓冲区优化
    client_body_buffer_size 128k;
    client_max_body_size 10m;
    client_header_buffer_size 3m;
    large_client_header_buffers 4 256k;
    
    # Gzip压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_comp_level 6;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
    
    # 包含站点配置
    include /etc/nginx/sites-available/task-platform;
}
EOF
    
    # 测试配置并重启
    nginx -t && systemctl restart nginx
    log "✅ Nginx优化配置完成"
}

# 函数：构建前端优化版本
build_frontend() {
    log "🎨 构建前端优化版本..."
    
    cd "$PROJECT_ROOT/web"
    
    # 安装依赖
    npm ci
    
    # 使用性能优化配置构建
    if [ "$ENV" = "prod" ]; then
        npm run build:perf
    else
        npm run build
    fi
    
    # 检查构建结果
    if [ -d "dist" ]; then
        log "✅ 前端构建完成"
        
        # 显示构建统计
        du -sh dist/
        find dist -name "*.js" -o -name "*.css" | head -10
    else
        log "❌ 前端构建失败"
        exit 1
    fi
}

# 函数：构建后端
build_backend() {
    log "🔨 构建后端服务..."
    
    cd "$PROJECT_ROOT"
    
    # 更新依赖
    go mod tidy
    
    # 构建
    if [ "$ENV" = "prod" ]; then
        CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/server cmd/server/main.go
    else
        go build -o bin/server cmd/server/main.go
    fi
    
    log "✅ 后端构建完成"
}

# 函数：配置服务管理
setup_services() {
    log "🛠️ 配置系统服务..."
    
    # 创建systemd服务文件
    cat > /etc/systemd/system/task-platform.service << EOF
[Unit]
Description=Task Platform Backend
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=$PROJECT_ROOT
ExecStart=$PROJECT_ROOT/bin/server -config=$PROJECT_ROOT/configs/config-optimized.yaml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=task-platform

# 性能优化
LimitNOFILE=65535
LimitNPROC=4096

# 环境变量
Environment=GIN_MODE=release
Environment=ENV=$ENV

[Install]
WantedBy=multi-user.target
EOF
    
    # 重新加载systemd
    systemctl daemon-reload
    
    # 启用并启动服务
    systemctl enable task-platform
    systemctl start task-platform
    
    # 检查服务状态
    sleep 3
    if systemctl is-active --quiet task-platform; then
        log "✅ 后端服务启动成功"
    else
        log "❌ 后端服务启动失败"
        journalctl -u task-platform --no-pager -n 20
    fi
}

# 函数：设置监控
setup_monitoring() {
    log "📊 设置性能监控..."
    
    # 创建监控目录
    mkdir -p /var/lib/task-platform/monitoring
    
    # 创建监控脚本
    cat > /usr/local/bin/task-platform-monitor.sh << 'EOF'
#!/bin/bash

# 任务交易平台性能监控脚本

LOG_FILE="/var/log/task-platform/monitoring.log"
METRICS_FILE="/var/lib/task-platform/monitoring/metrics.json"

# 获取系统指标
get_system_metrics() {
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    # CPU使用率
    local cpu_usage=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/%us,//')
    
    # 内存使用率
    local mem_info=$(free -m | awk 'NR==2{printf "%.2f", $3*100/$2}')
    
    # 磁盘使用率
    local disk_usage=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
    
    # 网络连接数
    local connections=$(netstat -an | grep :8080 | wc -l)
    
    # 负载平均值
    local load_avg=$(uptime | awk -F'load average:' '{print $2}' | awk '{print $1}' | sed 's/,//')
    
    # 生成JSON
    cat << METRICS
{
  "timestamp": "$timestamp",
  "cpu_usage_percent": $cpu_usage,
  "memory_usage_percent": $mem_info,
  "disk_usage_percent": $disk_usage,
  "active_connections": $connections,
  "load_average": $load_avg
}
METRICS
}

# 主监控循环
while true; do
    get_system_metrics > "$METRICS_FILE.tmp"
    mv "$METRICS_FILE.tmp" "$METRICS_FILE"
    
    echo "$(date): Metrics collected" >> "$LOG_FILE"
    sleep 30
done
EOF
    
    chmod +x /usr/local/bin/task-platform-monitor.sh
    
    # 创建监控服务
    cat > /etc/systemd/system/task-platform-monitor.service << EOF
[Unit]
Description=Task Platform Monitor
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/task-platform-monitor.sh
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl enable task-platform-monitor
    systemctl start task-platform-monitor
    
    log "✅ 性能监控设置完成"
}

# 函数：运行性能测试
run_performance_tests() {
    log "🧪 运行性能测试..."
    
    # 安装Apache Bench
    if command -v apt-get &> /dev/null; then
        apt-get install -y apache2-utils
    elif command -v yum &> /dev/null; then
        yum install -y httpd-tools
    fi
    
    # 测试API端点
    local endpoints=(
        "http://49.234.39.189:8080/health"
        "http://49.234.39.189:8080/api/"
        "http://49.234.39.189:8080/api/tasks"
    )
    
    for endpoint in "${endpoints[@]}"; do
        log "测试端点: $endpoint"
        ab -n 1000 -c 10 "$endpoint" | tee -a "$LOG_FILE"
        sleep 2
    done
    
    log "✅ 性能测试完成"
}

# 函数：生成部署报告
generate_report() {
    log "📋 生成部署报告..."
    
    local report_file="/var/log/task-platform/deployment-report-$ENV-$(date +%Y%m%d-%H%M%S).md"
    
    cat > "$report_file" << EOF
# 任务交易平台性能优化部署报告

## 部署信息
- **环境**: $ENV
- **部署时间**: $(date)
- **服务器IP**: 49.234.39.189
- **项目路径**: $PROJECT_ROOT

## 服务状态
$(systemctl is-active nginx) - Nginx Web服务器
$(systemctl is-active redis) - Redis缓存服务
$(systemctl is-active mysql) - MySQL数据库
$(systemctl is-active task-platform) - 后端API服务

## 性能配置
- **数据库**: 连接池100，索引优化，查询缓存
- **Redis**: 内存限制2GB，LRU淘汰策略
- **Nginx**: Gzip压缩，静态资源缓存，连接复用
- **前端**: 代码分割，懒加载，构建优化

## 监控指标
- **Prometheus**: 端口9090
- **健康检查**: /health
- **性能指标**: /api/performance/metrics

## 测试结果
请查看日志文件获取详细的性能测试结果。

## 优化建议
1. 定期监控内存使用情况
2. 根据实际负载调整连接池大小
3. 定期清理过期缓存
4. 监控慢查询日志
EOF
    
    log "📄 部署报告已生成: $report_file"
}

# 主执行流程
main() {
    log "🚀 开始性能优化部署..."
    
    # 检查权限
    if [ "$EUID" -ne 0 ]; then
        log "❌ 请使用root权限执行此脚本"
        exit 1
    fi
    
    # 安装依赖
    install_dependencies
    
    # 优化系统
    optimize_system
    
    # 配置数据库
    setup_database
    
    # 配置Redis
    setup_redis
    
    # 配置Nginx
    setup_nginx
    
    # 构建应用
    build_frontend
    build_backend
    
    # 配置服务
    setup_services
    
    # 设置监控
    setup_monitoring
    
    # 性能测试
    if [ "$ENV" = "prod" ]; then
        run_performance_tests
    fi
    
    # 生成报告
    generate_report
    
    log "🎉 性能优化部署完成！"
    log "📊 访问地址: http://49.234.39.189"
    log "📈 监控地址: http://49.234.39.189:9090/metrics"
    log "🔧 管理面板: http://49.234.39.189:21452/f97c6b7e"
}

# 执行主函数
main "$@"