@echo off
echo 正在启动RAG问答助手前端...
echo.
echo 请确保后端服务已启动在 http://localhost:8080
echo 前端服务将启动在 http://localhost:3000
echo.
echo 按任意键继续...
pause >nul

echo 安装依赖...
call npm install

echo.
echo 启动开发服务器...
call npm run dev
