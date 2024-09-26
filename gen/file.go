package gen

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"github.com/svc0a/mgo/tagx"
	"github.com/svc0a/reflect2"
	"go/ast"
	"go/token"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"
	"os"
	"path/filepath"
	"strings"
)

type FileX interface {
	Export() fileObject
	Generate() error
	FileImport() *tImport
}

func DefineFile(file string, define tagx.Service) FileX {
	fx := &fileX{
		fileObject: fileObject{
			path:      file,
			objects:   map[string]object{},
			variables: map[string]variable{},
		},
		define:     define,
		comments:   map[token.Pos]int{},
		fileImport: nil,
	}
	{
		err := fx.scan()
		if err != nil {
			logrus.Fatal(err)
			return nil
		}
	}
	if fx.fileImport != nil {
		{
			err := fx.getModuleName()
			if err != nil {
				logrus.Fatal(err)
				return nil
			}
		}
		{
			err := fx.prepareRegisterContent()
			if err != nil {
				logrus.Fatal(err)
				return nil
			}
		}
		return fx
	}
	{
		err := fx.prepareContent()
		if err != nil {
			logrus.Fatal(err)
			return nil
		}
	}
	{
		err := fx.prepareComment()
		if err != nil {
			logrus.Fatal(err)
			return nil
		}
	}
	{
		err := fx.deleteInit()
		if err != nil {
			logrus.Fatal(err)
			return nil
		}
	}
	return fx
}

type fileX struct {
	fileObject
	node       *ast.File
	fileImport *tImport
	comments   map[token.Pos]int
	define     tagx.Service
}

func (fx *fileX) Export() fileObject {
	return fx.fileObject
}

func (fx *fileX) Generate() error {
	return writeFile(fx.path, fx.content)
}

func (fx *fileX) FileImport() *tImport {
	return fx.fileImport
}

