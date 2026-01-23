package main

import (
	"fmt"
	"generate_api_docs_mLua/pkg/document"
	"generate_api_docs_mLua/pkg/generator"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rootDir := "RootDesk/MyDesk"
	outputDir := "document/api"

	// 1. 모든 mlua 파일을 찾고 파싱 준비
	filesToParse, err := findLuaFiles(rootDir)
	if err != nil {
		fmt.Printf("파일 검색 중 오류 발생: %v\n", err)
		return
	}

	// 2. 모든 파일을 파싱하여 문서 정보와 타입-파일 경로 맵 생성
	docs := make(map[string]*document.Documentation)
	typeLinks := make(generator.TypeLinkInfo)

	for _, file := range filesToParse {
		doc, err := document.ParseFile(file)
		if err != nil {
			fmt.Printf("파일 파싱 오류 %s: %v\n", file, err)
			continue
		}

		baseName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		docs[file] = doc

		// 타입과 문서 파일 경로를 매핑. e.g., "Channel" -> "./../struct/Channel.md"
		if doc.DocType == "Event" || doc.DocType == "Struct" {
			outPath := getOutputPath(doc, outputDir, baseName+".md")
			// 상대 경로 계산을 위해 임시로 기준 파일 경로 제공
			refPath := getOutputPath(doc, outputDir, "ref.md")
			relPath, _ := filepath.Rel(filepath.Dir(refPath), outPath)
			typeLinks[baseName] = strings.ReplaceAll(relPath, "\\", "/") // 경로 구분자를 URL 친화적으로 변경
		}
	}

	// 3. 각 문서에 대해 Markdown 생성 및 파일 저장
	for file, doc := range docs {
		baseName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		outPath := getOutputPath(doc, outputDir, baseName+".md")

		// Markdown 내용 생성 (CSS 전달 필요 없음)
		mdContent, err := generator.Generate(doc, typeLinks)
		if err != nil {
			fmt.Printf("문서 생성 오류 %s: %v\n", file, err)
			continue
		}

		// 디렉토리 생성 및 파일 쓰기
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			fmt.Printf("디렉토리 생성 오류 %s: %v\n", filepath.Dir(outPath), err)
			continue
		}
		if err := os.WriteFile(outPath, []byte(mdContent), 0644); err != nil {
			fmt.Printf("파일 쓰기 오류 %s: %v\n", outPath, err)
			continue
		}
		fmt.Printf("문서 생성 완료: %s\n", outPath)
	}

	fmt.Println("모든 문서 생성이 완료되었습니다.")
}

// findLuaFiles는 지정된 루트 디렉토리에서 모든 .mlua 파일을 찾습니다.
func findLuaFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".mlua") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// getOutputPath는 문서 타입에 따라 최종 파일 경로를 결정합니다.
func getOutputPath(doc *document.Documentation, baseDir, fileName string) string {
	docTypeDir := "etc"
	if doc.DocType != "" {
		docTypeDir = strings.ToLower(doc.DocType)
	}
	return filepath.Join(baseDir, docTypeDir, fileName)
}
