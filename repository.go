package importshttp

// Repository represents a supported location for downloading a package.
type Repository interface {
	RepositoryVCS() VCS
	RepositoryRoot() string
}

// SourceRepository represents a repository where its source code can be found.
type SourceRepository interface {
	SourceURL() string
	SourceDirTemplateURL() string
	SourceFileTemplateURL() string
}
