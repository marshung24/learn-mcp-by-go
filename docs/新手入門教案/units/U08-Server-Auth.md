# U08｜Server 認證

> 為 HTTP MCP Server 加入認證機制，保護 API 存取安全。
>
> 預估時數：90 min
> 前置依賴：U07

---

## ① 為什麼先教這個？

HTTP 模式讓 Server 可透過網路存取，但也帶來安全風險：
- 任何人都能呼叫你的 API
- 無法區分合法與非法使用者
- 敏感操作可能被濫用

認證機制解決這些問題：
- **驗證身份**：確認請求者是誰
- **控制存取**：只允許授權的使用者
- **保護資源**：防止未授權的 API 呼叫

本單元學習兩種常用認證方式：
1. **Basic Auth**：最簡單的帳密驗證（適合內部服務）
2. **Bearer Token**：業界標準的 Token 驗證（適合生產環境）

---

## ② 對應素材

| 類型 | 路徑/連結 | 說明 |
|------|----------|------|
| 範例程式碼 | `examples/U08-auth/main.go` | 整合認證的 Server |
| 範例程式碼 | `examples/U08-auth/middleware.go` | 認證 middleware |
| 測試腳本 | `examples/U08-auth/test.sh` | curl 測試指令 |
| Go SDK 範例 | [examples/auth/server](https://github.com/modelcontextprotocol/go-sdk/tree/main/examples/auth/server) | 官方 OAuth 範例（進階參考） |

---

## ③ 核心觀念

### 1. 認證 vs 授權

| 概念 | 英文 | 說明 | 比喻 |
|------|------|------|------|
| **認證** | Authentication | 確認「你是誰」 | 出示身份證 |
| **授權** | Authorization | 確認「你能做什麼」 | 檢查門禁權限 |

本單元聚焦「認證」；授權（如角色權限）屬進階主題。

### 2. 常見認證方式比較

| 方式 | 難度 | 安全性 | 適用場景 |
|------|------|--------|---------|
| **Basic Auth** | ⭐ | 低（需 HTTPS） | 內部服務、開發測試 |
| **Bearer Token** | ⭐⭐ | 中 | API 服務、生產環境 |
| **OAuth 2.0** | ⭐⭐⭐⭐ | 高 | 第三方整合、企業應用 |

### 3. Basic Auth 原理

HTTP Basic Authentication 將帳密編碼後放入 Header：

```
Authorization: Basic base64(username:password)
```

範例：
```bash
# username=admin, password=secret
# base64("admin:secret") = "YWRtaW46c2VjcmV0"

curl -H "Authorization: Basic YWRtaW46c2VjcmV0" ...

# 或使用 curl 內建支援
curl -u admin:secret ...
```

⚠️ **安全提醒**：Base64 不是加密，**必須搭配 HTTPS 使用**。

### 4. Bearer Token 原理

Bearer Token 是一種「持有者令牌」，持有 Token 即可存取：

```
Authorization: Bearer <token>
```

範例：
```bash
curl -H "Authorization: Bearer abc123xyz" ...
```

Token 的特性：
- **不透明**：Client 不需知道 Token 內容
- **可驗證**：Server 檢查 Token 是否有效
- **可過期**：設定有效期限增加安全性

### 5. Go SDK 內建支援

MCP Go SDK 提供 `RequireBearerToken` middleware：

```go
// TokenVerifier 介面
type TokenVerifier func(ctx context.Context, token string) (*TokenInfo, error)

// 使用方式
handler := mcp.RequireBearerToken(mcpHandler, myVerifier, "read")
```

### 6. Middleware Pattern

Middleware 是「包裝」handler 的設計模式，用於統一處理橫切關注點：

```
請求 → [認證 Middleware] → [業務 Handler] → 回應
         ↓ 驗證失敗
       401 Unauthorized
```

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !isAuthenticated(r) {
            http.Error(w, "Unauthorized", 401)
            return
        }
        next.ServeHTTP(w, r)  // 驗證通過，繼續處理
    })
}
```

---

## ④ 動手做

### [必做] 4.1 建立專案結構

```bash
cd ~/mcp-test
mkdir -p examples/U08-auth
```

### [必做] 4.2 實作 Basic Auth Middleware

建立 `examples/U08-auth/middleware.go`：

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ==================== Basic Auth Middleware ====================

// UserInfo 儲存使用者驗證結果的資訊
type UserInfo struct {
	UserID   string
	Username string
	Roles    []string
}

// CredentialVerifier 定義帳密驗證函數的簽章
type CredentialVerifier func(ctx context.Context, username, password string) (*UserInfo, error)

// BasicAuthMiddleware 驗證 HTTP Basic Authentication
// - verifier: 帳密驗證函數（可查 DB、LDAP、外部服務等）
// - next: 驗證成功後要執行的 handler
func BasicAuthMiddleware(verifier CredentialVerifier, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 從 Authorization header 解析 Basic Auth 帳密
		username, password, ok := r.BasicAuth()

		// 檢查是否有提供帳密
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="MCP Server"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 呼叫驗證函數驗證帳密（可查 DB、LDAP、外部服務等）
		userInfo, err := verifier(r.Context(), username, password)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="MCP Server"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 將 UserInfo 存入 context，供後續 handler 使用
		ctx := context.WithValue(r.Context(), userInfoKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ==================== Bearer Token Middleware ====================

// TokenInfo 儲存 Token 驗證結果的資訊
type TokenInfo struct {
	UserID    string
	Scopes    []string
	ExpiresAt int64
}

// TokenVerifier 定義 Token 驗證函數的簽章
type TokenVerifier func(ctx context.Context, token string) (*TokenInfo, error)

// BearerTokenMiddleware 驗證 Bearer Token
// - verifier: Token 驗證函數
// - requiredScope: 必要的權限範圍（空字串表示不檢查 scope）
// - next: 驗證成功後要執行的 handler
func BearerTokenMiddleware(verifier TokenVerifier, requiredScope string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 取得 Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 解析 "Bearer <token>" 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format (expected: Bearer <token>)", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// 呼叫驗證函數驗證 Token
		tokenInfo, err := verifier(r.Context(), token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// 檢查 Scope（如有指定）
		if requiredScope != "" && !hasScope(tokenInfo.Scopes, requiredScope) {
			http.Error(w, fmt.Sprintf("Insufficient scope: requires '%s'", requiredScope), http.StatusForbidden)
			return
		}

		// 將 TokenInfo 存入 context，供後續 handler 使用
		ctx := context.WithValue(r.Context(), tokenInfoKey, tokenInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ==================== Context Helpers ====================

// contextKey 是用於 context 存取的自訂型別，避免 key 衝突
type contextKey string

// context keys
const (
	tokenInfoKey contextKey = "tokenInfo"
	userInfoKey  contextKey = "userInfo"
)

// GetTokenInfo 從 context 取得 TokenInfo（Bearer Token 認證用）
func GetTokenInfo(ctx context.Context) *TokenInfo {
	// ctx.Value(tokenInfoKey) ... (略)
}

// GetUserInfo 從 context 取得 UserInfo（Basic Auth 認證用）
func GetUserInfo(ctx context.Context) *UserInfo {
	// ctx.Value(userInfoKey) ... (略)
}

// ==================== Helpers ====================

// hasScope 檢查 scopes 陣列中是否包含指定的 scope
func hasScope(scopes []string, required string) bool {
	// for _, s := range scopes ... (略)
}
```

