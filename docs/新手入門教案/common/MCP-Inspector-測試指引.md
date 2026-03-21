# MCP Inspector 測試指引

> 各單元共用的 MCP Inspector 操作流程。
> 單元 MD 中以連結指向本檔，不重複撰寫操作步驟。
> 完整 Inspector 安裝與進階設定另見 [docs/setup/mcp-inspector.md](../../setup/mcp-inspector.md)。

---

## Stdio 模式（U01-U06、U09 適用）

### 啟動方式

```bash
npx @modelcontextprotocol/inspector go run ./examples/<UNIT_DIR>/main.go
```

瀏覽器會自動開啟 Inspector 介面。

### 操作流程

1. 左側列出 Server 提供的 capabilities（Tools / Resources / Prompts）
2. 點擊對應頁籤，選擇要測試的項目
3. 填入參數並執行
4. 確認回傳結果符合預期

### 各單元指令速查

| 單元 | 指令 | 測試重點 |
|------|------|---------|
| U01 | `npx @modelcontextprotocol/inspector go run ./examples/U01-hello/main.go` | Server 資訊可見 |
| U03 | `npx @modelcontextprotocol/inspector go run ./examples/U03-tools/main.go` | Tools：`add`、`greet` |
| U04 | `npx @modelcontextprotocol/inspector go run ./examples/U04-resources/main.go` | Resources：readme、time、config |
| U05 | `npx @modelcontextprotocol/inspector go run ./examples/U05-prompts/main.go` | Prompts：code-review、summarize |
| U06 | `npx @modelcontextprotocol/inspector go run ./examples/U06-combined/main.go` | Tools + Resources + Prompts 三合一 |
| U09 | `npx @modelcontextprotocol/inspector go run ./examples/U09-weather-mvp/main.go` | Tools：get_forecast、get_alerts |

---

## HTTP 模式（U07 起適用）

### 啟動方式

需要兩個終端：

```bash
# 終端 1：啟動 Server
go run ./examples/<UNIT_DIR>/

# 終端 2：啟動 Inspector（不帶參數）
npx @modelcontextprotocol/inspector
```

### 操作流程

1. 在 Inspector 介面右上角點擊連線設定
2. 選擇 **HTTP** 模式（非預設的 Stdio）
3. 輸入 Server URL（如 `http://localhost:8080/mcp`）
4. 點擊連線
5. 測試各功能

---

## 需要認證的 Server（U08 適用）

MCP Inspector 目前不直接支援認證 Header，可用以下方式：

### 方式 1：暫時關閉認證

```bash
AUTH_MODE=none go run ./examples/U08-auth/

# 另一終端
npx @modelcontextprotocol/inspector
# 在 Inspector 中連線 http://localhost:8080/mcp
```

### 方式 2：使用 curl 測試完整認證流程

```bash
# 啟動認證模式
AUTH_MODE=bearer go run ./examples/U08-auth/

# 手動測試 MCP 方法
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token-123" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{...}}'
```

---

## 單元內引用寫法

後續單元（U03 起）引用本指引的標準寫法：

```markdown
### [必做] 4.X 使用 MCP Inspector 測試

> 操作方式見 [MCP Inspector 測試指引](../common/MCP-Inspector-測試指引.md)

​```bash
npx @modelcontextprotocol/inspector go run ./examples/U0X-xxx/main.go
​```

本單元測試重點：
- `Tools` 頁籤中的 `add` / `greet`
```

---

## 常見問題

| 現象 | 原因 | 解法 |
|------|------|------|
| Inspector 無回應 | Server 未啟動或已 crash | 確認終端有 Server 輸出 |
| 連線失敗（HTTP 模式） | URL 錯誤或 CORS 問題 | 確認 port 與路徑；參考 U07 踩坑 2 |
| 看不到 Tool/Resource | capabilities 未正確設定 | 參考 U06 踩坑 1 |
| `npx` 找不到套件 | Node.js 未安裝或版本過舊 | 需 Node.js 18+，見 U00 課前準備 |
