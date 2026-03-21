#!/bin/bash
# HTTP MCP Server 測試腳本
# 依序測試 MCP over HTTP 的核心端點，驗證 Server 是否正常運作。
#
# 使用方式：
#   1. 先啟動 Server: go run ./examples/U07-http-transport/
#   2. 另開終端執行: ./examples/U07-http-transport/test.sh
#
# 測試案例：
#   1. Health Check     — GET /health，預期 200 OK
#   2. Initialize       — POST /mcp，建立 JSON-RPC session，預期回傳 serverInfo
#   3. List Tools       — POST /mcp，列舉已註冊 Tool，預期回傳 add tool
#   4. Call Tool (整數)  — POST /mcp，呼叫 add(10,20)，預期回傳 30
#   5. Call Tool (浮點)  — POST /mcp，呼叫 add(3.14,2.86)，預期回傳 6
#   6. Error Case       — POST /mcp，呼叫不存在的 method，預期回傳 JSON-RPC error

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "=========================================="
echo "HTTP MCP Server 測試"
echo "目標: $BASE_URL"
echo "=========================================="
echo ""

# 測試 1: Health Check — 確認 Server 存活
echo "=== 1. Health Check ==="
curl -s "$BASE_URL/health"
echo -e "\n"

# 測試 2: Initialize — 建立 JSON-RPC session，預期回傳 serverInfo
echo "=== 2. Initialize ==="
curl -s -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {}
  }' | jq .
echo ""

# 測試 3: List Tools — 列舉已註冊的 Tool，預期回傳 add
echo "=== 3. List Tools ==="
curl -s -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }' | jq .
echo ""

# 測試 4: Call Tool — 整數加法，預期 "10 + 20 = 30"
echo "=== 4. Call Tool (add 10 + 20) ==="
curl -s -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "add",
      "arguments": {"a": 10, "b": 20}
    }
  }' | jq .
echo ""

# 測試 5: Call Tool — 浮點數加法，預期 "3.14 + 2.86 = 6"
echo "=== 5. Call Tool (add 3.14 + 2.86) ==="
curl -s -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "add",
      "arguments": {"a": 3.14, "b": 2.86}
    }
  }' | jq .
echo ""

# 測試 6: Error Case — 呼叫不存在的 method，預期回傳 error code -32601
echo "=== 6. Error Case (unknown method) ==="
curl -s -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 5,
    "method": "unknown/method",
    "params": {}
  }' | jq .
echo ""

echo "=========================================="
echo "測試完成"
echo "=========================================="
