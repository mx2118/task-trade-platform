.PHONY: build run test clean docker-build docker-run migrate seed swagger

# 变量定义
APP_NAME=task-platform-api
VERSION=1.0.0
BUILD_DIR=build

# 构建
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/server/main.go
	@echo "Build completed!"

# 运行
run:
	@echo "Running $(APP_NAME)..."
	@go run cmd/server/main.go -config=./configs/config.yaml

# 测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 支付测试
test-payment:
	@echo "Running payment tests..."
	@go test -v ./pkg/payment/...

# 单元测试
test-unit:
	@echo "Running unit tests..."
	@go test -v -short ./...

# 集成测试
test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...

# 清理
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@go clean

# 数据库迁移
migrate:
	@echo "Running database migrations..."
	@go run scripts/migrate.go

# 数据填充
seed:
	@echo "Seeding database..."
	@go run scripts/seed.go

# 生成Swagger文档
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs

# Docker构建
docker-build:
	@echo "Building Docker image..."
	@docker build -f deployments/docker/Dockerfile -t $(APP_NAME):$(VERSION) .

# Docker运行
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(APP_NAME):$(VERSION)

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码检查
lint:
	@echo "Running linter..."
	@golangci-lint run

# 安装依赖
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# 生成mock文件
mock:
	@echo "Generating mock files..."
	@mockgen -source=internal/repository/interfaces/user.go -destination=tests/mock/user_mock.go
	@mockgen -source=internal/repository/interfaces/task.go -destination=tests/mock/task_mock.go

# 性能测试
bench:
	@echo "Running benchmark tests..."
	@go test -bench=. -benchmem ./...

# 竞争检测
race:
	@echo "Running race condition tests..."
	@go test -race ./...

# 代码覆盖率
coverage:
	@echo "Running test coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# 安装开发工具
install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 创建目录结构
init-dirs:
	@echo "Creating directory structure..."
	@mkdir -p {cmd/server,internal/{api/{v1/{handlers,middleware,routes,validators},docs},config,models,repository/{interfaces,mysql,redis},service/{interfaces},pkg/{utils,logger,errors,queue},domain/{user,task,order}},pkg/{config,database,middleware,validator},scripts,configs,docs,deployments/{docker,kubernetes,nginx},tests/{integration,unit,mock}}