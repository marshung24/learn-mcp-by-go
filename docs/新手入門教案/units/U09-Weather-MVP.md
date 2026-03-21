# U09｜Weather MVP

> 整合所有學習內容，完成一個完整的天氣查詢 MCP Server。
>
> 預估時數：120 min
> 前置依賴：U08

---

## ① 為什麼先教這個？

這是整個課程的總結單元。透過實作一個真實的天氣查詢服務，將所有學過的知識整合起來：
- **Tools**：查詢天氣預報、天氣警報
- **Resources**：提供使用說明
- **Prompts**：天氣報告模板
- **HTTP Transport**：支援遠端存取

完成這個 MVP 後，你就具備獨立開發 MCP Server 的能力。

---

## ② 對應素材

| 類型 | 路徑/連結 | 說明 |
|------|----------|------|
| 範例程式碼 | `examples/U09-weather-mvp/` | 完整專案 |
| API 文件 | https://api.weather.gov | 美國國家氣象局 API |
| API 說明 | https://www.weather.gov/documentation/services-web-api | API 使用說明 |

---

## ③ 核心觀念

### 1. NWS API 介紹

美國國家氣象局（National Weather Service）提供免費的天氣 API：
- 不需 API Key
- 需設定 User-Agent header
- 僅支援美國地區

### 2. API 端點

| 端點 | 用途 | 範例 |
|------|------|------|
| `/points/{lat},{lon}` | 取得格點資料（含 forecast URL） | `/points/39.7456,-104.9897` |
| `/gridpoints/{wfo}/{x},{y}/forecast` | 取得天氣預報 | 從 points 回應取得 |
| `/alerts/active/area/{state}` | 取得州的天氣警報 | `/alerts/active/area/CO` |

### 3. API 回應結構

**Points 回應**
```json
{
  "properties": {
    "forecast": "https://api.weather.gov/gridpoints/BOU/54,62/forecast",
    "forecastHourly": "...",
    "forecastGridData": "..."
  }
}
```

**Forecast 回應**
```json
{
  "properties": {
    "periods": [
      {
        "name": "Tonight",
        "temperature": 45,
        "temperatureUnit": "F",
        "windSpeed": "5 mph",
        "windDirection": "SW",
        "detailedForecast": "Partly cloudy..."
      }
    ]
  }
}
```

### 4. 錯誤處理策略

外部 API 呼叫需要完善的錯誤處理：
- 網路錯誤（逾時、連線失敗）
- API 錯誤（404、500）
- 資料解析錯誤

### 5. 專案結構

```
examples/U09-weather-mvp/
├── main.go              # 入口點
├── nws/
│   └── client.go        # NWS API Client
├── tools/
│   ├── forecast.go      # get_forecast Tool
│   └── alerts.go        # get_alerts Tool
├── resources/
│   └── help.go          # 使用說明
├── prompts/
│   └── report.go        # 天氣報告模板
└── config/
    └── config.go        # 設定管理
```

---

## ④ 動手做

### [必做] 4.1 建立專案結構

```bash
cd ~/mcp-test
mkdir -p examples/U09-weather-mvp
```

### [必做] 4.2 實作 NWS API Client

建立 `examples/U09-weather-mvp/main.go`（包含所有功能）：

