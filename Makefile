# Makefile for gin-template project

.PHONY: help build test clean lint fmt security install-tools run dev

# 默认目标
help: ## 显示帮助信息
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# 构建
build: ## 构建应用
	@echo "Building application..."
	go build -o bin/gin-template main.go

# 构建生产版本
build-prod: ## 构建生产版本
	@echo "Building production binary..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/gin-template-linux main.go
	CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o bin/gin-template-darwin main.go
	CGO_ENABLED=0 GOOS=windows go build -a -installsuffix cgo -o bin/gin-template-windows.exe main.go

# 运行
run: ## 运行应用
	@echo "Running application..."
	go run main.go

# 开发模式
dev: ## 开发模式运行（热重载）
	@echo "Running in development mode..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Installing air for hot reload..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

# 测试
test: ## 运行测试
	@echo "Running tests..."
	go test -v ./...

# 测试覆盖率
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 基准测试
bench: ## 运行基准测试
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# 代码格式化
fmt: ## 格式化代码
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# 代码检查
lint: ## 运行代码检查
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
		$$(go env GOPATH)/bin/golangci-lint run; \
	fi

# 静态分析
staticcheck: ## 运行静态分析
	@echo "Running staticcheck..."
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
		staticcheck ./...; \
	fi

# 安全扫描
security: ## 运行安全扫描
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

# 依赖检查
deps: ## 检查依赖
	@echo "Checking dependencies..."
	@if command -v go-mod-tidy >/dev/null 2>&1; then \
		go-mod-tidy --check; \
	else \
		go mod tidy; \
		go mod verify; \
	fi

# 更新依赖
update-deps: ## 更新依赖
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# 安装开发工具
install-tools: ## 安装开发工具
	@echo "Installing development tools..."
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
# 	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/matm/gocov-html@latest

# 代码质量检查
quality: fmt lint staticcheck security ## 运行所有代码质量检查

# 生成文档
docs: ## 生成文档
	@echo "Generating documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g main.go -o docs; \
	else \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g main.go -o docs; \
	fi

# 清理
clean: ## 清理构建文件和临时文件
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f *.log
	go clean -cache
	go clean -testcache

# Docker 构建
docker-build: ## 构建 Docker 镜像
	@echo "Building Docker image..."
	docker build -t gin-template:latest .

# Docker 运行
docker-run: ## 运行 Docker 容器
	@echo "Running Docker container..."
	docker run -p 8080:8080 gin-template:latest

# 代码分析
analyze: ## 运行完整的代码分析
	@echo "Running complete code analysis..."
	@echo "1. Running tests..."
	make test-coverage
	@echo "2. Running benchmarks..."
	make bench
	@echo "3. Running linter..."
	make lint
	@echo "4. Running security scan..."
	make security
	@echo "5. Running static analysis..."
	make staticcheck
	@echo "Analysis complete!"

# 预提交检查
pre-commit: quality test ## Git 预提交检查
	@echo "Pre-commit checks passed!"

# 安装 Git hooks
install-hooks: ## 安装 Git hooks
	@echo "Installing Git hooks..."
	cp scripts/pre-commit .git/hooks/
	chmod +x .git/hooks/pre-commit

# 发布准备
release: clean test quality build-prod ## 发布准备
	@echo "Release ready!"
	@echo "Binaries created in bin/ directory"

# 监控
monitor: ## 启动监控（需要 Prometheus 和 Grafana）
	@echo "Starting monitoring stack..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f monitoring/docker-compose.yml up -d; \
	else \
		echo "Docker Compose not found. Please install Docker Compose."; \
	fi

# 性能测试
load-test: ## 运行负载测试
	@echo "Running load tests..."
	@if command -v hey >/dev/null 2>&1; then \
		hey -n 1000 -c 10 http://localhost:8080/api/health/; \
	else \
		echo "Installing hey..."; \
		go install github.com/rakyll/hey@latest; \
		hey -n 1000 -c 10 http://localhost:8080/api/health/; \
	fi

