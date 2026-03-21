// Package main 實作 Weather MVP 範例
// 這是 U09 Weather MVP 的範例程式碼
// 整合所有學習內容，完成一個完整的天氣查詢 MCP Server
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ==================== 常數定義 ====================

const (
	// NWSAPIBase NWS API 基礎 URL
	NWSAPIBase = "https://api.weather.gov"
	// UserAgent 請求時的 User-Agent header（NWS API 必須設定）
	UserAgent = "MCP-Weather-Demo/1.0 (learning@example.com)"
)

// ==================== NWS API Client ====================

// NWSClient 封裝 NWS API 呼叫邏輯
type NWSClient struct {
	httpClient *http.Client
}

// NewNWSClient 建立新的 NWS API Client
// 輸出: 包含 30 秒逾時設定的 NWSClient 實例
func NewNWSClient() *NWSClient {
	// 設定 HTTP Client，含逾時保護
	return &NWSClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest 發送 HTTP GET 請求到 NWS API
func (c *NWSClient) makeRequest(ctx context.Context, url string) ([]byte, error) {
	// 建立請求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("建立請求失敗: %w", err)
	}

	// 設定必要的 headers
	// NWS API 要求設定 User-Agent，否則會回傳 403
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/geo+json")

	// 發送請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("請求失敗: %w", err)
	}
	defer resp.Body.Close()

	// 檢查回應狀態碼
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 錯誤 %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// ==================== API Response 結構 ====================

// PointsResponse /points 端點的回應結構
type PointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// ForecastResponse /forecast 端點的回應結構
type ForecastResponse struct {
	Properties struct {
		Periods []ForecastPeriod `json:"periods"`
	} `json:"properties"`
}

// ForecastPeriod 單一預報時段
type ForecastPeriod struct {
	Name             string `json:"name"`
	Temperature      int    `json:"temperature"`
	TemperatureUnit  string `json:"temperatureUnit"`
	WindSpeed        string `json:"windSpeed"`
	WindDirection    string `json:"windDirection"`
	DetailedForecast string `json:"detailedForecast"`
}

// AlertsResponse /alerts 端點的回應結構
type AlertsResponse struct {
	Features []AlertFeature `json:"features"`
}

// AlertFeature 單一警報
type AlertFeature struct {
	Properties AlertProperties `json:"properties"`
}

