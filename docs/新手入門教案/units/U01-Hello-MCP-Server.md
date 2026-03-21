# U01｜Hello MCP Server

> 建立並執行第一個 MCP Server，設定 Claude Desktop / Claude Code 連接。
>
> 預估時數：90 min
> 前置依賴：U00

---

## ① 為什麼先教這個？

MCP Server 是與 AI 助手溝通的橋樑。在學習任何 capability（Tool、Resource、Prompt）之前，必須先能建立一個「會呼吸」的 Server——能啟動、能回應基本請求、能連接 Client。

這就像學習 Web 開發時，先跑通 "Hello World" 頁面，確認環境正確，才有信心繼續學習路由、資料庫等進階功能。

---

## ② 對應素材

| 類型 | 路徑/連結 | 說明 |
|------|----------|------|
| 範例程式碼 | `examples/U01-hello/main.go` | 最小 MCP Server |
| 官方文件 | https://modelcontextprotocol.io/docs/develop/build-server | 建置指南 |
| Go SDK | https://github.com/modelcontextprotocol/go-sdk | SDK 原始碼 |

---

## ③ 核心觀念

### 1. MCP 是什麼？

**Model Context Protocol (MCP)** 是一個開放標準，讓 AI 助手（如 Claude）能夠：
- 呼叫外部工具（Tools）
- 讀取資料來源（Resources）
- 使用預設提示詞（Prompts）

```
┌─────────────┐         JSON-RPC          ┌─────────────┐
│   Claude    │ ◄─────────────────────────► │ MCP Server  │
│  (Client)   │        over Stdio/HTTP     │  (你開發的)  │
└─────────────┘                            └─────────────┘
```

### 2. Server vs Client

| 角色 | 說明 | 範例 |
|------|------|------|
| **Client** | 呼叫 MCP 功能的一方 | Claude Desktop、Claude Code、MCP Inspector |
| **Server** | 提供 MCP 功能的一方 | 你開發的天氣服務、資料庫查詢服務 |

### 3. 三種 Capabilities

| Capability | 用途 | 類比 |
|------------|------|------|
| **Tools** | 讓 AI 執行動作（函式呼叫） | 函式 / API endpoint |
| **Resources** | 讓 AI 讀取資料 | 檔案 / 資料庫 |
| **Prompts** | 預設的提示詞模板 | 範本 / 快捷鍵 |

### 4. Transport 方式

| Transport | 說明 | 適用場景 |
|-----------|------|---------|
| **Stdio** | 透過標準輸入/輸出通訊 | 本地開發、Claude Desktop |
| **HTTP** | 透過 HTTP 請求通訊 | 遠端部署、多 Client |

### 5. JSON-RPC 2.0

MCP 使用 JSON-RPC 2.0 作為通訊協定：

```json
// 請求
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": { ... }
}

// 回應
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": { ... }
}
```

### 6. Implementation 結構

```go
server := mcp.NewServer(&mcp.Implementation{
    Name:    "my-server",    // Server 名稱
    Version: "1.0.0",        // 版本號
}, nil)
```

### 7. Logging 規則（重要！）

**Stdio 模式下，stdout 用於 JSON-RPC 通訊，禁止使用 `fmt.Println()`！**

| 輸出方式 | Stdio 模式 | HTTP 模式 |
|----------|-----------|-----------|
| `fmt.Println()` | ❌ 會破壞通訊 | ✅ 可用 |
| `log.Println()` | ✅ 輸出到 stderr | ✅ 可用 |
| `fmt.Fprintln(os.Stderr, ...)` | ✅ 明確輸出到 stderr | ✅ 可用 |

---

## ④ 動手做

### [必做] 4.1 建立專案結構

```bash
cd ~/mcp-test  # 或你的專案目錄
mkdir -p examples/U01-hello
cd examples/U01-hello
```

### [必做] 4.2 撰寫最小 MCP Server

建立 `main.go`：

