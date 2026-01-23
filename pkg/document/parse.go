package document

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	reDocType = regexp.MustCompile(`^@(Logic|Component|Event|Struct|BTNode|Item|State)\b`)

	reDesc        = regexp.MustCompile(`---@description\s*"([^"]+)"`)
	reExecSpace   = regexp.MustCompile(`@ExecSpace\("([^"]+)"\)`)
	reEventSender = regexp.MustCompile(`@EventSender\("([^"]+)",\s*"([^"]+)"\)`)
	reParam       = regexp.MustCompile(`---@param\s+([a-zA-Z_<>|]+)\s+([a-zA-Z0-9_]+)\s*(.*)`)

	// `readonly` 키워드를 선택적으로 포함하도록 수정
	rePropertyCore = regexp.MustCompile(`(?:readonly\s+)?property\s+([a-zA-Z_<>]+)\s+([a-zA-Z0-9_]+)\s*=\s*"?([^"]+)"?`)
	reMethodCore   = regexp.MustCompile(`method\s+([a-zA-Z_<>]+)\s+([a-zA-Z0-9_]+)\s*\(([^)]*)\)`)
	reHandlerCore  = regexp.MustCompile(`handler\s+([a-zA-Z0-9_]+)\s*\(([^)]*)\)`)
)

func parseCommonAttributes(commentBlock string) (desc, execSpace string, params []ParamInfo) {
	if match := reDesc.FindStringSubmatch(commentBlock); len(match) > 1 {
		desc = match[1]
	}
	if match := reExecSpace.FindStringSubmatch(commentBlock); len(match) > 1 {
		execSpace = match[1]
	}

	paramMatches := reParam.FindAllStringSubmatch(commentBlock, -1)
	for _, pMatch := range paramMatches {
		params = append(params, ParamInfo{
			Type:        pMatch[1],
			Name:        pMatch[2],
			Description: strings.TrimSpace(pMatch[3]),
		})
	}
	return
}

func Parse(content string) (*Documentation, error) {
	docs := &Documentation{}
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")

	var commentBlock []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if docs.DocType == "" {
			if match := reDocType.FindStringSubmatch(trimmedLine); len(match) > 1 {
				docs.DocType = match[1]
				continue
			}
		}

		// 문서 주석 라인 (---@ 또는 @로 시작)을 수집
		isCommentLine := strings.HasPrefix(trimmedLine, "---@") || strings.HasPrefix(trimmedLine, "@ExecSpace")
		if isCommentLine {
			commentBlock = append(commentBlock, trimmedLine)
			continue
		}

		// 실제 코드 라인(property, method, handler)을 만나면 블록 처리
		isCodeLine := rePropertyCore.MatchString(trimmedLine) ||
			reMethodCore.MatchString(trimmedLine) ||
			reHandlerCore.MatchString(trimmedLine)

		if isCodeLine {
			// 현재까지 수집된 주석과 코드 라인을 합쳐서 파싱
			parseBlock(strings.Join(commentBlock, "\n"), trimmedLine, docs)
			// 다음 블록을 위해 주석 블록을 초기화
			commentBlock = nil
		}
		// 관련 없는 라인은 무시
	}

	return docs, nil
}

// parseBlock은 수집된 주석 블록과 실제 코드 한 줄을 받아 처리합니다.
func parseBlock(comment string, code string, docs *Documentation) {
	desc, execSpace, params := parseCommonAttributes(comment)

	if propMatch := rePropertyCore.FindStringSubmatch(code); len(propMatch) > 0 {
		docs.Properties = append(docs.Properties, PropertyDoc{
			Description:  desc,
			ExecSpace:    execSpace,
			Type:         propMatch[1],
			Name:         propMatch[2],
			DefaultValue: strings.Trim(propMatch[3], `"`),
		})
	} else if methodMatch := reMethodCore.FindStringSubmatch(code); len(methodMatch) > 0 {
		// method에 붙은 @ExecSpace는 주석이 아닌 코드 라인과 붙어있을 수 있음
		if execSpace == "" {
			if match := reExecSpace.FindStringSubmatch(code); len(match) > 1 {
				execSpace = match[1]
			}
		}
		docs.Methods = append(docs.Methods, MethodDoc{
			Description: desc,
			ExecSpace:   execSpace,
			Params:      params,
			ReturnType:  methodMatch[1],
			Name:        methodMatch[2],
		})
	} else if handlerMatch := reHandlerCore.FindStringSubmatch(code); len(handlerMatch) > 0 {
		// handler도 마찬가지
		if execSpace == "" {
			if match := reExecSpace.FindStringSubmatch(code); len(match) > 1 {
				execSpace = match[1]
			}
		}
		eventSender := ""
		if match := reEventSender.FindStringSubmatch(comment); len(match) > 0 {
			eventSender = fmt.Sprintf("%s, %s", match[1], match[2])
		}
		docs.Handlers = append(docs.Handlers, HandlerDoc{
			Description: desc,
			ExecSpace:   execSpace,
			EventSender: eventSender,
			Name:        handlerMatch[1],
			Params:      params,
		})
	}
}

func ParseFile(filepath string) (*Documentation, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return Parse(string(content))
}
