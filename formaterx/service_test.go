package formaterx

import (
	"testing"
)

func TestDefine(t *testing.T) {
	Define(WithDir("../")).Console()
}