```go
// Package main 實作最小可執行的 MCP Server
// 這是 U01 Hello MCP Server 的範例程式碼
package main

import (
	// "context" ... (略)（標準庫：context, log, os）
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// 設定 log 輸出到 stderr
	// 重要：Stdio 模式下，stdout 專門用於 JSON-RPC 通訊，禁止使用 fmt.Println()
	log.SetOutput(os.Stderr)
	log.Println("Starting hello-mcp server...")

	// 建立 MCP Server
	// Implementation 包含 Server 的基本資訊：名稱和版本
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "hello-mcp",
		Version: "1.0.0",
	}, nil)

	log.Println("Server created, waiting for connections...")

	// 使用 Stdio Transport 執行
	// Server 會透過標準輸入/輸出與 Client 通訊
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
```

### [必做] 4.3 編譯並測試

```bash
# 回到專案根目錄
cd ~/mcp-test

# 下載依賴
go mod tidy

# 編譯（確認無錯誤）
go build -o hello-mcp ./examples/U01-hello/

# 執行（會等待輸入，這是正常的）
./hello-mcp

# 按 Ctrl+C 結束
```

### [必做] 4.4 使用 MCP Inspector 測試

```bash
# 複製這行指令直接測試
npx @modelcontextprotocol/inspector go run ./examples/U01-hello/main.go
```

瀏覽器會自動開啟 Inspector 介面，可檢視 Server 資訊。

### [延伸] 4.5 設定 Claude Desktop

編輯設定檔（macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`）：

```json
{
  "mcpServers": {
    "hello-mcp": {
      "command": "/Users/YOUR_USERNAME/mcp-test/hello-mcp"
    }
  }
}
```

**驗證**：重啟 Claude Desktop → 點擊 "+" → 確認看到 "hello-mcp"

### [延伸] 4.6 設定 Claude Code

```bash
# 新增 MCP Server
claude mcp add hello-mcp -- /Users/YOUR_USERNAME/mcp-test/hello-mcp

# 驗證
claude mcp list  # 應顯示 hello-mcp
```

---

## ⑤ 踩坑提示

### 踩坑 1：Claude Desktop 找不到 Server

**現象**
Claude Desktop 的連接器清單中沒有出現你的 Server。

**原因**
- config.json 格式錯誤
- 路徑不是絕對路徑
- 執行檔不存在或無執行權限

**解法**
```bash
# 1. 確認執行檔存在且有權限
ls -la /path/to/your/hello-mcp
chmod +x /path/to/your/hello-mcp

# 2. 重啟 Claude Desktop
```

---

### 踩坑 2：Server 啟動後立即結束

**現象**
執行 `./hello-mcp` 後，程式立即結束，沒有等待輸入。

**原因**
可能是 `server.Run()` 回傳了錯誤。

**解法**
```go
// 確認有正確處理錯誤
if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
    log.Fatalf("Server error: %v", err)
}
```

---

### 踩坑 3：fmt.Println 導致連線失敗

**現象**
Server 有輸出，但 Claude 無法連線，或連線後立即斷開。

**原因**
Stdio 模式下，stdout 專門用於 JSON-RPC 通訊。

**解法**
```go
// ❌ 錯誤
fmt.Println("Debug message")

// ✅ 正確
log.Println("Debug message")  // 輸出到 stderr
```

---

## ⑥ 驗收標準（DoD）

| 項目 | 驗證方式 | 預期結果 |
|------|---------|---------|
| 編譯成功 | `go build -o hello-mcp ./examples/U01-hello/` | 無錯誤，產生執行檔 |
| 執行不報錯 | `./hello-mcp` | 程式持續執行，等待輸入 |
| Client 可見 | Claude Desktop "+" 或 `claude mcp list` | 看到 "hello-mcp" |

**自我檢核清單**
- [ ] `go build` 編譯成功
- [ ] 執行檔可正常啟動（不立即結束）
- [ ] Claude Desktop 或 Claude Code 能看到 "hello-mcp"
- [ ] 沒有使用 `fmt.Println()`（僅使用 `log.Println()`）

---

## 下一步

Hello Server 就緒後，前往 [U02｜MCP Inspector](U02-MCP-Inspector.md) 學習使用官方測試工具。
