# U05｜Prompt Capability

> 學習建立可重用的提示詞模板，提升與 AI 互動效率。
>
> 預估時數：60 min
> 前置依賴：U02

---

## ① 為什麼先教這個？

Prompt 讓你能預先定義常用的提示詞模板。想像你經常需要請 AI 做 Code Review、摘要文章、或生成報告——每次都要打一長串指令很累。透過 Prompt，你可以定義模板，使用時只需填入參數，大幅提升效率。

Prompt 就像是 AI 互動的「快捷鍵」或「巨集」。

---

## ② 對應素材

| 類型 | 路徑/連結 | 說明 |
|------|----------|------|
| 範例程式碼 | `examples/U05-prompts/main.go` | 含 code-review、summarize 兩個 Prompt |
| SDK 文件 | https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk/mcp | API 文件 |

---

## ③ 核心觀念

### 1. Prompt 的組成元素

```go
server.AddPrompt(&mcp.Prompt{
    Name:        "code-review",             // 唯一名稱
    Description: "對程式碼進行 Code Review", // 描述
    Arguments: []*mcp.PromptArgument{        // 參數列表
        {
            Name:        "code",
            Description: "要 Review 的程式碼",
            Required:    true,
        },
    },
}, promptHandler)
```

| 元素 | 說明 |
|------|------|
| Name | Prompt 唯一名稱 |
| Description | 用途說明 |
| Arguments | 可帶入的參數列表 |

### 2. PromptArgument 結構

```go
type PromptArgument struct {
    Name        string  // 參數名稱
    Description string  // 參數說明
    Required    bool    // 是否必填
}
```

### 3. Prompt Handler 簽名

```go
func promptHandler(ctx context.Context, req *mcp.GetPromptRequest) (
    *mcp.GetPromptResult,  // 回傳結果
    error,                 // 錯誤
)
```

### 4. GetPromptResult 結構

```go
return &mcp.GetPromptResult{
    Description: "Code Review 提示詞",
    Messages: []*mcp.PromptMessage{
        {
            Role:    "user",                                  // 訊息角色
            Content: &mcp.TextContent{Text: promptText},      // 訊息內容
        },
    },
}, nil
```

### 5. Role 類型

| Role | 說明 |
|------|------|
| `"user"` | 使用者訊息 |
| `mcp.RoleAssistant` | AI 回覆（少用） |

### 6. Prompt vs Tool 的差異

| 特性 | Tool | Prompt |
|------|------|--------|
| 用途 | 執行動作 | 提供模板 |
| 觸發方式 | AI 自動選擇 | 使用者主動選用 |
| 回傳內容 | 執行結果 | 訊息內容（交給 AI 處理） |

---

## ④ 動手做

### [必做] 4.1 建立專案結構

```bash
cd ~/mcp-test
mkdir -p examples/U05-prompts
```

### [必做] 4.2 實作 code-review Prompt

建立 `examples/U05-prompts/main.go`：

```go
package main

// import "context", "fmt", "log" ... (略)

// codeReviewHandler 處理 code-review Prompt 的請求
// 產生程式碼審查的提示詞模板
func codeReviewHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 從 Arguments 取得參數
	code := ""
	if req.Params.Arguments != nil {
		if c, ok := req.Params.Arguments["code"]; ok {
			code = c
		}
	}

	// 驗證必填參數
	// if code == "" { ... (略)

	// 組合提示詞文字：將審查項目與程式碼組合成完整 Prompt
	// promptText := fmt.Sprintf(`請對以下 ... (略)

	// 回傳 GetPromptResult：包含 Description 與 PromptMessage 列表
	return &mcp.GetPromptResult{
		Description: "Code Review 提示詞",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

// main 啟動 MCP Server 並註冊所有 Prompt
// 輸出: 透過 Stdio Transport 提供 Prompt 服務
func main() {
	// Server 啟動骨架（見結構示範檔 §1）log.SetOutput ... (略)

	// ===== 以下為本單元新增內容 =====

	// 註冊 code-review Prompt
	// Name: Prompt 的唯一名稱
	// Description: 說明此 Prompt 的用途
	// Arguments: 可帶入的參數列表
	server.AddPrompt(&mcp.Prompt{
		Name:        "code-review",
		Description: "對程式碼進行專業的 Code Review",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "code",             // 參數名稱（取用時的 key）
				Description: "要 Review 的程式碼", // 參數說明
				Required:    true,                // 是否必填
			},
		},
	}, codeReviewHandler)

	// 使用 Stdio Transport 執行（參考 U01）
	// if err := server.Run(context. ... (略)
}
```

### [必做] 4.3 新增 summarize Prompt

```go
// summarizeHandler 處理 summarize Prompt 的請求
// 產生文字摘要的提示詞模板
func summarizeHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 取得參數（與 codeReviewHandler 相同模式）
	// text := "" ... (略)
	maxWords := "100" // 選填參數給預設值
	// if req.Params.Arguments != nil { ... (略)

	// 驗證必填參數
	// if text == "" { ... (略)

	// 組合提示詞文字：將摘要規則與原文組合成完整 Prompt
	// promptText := fmt.Sprintf(`請將以下 ... (略)

	// 回傳 GetPromptResult：包含 Description 與 PromptMessage 列表
	return &mcp.GetPromptResult{
		Description: "文字摘要提示詞",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

// 在 main() 中註冊 summarize Prompt
server.AddPrompt(&mcp.Prompt{
    Name:        "summarize",
    Description: "將文字摘要為重點",
    Arguments: []*mcp.PromptArgument{
        {
            Name:        "text",
            Description: "要摘要的文字",
            Required:    true,              // 必填
        },
        {
            Name:        "max_words",
            Description: "摘要字數上限（預設 100）",
            Required:    false,             // 選填
        },
    },
}, summarizeHandler)
```

### [必做] 4.4 使用 MCP Inspector 測試

> 操作方式見 [MCP Inspector 測試指引](../common/MCP-Inspector-測試指引.md)

```bash
npx @modelcontextprotocol/inspector go run ./examples/U05-prompts/main.go
```

本單元測試重點：
- `Prompts` 頁籤中的 `code-review` 和 `summarize`
- 填入參數並確認回傳的 Prompt 文字正確

### [延伸] 4.5 設定 Claude Desktop 並測試

> 設定流程見 [Claude Desktop 設定指引](../common/Claude-Desktop-設定指引.md)

```bash
go build -o prompts-demo ./examples/U05-prompts/
```

config.json 中的 Server 名稱為 `prompts-demo`，command 指向編譯產物的絕對路徑。

重啟後測試：使用 Prompt 模板。

### [延伸] 4.6 實作多語言支援

```go
// translateHandler 處理 translate Prompt 的請求
// 產生翻譯的提示詞模板
func translateHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 取得參數（與 codeReviewHandler 相同模式）
	// text := "" ... (略)

	// 驗證必填參數
	// if text == "" { ... (略)

	// 組合提示詞文字：將翻譯要求與原文組合成完整 Prompt
	// promptText := fmt.Sprintf(`請將以下 ... (略)

	// 回傳 GetPromptResult：包含 Description 與 PromptMessage 列表
	// return &mcp.GetPromptResult{ ... (略)
}

