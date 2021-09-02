package githubvcs

import (
	"fmt"
	"net/url"
	"strings"

	"go.dpb.io/importshttp"
)

// RepositoryFactory will create a GitHub-specific Repository from a given RepositoryConfig.
//
// URLs will be autodetected if they match the configured Server and look like a GitHub-style path. Specifying "github"
// as the VCS will force this factory to be used (and cause a failure if the match fails).
//
// For URL-based configuration, specify the full repository URL to the branch (otherwise DefaultRef will be used for
// the tree). Example:
//
//     https://github.com/dpb587/go-importshttp/tree/main
//
// For property-based configuration, the lowercase-form of Repository fields are required. Example:
//
//     { "server": "https://github.com",
//       "owner": "dpb587",
//       "repository": "go-importshttp",
//       "ref": "main" }
type RepositoryFactory struct {
	Server     string
	DefaultRef string
}

func (rf RepositoryFactory) NewRepository(config importshttp.RepositoryConfig) (importshttp.Repository, error) {
	vcs, vcsknown := config.VCS()
	if vcsknown && vcs != RepositoryService && vcs != importshttp.GitVCS {
		return nil, importshttp.ErrRepositoryConfigNotSupported
	}

	url, urlknown := config.URL()
	if urlknown {
		urlmatch := rf.matchServer(url)
		if !urlmatch && (!vcsknown || vcs == importshttp.GitVCS) {
			return nil, importshttp.ErrRepositoryConfigNotSupported
		}

		return rf.newFromURL(url)
	}

	props, propsknown := config.Properties()
	if !propsknown {
		return nil, fmt.Errorf("received invalid %T", config)
	}

	return rf.newFromProperties(props)
}

func (rf RepositoryFactory) matchServer(parsed *url.URL) bool {
	serverSplit := strings.SplitN(rf.Server, "://", 2)

	if serverSplit[1] != parsed.Host {
		return false
	} else if len(parsed.Scheme) == 0 {
		return true
	} else if serverSplit[0] == parsed.Scheme {
		return true
	}

	return false
}

func (rf RepositoryFactory) newFromURL(parsed *url.URL) (importshttp.Repository, error) {
	pathSplit := strings.SplitN(parsed.Path, "/", 6)
	pathSplitLen := len(pathSplit)
	if pathSplitLen < 3 || pathSplitLen == 4 || pathSplitLen == 6 || (pathSplitLen == 5 && pathSplit[3] != "tree") {
		return nil, fmt.Errorf("expected github-style path of `/{owner}/{repository}(/tree/{ref})?` but got %s", parsed.Path)
	}

	repo := Repository{
		Server:     strings.TrimSuffix(parsed.ResolveReference(&url.URL{Path: "/"}).String(), "/"),
		Owner:      pathSplit[1],
		Repository: pathSplit[2],
		Ref:        rf.DefaultRef,
	}

	if strings.HasPrefix(repo.Server, "//") {
		// TODO fix weird edge case of default non-desired schema factory claiming generic github service config
		repo.Server = fmt.Sprintf("%s:%s", strings.SplitN(rf.Server, "://", 2)[0], repo.Server)
	}

	if pathSplitLen > 4 {
		repo.Ref = pathSplit[4]
	}

	return repo, nil
}

func (rf RepositoryFactory) newFromProperties(props map[string]string) (importshttp.Repository, error) {
	var repo = Repository{
		Server: rf.Server,
		Ref:    rf.DefaultRef,
	}

	if val, ok := props["server"]; ok {
		repo.Server = val
	}

	if val, ok := props["owner"]; ok {
		repo.Owner = val
	} else {
		return nil, fmt.Errorf("missing property: owner")
	}

	if val, ok := props["repository"]; ok {
		repo.Repository = val
	} else {
		return nil, fmt.Errorf("missing property: repository")
	}

	if val, ok := props["ref"]; ok {
		repo.Ref = val
	}

	return repo, nil
}
