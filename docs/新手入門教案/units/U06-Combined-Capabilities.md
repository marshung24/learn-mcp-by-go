# U06｜Combined Capabilities

> 整合 Tools、Resources、Prompts 三種功能到同一個 MCP Server。
>
> 預估時數：90 min
> 前置依賴：U03, U04, U05

---

## ① 為什麼先教這個？

實際的 MCP Server 通常不只有單一功能。一個完整的服務可能同時需要：
- **Tool**：執行操作（如計算、API 呼叫）
- **Resource**：提供資料（如設定檔、說明文件）
- **Prompt**：預設模板（如常用指令）

本單元練習將三種功能整合到同一個 Server，並學習如何組織程式碼結構。

---

## ② 對應素材

| 類型 | 路徑/連結 | 說明 |
|------|----------|------|
| 範例程式碼 | `examples/U06-combined/main.go` | 整合版 Server |
| 專案結構 | `examples/U06-combined/` | 分檔組織範例 |

---

## ③ 核心觀念

### 1. ServerCapabilities 設定

要同時啟用三種功能，需在建立 Server 時設定 capabilities：

```go
server := mcp.NewServer(&mcp.Implementation{
    Name:    "combined-demo",
    Version: "1.0.0",
}, &mcp.ServerOptions{
    Capabilities: &mcp.ServerCapabilities{
        Tools:     &mcp.ServerCapabilitiesTools{},
        Resources: &mcp.ServerCapabilitiesResources{},
        Prompts:   &mcp.ServerCapabilitiesPrompts{},
    },
})
```

### 2. 功能組織策略

| 策略 | 適用場景 | 優點 |
|------|---------|------|
| **單一檔案** | 小型專案、少量功能 | 簡單直接 |
| **分檔組織** | 中型專案、功能較多 | 易於維護 |
| **分套件** | 大型專案、團隊開發 | 模組化、可重用 |

### 3. 建議的專案結構

```
examples/U06-combined/
├── main.go              # 入口點
├── tools/
│   └── calculator.go    # Tool 定義
├── resources/
│   └── help.go          # Resource 定義
├── prompts/
│   └── templates.go     # Prompt 定義
└── internal/
    └── common.go        # 共用函式
```

### 4. 功能間的配合

三種 capabilities 可以互相配合：
- **Resource** 提供說明文件 → **Prompt** 引用說明 → **Tool** 執行操作
- **Tool** 產生結果 → 存成 **Resource** → 供後續讀取

---

## ④ 動手做

### [必做] 4.1 建立專案結構

```bash
cd ~/mcp-test
mkdir -p examples/U06-combined
```

### [必做] 4.2 建立整合版 main.go

建立 `examples/U06-combined/main.go`：

```go
// Package main 實作 Combined Capabilities 範例
// 整合 Tools、Resources、Prompts 三種功能到同一個 MCP Server
package main

import (
	// 定位線索 "context" ... (略)
)

// ==================== 全域變數 ====================

// startTime 記錄 Server 啟動時間，用於計算 uptime
var startTime time.Time

// ==================== Tools ====================

// CalculateInput 定義 calculate Tool 的輸入參數
type CalculateInput struct {
	// 定位線索（見結構示範檔 §3）A float64 `json:"a"` ... (略)
}

// calculateHandler 處理四則運算
func calculateHandler(ctx context.Context, req *mcp.CallToolRequest, input CalculateInput) (
	*mcp.CallToolResult, any, error,
) {
	// 定位線索（實作細節見 U03 或 examples/U06-combined/main.go）var result float64 ... (略)
}

// ==================== Resources ====================

// helpResourceHandler 提供使用說明
func helpResourceHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 定位線索（實作細節見 U04 或 examples/U06-combined/main.go）helpContent := `# Combined ... (略)
}

// statusResourceHandler 提供系統狀態
func statusResourceHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 定位線索（實作細節見 U04 或 examples/U06-combined/main.go）uptime := time.Since(startTime) ... (略)
}

// ==================== Prompts ====================

// mathProblemHandler 產生數學問題解題模板
func mathProblemHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 定位線索（實作細節見 U05 或 examples/U06-combined/main.go）problem := "" ... (略)
}

// explainCodeHandler 產生程式碼解釋模板
func explainCodeHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 定位線索（實作細節見 U05 或 examples/U06-combined/main.go）code := "" ... (略)
}

// ==================== Main ====================

