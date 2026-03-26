#!/bin/bash

# ============================================================
# MCP Server 認證測試腳本
# ============================================================
#
# 用法：
#   chmod +x test.sh
#   ./test.sh
#
# 測試前請先啟動 Server：
#   AUTH_MODE=none   go run .   # 無認證
#   AUTH_MODE=basic  go run .   # Basic Auth
#   AUTH_MODE=bearer go run .   # Bearer Token
# ============================================================

BASE_URL="http://localhost:8080"

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "MCP Server 認證測試"
echo "=========================================="
echo ""

# ===== 1. Health Check（不需認證）=====
echo -e "${YELLOW}=== 1. Health Check（不需認證）===${NC}"
HEALTH_RESULT=$(curl -s "$BASE_URL/health")
echo "Response: $HEALTH_RESULT"
if [ "$HEALTH_RESULT" = "OK" ]; then
    echo -e "${GREEN}✓ Health check passed${NC}"
else
    echo -e "${RED}✗ Health check failed${NC}"
fi
echo ""

# ===== 2. 無認證請求 =====
echo -e "${YELLOW}=== 2. 無認證請求 ===${NC}"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}')
echo "HTTP Status: $HTTP_CODE"
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Request succeeded (AUTH_MODE=none)${NC}"
elif [ "$HTTP_CODE" = "401" ]; then
    echo -e "${YELLOW}→ Authentication required (expected if AUTH_MODE is set)${NC}"
else
    echo -e "${RED}✗ Unexpected status code${NC}"
fi
echo ""

# ===== 3. Basic Auth 測試（錯誤密碼）=====
echo -e "${YELLOW}=== 3. Basic Auth（錯誤密碼）===${NC}"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -u admin:wrongpassword -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}')
echo "HTTP Status: $HTTP_CODE"
if [ "$HTTP_CODE" = "401" ]; then
    echo -e "${GREEN}✓ Correctly rejected (expected 401)${NC}"
elif [ "$HTTP_CODE" = "200" ]; then
    echo -e "${YELLOW}→ Request succeeded (AUTH_MODE=none or different auth)${NC}"
fi
echo ""

# ===== 4. Basic Auth 測試（正確帳密: admin）=====
echo -e "${YELLOW}=== 4. Basic Auth（正確帳密: admin:secret）===${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/mcp" \
  -u admin:secret \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}')
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -u admin:secret -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}')
echo "HTTP Status: $HTTP_CODE"
echo "Response (truncated): ${RESPONSE:0:200}..."
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Authentication successful (admin)${NC}"
fi
echo ""

# ===== 4b. Basic Auth 測試（正確帳密: guest）=====
echo -e "${YELLOW}=== 4b. Basic Auth（正確帳密: guest:guest123）===${NC}"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -u guest:guest123 -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}')
echo "HTTP Status: $HTTP_CODE"
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Authentication successful (guest)${NC}"
elif [ "$HTTP_CODE" = "401" ]; then
    echo -e "${YELLOW}→ Rejected (AUTH_MODE=none or different auth)${NC}"
fi
echo ""

# ===== 5. Bearer Token 測試（無效 Token）=====
echo -e "${YELLOW}=== 5. Bearer Token（無效 Token）===${NC}"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid-token-xyz" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}')
echo "HTTP Status: $HTTP_CODE"
if [ "$HTTP_CODE" = "401" ]; then
    echo -e "${GREEN}✓ Correctly rejected (expected 401)${NC}"
elif [ "$HTTP_CODE" = "200" ]; then
    echo -e "${YELLOW}→ Request succeeded (AUTH_MODE=none or different auth)${NC}"
fi
echo ""

# ===== 6. Bearer Token 測試（有效 Token）=====
echo -e "${YELLOW}=== 6. Bearer Token（有效 Token: test-token-123）===${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token-123" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}')
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token-123" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}')
echo "HTTP Status: $HTTP_CODE"
echo "Response (truncated): ${RESPONSE:0:200}..."
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Authentication successful${NC}"
fi
echo ""

# ===== 7. 呼叫 Tool（使用 Bearer Token）=====
echo -e "${YELLOW}=== 7. 呼叫 Tool 'add'（使用 Bearer Token）===${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token-123" \
  -d '{
    "jsonrpc":"2.0",
    "id":2,
    "method":"tools/call",
    "params":{"name":"add","arguments":{"a":10,"b":20}}
  }')
echo "Response:"
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"
echo ""

# ===== 8. 呼叫 Tool（使用 Basic Auth）=====
echo -e "${YELLOW}=== 8. 呼叫 Tool 'add'（使用 Basic Auth）===${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/mcp" \
  -u admin:secret \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "id":3,
    "method":"tools/call",
    "params":{"name":"add","arguments":{"a":100,"b":200}}
  }')
echo "Response:"
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"
echo ""

# ===== 9. 列出 Tools =====
echo -e "${YELLOW}=== 9. 列出 Tools（使用 Bearer Token）===${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/mcp" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token-123" \
  -d '{"jsonrpc":"2.0","id":4,"method":"tools/list","params":{}}')
echo "Response:"
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"
echo ""

echo "=========================================="
echo "測試完成"
echo "=========================================="
echo ""
echo "測試指南："
echo "  AUTH_MODE=none   → 所有請求都成功（200）"
echo "  AUTH_MODE=basic  → 需要 -u admin:secret 或 -u guest:guest123"
echo "  AUTH_MODE=bearer → 需要 Authorization: Bearer test-token-123"
echo ""
echo "範例："
echo "  # 啟動無認證 Server"
echo "  AUTH_MODE=none go run ."
echo ""
echo "  # 啟動 Basic Auth Server（支援多組帳密）"
echo "  AUTH_MODE=basic go run ."
echo ""
echo "  # 啟動 Bearer Token Server"
echo "  AUTH_MODE=bearer go run ."
