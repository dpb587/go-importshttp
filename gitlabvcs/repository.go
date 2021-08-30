package gitlabvcs

import (
	"fmt"
	"net/url"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/internal/stringsutil"
)

const RepositoryService = "gitlab"

type Repository struct {
	Server     string
	Namespace  string
	Repository string
	Ref        string
}

var _ importshttp.Repository = Repository{}
var _ importshttp.SourceRepository = Repository{}

func (r Repository) RepositoryVCS() importshttp.VCS {
	return importshttp.GitVCS
}

func (r Repository) RepositoryRoot() string {
	return fmt.Sprintf(
		"%s/%s/%s",
		r.resolvedServer(),
		r.Namespace,
		r.Repository,
	)
}

func (r Repository) SourceURL() string {
	return r.RepositoryRoot()
}

func (r Repository) SourceDirTemplateURL() string {
	return fmt.Sprintf("%s/-/tree/%s{/dir}", r.RepositoryRoot(), url.PathEscape(r.resolvedRef()))
}

func (r Repository) SourceFileTemplateURL() string {
	return fmt.Sprintf("%s/-/blob/%s{/dir}/{file}#L{line}", r.RepositoryRoot(), url.PathEscape(r.resolvedRef()))
}

func (r Repository) resolvedServer() string {
	return stringsutil.Coalesce(r.Server, DefaultServer)
}

func (r Repository) resolvedRef() string {
	return stringsutil.Coalesce(r.Ref, DefaultRef)
}
