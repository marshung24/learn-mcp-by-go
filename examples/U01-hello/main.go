// Package main 實作最小可執行的 MCP Server
// 這是 U01 Hello MCP Server 的範例程式碼
package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// 設定 log 輸出到 stderr
	// 重要：Stdio 模式下，stdout 專門用於 JSON-RPC 通訊，禁止使用 fmt.Println()
	log.SetOutput(os.Stderr)
	log.Println("Starting hello-mcp server...")

	// 建立 MCP Server
	// Implementation 包含 Server 的基本資訊：名稱和版本
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "hello-mcp",
		Version: "1.0.0",
	}, nil)

	log.Println("Server created, waiting for connections...")

	// 使用 Stdio Transport 執行
	// Server 會透過標準輸入/輸出與 Client 通訊
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
