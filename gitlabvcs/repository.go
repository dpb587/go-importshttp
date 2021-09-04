package gitlabvcs

import (
	"fmt"
	"net/url"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/internal/stringsutil"
)

const RepositoryService = "gitlab"

// Repository is a GitLab-specific repository.
type Repository struct {
	Insecure  bool
	Host      string
	Namespace string
	Project   string
}

var _ importshttp.Repository = Repository{}

func (r Repository) RepositoryVCS() importshttp.VCS {
	return importshttp.GitVCS
}

func (r Repository) RepositoryRoot() string {
	schema := "https"
	if r.Insecure {
		schema = "http"
	}

	return fmt.Sprintf(
		"%s://%s/%s/%s",
		schema,
		r.resolvedHost(),
		r.Namespace,
		r.Project,
	)
}

func (r Repository) resolvedHost() string {
	return stringsutil.Coalesce(r.Host, DefaultHost)
}

// RepositoryRef is a GitLab-specific repository which has a known branch.
type RepositoryRef struct {
	Repository

	// Ref is the branch name where files can be found. This should be specified, but will default to DefaultRef.
	Ref string
}

var _ importshttp.Repository = RepositoryRef{}
var _ importshttp.SourceRepository = RepositoryRef{}

func (rr RepositoryRef) SourceURL() string {
	return rr.RepositoryRoot()
}

func (rr RepositoryRef) SourceDirTemplateURL() string {
	return fmt.Sprintf("%s/-/tree/%s{/dir}", rr.RepositoryRoot(), url.PathEscape(rr.resolvedRef()))
}

func (rr RepositoryRef) SourceFileTemplateURL() string {
	return fmt.Sprintf("%s/-/blob/%s{/dir}/{file}#L{line}", rr.RepositoryRoot(), url.PathEscape(rr.resolvedRef()))
}

func (rr RepositoryRef) resolvedRef() string {
	return stringsutil.Coalesce(rr.Ref, DefaultRef)
}
