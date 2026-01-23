package generator

import _ "embed"

//go:embed style.css
var StyleContent string

//go:embed function_doc.tmpl
var DocumentTemplate string

var Badges = map[string]string{
	"ServerOnly": ` <img src="https://img.shields.io/badge/ServerOnly-da70d6" alt="ServerOnly" style="vertical-align: middle; margin-left: 8px;">`,
	"ClientOnly": ` <img src="https://img.shields.io/badge/ClientOnly-87ceeb" alt="ClientOnly" style="vertical-align: middle; margin-left: 8px;">`,
	"Server":     ` <img src="https://img.shields.io/badge/Server-ffa500" alt="Server" style="vertical-align: middle; margin-left: 8px;">`,
	"Client":     ` <img src="https://img.shields.io/badge/Client-90ee90" alt="Client" style="vertical-align: middle; margin-left: 8px;">`,
	"All":        ` <img src="https://img.shields.io/badge/All-d3d3d3" alt="All" style="vertical-align: middle; margin-left: 8px;">`, // All 뱃지 추가
}
