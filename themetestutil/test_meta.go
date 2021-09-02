package themetestutil

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"testing"

	"go.dpb.io/importshttp"
)

// TestHandler verifies a theme contains the correct meta entries for go-import.
func TestGoImports(t *testing.T, theme importshttp.Theme) {
	var out bytes.Buffer

	err := theme.PackageTemplate.Execute(
		&out,
		map[string]interface{}{
			"Site":  importshttp.Site{},
			"Theme": importshttp.Theme{},
			"Package": importshttp.Package{
				Import:           "go.example.com/test-group-1/test-module-1",
				ImportSubpackage: "test-subpackage-1",
				Repository: importshttp.CustomRepository{
					VCS:  importshttp.FossilVCS,
					Root: "vcs.example.com/test-owner-1/test-repository-1",
				},
			},
			"SubpackageList": nil,
		})
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	metas, err := parseMetaGoImports(bytes.NewReader(out.Bytes()))
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	if _e, _a := 1, len(metas); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}

	meta0 := metas[0]
	if _e, _a := "go.example.com/test-group-1/test-module-1", meta0.Prefix; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := importshttp.FossilVCS, meta0.VCS; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "vcs.example.com/test-owner-1/test-repository-1", meta0.RepoRoot; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

// based on https://github.com/golang/go/blob/2ebe77a2fda1ee9ff6fd9a3e08933ad1ebaea039/src/cmd/go/internal/vcs/discovery.go#L32-L64
func parseMetaGoImports(r io.Reader) ([]importshttp.GoImportSpec, error) {
	attrValue := func(attrs []xml.Attr, name string) string {
		for _, a := range attrs {
			if strings.EqualFold(a.Name.Local, name) {
				return a.Value
			}
		}
		return ""
	}

	d := xml.NewDecoder(r)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch strings.ToLower(charset) {
		case "utf-8", "ascii":
			return input, nil
		default:
			return nil, fmt.Errorf("can't decode XML document using charset %q", charset)
		}
	}

	d.Strict = false
	var imports []importshttp.GoImportSpec
	for {
		t, err := d.RawToken()
		if err != nil {
			if err != io.EOF && len(imports) == 0 {
				return nil, err
			}
			break
		}
		if e, ok := t.(xml.StartElement); ok && strings.EqualFold(e.Name.Local, "body") {
			break
		}
		if e, ok := t.(xml.EndElement); ok && strings.EqualFold(e.Name.Local, "head") {
			break
		}
		e, ok := t.(xml.StartElement)
		if !ok || !strings.EqualFold(e.Name.Local, "meta") {
			continue
		}
		if attrValue(e.Attr, "name") != "go-import" {
			continue
		}
		if f := strings.Fields(attrValue(e.Attr, "content")); len(f) == 3 {
			imports = append(imports, importshttp.GoImportSpec{
				Prefix:   f[0],
				VCS:      importshttp.VCS(f[1]),
				RepoRoot: f[2],
			})
		}
	}

	return imports, nil
}
