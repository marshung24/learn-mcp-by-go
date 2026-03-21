// Package main 實作 Combined Capabilities 範例
// 這是 U06 Combined Capabilities 的範例程式碼
// 整合 Tools、Resources、Prompts 三種功能到同一個 MCP Server
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ==================== 全域變數 ====================

// startTime 記錄 Server 啟動時間，用於計算 uptime
var startTime time.Time

// ==================== Tools ====================

// CalculateInput 定義 calculate Tool 的輸入參數
type CalculateInput struct {
	A        float64 `json:"a" jsonschema:"第一個數字"`
	B        float64 `json:"b" jsonschema:"第二個數字"`
	Operator string  `json:"operator" jsonschema:"運算符號(+, -, *, /)"`
}

// calculateHandler 處理四則運算
// 輸入: CalculateInput（兩個數字與運算符號）
// 輸出: 運算結果文字，或錯誤訊息
func calculateHandler(ctx context.Context, req *mcp.CallToolRequest, input CalculateInput) (
	*mcp.CallToolResult, any, error,
) {
	// 宣告結果與錯誤訊息變數
	var result float64
	var errMsg string

	// 依運算符號執行對應運算
	switch input.Operator {
	case "+":
		result = input.A + input.B
	case "-":
		result = input.A - input.B
	case "*":
		result = input.A * input.B
	case "/":
		if input.B == 0 {
			errMsg = "錯誤：除數不能為零"
		} else {
			result = input.A / input.B
		}
	default:
		errMsg = fmt.Sprintf("不支援的運算符號：%s", input.Operator)
	}

	// 若有錯誤，回傳錯誤結果
	if errMsg != "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: errMsg},
			},
			IsError: true,
		}, nil, nil
	}

	// 回傳運算結果
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%v %s %v = %v", input.A, input.Operator, input.B, result),
			},
		},
	}, nil, nil
}

// ==================== Resources ====================

// helpResourceHandler 提供使用說明
// 輸入: ReadResourceRequest
// 輸出: Markdown 格式的使用說明文件
func helpResourceHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 定義使用說明內容
	helpContent := `# Combined MCP Server 使用說明

## 可用功能

### Tools
- **calculate**: 四則運算計算機
  - 參數：a (數字), b (數字), operator (+, -, *, /)
  - 範例：「計算 10 除以 3」

### Resources
- **help**: 此說明文件
- **status**: 系統狀態資訊

### Prompts
- **math-problem**: 數學問題解題模板
- **explain-code**: 程式碼解釋模板

## 使用方式
直接向 AI 描述你的需求即可。

## 關於這個 Server
這是一個整合了 Tools、Resources、Prompts 三種 capabilities 的示範 Server。
`
	// 回傳 Markdown 格式的說明文件
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "file://help",
				MIMEType: "text/markdown",
				Text:     helpContent,
			},
		},
	}, nil
}

// statusResourceHandler 提供系統狀態
// 輸入: ReadResourceRequest
// 輸出: JSON 格式的系統狀態資訊（含 uptime、capabilities 等）
func statusResourceHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 計算執行時間
	uptime := time.Since(startTime)

	// 組裝狀態資訊 map
	status := map[string]interface{}{
		"server":    "combined-demo",
		"version":   "1.0.0",
		"uptime":    uptime.String(),
		"timestamp": time.Now().Format(time.RFC3339),
		"capabilities": map[string]bool{
			"tools":     true,
			"resources": true,
			"prompts":   true,
		},
		"registeredTools":     []string{"calculate"},
		"registeredResources": []string{"help", "status"},
		"registeredPrompts":   []string{"math-problem", "explain-code"},
	}

	// map 結構固定，序列化不會失敗，故忽略 error
	jsonBytes, _ := json.MarshalIndent(status, "", "  ")

	// 回傳 JSON 格式的狀態資訊
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "data://status",
				MIMEType: "application/json",
				Text:     string(jsonBytes),
			},
		},
	}, nil
}

// ==================== Prompts ====================

// mathProblemHandler 產生數學問題解題模板
// 輸入: GetPromptRequest（arguments: problem）
// 輸出: 包含解題步驟引導的 PromptMessage
func mathProblemHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 從 arguments 取得數學問題
	problem := ""
	if req.Params.Arguments != nil {
		if p, ok := req.Params.Arguments["problem"]; ok {
			problem = p
		}
	}

	// 若未提供問題，使用預設提示文字
	if problem == "" {
		problem = "（請輸入數學問題）"
	}

	// 組裝解題模板
	promptText := fmt.Sprintf(`請解答以下數學問題，並說明解題步驟：

問題：%s

請提供：
1. 解題思路
2. 詳細步驟
3. 最終答案
4. 驗算過程（如適用）`, problem)

	// 回傳組裝好的 Prompt 結果
	return &mcp.GetPromptResult{
		Description: "數學問題解題",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

// explainCodeHandler 產生程式碼解釋模板
// 輸入: GetPromptRequest（arguments: code, language）
// 輸出: 包含程式碼解釋引導的 PromptMessage
func explainCodeHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 從 arguments 取得程式碼與語言
	code := ""
	language := "自動偵測"

	if req.Params.Arguments != nil {
		if c, ok := req.Params.Arguments["code"]; ok {
			code = c
		}
		if l, ok := req.Params.Arguments["language"]; ok && l != "" {
			language = l
		}
	}

	// 若未提供程式碼，使用預設提示文字
	if code == "" {
		code = "（請在此貼上程式碼）"
	}

	// 組裝程式碼解釋模板
	promptText := fmt.Sprintf(`請解釋以下 %s 程式碼：

%s

請說明：
1. 程式碼的主要功能
2. 逐行或逐區塊解釋
3. 使用的程式技巧或模式
4. 可能的改進建議`, language, "```\n"+code+"\n```")

	// 回傳組裝好的 Prompt 結果
	return &mcp.GetPromptResult{
		Description: "程式碼解釋",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
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
