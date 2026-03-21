# Claude Desktop 安裝設定

> Claude Desktop 是 Anthropic 官方的桌面應用程式，可作為 MCP Client 與你的 MCP Server 互動。

---

## 安裝方式

1. 前往 https://claude.ai/download
2. 下載對應作業系統的安裝檔
3. 執行安裝程式
4. 登入 Claude 帳號

---

## 設定檔位置

| 系統 | 路徑 |
|------|------|
| macOS | `~/Library/Application Support/Claude/claude_desktop_config.json` |
| Windows | `%AppData%\Claude\claude_desktop_config.json` |

### 確認目錄存在

**macOS**
```bash
ls -la ~/Library/Application\ Support/Claude/
# 如果目錄不存在，先開啟一次 Claude Desktop 讓它建立
```

**Windows (PowerShell)**
```powershell
ls $env:AppData\Claude\
```

---

## 設定 MCP Server

編輯 `claude_desktop_config.json`，加入 MCP Server 設定：

```json
{
  "mcpServers": {
    "my-server": {
      "command": "go",
      "args": ["run", "/path/to/your/main.go"]
    }
  }
}
```

### 設定範例：Hello Server

```json
{
  "mcpServers": {
    "hello-mcp": {
      "command": "go",
      "args": ["run", "/Users/username/mcp-test/examples/U01-hello/main.go"]
    }
  }
}
```

### 設定範例：編譯後的執行檔

```json
{
  "mcpServers": {
    "my-server": {
      "command": "/path/to/compiled/server"
    }
  }
}
```

### 設定範例：HTTP 模式

```json
{
  "mcpServers": {
    "http-server": {
      "command": "npx",
      "args": ["mcp-remote", "http://localhost:8080/mcp"]
    }
  }
}
```

---

## 套用設定

修改設定檔後，需要**完全關閉並重新啟動** Claude Desktop：

**macOS**
```bash
# 強制關閉 Claude Desktop
pkill -f "Claude"

# 重新開啟
open -a "Claude"
```

**Windows**
- 從系統匣（System Tray）右鍵點選 Claude 圖示 → 結束
- 重新啟動 Claude Desktop

---

## 驗證連線

1. 開啟 Claude Desktop
2. 開始新對話
3. 輸入測試訊息，例如：「請列出你可以使用的工具」
4. 如果 MCP Server 設定正確，Claude 會列出可用的 Tools

---

## 常見問題

### Server 無法連線

**可能原因**：
- 設定檔路徑錯誤
- JSON 格式有誤
- Server 程式無法執行

**除錯步驟**：
```bash
# 1. 驗證 JSON 格式
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json | jq .

# 2. 手動測試 Server 能否執行
go run /path/to/your/main.go

# 3. 檢查 Claude Desktop 日誌
# macOS: ~/Library/Logs/Claude/
```

### 修改設定後沒有生效

**原因**：Claude Desktop 需要完全重啟。

**解法**：確保從系統匣完全關閉，而非只是關閉視窗。

### Server 執行時錯誤

**排查方式**：
1. 先用 MCP Inspector 測試 Server 是否正常
2. 確認 Server 沒有使用 `fmt.Println()`（會破壞 Stdio 通訊）
3. 使用 `log.Println()` 輸出到 stderr

---

## 進階設定

### 環境變數

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

### 工作目錄

```json
{
  "mcpServers": {
    "my-server": {
      "command": "./server",
      "cwd": "/path/to/project"
    }
  }
}
```

---

## 參考資源

- [Claude Desktop 下載](https://claude.ai/download)
- [MCP Server 設定文件](https://modelcontextprotocol.io/docs/quickstart)
