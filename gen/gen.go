package gen

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"github.com/svc0a/mgo/tagx"
	"github.com/svc0a/reflect2"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	sourceKey     = "source"
	commentLabel  = "@mongoGenerated"
	reflectImport = "github.com/svc0a/reflect2"
)

type GenI interface {
	Generate() error
}

type module struct {
	name string
	dir  string
}

type variable struct {
	source       object
	name         string
	comment      string
	commentIndex int
	codePos      token.Pos
	commentPos   token.Pos
	shift        token.Pos
}

type fileObject struct {
	dir             string
	pkg             string
	objects         map[string]object
	registerContent string
	content         string
	path            string
	variables       map[string]variable
}

type object struct {
	path       string
	pkg        string
	name       string
	types      reflect2.Type
	fields     map[string]string
	samePkg    bool // same package
	fullName   string
	callerName string
}

type tImport struct {
	module   module
	filePath string
	dir      string
	tImport  string
}

type genI struct {
	callerFile    string
	dirPath       string
	files         []string
	fileObjects   map[string]fileObject
	xImports      map[string]tImport
	callerContent string
}

func Define(dirPath string) GenI {
	i := &genI{
		dirPath:     dirPath,
		files:       []string{},
		fileObjects: map[string]fileObject{},
		xImports:    map[string]tImport{},
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
			logrus.Fatal("tagx: no files found")
			return nil
		}
	}
	for _, filePath := range i.files {
		err1 := i.scanFile(filePath)
		if err1 != nil {
			logrus.Fatal(err1)
			return nil
		}
	}
	if len(i.xImports) != 0 {
		for _, v1 := range i.xImports {
			err2 := i.getModuleName(v1)
			if err2 != nil {
				logrus.Fatal(err2)
				return nil
			}
		}
		{
			err2 := i.appendImports()
			if err2 != nil {
				logrus.Fatal(err2)
				return nil
			}
		}
		{
			for _, o := range i.fileObjects {
				err1 := i.prepareRegisterContent(o)
				if err1 != nil {
					logrus.Fatal(err1)
					return nil
				}
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
	{
		for _, o := range i.fileObjects {
			err1 := i.prepareContent(o)
			if err1 != nil {
				logrus.Fatal(err1)
				return nil
			}
		}
	}
	return i
}

func (g *genI) scanDir() error {
	err := filepath.Walk(g.dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			g.files = append(g.files, path)
		}
		return nil
	})
	return err
}

