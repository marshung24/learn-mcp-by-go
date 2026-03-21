# Claude Code 安裝設定

> Claude Code 是 Anthropic 官方的 CLI 工具，可在終端機中與 Claude 互動，並支援 MCP Server。

---

## 安裝方式

### macOS (Homebrew)

```bash
brew install claude-code
```

### npm（跨平台）

```bash
npm install -g @anthropic-ai/claude-code
```

### 驗證安裝

```bash
claude --version
```

---

## 設定檔位置

| 範圍 | 路徑 | 說明 |
|------|------|------|
| 專案層級 | `.mcp.json`（專案根目錄） | 僅對該專案生效 |
| 使用者層級 | `~/.claude/settings.json` | 對所有專案生效 |

---

## MCP Server 管理指令

### 查看目前的 MCP Server

```bash
claude mcp list
```

### 新增 MCP Server

```bash
# 基本語法
claude mcp add <name> -- <command> [args...]

# 範例：新增 Go Server
claude mcp add hello-server -- go run ./examples/U01-hello/main.go

# 範例：新增編譯後的執行檔
claude mcp add my-server -- /path/to/server
```

### 移除 MCP Server

```bash
claude mcp remove <name>
```

---

## 設定檔格式

### 專案層級 `.mcp.json`

在專案根目錄建立 `.mcp.json`：

```json
{
  "mcpServers": {
    "hello-mcp": {
      "command": "go",
      "args": ["run", "./examples/U01-hello/main.go"]
    },
    "tools-server": {
      "command": "go",
      "args": ["run", "./examples/U03-tools/main.go"]
    }
  }
}
```

### 使用者層級 `~/.claude/settings.json`

```json
{
  "mcpServers": {
    "global-server": {
      "command": "/usr/local/bin/my-mcp-server"
    }
  }
}
```

---

## 使用方式

### 啟動互動模式

```bash
# 在專案目錄中啟動（會自動載入 .mcp.json）
cd /path/to/project
claude

# 或指定工作目錄
claude --cwd /path/to/project
```

### 單次查詢

```bash
claude -m "請列出你可以使用的工具"
```

### 執行特定任務

```bash
claude -m "使用 add 工具計算 5 + 3"
```

---

## 環境變數設定

### 在設定檔中指定

```json
{
  "mcpServers": {
    "my-server": {
      "command": "go",
      "args": ["run", "./main.go"],
      "env": {
        "API_KEY": "your-api-key",
        "DEBUG": "true"
      }
    }
  }
}
```

### 使用系統環境變數

```bash
export MY_API_KEY="your-key"
claude
```

---

## 常見問題

### claude: command not found

**原因**：Claude Code 未加入 PATH。

**解法**：
```bash
# npm 全域安裝後，確認 npm bin 目錄在 PATH 中
export PATH=$PATH:$(npm prefix -g)/bin

# 或重新安裝
npm install -g @anthropic-ai/claude-code
```

### MCP Server 無法載入

**排查步驟**：
```bash
# 1. 確認 .mcp.json 存在且格式正確
cat .mcp.json | jq .

# 2. 確認 Server 可以單獨執行
go run ./main.go

# 3. 使用 MCP Inspector 測試
npx @modelcontextprotocol/inspector go run ./main.go
```

### 設定檔優先順序

當專案層級和使用者層級都有設定時：
1. 專案層級 `.mcp.json` 優先
2. 使用者層級 `~/.claude/settings.json` 作為備援

---

## 進階用法

### 指定 Server 執行

```bash
# 僅載入特定 Server
claude --mcp-server hello-mcp
```

### 除錯模式

```bash
# 顯示詳細日誌
claude --verbose

# 顯示 MCP 通訊紀錄
claude --debug-mcp
```

### 與 HTTP Server 連接

```bash
# 先啟動 HTTP Server
go run ./examples/U07-http-transport/main.go &

# 在 .mcp.json 中設定
# {
#   "mcpServers": {
#     "http-server": {
#       "command": "npx",
#       "args": ["mcp-remote", "http://localhost:8080/mcp"]
#     }
#   }
# }
```

---

## 參考資源

- [Claude Code 官方文件](https://docs.anthropic.com/en/docs/claude-code)
- [Claude Code GitHub](https://github.com/anthropics/claude-code)
