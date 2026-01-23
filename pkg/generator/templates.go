package generator

import _ "embed"

//go:embed style.css
var StyleContent string

//go:embed function_doc.tmpl
var DocumentTemplate string

var Badges = map[string]string{
	"ServerOnly":  ` <img src="https://img.shields.io/badge/ServerOnly-da70d6" alt="ServerOnly" style="vertical-align: middle; margin-left: 8px;">`,
	"ClientOnly":  ` <img src="https://img.shields.io/badge/ClientOnly-87ceeb" alt="ClientOnly" style="vertical-align: middle; margin-left: 8px;">`,
	"Server":      ` <img src="https://img.shields.io/badge/Server-ffa500" alt="Server" style="vertical-align: middle; margin-left: 8px;">`,
	"Client":      ` <img src="https://img.shields.io/badge/Client-90ee90" alt="Client" style="vertical-align: middle; margin-left: 8px;">`,
	"All":         ` <img src="https://img.shields.io/badge/All-d3d3d3" alt="All" style="vertical-align: middle; margin-left: 8px;">`, // All 뱃지 추가
	"Entity":      ` <img src="https://img.shields.io/badge/Entity-ff6b6b" alt="Entity" style="vertical-align: middle; margin-left: 8px;">`,
	"LocalPlayer": ` <img src="https://img.shields.io/badge/LocalPlayer-4ecdc4" alt="LocalPlayer" style="vertical-align: middle; margin-left: 8px;">`,
	"Logic":       ` <img src="https://img.shields.io/badge/Logic-95e1d3" alt="Logic" style="vertical-align: middle; margin-left: 8px;">`,
	"Self":        ` <img src="https://img.shields.io/badge/Self-f38181" alt="Self" style="vertical-align: middle; margin-left: 8px;">`,
	"Model":       ` <img src="https://img.shields.io/badge/Model-aa96da" alt="Model" style="vertical-align: middle; margin-left: 8px;">`,
	"Service":     ` <img src="https://img.shields.io/badge/Service-fcbad3" alt="Service" style="vertical-align: middle; margin-left: 8px;">`,
}
