# MCP Server 開發新手入門教案

> 以 Golang 建立 MCP Server，從 Hello World 到天氣查詢 MVP。
>
> 版本：v1.3.0
> 最後更新：2026-03-27

---

## 假設清單

> 以下假設待使用者確認：

| 項目 | 假設值 | 備註 |
|------|--------|------|
| 學員背景 | 學過程式基礎，熟悉至少一種程式語言 | 不要求 Go 經驗 |
| Go 經驗 | 無或初學 | 教案包含 Go 基礎補充 |
| MCP 經驗 | 無 | 從零開始 |
| 開發環境 | macOS / Windows / Linux | 需支援 Go 1.21+ |
| 總學習時數 | 14-18 小時 | 自學導讀模式（含認證單元） |
| MVP 案例 | 美國國家氣象局 API (NWS) | 僅支援美國地區座標，選擇理由：免費、無需 API Key、文件完整 |

---

## §1 教案資訊

| 項目 | 說明 |
|------|------|
| 對象 | 有基礎程式概念（變數、迴圈、函式）與命令列操作經驗，但**尚未做過 MCP 開發**的初學者 |
| 目標 | 獨立完成一個具備 Tools、Resources、Prompts 三種 capabilities，支援 Stdio/HTTP 兩種傳輸模式，並具備認證機制的天氣查詢 MCP Server |
| 技術棧 | Go 1.21+、MCP Go SDK v1.2.0+、MCP Inspector（測試用） |
| 時數 | 課前準備（自學）+ 9 單元（各 45–120 分鐘） |
| 教學方式 | 觀念講解 → 範例實作 → 動手做 → 自我驗收 |
| 先備知識 | 基礎程式概念（任一語言皆可）、HTTP 概念（GET / POST）、JSON 格式 |
| 不含範圍 | MCP Client 開發、OAuth 2.0 完整流程、WebSocket Transport、雲端進階部署 |

---

## §2 學習成果（結訓時能做到）

1. **能說明** MCP 協定的核心概念與三種 capabilities（Tools、Resources、Prompts）的差異
2. **能建立** 最小可執行的 MCP Server 並與 Claude Desktop/Code 連接
3. **能使用** MCP Inspector 測試與除錯 Server
4. **能實作** Tool capability，定義輸入參數與回傳結果
5. **能實作** Resource capability，提供靜態或動態資料給 Client
6. **能實作** Prompt capability，建立可重用的提示詞模板
7. **能切換** Stdio 與 HTTP 兩種傳輸模式
8. **能實作** HTTP Server 認證機制（Basic Auth、Bearer Token）
9. **能整合** 三種 capabilities 完成一個 MVP 等級的天氣查詢服務

---

## §3 閱讀指引

### 練習標記

| 標記 | 意義 |
|------|------|
| **[必做]** | 核心練習，必須完成才能進入下一單元 |
| **[延伸]** | 加深理解的額外練習，時間允許時進行 |

### 依賴標記

- `⟵ 需完成 UXX` = 需先完成指定單元的對應練習
- `⟵ 需完成 UXX 延伸` = 特指依賴某單元的「延伸挑戰」成果

### 單元編號格式

| 編號 | 說明 |
|------|------|
| U00 | 課前準備（自學） |
| U01-U08 | 核心單元 |

### 測試工具定位

| 工具 | 定位 | 使用時機 |
|------|------|---------|
| **MCP Inspector** | 主要測試工具 | 開發中的功能測試、除錯、查看 JSON-RPC |
| **Claude Desktop/Code** | 整合驗證 | 確認 Server 能被 AI 正常使用（設定後簡單驗證即可） |

### 教案與範例的關係

| 文件 | 閱讀目的 |
|------|---------|
| 各單元 MD（`units/UXX-*.md`） | 理解原理、架構、流程與關鍵邏輯 |
| 範例程式碼（`examples/UXX-*/`） | 取得完整可執行的實作（Single Source of Truth） |
| 共用指引（`common/`） | Inspector 操作流程、Claude Desktop 設定流程 |

- 教案 MD 中的程式碼**刻意省略非關鍵部分**（以 `// 定位線索 ... (略)` 示意），聚焦架構與重點
- 教案與範例使用**一致的結構單位註解**，可依註解互相對照定位
- 建議閱讀順序：先讀教案理解「為什麼」→ 再看 `examples/` 理解「怎麼做」

### 單元內部結構

每個單元依序包含六個子區塊：
1. **為什麼先教這個？** — 教學動機
2. **對應素材** — 檔案、指令、連結
3. **核心觀念** — 3-7 項最小可行知識
4. **動手做** — 必做 + 延伸挑戰
5. **踩坑提示** — 現象 → 原因 → 解法
6. **驗收標準（DoD）** — 可觀察的完成條件

