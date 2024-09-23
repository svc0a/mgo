package tagx

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/svc0a/reflect2"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TagI interface {
	Generate() error
}

type object struct {
	name  string
	file  string
	types reflect2.Type
	code  string
}

type tagI struct {
	dirPath string
	files   []string
	objects map[string]object // map[objectName]object
}

func Define(dirPath string) (TagI, error) {
	i := &tagI{
		dirPath: dirPath,
		files:   []string{},
		objects: map[string]object{},
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
	for _, o := range i.objects {
		err := i.prepareCode(o)
		if err != nil {
			return nil, err
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
				structType, isStruct := typeSpec.Type.(*ast.StructType)
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
				objName := fmt.Sprintf("%s.%s", pkgName, typeSpec.Name.Name)
				object1 := object{
					name: objName,
					file: filePath,
				}
				// 扫描字段，检查是否为组合字段
				for _, field := range structType.Fields.List {
					fieldType := fmt.Sprintf("%s", field.Type)
					logrus.Info(fieldType)
					// 检查字段类型是否是结构体（组合）
					if _, ok := field.Type.(*ast.StructType); ok {
						logrus.WithField("field", field.Names[0].Name).Info("Composite field found")
					} else {
						logrus.WithField("field", field.Names[0].Name).Info("Non-composite field found")
					}
				}
				//t, err1 := reflect2.TypeByName(objName)
				//if err1 == nil {
				//	object1.types = t
				//}
				i.objects[objName] = object1
			}
		}
		return true
	})
	return nil
}

// fields @generated sql keys mapping1
//
//	var fields = struct {
//		ID              string
//		Balance         string
//		Balance_Balance string
//	}{ID: "_id", Balance: "balance", Balance_Balance: "balance.balance"}
func (i *tagI) prepareCode(in object) error {
	if in.types == nil {
		return nil
	}
	m := defineByType(in.types).Export()
	code := ""
	{
		for k, v := range m {
			logrus.WithField("k", k).WithField("v", v).Info()
		}
	}
	in.code = code
	i.objects[in.name] = in
	return nil
}

func (i *tagI) Generate() error {
	for _, o := range i.objects {
		err := i.generate(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *tagI) generate(in object) error {
	// 以追加模式打开文件
	f, err := os.OpenFile(in.file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err1 := f.Close()
		if err1 != nil {
			logrus.Error(err1)
		}
	}(f)
	// 写入代码
	if _, err := f.WriteString(in.code); err != nil {
		return err
	}

	logrus.WithField("file", in.file).Info("Code appended successfully.")
	return nil
}
