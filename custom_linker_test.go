package importshttp

import "testing"

func customLinkerPkgGoDev() CustomLinker {
	return CustomLinker{
		Ordering:    25,
		Label:       "Reference",
		PkgTemplate: "https://pkg.go.dev{/pkg}",
		DirTemplate: "https://pkg.go.dev{/pkg}{/dir}",
	}
}

func Test_CustomLinker_Package(t *testing.T) {
	link := customLinkerPkgGoDev().Link(Package{
		Import: "go.example.com/pkg",
	})
	if _e, _a := 25, link.Ordering; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "Reference", link.Label; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "https://pkg.go.dev/go.example.com/pkg", link.URL; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_CustomLinker_PackageDir(t *testing.T) {
	link := customLinkerPkgGoDev().Link(Package{
		Import:           "go.example.com/pkg",
		ImportSubpackage: "subpkg1/subpkg2",
	})
	if _e, _a := 25, link.Ordering; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "Reference", link.Label; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "https://pkg.go.dev/go.example.com/pkg/subpkg1/subpkg2", link.URL; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}