// AlertProperties 警報詳細資訊
type AlertProperties struct {
	Event       string `json:"event"`
	AreaDesc    string `json:"areaDesc"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Instruction string `json:"instruction"`
}

// ==================== 兩階段 API 呼叫（核心流程） ====================

// GetForecast 取得指定座標的天氣預報
func (c *NWSClient) GetForecast(ctx context.Context, lat, lon float64) ([]ForecastPeriod, error) {
	// Step 1: 取得 points 資料（包含 forecast URL）
	pointsURL := fmt.Sprintf("%s/points/%.4f,%.4f", NWSAPIBase, lat, lon)
	log.Printf("Fetching points: %s", pointsURL)

	pointsData, err := c.makeRequest(ctx, pointsURL)
	if err != nil {
		return nil, err
	}

	var points PointsResponse
	if err := json.Unmarshal(pointsData, &points); err != nil {
		return nil, fmt.Errorf("解析 points 資料失敗: %w", err)
	}

	if points.Properties.Forecast == "" {
		return nil, fmt.Errorf("無法取得預報 URL（座標可能不在美國境內）")
	}

	// Step 2: 取得 forecast 資料
	log.Printf("Fetching forecast: %s", points.Properties.Forecast)
	forecastData, err := c.makeRequest(ctx, points.Properties.Forecast)
	if err != nil {
		return nil, err
	}

	var forecast ForecastResponse
	if err := json.Unmarshal(forecastData, &forecast); err != nil {
		return nil, fmt.Errorf("解析預報資料失敗: %w", err)
	}

	return forecast.Properties.Periods, nil
}

// GetAlerts 取得指定州的天氣警報
func (c *NWSClient) GetAlerts(ctx context.Context, state string) ([]AlertFeature, error) {
	// 轉換為大寫
	state = strings.ToUpper(state)
	alertsURL := fmt.Sprintf("%s/alerts/active/area/%s", NWSAPIBase, state)
	log.Printf("Fetching alerts: %s", alertsURL)

	alertsData, err := c.makeRequest(ctx, alertsURL)
	if err != nil {
		return nil, err
	}

	var alerts AlertsResponse
	if err := json.Unmarshal(alertsData, &alerts); err != nil {
		return nil, fmt.Errorf("解析警報資料失敗: %w", err)
	}

	return alerts.Features, nil
}

// ==================== Tools ====================

// 全域 NWS Client 實例
var nwsClient = NewNWSClient()

// ForecastInput get_forecast Tool 的輸入參數
type ForecastInput struct {
	Latitude  float64 `json:"latitude" jsonschema:"緯度（僅支援美國地區，如 34.0522）"`
	Longitude float64 `json:"longitude" jsonschema:"經度（僅支援美國地區，如 -118.2437）"`
}

// forecastHandler 處理天氣預報查詢
func forecastHandler(ctx context.Context, req *mcp.CallToolRequest, input ForecastInput) (
	*mcp.CallToolResult, any, error,
) {
	// 呼叫 NWS API
	periods, err := nwsClient.GetForecast(ctx, input.Latitude, input.Longitude)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("取得預報失敗: %s", err.Error())},
			},
			IsError: true,
		}, nil, nil
	}

	// 檢查是否有資料
	if len(periods) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "此地區沒有可用的預報資料"},
			},
		}, nil, nil
	}

	// 格式化輸出（只顯示前 5 個時段）
	var forecasts []string
	for i, p := range periods {
		if i >= 5 {
			break
		}
		forecast := fmt.Sprintf(`**%s**
溫度: %d°%s
風向風速: %s %s
預報: %s`, p.Name, p.Temperature, p.TemperatureUnit,
			p.WindSpeed, p.WindDirection, p.DetailedForecast)
		forecasts = append(forecasts, forecast)
	}

	result := fmt.Sprintf("## 天氣預報 (%.4f, %.4f)\n\n%s",
		input.Latitude, input.Longitude, strings.Join(forecasts, "\n\n---\n\n"))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

// AlertsInput get_alerts Tool 的輸入參數
type AlertsInput struct {
	State string `json:"state" jsonschema:"美國州代碼（兩個字母，如 CA、NY、TX）"`
}

// alertsHandler 處理天氣警報查詢
func alertsHandler(ctx context.Context, req *mcp.CallToolRequest, input AlertsInput) (
	*mcp.CallToolResult, any, error,
) {
	// 驗證輸入
	if len(input.State) != 2 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "錯誤：請提供兩個字母的美國州代碼（如 CA、NY、TX）"},
			},
			IsError: true,
		}, nil, nil
	}

	// 呼叫 NWS API
	alerts, err := nwsClient.GetAlerts(ctx, input.State)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("取得警報失敗: %s", err.Error())},
			},
			IsError: true,
		}, nil, nil
	}

	// 檢查是否有警報
	if len(alerts) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("%s 目前沒有有效的天氣警報", strings.ToUpper(input.State))},
			},
		}, nil, nil
	}

	// 格式化輸出
	var alertTexts []string
	for _, a := range alerts {
		p := a.Properties
		instruction := p.Instruction
		if instruction == "" {
			instruction = "無特定指示"
		}

		// 截斷過長的描述
		description := p.Description
		if len(description) > 500 {
			description = description[:500] + "..."
		}

		alertText := fmt.Sprintf(`**%s**
地區: %s
嚴重程度: %s
說明: %s
指示: %s`, p.Event, p.AreaDesc, p.Severity, description, instruction)
		alertTexts = append(alertTexts, alertText)
	}

	result := fmt.Sprintf("## %s 天氣警報（共 %d 個）\n\n%s",
		strings.ToUpper(input.State), len(alerts), strings.Join(alertTexts, "\n\n---\n\n"))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

// ==================== Resources ====================

// helpResourceHandler 提供使用說明
func helpResourceHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	helpContent := `# Weather MCP Server 使用說明

## 概述
這是一個整合美國國家氣象局 (NWS) API 的 MCP Server，提供天氣預報和天氣警報查詢功能。

## 重要限制
⚠️ **僅支援美國地區**：NWS API 只提供美國境內的氣象資料。

## 可用功能

### Tools

#### get_forecast
查詢指定經緯度的天氣預報。

**參數：**
- latitude: 緯度（如 34.0522）
- longitude: 經度（如 -118.2437）

**範例城市座標：**
- 洛杉磯 (Los Angeles): 34.0522, -118.2437
- 紐約 (New York): 40.7128, -74.0060
- 芝加哥 (Chicago): 41.8781, -87.6298
- 休士頓 (Houston): 29.7604, -95.3698
- 丹佛 (Denver): 39.7392, -104.9903
- 西雅圖 (Seattle): 47.6062, -122.3321
- 邁阿密 (Miami): 25.7617, -80.1918

#### get_alerts
查詢指定州的天氣警報。

**參數：**
- state: 美國州代碼（兩個字母）

**常見州代碼：**
- CA: 加州 (California)
- NY: 紐約州 (New York)
- TX: 德州 (Texas)
- FL: 佛羅里達州 (Florida)
- CO: 科羅拉多州 (Colorado)
- WA: 華盛頓州 (Washington)

### Prompts

#### weather-report
產生天氣報告模板，用於撰寫完整的天氣報告。

## 使用範例
- 「查詢洛杉磯的天氣」
- 「紐約今天天氣如何？」（座標：40.7128, -74.0060）
- 「德州有什麼天氣警報？」
- 「加州現在有暴風雨警報嗎？」

## 注意事項
- 溫度單位為華氏度 (°F)
- 預報通常顯示未來 5 個時段
- 警報資料為即時更新
`
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "file://weather-help",
				MIMEType: "text/markdown",
				Text:     helpContent,
			},
		},
	}, nil
}

// ==================== Prompts ====================

// weatherReportHandler 產生天氣報告模板
func weatherReportHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 從 Arguments 取得 location 參數
	location := ""
	if req.Params.Arguments != nil {
		if l, ok := req.Params.Arguments["location"]; ok {
			location = l
		}
	}

	// 未指定時使用預設提示文字
	if location == "" {
		location = "（請指定地點）"
	}

	promptText := fmt.Sprintf(`請根據以下資訊撰寫一份天氣報告：

地點：%s

請包含：
1. 當前天氣概況
2. 未來幾天的天氣趨勢
3. 戶外活動建議
4. 穿著建議
5. 任何需要注意的天氣警報

請使用友善、易懂的語言撰寫，並將華氏溫度轉換為攝氏溫度供參考。`, location)

	return &mcp.GetPromptResult{
		Description: "天氣報告撰寫模板",
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
	// 設定 log 輸出到 stderr
	log.SetOutput(os.Stderr)
	log.Println("Starting Weather MCP Server...")
	log.Println("Note: This server only supports US locations (NWS API)")

	// 建立 MCP Server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "weather-mvp",
		Version: "1.0.0",
	}, nil)

	// ==================== 註冊 Tools ====================
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_forecast",
		Description: "查詢指定經緯度的天氣預報（僅支援美國地區）。輸入緯度和經度，回傳未來天氣預報。",
	}, forecastHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_alerts",
		Description: "查詢指定美國州的天氣警報。輸入州代碼（如 CA、NY、TX），回傳目前有效的天氣警報。",
	}, alertsHandler)
	log.Println("Tools registered: get_forecast, get_alerts")

	// ==================== 註冊 Resources ====================
	server.AddResource(&mcp.Resource{
		URI:         "file://weather-help",
		Name:        "Weather Help",
		Description: "Weather MCP Server 使用說明，包含城市座標和州代碼參考",
		MIMEType:    "text/markdown",
	}, helpResourceHandler)
	log.Println("Resource registered: weather-help")

	// ==================== 註冊 Prompts ====================
	server.AddPrompt(&mcp.Prompt{
		Name:        "weather-report",
		Description: "天氣報告撰寫模板，協助撰寫完整的天氣報告",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "location",
				Description: "要撰寫報告的地點",
				Required:    true,
			},
		},
	}, weatherReportHandler)
	log.Println("Prompt registered: weather-report")

	log.Println("Weather MCP Server ready!")
	log.Println("Waiting for connections...")

	// 使用 Stdio Transport 執行
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
