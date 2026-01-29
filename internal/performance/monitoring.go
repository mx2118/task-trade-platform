package performance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	
	"task-platform-api/internal/models"
)

// PerformanceMonitor 性能监控器
type PerformanceMonitor struct {
	logger *zap.Logger
	
	// Prometheus 指标
	requestsTotal       *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
	requestSize         *prometheus.HistogramVec
	responseSize        *prometheus.HistogramVec
	activeConnections   prometheus.Gauge
	memoryUsage         prometheus.Gauge
	goroutines          prometheus.Gauge
	dbConnections       prometheus.Gauge
	
	// 自定义指标
	customMetrics       sync.Map
	
	// 配置
	config              MonitoringConfig
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled             bool          `yaml:"enabled"`
	Port               string        `yaml:"port"`
	MetricsPath        string        `yaml:"metrics_path"`
	HealthCheckPath    string        `yaml:"health_check_path"`
	CollectInterval    time.Duration `yaml:"collect_interval"`
	SlowQueryThreshold time.Duration `yaml:"slow_query_threshold"`
}

// CustomMetric 自定义指标
type CustomMetric struct {
	Name        string
	Value       float64
	Timestamp   time.Time
	Labels      map[string]string
	Description string
}

// NewPerformanceMonitor 创建性能监控器
func NewPerformanceMonitor(logger *zap.Logger, config MonitoringConfig) *PerformanceMonitor {
	pm := &PerformanceMonitor{
		logger: logger,
		config: config,
	}
	
	if config.Enabled {
		pm.initPrometheusMetrics()
	}
	
	return pm
}

// initPrometheusMetrics 初始化Prometheus指标
func (pm *PerformanceMonitor) initPrometheusMetrics() {
	pm.requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)
	
	pm.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	
	pm.requestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "endpoint"},
	)
	
	pm.responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "endpoint"},
	)
	
	pm.activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)
	
	pm.memoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
	)
	
	pm.goroutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines",
			Help: "Number of goroutines",
		},
	)
	
	pm.dbConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections",
			Help: "Number of database connections",
		},
	)
	
	// 注册指标
	prometheus.MustRegister(
		pm.requestsTotal,
		pm.requestDuration,
		pm.requestSize,
		pm.responseSize,
		pm.activeConnections,
		pm.memoryUsage,
		pm.goroutines,
		pm.dbConnections,
	)
}

// Middleware 创建性能监控中间件
func (pm *PerformanceMonitor) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !pm.config.Enabled {
			c.Next()
			return
		}
		
		start := time.Now()
		pm.activeConnections.Inc()
		
		// 获取请求大小
		var requestSize float64
		if c.Request.ContentLength > 0 {
			requestSize = float64(c.Request.ContentLength)
		}
		
		// 处理请求
		c.Next()
		
		// 记录指标
		duration := time.Since(start).Seconds()
		statusCode := fmt.Sprintf("%d", c.Writer.Status())
		method := c.Request.Method
		endpoint := c.FullPath()
		
		if endpoint == "" {
			endpoint = "unknown"
		}
		
		pm.requestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
		pm.requestDuration.WithLabelValues(method, endpoint).Observe(duration)
		pm.requestSize.WithLabelValues(method, endpoint).Observe(requestSize)
		
		// 获取响应大小
		if c.Writer.Size() > 0 {
			pm.responseSize.WithLabelValues(method, endpoint).Observe(float64(c.Writer.Size()))
		}
		
		pm.activeConnections.Dec()
		
		// 记录慢请求
		if time.Duration(duration*1000) > pm.config.SlowQueryThreshold {
			pm.logger.Warn("Slow request detected",
				zap.String("method", method),
				zap.String("endpoint", endpoint),
				zap.Float64("duration", duration),
				zap.Int("status", c.Writer.Status()),
			)
		}
	}
}

// CollectSystemMetrics 收集系统指标
func (pm *PerformanceMonitor) CollectSystemMetrics(ctx context.Context) {
	if !pm.config.Enabled {
		return
	}
	
	ticker := time.NewTicker(pm.config.CollectInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			pm.updateSystemMetrics()
		case <-ctx.Done():
			return
		}
	}
}

// updateSystemMetrics 更新系统指标
func (pm *PerformanceMonitor) updateSystemMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	pm.memoryUsage.Set(float64(m.Alloc))
	pm.goroutines.Set(float64(runtime.NumGoroutine()))
}

// RecordCustomMetric 记录自定义指标
func (pm *PerformanceMonitor) RecordCustomMetric(name, description string, value float64, labels map[string]string) {
	if !pm.config.Enabled {
		return
	}
	
	metric := CustomMetric{
		Name:        name,
		Value:       value,
		Timestamp:   time.Now(),
		Labels:      labels,
		Description: description,
	}
	
	pm.customMetrics.Store(name, metric)
}

// GetCustomMetrics 获取自定义指标
func (pm *PerformanceMonitor) GetCustomMetrics() map[string]CustomMetric {
	metrics := make(map[string]CustomMetric)
	
	pm.customMetrics.Range(func(key, value interface{}) bool {
		metrics[key.(string)] = value.(CustomMetric)
		return true
	})
	
	return metrics
}

// SetDBConnections 设置数据库连接数
func (pm *PerformanceMonitor) SetDBConnections(count int) {
	if pm.config.Enabled && pm.dbConnections != nil {
		pm.dbConnections.Set(float64(count))
	}
}

