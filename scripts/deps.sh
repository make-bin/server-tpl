#!/bin/bash

# 依赖管理脚本
# Dependency Management Script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "${SCRIPT_DIR}")"

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

# 显示帮助信息
show_help() {
    cat << EOF
依赖管理脚本 - Dependency Management Script

用法: $0 [命令]

命令:
  init        初始化依赖管理 (Initialize dependency management)
  update      更新依赖版本 (Update dependencies)
  vendor      生成vendor目录 (Generate vendor directory)
  clean       清理vendor目录 (Clean vendor directory)
  verify      验证依赖完整性 (Verify dependencies)
  check       检查过期依赖 (Check outdated dependencies)
  build       使用vendor构建 (Build with vendor)
  test        使用vendor测试 (Test with vendor)
  security    安全性检查 (Security check)
  help        显示此帮助信息 (Show this help)

示例:
  $0 init     # 初始化依赖管理
  $0 update   # 更新所有依赖
  $0 vendor   # 重新生成vendor目录
  $0 build    # 使用vendor构建项目

EOF
}

# 初始化依赖管理
init_deps() {
    log_info "初始化依赖管理..."
    
    cd "${PROJECT_ROOT}"
    
    # 检查go.mod文件
    if [[ ! -f "go.mod" ]]; then
        log_error "go.mod文件不存在，请先运行 go mod init"
        exit 1
    fi
    
    # 整理依赖
    log_info "整理依赖关系..."
    go mod tidy
    
    # 下载依赖
    log_info "下载依赖..."
    go mod download
    
    # 生成vendor目录
    log_info "生成vendor目录..."
    go mod vendor
    
    # 验证依赖
    log_info "验证依赖完整性..."
    go mod verify
    
    log_success "依赖管理初始化完成"
}

# 更新依赖
update_deps() {
    log_info "更新依赖版本..."
    
    cd "${PROJECT_ROOT}"
    
    # 备份当前go.mod和go.sum
    if [[ -f "go.mod" ]]; then
        cp go.mod go.mod.backup
        log_info "已备份go.mod到go.mod.backup"
    fi
    
    if [[ -f "go.sum" ]]; then
        cp go.sum go.sum.backup
        log_info "已备份go.sum到go.sum.backup"
    fi
    
    # 更新所有依赖
    log_info "更新所有依赖到最新版本..."
    go get -u ./...
    
    # 整理依赖
    go mod tidy
    
    # 重新生成vendor
    vendor_deps
    
    # 验证更新
    log_info "验证更新后的依赖..."
    go mod verify
    
    log_success "依赖更新完成"
    log_warning "请测试应用以确保更新后的依赖工作正常"
}

# 生成vendor目录
vendor_deps() {
    log_info "生成vendor目录..."
    
    cd "${PROJECT_ROOT}"
    
    # 清理现有vendor目录
    if [[ -d "vendor" ]]; then
        log_info "清理现有vendor目录..."
        rm -rf vendor
    fi
    
    # 生成新的vendor目录
    log_info "创建vendor目录..."
    go mod vendor
    
    # 检查vendor目录大小
    if [[ -d "vendor" ]]; then
        vendor_size=$(du -sh vendor | cut -f1)
        log_success "vendor目录生成完成，大小: ${vendor_size}"
    else
        log_error "vendor目录生成失败"
        exit 1
    fi
}

# 清理vendor目录
clean_deps() {
    log_info "清理vendor目录..."
    
    cd "${PROJECT_ROOT}"
    
    if [[ -d "vendor" ]]; then
        rm -rf vendor
        log_success "vendor目录已清理"
    else
        log_info "vendor目录不存在，无需清理"
    fi
    
    # 清理模块缓存
    log_info "清理模块缓存..."
    go clean -modcache
    
    log_success "清理完成"
}

# 验证依赖
verify_deps() {
    log_info "验证依赖完整性..."
    
    cd "${PROJECT_ROOT}"
    
    # 验证go.mod和go.sum
    log_info "验证go.mod和go.sum..."
    go mod verify
    
    # 检查vendor目录
    if [[ -d "vendor" ]]; then
        log_info "检查vendor目录完整性..."
        
        # 比较vendor和go.mod的一致性
        temp_vendor=$(mktemp -d)
        cp -r vendor "${temp_vendor}/vendor_current"
        
        go mod vendor -o "${temp_vendor}/vendor_new"
        
        if diff -r "${temp_vendor}/vendor_current" "${temp_vendor}/vendor_new" > /dev/null; then
            log_success "vendor目录与go.mod一致"
        else
            log_warning "vendor目录与go.mod不一致，建议重新生成vendor"
        fi
        
        rm -rf "${temp_vendor}"
    else
        log_warning "vendor目录不存在"
    fi
    
    log_success "依赖验证完成"
}

