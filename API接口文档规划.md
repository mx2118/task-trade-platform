# 任务平台API接口文档

## 1. 接口规范

### 1.1 请求格式
- **协议**: HTTPS
- **方法**: RESTful API
- **内容类型**: `application/json`
- **字符编码**: UTF-8

### 1.2 响应格式
```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": 1640995200
}
```

### 1.3 状态码说明
- `200`: 请求成功
- `201`: 创建成功
- `400`: 请求参数错误
- `401`: 未授权
- `403`: 权限不足
- `404`: 资源不存在
- `409`: 资源冲突
- `422`: 数据验证失败
- `500`: 服务器错误

### 1.4 分页格式
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

## 2. 认证接口

### 2.1 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "user123",
  "password": "password123",
  "phone": "13800138000",
  "email": "user@example.com",
  "sms_code": "123456",
  "agreement": true
}
```

**响应示例**:
```json
{
  "code": 201,
  "message": "注册成功",
  "data": {
    "user_id": 1001,
    "username": "user123",
    "phone": "13800138000",
    "email": "user@example.com"
  }
}
```

### 2.2 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "user123",
  "password": "password123"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "refresh_token_here",
    "expires_in": 86400,
    "user": {
      "id": 1001,
      "username": "user123",
      "nickname": "用户昵称",
      "avatar": "https://example.com/avatar.jpg",
      "role": "user",
      "level": 1,
      "credit_score": 5.0
    }
  }
}
```

### 2.3 刷新令牌
```http
POST /api/v1/auth/refresh
Authorization: Bearer <refresh_token>
```

### 2.4 发送短信验证码
```http
POST /api/v1/auth/send-sms
Content-Type: application/json

{
  "phone": "13800138000",
  "type": "register"
}
```

### 2.5 用户登出
```http
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

## 3. 用户管理接口

### 3.1 获取用户信息
```http
GET /api/v1/user/profile
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1001,
    "username": "user123",
    "nickname": "用户昵称",
    "avatar": "https://example.com/avatar.jpg",
    "phone": "13800138000",
    "email": "user@example.com",
    "gender": 1,
    "bio": "个人简介",
    "level": 1,
    "credit_score": 5.0,
    "profile": {
      "real_name": "张三",
      "id_card": "330102199001011234",
      "location": "浙江省杭州市",
      "work_type": "全职",
      "experience": "3年开发经验",
      "balance": 1000.00
    },
    "skills": [
      {
        "id": 1,
        "category_id": 1,
        "skill_name": "JavaScript",
        "proficiency": 5
      }
    ]
  }
}
```

### 3.2 更新用户信息
```http
PUT /api/v1/user/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "nickname": "新昵称",
  "gender": 1,
  "bio": "更新的个人简介",
  "location": "浙江省杭州市"
}
```

### 3.3 上传头像
```http
POST /api/v1/user/upload-avatar
Authorization: Bearer <token>
Content-Type: multipart/form-data