### [必做] 4.3 實作主程式

建立 `examples/U08-auth/main.go`：

```go
package main

import (
	"context"
	"crypto/subtle"
	"errors"
	// "fmt" / "log" / "net/http" / "os" / "time" ... (略)

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ==================== 設定 ====================

// CredentialInfo 儲存帳密資訊（模擬資料庫中的使用者資料）
type CredentialInfo struct {
	PasswordHash string
	UserID       string
	Roles        []string
}

// AuthConfig 定義認證相關設定
type AuthConfig struct {
	Mode        string
	Credentials map[string]*CredentialInfo
	Tokens      map[string]*TokenInfo
}

// loadConfig 從環境變數載入認證設定
func loadConfig() *AuthConfig {
	// config := &AuthConfig{Mode: getEnv("AUTH_MODE", "none")} ... (略)
}

// getEnv 取得環境變數，若未設定則回傳預設值
func getEnv(key, fallback string) string {
	// if v := os.Getenv(key); v != "" ... (略)
}

// ==================== 帳密驗證器 ====================

// createCredentialVerifier 根據設定建立帳密驗證函數
//
// 本範例使用記憶體中的帳密對照表，實際環境可改為：
//   - 查詢資料庫（SELECT * FROM users WHERE username = ?）
//   - 呼叫 LDAP/Active Directory
//   - 呼叫外部認證服務
func createCredentialVerifier(config *AuthConfig) CredentialVerifier {
	return func(ctx context.Context, username, password string) (*UserInfo, error) {
		// cred, exists := config.Credentials[username] ... (略)

		// 使用 constant-time comparison 防止 timing attack
		// 實際環境應使用 bcrypt.CompareHashAndPassword
		if subtle.ConstantTimeCompare([]byte(password), []byte(cred.PasswordHash)) != 1 {
			return nil, errors.New("invalid password")
		}

		// return &UserInfo{UserID: cred.UserID} ... (略)
	}
}

// ==================== Token 驗證器 ====================

// createTokenVerifier 根據設定建立 Token 驗證函數
//
// 本範例使用記憶體中的 Token 對照表，實際環境可改為：
//   - 查詢資料庫
//   - 呼叫 OAuth Token Introspection 端點
//   - 驗證 JWT 簽章
func createTokenVerifier(config *AuthConfig) TokenVerifier {
	return func(ctx context.Context, token string) (*TokenInfo, error) {
		// info, exists := config.Tokens[token] ... (略)

		// 檢查 Token 是否已過期
		if time.Now().Unix() > info.ExpiresAt {
			return nil, errors.New("token expired")
		}

		// return info, nil ... (略)
	}
}

// ==================== Tools ====================

// AddInput 定義 add Tool 的輸入參數結構
type AddInput struct {
	// A float64 / B float64 ... (略)
}

// addHandler 處理 add Tool 的請求，計算兩數之和
func addHandler(ctx context.Context, req *mcp.CallToolRequest, input AddInput) (
	*mcp.CallToolResult, any, error,
) {
	// 提示：認證結果可從 context 取回（GetUserInfo / GetTokenInfo）
	// if userInfo := GetUserInfo(ctx); userInfo != nil ... (略)
}

// ==================== Main ====================

func main() {
	// 設定 log 輸出到 stderr（避免干擾 JSON-RPC 通訊）
	log.SetOutput(os.Stderr)

	// 載入設定
	config := loadConfig()
	port := getEnv("PORT", "8080")

	// 建立 MCP Server（前面單元已教）
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "auth-demo",
		Version: "1.0.0",
	}, nil)

	// 註冊 Tool（見結構示範檔 §3）
	// mcp.AddTool(server, &mcp.Tool{Name: "add"}, addHandler) ... (略)

	// 建立 HTTP Handler（使用 SDK 提供的 StreamableHTTPHandler）
	mcpHandler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	// 根據設定套用認證 Middleware
	var handler http.Handler = mcpHandler

	switch config.Mode {
	case "basic":
		// Basic Auth 模式：使用帳密驗證
		log.Println("Auth mode: Basic")
		// log.Println("Available users:") ... (略)
		verifier := createCredentialVerifier(config)
		handler = BasicAuthMiddleware(verifier, handler)

	case "bearer":
		// Bearer Token 模式：使用 Token 驗證
		log.Println("Auth mode: Bearer Token")
		// log.Println("Available tokens:") ... (略)
		verifier := createTokenVerifier(config)
		handler = BearerTokenMiddleware(verifier, "", handler)

	default:
		// 無認證模式
		log.Println("Auth mode: None (WARNING: no authentication)")
	}

	// 設定路由：/mcp 端點需要認證
	http.Handle("/mcp", handler)

	// Health check 端點：不需認證（供負載均衡器使用）
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 啟動 HTTP Server
	// log.Println("=================") ... (略)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
```