```go
package main

// import / const 宣告
// import ("context" ... (略)

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
func NewNWSClient() *NWSClient {
	// return &NWSClient{ ... (略)
}

// makeRequest 發送 HTTP GET 請求到 NWS API
// 輸入: ctx, url / 輸出: response body bytes, error
func (c *NWSClient) makeRequest(ctx context.Context, url string) ([]byte, error) {
	// 建立請求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("建立請求失敗: %w", err)
	}

	// 設定必要的 headers
	// ⚠️ NWS API 要求設定 User-Agent，否則會回傳 403
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
// 標準 JSON struct 定義，對應 NWS API 回應格式
// 完整定義請參考 examples/U09-weather-mvp/main.go

// PointsResponse /points 端點的回應結構
type PointsResponse struct { // 結構定義 ... (略) }
// ForecastResponse /forecast 端點的回應結構
type ForecastResponse struct { // 結構定義 ... (略) }
// ForecastPeriod 單一預報時段
type ForecastPeriod struct { // 結構定義 ... (略) }
// AlertsResponse /alerts 端點的回應結構
type AlertsResponse struct { // 結構定義 ... (略) }
// AlertFeature 單一警報
type AlertFeature struct { // 結構定義 ... (略) }
// AlertProperties 警報詳細資訊
type AlertProperties struct { // 結構定義 ... (略) }

// ==================== 兩階段 API 呼叫（核心流程） ====================

// GetForecast 取得指定座標的天氣預報
// 輸入: ctx, lat, lon / 輸出: 預報時段列表, error
func (c *NWSClient) GetForecast(ctx context.Context, lat, lon float64) ([]ForecastPeriod, error) {
	// Step 1: 取得 points 資料（包含 forecast URL）
	// NWS API 非直覺流程：不能直接用座標查預報，必須先查 points 取得 forecast URL
	pointsURL := fmt.Sprintf("%s/points/%.4f,%.4f", NWSAPIBase, lat, lon)
	pointsData, err := c.makeRequest(ctx, pointsURL)
	if err != nil {
		return nil, err
	}

	// 解析 points 回應，取得 forecast URL
	var points PointsResponse
	if err := json.Unmarshal(pointsData, &points); err != nil {
		return nil, fmt.Errorf("解析 points 資料失敗: %w", err)
	}

	// 檢查 forecast URL 是否存在（座標不在美國時會為空）
	if points.Properties.Forecast == "" {
		return nil, fmt.Errorf("無法取得預報 URL（座標可能不在美國境內）")
	}

	// Step 2: 用取得的 forecast URL 查詢實際預報資料
	forecastData, err := c.makeRequest(ctx, points.Properties.Forecast)
	if err != nil {
		return nil, err
	}

	// 解析預報回應
	var forecast ForecastResponse
	if err := json.Unmarshal(forecastData, &forecast); err != nil {
		return nil, fmt.Errorf("解析預報資料失敗: %w", err)
	}

	return forecast.Properties.Periods, nil
}

// GetAlerts 取得指定州的天氣警報（單階段呼叫，流程同 makeRequest + Unmarshal）
func (c *NWSClient) GetAlerts(ctx context.Context, state string) ([]AlertFeature, error) {
	// 與 GetForecast 相同模式：makeRequest → Unmarshal ... (略) state = strings.ToUpper(state)
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
		// 錯誤回傳：設定 IsError: true，回傳友善訊息（不是 Go error）
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("取得預報失敗: %s", err.Error())},
			},
			IsError: true,
		}, nil, nil
	}

	// 檢查是否有資料
	// if len(periods) == 0 { ... (略)

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

// alertsHandler 處理天氣警報查詢（與 forecastHandler 相同模式）
func alertsHandler(ctx context.Context, req *mcp.CallToolRequest, input AlertsInput) (
	*mcp.CallToolResult, any, error,
) {
	// 呼叫 API → 錯誤處理 → 格式化輸出，模式同 forecastHandler ... (略) if len(input.State) != 2 {
}

// ==================== Resources ====================

// helpResourceHandler 提供使用說明
func helpResourceHandler(ctx context.Context, req *mcp.ReadResourceRequest) (
	*mcp.ReadResourceResult, error,
) {
	// 回傳 Markdown 格式的使用說明
	helpContent := `# Weather MCP Server 使用說明 ...`
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{URI: "file://weather-help", MIMEType: "text/markdown", Text: helpContent},
		},
	}, nil
}

// ==================== Prompts ====================

// weatherReportHandler 產生天氣報告模板
func weatherReportHandler(ctx context.Context, req *mcp.GetPromptRequest) (
	*mcp.GetPromptResult, error,
) {
	// 從 Arguments 取得 location，組合 prompt 文字 ... (略) location := ""
}

// ==================== Main ====================

