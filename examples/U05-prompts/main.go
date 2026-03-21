// Package main 實作 Prompt Capability 範例
// 這是 U05 Prompt Capability 的範例程式碼，展示如何建立可重用的提示詞模板
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ==================== Prompt Handlers ====================

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
	if code == "" {
		code = "（請在此貼上程式碼）"
	}

	// 組合提示詞文字：將審查項目與程式碼組合成完整 Prompt
	promptText := fmt.Sprintf(`請對以下程式碼進行 Code Review，包含：

1. **潛在 Bug**：找出可能的邏輯錯誤或例外情況
2. **程式碼風格**：檢查命名、縮排、註解
3. **效能考量**：指出可能的效能問題
4. **安全性**：檢查潛在的安全漏洞
5. **改進建議**：提供具體的改進方案

程式碼：
%s`, "```\n"+code+"\n```")

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

// summarizeHandler 處理 summarize Prompt 的請求
// 產生文字摘要的提示詞模板
func summarizeHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 取得參數
	text := ""
	maxWords := "100"

	if req.Params.Arguments != nil {
		if t, ok := req.Params.Arguments["text"]; ok {
			text = t
		}
		if m, ok := req.Params.Arguments["max_words"]; ok && m != "" {
			maxWords = m
		}
	}

	// 驗證必填參數
	if text == "" {
		text = "（請在此貼上要摘要的文字）"
	}

	// 組合提示詞文字：將摘要規則與原文組合成完整 Prompt
	promptText := fmt.Sprintf(`請將以下文字摘要為 %s 字以內的重點：

1. 保留關鍵資訊
2. 使用條列式呈現
3. 語言簡潔明瞭

原文：
%s`, maxWords, text)

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

// translateHandler 處理 translate Prompt 的請求
// 產生翻譯的提示詞模板
func translateHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 取得參數
	text := ""
	targetLang := "繁體中文"

	if req.Params.Arguments != nil {
		if t, ok := req.Params.Arguments["text"]; ok {
			text = t
		}
		if l, ok := req.Params.Arguments["target_language"]; ok && l != "" {
			targetLang = l
		}
	}

	// 驗證必填參數
	if text == "" {
		text = "（請在此貼上要翻譯的文字）"
	}

	// 組合提示詞文字：將翻譯要求與原文組合成完整 Prompt
	promptText := fmt.Sprintf(`請將以下文字翻譯為 %s：

翻譯要求：
1. 保持原意
2. 語句通順
3. 符合目標語言習慣

原文：
%s`, targetLang, text)

	// 回傳 GetPromptResult：包含 Description 與 PromptMessage 列表
	return &mcp.GetPromptResult{
		Description: "翻譯提示詞",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

// debugHandler 處理 debug Prompt 的請求
// 產生錯誤分析的提示詞模板
func debugHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 取得參數
	errorMsg := ""
	stackTrace := ""
	codeContext := ""

	if req.Params.Arguments != nil {
		if e, ok := req.Params.Arguments["error"]; ok {
			errorMsg = e
		}
		if s, ok := req.Params.Arguments["stack_trace"]; ok {
			stackTrace = s
		}
		if c, ok := req.Params.Arguments["context"]; ok {
			codeContext = c
		}
	}

	// 設定預設值
	if errorMsg == "" {
		errorMsg = "（請貼上錯誤訊息）"
	}
	if stackTrace == "" {
		stackTrace = "（選填：貼上 Stack Trace）"
	}
	if codeContext == "" {
		codeContext = "（選填：提供相關程式碼或上下文）"
	}

	// 組合提示詞文字：將錯誤資訊與分析要求組合成完整 Prompt
	promptText := fmt.Sprintf(`請幫我分析以下錯誤：

## 錯誤訊息
%s

## Stack Trace
%s

## 相關程式碼/上下文
%s

請提供：
1. 錯誤原因分析
2. 可能的解決方案（依優先順序列出）
3. 預防措施建議`, errorMsg, stackTrace, codeContext)

	// 回傳 GetPromptResult：包含 Description 與 PromptMessage 列表
	return &mcp.GetPromptResult{
		Description: "錯誤分析提示詞",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

// ==================== Main ====================

// main 啟動 MCP Server 並註冊所有 Prompt
// 輸出: 透過 Stdio Transport 提供 Prompt 服務
func main() {
	// 設定 log 輸出到 stderr
	log.SetOutput(os.Stderr)
	log.Println("Starting prompts-demo server...")

	// 建立 MCP Server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "prompts-demo",
		Version: "1.0.0",
	}, nil)

	// 註冊 code-review Prompt
	// Name: Prompt 的唯一名稱
	// Description: 說明此 Prompt 的用途
	// Arguments: 可帶入的參數列表
	server.AddPrompt(&mcp.Prompt{
		Name:        "code-review",
		Description: "對程式碼進行專業的 Code Review",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "code",
				Description: "要 Review 的程式碼",
				Required:    true,
			},
		},
	}, codeReviewHandler)
	log.Println("Prompt 'code-review' registered")

	// 註冊 summarize Prompt
	server.AddPrompt(&mcp.Prompt{
		Name:        "summarize",
		Description: "將文字摘要為重點",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "text",
				Description: "要摘要的文字",
				Required:    true,
			},
			{
				Name:        "max_words",
				Description: "摘要字數上限（預設 100）",
				Required:    false,
			},
		},
	}, summarizeHandler)
	log.Println("Prompt 'summarize' registered")

	// 註冊 translate Prompt
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
	log.Println("Prompt 'translate' registered")

	// 註冊 debug Prompt
	server.AddPrompt(&mcp.Prompt{
		Name:        "debug",
		Description: "分析錯誤訊息並提供解決方案",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "error",
				Description: "錯誤訊息",
				Required:    true,
			},
			{
				Name:        "stack_trace",
				Description: "Stack Trace（選填）",
				Required:    false,
			},
			{
				Name:        "context",
				Description: "相關程式碼或上下文（選填）",
				Required:    false,
			},
		},
	}, debugHandler)
	log.Println("Prompt 'debug' registered")

	log.Println("All prompts registered, waiting for connections...")

	// 使用 Stdio Transport 執行
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
