// Package main 實作 Resource Capability 範例
// 這是 U04 Resource Capability 的範例程式碼，展示如何提供資料給 AI 讀取
package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ==================== Resource Handlers ====================

// readmeHandler 處理靜態 README Resource 的讀取請求
// 靜態 Resource：內容固定不變
func readmeHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 定義靜態內容
	content := `# MCP Resources Demo

這是一個示範 Resource 功能的 MCP Server。

## 功能
- 提供靜態文件（README）
- 提供動態資料（當前時間）
- 讀取外部檔案（設定檔）
- 列出目錄內容

## 使用方式
透過 Claude 或 MCP Inspector 讀取這些 Resource 即可。

## 可用 Resources
- file://readme - 此說明文件
- data://current-time - 當前系統時間
- file://config - 專案設定檔
- data://files - 目錄檔案列表
`
	// 回傳 ReadResourceResult，包含 TextResourceContents
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "file://readme",
				MIMEType: "text/markdown",
				Text:     content,
			},
		},
	}, nil
}

// timeHandler 處理動態時間 Resource 的讀取請求
// 動態 Resource：每次讀取時重新計算內容
func timeHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 取得當前時間
	now := time.Now()

	// 建立時間資訊結構
	content := map[string]interface{}{
		"timestamp":   now.Unix(),
		"formatted":   now.Format("2006-01-02 15:04:05"),
		"timezone":    now.Location().String(),
		"weekday":     now.Weekday().String(),
		"day_of_year": now.YearDay(),
	}

	// 轉換為 JSON 格式
	jsonBytes, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return nil, err
	}

	// 回傳動態內容——每次呼叫都會取得最新時間
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "data://current-time",
				MIMEType: "application/json",
				Text:     string(jsonBytes),
			},
		},
	}, nil
}

// configHandler 處理設定檔 Resource 的讀取請求
// 從外部檔案讀取內容
func configHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 取得執行檔所在目錄
	// 注意：實際部署時可能需要調整路徑處理方式
	exePath, err := os.Executable()
	if err != nil {
		// 如果無法取得執行檔路徑，嘗試使用當前目錄
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)
	filePath := filepath.Join(exeDir, "data", "sample.json")

	// 如果找不到相對於執行檔的路徑，嘗試當前目錄
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 嘗試從當前工作目錄讀取
		wd, _ := os.Getwd()
		filePath = filepath.Join(wd, "examples", "U04-resources", "data", "sample.json")
	}

	// 讀取檔案內容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 讀取失敗時回傳錯誤訊息，但不中斷程式
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "file://config",
					MIMEType: "text/plain",
					Text:     "無法讀取設定檔：" + err.Error() + "\n\n嘗試的路徑：" + filePath,
				},
			},
		}, nil
	}

	// 成功讀取，回傳檔案內容
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "file://config",
				MIMEType: "application/json",
				Text:     string(content),
			},
		},
	}, nil
}

// filesHandler 處理目錄列表 Resource 的讀取請求
// 列出當前目錄的檔案資訊
func filesHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 讀取當前目錄內容
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	// 建立檔案資訊列表
	var files []map[string]interface{}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, map[string]interface{}{
			"name":    entry.Name(),
			"isDir":   entry.IsDir(),
			"size":    info.Size(),
			"modTime": info.ModTime().Format(time.RFC3339),
		})
	}

	// 轉換為 JSON 格式
	jsonBytes, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return nil, err
	}

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

// ==================== Main ====================

func main() {
	// 設定 log 輸出到 stderr
	log.SetOutput(os.Stderr)
	log.Println("Starting resources-demo server...")

	// 建立 MCP Server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "resources-demo",
		Version: "1.0.0",
	}, nil)

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
	log.Println("Resource 'readme' registered (static)")

	// 註冊動態 Resource：當前時間
	// 使用 data:// 表示資料類資源
	server.AddResource(&mcp.Resource{
		URI:         "data://current-time",
		Name:        "Current Time",
		Description: "當前系統時間（動態更新）",
		MIMEType:    "application/json",
	}, timeHandler)
	log.Println("Resource 'current-time' registered (dynamic)")

	// 註冊檔案 Resource：設定檔
	server.AddResource(&mcp.Resource{
		URI:         "file://config",
		Name:        "Config File",
		Description: "專案設定檔（sample.json）",
		MIMEType:    "application/json",
	}, configHandler)
	log.Println("Resource 'config' registered (file)")

	// 註冊動態 Resource：目錄列表
	server.AddResource(&mcp.Resource{
		URI:         "data://files",
		Name:        "Directory Listing",
		Description: "當前目錄的檔案列表",
		MIMEType:    "application/json",
	}, filesHandler)
	log.Println("Resource 'files' registered (dynamic)")

	log.Println("All resources registered, waiting for connections...")

	// 使用 Stdio Transport 執行
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
