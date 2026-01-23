package generator

import (
	"bytes"
	"fmt"
	"generate_api_docs_mLua/pkg/document"
	tmpl_pkg "generate_api_docs_mLua/pkg/template"
	"html"
	"strings"
	"text/template"
)

// TypeLinkInfo는 타입 이름과 해당 타입의 문서 파일 경로를 매핑합니다.
type TypeLinkInfo map[string]string

type TemplateData struct {
	ReturnType        string
	FunctionName      string
	FunctionParamsStr string
	BadgeHTML         string
	Description       string
	Params            []ParamTemplateData
}

type ParamTemplateData struct {
	Name string
	Desc string
}

func Generate(doc *document.Documentation, typeLinks TypeLinkInfo, cssContent string) (string, error) {
	var mdBuilder strings.Builder

	// CSS 스타일 추가
	mdBuilder.WriteString("<style>\n")
	mdBuilder.WriteString(cssContent)
	mdBuilder.WriteString("</style>\n\n")

	// 문서 타입과 파일 이름 헤더 추가
	if doc.DocType != "" {
		mdBuilder.WriteString(fmt.Sprintf("# %s\n\n", doc.DocType))
	}

	// Properties 렌더링
	if len(doc.Properties) > 0 {
		mdBuilder.WriteString("## Properties\n\n")
		for _, p := range doc.Properties {
			mdBuilder.WriteString(fmt.Sprintf("- **%s** (`%s`): %s", p.Name, p.Type, p.Description))
			if p.DefaultValue != "" {
				mdBuilder.WriteString(fmt.Sprintf(" (기본값: `%s`)", p.DefaultValue))
			}
			mdBuilder.WriteString("\n")
		}
		mdBuilder.WriteString("\n")
	}

	// Methods 렌더링
	if len(doc.Methods) > 0 {
		mdBuilder.WriteString("## Methods\n\n")
		for _, m := range doc.Methods {
			html, err := renderFunctionDoc(m.Name, m.ReturnType, m.Description, m.ExecSpace, m.Params, typeLinks)
			if err != nil {
				return "", fmt.Errorf("method %s 렌더링 실패: %w", m.Name, err)
			}
			mdBuilder.WriteString(html)
		}
	}

	// Handlers 렌더링
	if len(doc.Handlers) > 0 {
		mdBuilder.WriteString("## Handlers\n\n")
		for _, h := range doc.Handlers {
			html, err := renderFunctionDoc(h.Name, "", h.Description, h.ExecSpace, h.Params, typeLinks)
			if err != nil {
				return "", fmt.Errorf("handler %s 렌더링 실패: %w", h.Name, err)
			}
			mdBuilder.WriteString(html)
		}
	}

	return mdBuilder.String(), nil
}

func renderFunctionDoc(name, returnType, desc, execSpace string, params []document.ParamInfo, typeLinks TypeLinkInfo) (string, error) {
	tmpl, err := template.New("doc").Parse(tmpl_pkg.DocumentTemplate)
	if err != nil {
		return "", err
	}

	// 파라미터 문자열 및 템플릿 데이터 생성
	var paramStrs []string
	var paramTemplates []ParamTemplateData
	for _, p := range params {
		linkedType := createLinkForType(p.Type, typeLinks)
		paramStrs = append(paramStrs, fmt.Sprintf("%s %s", linkedType, p.Name))
		paramTemplates = append(paramTemplates, ParamTemplateData{
			Name: p.Name,
			Desc: fmt.Sprintf("`%s` | %s", linkedType, html.EscapeString(p.Description)),
		})
	}

	badge, exists := tmpl_pkg.Badges[execSpace]
	if !exists && execSpace != "All" {
		badge = tmpl_pkg.Badges["Server"] // 기본값
	}

	data := TemplateData{
		ReturnType:        createLinkForType(returnType, typeLinks),
		FunctionName:      name,
		FunctionParamsStr: strings.Join(paramStrs, ", "),
		BadgeHTML:         badge,
		Description:       desc,
		Params:            paramTemplates,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// createLinkForType은 주어진 타입에 대한 마크다운 링크를 생성합니다.
func createLinkForType(typeName string, typeLinks TypeLinkInfo) string {
	// "table<string, Channel>" 같은 복합 타입 처리
	baseType := strings.TrimRight(strings.TrimLeft(typeName, "table<"), ">")
	if strings.Contains(baseType, ",") {
		parts := strings.Split(baseType, ",")
		baseType = strings.TrimSpace(parts[1]) // value 타입에 대해서만 링크 생성
	}

	if link, ok := typeLinks[baseType]; ok {
		return fmt.Sprintf("[%s](%s)", typeName, link)
	}
	return typeName // 링크가 없으면 타입 이름만 반환
}