func (g *genI) scanFile(filePath1 string) error {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filePath1, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	pkgName := node.Name.Name
	dir := filepath.Dir(filePath1)
	fileObject1, ok := g.fileObjects[filePath1]
	if !ok {
		fileObject1 = fileObject{
			path:      filePath1,
			dir:       dir,
			pkg:       pkgName,
			objects:   map[string]object{},
			variables: map[string]variable{},
		}
	}
	commentIndex := 0
	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				isTag := func() bool {
					if genDecl.Doc == nil {
						return false
					}
					for _, comment := range genDecl.Doc.List {
						if strings.Contains(comment.Text, commentLabel) {
							return true
						}
					}
					return false
				}()
				if !isTag {
					continue
				}
				valueSpec, ok1 := spec.(*ast.ValueSpec)
				if !ok1 {
					continue
				}
				for _, val := range valueSpec.Values {
					compositeLit1, ok2 := val.(*ast.CompositeLit)
					if ok2 {
						type1 := compositeLit1.Type
						structType1, ok3 := type1.(*ast.StructType)
						if ok3 {
							for _, field := range structType1.Fields.List {
								for _, name1 := range field.Names {
									if name1.Name == sourceKey {
										decl1 := name1.Obj.Decl
										decl2, ok4 := decl1.(*ast.Field)
										if ok4 {
											type2 := decl2.Type
											{
												ident1, ok5 := type2.(*ast.Ident)
												if ok5 {
													{
														obj := object{
															fullName:   fmt.Sprintf("%s.%s", pkgName, ident1.Name),
															pkg:        pkgName,
															path:       filePath1,
															name:       ident1.Name,
															samePkg:    true,
															callerName: valueSpec.Names[0].Name,
														}
														t, err2 := reflect2.TypeByName(obj.fullName)
														if err2 == nil {
															obj.types = t
														} else {
															g.xImports[dir] = tImport{
																dir:      dir,
																filePath: filePath1,
															}
														}
														fileObject1.objects[obj.fullName] = obj
														{
															fileObject1.variables[valueSpec.Names[0].Name] = variable{
																source:       obj,
																name:         valueSpec.Names[0].Name,
																comment:      genDecl.Doc.List[0].Text,
																codePos:      valueSpec.Pos(),
																commentPos:   genDecl.Doc.List[0].Slash,
																commentIndex: commentIndex,
															}
															commentIndex++
														}
													}
												}
											}
											{
												selectorExpr1, ok5 := type2.(*ast.SelectorExpr)
												if ok5 {
													x1, ok6 := selectorExpr1.X.(*ast.Ident)
													if ok6 {
														obj := object{
															fullName:   fmt.Sprintf("%s.%s", x1.Name, selectorExpr1.Sel.Name),
															pkg:        x1.Name,
															path:       filePath1,
															name:       selectorExpr1.Sel.Name,
															samePkg:    false,
															callerName: valueSpec.Names[0].Name,
														}
														t, err2 := reflect2.TypeByName(obj.fullName)
														if err2 == nil {
															obj.types = t
														} else {
															g.xImports[dir] = tImport{
																dir:      dir,
																filePath: filePath1,
															}
														}
														fileObject1.objects[obj.fullName] = obj
														{
															fileObject1.variables[valueSpec.Names[0].Name] = variable{
																source:       obj,
																name:         valueSpec.Names[0].Name,
																comment:      genDecl.Doc.List[0].Text,
																codePos:      valueSpec.Pos(),
																commentPos:   genDecl.Doc.List[0].Slash,
																commentIndex: commentIndex,
															}
															commentIndex++
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	if len(fileObject1.objects) > 0 {
		g.fileObjects[filePath1] = fileObject1
	}
	return nil
}

func (g *genI) prepareRegisterContent(in fileObject) error {
	if in.objects == nil || len(in.objects) == 0 {
		return nil
	}
	fileSet1 := token.NewFileSet()
	node, err := parser.ParseFile(fileSet1, in.path, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return err
	}
	fields := []string{}
	for _, obj := range in.objects {
		fields = append(fields, obj.callerName)
	}
	g.addDynamicInitFunction(fileSet1, node, fields)
	// 使用 bytes.Buffer 将内容写入内存
	var buf bytes.Buffer
	if err1 := printer.Fprint(&buf, fileSet1, node); err1 != nil {
		return fmt.Errorf("error printing AST to buffer: %w", err1)
	}
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("格式化代码时出错: %w", err)
	}
	f := g.fileObjects[in.path]
	f.registerContent = string(formattedCode)
	g.fileObjects[in.path] = f
	return nil
}

func (g *genI) addDynamicInitFunction(fileSet1 *token.FileSet, f *ast.File, fields []string) {
	{
		exists := false
		for _, imp := range f.Imports {
			if strings.Trim(imp.Path.Value, "\"") == reflectImport {
				exists = true
				break
			}
		}
		if !exists {
			_ = astutil.AddImport(fileSet1, f, reflectImport)
		}
	}
	var stmts []ast.Stmt
	{
		// 定义变量 mongoInit 并初始化为 "@mongoInit"
		mongoVar := &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{Name: "_"}, // 变量名
			},
			Tok: token.ASSIGN, // 使用 := 进行声明和赋值
			Rhs: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,   // 字符串类型
					Value: `"@mongoInit"`, // 字符串值
				},
			},
		}
		stmts = append(stmts, mongoVar)
	}
	// 动态生成 reflect2.Register2 调用
	for _, field := range fields {
		stmt := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("reflect2"),
					Sel: ast.NewIdent("Register2"),
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent(field),
						Sel: ast.NewIdent("source"),
					},
				},
			},
		}
		stmts = append(stmts, stmt)
	}

	// 检查是否已有 init 方法
	_, initFuncExists := g.checkInit(f)

	// 如果没有找到 init 方法，则创建新的 init 方法
	if !initFuncExists {
		initFunc := &ast.FuncDecl{
			Name: ast.NewIdent("init"), // init 方法的名称
			Type: &ast.FuncType{
				Params:  &ast.FieldList{}, // 无参数
				Results: nil,              // 无返回值
			},
			Body: &ast.BlockStmt{
				List: stmts, // 动态生成的语句列表
			},
		}
		f.Decls = append(f.Decls, initFunc)
	}
}

func (g *genI) checkInit(f *ast.File) (int, bool) {
	for j, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == "init" {
			// 遍历 init 函数的 Body，检查是否有 mongoInit 变量
			for _, stmt := range funcDecl.Body.List {
				if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
					for i, _ := range assignStmt.Lhs {
						// 检查右侧的值
						if len(assignStmt.Rhs) > i {
							if basicLit, ok := assignStmt.Rhs[i].(*ast.BasicLit); ok {
								// 检查值是否为 "@mongoInit"
								if basicLit.Kind == token.STRING && basicLit.Value == `"@mongoInit"` {
									return j, true
								}
							}
						}
					}
				}
			}
		}
	}
	return 0, false
}

