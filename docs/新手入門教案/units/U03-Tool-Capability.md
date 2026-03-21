# U03｜Tool Capability

> 學習 MCP 最常用的功能——讓 AI 能「做事」。
>
> 預估時數：120 min
> 前置依賴：U02

---

## ① 為什麼先教這個？

Tool 是 MCP 最核心、最常用的功能。它讓 AI 能夠執行你定義的任意函式：計算數學、查詢資料庫、呼叫 API、操作檔案系統⋯⋯幾乎任何程式能做的事，都能透過 Tool 讓 AI 執行。

學會 Tool 後，你就能讓 Claude 成為你的程式助手，不只能回答問題，還能實際執行操作。

---

## ② 對應素材

| 類型 | 路徑/連結 | 說明 |
|------|----------|------|
| 範例程式碼 | `examples/U03-tools/main.go` | 含 add、greet、calculate 三個 Tool（Input struct 亦定義於此） |
| SDK 文件 | https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk/mcp | API 文件 |

---

## ③ 核心觀念

### 1. Tool 的組成元素

```go
mcp.AddTool(server, &mcp.Tool{
    Name:        "add",           // 工具名稱（唯一識別）
    Description: "加法計算",        // 描述（AI 用來判斷何時使用）
}, handlerFunc)
```

| 元素 | 說明 | 重要性 |
|------|------|--------|
| Name | 工具的唯一名稱 | 必填，用於呼叫識別 |
| Description | 工具的用途說明 | 必填，AI 依此判斷何時使用 |
| Handler | 實際執行的函式 | 必填，處理邏輯 |
| InputSchema | 輸入參數定義 | 自動從 struct 產生 |

### 2. Input Schema

使用 Go struct 定義輸入參數，SDK 會自動轉換為 JSON Schema：

```go
type AddInput struct {
    A float64 `json:"a" jsonschema:"第一個數字"`
    B float64 `json:"b" jsonschema:"第二個數字"`
}
```

| Tag | 說明 |
|-----|------|
| `json:"xxx"` | JSON 欄位名稱 |
| `jsonschema:"xxx"` | 欄位說明（給 AI 看） |
| `jsonschema:"required"` | 標記為必填欄位 |

### 3. Handler 函式簽名

```go
func handlerFunc(ctx context.Context, req *mcp.CallToolRequest, input InputType) (
    *mcp.CallToolResult,  // 回傳結果
    any,                  // metadata（通常為 nil）
    error,                // 錯誤
)
```

### 4. CallToolResult 結構

```go
return &mcp.CallToolResult{
    Content: []mcp.Content{
        &mcp.TextContent{Text: "計算結果：8"},
    },
}, nil, nil
```

支援的 Content 類型：
| 類型 | 說明 |
|------|------|
| `TextContent` | 純文字 |
| `ImageContent` | 圖片（base64 或 URL） |
| `EmbeddedResource` | 嵌入的 Resource |

### 5. mcp.AddTool 註冊流程

```go
// 1. 定義輸入結構
type GreetInput struct { /* json + jsonschema tag ... (略) */ }

// 2. 定義 Handler
func greetHandler(ctx context.Context, req *mcp.CallToolRequest, input GreetInput) (
    *mcp.CallToolResult, any, error,
) { /* 組合問候訊息並回傳 ... (略) */ }

// 3. 註冊到 Server
mcp.AddTool(server, &mcp.Tool{ /* Name, Description ... (略) */ }, greetHandler)
```

> 完整範例見 4.3

### 6. 錯誤處理最佳實踐

```go
func divideHandler(ctx context.Context, req *mcp.CallToolRequest, input DivideInput) (
    *mcp.CallToolResult, any, error,
) {
    // 輸入驗證
    if input.B == 0 {
        return &mcp.CallToolResult{
            Content: []mcp.Content{
                &mcp.TextContent{Text: "錯誤：除數不能為零"},
            },
            IsError: true,  // 標記為錯誤
        }, nil, nil
    }

    result := input.A / input.B
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: fmt.Sprintf("結果：%v", result)},
        },
    }, nil, nil
}
```

---

## ④ 動手做

### [必做] 4.1 建立專案結構

```bash
cd ~/mcp-test
mkdir -p examples/U03-tools
```

### [必做] 4.2 實作 add Tool

建立 `examples/U03-tools/main.go`：

