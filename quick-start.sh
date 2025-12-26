#!/bin/bash

# Go Demo 项目快速启动脚本
# 用途：一键启动 Docker 部署

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 打印标题
print_header() {
    echo -e "${GREEN}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║          Go Demo 项目 Docker 快速启动脚本                  ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 未安装，请先安装 $1"
        exit 1
    fi
}

# 检查 Docker 服务
check_docker_service() {
    if ! docker info &> /dev/null; then
        print_error "Docker 服务未运行，请启动 Docker"
        exit 1
    fi
}

# 检查端口占用
check_port() {
    local port=$1
    local service=$2
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        print_warning "端口 $port 已被占用 ($service)"
        read -p "是否继续？这可能导致服务启动失败 (y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# 显示菜单
show_menu() {
    echo -e "${BLUE}请选择部署方式：${NC}"
    echo "1) 完整部署 (应用 + MySQL + Redis + Nginx)"
    echo "2) 简化部署 (应用 + MySQL + Redis)"
    echo "3) 仅启动依赖服务 (MySQL + Redis)"
    echo "4) 停止所有服务"
    echo "5) 查看服务状态"
    echo "6) 查看服务日志"
    echo "7) 重启服务"
    echo "8) 清理所有数据（危险操作）"
    echo "9) 退出"
    echo
}

# 完整部署
full_deploy() {
    print_info "开始完整部署..."
    cd deployments
    
    print_info "拉取最新镜像..."
    docker-compose pull
    
    print_info "构建应用镜像..."
    docker-compose build
    
    print_info "启动所有服务..."
    docker-compose up -d
    
    print_info "等待服务启动..."
    sleep 10
    
    print_success "服务启动完成！"
    show_services_info
}

# 简化部署
simple_deploy() {
    print_info "开始简化部署..."
    cd deployments
    
    print_info "拉取最新镜像..."
    docker-compose -f docker-compose.simple.yml pull
    
    print_info "构建应用镜像..."
    docker-compose -f docker-compose.simple.yml build
    
    print_info "启动服务..."
    docker-compose -f docker-compose.simple.yml up -d
    
    print_info "等待服务启动..."
    sleep 10
    
    print_success "服务启动完成！"
    show_services_info_simple
}

# 仅启动依赖服务
deps_only() {
    print_info "启动依赖服务 (MySQL + Redis)..."
    cd deployments
    
    docker-compose up -d mysql redis
    
    print_info "等待服务启动..."
    sleep 5
    
    print_success "依赖服务启动完成！"
    echo
    print_info "MySQL: localhost:3306"
    print_info "Redis: localhost:6379"
    echo
    print_info "现在可以在本地运行应用："
    echo -e "${YELLOW}cd .. && go run main.go server --config=./configs/config.dev.yaml${NC}"
}

# 停止所有服务
stop_services() {
    print_info "停止所有服务..."
    cd deployments
    
    if [ -f docker-compose.yml ]; then
        docker-compose down
    fi
    
    if [ -f docker-compose.simple.yml ]; then
        docker-compose -f docker-compose.simple.yml down
    fi
    
    print_success "所有服务已停止"
}

# 查看服务状态
show_status() {
    print_info "服务状态："
    cd deployments
    
    if docker-compose ps 2>/dev/null | grep -q "Up"; then
        docker-compose ps
    elif docker-compose -f docker-compose.simple.yml ps 2>/dev/null | grep -q "Up"; then
        docker-compose -f docker-compose.simple.yml ps
    else
        print_warning "没有运行中的服务"
    fi
}

# 查看日志
show_logs() {
    cd deployments
    
    echo -e "${BLUE}选择要查看的服务日志：${NC}"
    echo "1) 应用 (app)"
    echo "2) MySQL"
    echo "3) Redis"
    echo "4) Nginx"
    echo "5) 所有服务"
    read -p "请选择 (1-5): " log_choice
    
    case $log_choice in
        1)
            docker-compose logs -f app
            ;;
        2)
            docker-compose logs -f mysql
            ;;
        3)
            docker-compose logs -f redis
            ;;
        4)
            docker-compose logs -f nginx
            ;;
        5)
            docker-compose logs -f
            ;;
        *)
            print_error "无效选择"
            ;;
    esac
}

# 重启服务
restart_services() {
    cd deployments
    
    echo -e "${BLUE}选择要重启的服务：${NC}"
    echo "1) 应用 (app)"
    echo "2) MySQL"
    echo "3) Redis"
    echo "4) Nginx"
    echo "5) 所有服务"
    read -p "请选择 (1-5): " restart_choice
    
    case $restart_choice in
        1)
            print_info "重启应用..."
            docker-compose restart app
            ;;
        2)
            print_info "重启 MySQL..."
            docker-compose restart mysql
            ;;
        3)
            print_info "重启 Redis..."
            docker-compose restart redis
            ;;
        4)
            print_info "重启 Nginx..."
            docker-compose restart nginx
            ;;
        5)
            print_info "重启所有服务..."
            docker-compose restart
            ;;
        *)
            print_error "无效选择"
            return
            ;;
    esac
    
    print_success "重启完成"
}