func (g *genI) prepareContent(in fileObject) error {
	if in.objects == nil || len(in.objects) == 0 {
		return nil
	}
	fileSet1 := token.NewFileSet()
	node, err := parser.ParseFile(fileSet1, in.path, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return err
	}
	// 遍历 AST 并打印每个节点的位置信息
	g.checkPosition(node, &in)
	g.addFields(node, in.variables)
	// 使用 bytes.Buffer 将内容写入内存
	var buf bytes.Buffer
	if err1 := printer.Fprint(&buf, fileSet1, node); err1 != nil {
		return fmt.Errorf("error printing AST to buffer: %w", err1)
	}
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("格式化代码时出错: %w", err)
	}
	{
		formattedCode, err = g.prepareComment(formattedCode, in)
		if err != nil {
			return err
		}
	}
	{
		formattedCode, err = g.deleteInit(formattedCode)
		if err != nil {
			return err
		}
	}
	f := g.fileObjects[in.path]
	f.content = string(formattedCode)
	g.fileObjects[in.path] = f
	return nil
}

func (g *genI) checkPosition(node ast.Node, in *fileObject) {
	// 遍历 AST 并打印每个节点的位置信息
	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				isTag := false
				if genDecl.Doc != nil {
					isTag = strings.Contains(genDecl.Doc.List[0].Text, commentLabel)
				}
				isSource := false
				valueSpec, ok1 := spec.(*ast.ValueSpec)
				if !ok1 {
					continue
				}
				variable1, ok2 := in.variables[valueSpec.Names[0].Name]
				if !ok2 {
					continue
				}
				for _, val := range valueSpec.Values {
					compositeLit1, ok2 := val.(*ast.CompositeLit)
					if ok2 {
						type1 := compositeLit1.Type
						structType1, ok3 := type1.(*ast.StructType)
						if ok3 {
							for _, field := range structType1.Fields.List {
								for _, name1 := range field.Names {
									if name1.Name == sourceKey {
										isSource = true
									}
								}
							}
						}
					}
				}
				if isSource || isTag {
					variable1.shift = genDecl.Pos() - variable1.codePos
					variable1.codePos = genDecl.Pos()
					if genDecl.Doc != nil {
						variable1.commentPos = genDecl.Doc.List[0].Slash
						variable1.comment = genDecl.Doc.List[0].Text
					}
					in.variables[valueSpec.Names[0].Name] = variable1
				}
			}
		}
		return true
	})
}

func (g *genI) addFields(f *ast.File, variables map[string]variable) {
	// 遍历文件声明
	for i, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
		if genDecl.Doc == nil {
			continue
		}
		contains := strings.Contains(genDecl.Doc.List[0].Text, commentLabel)
		if !contains {
			continue
		}
		// 遍历变量声明
		for j, spec := range genDecl.Specs {
			valueSpec, ok1 := spec.(*ast.ValueSpec)
			if !ok1 {
				continue
			}
			variable1, ok2 := variables[valueSpec.Names[0].Name]
			if !ok2 {
				continue
			}
			fields := tagx.Define(variable1.source.types).Export()
			fields2 := map[string]string{}
			err := copier.Copy(&fields2, &fields)
			if err != nil {
				logrus.Fatal(err)
				return
			}
			lit, ok2 := valueSpec.Values[0].(*ast.CompositeLit)
			if !ok2 {
				continue
			}
			structType1, ok3 := lit.Type.(*ast.StructType)
			if !ok3 {
				continue
			}
			{
				for k, field1 := range structType1.Fields.List {
					_, ok4 := fields[field1.Names[0].Name]
					if !ok4 {
						continue
					}
					field1.Type = ast.NewIdent("string")
					structType1.Fields.List[k] = field1
					delete(fields, field1.Names[0].Name)
				}
				for k, _ := range fields {
					newField := &ast.Field{
						Names: []*ast.Ident{ast.NewIdent(k)},
						Type:  ast.NewIdent("string"),
					}
					structType1.Fields.List = append(structType1.Fields.List, newField)
				}
			}
			{
				for _, elt1 := range lit.Elts {
					elt2, ok := elt1.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					key1, ok := elt2.Key.(*ast.Ident)
					if !ok {
						continue
					}
					_, ok = fields2[key1.Name]
					if ok {
						delete(fields2, key1.Name)
					}
				}
				for k, v := range fields2 {
					keyValue := &ast.KeyValueExpr{
						Key:   ast.NewIdent(k),
						Value: ast.NewIdent(fmt.Sprintf(`"%s"`, v)),
					}
					if lit.Elts == nil {
						lit.Elts = []ast.Expr{}
					}
					lit.Elts = append(lit.Elts, keyValue)
				}
			}
			{
				{
					comment := &ast.Comment{
						Text: variable1.comment,
					}
					valueSpec.Doc = &ast.CommentGroup{
						List: []*ast.Comment{comment},
					}
				}
				lit.Type = structType1
				valueSpec.Values[0] = lit
				genDecl.Specs[j] = valueSpec
				f.Decls[i] = genDecl
			}
		}
	}
}

