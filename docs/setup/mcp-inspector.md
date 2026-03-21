# MCP Inspector 安裝設定

> MCP Inspector 是官方提供的互動式測試工具，用於開發與除錯 MCP Server。

---

## 前置需求

MCP Inspector 需要 Node.js 18+ 環境。

### 安裝 Node.js

**macOS (Homebrew)**
```bash
brew install node
```

**其他平台**
前往 https://nodejs.org/ 下載安裝

**驗證安裝**
```bash
node --version  # 需顯示 v18.x 或更高
npm --version
```

---

## 使用方式

MCP Inspector 透過 `npx` 執行，無需額外安裝。

### Stdio 模式（直接啟動 Server）

```bash
# 基本用法
npx @modelcontextprotocol/inspector go run ./main.go

# 範例：啟動 Hello Server
npx @modelcontextprotocol/inspector go run ./examples/U01-hello/main.go
```

### HTTP 模式（連接已運行的 Server）

```bash
# 先啟動 HTTP Server（另開終端）
go run ./examples/U07-http-transport/main.go

# 使用 Inspector 連接
npx @modelcontextprotocol/inspector --url http://localhost:8080/mcp
```

---

## 介面說明

執行後開啟瀏覽器前往 `http://localhost:5173`

```
┌─────────────────────────────────────────────────────────┐
│                    MCP Inspector                        │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │   Server    │    │   Tools     │    │  Resources  │ │
│  │   連線狀態   │    │   測試面板   │    │   瀏覽面板   │ │
│  └─────────────┘    └─────────────┘    └─────────────┘ │
│                           │                             │
│                    JSON-RPC 通訊                        │
│                           │                             │
│                    ┌─────────────┐                      │
│                    │ MCP Server  │                      │
│                    │ (你的程式)   │                      │
│                    └─────────────┘                      │
└─────────────────────────────────────────────────────────┘
```

### 主要功能

| 功能 | 說明 |
|------|------|
| **Server Info** | 檢視 Server 基本資訊（名稱、版本、capabilities） |
| **Tools** | 列出並測試所有 Tools，可輸入參數並執行 |
| **Resources** | 瀏覽並讀取所有 Resources |
| **Prompts** | 檢視並測試 Prompt 模板 |
| **Logs** | 查看通訊日誌，除錯用 |

---

## Inspector vs Claude Desktop

| 特性 | MCP Inspector | Claude Desktop |
|------|--------------|----------------|
| 用途 | 開發與除錯 | 實際使用 |
| 啟動速度 | 快（直接啟動 Server） | 慢（需重啟 App） |
| 請求可見度 | 完整 JSON 顯示 | 僅看到結果 |
| 互動方式 | 手動發送請求 | AI 自動呼叫 |
| 適用場景 | 開發中測試 | 整合測試、正式使用 |

---

## 開發流程建議

```
1. 撰寫/修改程式碼
       │
       ▼
2. 用 Inspector 快速測試
       │
       ├─ 失敗 → 查看 Logs → 修正 → 回到 1
       │
       ▼ 成功
3. 用 Claude Desktop/Code 整合測試
       │
       ▼
4. 完成
```

---

## 常見問題

### Inspector 無法啟動

**現象**
```bash
npm ERR! could not determine executable to run
```

**原因**：Node.js 未安裝或版本過舊。

**解法**
```bash
# 確認 Node.js 版本
node --version

# 需要 18+，如版本過舊：
brew upgrade node
```

### Server 連線失敗

**現象**：Inspector 介面顯示「Disconnected」或紅色狀態。

**原因**：
- Server 程式有錯誤，啟動失敗
- Server 使用了 `fmt.Println()`，破壞了 Stdio 通訊

**解法**
```bash
# 1. 先單獨測試 Server 能否編譯
go build ./main.go

# 2. 檢查是否有 fmt.Println（應改用 log.Println）
grep -r "fmt.Println" ./
```

### Tool 執行無回應

**原因**：
- Tool handler 內有無限迴圈
- Tool handler 沒有正確回傳結果

**解法**：確認 handler 有正確 return

### 參數格式錯誤

**現象**：`Error: Invalid params`

**解法**：
1. 在「Tools」頁面查看 Tool 的 Input Schema
2. 確認參數名稱與型別正確
3. JSON 字串需要雙引號：`{"name": "test"}` 而非 `{name: "test"}`

---

## 參考資源

- [MCP Inspector 官方文件](https://modelcontextprotocol.io/docs/tools/inspector)
- [Inspector GitHub](https://github.com/modelcontextprotocol/inspector)