# 数据库迁移
migrate: ## 运行数据库迁移
	@echo "Running database migrations..."
	@if [ -f migrations/migrate.go ]; then \
		go run migrations/migrate.go; \
	else \
		echo "No migration file found."; \
	fi

# 备份
backup: ## 备份数据库
	@echo "Creating database backup..."
	@if command -v pg_dump >/dev/null 2>&1; then \
		pg_dump $$DATABASE_URL > backup_$(date +%Y%m%d_%H%M%S).sql; \
	else \
		echo "pg_dump not found. Please install PostgreSQL client tools."; \
	fi

# 初始化项目
init: install-tools deps ## 初始化新项目环境
	@echo "Project initialized!"
	@echo "Run 'make dev' to start development server."

# 版本信息
version: ## 显示版本信息
	@echo "gin-template v2.0.0"
	@echo "Go version: $$(go version)"
	@echo "Git commit: $$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build time: $$(date)"

# 检查工具
check-tools: ## 检查必需的工具
	@echo "Checking required tools..."
	@for tool in go docker git; do \
		if command -v $$tool >/dev/null 2>&1; then \
			echo "✓ $$tool is installed"; \
		else \
			echo "✗ $$tool is not installed"; \
		fi; \
	done

# 工程化命令
dev-setup: ## 快速搭建开发环境
	@echo "Setting up development environment..."
	@./scripts/quick-start.sh start

dev-stop: ## 停止开发环境
	@echo "Stopping development environment..."
	@./scripts/quick-start.sh stop

dev-clean: ## 清理开发环境
	@echo "Cleaning development environment..."
	@./scripts/quick-start.sh clean

# 部署命令
deploy-prod: ## 部署到生产环境
	@echo "Deploying to production..."
	@./scripts/deploy.sh deploy

deploy-update: ## 更新部署
	@echo "Updating deployment..."
	@./scripts/deploy.sh update

deploy-backup: ## 备份数据
	@echo "Creating backup..."
	@./scripts/deploy.sh backup

deploy-rollback: ## 回滚部署
	@echo "Rolling back deployment..."
	@./scripts/deploy.sh rollback

deploy-health: ## 检查部署健康状态
	@echo "Checking deployment health..."
	@./scripts/deploy.sh health

# 监控命令
monitor-start: ## 启动监控服务
	@echo "Starting monitoring services..."
	docker-compose up -d prometheus grafana

monitor-stop: ## 停止监控服务
	@echo "Stopping monitoring services..."
	docker-compose down prometheus grafana

monitor-logs: ## 查看监控日志
	@echo "Showing monitoring logs..."
	docker-compose logs -f prometheus grafana

# 容器管理
docker-start: ## 启动所有服务
	@echo "Starting all services..."
	docker-compose up -d

docker-stop: ## 停止所有服务
	@echo "Stopping all services..."
	docker-compose down

docker-restart: ## 重启所有服务
	@echo "Restarting all services..."
	docker-compose restart

docker-logs: ## 查看服务日志
	@echo "Showing service logs..."
	docker-compose logs -f

docker-status: ## 查看服务状态
	@echo "Showing service status..."
	docker-compose ps

# 配置管理
config-validate: ## 验证配置文件
	@echo "Validating configuration..."
	@if [ -f config/environments.yml ]; then \
		echo "✓ Configuration file exists"; \
	else \
		echo "✗ Configuration file missing"; \
	fi

config-dev: ## 设置开发环境配置
	@echo "Setting up development configuration..."
	@cp config/environments.yml config/application.yml
	@echo "Development configuration loaded"

config-prod: ## 设置生产环境配置
	@echo "Setting up production configuration..."
	@cp config/environments.yml config/application.yml
	@echo "Production configuration loaded"

# 数据库管理
db-start: ## 启动数据库服务
	@echo "Starting database..."
	docker-compose up -d postgres

db-stop: ## 停止数据库服务
	@echo "Stopping database..."
	docker-compose stop postgres