# 清理数据
cleanup() {
    print_warning "⚠️  警告：此操作将删除所有容器、镜像和数据卷！"
    print_warning "⚠️  所有数据库数据将被永久删除！"
    read -p "确定要继续吗？(输入 'yes' 确认): " confirm
    
    if [ "$confirm" != "yes" ]; then
        print_info "操作已取消"
        return
    fi
    
    print_info "停止并删除所有容器..."
    cd deployments
    docker-compose down -v
    docker-compose -f docker-compose.simple.yml down -v 2>/dev/null || true
    
    print_info "删除应用镜像..."
    docker rmi go-demo:latest 2>/dev/null || true
    docker rmi deployments-app 2>/dev/null || true
    docker rmi deployments_app 2>/dev/null || true
    
    print_info "清理未使用的资源..."
    docker system prune -f
    
    print_success "清理完成"
}

# 显示服务信息
show_services_info() {
    echo
    print_success "═══════════════════════════════════════════════════════"
    print_success "服务访问地址："
    echo -e "${GREEN}  • 应用 API:${NC}      http://localhost:8080"
    echo -e "${GREEN}  • Nginx 代理:${NC}    http://localhost"
    echo -e "${GREEN}  • Swagger 文档:${NC}  http://localhost:8080/swagger/index.html"
    echo -e "${GREEN}  • 健康检查:${NC}      http://localhost:8080/health"
    echo
    print_success "数据库连接信息："
    echo -e "${GREEN}  • MySQL:${NC}         localhost:3306"
    echo -e "${GREEN}  • Redis:${NC}         localhost:6379"
    print_success "═══════════════════════════════════════════════════════"
    echo
    
    print_info "测试服务："
    echo -e "${YELLOW}curl http://localhost:8080/health${NC}"
    echo
    
    print_info "查看日志："
    echo -e "${YELLOW}cd deployments && docker-compose logs -f app${NC}"
    echo
}

# 显示简化部署服务信息
show_services_info_simple() {
    echo
    print_success "═══════════════════════════════════════════════════════"
    print_success "服务访问地址："
    echo -e "${GREEN}  • 应用 API:${NC}      http://localhost:8080"
    echo -e "${GREEN}  • Swagger 文档:${NC}  http://localhost:8080/swagger/index.html"
    echo -e "${GREEN}  • 健康检查:${NC}      http://localhost:8080/health"
    echo
    print_success "数据库连接信息："
    echo -e "${GREEN}  • MySQL:${NC}         localhost:3306"
    echo -e "${GREEN}  • Redis:${NC}         localhost:6379"
    print_success "═══════════════════════════════════════════════════════"
    echo
    
    print_info "测试服务："
    echo -e "${YELLOW}curl http://localhost:8080/health${NC}"
    echo
}

# 健康检查
health_check() {
    print_info "执行健康检查..."
    
    # 等待服务启动
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            print_success "应用健康检查通过！"
            curl -s http://localhost:8080/health | jq . 2>/dev/null || curl -s http://localhost:8080/health
            return 0
        fi
        
        print_info "等待服务启动... ($attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
    
    print_error "健康检查失败，请查看日志"
    return 1
}

# 主函数
main() {
    print_header
    
    # 检查必要的命令
    print_info "检查环境..."
    check_command docker
    check_command docker-compose
    check_docker_service
    
    print_success "环境检查通过"
    echo
    
    # 检查端口占用
    check_port 8080 "应用"
    check_port 3306 "MySQL"
    check_port 6379 "Redis"
    check_port 80 "Nginx"
    
    # 显示菜单并处理选择
    while true; do
        show_menu
        read -p "请选择 (1-9): " choice
        echo
        
        case $choice in
            1)
                full_deploy
                health_check
                ;;
            2)
                simple_deploy
                health_check
                ;;
            3)
                deps_only
                ;;
            4)
                stop_services
                ;;
            5)
                show_status
                ;;
            6)
                show_logs
                ;;
            7)
                restart_services
                ;;
            8)
                cleanup
                ;;
            9)
                print_info "退出"
                exit 0
                ;;
            *)
                print_error "无效选择，请重新选择"
                ;;
        esac
        
        echo
        read -p "按回车键继续..." dummy
        clear
        print_header
    done
}

# 运行主函数
main
