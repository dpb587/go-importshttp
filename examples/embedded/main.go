// main shows how to configure this programmatically and attach the handler only to its package-specific endpoints.
package main

import (
	"net/http"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/githubvcs"
	"go.dpb.io/importshttp/themeless"
)

func main() {
	// pretend we have an existing server and pages
	mux := http.NewServeMux()

	// add our go-specific handlers
	registerGoPackages(mux)

	// back to regular server usage
	http.ListenAndServe("127.0.0.1:8080", mux)
}

func registerGoPackages(mux *http.ServeMux) {
	site := importshttp.Site{
		PackagePathPrefix: "example.com",
		Links: importshttp.LinkList{
			{
				Label: "Back to example.com",
				URL:   "https://example.com",
			},
		},
	}

	pkgs := importshttp.PackageList{
		{
			Import: "example.com/devops/sretools",
			Repository: githubvcs.Repository{
				Host:       "github.example.com",
				Owner:      "devops-team",
				Repository: "go-sretools",
			},
		},
		{
			Import: "example.com/dba/queryauditor",
			Repository: githubvcs.Repository{
				Host:       "github.example.com",
				Owner:      "dbtools",
				Repository: "go-queryauditor",
			},
		},
	}

	// TODO linkers

	handler := importshttp.NewHandler(site, themeless.Theme, pkgs)

	// only handle the literal package paths we need
	for _, pkg := range pkgs {
		mux.Handle(site.PackageURL(pkg.Path()), handler)
	}
}