---

## §4 課前準備

詳見 [U00 課前準備](units/U00-課前準備.md)

**安裝清單摘要**：

| 軟體 | 版本需求 | 安裝指南 | 備註 |
|------|---------|---------|------|
| Go | 1.21+ | [golang.md](../setup/golang.md) | 必要 |
| Git | — | — | 必要 |
| VS Code + Go 擴充 | — | — | 建議 |
| Node.js | 18+ | — | MCP Inspector 需要 |
| MCP Inspector | — | [mcp-inspector.md](../setup/mcp-inspector.md) | 主要測試工具 |
| Claude Desktop | — | [claude-desktop.md](../setup/claude-desktop.md) | 選用，整合驗證 |
| Claude Code | — | [claude-code.md](../setup/claude-code.md) | 選用，整合驗證 |

**驗證指令**：
```bash
go version          # 確認顯示 1.21+
go env GOPATH       # 確認路徑有效
node --version      # 確認顯示 18+
```

---

## §5 課程模組

### 地圖層（單元總覽）

| 單元 | 名稱 | 時數 | 一句話目標 | 前置依賴 |
|------|------|------|-----------|----------|
| [U00](units/U00-課前準備.md) | 課前準備 | 30-45 min | 安裝開發環境並驗證可用 | — |
| [U01](units/U01-Hello-MCP-Server.md) | Hello MCP Server | 90 min | 能建立並執行第一個 MCP Server | U00 |
| [U02](units/U02-MCP-Inspector.md) | MCP Inspector | 45 min | 能使用 Inspector 測試與除錯 Server | U01 |
| [U03](units/U03-Tool-Capability.md) | Tool Capability | 120 min | 能定義並測試 Tool 功能 | U02 |
| [U04](units/U04-Resource-Capability.md) | Resource Capability | 90 min | 能提供 Resource 給 Client 讀取 | U02 |
| [U05](units/U05-Prompt-Capability.md) | Prompt Capability | 60 min | 能建立 Prompt 模板 | U02 |
| [U06](units/U06-Combined-Capabilities.md) | Combined Capabilities | 90 min | 能整合三種功能於同一 Server | U03, U04, U05 |
| [U07](units/U07-HTTP-Transport.md) | HTTP Transport | 90 min | 能將 Stdio 切換為 HTTP 模式 | U06 |
| [U08](units/U08-Server-Auth.md) | Server 認證 | 90 min | 能為 HTTP Server 加入認證機制 | U07 |
| [U09](units/U09-Weather-MVP.md) | Weather MVP | 120 min | 能完成完整的天氣查詢服務 | U08 |

### 學習路徑圖

```
U00 課前準備
    │
    ▼
U01 Hello MCP Server（含 Claude Desktop/Code 簡易設定）
    │
    ▼
U02 MCP Inspector（主要測試工具）
    │
    ├───────────────┬───────────────┐
    ▼               ▼               ▼
U03 Tools      U04 Resources    U05 Prompts
    │               │               │
    └───────────────┴───────────────┘
                    │
                    ▼
             U06 Combined
                    │
                    ▼
           U07 HTTP Transport
                    │
                    ▼
            U08 Server 認證
                    │
                    ▼
            U09 Weather MVP
```

> **設計說明**：單元順序遵循「漸進式推演」原則——先跑起最小 Server（U01），學會測試工具（U02），再分別學習三種 capabilities（U03-U05），接著整合（U06），最後加入 HTTP 模式並完成 MVP（U07-U08）。每個單元都有可見的產出，避免「學了半天什麼都看不到」的挫折。

---

## §6 里程碑檢核

| 里程碑 | 涵蓋單元 | 達成標準 | 未達標補救 |
|--------|----------|----------|-----------|
| **M1：環境就緒** | U00, U01, U02 | 能執行 Hello Server 並用 Inspector 測試 | 重新檢查安裝步驟；查閱踩坑提示 |
| **M2：Capabilities 掌握** | U03-U05 | 能獨立實作三種 capabilities 並用 Inspector 驗證 | 重做對應單元練習；參考範例程式碼 |
| **M3：整合完成** | U06, U07 | 能建立整合 Server 並支援 HTTP | 檢查 ServerCapabilities 設定 |
| **M4：認證就緒** | U08 | 能為 HTTP Server 加入 Basic Auth 或 Bearer Token | 檢查 middleware 順序與環境變數 |
| **M5：MVP 交付** | U09 | Weather MVP 所有功能正常 | 逐一測試各功能；檢查 API 整合 |

---

## §7 自學流程建議

每個單元建議的學習節奏：

