package importshttp

import (
	"errors"
	"net/url"
)

// ErrRepositoryConfigNotSupported indicates the config is not supported (e.g. expects git, but got svn; or server does
// not match).
var ErrRepositoryConfigNotSupported = errors.New("repository config not supported")

// RepositoryFactory supports converting an untyped RepositoryConfig to a Repository.
type RepositoryFactory interface {
	NewRepository(config RepositoryConfig) (Repository, error)
}

// BestEffortRepositoryFactory uses multiple repository factories in attempting to convert a RepositoryConfig, ignoring
// some expected conversion errors in favor of trying another factory.
type BestEffortRepositoryFactory []RepositoryFactory

var _ RepositoryFactory = BestEffortRepositoryFactory{}

func (rfl BestEffortRepositoryFactory) NewRepository(config RepositoryConfig) (Repository, error) {
	for _, rf := range rfl {
		repository, err := rf.NewRepository(config)
		if err == nil {
			return repository, nil
		} else if err == ErrRepositoryConfigNotSupported {
			continue
		}

		return nil, err
	}

	return nil, ErrRepositoryConfigNotSupported
}

// RepositoryConfig contains configuration for building a Repository - either by a URI or key-values. Unlike the package
// value which is passed to `go get`, we assume it can be a full URL since this prefers more details (e.g. to get hints
// about the branch).
type RepositoryConfig struct {
	vcs        VCS
	url        *url.URL
	properties map[string]string
}

func NewRepositoryConfigURL(vcs VCS, url *url.URL) RepositoryConfig {
	return RepositoryConfig{
		vcs: vcs,
		url: url,
	}
}

func NewRepositoryConfigProperties(vcs VCS, properties map[string]string) RepositoryConfig {
	return RepositoryConfig{
		vcs:        vcs,
		properties: properties,
	}
}

func (rc RepositoryConfig) VCS() (VCS, bool) {
	return rc.vcs, rc.vcs != UnknownVCS
}

func (rc RepositoryConfig) URL() (*url.URL, bool) {
	return rc.url, rc.url != nil
}

func (rc RepositoryConfig) Properties() (map[string]string, bool) {
	return rc.properties, rc.properties != nil
}