avatar: <file>
```

### 3.4 获取用户技能
```http
GET /api/v1/user/skills
Authorization: Bearer <token>
```

### 3.5 添加用户技能
```http
POST /api/v1/user/skills
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": 1,
  "skill_name": "JavaScript",
  "proficiency": 5
}
```

## 4. 任务管理接口

### 4.1 获取任务分类
```http
GET /api/v1/categories
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "技术开发",
      "icon": "https://example.com/icon1.png",
      "description": "软件开发相关任务",
      "sort_order": 1
    },
    {
      "id": 2,
      "name": "设计创意",
      "icon": "https://example.com/icon2.png",
      "description": "设计相关任务",
      "sort_order": 2
    }
  ]
}
```

### 4.2 任务列表
```http
GET /api/v1/tasks?page=1&page_size=20&category_id=1&keyword=Web开发&location=杭州&min_budget=100&max_budget=1000
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 20, 最大: 100)
- `category_id`: 分类ID
- `keyword`: 关键词搜索
- `location`: 地区筛选
- `min_budget`: 最低预算
- `max_budget`: 最高预算
- `status`: 任务状态 (0-草稿,1-待接单,2-进行中,3-已完成)
- `sort`: 排序方式 (latest-最新, budget_low-预算从低到高, budget_high-预算从高到低)

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1001,
        "title": "企业网站开发",
        "description": "需要开发一个企业展示网站...",
        "budget": 5000.00,
        "budget_type": "fixed",
        "category": {
          "id": 1,
          "name": "技术开发"
        },
        "publisher": {
          "id": 1001,
          "nickname": "企业用户",
          "avatar": "https://example.com/avatar.jpg",
          "credit_score": 4.8
        },
        "location": "浙江省杭州市",
        "worker_count": 1,
        "current_count": 0,
        "deadline": "2024-02-01T23:59:59Z",
        "status": 1,
        "created_at": "2024-01-15T10:00:00Z",
        "applications_count": 5
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 50,
      "total_pages": 3
    }
  }
}
```

### 4.3 任务详情
```http
GET /api/v1/tasks/:id
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1001,
    "title": "企业网站开发",
    "description": "详细的任务描述...",
    "budget": 5000.00,
    "budget_type": "fixed",
    "requirements": "具体要求...",
    "category": {
      "id": 1,
      "name": "技术开发"
    },
    "publisher": {
      "id": 1001,
      "nickname": "企业用户",
      "avatar": "https://example.com/avatar.jpg",
      "credit_score": 4.8,
      "completed_orders": 25
    },
    "location": "浙江省杭州市",
    "worker_count": 1,
    "current_count": 0,
    "deadline": "2024-02-01T23:59:59Z",
    "status": 1,
    "attachments": [
      {
        "id": 1,
        "filename": "需求文档.pdf",
        "url": "https://example.com/file.pdf",
        "size": 1024000
      }
    ],
    "tags": ["Web开发", "响应式", "企业网站"],
    "created_at": "2024-01-15T10:00:00Z",
    "views": 156,
    "applications": [
      {
        "id": 1,
        "applicant": {
          "id": 1002,
          "nickname": "开发者小王",
          "avatar": "https://example.com/avatar2.jpg",
          "credit_score": 4.9
        },
        "message": "我有丰富的网站开发经验...",
        "quoted_price": 4800.00,
        "attachments": [],
        "created_at": "2024-01-16T09:00:00Z"
      }
    ]
  }
}
```

### 4.4 创建任务
```http
POST /api/v1/tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": 1,
  "title": "任务标题",
  "description": "任务描述",
  "budget": 5000.00,
  "budget_type": "fixed",
  "deadline": "2024-02-01T23:59:59Z",
  "location": "浙江省杭州市",
  "worker_count": 1,
  "requirements": "具体要求",
  "attachments": [
    {
      "filename": "需求文档.pdf",
      "url": "https://example.com/file.pdf"
    }
  ],
  "tags": ["Web开发", "响应式"]
}
```

### 4.5 申请任务
```http
POST /api/v1/tasks/:id/apply
Authorization: Bearer <token>
Content-Type: application/json

{
  "message": "申请留言",
  "quoted_price": 4800.00,
  "attachments": [
    {
      "filename": "作品集.pdf",
      "url": "https://example.com/portfolio.pdf"
    }
  ]
}
```

### 4.6 获取我发布的任务
```http
GET /api/v1/tasks/my-published?status=1&page=1
Authorization: Bearer <token>
```

### 4.7 获取我申请的任务
```http
GET /api/v1/tasks/my-applied?status=0&page=1
Authorization: Bearer <token>
```

## 5. 订单管理接口

### 5.1 订单列表
```http
GET /api/v1/orders?page=1&status=1
Authorization: Bearer <token>
```

**查询参数**:
- `page`: 页码
- `status`: 订单状态 (0-待支付,1-已支付,2-进行中,3-待验收,4-已完成,5-已取消)
- `role`: 角色 (publisher-发布者, worker-工作者)

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1001,
        "order_no": "ORD202401150001",
        "task": {
          "id": 1001,
          "title": "企业网站开发"
        },
        "publisher": {
          "id": 1001,
          "nickname": "企业用户"
        },
        "worker": {
          "id": 1002,
          "nickname": "开发者小王"
        },
        "amount": 5000.00,
        "status": 2,
        "accepted_at": "2024-01-16T10:00:00Z",
        "deadline": "2024-02-01T23:59:59Z",
        "created_at": "2024-01-15T15:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 10,
      "total_pages": 1
    }
  }
}
```

### 5.2 订单详情
```http
GET /api/v1/orders/:id
Authorization: Bearer <token>
```

### 5.3 接受订单
```http
POST /api/v1/orders/:id/accept
Authorization: Bearer <token>
Content-Type: application/json

{
  "application_id": 123
}
```

### 5.4 完成订单
```http
POST /api/v1/orders/:id/complete
Authorization: Bearer <token>
Content-Type: application/json

{
  "completion_note": "工作完成说明",
  "work_evidence": [
    {
      "filename": "交付文件.zip",
      "url": "https://example.com/delivery.zip"
    }
  ]
}
```

### 5.5 确认订单完成
```http
POST /api/v1/orders/:id/confirm
Authorization: Bearer <token>
Content-Type: application/json

{
  "rating": 5,
  "comment": "工作完成得很好",
  "tags": ["专业", "准时"]
}
```

### 5.6 取消订单
```http
POST /api/v1/orders/:id/cancel
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "取消原因"
}
```

## 6. 支付接口

### 6.1 创建支付
```http
POST /api/v1/payments/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "order_id": 1001,
  "payment_method": "alipay", // alipay, wechat, balance
  "return_url": "https://example.com/payment/return",
  "notify_url": "https://example.com/payment/notify"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "payment_no": "PAY202401150001",
    "payment_url": "https://openapi.alipay.com/gateway.do?...",
    "qr_code": "data:image/png;base64,xxx"
  }
}
```

