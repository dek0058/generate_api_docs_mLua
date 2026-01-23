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

	filesToParse, err := findLuaFiles(rootDir)
	if err != nil {
		fmt.Printf("파일 검색 중 오류 발생: %v\n", err)
		return
	}

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

		if doc.DocType == "Event" || doc.DocType == "Struct" {
			outPath := getOutputPath(doc, outputDir, baseName+".md")
			refDir := filepath.Join(outputDir, "logic")
			relPath, _ := filepath.Rel(refDir, outPath)
			typeLinks[baseName] = strings.ReplaceAll(relPath, "\\", "/")
		}
	}

	for file, doc := range docs {
		baseName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		outPath := getOutputPath(doc, outputDir, baseName+".md")

		// 원본 mlua 파일에 대한 상대 경로 계산
		relPathToSource, err := filepath.Rel(filepath.Dir(outPath), file)
		if err != nil {
			fmt.Printf("상대 경로 계산 오류: %v\n", err)
			// 실패 시 대체 경로 사용 (루트 기준)
			relPathToSource = file
		}
		// URL 경로 형식으로 변경
		relPathToSource = strings.ReplaceAll(relPathToSource, "\\", "/")

		// Generate 함수에 문서 제목과 원본 파일 링크를 전달
		mdContent, err := generator.Generate(doc, baseName, relPathToSource, typeLinks)
		if err != nil {
			fmt.Printf("문서 생성 오류 %s: %v\n", file, err)
			continue
		}

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

func getOutputPath(doc *document.Documentation, baseDir, fileName string) string {
	docTypeDir := "etc"
	if doc.DocType != "" {
		docTypeDir = strings.ToLower(doc.DocType)
	}
	return filepath.Join(baseDir, docTypeDir, fileName)
}
