#!/bin/bash
# 个人记账系统 - 本地部署启动脚本

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$PROJECT_ROOT"

# 清理函数：退出时停止后端
cleanup() {
    if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
        echo ""
        echo "正在停止后端服务..."
        kill "$BACKEND_PID" 2>/dev/null || true
    fi
    exit 0
}
trap cleanup EXIT INT TERM

echo "=========================================="
echo "  个人记账系统 - 启动部署"
echo "=========================================="

# 检查依赖
echo ""
echo "[1/3] 检查依赖..."
if ! command -v go &> /dev/null; then
    echo "错误: 未安装 Go，请先安装 Go 1.21+"
    exit 1
fi
if ! command -v npm &> /dev/null; then
    echo "错误: 未安装 Node.js/npm，请先安装"
    exit 1
fi

# 启动后端
echo ""
echo "[2/3] 启动后端服务 (端口 8080)..."
cd "$PROJECT_ROOT/backend"
go mod download 2>/dev/null || true
go build -o finance-server . 2>/dev/null || true
if [ -f finance-server ]; then
    ./finance-server &
else
    go run main.go &
fi
BACKEND_PID=$!
cd "$PROJECT_ROOT"

# 等待后端就绪
echo "等待后端就绪..."
for i in $(seq 1 20); do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "后端已就绪 ✓"
        break
    fi
    if [ $i -eq 20 ]; then
        echo "警告: 后端启动超时，请检查 8080 端口"
    fi
    sleep 1
done

# 启动前端
echo ""
echo "[3/3] 启动前端开发服务器 (端口 3000)..."
cd "$PROJECT_ROOT/frontend"
if [ ! -d node_modules ]; then
    echo "安装前端依赖..."
    npm install
fi
echo ""
echo "=========================================="
echo "  部署完成！请在浏览器访问："
echo "  http://localhost:3000"
echo "=========================================="
npm start
