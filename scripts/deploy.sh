#!/bin/bash

# Gin Template 自动化部署脚本
# 作者: TheIng
# 版本: 1.0.0

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 命令未找到，请先安装"
        exit 1
    fi
}

# 检查环境
check_environment() {
    log_info "检查部署环境..."
    
    # 检查必要的命令
    check_command "docker"
    check_command "docker-compose"
    check_command "git"
    
    # 检查 Docker 服务状态
    if ! docker info &> /dev/null; then
        log_error "Docker 服务未运行"
        exit 1
    fi
    
    log_success "环境检查通过"
}

# 备份数据
backup_data() {
    log_info "备份数据..."
    
    BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    
    # 备份数据库（如果存在）
    if docker-compose ps postgres | grep -q "Up"; then
        log_info "备份数据库..."
        docker-compose exec -T postgres pg_dump -U postgres gin_template > "$BACKUP_DIR/database.sql"
        log_success "数据库备份完成"
    fi
    
    # 备份 Redis（如果存在）
    if docker-compose ps redis | grep -q "Up"; then
        log_info "备份 Redis 数据..."
        docker-compose exec -T redis redis-cli BGSAVE
        sleep 2
        docker cp $(docker-compose ps -q redis):/data/dump.rdb "$BACKUP_DIR/redis.rdb"
        log_success "Redis 备份完成"
    fi
    
    # 备份配置文件
    cp -r ./config "$BACKUP_DIR/" 2>/dev/null || true
    cp .env "$BACKUP_DIR/" 2>/dev/null || true
    
    log_success "数据备份完成: $BACKUP_DIR"
}

# 更新代码
update_code() {
    log_info "更新代码..."
    
    # 拉取最新代码
    git fetch origin
    git reset --hard origin/main
    
    log_success "代码更新完成"
}

# 构建镜像
build_images() {
    log_info "构建 Docker 镜像..."
    
    # 构建应用镜像
    docker-compose build --no-cache app
    
    log_success "镜像构建完成"
}

# 运行测试
run_tests() {
    log_info "运行测试..."
    
    # 运行单元测试
    go test -v ./...
    
    log_success "测试通过"
}

# 部署服务
deploy_services() {
    log_info "部署服务..."
    
    # 停止现有服务
    docker-compose down
    
    # 启动服务
    docker-compose up -d
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 30
    
    # 检查服务状态
    check_service_health
    
    log_success "服务部署完成"
}

# 检查服务健康状态
check_service_health() {
    log_info "检查服务健康状态..."
    
    # 检查应用服务
    if curl -f http://localhost:8080/api/health/ &> /dev/null; then
        log_success "应用服务正常"
    else
        log_error "应用服务异常"
        return 1
    fi
    
    # 检查数据库服务
    if docker-compose ps postgres | grep -q "Up"; then
        log_success "数据库服务正常"
    else
        log_error "数据库服务异常"
        return 1
    fi
    
    # 检查 Redis 服务
    if docker-compose ps redis | grep -q "Up"; then
        log_success "Redis 服务正常"
    else
        log_error "Redis 服务异常"
        return 1
    fi
}

# 清理旧镜像
cleanup() {
    log_info "清理旧镜像..."
    
    # 清理未使用的镜像
    docker image prune -f
    
    # 清理旧的备份（保留最近7天）
    find ./backups -type d -mtime +7 -exec rm -rf {} \; 2>/dev/null || true
    
    log_success "清理完成"
}

# 回滚部署
rollback() {
    log_warning "开始回滚部署..."
    
    # 获取最近的备份
    LATEST_BACKUP=$(ls -t ./backups/ | head -n 1)
    
    if [ -z "$LATEST_BACKUP" ]; then
        log_error "没有找到备份文件"
        exit 1
    fi
    
    BACKUP_DIR="./backups/$LATEST_BACKUP"
    
    # 恢复数据库
    if [ -f "$BACKUP_DIR/database.sql" ]; then
        log_info "恢复数据库..."
        docker-compose exec -T postgres psql -U postgres -d gin_template < "$BACKUP_DIR/database.sql"
    fi
    
    # 恢复 Redis
    if [ -f "$BACKUP_DIR/redis.rdb" ]; then
        log_info "恢复 Redis..."
        docker cp "$BACKUP_DIR/redis.rdb" $(docker-compose ps -q redis):/data/dump.rdb
        docker-compose restart redis
    fi
    
    # 回滚代码到上一个版本
    git reset --hard HEAD~1
    
    # 重新部署
    deploy_services
    
    log_success "回滚完成"
}

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  deploy      完整部署流程"
    echo "  update      仅更新代码和重新部署"
    echo "  backup      仅备份数据"
    echo "  rollback    回滚到上一个版本"
    echo "  health      检查服务健康状态"
    echo "  cleanup     清理旧镜像和备份"
    echo "  help        显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 deploy   # 完整部署"
    echo "  $0 update   # 快速更新"
    echo "  $0 backup   # 备份数据"
}

# 主函数
main() {
    case "${1:-deploy}" in
        "deploy")
            check_environment
            backup_data
            update_code
            build_images
            run_tests
            deploy_services
            cleanup
            log_success "部署完成！"
            ;;
        "update")
            check_environment
            backup_data
            update_code
            build_images
            deploy_services
            log_success "更新完成！"
            ;;
        "backup")
            backup_data
            ;;
        "rollback")
            rollback
            ;;
        "health")
            check_service_health
            ;;
        "cleanup")
            cleanup
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
