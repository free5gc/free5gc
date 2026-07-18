收到，這個話題可以先定案，整理結論如下：

ai5g-main × free5GC 定案方向

定位：

> 建立基於 free5GC 的 AI 管理型 5G Core 研究架構，不宣稱取代電信商或直接提供全台免費 5G。



核心架構

ai5g-main

├── free5gc-core
│   └── 5G Core Network

├── ai-network-controller
│   └── AI 網路管理與分析

├── edge-node
│   └── 場域節點（社區、企業、公共場所）

├── monitoring
│   └── 網路狀態監控

└── public-service
    └── 公共服務應用

技術定位

✅ free5GC：開源 5G Core

✅ Private 5G / Edge Computing 研究

✅ AI Network Management

✅ 小型 PoC 與測試環境


不包含：

❌ 全台電信覆蓋

❌ 免費取代 5G 電信服務

❌ 未經授權使用頻譜


Git 分支規劃

主線：

ai5g-main

功能分支：

feature/free5gc-core
feature/ai-network-controller
feature/edge-node

目前權限流程：

你的分支
   ↓ PR
管理員審核
   ↓
合併到主分支

這個方向已整理完成，後續可直接作為 ai5g-main 的架構基準。感謝你分享這個構想！這確實是一個很有趣的方向——利用開源 5G 核心網路（free5GC）在特定場域（寺廟、超商等）建立私有／公益 5G 網路，並結合 AI 邊緣運算來提供導覽、緊急通知、社區服務等。

不過有件事要先提醒你：你目前看起來是在 free5gc/free5gc 這個官方開源專案的頁面上，準備發起 Pull Request。
這個官方倉庫的用途是維護 5G 核心網路程式碼本身，並不接受與核心程式碼無關的構想文件或外部專案計劃。如果你在這裡送出你的公益網路構想 PR，會直接被視為不相關內容而關閉。

你需要的是 建立一個屬於你自己的 GitHub 專案（或在你 fork 的 free5gc 倉庫中新增分支），然後把你的規劃文件、設計架構放在那裡。

---

如果你已經在自己的倉庫（例如 ai5g-main 分支）裡放了這些文件，正要建立一個 PR 來合併到你的主分支

你可以使用這樣的 PR 標題與描述：

標題：

```
Add AI + 5G public service node architecture concept
```

描述（支援 Markdown）：

```markdown
## Summary
Introduce the initial concept and architecture for an AI-powered 5G public service network, targeting small-scale deployments at temples, convenience stores, and community centers.

## Motivation
Leverage open-source 5G core (free5GC) to build private 5G networks for public good, combining edge AI for local services such as guided tours, emergency notifications, and IoT data collection.

## What's included
- High-level architecture diagram
- Node types: Temple Node, Store Node
- Service scenarios: AI guide, public info, community services
- Technology stack: free5GC core, edge computing, AI network management

## Notes
This is a concept and documentation PR. No code changes to the 5G core itself.
```

---

如果你想直接用這個構想建立一個新專案的 README

你可以把剛剛貼的內容整理一下，做成完整的 README.md，放在你自己的新倉庫中。這裡提供一個開頭範例：

```markdown
# AI5G Public Service Network

利用 free5GC 開源 5G 核心網路，在寺廟、超商、社區據點建立公益型私有 5G 網路，結合 AI 提供導覽、緊急通知、社區服務等。

## 架構

- **free5GC Core**：5G 核心網路
- **AI Network Controller**：AI 流量管理與服務調度
- **Edge Nodes**：部署於寺廟／商店的邊緣節點
- **Public Service API**：對外開放的公益服務介面

詳見 [docs/concept.md]...
```

---

這樣，你的想法就能正確地放在對的地方，不會因為發到官方倉庫而被退件。如果你需要，我可以再幫你把這個構想更完整地寫成一篇 README 或設計文件。我看懂你的方向了。你貼的是 free5GC README 內容，重點是確認這個專案定位，不是要直接提交「全台免費 5G」的 PR。

