package generator

import (
	"fmt"
	"generate_api_docs_mLua/pkg/document"
	"strings"
)

// TypeLinkInfo는 타입 이름과 해당 타입의 문서 파일 경로를 매핑합니다.
type TypeLinkInfo map[string]string

func Generate(doc *document.Documentation, docTitle, sourceLink string, typeLinks TypeLinkInfo) (string, error) {
	var mdBuilder strings.Builder

	// 임베드된 CSS 스타일 추가
	mdBuilder.WriteString("<style>\n")
	mdBuilder.WriteString(StyleContent)
	mdBuilder.WriteString("</style>\n\n")

	// 문서 제목에 원본 파일 링크 추가
	mdBuilder.WriteString(fmt.Sprintf("# [%s](%s)\n\n", docTitle, sourceLink))

	// Properties 렌더링
	if len(doc.Properties) > 0 {
		mdBuilder.WriteString("## Properties\n\n")
		mdBuilder.WriteString(`<table class="doc-table"><thead><tr><th>Property</th><th>Type</th><th>Description</th></tr></thead><tbody>`)
		for _, p := range doc.Properties {
			badge, _ := Badges[p.ExecSpace]
			desc := p.Description
			if p.DefaultValue != "" {
				desc += fmt.Sprintf(" (기본값: `%s`)", p.DefaultValue)
			}
			mdBuilder.WriteString(fmt.Sprintf(
				`<tr><td><strong>%s</strong>%s</td><td><code>%s</code></td><td>%s</td></tr>`,
				p.Name, badge, p.Type, desc,
			))
		}
		mdBuilder.WriteString(`</tbody></table>`)
		mdBuilder.WriteString("\n\n")
	}

	// Methods 렌더링
	if len(doc.Methods) > 0 {
		mdBuilder.WriteString("## Methods\n\n")
		for _, m := range doc.Methods {
			html := renderFunctionDoc(m.Name, m.ReturnType, m.Description, m.ExecSpace, m.Params, typeLinks)
			mdBuilder.WriteString(html)
		}
		mdBuilder.WriteString("\n\n")
	}

	// Handlers 렌더링
	if len(doc.Handlers) > 0 {
		mdBuilder.WriteString("## Handlers\n\n")
		for _, h := range doc.Handlers {
			html := renderHandlerDoc(h, typeLinks)
			mdBuilder.WriteString(html)
		}
		mdBuilder.WriteString("\n")
	}

	return mdBuilder.String(), nil
}

// renderFunctionDoc은 함수/핸들러 문서를 인라인 스타일 HTML로 생성합니다.
func renderFunctionDoc(name, returnType, desc, execSpace string, params []document.ParamInfo, typeLinks TypeLinkInfo) string {
	var paramsStrBuilder strings.Builder
	for i, p := range params {
		linkedType := createLinkForType(p.Type, typeLinks)
		paramsStrBuilder.WriteString(fmt.Sprintf("%s %s", linkedType, p.Name))
		if i < len(params)-1 {
			paramsStrBuilder.WriteString(", ")
		}
	}

	badge, _ := Badges[execSpace]

	var bodyContent string
	if desc != "" {
		bodyContent += fmt.Sprintf(`<tr><td style="background-color:#fff;border-color:#ccc;border-style:solid;border-width:0px;color:#333;overflow:hidden;padding:10px 5px;text-align:left;vertical-align:top;word-break:normal">%s</td></tr>`, desc)
	}

	if len(params) > 0 {
		for _, p := range params {
			// Only show parameter details if there's a description
			if p.Description != "" {
				linkedType := createLinkForType(p.Type, typeLinks)
				bodyContent += fmt.Sprintf(`<tr><td style="background-color:#fafafa; border-top: 1px solid #eee; padding: 10px 5px 10px 15px;"><code style="background-color: #e1e4e8; padding: 2px 5px; border-radius: 4px;">%s</code><span style="color: #57606a;"> &nbsp;|&nbsp; <code>%s</code> | %s</span></td></tr>`, p.Name, linkedType, p.Description)
			}
		}
	}

	return fmt.Sprintf(
		`<table style="width:100%%;border-collapse:collapse; border-color:#ccc;border-spacing:0;border-style:solid;border-width:1px; margin-bottom: 16px;"><thead><tr><th style="background-color:#f0f0f0;border-color:#ccc;border-style:solid;border-width:0px;color:#333;overflow:hidden;padding:10px 5px;text-align:left;vertical-align:top;word-break:normal"><span style="color:#3167ad">%s</span> <span style="font-weight:bold">%s</span>(%s)%s</th></tr></thead><tbody>%s</tbody></table>`,
		returnType, name, paramsStrBuilder.String(), badge, bodyContent)
}

