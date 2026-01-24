# generate_api_docs_mLua&nbsp;[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT) ![Go](https://img.shields.io/badge/Language-Go-00ADD8?style=flat&logo=go&logoColor=white) ![Go Version](https://img.shields.io/badge/Version-1.25.4-00ADD8?style=flat&logo=go&logoColor=white)

MaplestoryWorlds Lua(.mlua)를 위한 API 문서 자동 생성기입니다. 소스 코드 내 특수 주석을 파싱하여 깔끔하고 탐색하기 쉬운 Markdown 형식의 API 문서를 생성합니다.

## ✨ 주요 기능

- `.mlua` 파일의 특수 주석(`@Logic`, `@Component` 등)을 분석하여 문서 생성
- `Properties`, `Methods`, `Handlers` 등 코드 구조를 자동으로 인식하고 분류
- `ExecSpace`, `EventSender` 등의 속성을 기반으로 시각적인 뱃지 생성
- 타입 정보를 분석하여 관련 문서로 바로 이동할 수 있는 하이퍼링크 자동 생성
- CSS를 포함한 독립적인 Markdown 파일을 생성하여 별도 설정 없이 깔끔한 스타일 적용

## 🚀 시작하기

### 1. 소스 코드 준비

문서를 생성할 `.mlua` 파일에 아래 형식에 맞춰 주석을 작성합니다.

```lua
@Logic
---@description "플레이어의 상태를 관리하는 로직"

---@description "플레이어가 스폰될 때 호출됩니다."
---@param string playerName "플레이어 이름"
---@param number playerHealth "초기 체력"
@EventSender("Logic", "GameManager")
handler OnPlayerSpawn(string playerName, number playerHealth)

---@description "플레이어에게 데미지를 입힙니다."
---@param number damage "입힐 데미지 양"
method void TakeDamage(number damage)

```

### 2. 문서 생성 실행

프로젝트 루트에서 `main.go`를 실행하면 `RootDesk/MyDesk` 디렉토리 내의 모든 `.mlua` 파일을 탐색하여 문서를 생성하고 `document/api` 폴더에 저장합니다.

```bash
go run cmd/main.go
```

## 📝 문서 생성 예시

- **입력** (`.mlua` 파일)

    ```lua
    ---@description "게임을 관리하는 로직 입니다."
    @Logic
    script GameLogic extends Logic

        ---@description "플레이어 접속 시 호출"
        ---@param string playerName "접속한 플레이어 이름"
        @EventSender("Logic", "AuthLogic")
        handler OnPlayerConnect(string playerName)

        ---@description "서버에 메시지를 전송합니다."
        ---@param string message "전송할 메시지"
        @ExecSpace("ServerOnly")
        method void SendMessageToServer(string message)
    end
    ```

- **출력** (생성된 `*.md` 파일)

    # GameLogic

    게임을 관리하는 로직 입니다.
    
    ## Handlers

    <table style="width: 100%; border-collapse: collapse; border: 1px solid #ccc; margin-bottom: 16px;">
        <thead>
            <tr>
                <th style="background-color: #f0f0f0; padding: 10px 5px; text-align: left; vertical-align: top;">
                    <span style="font-weight: bold;">OnPlayerConnect</span>(<a href="#" style="text-decoration: none; color: #3167ad;">string</a> playerName) <img src="https://img.shields.io/badge/Logic-95e1d3" alt="Logic" style="vertical-align: middle; margin-left: 8px;">
                </th>
            </tr>
        </thead>
        <tbody>
            <tr><td style="background-color: #fff; padding: 10px 5px; text-align: left; vertical-align: top;">플레이어 접속 시 호출</td></tr>
            <tr><td style="background-color: #fafafa; border-top: 1px solid #eee; padding: 10px 5px 10px 15px; text-align: left; vertical-align: top;"><strong>Logic:</strong> AuthLogic</td></tr>
            <tr><td style="background-color: #fafafa; border-top: 1px solid #eee; padding: 10px 5px 10px 15px; text-align: left; vertical-align: top;"><code style="background-color: #e1e4e8; padding: 2px 5px; border-radius: 4px; font-family: monospace;">playerName</code><span style="color: #57606a;"> &nbsp;|&nbsp; 접속한 플레이어 이름</span></td></tr>
        </tbody>
    </table>

    ## Methods

    <table style="width: 100%; border-collapse: collapse; border: 1px solid #ccc; margin-bottom: 16px;">
        <thead>
            <tr>
                <th style="background-color: #f0f0f0; padding: 10px 5px; text-align: left; vertical-align: top;">
                    <span style="color: #3167ad;">void</span> <span style="font-weight: bold;">SendMessageToServer</span>(<a href="#" style="text-decoration: none; color: #3167ad;">string</a> message) <img src="https://img.shields.io/badge/ServerOnly-da70d6" alt="ServerOnly" style="vertical-align: middle; margin-left: 8px;">
                </th>
            </tr>
        </thead>
        <tbody>
            <tr><td style="background-color: #fff; padding: 10px 5px; text-align: left; vertical-align: top;">서버에 메시지를 전송합니다.</td></tr>
            <tr><td style="background-color: #fafafa; border-top: 1px solid #eee; padding: 10px 5px 10px 15px; text-align: left; vertical-align: top;"><code style="background-color: #e1e4e8; padding: 2px 5px; border-radius: 4px; font-family: monospace;">message</code><span style="color: #57606a;"> &nbsp;|&nbsp; 전송할 메시지</span></td></tr>
        </tbody>
    </table>

## 📂 프로젝트 구조

```
.
├─ cmd/main.go                # 프로그램 진입점
└─ pkg/
    ├─ document/               # 소스 코드 파싱 및 구조화
    │   ├─ parse.go
    │   └─ struct.go
    └─ generator/              # Markdown 문서 생성
        ├─ generate.go
        ├─ templates.go
        └─ style.css
```