| 階段 | 比重 | 說明 |
|------|------|------|
| 閱讀 | 15-20% | 先看核心觀念，理解「為什麼」 |
| 執行範例 | 20-30% | 跑通範例程式碼，觀察輸出 |
| 必做練習 | 40-50% | 動手實作，用 MCP Inspector 測試 |
| 自我驗收 | 10% | 對照 DoD 確認完成度 |
| 延伸挑戰 | 選做 | 時間允許時深化理解 |

**標準測試流程**：
```
1. 撰寫程式碼
2. 使用 MCP Inspector 測試功能
3. 查看 Logs 修正問題
4. （選做）用 Claude Desktop/Code 整合驗證
```

**每日建議**：專注 1-2 個單元，避免一次學太多導致消化不良。

---

## §8 評量機制

| 維度 | 比重 | 量化門檻 |
|------|------|---------|
| 必做練習完成度 | 50% | 100% 完成 |
| MVP 功能完整度 | 30% | 4/4 功能可運作 |
| 程式碼品質 | 10% | 無 panic、有基本錯誤處理 |
| 文件理解 | 10% | 能說明 MCP 三種 capabilities 差異 |

**及格標準**：必做練習 100% + MVP 功能 3/4 以上

---

## §9 參考文件

### 專案內文件

| 文件 | 用途 |
|------|------|
| [專案 README](../../README.md) | 專案說明與快速開始 |
| [教案綱要](新手入門教案綱要.md) | 教案編寫規格（給教案作者） |

### 安裝設定指南

| 文件 | 說明 |
|------|------|
| [Go 語言安裝](../setup/golang.md) | Go 安裝、環境變數、模組初始化 |
| [MCP Inspector](../setup/mcp-inspector.md) | Inspector 安裝與使用方式 |
| [Claude Desktop](../setup/claude-desktop.md) | Desktop 安裝與 MCP Server 設定 |
| [Claude Code](../setup/claude-code.md) | CLI 安裝與 mcp 指令使用 |

### 共用操作指引

| 文件 | 說明 |
|------|------|
| [MCP Inspector 測試指引](common/MCP-Inspector-測試指引.md) | 各單元共用的 Inspector 操作流程（Stdio / HTTP / 認證） |
| [Claude Desktop 設定指引](common/Claude-Desktop-設定指引.md) | 各單元共用的 Desktop 設定流程 + 各單元對照表 |
| [結構示範檔](common/結構示範檔.md) | Server / Handler 程式骨架速查（Stdio、HTTP、Tool、Resource、Prompt） |

### 外部參考資源

| 主題 | 資源 |
|------|------|
| MCP 官方文件 | [modelcontextprotocol.io/docs](https://modelcontextprotocol.io/docs) |
| MCP Inspector | [modelcontextprotocol.io/docs/tools/inspector](https://modelcontextprotocol.io/docs/tools/inspector) |
| Go SDK GitHub | [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) |
| Build Server 教學 | [modelcontextprotocol.io/docs/develop/build-server](https://modelcontextprotocol.io/docs/develop/build-server) |
| Go 語言入門 | [go.dev/tour](https://go.dev/tour/) |
| JSON-RPC 2.0 規範 | [jsonrpc.org/specification](https://www.jsonrpc.org/specification) |

---

## §10 延伸學習方向

- **MCP Client 開發**：學習如何建立自己的 MCP Client
- **進階 Transport**：WebSocket、自訂 Transport
- **OAuth 2.0 整合**：在 U08 基礎上整合完整 OAuth 流程（需外部授權伺服器）
- **多語言 SDK**：了解 Python、TypeScript 版本的差異
- **生產環境部署**：Docker 化、監控、日誌收集、HTTPS 設定

---

## §11 維護與版本管理

### 版本紀錄

| 版本 | 日期 | 變更說明 |
|------|------|---------|
| v1.3.0 | 2026-03-27 | 教案精簡：抽出共用指引(common/)、省略註解格式統一、跨單元樣板去重 |
| v1.2.0 | 2026-03-25 | 新增 U08 Server 認證單元；原 U08 改為 U09 |
| v1.1.0 | 2026-03-21 | 重整章節編號、MCP Inspector 升格為 U02、精簡 Claude Desktop/Code 內容 |
| v1.0.0 | 2026-03-21 | 初版發布 |

### 開課前檢查清單

- [ ] Go SDK 版本是否仍為最新穩定版？
- [ ] MCP Inspector 是否仍可正常運作？
- [ ] 範例程式碼是否能正常編譯執行？
- [ ] Claude Desktop/Code 設定方式是否有變更？
- [ ] 外部 API（NWS）是否仍可用？
- [ ] 所有連結是否有效？
