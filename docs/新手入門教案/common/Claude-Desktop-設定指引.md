# Claude Desktop 設定指引

> 各單元共用的 Claude Desktop MCP Server 設定流程。
> 單元 MD 中以連結指向本檔，不重複撰寫。
> 完整 Claude Desktop 設定另見 [docs/setup/claude-desktop.md](../../setup/claude-desktop.md)。

---

## 設定流程

### Step 1：編譯 Server

```bash
go build -o <BINARY_NAME> ./examples/<UNIT_DIR>/
```

### Step 2：編輯設定檔

macOS 設定檔路徑：

```
~/Library/Application Support/Claude/claude_desktop_config.json
```

設定格式：

```json
{
  "mcpServers": {
    "<SERVER_NAME>": {
      "command": "/完整路徑/<BINARY_NAME>"
    }
  }
}
```

> **注意**：`command` 必須使用**絕對路徑**，不支援 `~` 或相對路徑。

### Step 3：重啟 Claude Desktop

設定修改後需完全重啟（Quit → 重新開啟），僅關閉視窗不會重新載入設定。

### Step 4：驗證

在 Claude Desktop 中：
1. 點擊輸入框旁的 "+" 按鈕
2. 確認列表中出現你的 Server 名稱
3. 嘗試與 Server 互動（如：請 Claude 呼叫你的 Tool）

---

## 各單元對照表

| 單元 | 編譯指令 | Server 名稱 | 測試建議 |
|------|---------|-------------|---------|
| U01 | `go build -o hello-mcp ./examples/U01-hello/` | hello-mcp | 確認 Server 出現在列表 |
| U03 | `go build -o tools-demo ./examples/U03-tools/` | tools-demo | 「請計算 123 + 456」、「跟 Alice 打招呼」 |
| U04 | `go build -o resources-demo ./examples/U04-resources/` | resources-demo | 讀取 Resource |
| U05 | `go build -o prompts-demo ./examples/U05-prompts/` | prompts-demo | 使用 Prompt 模板 |
| U06 | `go build -o combined-demo ./examples/U06-combined/` | combined-demo | 測試 Tool + Resource + Prompt |
| U09 | `go build -o weather-mvp ./examples/U09-weather-mvp/` | weather-mvp | 「查詢加州天氣警報」（NWS API 僅支援美國地區） |

> U07、U08 使用 HTTP Transport，不透過 Claude Desktop 的 stdio 機制，故不在此列。

---

## 單元內引用寫法

後續單元（U03 起）引用本指引的標準寫法：

```markdown
### [延伸] 4.X 設定 Claude Desktop 並測試

> 設定流程見 [Claude Desktop 設定指引](../common/Claude-Desktop-設定指引.md)

​```bash
go build -o xxx-demo ./examples/U0X-xxx/
​```

config.json 中的 Server 名稱為 `xxx-demo`，command 指向編譯產物的絕對路徑。

重啟後測試：
- 「請計算 123 + 456」→ 應呼叫 `add` Tool
- 「跟 Alice 打招呼」→ 應呼叫 `greet` Tool
```

---

## 常見問題

| 現象 | 原因 | 解法 |
|------|------|------|
| Server 未出現在列表 | 設定檔 JSON 格式錯誤 | 用 `jq . config.json` 驗證格式 |
| Server 出現但無法使用 | command 路徑錯誤 | 用 `ls -la /path/to/binary` 確認檔案存在且可執行 |
| 修改設定後無效 | 未完全重啟 | Quit Claude Desktop 再重新開啟 |
| 多個 Server 衝突 | mcpServers 中 key 重複 | 確保每個 Server 使用不同 key |