func (g *genI) prepareComment(b []byte, in fileObject) ([]byte, error) {
	fileSet1 := token.NewFileSet()
	node, err := parser.ParseFile(fileSet1, "", b, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	g.checkPosition(node, &in)
	g.addComment(node, in.variables)
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

func (g *genI) deleteInit(b []byte) ([]byte, error) {
	fileSet1 := token.NewFileSet()
	node, err := parser.ParseFile(fileSet1, "", b, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	funcIndex, isExisted := g.checkInit(node)
	if isExisted {
		node.Decls = append(node.Decls[:funcIndex], node.Decls[funcIndex+1:]...)
		astutil.DeleteImport(fileSet1, node, reflectImport)
	}
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

func (g *genI) addComment(f *ast.File, variables map[string]variable) {
	for k, comment := range f.Comments {
		for _, v := range variables {
			if v.commentIndex != k {
				continue
			}
			comment.List[0].Slash = v.commentPos + v.shift + func() token.Pos {
				if k == 0 {
					return token.Pos(0)
				}
				return token.Pos(len(v.comment))
			}()
			f.Comments[k] = comment
		}
	}
}

func (g *genI) register() error {
	{
		err := g.writeFile(g.callerFile, g.callerContent)
		if err != nil {
			return err
		}
	}
	for _, f := range g.fileObjects {
		if f.objects == nil || len(f.objects) == 0 || f.registerContent == "" {
			continue
		}
		err := g.writeFile(f.path, f.registerContent)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *genI) Generate() error {
	for _, f := range g.fileObjects {
		if f.objects == nil || len(f.objects) == 0 || f.content == "" {
			continue
		}
		err := g.writeFile(f.path, f.content)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *genI) writeFile(filename1, content string) error {
	file, err := os.Create(filename1)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (g *genI) getModuleName(v2 tImport) error {
	dir := v2.dir
	for {
		modFile := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modFile); err == nil {
			data, err := os.ReadFile(modFile)
			if err != nil {
				return err
			}
			modFileParsed, err := modfile.Parse("go.mod", data, nil)
			if err != nil {
				return err
			}
			v2.module = module{
				name: modFileParsed.Module.Mod.Path,
				dir:  dir,
			}
			v2.tImport = filepath.Join(v2.module.name, strings.ReplaceAll(v2.dir, v2.module.dir, ""))
			g.xImports[v2.dir] = v2
			return nil
		}
		if dir == "/" {
			break
		}
		dir = filepath.Dir(dir)
	}
	return errors.New("go.mod not found")
}

func (g *genI) getCallerFile() error {

	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return fmt.Errorf("unable to retrieve caller information")
	}
	g.callerFile = file
	return nil
}

func (g *genI) appendImports() error {
	importPaths := []string{}
	for _, tImport1 := range g.xImports {
		importPaths = append(importPaths, tImport1.tImport)
	}
	filePath := g.callerFile

	fileSet1 := token.NewFileSet()
	node, err := parser.ParseFile(fileSet1, filePath, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return err
	}
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
	// 使用 bytes.Buffer 将内容写入内存
	var buf bytes.Buffer
	if err1 := printer.Fprint(&buf, fileSet1, node); err1 != nil {
		return fmt.Errorf("error printing AST to buffer: %w", err1)
	}
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("格式化代码时出错: %w", err)
	}
	g.callerContent = string(formattedCode)
	return nil
}
