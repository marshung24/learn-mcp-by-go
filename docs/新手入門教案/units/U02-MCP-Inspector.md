# U02｜MCP Inspector

> 學習使用官方 MCP Inspector 工具測試與除錯 MCP Server。
>
> 預估時數：45 min
> 前置依賴：U01

---

## ① 為什麼先教這個？

在開發 MCP Server 時，如果每次都要透過 Claude Desktop 或 Claude Code 來測試，效率會很低。MCP Inspector 是官方提供的互動式測試工具，可以：

- **即時測試**：不需重啟 Claude，直接與 Server 互動
- **視覺化介面**：清楚看到請求與回應的 JSON 結構
- **快速迭代**：修改程式碼後立即測試，加速開發流程
- **除錯利器**：當 Claude 連不上或行為異常時，用 Inspector 隔離問題

這就像 Web 開發時使用 Postman 測試 API 一樣重要。

---

## ② 對應素材

| 類型 | 路徑/連結 | 說明 |
|------|----------|------|
| 官方文件 | https://modelcontextprotocol.io/docs/tools/inspector | MCP Inspector 文件 |
| GitHub | https://github.com/modelcontextprotocol/inspector | Inspector 原始碼 |
| 範例 Server | `examples/U01-hello/main.go` | 用於測試的 Hello Server |

---

## ③ 核心觀念

### 1. MCP Inspector 是什麼？

MCP Inspector 是一個基於 Web 的互動式開發工具，用於測試和除錯 MCP Server。

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

### 2. Inspector vs Claude Desktop

| 特性 | MCP Inspector | Claude Desktop |
|------|--------------|----------------|
| 用途 | 開發與除錯 | 實際使用 |
| 啟動速度 | 快（直接啟動 Server） | 慢（需重啟 App） |
| 請求可見度 | 完整 JSON 顯示 | 僅看到結果 |
| 互動方式 | 手動發送請求 | AI 自動呼叫 |
| 適用場景 | 開發中測試 | 整合測試、正式使用 |

### 3. Inspector 主要功能

| 功能 | 說明 |
|------|------|
| **Server Info** | 檢視 Server 基本資訊（名稱、版本、capabilities） |
| **Tools** | 列出並測試所有 Tools，可輸入參數並執行 |
| **Resources** | 瀏覽並讀取所有 Resources |
| **Prompts** | 檢視並測試 Prompt 模板 |
| **Logs** | 查看通訊日誌，除錯用 |

### 4. 兩種啟動模式

| 模式 | 指令 | 說明 |
|------|------|------|
| **Stdio** | `npx @anthropic-ai/mcp-inspector go run ./main.go` | 直接啟動 Server |
| **HTTP** | `npx @anthropic-ai/mcp-inspector --url http://localhost:8080` | 連接已運行的 HTTP Server |

---

## ④ 動手做

### [必做] 4.1 安裝 MCP Inspector

MCP Inspector 透過 npx 執行，需要 Node.js 環境：

```bash
# 確認 Node.js 已安裝
node --version  # 需 18+

# 如未安裝，使用 Homebrew (macOS)
brew install node

# 或下載安裝
# https://nodejs.org/
```

### [必做] 4.2 啟動 Inspector 測試 Hello Server

```bash
# 進入專案目錄
cd ~/mcp-test

# 使用 Inspector 啟動 Hello Server
npx @modelcontextprotocol/inspector go run ./examples/U01-hello/main.go
```

執行後會看到類似輸出：

```
MCP Inspector is running at http://localhost:5173
```

### [必做] 4.3 瀏覽 Inspector 介面

1. 開啟瀏覽器，前往 `http://localhost:5173`
2. 左側面板顯示連線狀態，應為綠色「Connected」
3. 點擊「Server Info」查看 Server 資訊：
   - Name: hello-mcp
   - Version: 1.0.0
   - Capabilities: （目前為空，因為還沒加功能）

### [必做] 4.4 測試 Tool（需先完成 U02）

如果已完成 U02 Tool Capability，可以測試 Tool：

1. 點擊左側「Tools」標籤
2. 選擇要測試的 Tool（例如 `add`）
3. 在參數欄位輸入 JSON：
   ```json
   {
     "a": 5,
     "b": 3
   }
   ```
