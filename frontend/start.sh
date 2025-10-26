#!/bin/bash

echo "正在启动RAG问答助手前端..."
echo ""
echo "请确保后端服务已启动在 http://localhost:8080"
echo "前端服务将启动在 http://localhost:3000"
echo ""
read -p "按回车键继续..."

echo "安装依赖..."
npm install

echo ""
echo "启动开发服务器..."
npm run dev
