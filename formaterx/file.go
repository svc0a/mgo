package formaterx

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/svc0a/mgo/filex"
	"github.com/svc0a/mgo/tagx"
	"go/ast"
	"go/token"
	"path/filepath"
)

type fileObject struct {
	path    string
	dir     string
	pkg     string
	content []byte
}

type FileX interface {
	Export() map[string]Result
}

func defineFile(file string, client tagx.Client) FileX {
	fx := &fileX{
		fileObject: fileObject{
			path: file,
		},
		client:  client,
		results: map[string]Result{},
	}
	{
		err := fx.scan()
		if err != nil {
			logrus.Fatal(err)
			return nil
		}
	}
	fx.format()
	return fx
}

type fileX struct {
	fileObject
	node    *ast.File
	client  tagx.Client
	results map[string]Result // map[pkg.struct.property]result
}

func (fx *fileX) format() {
	for k, v := range fx.results {
		v.expected = fx.client.Format(v.value)
		if v.expected != v.value {
			v.err = errors.New(fmt.Sprintf("format error: value=[%v],expect=[%v]", v.value, v.expected))
		}
		fx.results[k] = v
	}
	return
}

func (fx *fileX) Export() map[string]Result {
	return fx.results
}

func (fx *fileX) scan() error {
	{
		dir := filepath.Dir(fx.path)
		fx.dir = dir
	}
	b, err := filex.ParseFile(fx.path, func(tf *token.FileSet, node *ast.File) {
		fx.node = node
		pkgName := fx.node.Name.Name
		fx.pkg = pkgName
		ast.Inspect(node, func(n ast.Node) bool {
			genDecl, ok := n.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				return true
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, field := range structType.Fields.List {
					result := Result{
						file:     fx.path,
						typeName: typeSpec.Name.Name,
						line:     tf.Position(field.Pos()),
						property: func() string {
							if len(field.Names) != 0 {
								return field.Names[0].Name
							}
							expr, ok := field.Type.(*ast.SelectorExpr)
							if !ok {
								return ""
							}
							pkg, ok := expr.X.(*ast.Ident)
							if !ok {
								return ""
							}
							return fmt.Sprintf("%s.%s", pkg.Name, expr.Sel.Name)
						}(),
						value: func() string {
							if field.Tag == nil {
								return ""
							}
							return fx.client.GetTag(field.Tag.Value)
						}(),
					}
					fx.results[fmt.Sprintf("%s.%s", result.typeName, result.property)] = result
				}
			}
			return true
		})
	})
	if err != nil {
		return err
	}
	fx.content = b
	return nil
}