// renderHandlerDoc은 핸들러 문서를 인라인 스타일 HTML로 생성합니다.
// EventSender 배지와 추가 정보를 처리합니다.
func renderHandlerDoc(h document.HandlerDoc, typeLinks TypeLinkInfo) string {
	var paramsStrBuilder strings.Builder
	for i, p := range h.Params {
		linkedType := createLinkForType(p.Type, typeLinks)
		paramsStrBuilder.WriteString(fmt.Sprintf("%s %s", linkedType, p.Name))
		if i < len(h.Params)-1 {
			paramsStrBuilder.WriteString(", ")
		}
	}

	// ExecSpace badge
	badge, _ := Badges[h.ExecSpace]
	
	// EventSender badge
	if h.EventSenderType != "" {
		eventSenderBadge, _ := Badges[h.EventSenderType]
		badge += eventSenderBadge
	}

	var bodyContent string
	
	// Description
	if h.Description != "" {
		bodyContent += fmt.Sprintf(`<tr><td style="background-color:#fff;border-color:#ccc;border-style:solid;border-width:0px;color:#333;overflow:hidden;padding:10px 5px;text-align:left;vertical-align:top;word-break:normal">%s</td></tr>`, h.Description)
	}
	
	// EventSender additional info for Logic and Service types
	if (h.EventSenderType == "Logic" || h.EventSenderType == "Service") && h.EventSenderValue != "" {
		bodyContent += fmt.Sprintf(`<tr><td style="background-color:#f9f9f9;border-color:#ccc;border-style:solid;border-width:0px;color:#333;overflow:hidden;padding:10px 5px;text-align:left;vertical-align:top;word-break:normal"><strong>%s:</strong> %s</td></tr>`, h.EventSenderType, h.EventSenderValue)
	}

	// Parameters
	if len(h.Params) > 0 {
		for _, p := range h.Params {
			// Only show parameter details if there's a description
			if p.Description != "" {
				linkedType := createLinkForType(p.Type, typeLinks)
				bodyContent += fmt.Sprintf(`<tr><td style="background-color:#fafafa; border-top: 1px solid #eee; padding: 10px 5px 10px 15px;"><code style="background-color: #e1e4e8; padding: 2px 5px; border-radius: 4px;">%s</code><span style="color: #57606a;"> &nbsp;|&nbsp; <code>%s</code> | %s</span></td></tr>`, p.Name, linkedType, p.Description)
			}
		}
	}

	// Build the header with return type if present
	var headerContent string
	if h.ReturnType != "" {
		headerContent = fmt.Sprintf(`<span style="color:#3167ad">%s</span> <span style="font-weight:bold">%s</span>(%s)%s`, h.ReturnType, h.Name, paramsStrBuilder.String(), badge)
	} else {
		headerContent = fmt.Sprintf(`<span style="font-weight:bold">%s</span>(%s)%s`, h.Name, paramsStrBuilder.String(), badge)
	}

	return fmt.Sprintf(
		`<table style="width:100%%;border-collapse:collapse; border-color:#ccc;border-spacing:0;border-style:solid;border-width:1px; margin-bottom: 16px;"><thead><tr><th style="background-color:#f0f0f0;border-color:#ccc;border-style:solid;border-width:0px;color:#333;overflow:hidden;padding:10px 5px;text-align:left;vertical-align:top;word-break:normal">%s</th></tr></thead><tbody>%s</tbody></table>`,
		headerContent, bodyContent)
}

func createLinkForType(typeName string, typeLinks TypeLinkInfo) string {
	baseType := typeName
	if strings.Contains(typeName, ",") {
		parts := strings.Split(typeName, ",")
		baseType = strings.TrimSpace(parts[len(parts)-1])
		baseType = strings.TrimSuffix(baseType, ">")
	}

	if link, ok := typeLinks[baseType]; ok {
		return fmt.Sprintf(`[%s](%s)`, typeName, link)
	}
	return typeName
}
