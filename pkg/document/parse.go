package document

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	reBlockStarter = regexp.MustCompile(`\b(property|method|handler)\b`)

	reDesc        = regexp.MustCompile(`---@description\s*"([^"]+)"`)
	reExecSpace   = regexp.MustCompile(`@ExecSpace\("([^"]+)"\)`)
	reEventSender = regexp.MustCompile(`@EventSender\("([^"]+)",\s*"([^"]+)"\)`)
	reParam       = regexp.MustCompile(`---@param\s+([a-zA-Z_<>|]+)\s+([a-zA-Z0-9_]+)\s*(.*)`)

	rePropertyCore = regexp.MustCompile(`property\s+([a-zA-Z_<>]+)\s+([a-zA-Z0-9_]+)\s*=\s*"([^"]+)"`)
	reMethodCore   = regexp.MustCompile(`method\s+([a-zA-Z_<>]+)\s+([a-zA-Z0-9_]+)\s*\(([^)]*)\)`)
	reHandlerCore  = regexp.MustCompile(`handler\s+([a-zA-Z0-9_]+)\s*\(([^)]*)\)`)
)

func parseCommonAttributes(block string) (desc, execSpace string, params []ParamInfo) {
	if match := reDesc.FindStringSubmatch(block); len(match) > 1 {
		desc = match[1]
	}
	if match := reExecSpace.FindStringSubmatch(block); len(match) > 1 {
		execSpace = match[1]
	} else {
		execSpace = "All" // 기본값
	}

	paramMatches := reParam.FindAllStringSubmatch(block, -1)
	for _, pMatch := range paramMatches {
		params = append(params, ParamInfo{
			Type:        pMatch[1],
			Name:        pMatch[2],
			Description: pMatch[3],
		})
	}
	return
}

func Parse(content string) (*Documentation, error) {
	// 파싱의 안정성을 위해 파일 끝에 더미 키워드를 추가합니다.
	contentStr := string(content) + "\nproperty"

	docs := &Documentation{}

	indices := reBlockStarter.FindAllStringIndex(contentStr, -1)
	if indices == nil {
		return docs, nil
	}

	var blocks []string
	for i, index := range indices {
		start := index[0]
		var end int
		if i+1 < len(indices) {
			end = indices[i+1][0]
		} else {
			end = len(contentStr)
		}

		realStart := start
		if i > 0 {
			prevEnd := indices[i-1][1]
			lineStart := strings.LastIndex(contentStr[:start], "\n") + 1
			if lineStart > prevEnd {
				realStart = lineStart
			}
		} else {
			realStart = 0
		}

		blockText := contentStr[realStart:end]
		blocks = append(blocks, blockText)
	}

	for _, block := range blocks {
		desc, execSpace, params := parseCommonAttributes(block)

		if propMatch := rePropertyCore.FindStringSubmatch(block); len(propMatch) > 0 {
			docs.Properties = append(docs.Properties, PropertyDoc{
				Description: desc, ExecSpace: execSpace,
				Type: propMatch[1], Name: propMatch[2], DefaultValue: propMatch[3],
			})
		} else if methodMatch := reMethodCore.FindStringSubmatch(block); len(methodMatch) > 0 {
			docs.Methods = append(docs.Methods, MethodDoc{
				Description: desc, ExecSpace: execSpace, Params: params,
				ReturnType: methodMatch[1], Name: methodMatch[2],
			})
		} else if handlerMatch := reHandlerCore.FindStringSubmatch(block); len(handlerMatch) > 0 {
			eventSender := ""
			if match := reEventSender.FindStringSubmatch(block); len(match) > 0 {
				eventSender = fmt.Sprintf("%s, %s", match[1], match[2])
			}
			docs.Handlers = append(docs.Handlers, HandlerDoc{
				Description: desc, ExecSpace: execSpace, EventSender: eventSender,
				Name: handlerMatch[1], Params: params,
			})
		}
	}

	return docs, nil
}

func ParseFile(filepath string) (*Documentation, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return Parse(string(content))
}