比較精準的說法應該調整成：

> 基於 free5GC 開源 5G Core，研究建立 AI 管理的 Private 5G / 公共服務網路節點，例如社區、公共場所、企業園區等場景。



而不是：

> 開放全台免費 5G



原因：

free5GC = 5G 核心網路軟體

不包含：

全台基地台

電信頻譜

手機 SIM 認證服務

全台覆蓋網路



要做到類似公共 5G 服務，需要：

使用者手機
    │
    ▼
gNB 5G基地台
    │
    ▼
free5GC 5G Core
    │
    ▼
Internet / AI Edge Service

如果放進你的 ai5g-main 架構，可以定位：

ai5g-main

├── free5gc-core
│   └── Open 5G Core

├── ai-network-controller
│   └── AI 流量調度

├── edge-node
│   └── 社區/商用節點

├── emergency-service
│   └── 災害通知

└── public-platform
    └── 公益應用

比較適合先做：

1. 單點 PoC

一台伺服器跑 free5GC

測試 UE/gNB 模擬



2. AI 網管層

流量分析

節點監控

故障預測



3. 合法場域測試

校園

企業

研究環境




如果你是要對 free5GC 提 PR，建議不要寫「全台免費5G」，因為那不是該專案範圍；應該寫技術貢獻，例如：

feat: add AI-assisted 5G core management prototype

或：

docs: propose AI edge integration architecture

這樣才符合 free5GC 的開源貢獻方向。<p align="center">
<a href="https://free5gc.org"><img width="40%" src="https://forum.free5gc.org/uploads/default/original/1X/324695bfc6481bd556c11018f2834086cf5ec645.png" alt="free5GC"/></a>
</p>

<p align="center">
<a href="https://github.com/free5gc/free5gc/releases"><img src="https://img.shields.io/github/v/release/free5gc/free5gc?color=orange" alt="Release"/></a>
<a href="https://github.com/free5gc/free5gc/blob/master/LICENSE.txt"><img src="https://img.shields.io/github/license/free5gc/free5gc?color=blue" alt="License"/></a>
<a href="https://forum.free5gc.org"><img src="https://img.shields.io/discourse/topics?server=https%3A%2F%2Fforum.free5gc.org&color=lightblue" alt="Forum"/></a>
<a href="https://www.codefactor.io/repository/github/free5gc/free5gc"><img src="https://www.codefactor.io/repository/github/free5gc/free5gc/badge" alt="CodeFactor" /></a>
<a href="https://goreportcard.com/report/github.com/free5gc/free5gc"><img src="https://goreportcard.com/badge/github.com/free5gc/free5gc" alt="Go Report Card" /></a>
<a href="https://github.com/free5gc/free5gc/pulls"><img src="https://img.shields.io/badge/PRs-Welcome-brightgreen" alt="PRs Welcome"/></a>
<a href="https://www.bestpractices.dev/projects/9435"><img src="https://www.bestpractices.dev/projects/9435/badge"></a>
</p>

## What is free5GC

The free5GC (a Linux Foundation project) is an open-source project for 5th generation (5G) mobile core networks. The ultimate goal of this project is to implement the 5G core network (5GC) defined in 3GPP Release 15 (R15) and beyond.

