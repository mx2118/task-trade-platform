package performance

import (
	"context"
	"fmt"
	"sync"
	"time"
	
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

// ConcurrencyManager 并发管理器
type ConcurrencyManager struct {
	limiter      *rate.Limiter
	redisClient  *redis.Client
	locks        sync.Map
	metrics      *ConcurrencyMetrics
}

// ConcurrencyMetrics 并发指标
type ConcurrencyMetrics struct {
	TotalRequests     int64
	RejectedRequests int64
	ActiveGoroutines int64
	Mutex            sync.Mutex
}

// DistributedLock 分布式锁接口
type DistributedLock interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error
	TryLock(ctx context.Context, key string) (bool, error)
}

// RedisDistributedLock Redis分布式锁实现
type RedisDistributedLock struct {
	client *redis.Client
}

// NewConcurrencyManager 创建并发管理器
func NewConcurrencyManager(rps float64, burst int, redisClient *redis.Client) *ConcurrencyManager {
	return &ConcurrencyManager{
		limiter:     rate.NewLimiter(rate.Limit(rps), burst),
		redisClient: redisClient,
		metrics:     &ConcurrencyMetrics{},
	}
}

// AllowRequest 检查请求是否允许通过
func (c *ConcurrencyManager) AllowRequest() bool {
	c.metrics.Mutex.Lock()
	c.metrics.TotalRequests++
	c.metrics.Mutex.Unlock()

	allowed := c.limiter.Allow()
	
	if !allowed {
		c.metrics.Mutex.Lock()
		c.metrics.RejectedRequests++
		c.metrics.Mutex.Unlock()
	}
	
	return allowed
}

// GetMetrics 获取并发指标
func (c *ConcurrencyManager) GetMetrics() ConcurrencyMetrics {
	c.metrics.Mutex.Lock()
	defer c.metrics.Mutex.Unlock()
	
	return ConcurrencyMetrics{
		TotalRequests:     c.metrics.TotalRequests,
		RejectedRequests: c.metrics.RejectedRequests,
		ActiveGoroutines: c.metrics.ActiveGoroutines,
	}
}

// SafeExecute 安全执行函数，包含panic恢复
func (c *ConcurrencyManager) SafeExecute(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	
	c.metrics.Mutex.Lock()
	c.metrics.ActiveGoroutines++
	c.metrics.Mutex.Unlock()
	
	defer func() {
		c.metrics.Mutex.Lock()
		c.metrics.ActiveGoroutines--
		c.metrics.Mutex.Unlock()
	}()
	
	return fn()
}

// WorkerPool 工作池
type WorkerPool struct {
	workerCount int
	taskQueue   chan func()
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan func(), queueSize),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start 启动工作池
func (p *WorkerPool) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// Stop 停止工作池
func (p *WorkerPool) Stop() {
	p.cancel()
	close(p.taskQueue)
	p.wg.Wait()
}

// Submit 提交任务
func (p *WorkerPool) Submit(task func()) bool {
	select {
	case p.taskQueue <- task:
		return true
	default:
		return false // 队列已满
	}
}

// worker 工作协程
func (p *WorkerPool) worker() {
	defer p.wg.Done()
	
	for {
		select {
		case task := <-p.taskQueue:
			if task != nil {
				task()
			}
		case <-p.ctx.Done():
			return
		}
	}
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	maxFailures   int
	resetTimeout  time.Duration
	failures      int
	lastFailTime  time.Time
	state         CircuitState
	mutex         sync.RWMutex
}

type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Call 执行函数调用
func (cb *CircuitBreaker) Call(fn func() error) error {
	if !cb.allowRequest() {
		return fmt.Errorf("circuit breaker is open")
	}
	
	err := fn()
	cb.recordResult(err)
	return err
}

// allowRequest 检查是否允许请求
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// recordResult 记录执行结果
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()
		
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		}
	} else {
		cb.failures = 0
		cb.state = StateClosed
	}
}

// RedisDistributedLock 实现
func NewRedisDistributedLock(client *redis.Client) *RedisDistributedLock {
	return &RedisDistributedLock{client: client}
}

