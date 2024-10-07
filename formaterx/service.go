package formaterx

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/svc0a/mgo/tagx"
	"go/token"
	"os"
	"path/filepath"
)

type Formater interface {
	Console()
	Execute()
}

type impl struct {
	dir    string
	files  map[string]map[string]Result // map[filename]map[pkg.struct.property]Result
	rule   string
	client tagx.Client
}

type Result struct {
	file     string
	line     token.Position
	property string // pkg.struct.property
	value    string
	expected string
	err      error
	typeName string
}

type Option func(*impl)

func WithDir(dir string) Option {
	return func(f *impl) {
		f.dir = dir
	}
}

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

func Define(options ...Option) Formater {
	f := &impl{
		dir:    "./",
		files:  make(map[string]map[string]Result),
		client: tagx.Mongo(),
	}
	for _, o := range options {
		o(f)
	}
	err := f.init()
	if err != nil {
		logrus.Fatal(err)
		return nil
	}
	return f
}

func (svc *impl) Console() {
	for _, v := range svc.files {
		for _, result := range v {
			if result.err != nil {
				logrus.WithField("file", fmt.Sprintf("%v:%v", result.file, result.line.Line)).WithField("property", result.property).WithField("value", result.value).WithField("expected", result.expected).Error(result.err)
			}
		}
	}
}

func (svc *impl) Execute() {
	return
}

func (svc *impl) init() error {
	{
		err := svc.scanDir()
		if err != nil {
			return err
		}
	}
	for file1 := range svc.files {
		results := defineFile(file1, svc.client).Export()
		svc.files[file1] = results
	}
	return nil
}

func (svc *impl) scanDir() error {
	err := filepath.Walk(svc.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			svc.files[path] = map[string]Result{}
		}
		return nil
	})
	return err
}
