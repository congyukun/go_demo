#!/bin/bash

# Go Demo 项目部署脚本
# 使用方法: ./scripts/deploy.sh [环境] [操作]
# 环境: dev|test|prod (默认: dev)
# 操作: up|down|restart|logs|status (默认: up)

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
PROJECT_NAME="go-demo"
COMPOSE_FILE="deployments/docker-compose.yml"
ENV_FILE=".env"

# 默认参数
ENVIRONMENT=${1:-dev}
ACTION=${2:-up}

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
    log_info "检查依赖..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    log_success "依赖检查通过"
}

# 创建环境文件
create_env_file() {
    if [ ! -f "$ENV_FILE" ]; then
        log_info "创建环境配置文件..."
        cat > "$ENV_FILE" << EOF
# Go Demo 环境配置
ENVIRONMENT=$ENVIRONMENT
PROJECT_NAME=$PROJECT_NAME

# 数据库配置
MYSQL_ROOT_PASSWORD=123456
MYSQL_DATABASE=go_demo
MYSQL_USER=demo_user
MYSQL_PASSWORD=demo_pass

# 应用配置
APP_PORT=8080
CONFIG_PATH=/app/configs/config.yaml

# JWT 配置
JWT_SECRET=your-secret-key-change-in-production

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
EOF
        log_success "环境配置文件创建完成: $ENV_FILE"
    else
        log_info "环境配置文件已存在: $ENV_FILE"
    fi
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    
    directories=(
        "logs"
        "logs/nginx"
        "data/mysql"
        "data/redis"
    )
    
    for dir in "${directories[@]}"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "$dir"
            log_info "创建目录: $dir"
        fi
    done
    
    log_success "目录创建完成"
}

# 构建应用
build_app() {
    log_info "构建应用..."
    
    # 检查 go.mod 文件
    if [ ! -f "go.mod" ]; then
        log_error "go.mod 文件不存在，请确保在项目根目录执行"
        exit 1
    fi
    
    # 使用 Docker Compose 构建
    docker-compose -f "$COMPOSE_FILE" build --no-cache app
    
    log_success "应用构建完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    # 启动所有服务
    docker-compose -f "$COMPOSE_FILE" up -d
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 10
    
    # 检查服务状态
    check_services_health
    
    log_success "服务启动完成"
}

# 停止服务
stop_services() {
    log_info "停止服务..."
    
    docker-compose -f "$COMPOSE_FILE" down
    
    log_success "服务停止完成"
}

# 重启服务
restart_services() {
    log_info "重启服务..."
    
    stop_services
    sleep 5
    start_services
    
    log_success "服务重启完成"
}

# 查看日志
show_logs() {
    local service=${3:-app}
    log_info "查看 $service 服务日志..."
    
    docker-compose -f "$COMPOSE_FILE" logs -f "$service"
}

# 检查服务健康状态
check_services_health() {
    log_info "检查服务健康状态..."
    
    # 检查应用健康状态
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            log_success "应用服务健康检查通过"
            break
        else
            log_warning "等待应用服务启动... ($attempt/$max_attempts)"
            sleep 2
            ((attempt++))
        fi
    done
    
    if [ $attempt -gt $max_attempts ]; then
        log_error "应用服务健康检查失败"
        show_service_status
        return 1
    fi
    
    # 检查数据库连接
    if docker-compose -f "$COMPOSE_FILE" exec -T mysql mysqladmin ping -h localhost > /dev/null 2>&1; then
        log_success "数据库服务健康检查通过"
    else
        log_warning "数据库服务可能未完全启动"
    fi
    
    # 检查 Redis 连接
    if docker-compose -f "$COMPOSE_FILE" exec -T redis redis-cli ping > /dev/null 2>&1; then
        log_success "Redis 服务健康检查通过"
    else
        log_warning "Redis 服务可能未完全启动"
    fi
}

# 显示服务状态
show_service_status() {
    log_info "服务状态:"
    docker-compose -f "$COMPOSE_FILE" ps
    
    echo ""
    log_info "服务访问地址:"
    echo "  应用服务: http://localhost:8080"
    echo "  健康检查: http://localhost:8080/health"
    echo "  Nginx: http://localhost:80"
    echo "  MySQL: localhost:3306"
    echo "  Redis: localhost:6379"
}

# 清理资源
cleanup() {
    log_info "清理资源..."
    
    # 停止并删除容器
    docker-compose -f "$COMPOSE_FILE" down -v --remove-orphans
    
    # 删除未使用的镜像
    docker image prune -f
    
    log_success "资源清理完成"
}

# 备份数据
backup_data() {
    log_info "备份数据..."
    
    local backup_dir="backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$backup_dir"
    
    # 备份 MySQL 数据
    docker-compose -f "$COMPOSE_FILE" exec -T mysql mysqldump -u root -p123456 go_demo > "$backup_dir/mysql_backup.sql"
    
    # 备份 Redis 数据
    docker-compose -f "$COMPOSE_FILE" exec -T redis redis-cli BGSAVE
    docker cp $(docker-compose -f "$COMPOSE_FILE" ps -q redis):/data/dump.rdb "$backup_dir/redis_backup.rdb"
    
    log_success "数据备份完成: $backup_dir"
}

# 显示帮助信息
show_help() {
    echo "Go Demo 项目部署脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 [环境] [操作] [服务]"
    echo ""
    echo "环境:"
    echo "  dev     开发环境 (默认)"
    echo "  test    测试环境"
    echo "  prod    生产环境"
    echo ""
    echo "操作:"
    echo "  up      启动服务 (默认)"
    echo "  down    停止服务"
    echo "  restart 重启服务"
    echo "  logs    查看日志"
    echo "  status  查看状态"
    echo "  build   构建应用"
    echo "  backup  备份数据"
    echo "  cleanup 清理资源"
    echo "  help    显示帮助"
    echo ""
    echo "服务 (仅用于 logs 操作):"
    echo "  app     应用服务 (默认)"
    echo "  mysql   数据库服务"
    echo "  redis   缓存服务"
    echo "  nginx   代理服务"
    echo ""
    echo "示例:"
    echo "  $0                    # 启动开发环境"
    echo "  $0 prod up           # 启动生产环境"
    echo "  $0 dev logs app      # 查看应用日志"
    echo "  $0 dev status        # 查看服务状态"
}

# 主函数
main() {
    echo "=========================================="
    echo "Go Demo 项目部署脚本"
    echo "环境: $ENVIRONMENT"
    echo "操作: $ACTION"
    echo "=========================================="
    
    case "$ACTION" in
        "up")
            check_dependencies
            create_env_file
            create_directories
            build_app
            start_services
            show_service_status
            ;;
        "down")
            stop_services
            ;;
        "restart")
            restart_services
            show_service_status
            ;;
        "logs")
            show_logs "$@"
            ;;
        "status")
            show_service_status
            ;;
        "build")
            check_dependencies
            build_app
            ;;
        "backup")
            backup_data
            ;;
        "cleanup")
            cleanup
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "未知操作: $ACTION"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"