```go
package main

import (
	// "context" ... (略)
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ==================== 輸入參數結構 ====================

// AddInput 定義 add Tool 的輸入參數
// json tag 指定 JSON 欄位名稱
// jsonschema tag 提供欄位說明給 AI 參考
type AddInput struct {
	A float64 `json:"a" jsonschema:"第一個數字"`
	B float64 `json:"b" jsonschema:"第二個數字"`
}

// ==================== Tool Handlers ====================

// addHandler 處理 add Tool 的請求
// 計算兩個數字的和
func addHandler(ctx context.Context, req *mcp.CallToolRequest, input AddInput) (
	*mcp.CallToolResult, any, error,
) {
	// 計算結果
	result := input.A + input.B
	// 回傳 CallToolResult，Content 包含 TextContent
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%v + %v = %v", input.A, input.B, result),
			},
		},
	}, nil, nil
}

// ==================== Main ====================

func main() {
	// Server 啟動骨架（見結構示範檔 §1）log.SetOutput ... (略)

	// ===== 以下為本單元新增內容 =====

	// 註冊 add Tool
	// Name: Tool 的唯一識別名稱
	// Description: 說明 Tool 的用途，AI 會根據這個描述決定何時使用
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "計算兩個數字的和。輸入兩個數字 a 和 b，回傳 a + b 的結果。",
	}, addHandler)

	// 使用 Stdio Transport 執行（見結構示範檔 §1）
	// if err := server.Run( ... (略)
}
```

### [必做] 4.3 新增 greet Tool

在 `main.go` 中新增：