// SetupRoutes 设置监控路由
func (pm *PerformanceMonitor) SetupRoutes(router *gin.Engine) {
	if !pm.config.Enabled {
		return
	}
	
	// Prometheus 指标端点
	router.GET(pm.config.MetricsPath, gin.WrapH(promhttp.Handler()))
	
	// 健康检查端点
	router.GET(pm.config.HealthCheckPath, pm.healthCheck)
	
	// 性能指标 API
	router.GET("/api/performance/metrics", pm.getPerformanceMetrics)
	router.GET("/api/performance/stats", pm.getPerformanceStats)
}

// healthCheck 健康检查
func (pm *PerformanceMonitor) healthCheck(c *gin.Context) {
	status := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
		"uptime":    time.Since(time.Now()).String(),
	}
	
	c.JSON(http.StatusOK, status)
}

// getPerformanceMetrics 获取性能指标
func (pm *PerformanceMonitor) getPerformanceMetrics(c *gin.Context) {
	metrics := map[string]interface{}{
		"prometheus": map[string]interface{}{
			"requests_total":      pm.requestsTotal,
			"request_duration":    pm.requestDuration,
			"memory_usage":        pm.memoryUsage,
			"goroutines":         pm.goroutines,
			"active_connections":  pm.activeConnections,
		},
		"custom": pm.GetCustomMetrics(),
		"system": pm.getSystemInfo(),
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Performance metrics retrieved successfully",
		"data":    metrics,
	})
}

// getPerformanceStats 获取性能统计
func (pm *PerformanceMonitor) getPerformanceStats(c *gin.Context) {
	stats := map[string]interface{}{
		"uptime":           time.Since(time.Now()).String(),
		"memory":           pm.getMemoryStats(),
		"goroutines":       runtime.NumGoroutine(),
		"gc_cycles":        pm.getGCCycles(),
		"custom_metrics":   len(pm.GetCustomMetrics()),
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Performance stats retrieved successfully",
		"data":    stats,
	})
}

// getSystemInfo 获取系统信息
func (pm *PerformanceMonitor) getSystemInfo() map[string]interface{} {
	return map[string]interface{}{
		"go_version":    runtime.Version(),
		"go_os":         runtime.GOOS,
		"go_arch":       runtime.GOARCH,
		"cpu_count":     runtime.NumCPU(),
		"process_id":    0, // 可以通过os.Getpid()获取
	}
}

// getMemoryStats 获取内存统计
func (pm *PerformanceMonitor) getMemoryStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return map[string]interface{}{
		"alloc":         m.Alloc,
		"total_alloc":   m.TotalAlloc,
		"sys":           m.Sys,
		"lookups":       m.Lookups,
		"mallocs":       m.Mallocs,
		"frees":         m.Frees,
		"heap_alloc":    m.HeapAlloc,
		"heap_sys":      m.HeapSys,
		"heap_idle":     m.HeapIdle,
		"heap_inuse":    m.HeapInuse,
		"heap_released": m.HeapReleased,
		"heap_objects":  m.HeapObjects,
		"stack_inuse":   m.StackInuse,
		"stack_sys":     m.StackSys,
	}
}

// getGCCycles 获取GC周期数
func (pm *PerformanceMonitor) getGCCycles() uint32 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.NumGC
}

// AlertManager 告警管理器
type AlertManager struct {
	monitor *PerformanceMonitor
	logger  *zap.Logger
	rules   []AlertRule
}

// AlertRule 告警规则
type AlertRule struct {
	Name        string
	Metric      string
	Threshold   float64
	Operator    string // "gt", "lt", "eq"
	Severity    string
	Description string
	Enabled     bool
}

// NewAlertManager 创建告警管理器
func NewAlertManager(monitor *PerformanceMonitor, logger *zap.Logger) *AlertManager {
	return &AlertManager{
		monitor: monitor,
		logger:  logger,
		rules: []AlertRule{
			{
				Name:        "High Memory Usage",
				Metric:      "memory_usage",
				Threshold:   0.8, // 80%
				Operator:    "gt",
				Severity:    "warning",
				Description: "Memory usage is above 80%",
				Enabled:     true,
			},
			{
				Name:        "Too Many Goroutines",
				Metric:      "goroutines",
				Threshold:   1000,
				Operator:    "gt",
				Severity:    "critical",
				Description: "Number of goroutines exceeds 1000",
				Enabled:     true,
			},
		},
	}
}

// CheckAlerts 检查告警
func (am *AlertManager) CheckAlerts() {
	for _, rule := range am.rules {
		if !rule.Enabled {
			continue
		}
		
		if am.evaluateRule(rule) {
			am.sendAlert(rule)
		}
	}
}

// evaluateRule 评估告警规则
func (am *AlertManager) evaluateRule(rule AlertRule) bool {
	// 这里应该根据具体的指标值进行评估
	// 简化实现
	return false
}

// sendAlert 发送告警
func (am *AlertManager) sendAlert(rule AlertRule) {
	am.logger.Error("Alert triggered",
		zap.String("rule", rule.Name),
		zap.String("severity", rule.Severity),
		zap.String("description", rule.Description),
	)
	
	// 这里可以集成其他告警渠道，如邮件、短信、钉钉等
}