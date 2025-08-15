#!/bin/bash

# 测试运行脚本
# 用法: ./scripts/test.sh [选项]
# 选项:
#   -v, --verbose    显示详细输出
#   -c, --coverage   生成覆盖率报告
#   -b, --benchmark  运行基准测试
#   -a, --all        运行所有测试（包括基准测试）

set -e

# 默认参数
VERBOSE=false
COVERAGE=false
BENCHMARK=false
ALL=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -b|--benchmark)
            BENCHMARK=true
            shift
            ;;
        -a|--all)
            ALL=true
            shift
            ;;
        *)
            echo "未知选项: $1"
            echo "用法: $0 [-v|--verbose] [-c|--coverage] [-b|--benchmark] [-a|--all]"
            exit 1
            ;;
    esac
done

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

# 检查Go环境
check_go_env() {
    print_info "检查Go环境..."
    if ! command -v go &> /dev/null; then
        print_error "Go未安装或不在PATH中"
        exit 1
    fi
    
    go_version=$(go version | awk '{print $3}')
    print_success "Go版本: $go_version"
}

# 更新依赖
update_deps() {
    print_info "更新依赖..."
    go mod tidy
    go mod vendor
    print_success "依赖更新完成"
}

# 运行单元测试
run_unit_tests() {
    print_info "运行单元测试..."
    
    local test_args="-mod=vendor"
    
    if [ "$VERBOSE" = true ]; then
        test_args="$test_args -v"
    fi
    
    if [ "$COVERAGE" = true ]; then
        test_args="$test_args -coverprofile=coverage.out"
    fi
    
    # 运行测试
    if go test $test_args ./...; then
        print_success "单元测试通过"
    else
        print_error "单元测试失败"
        exit 1
    fi
}

# 生成覆盖率报告
generate_coverage_report() {
    if [ "$COVERAGE" = true ]; then
        print_info "生成覆盖率报告..."
        
        if [ -f "coverage.out" ]; then
            # 生成HTML报告
            go tool cover -html=coverage.out -o coverage.html
            print_success "覆盖率报告已生成: coverage.html"
            
            # 显示覆盖率摘要
            print_info "覆盖率摘要:"
            go tool cover -func=coverage.out | tail -1
        else
            print_warning "未找到覆盖率文件 coverage.out"
        fi
    fi
}

# 运行基准测试
run_benchmark_tests() {
    if [ "$BENCHMARK" = true ] || [ "$ALL" = true ]; then
        print_info "运行基准测试..."
        
        local bench_args="-mod=vendor -bench=."
        
        if [ "$VERBOSE" = true ]; then
            bench_args="$bench_args -v"
        fi
        
        # 运行基准测试
        if go test $bench_args ./pkg/domain/service/... ./pkg/utils/...; then
            print_success "基准测试完成"
        else
            print_warning "基准测试失败"
        fi
    fi
}

# 运行代码质量检查
run_code_quality_checks() {
    print_info "运行代码质量检查..."
    
    # 检查是否有golangci-lint
    if command -v golangci-lint &> /dev/null; then
        if golangci-lint run; then
            print_success "代码质量检查通过"
        else
            print_warning "代码质量检查发现问题"
        fi
    else
        print_warning "golangci-lint未安装，跳过代码质量检查"
    fi
}

# 主函数
main() {
    print_info "开始测试流程..."
    
    check_go_env
    update_deps
    run_unit_tests
    generate_coverage_report
    run_benchmark_tests
    run_code_quality_checks
    
    print_success "测试流程完成"
}

# 执行主函数
main "$@"