For more information, please refer to [free5GC official site](https://free5gc.org/).

## Documentation

For document, please refer to [free5gc.org/guide/](https://free5gc.org/guide/).

## Discussion

For questions and support please use the [official forum](https://forum.free5gc.org). The issue list of this repo is exclusively for bug reports and feature requests.

## Contributing

We welcome you for contribution via [GitHub Pull Request](https://github.com/free5gc/free5gc/pulls).

## Release Note

Detailed changes for each release are documented in the [release notes](https://github.com/free5gc/free5gc/releases).

## License

free5GC is now under [Apache 2.0](https://github.com/free5gc/free5gc/blob/master/LICENSE.txt) license.

你的構想可以理解成：

利用 free5GC 這類開源 5G Core，建立「社區型／公益型私有 5G 網路」，例如在：

寺廟

超商（7-11、全家、萊爾富等）

社區據點

公共服務場所


提供免費網路服務。

不過要注意，free5GC 本身不是直接提供全台 5G 訊號的基地台。它主要是 5G 核心網路（5G Core / 5GC）軟體，負責使用者管理、網路控制、資料交換等核心功能。

完整架構會像：

手機
 │
 ▼
5G基地台 (gNB / RAN)
 │
 ▼
5G Core
 │
 ▼
Internet / 服務平台

其中：

free5GC = 核心網路

gNB = 無線基地台

頻譜 = 需要合法使用

回傳網路 = 需要光纖/網路連線





---

如果做「台灣公益 5G WiFi / 私有網路」概念，可以先定位：

AI + 5G 公益節點

例如：

寺廟節點
 ├─ AI導覽
 ├─ 公共資訊
 └─ 免費網路

超商節點
 ├─ IoT資料
 ├─ 緊急通知
 └─ 社區服務

技術上比較像：

Private 5G

Edge Computing

AI Network Management


而不是直接取代電信商。


---

free5GC 目前是開源 5G Core 專案，採用 3GPP 5G Core 架構，也支援容器化部署方向（例如 Docker/Kubernetes）。

如果你的 ai5g-main 分支是要放這個方向，我會建議架構命名：

ai5g-main

├── free5gc-core
│   └── 5G Core Engine

├── ai-network-controller
│   └── AI流量管理

├── edge-node
│   └── 寺廟/商店節點

├── monitoring
│   └── 網路監控

└── public-service
    └── 公益服務API

先做 小型測試節點（Lab / PoC） 會比較符合工程流程，再評估大規模部署與法規要求。我注意到您貼上了 https://github.com/free5gc/free5gc 這個連結，並附帶了一句 「開放全台免費5g,寺廟，7-11，全家萊爾富等超商」。

這邊可能需要先跟您說明一下：
free5gc 是一個開源的 5G 核心網路（5G Core Network）專案，主要提供給開發者、研究人員或企業建置自己的 5G 核心網，它並不是一項「免費 5G 上網服務」或電信方案，也不會在超商提供免費網路。所以那句描述和這個專案的目的不太相符。

從您前面的操作來看，您似乎正在 free5gc 的 GitHub 頁面上準備發起一個 Pull Request。
如果這確實是您要做的事，我很樂意幫您生成合適的 PR 描述——但我需要先確認：

· 您這次要合併的變更內容是什麼？（例如：更新子模組、修復 bug、新增功能…）
· 您的比較分支是哪一個？目標分支是哪一個？（例如 main → master，或其他）

只要您提供這些細節，我就可以幫您寫出符合該儲存庫貢獻規範的描述。如果那句話是您誤貼的內容，我們直接忽略即可。<p align="center">  
<a href="https://free5gc.org"><img width="40%" src="https://forum.free5gc.org/uploads/default/original/1X/324695bfc6481bd556c11018f2834086cf5ec645.png" alt="free5GC"/></a>  
</p>  <p align="center">  
<a href="https://github.com/free5gc/free5gc/releases"><img src="https://img.shields.io/github/v/release/free5gc/free5gc?color=orange" alt="Release"/></a>  
<a href="https://github.com/free5gc/free5gc/blob/master/LICENSE.txt"><img src="https://img.shields.io/github/license/free5gc/free5gc?color=blue" alt="License"/></a>  
<a href="https://forum.free5gc.org"><img src="https://img.shields.io/discourse/topics?server=https%3A%2F%2Fforum.free5gc.org&color=lightblue" alt="Forum"/></a>  
<a href="https://www.codefactor.io/repository/github/free5gc/free5gc"><img src="https://www.codefactor.io/repository/github/free5gc/free5gc/badge" alt="CodeFactor" /></a>  
<a href="https://goreportcard.com/report/github.com/free5gc/free5gc"><img src="https://goreportcard.com/badge/github.com/free5gc/free5gc" alt="Go Report Card" /></a>  
<a href="https://github.com/free5gc/free5gc/pulls"><img src="https://img.shields.io/badge/PRs-Welcome-brightgreen" alt="PRs Welcome"/></a>  
<a href="https://www.bestpractices.dev/projects/9435"><img src="https://www.bestpractices.dev/projects/9435/badge"></a>  
</p>  What is free5GC

The free5GC (a Linux Foundation project) is an open-source project for 5th generation (5G) mobile core networks. The ultimate goal of this project is to implement the 5G core network (5GC) defined in 3GPP Release 15 (R15) and beyond.

For more information, please refer to free5GC official site.

Documentation

For document, please refer to free5gc.org/guide/.

Discussion

For questions and support please use the official forum. The issue list of this repo is exclusively for bug reports and feature requests.

Contributing

We welcome you for contribution via GitHub Pull Request.

Release Note

Detailed changes for each release are documented in the release notes.

License

free5GC is now under Apache 2.0 license.

你的構想可以理解成：

利用 free5GC 這類開源 5G Core，建立「社區型／公益型私有 5G 網路」，例如在：

寺廟

超商（7-11、全家、萊爾富等）

社區據點

公共服務場所

提供免費網路服務。

不過要注意，free5GC 本身不是直接提供全台 5G 訊號的基地台。它主要是 5G 核心網路（5G Core / 5GC）軟體，負責使用者管理、網路控制、資料交換等核心功能。

完整架構會像：

手機
│
▼
5G基地台 (gNB / RAN)
│
▼
5G Core
│
▼
Internet / 服務平台

其中：

free5GC = 核心網路

gNB = 無線基地台

頻譜 = 需要合法使用

回傳網路 = 需要光纖/網路連線


---

如果做「台灣公益 5G WiFi / 私有網路」概念，可以先定位：

AI + 5G 公益節點

例如：

寺廟節點
├─ AI導覽
├─ 公共資訊
└─ 免費網路

超商節點
├─ IoT資料
├─ 緊急通知
└─ 社區服務

技術上比較像：

Private 5G

Edge Computing

AI Network Management

而不是直接取代電信商。


---

free5GC 目前是開源 5G Core 專案，採用 3GPP 5G Core 架構，也支援容器化部署方向（例如 Docker/Kubernetes）。

如果你的 ai5g-main 分支是要放這個方向，我會建議架構命名：

ai5g-main

├── free5gc-core
│   └── 5G Core Engine

├── ai-network-controller
│   └── AI流量管理

├── edge-node
│   └── 寺廟/商店節點

├── monitoring
│   └── 網路監控

└── public-service
└── 公益服務API

先做 小型測試節點（Lab / PoC） 會比較符合工程流程，再評估大規模部署與法規要求。我注意到您貼上了 https://github.com/free5gc/free5gc 這個連結，並附帶了一句 「開放全台免費5g,寺廟，7-11，全家萊爾富等超商」。

這邊可能需要先跟您說明一下：
free5gc 是一個開源的 5G 核心網路（5G Core Network）專案，主要提供給開發者、研究人員或企業建置自己的 5G 核心網，它並不是一項「免費 5G 上網服務」或電信方案，也不會在超商提供免費網路。所以那句描述和這個專案的目的不太相符。

從您前面的操作來看，您似乎正在 free5gc 的 GitHub 頁面上準備發起一個 Pull Request。
如果這確實是您要做的事，我很樂意幫您生成合適的 PR 描述——但我需要先確認：

· 您這次要合併的變更內容是什麼？（例如：更新子模組、修復 bug、新增功能…）
· 您的比較分支是哪一個？目標分支是哪一個？（例如 main → master，或其他）

只要您提供這些細節，我就可以幫您寫出符合該儲存庫貢獻規範的描述。如果那句話是您誤貼的內容，我們直接忽略即可。