### [必做] 4.4 建立測試腳本

建立 `examples/U08-auth/test.sh`：

測試腳本涵蓋三種認證模式的完整測試案例。

**三模式測試矩陣**

| 測試案例 | AUTH_MODE=none | AUTH_MODE=basic | AUTH_MODE=bearer |
|---------|:-:|:-:|:-:|
| Health Check（無認證） | 200 | 200 | 200 |
| 無認證請求 | 200 | 401 | 401 |
| Basic Auth 錯誤密碼 | 200 | 401 | 401 |
| Basic Auth 正確帳密 | 200 | 200 | 401 |
| Bearer Token 無效 | 200 | 401 | 401 |
| Bearer Token 有效 | 200 | 401 | 200 |
| Tool 呼叫（帶認證） | 200 | 200 | 200 |

**代表性 curl 範例（Bearer Token 呼叫 Tool）**

```bash
curl -s -X POST "http://localhost:8080/mcp" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token-123" \
  -d '{
    "jsonrpc":"2.0",
    "id":2,
    "method":"tools/call",
    "params":{"name":"add","arguments":{"a":10,"b":20}}
  }' | jq .
```

完整腳本請參考 `examples/U08-auth/test.sh`。

### [必做] 4.5 執行並測試

**測試 1：無認證模式**
```bash
# 終端 1：啟動 Server（無認證）
AUTH_MODE=none go run ./examples/U08-auth/

# 終端 2：執行測試
chmod +x examples/U08-auth/test.sh
./examples/U08-auth/test.sh

# 預期：所有請求都成功（200）
```

