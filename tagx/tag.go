package tagx

import (
	"bytes"
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
	"text/template"
)

type TagI interface {
	WithModel(model ...any) TagI
	Generate() error
}

type fileObject struct {
	dir     string
	pkg     string
	objects map[string]object // map[objectName]object
	content string
}

type object struct {
	file   string
	path   string
	pkg    string
	name   string
	types  reflect2.Type
	fields map[string]string
}

type tagI struct {
	dirPath     string
	files       []string
	fileObjects map[string]fileObject // map[dir]fileObject
	failObjects []string
}

func Define(dirPath string) TagI {
	i := &tagI{
		dirPath:     dirPath,
		files:       []string{},
		fileObjects: map[string]fileObject{},
	}
	err := i.scanDir()
	if err != nil {
		logrus.Fatal(err)
		return nil
	}
	if len(i.files) == 0 {
		logrus.Fatal("tagx: no files found")
		return nil
	}
	for _, filePath := range i.files {
		err1 := i.scanFile(filePath)
		if err1 != nil {
			logrus.Fatal(err1)
			return nil
		}
	}
	if len(i.failObjects) != 0 {
		logrus.WithField("reminder", "please add model to generator").Fatal(strings.Join(i.failObjects, "."))
		return nil
	}
	for _, v1 := range i.fileObjects {
		for _, v2 := range v1.objects {
			i.prepareFields(v2)
		}
	}
	for _, o := range i.fileObjects {
		err1 := i.prepareContent(o)
		if err1 != nil {
			logrus.Fatal(err1)
			return nil
		}
	}
	return i
}

func (i *tagI) WithModel(model ...any) TagI {
	return i
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
	dir := filepath.Dir(filePath)
	fileObject1, ok := i.fileObjects[dir]
	if !ok {
		fileObject1 = fileObject{
			dir:     dir,
			pkg:     pkgName,
			objects: map[string]object{},
		}
	}
	// 遍历 AST，找到结构体及其关联注释
	ast.Inspect(node, func(n ast.Node) bool {
		// 检查是否为类型声明（包括结构体）
		genDecl, ok := n.(*ast.GenDecl)
		if ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				typeSpec, ok1 := spec.(*ast.TypeSpec)
				if !ok1 {
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
				objName := fmt.Sprintf("%s.%s", pkgName, typeSpec.Name.Name)
				object1 := object{
					file: filePath,
					path: dir,
					pkg:  pkgName,
					name: typeSpec.Name.Name,
				}
				t, err1 := reflect2.TypeByName(objName)
				if err1 == nil {
					object1.types = t
				} else {
					i.failObjects = append(i.failObjects, fmt.Sprintf("WithModel(%s.%s{})", pkgName, typeSpec.Name.Name))
				}
				fileObject1.objects[object1.name] = object1
			}
		}
		return true
	})
	if len(fileObject1.objects) > 0 {
		i.fileObjects[dir] = fileObject1
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
func (i *tagI) prepareContent(in fileObject) error {
	type TplField struct {
		Key string
		Val string
	}
	type TplObject struct {
		Name   string
		Fields []TplField
	}
	type TplModel struct {
		Pkg     string
		Objects []TplObject
	}

	objs := []TplObject{}
	// 遍历对象并构建 objs
	for _, v := range in.objects {
		tplObject := TplObject{
			Name:   v.name,
			Fields: []TplField{},
		}
		for k1, v1 := range v.fields {
			tplObject.Fields = append(tplObject.Fields, TplField{
				Key: k1,
				Val: v1,
			})
		}
		objs = append(objs, tplObject)
	}

	// 定义模板内容
	tmpl := `package {{ .Pkg }}

    {{ range .Objects }}
		var {{ .Name }}Fields = struct {
			{{ range .Fields }}
				{{ .Key }}              string
			{{ end }}
			Balance         string
			Balance_Balance string
		}{
			{{ range .Fields }}
				{{ .Key }}: 	"{{ .Val }}",
			{{ end }}
			ID: 	"_id",
			Balance: "balance",
			Balance_Balance: "balance.balance",
		}
    {{ end }}
}
`

	// 创建模板
	t, err := template.New("v1File").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板时出错: %w", err)
	}
	// 定义数据结构
	data := TplModel{
		Pkg:     in.pkg,
		Objects: objs,
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return fmt.Errorf("渲染模板时出错: %w", err)
	}
	f := i.fileObjects[in.dir]
	f.content = buf.String()
	i.fileObjects[in.dir] = f
	return nil
}

func (i *tagI) Generate() error {
	for _, f := range i.fileObjects {
		err := i.generate(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *tagI) generate(in fileObject) error {

}

func (i *tagI) prepareFields(v2 object) {
	t := defineByType(v2.types)
	v2.fields = t.Export()
	i.fileObjects[v2.path].objects[v2.name] = v2
}