func (fx *fileX) scan() error {
	{
		dir := filepath.Dir(fx.path)
		fx.dir = dir
	}
	b, err := parseFile(fx.path, func(tf *token.FileSet, node *ast.File) {
		fx.node = node
		{
			pkgName := fx.node.Name.Name
			fx.pkg = pkgName
		}
		for i, comment := range fx.node.Comments {
			fx.comments[comment.List[0].Slash] = i
		}
		ast.Inspect(fx.node, func(n ast.Node) bool {
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
																fullName:   fmt.Sprintf("%s.%s", fx.pkg, ident1.Name),
																pkg:        fx.pkg,
																path:       fx.path,
																name:       ident1.Name,
																samePkg:    true,
																callerName: valueSpec.Names[0].Name,
															}
															t, err2 := reflect2.TypeByName(obj.fullName)
															if err2 == nil {
																obj.types = t.Type1()
															} else {
																fx.fileImport = &tImport{
																	dir:      fx.dir,
																	filePath: fx.path,
																}
															}
															fx.objects[obj.fullName] = obj
															{
																fx.variables[valueSpec.Names[0].Name] = variable{
																	source:       obj,
																	name:         valueSpec.Names[0].Name,
																	comment:      genDecl.Doc.List[0].Text,
																	commentIndex: fx.comments[genDecl.Doc.List[0].Slash],
																}
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
																path:       fx.path,
																name:       selectorExpr1.Sel.Name,
																samePkg:    false,
																callerName: valueSpec.Names[0].Name,
															}
															t, err2 := reflect2.TypeByName(obj.fullName)
															if err2 == nil {
																obj.types = t.Type1()
															} else {
																fx.fileImport = &tImport{
																	dir:      fx.dir,
																	filePath: fx.path,
																}
															}
															fx.objects[obj.fullName] = obj
															{
																fx.variables[valueSpec.Names[0].Name] = variable{
																	source:       obj,
																	name:         valueSpec.Names[0].Name,
																	comment:      genDecl.Doc.List[0].Text,
																	commentIndex: fx.comments[genDecl.Doc.List[0].Slash],
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

func (fx *fileX) getModuleName() error {
	v2 := fx.fileImport
	dir := v2.dir
	for count := 0; count < 10; count++ {
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
			fx.fileImport = v2
			return nil
		}
		if dir == "/" {
			break
		}
		dir = filepath.Dir(dir)
	}
	return errors.New("go.mod not found")
}

func (fx *fileX) prepareRegisterContent() error {
	if fx.objects == nil || len(fx.objects) == 0 {
		return nil
	}
	fields := []string{}
	for _, obj := range fx.objects {
		fields = append(fields, obj.callerName)
	}
	formattedCode, err := parse(fx.content, func(tf *token.FileSet, node *ast.File) {
		fx.addDynamicInitFunction(tf, node, fields)
	})
	if err != nil {
		return err
	}
	fx.registerContent = formattedCode
	return nil
}

func (fx *fileX) addDynamicInitFunction(fileSet1 *token.FileSet, f *ast.File, fields []string) {
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
	_, initFuncExists := fx.checkInit(f)
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

func (fx *fileX) checkInit(f *ast.File) (int, bool) {
	for j, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == "init" {
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

func (fx *fileX) prepareContent() error {
	if fx.objects == nil || len(fx.objects) == 0 {
		return nil
	}
	formattedCode, err := parse(fx.content, func(tf *token.FileSet, node *ast.File) {
		variables := fx.variables
		// 遍历文件声明
		for i, decl := range node.Decls {
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
				fields := fx.define.Register(variable1.source.types).Export()
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
					node.Decls[i] = genDecl
				}
			}
		}
	})
	if err != nil {
		return err
	}
	fx.content = formattedCode
	return nil
}

func (fx *fileX) prepareComment() error {
	content, err := parse(fx.content, func(tf *token.FileSet, node *ast.File) {
		{
			ast.Inspect(node, func(n ast.Node) bool {
				genDecl, ok := n.(*ast.GenDecl)
				if ok && genDecl.Tok == token.VAR {
					for _, spec := range genDecl.Specs {
						isSource := false
						valueSpec, ok1 := spec.(*ast.ValueSpec)
						if !ok1 {
							continue
						}
						variable1, ok2 := fx.variables[valueSpec.Names[0].Name]
						if !ok2 {
							continue
						}
						for _, val := range valueSpec.Values {
							compositeLit1, ok2 := val.(*ast.CompositeLit)
							if !ok2 {
								continue
							}
							type1 := compositeLit1.Type
							structType1, ok3 := type1.(*ast.StructType)
							if !ok3 {
								continue
							}
							for _, field := range structType1.Fields.List {
								for _, name1 := range field.Names {
									if name1.Name == sourceKey {
										isSource = true
									}
								}
							}
						}
						if !isSource {
							continue
						}
						variable1.codePos = tf.Position(genDecl.Pos())
						fx.variables[valueSpec.Names[0].Name] = variable1
					}
				}
				return true
			})
		}
		{
			variables := fx.variables
			file := tf.File(node.Pos()) // 获取文件对象
			for k, comment := range node.Comments {
				contains := strings.Contains(comment.List[0].Text, commentLabel)
				if !contains {
					continue
				}
				for _, v := range variables {
					if v.commentIndex != k {
						continue
					}
					// 使用 FileSet 的 PositionFor 方法来获取某一行的起始位置
					pos := file.LineStart(v.codePos.Line - 1) // 获取第 `line` 行的起始位置 (Pos)
					comment.List[0].Slash = pos
					node.Comments[k] = comment
				}
			}
		}
	})
	if err != nil {
		return err
	}
	fx.content = content
	return nil
}

func (fx *fileX) deleteInit() error {
	content, err := parse(fx.content, func(tf *token.FileSet, node *ast.File) {
		funcIndex, isExisted := fx.checkInit(node)
		if isExisted {
			node.Decls = append(node.Decls[:funcIndex], node.Decls[funcIndex+1:]...)
			astutil.DeleteImport(tf, node, reflectImport)
		}
	})
	if err != nil {
		return err
	}
	fx.content = content
	return nil
}
