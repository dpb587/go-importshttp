package importshttp

import (
	"fmt"
	"strings"
)

// CustomLinker generates a link using a package or directory template.
type CustomLinker struct {
	// Ordering which will be propagated to any generated Link.
	Ordering int

	// Label which will be propagated to any generated Link.
	Label string

	// PkgTemplate may use the `{/pkg}` placeholder.
	PkgTemplate string

	// DirTemplate may use the `{/pkg}` and `{/dir}` placeholders. If empty, PkgTemplate is used.
	DirTemplate string
}

func (l CustomLinker) Link(pkg Package) *Link {
	link := &Link{
		Ordering: l.Ordering,
		Label:    l.Label,
	}

	if len(pkg.ImportSubpackage) > 0 && len(l.DirTemplate) > 0 {
		link.URL = strings.Replace(strings.Replace(l.PkgTemplate, "{/pkg}", fmt.Sprintf("/%s", pkg.Import), 1), "{/dir}", fmt.Sprintf("/%s", pkg.ImportSubpackage), 1)
	} else {
		link.URL = strings.Replace(l.PkgTemplate, "{/pkg}", fmt.Sprintf("/%s", pkg.Import), 1)
	}

	return link
}
