# U04｜Resource Capability

> 學習如何讓 AI 「讀取資料」——提供檔案、設定、日誌等資料給 AI 分析。
>
> 預估時數：90 min
> 前置依賴：U02

---

## ① 為什麼先教這個？

Tool 讓 AI 能「做事」，Resource 則讓 AI 能「讀資料」。想像你有一份日誌檔、一個設定檔、或是資料庫查詢結果，透過 Resource，AI 可以直接讀取這些資料並進行分析、摘要或回答問題。

Resource 與 Tool 的差異在於：Tool 是主動執行動作，Resource 是被動提供資料。

---

## ② 對應素材

| 類型 | 路徑/連結 | 說明 |
|------|----------|------|
| 範例程式碼 | `examples/U04-resources/main.go` | 含靜態與動態 Resource |
| 資料檔案 | `examples/U04-resources/data/sample.json` | 範例 JSON 資料 |
| SDK 文件 | https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk/mcp | API 文件 |

---

## ③ 核心觀念

### 1. Resource 的組成元素

```go
server.AddResource(&mcp.Resource{
    URI:         "file://readme",           // 唯一識別 URI
    Name:        "README",                  // 顯示名稱
    Description: "專案說明文件",             // 描述
    MIMEType:    "text/plain",             // MIME 類型
}, readHandler)
```

| 元素 | 說明 | 範例 |
|------|------|------|
| URI | 唯一識別符 | `file://readme`, `data://users` |
| Name | 顯示名稱 | "README", "User List" |
| Description | 用途說明 | "專案說明文件" |
| MimeType | 內容類型 | `text/plain`, `application/json` |

### 2. 靜態 vs 動態 Resource

| 類型 | 說明 | 適用場景 |
|------|------|---------|
| **靜態** | 內容固定不變 | 設定檔、README、範本 |
| **動態** | 每次讀取時重新計算 | 系統狀態、即時資料、日誌 |

### 3. URI 格式慣例

| Scheme | 用途 | 範例 |
|--------|------|------|
| `file://` | 檔案類資源 | `file://config.json` |
| `data://` | 資料類資源 | `data://users` |
| `log://` | 日誌類資源 | `log://app` |
| `custom://` | 自訂類型 | `custom://anything` |

### 4. ReadResourceResult 結構

```go
// 回傳 ReadResourceResult，包含一或多筆 ResourceContents
return &mcp.ReadResourceResult{
    Contents: []*mcp.ResourceContents{
        {
            URI:      "file://readme",           // 必須與註冊時的 URI 一致
            MIMEType: "text/plain",              // 對應內容格式
            Text:     "這是 README 內容",          // 純文字內容（TextResourceContents）
        },
    },
}, nil
```

支援的 Contents 類型：
| 類型 | 說明 |
|------|------|
| `TextResourceContents` | 純文字內容 |
| `BlobResourceContents` | 二進位內容（base64） |

### 5. 常見 MIME Types

| MIME Type | 用途 |
|-----------|------|
| `text/plain` | 純文字 |
| `text/markdown` | Markdown |
| `application/json` | JSON |
| `text/html` | HTML |
| `text/csv` | CSV |

### 6. Resource Handler 簽名

```go
func resourceHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
    *mcp.ReadResourceResult,  // 回傳結果
    error,                    // 錯誤
)
```

---

## ④ 動手做

### [必做] 4.1 建立專案結構

```bash
cd ~/mcp-test
mkdir -p examples/U04-resources/data
```

### [必做] 4.2 建立範例資料檔

建立 `examples/U04-resources/data/sample.json`：

```json
{
  "name": "MCP Learning Project",
  "version": "1.0.0",
  "description": "學習 MCP Server 開發的範例專案",
  "features": [
    "Tools",
    "Resources",
    "Prompts"
  ]
}
```

### [必做] 4.3 實作靜態 Resource

建立 `examples/U04-resources/main.go`：

```go
// Package main 實作 Resource Capability 範例
// 這是 U04 Resource Capability 的範例程式碼，展示如何提供資料給 AI 讀取
package main

import (
	// "context" ... (略)
)

// ==================== Resource Handlers ====================

// readmeHandler 處理靜態 README Resource 的讀取請求
// 靜態 Resource：內容固定不變
func readmeHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 定義靜態內容
	content := `# MCP Resources Demo