func main() {
	// 記錄啟動時間
	startTime = time.Now()

	// 設定 log 輸出到 stderr
	log.SetOutput(os.Stderr)
	log.Println("Starting combined-demo server...")

	// 建立 MCP Server
	// 設定 ServerCapabilities 以啟用所有功能
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "combined-demo",
		Version: "1.0.0",
	}, nil)

	// ==================== 註冊 Tools ====================
	mcp.AddTool(server, &mcp.Tool{
		Name:        "calculate",
		Description: "四則運算計算機，支援加減乘除",
	}, calculateHandler)
	log.Println("Tool 'calculate' registered")

	// ==================== 註冊 Resources ====================
	server.AddResource(&mcp.Resource{
		URI:         "file://help",
		Name:        "Help",
		Description: "使用說明文件",
		MIMEType:    "text/markdown",
	}, helpResourceHandler)

	server.AddResource(&mcp.Resource{
		URI:         "data://status",
		Name:        "Status",
		Description: "系統狀態資訊",
		MIMEType:    "application/json",
	}, statusResourceHandler)
	log.Println("Resources registered: help, status")

	// ==================== 註冊 Prompts ====================
	server.AddPrompt(&mcp.Prompt{
		Name:        "math-problem",
		Description: "數學問題解題模板",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "problem",
				Description: "要解答的數學問題",
				Required:    true,
			},
		},
	}, mathProblemHandler)

	server.AddPrompt(&mcp.Prompt{
		Name:        "explain-code",
		Description: "程式碼解釋模板",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "code",
				Description: "要解釋的程式碼",
				Required:    true,
			},
			{
				Name:        "language",
				Description: "程式語言（選填）",
				Required:    false,
			},
		},
	}, explainCodeHandler)
	log.Println("Prompts registered: math-problem, explain-code")

	log.Printf("Server ready with Tools, Resources, and Prompts")
	log.Println("Waiting for connections...")

	// 使用 Stdio Transport 執行
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
```

### [必做] 4.3 使用 MCP Inspector 測試

> 操作方式見 [MCP Inspector 測試指引](../common/MCP-Inspector-測試指引.md)

```bash
npx @modelcontextprotocol/inspector go run ./examples/U06-combined/main.go
```

本單元測試重點（三種 capabilities 皆需驗證）：
- **Tools**：測試 `calculate` 工具
- **Resources**：讀取 `help` 和 `status`
- **Prompts**：使用 `math-problem` 和 `explain-code`

### [延伸] 4.4 設定 Claude Desktop 並測試

> 設定流程見 [Claude Desktop 設定指引](../common/Claude-Desktop-設定指引.md)

```bash
go build -o combined-demo ./examples/U06-combined/
```

config.json 中的 Server 名稱為 `combined-demo`，command 指向編譯產物的絕對路徑。

重啟後測試：
1. **Tool 測試**：「請計算 100 除以 7」
2. **Resource 測試**：讀取 help 或 status Resource
3. **Prompt 測試**：使用 math-problem 或 explain-code Prompt

### [延伸] 4.5 分檔組織

將功能拆分到不同檔案：

**tools/calculator.go**
```go
package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CalculateInput struct {
	A        float64 `json:"a" jsonschema:"第一個數字"`
	B        float64 `json:"b" jsonschema:"第二個數字"`
	Operator string  `json:"operator" jsonschema:"運算符號"`
}

func CalculateHandler(ctx context.Context, req *mcp.CallToolRequest, input CalculateInput) (
	*mcp.CallToolResult, any, error,
) {
	// 定位線索 實作 var result float64 ... (略)
}

func RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "calculate",
		Description: "計算機",
	}, CalculateHandler)
}
```

**main.go（簡化版）**
```go
package main

import (
	"myproject/tools"
	"myproject/resources"
	"myproject/prompts"
)

func main() {
	server := mcp.NewServer(...)

	tools.RegisterTools(server)
	resources.RegisterResources(server)
	prompts.RegisterPrompts(server)

	server.Run(...)
}
```

---

## ⑤ 踩坑提示

### 踩坑 1：部分功能未生效

**現象**
只有 Tool 可用，Resource 和 Prompt 沒反應。

**原因**
ServerCapabilities 未正確設定（雖然 SDK 可能自動處理）。

**解法**
確認 SDK 版本，或明確設定 capabilities：
```go
server := mcp.NewServer(&mcp.Implementation{...}, &mcp.ServerOptions{
    Capabilities: &mcp.ServerCapabilities{
        Tools:     &mcp.ServerCapabilitiesTools{},
        Resources: &mcp.ServerCapabilitiesResources{},
        Prompts:   &mcp.ServerCapabilitiesPrompts{},
    },
})
```

---

### 踩坑 2：功能名稱衝突

**現象**
註冊多個同名功能時出現非預期行為。

**原因**
Tool、Resource、Prompt 的名稱在各自類型中必須唯一。

**解法**
```go
// ✅ 不同類型可以同名（但不建議）
mcp.AddTool(server, &mcp.Tool{Name: "help", ...}, ...)
server.AddResource(&mcp.Resource{Name: "help", ...}, ...)

// ✅ 建議使用有區別的命名
mcp.AddTool(server, &mcp.Tool{Name: "get-help", ...}, ...)
server.AddResource(&mcp.Resource{Name: "help-doc", ...}, ...)
```

---

### 踩坑 3：import cycle

**現象**
分檔後編譯出現 import cycle 錯誤。

**原因**
套件間互相 import。

**解法**
- 使用 internal 套件放共用邏輯
- 避免套件間互相依賴
- 主程式 import 子套件，子套件不 import 主程式

---

## ⑥ 驗收標準（DoD）

| 項目 | 驗證方式 | 預期結果 |
|------|---------|---------|
| calculate Tool | 「計算 50 * 3」 | 回傳 150 |
| help Resource | 讀取 help | 顯示使用說明 |
| status Resource | 讀取 status | 顯示 JSON 狀態 |
| math-problem Prompt | 選用並輸入問題 | 產生解題提示詞 |
| explain-code Prompt | 選用並貼程式碼 | 產生解釋提示詞 |

**自我檢核清單**
- [ ] 三種 capabilities 都能正常運作
- [ ] 功能之間不互相干擾
- [ ] 程式碼組織清晰
- [ ] status Resource 顯示正確的 uptime

---

## 下一步

整合完成後，前往 [U07｜HTTP Transport](U07-HTTP-Transport.md) 學習如何從 Stdio 切換到 HTTP 模式。
