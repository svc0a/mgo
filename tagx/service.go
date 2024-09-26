package tagx

import "github.com/svc0a/reflect2"

type Service interface {
	Register(in reflect2.Type) Service
	Export() map[string]string
}