4. 點擊「Execute」
5. 查看下方的回應結果

### [必做] 4.5 查看通訊日誌

1. 點擊「Logs」標籤
2. 可以看到所有 JSON-RPC 請求與回應
3. 這對於除錯非常有用：
   - 檢查請求格式是否正確
   - 確認回應內容
   - 找出錯誤訊息

### [延伸] 4.6 測試 HTTP 模式

如果已完成 U06 HTTP Transport：

```bash
# 先啟動 HTTP Server（另開終端）
go run ./examples/U07-http-transport/main.go

# 使用 Inspector 連接 HTTP Server
npx @modelcontextprotocol/inspector --url http://localhost:8080/mcp
```

### [延伸] 4.7 搭配 --stdio 參數除錯

當 Server 有問題時，可以直接用 Inspector 看錯誤：

```bash
# 啟動時會顯示 Server 的 stderr 輸出
npx @modelcontextprotocol/inspector go run ./main.go 2>&1
```

---

## ⑤ 踩坑提示

### 踩坑 1：Inspector 無法啟動

**現象**
```bash
$ npx @modelcontextprotocol/inspector go run ./main.go
npm ERR! could not determine executable to run
```

**原因**
Node.js 未安裝或版本過舊。

**解法**
```bash
# 確認 Node.js 版本
node --version

# 需要 18+，如版本過舊：
brew upgrade node
# 或重新安裝
```

---

### 踩坑 2：Server 連線失敗

**現象**
Inspector 介面顯示「Disconnected」或紅色狀態。

**原因**
- Server 程式有錯誤，啟動失敗
- Server 使用了 `fmt.Println()`，破壞了 Stdio 通訊

**解法**
```bash
# 1. 先單獨測試 Server 能否編譯
go build ./examples/U01-hello/

# 2. 檢查是否有 fmt.Println
grep -r "fmt.Println" ./examples/U01-hello/

# 3. 改用 log.Println 輸出到 stderr
```

---

### 踩坑 3：Tool 執行無回應

**現象**
點擊 Execute 後，介面卡住或超時。

**原因**
- Tool handler 內有無限迴圈
- Tool handler 沒有正確回傳結果

**解法**
```go
// 確認 Tool handler 有正確回傳
func handleAdd(ctx context.Context, params AddParams) (*mcp.CallToolResult, error) {
    result := params.A + params.B
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            mcp.TextContent{
                Text: fmt.Sprintf("%d", result),
            },
        },
    }, nil  // 別忘了 return！
}
```

---

### 踩坑 4：參數格式錯誤

**現象**
```
Error: Invalid params
```

**原因**
輸入的 JSON 參數格式與 Tool 定義不符。

**解法**
1. 在「Tools」頁面查看 Tool 的 Input Schema
2. 確認參數名稱與型別正確
3. JSON 字串需要雙引號：`{"name": "test"}` 而非 `{name: "test"}`

---

## ⑥ 驗收標準（DoD）

| 項目 | 驗證方式 | 預期結果 |
|------|---------|---------|
| Inspector 啟動 | 執行 npx 指令 | 瀏覽器可開啟 localhost:5173 |
| Server 連線 | Inspector 介面 | 顯示綠色「Connected」 |
| Server Info | 點擊 Server Info | 顯示正確的 Name 和 Version |
| 日誌可見 | 點擊 Logs | 可看到 JSON-RPC 訊息 |

**自我檢核清單**
- [ ] Node.js 18+ 已安裝
- [ ] `npx @modelcontextprotocol/inspector` 可正常執行
- [ ] Inspector 能成功連線 Hello Server
- [ ] 能在 Logs 頁面看到通訊紀錄

---

## 開發流程建議

學會 Inspector 後，建議採用以下開發流程：

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

這樣可以大幅縮短開發迭代時間。

---

## 下一步

Inspector 測試就緒後，繼續學習：
- [U03｜Tool Capability](U03-Tool-Capability.md) — 讓 AI「做事」
- [U04｜Resource Capability](U04-Resource-Capability.md) — 讓 AI「讀資料」
- [U05｜Prompt Capability](U05-Prompt-Capability.md) — 讓 AI「套模板」
