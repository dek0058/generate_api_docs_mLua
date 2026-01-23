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

type FuncTmplData struct {
	ReturnType        string
	FunctionName      string
	FunctionParamsStr template.HTML
	BadgeHTML         template.HTML
	Description       string
	Params            []document.ParamInfo
}

func Generate(doc *document.Documentation, docTitle, sourceLink string, typeLinks TypeLinkInfo) (string, error) {
	var mdBuilder strings.Builder

	// CSS 스타일은 GitHub에서 지원하지 않으므로 제거
	// 문서 제목에 원본 파일 링크 추가
	mdBuilder.WriteString(fmt.Sprintf("# [%s](%s)\n\n", docTitle, sourceLink))

	// Properties 렌더링
	if len(doc.Properties) > 0 {
		mdBuilder.WriteString("## Properties\n\n")
		// 테이블에 inline style 적용
		mdBuilder.WriteString(`<table style="width: 100%; border-collapse: collapse; border: 1px solid #ccc; margin-bottom: 16px;">`)
		mdBuilder.WriteString(`<thead><tr>`)
		mdBuilder.WriteString(`<th style="background-color: #f0f0f0; padding: 10px 5px; text-align: left; vertical-align: top;">Property</th>`)
		mdBuilder.WriteString(`<th style="background-color: #f0f0f0; padding: 10px 5px; text-align: left; vertical-align: top;">Type</th>`)
		mdBuilder.WriteString(`<th style="background-color: #f0f0f0; padding: 10px 5px; text-align: left; vertical-align: top;">Description</th>`)
		mdBuilder.WriteString(`</tr></thead><tbody>`)
		for _, p := range doc.Properties {
			badge, _ := Badges[p.ExecSpace]
			desc := p.Description
			if p.DefaultValue != "" {
				desc += fmt.Sprintf(" (기본값: `%s`)", p.DefaultValue)
			}
			mdBuilder.WriteString(fmt.Sprintf(
				`<tr><td style="background-color: #fff; padding: 10px 5px; text-align: left; vertical-align: top;"><strong>%s</strong>%s</td><td style="background-color: #fff; padding: 10px 5px; text-align: left; vertical-align: top;"><code>%s</code></td><td style="background-color: #fff; padding: 10px 5px; text-align: left; vertical-align: top;">%s</td></tr>`,
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

	tmpl, err := template.New("function").Parse(DocumentTemplateInline)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

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

	// 핸들러는 반환 타입이 없을 수도 있음
	var returnTypeSpan string
	if h.ReturnType != "" {
		returnTypeSpan = fmt.Sprintf(`<span style="color: #3167ad;">%s</span> `, h.ReturnType)
	}

	// 헤더 생성
	header := fmt.Sprintf(`%s<span style="font-weight: bold;">%s</span>(%s)%s`,
		returnTypeSpan, h.Name, paramsStrBuilder.String(), badge)

	// 본문 내용 생성
	var bodyContent strings.Builder

	if h.Description != "" {
		bodyContent.WriteString(fmt.Sprintf(
			`<tr><td style="background-color: #fff; padding: 10px 5px; text-align: left; vertical-align: top;">%s</td></tr>`,
			h.Description,
		))
	}

	// EventSender 추가 정보 (Logic, Service)
	if (h.EventSenderType == "Logic" || h.EventSenderType == "Service") && h.EventSenderValue != "" {
		bodyContent.WriteString(fmt.Sprintf(
			`<tr><td style="background-color: #fafafa; border-top: 1px solid #eee; padding: 10px 5px 10px 15px; text-align: left; vertical-align: top;"><strong>%s:</strong> %s</td></tr>`,
			h.EventSenderType, h.EventSenderValue,
		))
	}

	// 파라미터 설명
	for _, p := range h.Params {
		if p.Description != "" {
			bodyContent.WriteString(fmt.Sprintf(
				`<tr><td style="background-color: #fafafa; border-top: 1px solid #eee; padding: 10px 5px 10px 15px; text-align: left; vertical-align: top;"><code style="background-color: #e1e4e8; padding: 2px 5px; border-radius: 4px; font-family: monospace;">%s</code><span style="color: #57606a;"> &nbsp;|&nbsp; %s</span></td></tr>`,
				p.Name, p.Description,
			))
		}
	}

	// 완전한 테이블 생성
	table := fmt.Sprintf(
		`<table style="width: 100%%; border-collapse: collapse; border: 1px solid #ccc; margin-bottom: 16px;"><thead><tr><th style="background-color: #f0f0f0; padding: 10px 5px; text-align: left; vertical-align: top;">%s</th></tr></thead><tbody>%s</tbody></table>`,
		header, bodyContent.String(),
	)

	return table
}

func createLinkForType(typeName string, typeLinks TypeLinkInfo) string {
	baseType := strings.TrimSuffix(strings.TrimSpace(strings.Split(typeName, ",")[0]), ">")

	if link, ok := typeLinks[baseType]; ok {
		// ��크가 있으면 a 태그로 감싸고 inline style 추가
		return fmt.Sprintf(`<a href="%s" style="text-decoration: none; color: #3167ad;">%s</a>`, link, typeName)
	}
	// 링크가 없으면 span으로 감싸서 스타일만 적용
	return fmt.Sprintf(`<span style="color: #3167ad;">%s</span>`, typeName)
}
