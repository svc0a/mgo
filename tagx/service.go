package tagx

import (
	"github.com/svc0a/reflect2"
	"reflect"
)

type Service interface {
	Register(in reflect.Type) Service
	Register2(in reflect2.Type) Service
	Export() map[string]string
}