> **備註**：此處先展示基本功能，不含輸入驗證。驗證邏輯見 [4.7 加入輸入驗證](#延伸-47-加入輸入驗證)。

```go
// GreetInput 定義 greet Tool 的輸入參數
type GreetInput struct {
	Name string `json:"name" jsonschema:"要問候的人的名字"`
}

// greetHandler 處理 greet Tool 的請求
// 向指定的人打招呼
func greetHandler(ctx context.Context, req *mcp.CallToolRequest, input GreetInput) (
	*mcp.CallToolResult, any, error,
) {
	// 組合問候訊息並回傳
	greeting := fmt.Sprintf("Hello, %s! 歡迎使用 MCP Tools!", input.Name)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: greeting},
		},
	}, nil, nil
}

// 在 main() 中註冊 greet Tool
mcp.AddTool(server, &mcp.Tool{
    Name:        "greet",
    Description: "向指定的人打招呼。輸入名字，回傳問候訊息。",
}, greetHandler)
```

### [必做] 4.4 使用 MCP Inspector 測試

> 操作方式見 [MCP Inspector 測試指引](../common/MCP-Inspector-測試指引.md)

```bash
npx @modelcontextprotocol/inspector go run ./examples/U03-tools/main.go
```

本單元測試重點：
- `Tools` 頁籤中的 `add` 和 `greet`
- 輸入參數並確認回傳結果正確

### [延伸] 4.5 設定 Claude Desktop 並測試

> 設定流程見 [Claude Desktop 設定指引](../common/Claude-Desktop-設定指引.md)

```bash
go build -o tools-demo ./examples/U03-tools/
```

config.json 中的 Server 名稱為 `tools-demo`，command 指向編譯產物的絕對路徑。

重啟後測試：
- 「請計算 123 + 456」→ 應呼叫 `add` Tool
- 「跟 Alice 打招呼」→ 應呼叫 `greet` Tool

### [延伸] 4.6 實作 calculate Tool

支援四則運算：

```go
// CalculateInput 定義 calculate Tool 的輸入參數
// 支援四則運算
type CalculateInput struct {
	A        float64 `json:"a" jsonschema:"第一個數字"`
	B        float64 `json:"b" jsonschema:"第二個數字"`
	Operator string  `json:"operator" jsonschema:"運算符號(+, -, *, /)"`
}

// calculateHandler 處理 calculate Tool 的請求
// 支援四則運算（+、-、*、/）
func calculateHandler(ctx context.Context, req *mcp.CallToolRequest, input CalculateInput) (
	*mcp.CallToolResult, any, error,
) {
	// 宣告計算結果與錯誤訊息變數
	var result float64
	var errMsg string

	// 根據運算符號執行對應運算
	switch input.Operator {
	case "+":
		result = input.A + input.B
	// case "-": ... (略)
	// case "*": ... (略)
	case "/":
		// 除法需要檢查除數是否為零
		if input.B == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "錯誤：除數不能為零"},
				},
				IsError: true, // 標記為錯誤，AI 會知道這是失敗結果
			}, nil, nil
		}
		result = input.A / input.B
	default:
		errMsg = fmt.Sprintf("錯誤：不支援的運算符號「%s」", input.Operator)
	}

	// 如果有錯誤訊息，回傳錯誤結果
	if errMsg != "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: errMsg},
			},
			IsError: true,
		}, nil, nil
	}

	// 回傳計算結果
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%v %s %v = %v", input.A, input.Operator, input.B, result),
			},
		},
	}, nil, nil
}

// 在 main() 中註冊 calculate Tool
mcp.AddTool(server, &mcp.Tool{
	Name:        "calculate",
	Description: "四則運算計算機。輸入兩個數字和運算符號（+、-、*、/），回傳計算結果。",
}, calculateHandler)
```

### [延伸] 4.7 加入輸入驗證

為 greet Tool 加入名字長度驗證：

```go
// greetHandler 處理 greet Tool 的請求
// 向指定的人打招呼，包含輸入驗證
func greetHandler(ctx context.Context, req *mcp.CallToolRequest, input GreetInput) (
	*mcp.CallToolResult, any, error,
) {
	// 輸入驗證：檢查名字是否為空
	if len(input.Name) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "錯誤：請提供名字"},
			},
			IsError: true,  // 標記為錯誤
		}, nil, nil
	}

	// 輸入驗證：檢查名字長度
	// if len(input.Name) > 100 { ... (略)

	// 組合問候訊息並回傳（同 4.3）
	// fmt.Sprintf("Hello, %s! ... (略)
}
```

---

## ⑤ 踩坑提示

### 踩坑 1：Tool 未出現在 Claude 選單

**現象**
重啟 Claude Desktop 後，連接器清單沒有顯示你註冊的 Tool。

**原因**
- `mcp.AddTool()` 未執行
- Server 名稱與 config.json 不符

**解法**
```bash
# 1. 確認 Server 有正確回應 tools/list
# 查看 Claude log
tail -f ~/Library/Logs/Claude/mcp*.log

# 2. 確認 config.json 的 Server 名稱正確
```

---

### 踩坑 2：輸入參數解析失敗

**現象**
Claude 呼叫 Tool 時，參數都是零值或空字串。

**原因**
- struct tag 格式錯誤
- JSON 欄位名稱不匹配

**解法**
```go
// ❌ 錯誤：缺少 json tag
type Input struct {
    Name string
}

// ✅ 正確：有 json tag
type Input struct {
    Name string `json:"name"`
}
```

---

### 踩坑 3：回傳結果為空

**現象**
Tool 執行了，但 Claude 顯示沒有結果。

**原因**
`CallToolResult.Content` 為空或 nil。

**解法**
```go
// ❌ 錯誤：Content 為 nil
return &mcp.CallToolResult{}, nil, nil

// ✅ 正確：有 Content
return &mcp.CallToolResult{
    Content: []mcp.Content{
        &mcp.TextContent{Text: "結果"},
    },
}, nil, nil
```

---

### 踩坑 4：Description 不清楚導致 AI 誤用

**現象**
Claude 在不該使用此 Tool 時使用了，或該使用時沒使用。

**原因**
Tool 的 Description 不夠明確。

**解法**
```go
// ❌ 模糊的描述
Description: "處理數字"

// ✅ 清晰的描述
Description: "計算兩個數字的和，輸入兩個數字 a 和 b，回傳 a + b 的結果"
```

---

## ⑥ 驗收標準（DoD）

| 項目 | 驗證方式 | 預期結果 |
|------|---------|---------|
| add Tool 註冊成功 | Claude 連接器顯示 | 看到 add Tool |
| add 計算正確 | 問「3 + 5 等於多少」 | Claude 呼叫 add，回傳 8 |
| greet Tool 運作 | 問「跟 Bob 打招呼」 | 回傳 "Hello, Bob!" |
| 錯誤處理正確 | 測試除以零 | 回傳友善錯誤訊息 |

**自我檢核清單**
- [ ] 編譯無錯誤
- [ ] add Tool 正確計算
- [ ] greet Tool 正確問候
- [ ] Description 足夠清晰
- [ ] 有基本的輸入驗證

---

## 下一步

掌握 Tool 後，前往 [U04｜Resource Capability](U04-Resource-Capability.md) 學習如何讓 AI 「讀取資料」。
