package importshttp

import (
	"fmt"
	"net/http"
	"strings"
)

// IsGoGetRequest checks if the request has the expected "go-get=1" param in the query string.
func IsGoGetRequest(r *http.Request) bool {
	return r.URL.RawQuery == "go-get=1" || strings.HasPrefix(r.URL.RawQuery, "go-get=1&") || strings.HasSuffix(r.URL.RawQuery, "&go-get=1") || strings.Contains(r.URL.RawQuery, "&go-get=1&")
}

type GoImportSpec struct {
	Prefix   string
	VCS      VCS
	RepoRoot string
}

func (v GoImportSpec) MetaContent() string {
	return fmt.Sprintf("%s %s %s", v.Prefix, v.VCS, v.RepoRoot)
}

type GoSourceSpec struct {
	RepoRootPrefix string
	RepoURL        string
	DirTemplate    string
	FileTemplate   string
}

func (v GoSourceSpec) MetaContent() string {
	return strings.TrimSpace(strings.Join(
		[]string{
			v.RepoRootPrefix,
			v.RepoURL,
			v.DirTemplate,
			v.FileTemplate,
		},
		" ",
	))
}
