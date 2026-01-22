package main

import (
	"fmt"
	"generate_api_docs_mLua/pkg/document"
)

func main() {
	testContent := `
---@description "서버 상태: 초기화 중"
property string State_Initializing = "Initializing"

---@description "서버의 현재 상태를 클라이언트에 전송합니다."
---@param string uid 요청한 유저의 고유 ID.
---@param bool force 강제로 상태를 다시 보낼지 여부.
@ExecSpace("Server")
method void FetchServerState(string uid, bool force)
    -- body
end

---@description "유저 퇴장 이벤트를 처리합니다. 핸들러도 파라미터를 가질 수 있습니다."
---@param UserLeaveEvent e 퇴장한 유저의 이벤트 객체.
@ExecSpace("ServerOnly")
@EventSender("Service", "UserService")
handler HandleUserLeaveEvent(UserLeaveEvent e)
    -- body
end
`
	docs, err := document.Parse(testContent)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("--- Properties ---")
	for _, p := range docs.Properties {
		fmt.Printf("  - Name: %s, Desc: %q\n", p.Name, p.Description)
	}

	fmt.Println("\n--- Methods ---")
	for _, m := range docs.Methods {
		fmt.Printf("  - Name: %s, Desc: %q\n", m.Name, m.Description)
		for _, p := range m.Params {
			fmt.Printf("    - Param: %s (%s) - %s\n", p.Name, p.Type, p.Description)
		}
	}

	fmt.Println("\n--- Handlers ---")
	for _, h := range docs.Handlers {
		fmt.Printf("  - Name: %s, Desc: %q\n", h.Name, h.Description)
		for _, p := range h.Params {
			fmt.Printf("    - Param: %s (%s) - %s\n", p.Name, p.Type, p.Description)
		}
	}
}
