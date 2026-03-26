// Package main 展示如何為 MCP HTTP Server 加入認證機制
//
// 本範例支援三種認證模式，透過環境變數 AUTH_MODE 切換：
//   - none: 不啟用認證（預設）
//   - basic: HTTP Basic Authentication
//   - bearer: Bearer Token Authentication
//
// 使用方式：
//
//	# 無認證模式
//	AUTH_MODE=none go run .
//
//	# Basic Auth 模式
//	AUTH_MODE=basic AUTH_USER=admin AUTH_PASS=secret go run .
//
//	# Bearer Token 模式
//	AUTH_MODE=bearer go run .
package main

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ==================== 設定 ====================

// CredentialInfo 儲存帳密資訊（模擬資料庫中的使用者資料）
type CredentialInfo struct {
	PasswordHash string   // 密碼（實際應存 hash，這裡簡化為明文）
	UserID       string   // 使用者 ID
	Roles        []string // 角色
}

// AuthConfig 定義認證相關設定
type AuthConfig struct {
	Mode        string                     // 認證模式: "none", "basic", "bearer"
	Credentials map[string]*CredentialInfo // Basic Auth 帳密對照表（模擬 DB）
	Tokens      map[string]*TokenInfo      // Bearer Token 對照表（模擬 DB）
}

// loadConfig 從環境變數載入認證設定
//
// 支援的環境變數：
//   - AUTH_MODE: 認證模式 (none/basic/bearer)，預設 none
//   - PORT: 監聽埠號，預設 8080
//
// 注意：實際環境中，帳密和 Token 應從資料庫載入
func loadConfig() *AuthConfig {
	config := &AuthConfig{
		Mode:        getEnv("AUTH_MODE", "none"),
		Credentials: make(map[string]*CredentialInfo),
		Tokens:      make(map[string]*TokenInfo),
	}

	// 預設的測試帳密（模擬從資料庫載入）
	// 實際環境應查詢資料庫，密碼應使用 bcrypt 等方式雜湊
	config.Credentials["admin"] = &CredentialInfo{
		PasswordHash: "secret", // 實際應存 bcrypt hash
		UserID:       "user-1",
		Roles:        []string{"admin", "user"},
	}
	config.Credentials["guest"] = &CredentialInfo{
		PasswordHash: "guest123",
		UserID:       "user-2",
		Roles:        []string{"user"},
	}

	// 預設的測試 Token（模擬從資料庫載入）
	config.Tokens["test-token-123"] = &TokenInfo{
		UserID:    "user-1",
		Scopes:    []string{"read", "write"},
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	config.Tokens["readonly-token"] = &TokenInfo{
		UserID:    "user-2",
		Scopes:    []string{"read"},
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	return config
}

// getEnv 取得環境變數，若未設定則回傳預設值
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
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
		// 從對照表查找帳號（實際應查詢資料庫）
		cred, exists := config.Credentials[username]
		if !exists {
			return nil, errors.New("user not found")
		}

		// 使用 constant-time comparison 防止 timing attack
		// 實際環境應使用 bcrypt.CompareHashAndPassword
		if subtle.ConstantTimeCompare([]byte(password), []byte(cred.PasswordHash)) != 1 {
			return nil, errors.New("invalid password")
		}

		// 回傳使用者資訊
		return &UserInfo{
			UserID:   cred.UserID,
			Username: username,
			Roles:    cred.Roles,
		}, nil
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
		// 從對照表查找 Token（實際應查詢資料庫）
		info, exists := config.Tokens[token]
		if !exists {
			return nil, errors.New("token not found")
		}

		// 檢查 Token 是否已過期
		if time.Now().Unix() > info.ExpiresAt {
			return nil, errors.New("token expired")
		}

		return info, nil
	}
}

// ==================== Tools ====================

// AddInput 定義 add Tool 的輸入參數結構
type AddInput struct {
	A float64 `json:"a" jsonschema:"description=第一個數字"`
	B float64 `json:"b" jsonschema:"description=第二個數字"`
}

// addHandler 處理 add Tool 的請求，計算兩數之和
func addHandler(ctx context.Context, req *mcp.CallToolRequest, input AddInput) (
	*mcp.CallToolResult, any, error,
) {
	// 記錄呼叫者資訊（支援 Basic Auth 和 Bearer Token 兩種認證方式）
	if userInfo := GetUserInfo(ctx); userInfo != nil {
		log.Printf("Tool 'add' called by user: %s (roles: %v)", userInfo.Username, userInfo.Roles)
	} else if tokenInfo := GetTokenInfo(ctx); tokenInfo != nil {
		log.Printf("Tool 'add' called by user: %s (scopes: %v)", tokenInfo.UserID, tokenInfo.Scopes)
	}

	// 執行加法運算
	result := input.A + input.B

	// 回傳結果
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
	// 設定 log 輸出到 stderr（避免干擾 JSON-RPC 通訊）
	log.SetOutput(os.Stderr)

	// 載入設定
	config := loadConfig()
	port := getEnv("PORT", "8080")

	// 建立 MCP Server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "auth-demo",
		Version: "1.0.0",
	}, nil)

	// 註冊 Tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "計算兩數之和",
	}, addHandler)
	log.Println("Tool registered: add")

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
		log.Println("Available users:")
		log.Println("  - admin:secret (roles: admin, user)")
		log.Println("  - guest:guest123 (roles: user)")
		verifier := createCredentialVerifier(config)
		handler = BasicAuthMiddleware(verifier, handler)

	case "bearer":
		// Bearer Token 模式：使用 Token 驗證
		log.Println("Auth mode: Bearer Token")
		log.Println("Available tokens:")
		log.Println("  - test-token-123 (scopes: read, write)")
		log.Println("  - readonly-token (scopes: read)")
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
	log.Println("========================================")
	log.Printf("HTTP MCP Server listening on :%s", port)
	log.Printf("Endpoints:")
	log.Printf("  POST /mcp    - MCP JSON-RPC endpoint")
	log.Printf("  GET  /health - Health check (no auth)")
	log.Println("========================================")
	log.Println("Press Ctrl+C to stop")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