**測試 2：Basic Auth 模式**
```bash
# 終端 1：啟動 Server（Basic Auth）
AUTH_MODE=basic AUTH_USER=admin AUTH_PASS=secret go run ./examples/U08-auth/

# 終端 2：執行測試
./examples/U08-auth/test.sh

# 預期結果：
# - 無認證請求 → 401 Unauthorized
# - 錯誤密碼   → 401 Unauthorized
# - 正確帳密   → 200 OK
```

**測試 3：Bearer Token 模式**
```bash
# 終端 1：啟動 Server（Bearer Token）
AUTH_MODE=bearer go run ./examples/U08-auth/

# 終端 2：執行測試
./examples/U08-auth/test.sh

# 預期結果：
# - 無效 Token → 401 Unauthorized
# - 有效 Token → 200 OK
```

### [必做] 4.6 使用 MCP Inspector 測試

MCP Inspector 目前不直接支援認證 Header，可透過以下方式測試：

**方式 1：暫時關閉認證**
```bash
AUTH_MODE=none go run ./examples/U08-auth/

# 另一終端
npx @modelcontextprotocol/inspector
# 在 Inspector 中連線 http://localhost:8080/mcp
```

**方式 2：使用 curl 測試完整流程**
```bash
# 使用認證模式啟動
AUTH_MODE=bearer go run ./examples/U08-auth/

# 手動測試各個 MCP 方法
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token-123" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | jq .
```

### [延伸] 4.7 使用 SDK 內建的 RequireBearerToken

MCP Go SDK 提供內建的 Bearer Token middleware，適合需要 OAuth Token Introspection 的場景：

```go
import "github.com/modelcontextprotocol/go-sdk/mcp"

// 實作 SDK 定義的 TokenVerifier 介面
func myVerifier(ctx context.Context, token string) (*mcp.TokenInfo, error) {
    // 實作 Token 驗證邏輯（可呼叫 OAuth Introspection 端點）
    if token == "valid-token" {
        return &mcp.TokenInfo{
            Active:  true,
            Subject: "user-1",
            Scopes:  []string{"read", "write"},
        }, nil
    }
    return nil, errors.New("invalid token")
}

func main() {
    // server := mcp.NewServer(...) ... (略)

    // 使用 SDK 內建的 RequireBearerToken
    // 第三個參數是要求的 scope（空字串表示不檢查）
    handler := mcp.RequireBearerToken(mcpHandler, myVerifier, "read")

    http.Handle("/mcp", handler)
    // http.ListenAndServe(":"+port, nil) ... (略)
}
```

