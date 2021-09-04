// defaults provides the most common, recommended configurations.
package defaults

import (
	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/githubvcs"
	"go.dpb.io/importshttp/gitlabvcs"
	"go.dpb.io/importshttp/themepro"
)

// Theme uses the themepro package.
var Theme = themepro.Theme

// Linkers supports "Reference" links to pkg.go.dev and "Source" links for source repositories.
var Linkers = importshttp.LinkerList{
	importshttp.CustomLinker{
		Ordering:    25,
		Label:       "Reference",
		PkgTemplate: "https://pkg.go.dev{/pkg}",
		DirTemplate: "https://pkg.go.dev{/pkg}{/dir}",
	},
	importshttp.SourceRepositoryLinker{
		Ordering: 50,
		Label:    "Source",
	},
}

// RepositoryFactory supports detecting GitHub and GitLab repositories (or explicit references for self-hosted
// installations), as well as a fallback for custom repositories with explicit VCS and reopsitory root.
var RepositoryFactory = importshttp.BestEffortRepositoryFactory{
	githubvcs.RepositoryFactory{
		Host: githubvcs.DefaultHost,
	},
	gitlabvcs.RepositoryFactory{
		Host: gitlabvcs.DefaultHost,
	},
	importshttp.CustomRepositoryFactory{},
}