func main() {
	// 設定 log 輸出到 stderr
	// log.SetOutput(os.Stderr) ... (略)

	// 建立 MCP Server（參考 U01）
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

	// ==================== 註冊 Resources ====================
	server.AddResource(&mcp.Resource{
		URI:         "file://weather-help",
		Name:        "Weather Help",
		Description: "Weather MCP Server 使用說明，包含城市座標和州代碼參考",
		MIMEType:    "text/markdown",
	}, helpResourceHandler)

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

	// 使用 Stdio Transport 執行（參考 U07）
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
```

### [必做] 4.3 使用 MCP Inspector 測試

> 操作方式見 [MCP Inspector 測試指引](../common/MCP-Inspector-測試指引.md)

```bash
npx @modelcontextprotocol/inspector go run ./examples/U09-weather-mvp/main.go
```

本單元測試重點：
- **Tools**：`get_forecast`（輸入美國座標如 34.0522, -118.2437）、`get_alerts`（輸入州代碼如 CA）
- **Resources**：讀取 `weather-help`
- **Prompts**：使用 `weather-report`

### [延伸] 4.4 設定 Claude Desktop 並測試

> 設定流程見 [Claude Desktop 設定指引](../common/Claude-Desktop-設定指引.md)

```bash
go build -o weather-mvp ./examples/U09-weather-mvp/
```

config.json 中的 Server 名稱為 `weather-mvp`，command 指向編譯產物的絕對路徑。

重啟後測試（注意：NWS API 僅支援美國地區）：
- 「查詢洛杉磯的天氣預報」→ 應呼叫 `get_forecast`
- 「加州有什麼天氣警報？」→ 應呼叫 `get_alerts`
- 使用 `weather-report` Prompt 產生完整天氣報告

### [延伸] 4.5 新增 HTTP 模式支援

參考 U06，為 Weather MVP 加入 HTTP Transport 支援，讓它可以同時支援 Stdio 和 HTTP 兩種模式：

```go
func main() {
	mode := os.Getenv("MODE")
	if mode == "http" {
		runHTTPServer()
	} else {
		runStdioServer()
	}
}

func runHTTPServer() {
	// HTTP 模式實作
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting HTTP server on :%s", port)
	// HTTP server 啟動邏輯 ... (略)
}

func runStdioServer() {
	// Stdio 模式（原本的程式碼）
	server.Run(context.Background(), &mcp.StdioTransport{})
}
```

### [延伸] 4.6 加入快取機制

```go
type Cache struct {
	data      map[string]cacheEntry
	mu        sync.RWMutex
	ttl       time.Duration
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.data[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}
```

### [延伸] 4.7 整合認證機制

參考 U08，為 Weather MVP 加入認證保護：

```go
func main() {
    // 建立 server 和註冊功能 ... (略)

    // 建立 HTTP Handler
    mcpHandler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
        return server
    }, nil)

    // 套用 Bearer Token 認證
    handler := BearerTokenMiddleware(tokenVerifier, "", mcpHandler)

    http.Handle("/mcp", handler)
    // 啟動 HTTP server ... (略)
}
```

### [延伸] 4.8 新增更多功能

- 支援更多 API 端點（小時預報、雷達資料）
- 整合其他天氣服務（OpenWeatherMap、WeatherAPI）
- 加入地點搜尋功能

---

## ⑤ 踩坑提示

### 踩坑 1：API 回傳 403

**現象**
NWS API 回傳 403 Forbidden。

**原因**
未設定 User-Agent header。

**解法**
```go
req.Header.Set("User-Agent", "YourApp/1.0 (your@email.com)")
```

---

### 踩坑 2：座標查無資料

**現象**
get_forecast 回傳「無法取得預報 URL」。

**原因**
座標不在美國境內，NWS API 僅支援美國。

**解法**
使用美國城市座標測試：
- 洛杉磯：34.0522, -118.2437
- 紐約：40.7128, -74.0060

---

### 踩坑 3：JSON 解析失敗

**現象**
「解析資料失敗」錯誤。

**原因**
API 回應格式可能有變動，或網路問題導致回應不完整。

**解法**
```go
// 加入詳細錯誤日誌
body, _ := io.ReadAll(resp.Body)
log.Printf("API Response: %s", string(body))
```

---

### 踩坑 4：逾時錯誤

**現象**
請求逾時，特別是首次請求。

**原因**
網路延遲或 API 回應慢。

**解法**
```go
httpClient: &http.Client{
    Timeout: 30 * time.Second,  // 增加逾時時間
}
```

---

## ⑥ 驗收標準（DoD）

| 項目 | 驗證方式 | 預期結果 |
|------|---------|---------|
| get_forecast | 查詢洛杉磯天氣 | 回傳 5 個時段的預報 |
| get_alerts | 查詢任一州警報 | 回傳警報列表或「無警報」 |
| weather-help | 讀取 Resource | 顯示完整使用說明 |
| weather-report | 使用 Prompt | 產生報告模板 |
| 錯誤處理 | 查詢非美國座標 | 回傳友善錯誤訊息 |

**自我檢核清單**
- [ ] get_forecast Tool 正常運作
- [ ] get_alerts Tool 正常運作
- [ ] weather-help Resource 內容完整
- [ ] weather-report Prompt 正常運作
- [ ] 錯誤情況有友善的訊息
- [ ] User-Agent 正確設定

---

## 恭喜完成！

你已經完成了 MCP Server 開發的完整學習路徑！

### 你學會了
1. MCP 協定的核心概念
2. Tools、Resources、Prompts 三種 capabilities
3. Stdio 與 HTTP 兩種 Transport
4. Basic Auth 與 Bearer Token 認證機制
5. 外部 API 整合
6. 錯誤處理與最佳實踐

### 下一步建議
1. 嘗試整合其他 API（GitHub、Slack、資料庫）
2. 學習 MCP Client 開發
3. 探索 OAuth 2.0 完整流程（需外部授權伺服器）
4. 將 Server 部署到雲端（加上 HTTPS）

### 參考資源
- [MCP 官方文件](https://modelcontextprotocol.io/docs)
- [Go SDK GitHub](https://github.com/modelcontextprotocol/go-sdk)
- [MCP 範例專案](https://github.com/modelcontextprotocol/quickstart-resources)