### [延伸] 4.8 Token 過期處理

```go
// TokenExpiredError 表示 Token 已過期
type TokenExpiredError struct {
    ExpiredAt time.Time
}

func (e *TokenExpiredError) Error() string {
    return fmt.Sprintf("token expired at %s", e.ExpiredAt.Format(time.RFC3339))
}

// 在 TokenVerifier 中檢查過期
func createTokenVerifier(tokens map[string]*TokenInfo) func(string) (*UserInfo, error) {
	return func(token string) (*UserInfo, error) {
		// （同 4.3 主程式）info, exists := tokens[token] ... (略)

		// 新增：自定義過期錯誤
		if time.Now().Unix() > info.ExpiresAt {
			return nil, &TokenExpiredError{Token: token, ExpiredAt: info.ExpiresAt}
		}
		// userInfo 組合 ... (略)
	}
}
```

### [延伸] 4.9 Scope-based 存取控制

為不同 Tool 設定不同的權限要求：

```go
// ScopedHandler 包裝 handler，加入 scope 檢查
func ScopedHandler(requiredScope string, handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenInfo := GetTokenInfo(r.Context())
        if tokenInfo == nil {
            http.Error(w, "No token info", http.StatusUnauthorized)
            return
        }

        if !hasScope(tokenInfo.Scopes, requiredScope) {
            http.Error(w,
                fmt.Sprintf("Insufficient scope: requires '%s'", requiredScope),
                http.StatusForbidden)
            return
        }

        handler.ServeHTTP(w, r)
    })
}

// 使用範例：建立不同權限的路由
http.Handle("/mcp/read", ScopedHandler("read", readHandler))
http.Handle("/mcp/write", ScopedHandler("write", writeHandler))
http.Handle("/mcp/admin", ScopedHandler("admin", adminHandler))
```

### [延伸] 4.10 CORS + Authentication

當 Client 是瀏覽器時，需要正確處理 CORS preflight：

```go
// CORSMiddleware 處理跨域請求
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // w.Header().Set("Access-Control-Allow-Origin", "*") ... (略)

        // OPTIONS preflight 不需認證，直接回傳（關鍵：避免被 Auth 擋掉）
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// 正確的 middleware 順序：CORS → Auth → Handler
handler := CORSMiddleware(
    BearerTokenMiddleware(verifier, "", mcpHandler),
)
```

---

## ⑤ 踩坑提示

### 踩坑 1：Basic Auth 帳密明文傳輸

**現象**
使用 Wireshark 等工具可看到 Base64 編碼的帳密。

**原因**
HTTP Basic Auth 的 Base64 只是編碼，不是加密，可輕易還原。

**解法**
```bash
# 生產環境必須使用 HTTPS

# 開發時可使用 mkcert 建立本地憑證
brew install mkcert      # macOS
mkcert -install
mkcert localhost 127.0.0.1

# 或使用 ngrok 建立臨時 HTTPS 通道
ngrok http 8080
```

---

### 踩坑 2：Token 驗證失敗但沒有詳細錯誤

**現象**
Client 只看到 "Unauthorized"，不知道哪裡出錯。

**原因**
安全考量，不應洩漏太多錯誤細節給 Client。

**解法**
```go
// Server 端記錄詳細 log，但 Client 只看到簡單訊息
func BearerTokenMiddleware(...) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 解析 Authorization header ... (略)
        tokenInfo, err := verifier(r.Context(), token)
        if err != nil {
            // 詳細 log 記錄在 Server（顯示 Token 前幾字元避免洩漏完整 Token）
            log.Printf("Token verification failed: %v (token: %s...)",
                err, token[:min(8, len(token))])

            // Client 只看到簡單訊息
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        // 驗證通過，繼續處理 ... (略)
    })
}
```

