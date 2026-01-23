package generator

import _ "embed"

//go:embed style.css
var StyleContent string

//go:embed function_doc.tmpl
var DocumentTemplate string

// GitHub Markdown용 inline style 템플릿
var DocumentTemplateInline = `<table style="width: 100%; border-collapse: collapse; border: 1px solid #ccc; margin-bottom: 16px;">
    <thead>
        <tr>
            <th style="background-color: #f0f0f0; padding: 10px 5px; text-align: left; vertical-align: top;">
                <span style="color: #3167ad;">{{.ReturnType}}</span> <span style="font-weight: bold;">{{.FunctionName}}</span>({{.FunctionParamsStr}}){{.BadgeHTML}}
            </th>
        </tr>
    </thead>
    <tbody>{{- if .Description}}
        <tr>
            <td style="background-color: #fff; padding: 10px 5px; text-align: left; vertical-align: top;">
                {{.Description}}
            </td>
        </tr>{{- end}}{{- range .Params}}{{- if .Description}}
        <tr>
            <td style="background-color: #fafafa; border-top: 1px solid #eee; padding: 10px 5px 10px 15px; text-align: left; vertical-align: top;">
                <code style="background-color: #e1e4e8; padding: 2px 5px; border-radius: 4px; font-family: monospace;">{{.Name}}</code>
                <span style="color: #57606a;"> &nbsp;|&nbsp; {{.Description}}</span>
            </td>
        </tr>{{- end}}{{- end}}
    </tbody>
</table>
`

var Badges = map[string]string{
	"ServerOnly": ` <img src="https://img.shields.io/badge/ServerOnly-da70d6" alt="ServerOnly" style="vertical-align: middle; margin-left: 8px;">`,
	"ClientOnly": ` <img src="https://img.shields.io/badge/ClientOnly-87ceeb" alt="ClientOnly" style="vertical-align: middle; margin-left: 8px;">`,
	"Server":     ` <img src="https://img.shields.io/badge/Server-ffa500" alt="Server" style="vertical-align: middle; margin-left: 8px;">`,
	"Client":     ` <img src="https://img.shields.io/badge/Client-90ee90" alt="Client" style="vertical-align: middle; margin-left: 8px;">`,
	"Logic":      ` <img src="https://img.shields.io/badge/Logic-95e1d3" alt="Logic" style="vertical-align: middle; margin-left: 8px;">`,
	"Service":    ` <img src="https://img.shields.io/badge/Service-f38181" alt="Service" style="vertical-align: middle; margin-left: 8px;">`,
}
