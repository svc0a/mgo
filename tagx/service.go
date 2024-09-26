package tagx

import (
	"github.com/svc0a/reflect2"
	"reflect"
)

type Client interface {
	Register(in reflect.Type) Client
	Register2(in reflect2.Type) Client
	Export() map[string]string
}