// 在 main() 中註冊 translate Prompt
server.AddPrompt(&mcp.Prompt{
    Name:        "translate",
    Description: "翻譯文字到指定語言",
    Arguments: []*mcp.PromptArgument{
        {
            Name:        "text",
            Description: "要翻譯的文字",
            Required:    true,
        },
        {
            Name:        "target_language",
            Description: "目標語言（預設：繁體中文）",
            Required:    false,
        },
    },
}, translateHandler)
```

### [延伸] 4.7 建立系統分析 Prompt

```go
// debugHandler 處理 debug Prompt 的請求
// 產生錯誤分析的提示詞模板
func debugHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 取得參數（多個選填參數的範例）
	// errorMsg := "" ... (略)

	// 設定預設值
	// if errorMsg == "" { ... (略)

	// 組合提示詞文字：將錯誤資訊與分析要求組合成完整 Prompt
	// promptText := fmt.Sprintf(`請幫我 ... (略)

	// 回傳 GetPromptResult：包含 Description 與 PromptMessage 列表
	// return &mcp.GetPromptResult{ ... (略)
}
```

---

## ⑤ 踩坑提示

### 踩坑 1：Prompt 參數未帶入

**現象**
Prompt 產生的訊息中，參數位置是空的。

**原因**
`req.Params.Arguments` 的 key 與定義的 `Name` 不一致。

**解法**
```go
// 定義時的 Name
Arguments: []*mcp.PromptArgument{
    {Name: "code", ...},  // 注意是 "code" 不是 "Code"
}

// 取用時的 key
code := req.Params.Arguments["code"]  // 必須完全一致
```

---

### 踩坑 2：Required 參數為空

**現象**
設定 `Required: true`，但使用者沒填時沒有報錯。

**原因**
MCP SDK 不一定會強制驗證，需自行檢查。

**解法**
```go
func handler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    code, ok := req.Params.Arguments["code"]
    if !ok || code == "" {
        return nil, fmt.Errorf("code 參數為必填")
    }
    // promptText := fmt.Sprintf( ... (略)
}
```

---

### 踩坑 3：Prompt 內容格式混亂

**現象**
生成的 Prompt 文字格式不如預期。

**原因**
字串拼接時沒處理好換行和縮排。

**解法**
```go
// 使用 raw string literal
promptText := `請對以下程式碼進行 Code Review：

程式碼：
` + "```\n" + code + "\n```"

// 或使用 fmt.Sprintf 搭配 %s
promptText := fmt.Sprintf(`請 Review：

%s`, "```\n"+code+"\n```")
```

---

## ⑥ 驗收標準（DoD）

| 項目 | 驗證方式 | 預期結果 |
|------|---------|---------|
| code-review Prompt | 在 Claude 中選用 | 產生 Code Review 格式的提示詞 |
| summarize Prompt | 在 Claude 中選用 | 產生摘要格式的提示詞 |
| 參數帶入正確 | 輸入程式碼後檢查 | 程式碼正確嵌入提示詞 |
| 選填參數 | 不填 max_words | 使用預設值 100 |

**自我檢核清單**
- [ ] code-review Prompt 可正常使用
- [ ] summarize Prompt 可正常使用
- [ ] 必填參數驗證正確
- [ ] 選填參數有預設值

---

## 下一步

掌握三種 Capabilities 後，前往 [U06｜Combined Capabilities](U06-Combined-Capabilities.md) 學習如何整合到同一個 Server。