### 6.2 支付回调
```http
POST /api/v1/payments/callback
Content-Type: application/json

{
  "payment_no": "PAY202401150001",
  "third_party_no": "2024011522001234567890",
  "status": "success",
  "amount": 5000.00
}
```

### 6.3 获取钱包信息
```http
GET /api/v1/payments/wallet
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "balance": 1000.00,
    "frozen_balance": 500.00,
    "total_income": 5000.00,
    "total_withdraw": 4000.00
  }
}
```

### 6.4 提现申请
```http
POST /api/v1/payments/withdraw
Authorization: Bearer <token>
Content-Type: application/json

{
  "amount": 500.00,
  "withdraw_method": "alipay", // alipay, wechat, bank
  "account_info": {
    "account_no": "user@example.com",
    "real_name": "张三"
  }
}
```

### 6.5 获取交易记录
```http
GET /api/v1/payments/transactions?page=1&type=all
Authorization: Bearer <token>
```

## 7. 评价接口

### 7.1 创建评价
```http
POST /api/v1/reviews
Authorization: Bearer <token>
Content-Type: application/json

{
  "order_id": 1001,
  "rating": 5,
  "comment": "服务非常好，专业可靠",
  "tags": ["专业", "准时", "质量好"]
}
```

### 7.2 获取收到的评价
```http
GET /api/v1/reviews/received?page=1
Authorization: Bearer <token>
```

### 7.3 获取给出的评价
```http
GET /api/v1/reviews/given?page=1
Authorization: Bearer <token>
```

## 8. 消息接口

### 8.1 消息列表
```http
GET /api/v1/messages?page=1&user_id=1002
Authorization: Bearer <token>
```

### 8.2 发送消息
```http
POST /api/v1/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "receiver_id": 1002,
  "order_id": 1001,
  "content": "消息内容",
  "type": "text" // text, image, file
}
```

### 8.3 标记消息已读
```http
PUT /api/v1/messages/:id/read
Authorization: Bearer <token>
```

### 8.4 获取通知列表
```http
GET /api/v1/messages/notifications?page=1
Authorization: Bearer <token>
```

### 8.5 标记通知已读
```http
PUT /api/v1/messages/notifications/:id/read
Authorization: Bearer <token>
```

## 9. 文件上传接口

### 9.1 文件上传
```http
POST /api/v1/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <file>
type: avatar|task|order|message
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "url": "https://example.com/uploads/2024/01/15/filename.jpg",
    "filename": "filename.jpg",
    "size": 1024000,
    "type": "image/jpeg"
  }
}
```

## 10. 管理员接口

### 10.1 用户管理
```http
GET /api/v1/admin/users?page=1&status=1
Authorization: Bearer <admin_token>
```

### 10.2 更新用户状态
```http
PUT /api/v1/admin/users/:id/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "status": 0 // 0-禁用,1-正常,2-待审核
}
```

### 10.3 获取待审核任务
```http
GET /api/v1/admin/tasks/pending?page=1
Authorization: Bearer <admin_token>
```

### 10.4 审核任务
```http
PUT /api/v1/admin/tasks/:id/approve
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "approved": true,
  "reason": "审核通过"
}
```

### 10.5 投诉管理
```http
GET /api/v1/admin/complaints?page=1&status=0
Authorization: Bearer <admin_token>
```

### 10.6 处理投诉
```http
PUT /api/v1/admin/complaints/:id/handle
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "status": 2, // 1-处理中,2-已解决,3-已驳回
  "admin_note": "处理结果说明"
}
```

## 11. 统计数据接口

### 11.1 仪表板数据
```http
GET /api/v1/admin/dashboard
Authorization: Bearer <admin_token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "users": {
      "total": 10000,
      "today_new": 50,
      "active": 2500
    },
    "tasks": {
      "total": 5000,
      "today_new": 20,
      "in_progress": 150,
      "completed": 4500
    },
    "orders": {
      "total": 4000,
      "today_new": 15,
      "completed": 3800,
      "total_amount": 2000000.00
    }
  }
}
```

## 12. WebSocket接口

### 12.1 连接WebSocket
```javascript
const ws = new WebSocket('wss://api.example.com/ws?token=<token>');

ws.onopen = function(event) {
    console.log('WebSocket连接已建立');
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    // 处理消息
};
```

### 12.2 消息格式
```json
{
  "type": "notification", // notification, message, order_update
  "data": {
    "title": "新任务通知",
    "content": "您有新的任务可以申请",
    "timestamp": 1640995200
  }
}
```

这份API文档提供了完整的接口规范，涵盖了任务平台的所有核心功能，便于前后端开发协作和第三方系统集成。