# 检查过期依赖
check_outdated() {
    log_info "检查过期依赖..."
    
    cd "${PROJECT_ROOT}"
    
    # 检查可用更新
    log_info "检查可用的依赖更新..."
    go list -u -m all | grep -v "^github.com/make-bin/server-tpl" | while read -r line; do
        if [[ $line == *" ["* ]]; then
            echo "$line"
        fi
    done
    
    log_success "过期依赖检查完成"
}

# 使用vendor构建
build_with_vendor() {
    log_info "使用vendor目录构建项目..."
    
    cd "${PROJECT_ROOT}"
    
    if [[ ! -d "vendor" ]]; then
        log_error "vendor目录不存在，请先运行: $0 vendor"
        exit 1
    fi
    
    # 使用vendor构建
    log_info "构建应用..."
    go build -mod=vendor -o bin/server ./cmd/server
    
    if [[ $? -eq 0 ]]; then
        log_success "构建成功"
        if [[ -f "bin/server" ]]; then
            file_size=$(ls -lah bin/server | awk '{print $5}')
            log_info "可执行文件大小: ${file_size}"
        fi
    else
        log_error "构建失败"
        exit 1
    fi
}

# 使用vendor测试
test_with_vendor() {
    log_info "使用vendor目录运行测试..."
    
    cd "${PROJECT_ROOT}"
    
    if [[ ! -d "vendor" ]]; then
        log_error "vendor目录不存在，请先运行: $0 vendor"
        exit 1
    fi
    
    # 使用vendor运行测试
    log_info "运行测试..."
    go test -mod=vendor ./...
    
    if [[ $? -eq 0 ]]; then
        log_success "所有测试通过"
    else
        log_error "部分测试失败"
        exit 1
    fi
}

# 安全性检查
security_check() {
    log_info "执行安全性检查..."
    
    cd "${PROJECT_ROOT}"
    
    # 检查已知漏洞
    if command -v govulncheck &> /dev/null; then
        log_info "使用govulncheck检查漏洞..."
        govulncheck ./...
    else
        log_warning "govulncheck未安装，跳过漏洞检查"
        log_info "安装govulncheck: go install golang.org/x/vuln/cmd/govulncheck@latest"
    fi
    
    # 检查许可证
    log_info "检查依赖许可证..."
    go list -m -json all | jq -r '.Path' | while read -r module; do
        if [[ -n "$module" && "$module" != "github.com/make-bin/server-tpl" ]]; then
            echo "模块: $module"
        fi
    done
    
    log_success "安全性检查完成"
}

# 显示依赖信息
show_info() {
    log_info "依赖信息概览..."
    
    cd "${PROJECT_ROOT}"
    
    echo ""
    echo "Go版本信息:"
    go version
    
    echo ""
    echo "模块信息:"
    if [[ -f "go.mod" ]]; then
        head -3 go.mod
    fi
    
    echo ""
    echo "直接依赖数量:"
    if [[ -f "go.mod" ]]; then
        grep -c "^\s*github.com\|^\s*golang.org\|^\s*gopkg.in\|^\s*gorm.io" go.mod || echo "0"
    fi
    
    echo ""
    echo "vendor目录状态:"
    if [[ -d "vendor" ]]; then
        vendor_size=$(du -sh vendor 2>/dev/null | cut -f1)
        vendor_modules=$(find vendor -name "go.mod" | wc -l)
        echo "  大小: ${vendor_size}"
        echo "  模块数量: ${vendor_modules}"
    else
        echo "  不存在"
    fi
    
    echo ""
}

# 主函数
main() {
    case "${1:-help}" in
        "init")
            init_deps
            ;;
        "update")
            update_deps
            ;;
        "vendor")
            vendor_deps
            ;;
        "clean")
            clean_deps
            ;;
        "verify")
            verify_deps
            ;;
        "check")
            check_outdated
            ;;
        "build")
            build_with_vendor
            ;;
        "test")
            test_with_vendor
            ;;
        "security")
            security_check
            ;;
        "info")
            show_info
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "未知命令: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
