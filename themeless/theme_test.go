package themeless

import (
	"testing"

	"go.dpb.io/importshttp/themetestutil"
)

func TestGoImports(t *testing.T) {
	themetestutil.TestGoImports(t, Theme)
}
