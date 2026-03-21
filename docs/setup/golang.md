# Go 語言安裝設定

> 本教學使用 Go 開發 MCP Server，需安裝 Go 1.21 或更高版本。

---

## 安裝方式

### macOS (Homebrew)

```bash
brew install go
```

### macOS/Linux (官方安裝包)

```bash
# 1. 下載安裝包
# 前往 https://go.dev/dl/ 下載對應系統的安裝包

# 2. 解壓縮到 /usr/local
sudo tar -C /usr/local -xzf go1.21.x.linux-amd64.tar.gz

# 3. 設定環境變數（加入 ~/.zshrc 或 ~/.bashrc）
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$(go env GOPATH)/bin

# 4. 重新載入設定
source ~/.zshrc
```

### Windows

1. 前往 https://go.dev/dl/ 下載 `.msi` 安裝檔
2. 執行安裝程式
3. 完成後重新開啟終端機

---

## 驗證安裝

```bash
# 確認版本（需顯示 1.21 或更高）
go version

# 確認環境變數
go env GOROOT
go env GOPATH
```

---

## Go 環境架構

```
$HOME/
├── go/                    # GOPATH（預設）
│   ├── bin/              # 編譯後的執行檔
│   ├── pkg/              # 快取的套件
│   └── src/              # 原始碼（舊式，現已少用）
└── your-project/         # 專案目錄（任意位置）
    ├── go.mod            # 模組定義
    ├── go.sum            # 依賴鎖定
    └── main.go           # 程式碼
```

---

## 環境變數說明

| 變數 | 說明 | 查看指令 |
|------|------|---------|
| `GOROOT` | Go 安裝路徑 | `go env GOROOT` |
| `GOPATH` | 工作區路徑 | `go env GOPATH` |
| `GOPROXY` | 套件下載代理 | `go env GOPROXY` |

---

## 初始化專案

```bash
# 建立專案目錄
mkdir -p ~/mcp-test && cd ~/mcp-test

# 初始化 Go Module
go mod init mcp-test

# 安裝 MCP Go SDK
go get github.com/modelcontextprotocol/go-sdk/mcp
```

---

## 測試安裝

建立 `main.go`：

```go
package main

import "fmt"

func main() {
    fmt.Println("Go is ready for MCP!")
}
```

執行：

```bash
go run main.go
# 預期輸出：Go is ready for MCP!
```

---

## 常見問題

### go: command not found

**原因**：Go 未加入 PATH 環境變數。

**解法**：
```bash
# 編輯 ~/.zshrc 或 ~/.bashrc，加入：
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$(go env GOPATH)/bin

# 重新載入
source ~/.zshrc
```

### go mod init 失敗

**現象**：
```bash
go: cannot determine module path for source directory
```

**解法**：指定模組名稱
```bash
go mod init my-module-name
# 或
go mod init github.com/username/project
```

### VS Code 無法找到 Go

**解法**：
1. 從終端機啟動 VS Code：`code .`
2. 或在 VS Code 設定中指定 Go 路徑：
   - 開啟設定 (Cmd+,)
   - 搜尋 "Go: GOROOT"
   - 設定為 `/usr/local/go`

---

## 參考資源

- [Go 官方安裝指南](https://go.dev/doc/install)
- [VS Code Go 擴充](https://marketplace.visualstudio.com/items?itemName=golang.Go)
