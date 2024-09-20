package examples

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestUser(t *testing.T) {
	projectDir := "." // 替换为你的项目目录

	// 遍历项目中的所有文件
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理 .go 文件
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			parseFileForComments(path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through project: %v", err)
	}
}

// 解析文件中的注释
func parseFileForComments(filePath string) {
	// 创建文件集
	fs := token.NewFileSet()

	// 解析 Go 文件，包含注释
	node, err := parser.ParseFile(fs, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Printf("Error parsing file %s: %v", filePath, err)
		return
	}

	fmt.Printf("Comments in file: %s\n", filePath)
	// 遍历并打印注释
	for _, comment := range node.Comments {
		for _, c := range comment.List {
			fmt.Println(c.Text) // 输出注释
		}
	}
}

func TestUser2(t *testing.T) {
	fs := token.NewFileSet()
	filePath := "user.go"

	// 解析 Go 文件并包含注释
	node, err := parser.ParseFile(fs, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历 AST，找到结构体及其关联注释
	ast.Inspect(node, func(n ast.Node) bool {
		// 检查是否为类型声明（包括结构体）
		genDecl, ok := n.(*ast.GenDecl)
		if ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if ok {
					// 检查是否为结构体类型
					_, isStruct := typeSpec.Type.(*ast.StructType)
					if isStruct {
						// 打印结构体名称
						fmt.Printf("Struct: %s\n", typeSpec.Name.Name)

						// 打印关联的注释
						if genDecl.Doc != nil {
							for _, comment := range genDecl.Doc.List {
								fmt.Printf("Comment: %s\n", comment.Text)
							}
						} else {
							fmt.Println("No associated comments.")
						}
					}
				}
			}
		}
		return true
	})
}
