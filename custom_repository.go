package importshttp

import (
	"fmt"
	"regexp"
)

// CustomRepository supports any VCS and repository root.
//
// Consider using service-specific Repository types, when applicable, since they include built-in SourceRepository
// behaviors.
type CustomRepository struct {
	VCS  VCS
	Root string
}

var _ Repository = CustomRepository{}

func (r CustomRepository) RepositoryVCS() VCS {
	return r.VCS
}

func (r CustomRepository) RepositoryRoot() string {
	return r.Root
}

// CustomSourceRepository supports configuring source URLs for a given Repository.
type CustomSourceRepository struct {
	Repository
	URL             string
	DirTemplateURL  string
	FileTemplateURL string
}

var _ Repository = CustomSourceRepository{}
var _ SourceRepository = CustomSourceRepository{}

func (r CustomSourceRepository) RepositoryVCS() VCS {
	return r.Repository.RepositoryVCS()
}

func (r CustomSourceRepository) RepositoryRoot() string {
	return r.Repository.RepositoryRoot()
}

func (r CustomSourceRepository) SourceURL() string {
	return r.URL
}
func (r CustomSourceRepository) SourceDirTemplateURL() string {
	return r.DirTemplateURL
}
func (r CustomSourceRepository) SourceFileTemplateURL() string {
	return r.FileTemplateURL
}

// CustomRepositoryFactory supports creating a CustomRepository from a RepositoryConfig.
//
// TODO this is not fully following specs yet; weird mix of go vs url pkg parsing; https://github.com/golang/go/blob/2ebe77a2fda1ee9ff6fd9a3e08933ad1ebaea039/src/cmd/go/internal/vcs/vcs_test.go
type CustomRepositoryFactory struct{}

var _ RepositoryFactory = CustomRepositoryFactory{}

// https://github.com/golang/go/blob/4c8ffb3baaabce1aa2139ce7739fec333ab80728/src/cmd/go/internal/vcs/vcs.go#L1232
var reCustomRepositoryVCS = regexp.MustCompile((`(?P<root>(?P<repo>([a-z0-9.\-]+\.)+[a-z0-9.\-]+(:[0-9]+)?(/~?[A-Za-z0-9_.\-]+)+?)\.(?P<vcs>bzr|fossil|git|hg|svn))(/~?[A-Za-z0-9_.\-]+)*$`))

func (rf CustomRepositoryFactory) NewRepository(config RepositoryConfig) (Repository, error) {
	vcs, vcsknown := config.VCS()

	url, urlknown := config.URL()
	if urlknown {
		if !vcsknown {
			match := reCustomRepositoryVCS.FindStringSubmatch(url.String())
			if len(match) == 0 {
				return nil, ErrRepositoryConfigNotSupported
			}

			return CustomRepository{
				VCS:  VCS(match[reCustomRepositoryVCS.SubexpIndex("vcs")]),
				Root: match[reCustomRepositoryVCS.SubexpIndex("repo")],
			}, nil
		}

		return CustomRepository{
			VCS:  vcs,
			Root: fmt.Sprintf("%s%s", url.Host, url.Path),
		}, nil
	}

	// TODO support properties w/ optional source repository fields

	return nil, ErrRepositoryConfigNotSupported
}
