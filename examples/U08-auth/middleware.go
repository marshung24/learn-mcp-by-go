// Package main 提供 MCP Server 的認證 middleware 實作
//
// 本檔案包含兩種常用的 HTTP 認證機制：
// - Basic Auth: 簡單的帳密驗證，適合內部服務
// - Bearer Token: Token-based 驗證，適合 API 服務
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
	UserID   string   // 使用者 ID
	Username string   // 使用者名稱
	Roles    []string // 角色（如 "admin", "user"）
}

// CredentialVerifier 定義帳密驗證函數的簽章
//
// 實作此函數以自訂帳密驗證邏輯，例如：
//   - 查詢資料庫
//   - 呼叫 LDAP/AD 服務
//   - 查詢外部認證服務
//
// 參數：
//   - ctx: 請求的 context
//   - username: 使用者輸入的帳號
//   - password: 使用者輸入的密碼
//
// 回傳：
//   - *UserInfo: 驗證成功時回傳使用者資訊
//   - error: 驗證失敗時回傳錯誤
type CredentialVerifier func(ctx context.Context, username, password string) (*UserInfo, error)

// BasicAuthMiddleware 驗證 HTTP Basic Authentication
//
// 使用方式：
//
//	verifier := func(ctx context.Context, username, password string) (*UserInfo, error) {
//	    // 驗證邏輯（查 DB、LDAP 等）
//	}
//	handler := BasicAuthMiddleware(verifier, myHandler)
//
// 參數：
//   - verifier: 帳密驗證函數
//   - next: 驗證成功後要執行的 handler
//
// 回傳：
//   - 包裝後的 http.Handler，會先驗證帳密再執行 next
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
	UserID    string   // 使用者 ID
	Scopes    []string // 權限範圍（如 "read", "write", "admin"）
	ExpiresAt int64    // 過期時間（Unix timestamp）
}

// TokenVerifier 定義 Token 驗證函數的簽章
//
// 實作此函數以自訂 Token 驗證邏輯，例如：
//   - 查詢資料庫
//   - 呼叫 OAuth Token Introspection 端點
//   - 驗證 JWT 簽章
type TokenVerifier func(ctx context.Context, token string) (*TokenInfo, error)

// BearerTokenMiddleware 驗證 Bearer Token
//
// 使用方式：
//
//	verifier := func(ctx context.Context, token string) (*TokenInfo, error) {
//	    // 驗證邏輯
//	}
//	handler := BearerTokenMiddleware(verifier, "read", myHandler)
//
// 參數：
//   - verifier: Token 驗證函數
//   - requiredScope: 必要的權限範圍（空字串表示不檢查 scope）
//   - next: 驗證成功後要執行的 handler
//
// 回傳：
//   - 包裝後的 http.Handler
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
//
// 使用方式：
//
//	func myHandler(w http.ResponseWriter, r *http.Request) {
//	    if info := GetTokenInfo(r.Context()); info != nil {
//	        log.Printf("User: %s", info.UserID)
//	    }
//	}
func GetTokenInfo(ctx context.Context) *TokenInfo {
	if info, ok := ctx.Value(tokenInfoKey).(*TokenInfo); ok {
		return info
	}
	return nil
}

// GetUserInfo 從 context 取得 UserInfo（Basic Auth 認證用）
//
// 使用方式：
//
//	func myHandler(w http.ResponseWriter, r *http.Request) {
//	    if info := GetUserInfo(r.Context()); info != nil {
//	        log.Printf("User: %s", info.Username)
//	    }
//	}
func GetUserInfo(ctx context.Context) *UserInfo {
	if info, ok := ctx.Value(userInfoKey).(*UserInfo); ok {
		return info
	}
	return nil
}

// ==================== Helpers ====================

// hasScope 檢查 scopes 陣列中是否包含指定的 scope
func hasScope(scopes []string, required string) bool {
	for _, s := range scopes {
		if s == required {
			return true
		}
	}
	return false
}

// ==================== CORS Middleware ====================

// CORSMiddleware 處理跨域請求（Cross-Origin Resource Sharing）
//
// 重要：此 middleware 必須放在 Auth middleware 之前，
// 因為瀏覽器的 preflight (OPTIONS) 請求不會帶 Authorization header。
//
// 使用方式：
//
//	handler := CORSMiddleware(AuthMiddleware(myHandler))
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 設定 CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 小時

		// OPTIONS preflight 請求直接回傳，不進入認證流程
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
