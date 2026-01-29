-- Redis缓存策略优化配置
-- 用于redis-cli EVAL执行或集成到Go代码中

-- 1. 内存优化配置
-- 建议在redis.conf中设置：
-- maxmemory 2gb
-- maxmemory-policy allkeys-lru
-- save 900 1
-- save 300 10
-- save 60 10000

-- 2. 键命名规范
local function generateKey(prefix, id)
    return string.format("task_platform:%s:%s", prefix, id)
end

-- 3. 用户缓存策略
-- 用户基础信息缓存（过期时间：30分钟）
local function cacheUser(userId, userData)
    local key = generateKey("user", userId)
    redis.call('SET', key, userData)
    redis.call('EXPIRE', key, 1800) -- 30分钟
end

-- 用户会话缓存（过期时间：7天）
local function cacheSession(sessionId, sessionData)
    local key = generateKey("session", sessionId)
    redis.call('SET', key, sessionData)
    redis.call('EXPIRE', key, 604800) -- 7天
end

-- 用户信用分数缓存（过期时间：1小时）
local function cacheUserCredit(userId, creditData)
    local key = generateKey("credit", userId)
    redis.call('SET', key, creditData)
    redis.call('EXPIRE', key, 3600) -- 1小时
end

-- 4. 任务缓存策略
-- 任务详情缓存（过期时间：15分钟）
local function cacheTask(taskId, taskData)
    local key = generateKey("task", taskId)
    redis.call('SET', key, taskData)
    redis.call('EXPIRE', key, 900) -- 15分钟
end

-- 任务列表缓存（过期时间：5分钟）
local function cacheTaskList(listType, page, data)
    local key = generateKey("task_list", string.format("%s:%s", listType, page))
    redis.call('SET', key, data)
    redis.call('EXPIRE', key, 300) -- 5分钟
end

-- 热门任务缓存（过期时间：10分钟）
local function cacheHotTasks(data)
    local key = generateKey("hot_tasks", "daily")
    redis.call('SET', key, data)
    redis.call('EXPIRE', key, 600) -- 10分钟
end

-- 5. 缓存预热脚本
local function warmupCache()
    -- 预热热门任务
    local hotTasksKey = generateKey("hot_tasks", "daily")
    if not redis.call('EXISTS', hotTasksKey) then
        -- 从数据库加载热门任务数据
        -- 这里应该调用数据库查询
    end
    
    -- 预热用户统计
    local statsKey = generateKey("stats", "daily")
    if not redis.call('EXISTS', statsKey) then
        local stats = {
            total_users = 1000,
            total_tasks = 500,
            active_users = 200,
            completed_tasks = 300
        }
        redis.call('HSET', statsKey, unpack(stats))
        redis.call('EXPIRE', statsKey, 86400) -- 24小时
    end
end

-- 6. 缓存失效策略
-- 用户相关缓存失效
local function invalidateUserCache(userId)
    local pattern = generateKey("user", userId)
    redis.call('DEL', pattern)
    
    local creditKey = generateKey("credit", userId)
    redis.call('DEL', creditKey)
    
    -- 失效用户相关的任务列表缓存
    local userTasksPattern = string.format("task_platform:task_list:user_%s:*", userId)
    local keys = redis.call('KEYS', userTasksPattern)
    if #keys > 0 then
        redis.call('DEL', unpack(keys))
    end
end

-- 任务相关缓存失效
local function invalidateTaskCache(taskId)
    local taskKey = generateKey("task", taskId)
    redis.call('DEL', taskKey)
    
    -- 失效所有任务列表缓存
    local listPattern = "task_platform:task_list:*"
    local keys = redis.call('KEYS', listPattern)
    if #keys > 0 then
        redis.call('DEL', unpack(keys))
    end
    
    -- 失效热门任务缓存
    local hotTasksKey = generateKey("hot_tasks", "daily")
    redis.call('DEL', hotTasksKey)
end

-- 7. 缓存统计和监控
local function getCacheStats()
    local info = redis.call('INFO', 'memory')
    local keyspace = redis.call('INFO', 'keyspace')
    
    return {
        memory_info = info,
        keyspace_info = keyspace,
        total_commands = redis.call('INFO', 'stats')
    }
end

-- 8. 批量操作优化
-- 批量获取用户信息
local function batchGetUsers(userIds)
    local keys = {}
    for i, userId in ipairs(userIds) do
        table.insert(keys, generateKey("user", userId))
    end
    
    return redis.call('MGET', unpack(keys))
end

-- 批量设置用户信息
local function batchSetUsers(users)
    local args = {}
    for _, user in ipairs(users) do
        local key = generateKey("user", user.id)
        table.insert(args, key)
        table.insert(args, user.data)
        table.insert(args, 'EX')
        table.insert(args, 1800) -- 30分钟过期
    end
    
    return redis.call('MSET', unpack(args))
end

-- 9. 分布式锁实现
local function acquireLock(lockKey, expireTime)
    lockKey = generateKey("lock", lockKey)
    local result = redis.call('SET', lockKey, "1", 'NX', 'EX', expireTime)
    return result == "OK"
end

local function releaseLock(lockKey)
    lockKey = generateKey("lock", lockKey)
    return redis.call('DEL', lockKey)
end

-- 10. 限流器实现
local function rateLimit(userId, limit, window)
    local key = generateKey("rate_limit", userId)
    local current = redis.call('INCR', key)
    
    if current == 1 then
        redis.call('EXPIRE', key, window)
    end
    
    return current <= limit
end

-- 11. 消息队列优化
-- 发布任务状态变更通知
local function publishTaskUpdate(taskId, status)
    local channel = "task_platform:task_updates"
    local message = {
        task_id = taskId,
        status = status,
        timestamp = redis.call('TIME')[1]
    }
    
    return redis.call('PUBLISH', channel, cjson.encode(message))
end

-- 12. 缓存预热模板
local CACHE_TEMPLATES = {
    user_info = {
        key_pattern = "user:{id}",
        ttl = 1800,
        warmup_query = "SELECT * FROM users WHERE status = 1 ORDER BY create_time DESC LIMIT 100"
    },
    
    task_list = {
        key_pattern = "task_list:{type}:{page}",
        ttl = 300,
        warmup_query = "SELECT * FROM tasks WHERE status = 1 ORDER BY create_time DESC LIMIT 50"
    },
    
    hot_tasks = {
        key_pattern = "hot_tasks:daily",
        ttl = 600,
        warmup_query = "SELECT * FROM tasks WHERE status = 1 ORDER BY view_count DESC, create_time DESC LIMIT 20"
    }
}

-- 返回主要函数供外部调用
return {
    cacheUser = cacheUser,
    cacheSession = cacheSession,
    cacheTask = cacheTask,
    cacheTaskList = cacheTaskList,
    invalidateUserCache = invalidateUserCache,
    invalidateTaskCache = invalidateTaskCache,
    batchGetUsers = batchGetUsers,
    batchSetUsers = batchSetUsers,
    acquireLock = acquireLock,
    releaseLock = releaseLock,
    rateLimit = rateLimit,
    publishTaskUpdate = publishTaskUpdate,
    warmupCache = warmupCache,
    getCacheStats = getCacheStats,
    templates = CACHE_TEMPLATES
}