# U07｜HTTP Transport

> 從 Stdio 切換到 HTTP 模式，支援遠端部署與多 Client 連線。
>
> 預估時數：90 min
> 前置依賴：U06

---

## ① 為什麼先教這個？

Stdio 模式適合本地開發和 Claude Desktop，但有限制：
- 只能在本機執行
- 一次只能服務一個 Client
- 無法遠端存取

HTTP 模式解決這些問題：
- 可部署到伺服器
- 支援多個 Client 同時連線
- 可透過網路存取

學會 HTTP Transport 後，就能將 MCP Server 部署到雲端，讓多個使用者共用。

---

## ② 對應素材

| 類型 | 路徑/連結 | 說明 |
|------|----------|------|
| 範例程式碼 | `examples/U07-http-transport/main.go` | HTTP 模式 Server |
| 測試腳本 | `examples/U07-http-transport/test.sh` | curl 測試指令 |
| MCP 規範 | https://modelcontextprotocol.io/docs | Transport 說明 |

---

## ③ 核心觀念

### 1. Transport 比較

| 特性 | Stdio | HTTP |
|------|-------|------|
| 通訊方式 | 標準輸入/輸出 | HTTP 請求 |
| 適用場景 | 本地開發、Claude Desktop | 遠端部署、Web 服務 |
| 多 Client | ❌ 不支援 | ✅ 支援 |
| 網路存取 | ❌ 僅本機 | ✅ 可遠端 |
| 狀態管理 | Process 生命週期 | 需自行管理 Session |

### 2. HTTP MCP 協定

MCP over HTTP 使用 JSON-RPC 2.0：

**請求**
```http
POST /mcp HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "add",
    "arguments": {"a": 1, "b": 2}
  }
}
```

**回應**
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [{"type": "text", "text": "1 + 2 = 3"}]
  }
}
```

### 3. Go SDK HTTP 支援現況

Go SDK 提供 `mcp.NewStreamableHTTPHandler()`，可直接建立 HTTP handler：

1. **使用 Streamable HTTP Transport**（本單元採用此方式，SDK 官方實作）
2. **自行包裝 HTTP Server**（適合需要深度客製的場景）
3. **使用第三方 HTTP adapter**

**生產環境建議**：如需成熟的 HTTP 功能（認證、限流、中介軟體等），可考慮使用 Web 框架搭配 MCP 適配器：
- [Gin](https://github.com/gin-gonic/gin) + [gin-mcp](https://github.com/metoro-io/gin-mcp)
- 其他 Go Web 框架 + 自行實作適配器

### 4. Session 管理

Stdio 模式下，每個 process 就是一個 session。HTTP 模式需要處理：
- Session 建立與銷毀
- 多 Client 隔離
- 狀態保持

### 5. CORS 設定

如果 Client 是瀏覽器，需要設定 CORS：

```go
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
```

---

## ④ 動手做

### [必做] 4.1 建立專案結構

```bash
cd ~/mcp-test
mkdir -p examples/U07-http-transport
```

### [必做] 4.2 使用官方 StreamableHTTPHandler

建立 `examples/U07-http-transport/main.go`：

```go
package main

import (
	// "context" ... (略)
	"net/http"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ==================== Tools ====================

// AddInput 定義 add Tool 的輸入參數
// 用於接收兩個浮點數進行加法運算
type AddInput struct {
	A float64 `json:"a" jsonschema:"第一個數字"`
	B float64 `json:"b" jsonschema:"第二個數字"`
}

// addHandler 處理加法運算
// addHandler ... (略)（見結構示範檔 §3）

// ==================== Main ====================

func main() {
	// 設定 log 輸出到 stderr
	// log.SetOutput(os.Stderr) ... (略)

	// 讀取環境變數設定 port
	// port := os.Getenv("PORT") ... (略)

	// 建立 MCP Server（見結構示範檔 §1）
	// server := mcp.NewServer( ... (略)

	// 註冊 Tool（見結構示範檔 §3）
	// mcp.AddTool(server, ... (略)

	// 使用官方的 StreamableHTTPHandler
	// 這個 handler 會自動處理 JSON-RPC、session 管理等
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		// 為每個請求返回 server 實例
		// 實際應用中可以根據 req (如 headers) 返回不同的 server
		return server
	}, nil)

	// 設定路由 - 將 /mcp 路徑綁定到 MCP handler
	http.Handle("/mcp", handler)

	// 首頁說明端點 - 顯示 Server 資訊與使用方式
	// http.HandleFunc("/", ... (略)

	// Health check 端點
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 輸出啟動資訊
	// log.Printf("HTTP MCP Server ... (略)

	// 啟動 HTTP Server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