---

### 踩坑 3：Timing Attack

**現象**
攻擊者透過回應時間差異猜測密碼。

**原因**
使用 `==` 比較字串時，在第一個不同字元就會返回 false，導致比較時間不一致。

**解法**
```go
import "crypto/subtle"

// 錯誤：有 timing 漏洞
if password == expected {
    // 驗證通過 ... (略)
}

// 正確：constant-time comparison
if subtle.ConstantTimeCompare([]byte(password), []byte(expected)) == 1 {
    // 驗證通過 ... (略)
}
```

---

### 踩坑 4：CORS Preflight 被認證擋掉

**現象**
瀏覽器 Console 顯示 CORS 錯誤，實際是 401。

**原因**
瀏覽器發送 OPTIONS preflight 請求時不會帶 Authorization header。

**解法**

解法見延伸 4.10 的 CORSMiddleware。關鍵是 CORS middleware 串在 Auth 之前、OPTIONS 直接放行：

```go
mux.Handle("/mcp", CORSMiddleware(authMiddleware(mcpHandler)))
```

---

### 踩坑 5：忘記處理 Health Check

**現象**
負載均衡器的 health check 因認證失敗而誤判服務不健康。

**原因**
Health check 端點也被認證 middleware 保護。

**解法**

解法：將 health check 端點獨立於認證 middleware 之外（如 4.3 主程式的做法）。

---

## ⑥ 驗收標準（DoD）

| 項目 | 驗證方式 | 預期結果 |
|------|---------|---------|
| Basic Auth（錯誤帳密） | `curl -u wrong:wrong` | 401 Unauthorized |
| Basic Auth（正確帳密） | `curl -u admin:secret` | 200 OK |
| Bearer Token（無效） | `Authorization: Bearer invalid` | 401 Unauthorized |
| Bearer Token（有效） | `Authorization: Bearer test-token-123` | 200 OK |
| Health Check | `curl /health` | 不受認證影響，回傳 OK |
| 環境變數切換 | `AUTH_MODE=basic/bearer/none` | 正確切換認證模式 |
| Tool 呼叫 | 帶認證呼叫 add Tool | 回傳正確計算結果 |

**自我檢核清單**
- [ ] Basic Auth middleware 正確驗證帳密
- [ ] Bearer Token middleware 正確驗證 Token
- [ ] 使用 `subtle.ConstantTimeCompare` 防止 timing attack
- [ ] 驗證失敗時回傳正確的 HTTP 狀態碼（401/403）
- [ ] Server log 記錄認證事件（成功/失敗）
- [ ] Health Check 端點不受認證影響
- [ ] 環境變數能正確切換認證模式

---

## 下一步

認證機制就緒後，前往 [U09｜Weather MVP](U09-Weather-MVP.md) 完成最終的天氣查詢服務。

你可以在 Weather MVP 的延伸練習中整合認證功能，打造一個安全的生產級 MCP Server。

---

## 延伸學習：OAuth 2.0

OAuth 2.0 是更完整的授權框架，適合以下場景：
- 第三方應用整合（如「使用 Google 登入」）
- 委派授權（使用者授權應用存取其資源）
- 細粒度的權限控制與 Token 生命週期管理

MCP Go SDK 提供實驗性 OAuth 支援，包含：
- `RequireBearerToken` middleware（Token Introspection）
- `ProtectedResourceMetadataHandler`（RFC 9728）
- Client 端 OAuth 流程（需 build tag）

**學習資源**：
- [MCP Go SDK Auth 範例](https://github.com/modelcontextprotocol/go-sdk/tree/main/examples/auth)
- [RFC 6749: OAuth 2.0](https://tools.ietf.org/html/rfc6749)
- [RFC 7662: Token Introspection](https://tools.ietf.org/html/rfc7662)

由於 OAuth 需要外部授權伺服器（如 Keycloak、Auth0），複雜度較高，列入進階學習方向。
