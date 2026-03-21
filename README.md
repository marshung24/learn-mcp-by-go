# Learn MCP by Go

> **專案狀態：實作中**
> 本專案提供可執行的 MCP Server 範例程式碼，配合新手入門教案使用。

使用 Golang 建立 MCP (Model Context Protocol) Server 的教學範本專案。

從最小化的 Hello World 範例開始，逐步學習 Resources、Tools、Prompts 三種 capabilities，支援 Stdio 與 HTTP 兩種傳輸模式，並包含 Server 認證機制。

---

## 說明

本專案是一個 **MCP Server 開發教學範本**，採用漸進式學習路徑：

1. **Hello World** - 最小可執行的 MCP Server
2. **MCP Inspector** - 學習主要測試與除錯工具
3. **Tools** - 學習工具（函式呼叫）功能
4. **Resources** - 學習資源（檔案類資料）功能
5. **Prompts** - 學習提示詞模板功能
6. **Combined** - 整合三種 capabilities
7. **HTTP Transport** - 從 Stdio 切換至 HTTP 模式
8. **Server 認證** - Basic Auth 與 Bearer Token
9. **MVP 整合** - 完整的天氣查詢服務範例

**測試策略**：開發期以 MCP Inspector 為主要測試工具，Claude Desktop / Claude Code 僅用於整合驗證。

### 關於技術選型

本專案採用 **原生 Go + MCP 官方 SDK** 進行教學展示，目的是讓學員深入理解 MCP 協定的運作原理與核心概念。

在實務專案中，若需考量程式碼的可維護性、擴展性及團隊協作效率，可評估採用以下方案：

| 方案 | 適用場景 | 說明 |
|------|---------|------|
| 原生 Go + 官方 SDK | 學習、原型驗證、輕量級服務 | 本專案採用，直接掌握底層機制 |
| Web 框架 + 官方 SDK | 生產環境、中大型專案 | 結合 [Gin](https://github.com/gin-gonic/gin)、[Echo](https://github.com/labstack/echo) 等框架，提供更好的路由管理、中介層支援與錯誤處理 |

> 建議先透過本教案掌握 MCP 核心概念後，再根據專案需求選擇適合的實作方式。

---

## 依賴

| 工具/套件 | 版本需求 | 說明 |
|----------|---------|------|
| Go | 1.25+ | Go 程式語言 |
| Node.js | 18+ | MCP Inspector 執行環境 |
| MCP Go SDK | v1.4.1+ | `github.com/modelcontextprotocol/go-sdk` |
| Claude Desktop | 最新版 | 整合驗證用（可選） |

---

## 安裝

```bash
# 1. 複製專案
git clone https://github.com/marshung24/learn-mcp-by-go.git
cd learn-mcp-by-go

# 2. 確認 Go 版本
go version  # 需 1.25+

# 3. 下載依賴
go mod tidy
```

---

## 用法

### Stdio 模式（本地執行）

```bash
# 執行特定範例
go run examples/U01-hello/main.go
```

### 使用 MCP Inspector 測試（推薦）

```bash
# 啟動 Inspector 測試 Server
npx @modelcontextprotocol/inspector go run examples/U01-hello/main.go
```

開啟瀏覽器 `http://localhost:5173` 即可互動測試 Tools、Resources、Prompts。

### HTTP 模式

```bash
# 啟動 HTTP Server
go run examples/U07-http-transport/main.go

# 另開終端測試
curl http://localhost:8080/mcp

# 或用 Inspector 連接 HTTP Server
npx @modelcontextprotocol/inspector --server-url http://localhost:8080/mcp
```

### 與 Claude Desktop 整合（簡易設定）

編輯 `~/Library/Application Support/Claude/claude_desktop_config.json`：

```json
{
  "mcpServers": {
    "learn-mcp": {
      "command": "go",
      "args": ["run", "/path/to/examples/U01-hello/main.go"]
    }
  }
}
```

> 完整功能測試請使用 MCP Inspector。

### 與 Claude Code 整合（簡易設定）

```bash
claude mcp add learn-mcp -- go run ./examples/U01-hello/main.go
```

> 完整功能測試請使用 MCP Inspector。

---

## 範例

| 範例 | 說明 | 對應單元 | 路徑 |
|-----|------|---------|-----|
| U01-hello | 最小 MCP Server | U01 | `examples/U01-hello/` |
| U03-tools | 工具功能實作 | U03 | `examples/U03-tools/` |
| U04-resources | 資源功能實作 | U04 | `examples/U04-resources/` |
| U05-prompts | 提示詞模板 | U05 | `examples/U05-prompts/` |
| U06-combined | 整合三種 capabilities | U06 | `examples/U06-combined/` |
| U07-http-transport | HTTP 傳輸模式 | U07 | `examples/U07-http-transport/` |
| U08-auth | Server 認證 | U08 | `examples/U08-auth/` |
| U09-weather-mvp | 天氣查詢 MVP | U09 | `examples/U09-weather-mvp/` |

---

## 限制

- **僅支援 Go 1.25+**：需使用較新的 Go 版本
- **MCP 規範版本**：基於 2025-11-25 規範（待驗證），部分實驗性功能可能變動
- **平台限制**：Claude Desktop 目前不支援 Linux
- **HTTP 模式**：目前 Go SDK 原生支援有限，需自行包裝（待驗證）

---

## 參考資料

| 資源 | 說明 |
|-----|------|
| [MCP 官方文件](https://modelcontextprotocol.io/docs) | Model Context Protocol 規範與教學 |
| [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) | 官方測試與除錯工具 |
| [Go SDK GitHub](https://github.com/modelcontextprotocol/go-sdk) | 官方 Go SDK 原始碼 |
| [Build Server 教學](https://modelcontextprotocol.io/docs/develop/build-server) | 官方 Server 建置指南 |
| [新手入門教案](docs/新手入門教案/README.md) | 本專案的完整學習路徑 |

---

## 其他

### 專案結構

```
learn-mcp-by-go/
├── README.md                 # 專案說明
├── go.mod                    # Go 模組定義
├── docs/                     # 文件
│   ├── setup/                # 工具安裝與設定指引
│   └── 新手入門教案/          # 教學文件
│       ├── README.md         # 教案總綱
│       ├── 新手入門教案綱要.md # 教案綱要
│       ├── common/           # 跨單元共用指引（結構示範檔、測試指引等）
│       └── units/            # 各單元教案（U00-U09）
└── examples/                 # 範例程式碼（Single Source of Truth）
    ├── U01-hello/            # U01 最小 MCP Server
    ├── U03-tools/            # U03 工具功能
    ├── U04-resources/        # U04 資源功能
    ├── U05-prompts/          # U05 提示詞模板
    ├── U06-combined/         # U06 整合範例
    ├── U07-http-transport/   # U07 HTTP 模式
    ├── U08-auth/             # U08 Server 認證
    └── U09-weather-mvp/      # U09 天氣查詢 MVP
```

### 貢獻指南

歡迎提交 Issue 或 Pull Request。

### 授權

MIT License
