package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func writeFile(filename string, content []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func parseFile(file string, cb func(tf *token.FileSet, node *ast.File)) ([]byte, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return parse(b, cb)
}

func parse(b []byte, cb func(tf *token.FileSet, node *ast.File)) ([]byte, error) {
	fileSet1 := token.NewFileSet()
	node, err := parser.ParseFile(fileSet1, "", b, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	cb(fileSet1, node)
	// 使用 bytes.Buffer 将内容写入内存
	var buf bytes.Buffer
	if err1 := printer.Fprint(&buf, fileSet1, node); err1 != nil {
		return nil, fmt.Errorf("error printing AST to buffer: %w", err1)
	}
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("格式化代码时出错: %w", err)
	}
	return formattedCode, nil
}
