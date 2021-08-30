package importshttp

import (
	"errors"
	"net/url"
)

var ErrRepositoryConfigNotSupported = errors.New("repository config not supported")
var ErrRepositoryConfigNotDetected = errors.New("repository config not detected")

// RepositoryFactory supports converting an untyped RepositoryConfig to a Repository.
type RepositoryFactory interface {
	NewRepository(config RepositoryConfig) (Repository, error)
}

// BestEffortRepositoryFactory attempts multiple repository factories, ignoring expected some expected errors.
type BestEffortRepositoryFactory []RepositoryFactory

var _ RepositoryFactory = BestEffortRepositoryFactory{}

func (rfl BestEffortRepositoryFactory) NewRepository(config RepositoryConfig) (Repository, error) {
	for _, rf := range rfl {
		repository, err := rf.NewRepository(config)
		if err == nil {
			return repository, nil
		} else if err == ErrRepositoryConfigNotSupported || err == ErrRepositoryConfigNotDetected {
			continue
		}

		return nil, err
	}

	return nil, ErrRepositoryConfigNotSupported
}

// RepositoryConfig contains configuration for building a Repository - either by a URI or key-values.
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
