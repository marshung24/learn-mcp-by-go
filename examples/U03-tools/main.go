// Package main 實作 Tool Capability 範例
// 這是 U03 Tool Capability 的範例程式碼，展示如何定義和使用 MCP Tools
package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

// GreetInput 定義 greet Tool 的輸入參數
type GreetInput struct {
	Name string `json:"name" jsonschema:"要問候的人的名字"`
}

// CalculateInput 定義 calculate Tool 的輸入參數
// 支援四則運算
type CalculateInput struct {
	A        float64 `json:"a" jsonschema:"第一個數字"`
	B        float64 `json:"b" jsonschema:"第二個數字"`
	Operator string  `json:"operator" jsonschema:"運算符號(+, -, *, /)"`
}

// ==================== Tool Handlers ====================

// addHandler 處理 add Tool 的請求
// 計算兩個數字的和
func addHandler(ctx context.Context, req *mcp.CallToolRequest, input AddInput) (
	*mcp.CallToolResult, any, error,
) {
	// 計算兩數之和
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
			IsError: true,
		}, nil, nil
	}

	// 輸入驗證：檢查名字長度
	if len(input.Name) > 100 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "錯誤：名字太長了（最多 100 字元）"},
			},
			IsError: true,
		}, nil, nil
	}

	// 組合問候訊息並回傳
	greeting := fmt.Sprintf("Hello, %s! 歡迎使用 MCP Tools!", input.Name)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: greeting,
			},
		},
	}, nil, nil
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
	case "-":
		result = input.A - input.B
	case "*":
		result = input.A * input.B
	case "/":
		// 除法需要檢查除數是否為零
		if input.B == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "錯誤：除數不能為零"},
				},
				IsError: true,
			}, nil, nil
		}
		result = input.A / input.B
	default:
		errMsg = fmt.Sprintf("錯誤：不支援的運算符號「%s」，請使用 +、-、*、/", input.Operator)
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

// ==================== Main ====================

func main() {
	// 設定 log 輸出到 stderr
	log.SetOutput(os.Stderr)
	log.Println("Starting tools-demo server...")

	// 建立 MCP Server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "tools-demo",
		Version: "1.0.0",
	}, nil)

	// 註冊 add Tool
	// Name: Tool 的唯一識別名稱
	// Description: 說明 Tool 的用途，AI 會根據這個描述決定何時使用
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "計算兩個數字的和。輸入兩個數字 a 和 b，回傳 a + b 的結果。",
	}, addHandler)
	log.Println("Tool 'add' registered")

	// 註冊 greet Tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "向指定的人打招呼。輸入名字，回傳問候訊息。",
	}, greetHandler)
	log.Println("Tool 'greet' registered")

	// 註冊 calculate Tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "calculate",
		Description: "四則運算計算機。輸入兩個數字和運算符號（+、-、*、/），回傳計算結果。",
	}, calculateHandler)
	log.Println("Tool 'calculate' registered")

	log.Println("All tools registered, waiting for connections...")

	// 使用 Stdio Transport 執行
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