...（省略靜態文字內容）
`
	// 回傳 ReadResourceResult，包含 TextResourceContents
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "file://readme",       // 必須與註冊時的 URI 一致
				MIMEType: "text/markdown",        // 對應內容格式
				Text:     content,                // 純文字內容
			},
		},
	}, nil
}

// ==================== Main ====================

func main() {
	// Server 啟動骨架（見結構示範檔 §1）log.SetOutput ... (略)

	// ===== 以下為本單元新增內容 =====

	// 註冊靜態 Resource：README
	// URI: 唯一識別符，使用 file:// 表示檔案類資源
	// Name: 顯示名稱
	// Description: 說明此 Resource 的用途
	// MIMEType: 內容類型
	server.AddResource(&mcp.Resource{
		URI:         "file://readme",
		Name:        "README",
		Description: "專案說明文件",
		MIMEType:    "text/markdown",
	}, readmeHandler)

	// 使用 Stdio Transport 執行
	// （見結構示範檔 §1）if err := server.Run( ... (略)
}
```

### [必做] 4.4 新增動態 Resource

加入顯示當前時間的 Resource：

```go
// timeHandler 處理動態時間 Resource 的讀取請求
// 動態 Resource：每次讀取時重新計算內容
func timeHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	now := time.Now()

	// 建立時間資訊結構
	content := map[string]interface{}{
		"timestamp":   now.Unix(),
		"formatted":   now.Format("2006-01-02 15:04:05"),
		"timezone":    now.Location().String(),
		// "weekday": now.Weekday().String(), ... (略)
	}

	// 轉換為 JSON 格式
	jsonBytes, err := json.MarshalIndent(content, "", "  ")
	// if err != nil { return nil, err } ... (略)

	// 回傳 ReadResourceResult（結構同 readmeHandler）
	// return &mcp.ReadResourceResult{ ... (略)
}

// 在 main() 中註冊動態 Resource：當前時間
// 使用 data:// 表示資料類資源
server.AddResource(&mcp.Resource{
    URI:         "data://current-time",
    Name:        "Current Time",
    Description: "當前系統時間（動態更新）",
    MIMEType:    "application/json",
}, timeHandler)
```

### [必做] 4.5 讀取外部檔案

加入讀取 sample.json 的 Resource：

```go
// configHandler 處理設定檔 Resource 的讀取請求
// 從外部檔案讀取內容
func configHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 取得執行檔所在目錄
	// 注意：實際部署時可能需要調整路徑處理方式
	// （路徑解析邏輯，詳見 example）exePath, err := os.Executable() ... (略)

	// 讀取檔案內容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 讀取失敗時回傳錯誤訊息，但不中斷程式
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "file://config",
					MIMEType: "text/plain",
					Text:     "無法讀取設定檔：" + err.Error(),
				},
			},
		}, nil
	}

	// 回傳 ReadResourceResult（結構同 readmeHandler）
	// return &mcp.ReadResourceResult{ ... (略)
}

// 在 main() 中註冊檔案 Resource：設定檔
server.AddResource(&mcp.Resource{
    URI:         "file://config",
    Name:        "Config File",
    Description: "專案設定檔（sample.json）",
    MIMEType:    "application/json",
}, configHandler)
```

### [必做] 4.6 使用 MCP Inspector 測試

> 操作方式見 [MCP Inspector 測試指引](../common/MCP-Inspector-測試指引.md)

```bash
npx @modelcontextprotocol/inspector go run ./examples/U04-resources/main.go
```

本單元測試重點：
- `Resources` 頁籤中的 `readme`、`current-time`、`config`
- 確認各 Resource 回傳內容正確

### [延伸] 4.7 設定 Claude Desktop 並測試

> 設定流程見 [Claude Desktop 設定指引](../common/Claude-Desktop-設定指引.md)

```bash
go build -o resources-demo ./examples/U04-resources/
```

config.json 中的 Server 名稱為 `resources-demo`，command 指向編譯產物的絕對路徑。

重啟後測試：在 Claude Desktop 中讀取 Resource。

### [延伸] 4.8 實作目錄列表 Resource

