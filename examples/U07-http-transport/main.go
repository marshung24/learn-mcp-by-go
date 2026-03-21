// Package main 實作 HTTP Transport 範例
// 這是 U07 HTTP Transport 的範例程式碼
// 展示如何使用 SDK 提供的 StreamableHTTPHandler 實作 HTTP 模式
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ==================== Tools ====================

// AddInput 定義 add Tool 的輸入參數
// 用於接收兩個浮點數進行加法運算
type AddInput struct {
	A float64 `json:"a" jsonschema:"第一個數字"`
	B float64 `json:"b" jsonschema:"第二個數字"`
}

// addHandler 處理加法運算
// 輸入: AddInput（兩個浮點數）
// 輸出: CallToolResult（文字格式的運算結果）
func addHandler(ctx context.Context, req *mcp.CallToolRequest, input AddInput) (
	*mcp.CallToolResult, any, error,
) {
	// 執行加法運算並格式化結果
	result := input.A + input.B
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
	// 設定 log 輸出到 stderr
	log.SetOutput(os.Stderr)

	// 讀取環境變數設定 port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 建立 MCP Server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "http-demo",
		Version: "1.0.0",
	}, nil)

	// 註冊 Tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "計算兩數之和",
	}, addHandler)

	log.Println("Tool 'add' registered")

	// 使用官方的 StreamableHTTPHandler
	// 這個 handler 會自動處理 JSON-RPC、session 管理等
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		// 為每個請求返回 server 實例
		// 實際應用中可以根據 req (如 headers) 返回不同的 server
		return server
	}, nil)

	// 設定路由
	http.Handle("/mcp", handler)

	// Health check 端點
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 首頁說明端點 - 顯示 Server 資訊與使用方式
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 檢查路徑是否為根路徑
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// 設定回應 Content-Type
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// 輸出 Server 資訊與使用說明
		fmt.Fprintf(w, `HTTP MCP Server Demo

Endpoints:
  POST /mcp     - MCP endpoint (使用 StreamableHTTPHandler)
  GET  /health  - Health check

使用 MCP Inspector 測試:
  1. 啟動 Inspector: npx @modelcontextprotocol/inspector
  2. 選擇 HTTP 模式
  3. 輸入 URL: http://localhost:%s/mcp
  4. 測試 add 工具

或使用 Claude Code:
  claude mcp add -t http http-demo http://localhost:%s/mcp
`, port, port)
	})

	// 輸出啟動資訊
	log.Printf("HTTP MCP Server listening on :%s", port)
	log.Printf("使用 StreamableHTTPHandler（官方實作）")
	log.Printf("Endpoints: POST /mcp, GET /health")
	log.Println("Press Ctrl+C to stop")

	// 啟動 HTTP Server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
