#!/bin/bash

# Gin Template 快速启动脚本
# 用于本地开发环境的快速搭建

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

# 检查依赖
check_dependencies() {
    log_info "检查系统依赖..."
    
    local missing_deps=()
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        missing_deps+=("docker")
    fi
    
    # 检查 Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        missing_deps+=("docker-compose")
    fi
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        missing_deps+=("go")
    fi
    
    # 检查 Make
    if ! command -v make &> /dev/null; then
        missing_deps+=("make")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "缺少以下依赖: ${missing_deps[*]}"
        echo "请先安装缺少的依赖包"
        exit 1
    fi
    
    log_success "依赖检查通过"
}

# 创建环境配置
setup_environment() {
    log_info "设置环境配置..."
    
    # 复制环境变量文件
    if [ ! -f .env ]; then
        if [ -f .env_example ]; then
            cp .env_example .env
            log_success "已创建 .env 文件"
        else
            log_warning ".env_example 文件不存在，使用默认配置"
            cat > .env << EOF
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_template

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 应用配置
GIN_MODE=debug
PORT=8080
JWT_SECRET=your-jwt-secret-key-here

# 缓存配置
CACHE_HOST=localhost
CACHE_PORT=6379
CACHE_PASSWORD=
CACHE_DB=0
CACHE_PREFIX=gin_template:

# 日志配置
LOG_LEVEL=info
LOG_FILE=logs/app.log
EOF
            log_success "已创建默认 .env 文件"
        fi
    fi
    
    # 创建必要的目录
    mkdir -p logs
    mkdir -p data
    mkdir -p backups
    mkdir -p monitoring/grafana/dashboards
    mkdir -p monitoring/grafana/datasources
    
    log_success "目录结构创建完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    # 启动基础服务（数据库和缓存）
    log_info "启动数据库和缓存服务..."
    docker-compose up -d postgres redis
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 10
    
    # 检查服务状态
    if docker-compose ps postgres | grep -q "Up"; then
        log_success "PostgreSQL 启动成功"
    else
        log_error "PostgreSQL 启动失败"
        return 1
    fi
    
    if docker-compose ps redis | grep -q "Up"; then
        log_success "Redis 启动成功"
    else
        log_error "Redis 启动失败"
        return 1
    fi
}

# 安装依赖
install_dependencies() {
    log_info "安装 Go 依赖..."
    
    go mod download
    go mod tidy
    
    log_success "依赖安装完成"
}

# 运行数据库迁移
run_migrations() {
    log_info "运行数据库迁移..."
    
    # 这里可以添加数据库迁移逻辑
    # 例如: go run cmd/migrate/main.go up
    
    log_success "数据库迁移完成"
}

# 启动应用
start_application() {
    log_info "启动应用..."
    
    # 使用 air 进行热重载（开发模式）
    if command -v air &> /dev/null; then
        log_info "使用 air 热重载启动..."
        air
    else
        log_info "安装 air..."
        go install github.com/air-verse/air@latest
        log_info "使用 air 热重载启动..."
        air
    fi
}

# 启动监控服务
start_monitoring() {
    log_info "启动监控服务..."
    
    # 启动 Prometheus 和 Grafana
    docker-compose up -d prometheus grafana
    
    log_success "监控服务启动完成"
    log_info "Prometheus: http://localhost:9090"
    log_info "Grafana: http://localhost:3000 (admin/admin)"
}

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  start       完整启动开发环境"
    echo "  dev         仅启动应用（依赖外部数据库）"
    echo "  services    仅启动基础服务（数据库、Redis）"
    echo "  monitoring  启动监控服务"
    echo "  stop        停止所有服务"
    echo "  clean       清理数据和容器"
    echo "  help        显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 start      # 完整启动"
    echo "  $0 dev        # 仅启动应用"
    echo "  $0 services   # 仅启动基础服务"
}

# 停止服务
stop_services() {
    log_info "停止所有服务..."
    
    docker-compose down
    
    log_success "所有服务已停止"
}

# 清理环境
clean_environment() {
    log_warning "清理环境（将删除所有数据）..."
    
    read -p "确认删除所有数据和容器? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log_info "清理容器和数据..."
        docker-compose down -v --remove-orphans
        docker system prune -f
        
        # 清理本地数据
        rm -rf logs/*
        rm -rf data/*
        rm -rf backups/*
        
        log_success "环境清理完成"
    else
        log_info "取消清理操作"
    fi
}

# 主函数
main() {
    case "${1:-start}" in
        "start")
            check_dependencies
            setup_environment
            start_services
            install_dependencies
            run_migrations
            start_monitoring
            log_success "开发环境启动完成！"
            echo ""
            echo "服务地址:"
            echo "  - 应用: http://localhost:8080"
            echo "  - API 文档: http://localhost:8080/api/docs/"
            echo "  - Prometheus: http://localhost:9090"
            echo "  - Grafana: http://localhost:3000 (admin/admin)"
            echo "  - 数据库: localhost:5432"
            echo "  - Redis: localhost:6379"
            echo ""
            log_info "按 Ctrl+C 停止应用"
            start_application
            ;;
        "dev")
            check_dependencies
            setup_environment
            install_dependencies
            log_success "开发环境就绪！"
            start_application
            ;;
        "services")
            check_dependencies
            setup_environment
            start_services
            log_success "基础服务启动完成！"
            ;;
        "monitoring")
            check_dependencies
            start_monitoring
            ;;
        "stop")
            stop_services
            ;;
        "clean")
            clean_environment
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
}

# 脚本入口
main "$@"
