package gen

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/svc0a/mgo/filex"
	"github.com/svc0a/mgo/tagx"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

var sourceKey = "source"
var commentLabel = "@qlGenerated"
var reflectImport = "github.com/svc0a/reflect2"

type module struct {
	name string
	dir  string
}

type variable struct {
	source       object
	name         string
	comment      string
	commentIndex int
	codePos      token.Position
}

type fileObject struct {
	dir             string
	pkg             string
	objects         map[string]object
	registerContent []byte
	content         []byte
	path            string
	variables       map[string]variable
}

type object struct {
	path       string
	pkg        string
	name       string
	types      reflect.Type
	fields     map[string]string
	samePkg    bool // same package
	fullName   string
	callerName string
}

type fileImport struct {
	module   module
	filePath string
	dir      string
	tImport  string
}

type Service interface {
	Generate() error
}

type impl struct {
	callerFile    string
	dirPath       string
	files         []string
	fileObjects   map[string]fileObject
	xImports      map[string]*fileImport
	callerContent []byte
	client        tagx.Client
}

type Option func(*impl)

func WithMongodb() Option {
	return func(i *impl) {
		i.client = tagx.Mongo()
	}
}

func WithPostgre() Option {
	return func(i *impl) {
		i.client = tagx.Postgre()
	}
}

func WithSource(source string) Option {
	return func(i *impl) {
		sourceKey = source
	}
}

func WithCommentLabel(label string) Option {
	return func(i *impl) {
		commentLabel = label
	}
}

func WithDir(dir string) Option {
	return func(i *impl) {
		i.dirPath = dir
	}
}

func Define(options ...Option) Service {
	i := &impl{
		dirPath:     "./",
		files:       []string{},
		fileObjects: map[string]fileObject{},
		xImports:    map[string]*fileImport{},
		client:      tagx.Mongo(),
	}
	for _, o := range options {
		o(i)
	}
	{
		err := i.getCallerFile()
		if err != nil {
			logrus.Fatal(err)
			return nil
		}
	}
	{
		err := i.scanDir()
		if err != nil {
			logrus.Fatal(err)
			return nil
		}
		if len(i.files) == 0 {
			logrus.Fatal("no files found")
			return nil
		}
	}
	for _, filePath := range i.files {
		fileX1 := DefineFile(filePath, i.client)
		if fileX1.FileImport() != nil {
			fileImport1 := fileX1.FileImport()
			i.xImports[fileImport1.dir] = fileImport1
		}
		i.fileObjects[filePath] = fileX1.Export()
	}
	if len(i.xImports) != 0 {
		{
			err2 := i.appendImports()
			if err2 != nil {
				logrus.Fatal(err2)
				return nil
			}
		}
		{
			err := i.register()
			if err != nil {
				logrus.Fatal(err)
				return nil
			}
		}
		logrus.Fatal("please run again")
	}
	return i
}

func (svc *impl) scanDir() error {
	err := filepath.Walk(svc.dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			svc.files = append(svc.files, path)
		}
		return nil
	})
	return err
}

func (svc *impl) register() error {
	{
		err := filex.RewriteFile(svc.callerFile, svc.callerContent)
		if err != nil {
			return err
		}
	}
	for _, f := range svc.fileObjects {
		if f.objects == nil || len(f.objects) == 0 || f.registerContent == nil {
			continue
		}
		err := filex.RewriteFile(f.path, f.registerContent)
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *impl) Generate() error {
	for _, f := range svc.fileObjects {
		if f.objects == nil || len(f.objects) == 0 || f.content == nil {
			continue
		}
		err := filex.RewriteFile(f.path, f.content)
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *impl) getCallerFile() error {
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return fmt.Errorf("unable to retrieve caller information")
	}
	svc.callerFile = file
	return nil
}

func (svc *impl) appendImports() error {
	importPaths := []string{}
	for _, tImport1 := range svc.xImports {
		importPaths = append(importPaths, tImport1.tImport)
	}
	b, err := filex.ParseFile(svc.callerFile, func(fileSet1 *token.FileSet, node *ast.File) {
		for _, decl := range node.Decls {
			decl1, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, importPath := range importPaths {
				exists := false
				for _, imp := range decl1.Specs {
					imp1, ok1 := imp.(*ast.ImportSpec)
					if ok1 {
						if strings.Trim(imp1.Path.Value, "\"") == importPath {
							exists = true
							break
						}
					}
				}
				if !exists {
					_ = astutil.AddNamedImport(fileSet1, node, "_", importPath)
				}
			}
		}
	})
	if err != nil {
		return err
	}
	svc.callerContent = b
	return nil
}
