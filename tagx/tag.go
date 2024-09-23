package tagx

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TagI interface {
	Generate() error
}

type object struct {
	name string
	pkg  string
}

type generatedFile struct {
	path    string
	content string
}

type fileObject struct {
	path    string
	pkg     string
	objects map[string]*object
	v1files generatedFile
	v2files generatedFile
}

type tagI struct {
	dirPath     string
	files       []string
	fileObjects map[string]fileObject
	v1InitFile  generatedFile
}

func GenerateV1(dirPath string) (TagI, error) {
	i := &tagI{
		dirPath:     dirPath,
		files:       []string{},
		fileObjects: map[string]fileObject{},
	}
	err := i.scanDir()
	if err != nil {
		return nil, err
	}
	if len(i.files) == 0 {
		return nil, errors.New("tagx: no files found")
	}
	for _, filePath := range i.files {
		err1 := i.scanFile(filePath)
		if err1 != nil {
			return nil, err1
		}
	}
	for path, fileObject1 := range i.fileObjects {
		err1 := i.prepareV1FileContent(path, fileObject1)
		if err1 != nil {
			return nil, err1
		}
	}
	for _, fileObject1 := range i.fileObjects {
		err1 := i.writeToFile(fileObject1.v1files.path, fileObject1.v1files.content)
		if err1 != nil {
			return nil, err1
		}
	}
	{
		err1 := i.prepareV1InitFileContent()
		if err1 != nil {
			return nil, err1
		}
	}
	{
		err1 := i.writeToFile(i.v1InitFile.path, i.v1InitFile.content)
		if err1 != nil {
			return nil, err1
		}
	}
	return i, nil
}

func (i *tagI) scanDir() error {
	// 遍历项目中的所有文件
	err := filepath.Walk(i.dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理 .go 文件
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			i.files = append(i.files, path)
		}
		return nil
	})
	return err
}

func (i *tagI) scanFile(filePath string) error {
	fs := token.NewFileSet()
	// 解析 Go 文件并包含注释
	node, err := parser.ParseFile(fs, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	pkgName := node.Name.Name
	fileObject1 := fileObject{
		path:    filePath,
		pkg:     pkgName,
		objects: map[string]*object{},
	}
	// 遍历 AST，找到结构体及其关联注释
	ast.Inspect(node, func(n ast.Node) bool {
		// 检查是否为类型声明（包括结构体）
		genDecl, ok := n.(*ast.GenDecl)
		if ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				// 检查是否为结构体类型
				_, isStruct := typeSpec.Type.(*ast.StructType)
				if !isStruct {
					continue
				}
				// 检查注释是否包含 @generated
				isTag := func() bool {
					if genDecl.Doc == nil {
						return false
					}
					for _, comment := range genDecl.Doc.List {
						if strings.Contains(comment.Text, "@generated") {
							return true
						}
					}
					return false
				}()
				if !isTag {
					continue
				}
				object1 := object{
					pkg:  pkgName,
					name: typeSpec.Name.Name,
				}
				fileObject1.objects[object1.name] = &object1
			}
		}
		return true
	})
	if len(fileObject1.objects) > 0 {
		i.fileObjects[filePath] = fileObject1
	}
	return nil
}

// fields @generated sql keys mapping1
//
//	var fields = struct {
//		ID              string
//		Balance         string
//		Balance_Balance string
//	}{ID: "_id", Balance: "balance", Balance_Balance: "balance.balance"}
func (i *tagI) prepareV1FileContent(path string, in fileObject) error {
	arr := []map[string]string{}
	dir := filepath.Dir(path)
	v1file := filepath.Join(dir, "mgo_generated_v1.go")
	// 遍历对象并构建 arr
	for _, v := range in.objects {
		arr = append(arr, map[string]string{
			"structName": v.name,
		})
	}

	// 定义模板内容
	tmpl := `package {{ .PackageName }}

  {{ range .Objects }}
	import "github.com/svc0a/mgo/tagx"
    {{ end }}

func init() {
    {{ range .Objects }}
    _ = tagx.DefineType[{{ .structName }}]()
    {{ end }}
}
`

	// 创建模板
	t, err := template.New("v1File").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板时出错: %w", err)
	}
	// 定义数据结构
	data := struct {
		PackageName string
		Objects     []map[string]string
	}{
		PackageName: in.pkg,
		Objects:     arr,
	}

	// 渲染模板并保存到变量
	var generatedCode string
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return fmt.Errorf("渲染模板时出错: %w", err)
	}

	generatedCode = buf.String()
	// 这里可以根据需要使用 generatedCode
	in.v1files = generatedFile{
		path:    v1file, // 或者给定一个临时路径
		content: generatedCode,
	}
	// 更新文件对象
	i.fileObjects[path] = in
	return nil
}
func (i *tagI) Generate() error {
	return nil
}

func (i *tagI) generate(in object) error {
	return nil
}

// writeToFile 函数用于将生成的代码写入文件
func (i *tagI) writeToFile(filename string, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	fmt.Printf("Code generated and saved to %s\n", filename)
	return nil
}

func (i *tagI) prepareV1InitFileContent() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	arr := []map[string]string{}
	// 遍历对象并构建 arr
	for _, v := range i.fileObjects {
		arr = append(arr, map[string]string{
			"pkg": v.pkg,
		})
	}

	// 定义模板内容
	tmpl := `package main

import "github.com/svc0a/mgo/tagx"

func init() {
    {{ range .Objects }}
   {{ .pkg }}.InitMgo()
    {{ end }}
}
`
	// 创建模板
	t, err := template.New("v1File").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板时出错: %w", err)
	}
	// 定义数据结构
	data := struct {
		Objects []map[string]string
	}{
		Objects: arr,
	}

	// 渲染模板并保存到变量
	var generatedCode string
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return fmt.Errorf("渲染模板时出错: %w", err)
	}

	generatedCode = buf.String()
	// 这里可以根据需要使用 generatedCode
	i.v1InitFile = generatedFile{
		path:    exePath, // 或者给定一个临时路径
		content: generatedCode,
	}
	return nil
}