```

**重點說明**：

1. **`mcp.NewStreamableHTTPHandler`**：SDK 提供的官方 HTTP handler
2. **自動處理**：JSON-RPC 解析、協定版本、session 管理都由 SDK 處理
3. **彈性**：回調函數可根據 HTTP request（如 headers）返回不同 server 實例

### [必做] 4.3 建立測試腳本

測試腳本 `examples/U07-http-transport/test.sh` 依序驗證 HTTP MCP Server 的核心端點：

| # | 測試案例 | 端點 | 預期結果 |
|---|---------|------|---------|
| 1 | Health Check | GET /health | 200 OK |
| 2 | Initialize | POST /mcp | 回傳 serverInfo |
| 3 | List Tools | POST /mcp | 回傳 add tool |
| 4 | Call Tool (整數) | POST /mcp | `10 + 20 = 30` |
| 5 | Call Tool (浮點) | POST /mcp | `3.14 + 2.86 = 6` |
| 6 | Error Case | POST /mcp | JSON-RPC error -32601 |

> 完整腳本見 `examples/U07-http-transport/test.sh`。

### [必做] 4.4 執行並測試

```bash
# 終端 1：啟動 Server
go run ./examples/U07-http-transport/

# 終端 2：執行測試
chmod +x examples/U07-http-transport/test.sh
./examples/U07-http-transport/test.sh
```

預期輸出：
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [{"type": "text", "text": "10 + 20 = 30"}]
  }
}
```

### [必做] 4.5 使用 MCP Inspector 測試（HTTP 模式）

```bash
# 先啟動 Server（終端 1）
go run ./examples/U07-http-transport/

# 啟動 Inspector（終端 2）
npx @modelcontextprotocol/inspector
```

在 Inspector 介面中：
1. 點擊右上角的連線設定
2. 選擇 **HTTP** 模式
3. 輸入 URL：`http://localhost:8080/mcp`
4. 點擊連線
5. 測試 `add` 工具

### [延伸] 4.6 使用環境變數設定

```go
func main() {
	// 讀取環境變數設定 port（已在 4.2 示範）
	// port := os.Getenv("PORT") ... (略)

	// 讀取 API Key，無設定時輸出警告
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("Warning: No API_KEY set, running without authentication")
	}

	// server := mcp.NewServer( ... (略)
}
```

### [延伸] 4.7 Docker 化

建立 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./examples/U07-http-transport/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

```bash
docker build -t mcp-http-server .
docker run -p 8080:8080 mcp-http-server
```

---

## ⑤ 踩坑提示

### 踩坑 1：連線被拒絕

**現象**
curl 回報 "Connection refused"。

**原因**
- Server 未啟動
- Port 被佔用
- 防火牆阻擋

**解法**
```bash
# 確認 Server 執行中
ps aux | grep http-demo

# 確認 Port 未被佔用
lsof -i :8080

# 換 Port
PORT=9090 go run ./examples/U07-http-transport/
```

---

### 踩坑 2：CORS 錯誤

**現象**
瀏覽器 Console 顯示 CORS 錯誤。

**原因**
未設定 CORS headers。

**解法**
```go
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 處理 preflight
if r.Method == "OPTIONS" {
    w.WriteHeader(http.StatusOK)
    return
}
```

---

### 踩坑 3：JSON 解析失敗

**現象**
Server 回傳 "Parse error"。

**原因**
請求 body 不是有效 JSON。

**解法**
```bash
# 確認 JSON 格式正確
echo '{"jsonrpc":"2.0","id":1,"method":"initialize"}' | jq .

# 使用正確的 Content-Type
curl -H "Content-Type: application/json" ...
```

---

## ⑥ 驗收標準（DoD）

| 項目 | 驗證方式 | 預期結果 |
|------|---------|---------|
| Server 啟動 | 執行後觀察 log | 顯示 "listening on :8080" |
| Health Check | `curl /health` | 回傳 "OK" |
| Initialize | POST /mcp | 回傳 serverInfo |
| tools/list | POST /mcp | 回傳 add tool |
| tools/call | POST /mcp | 回傳計算結果 |

**自我檢核清單**
- [ ] Server 在 8080 port 啟動
- [ ] curl 測試腳本全部通過
- [ ] CORS 設定正確
- [ ] 錯誤處理正確（錯誤的 method 回傳 -32601）

---

## 下一步

HTTP Transport 就緒後，前往 [U08｜Server 認證](U08-Server-Auth.md) 學習如何保護你的 API。
