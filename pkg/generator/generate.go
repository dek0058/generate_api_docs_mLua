package generator

import (
	"bytes"
	"fmt"
	"generate_api_docs_mLua/pkg/document"
	"html/template"
	"strings"
)

// TypeLinkInfo는 타입 이름과 해당 타입의 문서 파일 경로를 매핑합니다.
type TypeLinkInfo map[string]string

// FuncTmplData는 function_doc.tmpl 템플릿에 전달될 데이터 구조체입니다.
type FuncTmplData struct {
	ReturnType        string
	FunctionName      string
	FunctionParamsStr template.HTML // HTML을 안전하게 렌더링하기 위해 template.HTML 사용
	BadgeHTML         template.HTML // HTML을 안전하게 렌더링하기 위�� template.HTML 사용
	Description       string
	Params            []document.ParamInfo
}

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
			html, err := renderFunctionDoc(m, typeLinks)
			if err != nil {
				return "", fmt.Errorf("method %s 렌더링 오류: %w", m.Name, err)
			}
			mdBuilder.WriteString(html)
		}
	}

	// Handlers 렌더링
	if len(doc.Handlers) > 0 {
		if len(doc.Methods) > 0 {
			mdBuilder.WriteString("\n\n")
		}
		mdBuilder.WriteString("## Handlers\n\n")
		for _, h := range doc.Handlers {
			html := renderHandlerDoc(h, typeLinks)
			mdBuilder.WriteString(html)
		}
		mdBuilder.WriteString("\n")
	}

	return mdBuilder.String(), nil
}

// renderFunctionDoc은 메서드 문서를 function_doc.tmpl 템플릿을 사용하여 생성합니다.
func renderFunctionDoc(m document.MethodDoc, typeLinks TypeLinkInfo) (string, error) {
	var paramsStrBuilder strings.Builder
	for i, p := range m.Params {
		paramsStrBuilder.WriteString(fmt.Sprintf("%s %s", createLinkForType(p.Type, typeLinks), p.Name))
		if i < len(m.Params)-1 {
			paramsStrBuilder.WriteString(", ")
		}
	}

	badge, _ := Badges[m.ExecSpace]

	data := FuncTmplData{
		ReturnType:        m.ReturnType,
		FunctionName:      m.Name,
		FunctionParamsStr: template.HTML(paramsStrBuilder.String()),
		BadgeHTML:         template.HTML(badge),
		Description:       m.Description,
		Params:            m.Params,
	}

	tmpl, err := template.New("function").Parse(DocumentTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// renderHandlerDoc은 핸들러 문서를 function_doc.tmpl 템플릿을 사용하여 생성합니다.
func renderHandlerDoc(h document.HandlerDoc, typeLinks TypeLinkInfo) string {
	var paramsStrBuilder strings.Builder
	for i, p := range h.Params {
		paramsStrBuilder.WriteString(fmt.Sprintf("%s %s", createLinkForType(p.Type, typeLinks), p.Name))
		if i < len(h.Params)-1 {
			paramsStrBuilder.WriteString(", ")
		}
	}

	badge, _ := Badges[h.ExecSpace]
	if h.EventSenderType != "" {
		eventSenderBadge, _ := Badges[h.EventSenderType]
		badge += eventSenderBadge
	}

	// 핸들러는 템플릿의 파라미터 상세 설명 부분을 커스터마이징해야 할 수 있음.
	// 여기서는 Logic/Service 추가 정보를 위해 별도 처리
	var bodyContent strings.Builder
	if h.Description != "" {
		bodyContent.WriteString(fmt.Sprintf(`<tr><td>%s</td></tr>`, h.Description))
	}

	// EventSender 추가 정보 (Logic, Service)
	if (h.EventSenderType == "Logic" || h.EventSenderType == "Service") && h.EventSenderValue != "" {
		bodyContent.WriteString(fmt.Sprintf(
			`<tr class="param-row"><td><strong>%s:</strong> %s</td></tr>`,
			h.EventSenderType, h.EventSenderValue,
		))
	}

	// 파라미터 설명
	for _, p := range h.Params {
		if p.Description != "" {
			bodyContent.WriteString(fmt.Sprintf(
				`<tr class="param-row"><td><code class="param-name">%s</code><span class="param-desc"> &nbsp;|&nbsp; %s</span></td></tr>`,
				p.Name, p.Description,
			))
		}
	}

	// 핸들러는 반환 타입이 없을 수도 있음
	var returnTypeSpan string
	if h.ReturnType != "" {
		returnTypeSpan = fmt.Sprintf(`<span class="return-type">%s</span> `, h.ReturnType)
	}

	header := fmt.Sprintf(`%s<span class="function-name">%s</span>(%s)%s`,
		returnTypeSpan, h.Name, paramsStrBuilder.String(), badge)

	table := fmt.Sprintf(`<table class="doc-table"><thead><tr><th>%s</th></tr></thead><tbody>%s</tbody></table>`, header, bodyContent.String())
	return table
}

// createLinkForType은 타입 이름에 적절한 HTML(링크 또는 span)을 적용합니다.
func createLinkForType(typeName string, typeLinks TypeLinkInfo) string {
	baseType := strings.TrimSuffix(strings.TrimSpace(strings.Split(typeName, ",")[0]), ">")

	if link, ok := typeLinks[baseType]; ok {
		// 링크가 있으면 a 태그로 감싸고 class 추가
		return fmt.Sprintf(`<a href="%s" class="param-type">%s</a>`, link, typeName)
	}
	// 링크가 없으면 span으로 감싸서 스타일만 적용
	return fmt.Sprintf(`<span class="param-type">%s</span>`, typeName)
}