```go
// filesHandler 處理目錄列表 Resource 的讀取請求
// 列出當前目錄的檔案資訊
func filesHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 讀取當前目錄內容
	entries, err := os.ReadDir(".")
	// if err != nil { return nil, err } ... (略)

	// 建立檔案資訊列表
	// （遍歷 entries）var files []map[string]interface{} ... (略)

	// 轉換為 JSON 格式
	// jsonBytes, err := json.MarshalIndent( ... (略)

	// 回傳目錄列表
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "data://files",
				MIMEType: "application/json",
				Text:     string(jsonBytes),
			},
		},
	}, nil
}

// 在 main() 中註冊動態 Resource：目錄列表
server.AddResource(&mcp.Resource{
    URI:         "data://files",
    Name:        "Directory Listing",
    Description: "當前目錄的檔案列表",
    MIMEType:    "application/json",
}, filesHandler)
```

### [延伸] 4.9 加入快取機制

```go
// 快取相關變數宣告
var (
	configCache     string
	configCacheTime time.Time
	cacheDuration   = 5 * time.Minute
)

// cachedConfigHandler 帶快取的設定檔 Resource Handler
// 在 cacheDuration 內回傳快取，過期後重新讀取檔案
func cachedConfigHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 檢查快取是否有效
	if time.Since(configCacheTime) < cacheDuration && configCache != "" {
		// 快取命中，直接回傳
		// return &mcp.ReadResourceResult{ ... (略)
	}

	// 快取過期或為空，重新讀取檔案
	// content, err := ioutil.ReadFile( ... (略)

	// 更新快取
	// configCache = string(content) ... (略)

	// 回傳新讀取的內容
	// return &mcp.ReadResourceResult{ ... (略)
}
```

---

## ⑤ 踩坑提示

### 踩坑 1：Resource 內容亂碼

**現象**
Claude 顯示的 Resource 內容出現亂碼或格式錯亂。

**原因**
MimeType 設定錯誤，例如 JSON 內容設成 `text/plain`。

**解法**
```go
// ❌ 錯誤
MIMEType: "text/plain"  // 對於 JSON 內容

// ✅ 正確
MIMEType: "application/json"  // 對於 JSON 內容
```

---

### 踩坑 2：讀取檔案失敗

**現象**
Resource Handler 回傳「找不到檔案」錯誤。

**原因**
路徑為相對路徑，但執行時的工作目錄不同。

**解法**
```go
// ❌ 相對路徑（不可靠）
content, err := ioutil.ReadFile("data/sample.json")

// ✅ 相對於執行檔的路徑
exePath, _ := os.Executable()
exeDir := filepath.Dir(exePath)
filePath := filepath.Join(exeDir, "data", "sample.json")
content, err := ioutil.ReadFile(filePath)

// ✅ 或使用絕對路徑
content, err := ioutil.ReadFile("/absolute/path/to/sample.json")
```

---

### 踩坑 3：URI 格式不一致

**現象**
註冊 Resource 時的 URI 與回傳時的 URI 不同，導致 Client 無法匹配。

**原因**
註冊時寫 `file://readme`，回傳時寫 `file:///readme`。

**解法**
```go
// 確保 URI 完全一致
const readmeURI = "file://readme"

server.AddResource(&mcp.Resource{
    URI: readmeURI,  // 使用常數
    // Name: "README", ... (略)
}, func(...) {
    return &mcp.ReadResourceResult{
        Contents: []*mcp.ResourceContents{
            {
                URI: readmeURI,  // 同一常數
                // MIMEType: "text/markdown", ... (略)
            },
        },
    }, nil
})
```

---

## ⑥ 驗收標準（DoD）

| 項目 | 驗證方式 | 預期結果 |
|------|---------|---------|
| README Resource | 在 Claude 中讀取 | 顯示 Markdown 格式的說明 |
| Time Resource | 在 Claude 中讀取 | 顯示當前時間 JSON |
| Config Resource | 在 Claude 中讀取 | 顯示 sample.json 內容 |
| 動態更新 | 連續讀取 Time | 時間有變化 |

**自我檢核清單**
- [ ] 三種 Resource 都能正常讀取
- [ ] MIME Type 設定正確
- [ ] 檔案路徑處理正確
- [ ] 動態 Resource 每次讀取都更新

---

## 下一步

掌握 Resource 後，前往 [U05｜Prompt Capability](U05-Prompt-Capability.md) 學習如何建立提示詞模板。
