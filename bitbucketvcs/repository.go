package bitbucketvcs

import (
	"fmt"
	"net/url"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/internal/stringsutil"
)

const RepositoryService = "bitbucket"

type Repository struct {
	Server     string
	Owner      string
	Repository string
	Ref        string
}

var _ importshttp.Repository = Repository{}
var _ importshttp.SourceRepository = Repository{}

func (r Repository) Service() string {
	return RepositoryService
}

func (r Repository) RepositoryVCS() importshttp.VCS {
	return "git"
}

func (r Repository) RepositoryRoot() string {
	return fmt.Sprintf(
		"%s/%s/%s",
		r.resolvedServer(),
		r.Owner,
		r.Repository,
	)
}

func (r Repository) SourceURL() string {
	return r.RepositoryRoot()
}

func (r Repository) SourceDirTemplateURL() string {
	return fmt.Sprintf("%s/src/%s{/dir}", r.RepositoryRoot(), url.PathEscape(r.resolvedRef()))
}

func (r Repository) SourceFileTemplateURL() string {
	return fmt.Sprintf("%s/src/%s{/dir}/{file}#{file}-{line}", r.RepositoryRoot(), url.PathEscape(r.resolvedRef()))
}

func (r Repository) resolvedServer() string {
	return stringsutil.Coalesce(r.Server, DefaultServer)
}

func (r Repository) resolvedRef() string {
	return stringsutil.Coalesce(r.Ref, DefaultRef)
}
