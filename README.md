# 任务交易平台 API

基于Go+Vue3.5技术栈的任务交易平台后端服务，支持微信/支付宝授权登录，集成收钱吧支付，实现任务发布、接取、验收、结算全流程。

## 🚀 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- Docker & Docker Compose (可选)

### 本地开发

1. **克隆项目**
```bash
git clone https://github.com/mx2118/task-trade-platform.git
cd task-trade-platform
```

2. **安装依赖**
```bash
make deps
```

3. **配置数据库**
```bash
# 启动数据库服务 (需要Docker)
docker-compose up -d mysql redis

# 或手动启动MySQL和Redis
```

4. **复制配置文件**
```bash
cp configs/config.yaml.example configs/config.yaml
# 修改配置文件中的数据库连接等信息
```

5. **数据库迁移**
```bash
make migrate
```

6. **启动服务**
```bash
make run
```

服务将在 `http://localhost:8080` 启动

### Docker 部署

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f api

# 停止服务
docker-compose down
```

## 📁 项目结构

```
task-trade-platform/
├── cmd/server/              # 应用入口
├── internal/               # 私有代码
│   ├── api/v1/            # API层
│   │   ├── handlers/       # HTTP处理器
│   │   ├── middleware/     # 中间件
│   │   └── routes/         # 路由定义
│   ├── config/             # 配置管理
│   ├── models/             # 数据模型
│   └── pkg/               # 内部包
├── configs/                # 配置文件
├── scripts/                # 脚本文件
├── deployments/            # 部署文件
├── web/                   # Vue3.5前端项目
└── tests/                  # 测试文件
```

## 🔧 开发工具

```bash
# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test

# 生成测试覆盖率报告
make coverage

# 生成API文档
make swagger

# 生成Mock文件
make mock
```

## 📊 核心功能

### 用户模块
- [x] 微信/支付宝授权登录
- [x] JWT令牌认证
- [x] 用户信息管理
- [x] 会话管理

### 任务模块
- [x] 任务发布和审核
- [x] 任务搜索和筛选
- [x] 任务接取和交付
- [x] 任务验收流程

### 支付模块
- [x] 收钱吧支付集成
- [x] 资金预缴和结算
- [x] 退款处理
- [x] 钱包管理

### 风控模块
- [x] 设备指纹识别
- [x] 异常行为监测
- [x] 信誉体系
- [x] 申诉处理

## 🔐 安全特性

- JWT无状态认证
- API接口签名验证
- 敏感数据加密存储
- 请求频率限制
- SQL注入防护
- XSS攻击防护

## 📈 监控和日志

- 结构化日志 (Zap)
- Prometheus指标监控
- Grafana可视化
- 链路追踪支持

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📝 API文档

启动服务后，访问以下地址查看API文档：

- Swagger UI: `http://localhost:8080/swagger/index.html`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3001` (admin/admin)

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🆘 支持

如果您遇到问题或有疑问，请：

1. 查看 [FAQ](docs/FAQ.md)
2. 搜索 [Issues](https://github.com/mx2118/task-trade-platform/issues)
3. 创建新的 [Issue](https://github.com/mx2118/task-trade-platform/issues/new)

## 🔄 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本更新信息。