db-reset: ## 重置数据库
	@echo "Resetting database..."
	docker-compose down postgres
	docker volume rm gin-template_postgres_data 2>/dev/null || true
	docker-compose up -d postgres

db-shell: ## 连接数据库 shell
	@echo "Connecting to database shell..."
	docker-compose exec postgres psql -U postgres gin_template

# 缓存管理
cache-start: ## 启动缓存服务
	@echo "Starting Redis..."
	docker-compose up -d redis

cache-stop: ## 停止缓存服务
	@echo "Stopping Redis..."
	docker-compose stop redis

cache-reset: ## 重置缓存
	@echo "Resetting Redis..."
	docker-compose exec redis redis-cli FLUSHALL

cache-shell: ## 连接 Redis shell
	@echo "Connecting to Redis shell..."
	docker-compose exec redis redis-cli

# 性能分析
profile-cpu: ## CPU 性能分析
	@echo "Starting CPU profiling..."
	go tool pprof http://localhost:8080/debug/pprof/profile

profile-memory: ## 内存性能分析
	@echo "Starting memory profiling..."
	go tool pprof http://localhost:8080/debug/pprof/heap

profile-goroutine: ## Goroutine 分析
	@echo "Starting goroutine profiling..."
	go tool pprof http://localhost:8080/debug/pprof/goroutine

# 日志管理
logs-app: ## 查看应用日志
	@echo "Showing application logs..."
	tail -f logs/app.log

logs-access: ## 查看访问日志
	@echo "Showing access logs..."
	docker-compose logs -f nginx

logs-all: ## 查看所有服务日志
	@echo "Showing all service logs..."
	docker-compose logs -f

# 安全检查
security-deps: ## 检查依赖安全
	@echo "Checking dependency security..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
		govulncheck ./...; \
	fi

security-licenses: ## 检查许可证
	@echo "Checking licenses..."
	@if command -v go-licenses >/dev/null 2>&1; then \
		go-licenses check ./...; \
	else \
		echo "Installing go-licenses..."; \
		go install github.com/google/go-licenses@latest; \
		go-licenses check ./...; \
	fi

# 完整工作流
workflow-dev: ## 完整的开发工作流
	@echo "Running complete development workflow..."
	make deps
	make quality
	make test-coverage
	make docs
	@echo "Development workflow completed!"

workflow-ci: ## 模拟 CI 工作流
	@echo "Running CI workflow..."
	make deps
	make test
	make lint
	make security
	make build
	@echo "CI workflow completed!"

workflow-release: ## 发布工作流
	@echo "Running release workflow..."
	make clean
	make test-coverage
	make quality
	make security
	make build-prod
	make docs
	@echo "Release workflow completed!"

# 帮助信息
help-extended: ## 显示扩展帮助信息
	@echo "=== Extended Help ==="
	@echo ""
	@echo "Development Commands:"
	@echo "  dev-setup      - Quick setup development environment"
	@echo "  dev-stop       - Stop development environment"
	@echo "  dev-clean      - Clean development environment"
	@echo ""
	@echo "Deployment Commands:"
	@echo "  deploy-prod    - Deploy to production"
	@echo "  deploy-update   - Update deployment"
	@echo "  deploy-backup  - Create backup"
	@echo "  deploy-rollback - Rollback deployment"
	@echo "  deploy-health   - Check deployment health"
	@echo ""
	@echo "Monitoring Commands:"
	@echo "  monitor-start   - Start monitoring services"
	@echo "  monitor-stop    - Stop monitoring services"
	@echo "  monitor-logs    - View monitoring logs"
	@echo ""
	@echo "Container Commands:"
	@echo "  docker-start    - Start all services"
	@echo "  docker-stop     - Stop all services"
	@echo "  docker-restart  - Restart all services"
	@echo "  docker-logs     - View service logs"
	@echo "  docker-status   - View service status"
	@echo ""
	@echo "For basic commands, run: make help"
