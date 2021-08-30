package importshttp

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
type CustomRepositoryFactory struct{}

var _ RepositoryFactory = CustomRepositoryFactory{}

func (rf CustomRepositoryFactory) NewRepository(config RepositoryConfig) (Repository, error) {
	vcs, vcsknown := config.VCS()
	if !vcsknown {
		return nil, ErrRepositoryConfigNotSupported
	}

	url, urlknown := config.URL()
	if urlknown {
		return CustomRepository{
			VCS:  vcs,
			Root: url.String(),
		}, nil
	}

	// TODO support properties w/ optional source repository fields

	return nil, ErrRepositoryConfigNotSupported
}