// Lock 获取锁
func (r *RedisDistributedLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	result := r.client.SetNX(ctx, key, "locked", ttl)
	return result.Val(), result.Err()
}

// Unlock 释放锁
func (r *RedisDistributedLock) Unlock(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// TryLock 尝试获取锁
func (r *RedisDistributedLock) TryLock(ctx context.Context, key string) (bool, error) {
	result := r.client.SetNX(ctx, key, "locked", time.Second*10)
	return result.Val(), result.Err()
}

// BatchProcessor 批处理器
type BatchProcessor struct {
	batchSize    int
	flushTimeout time.Duration
	buffer       []interface{}
	mutex        sync.Mutex
	flushFn      func([]interface{}) error
	ticker       *time.Ticker
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewBatchProcessor 创建批处理器
func NewBatchProcessor(batchSize int, flushTimeout time.Duration, flushFn func([]interface{}) error) *BatchProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	
	bp := &BatchProcessor{
		batchSize:    batchSize,
		flushTimeout: flushTimeout,
		buffer:       make([]interface{}, 0, batchSize),
		flushFn:      flushFn,
		ticker:       time.NewTicker(flushTimeout),
		ctx:          ctx,
		cancel:       cancel,
	}
	
	// 启动定时刷新协程
	go bp.flushLoop()
	
	return bp
}

// Add 添加数据到批处理
func (bp *BatchProcessor) Add(data interface{}) error {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	
	bp.buffer = append(bp.buffer, data)
	
	if len(bp.buffer) >= bp.batchSize {
		return bp.flush()
	}
	
	return nil
}

// flush 立即刷新缓冲区
func (bp *BatchProcessor) flush() error {
	if len(bp.buffer) == 0 {
		return nil
	}
	
	// 复制缓冲区数据
	data := make([]interface{}, len(bp.buffer))
	copy(data, bp.buffer)
	
	// 清空缓冲区
	bp.buffer = bp.buffer[:0]
	
	// 执行刷新函数
	return bp.flushFn(data)
}

// flushLoop 定时刷新循环
func (bp *BatchProcessor) flushLoop() {
	for {
		select {
		case <-bp.ticker.C:
			bp.flush()
		case <-bp.ctx.Done():
			bp.ticker.Stop()
			bp.flush() // 最后一次刷新
			return
		}
	}
}

// Stop 停止批处理器
func (bp *BatchProcessor) Stop() {
	bp.cancel()
}

// CacheWarmer 缓存预热器
type CacheWarmer struct {
	redisClient *redis.Client
	workerPool  *WorkerPool
}

// NewCacheWarmer 创建缓存预热器
func NewCacheWarmer(redisClient *redis.Client) *CacheWarmer {
	return &CacheWarmer{
		redisClient: redisClient,
		workerPool:  NewWorkerPool(5, 100),
	}
}

// WarmupCache 预热缓存
func (cw *CacheWarmer) WarmupCache(ctx context.Context, warmupTasks []func() error) error {
	cw.workerPool.Start()
	defer cw.workerPool.Stop()
	
	var wg sync.WaitGroup
	errors := make(chan error, len(warmupTasks))
	
	for _, task := range warmupTasks {
		wg.Add(1)
		cw.workerPool.Submit(func() {
			defer wg.Done()
			if err := task(); err != nil {
				errors <- err
			}
		})
	}
	
	wg.Wait()
	close(errors)
	
	// 检查是否有错误
	for err := range errors {
		if err != nil {
			return fmt.Errorf("缓存预热失败: %w", err)
		}
	}
	
	return nil
}

// RateLimiter 限流器
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mutex    sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	rl.mutex.RLock()
	limiter, exists := rl.limiters[key]
	rl.mutex.RUnlock()
	
	if !exists {
		rl.mutex.Lock()
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
		rl.mutex.Unlock()
	}
	
	return limiter.Allow()
}

// CleanupExpiredLimiters 清理过期的限流器
func (rl *RateLimiter) CleanupExpiredLimiters() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	// 这里可以根据需要实现清理逻辑
	// 例如：定期清理长时间未使用的限流器
}