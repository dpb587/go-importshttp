package gitlabvcs

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"go.dpb.io/importshttp"
)

// RepositoryFactory will create a GitLab-specific Repository from a given RepositoryConfig.
//
// URLs will be autodetected if they match the configured Server and look like a GitLab-style path. Specifying "gitlab"
// as the VCS will force this factory to be used (and cause a failure if the match fails).
//
// For URL-based configuration, specify the full repository URL to the branch (otherwise DefaultRef will be used for
// the tree). URLs containing more than two segments before "/-/" will interpret extra segments as subgroups. Example:
//
//     https://gitlab.com/dpb587/go-importshttp/-/tree/main
//     https://gitlab.com/my-awesome-group/my-subgroup/my-project/-/tree/main
//
// For property-based configuration, the lowercase-form of Repository fields are required. Example:
//
//     { "server": "https://gitlab.com",
//       "namespace": "my-awesome-group/my-subgroup",
//       "project": "my-project",
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

var rePathTree = regexp.MustCompile(`/((?U).+)/([^/]+)(/-/([^/]+)/(.+))?$`)

func (rf RepositoryFactory) newFromURL(parsed *url.URL) (importshttp.Repository, error) {
	matches := rePathTree.FindStringSubmatch(parsed.Path)
	if len(matches) == 0 || (len(matches[4]) > 0 && matches[4] != "tree") {
		return nil, fmt.Errorf("expected gitlab-style path of `/{owner}(/{subgroup})...?/{repository}(/-/tree/{ref})?` but got %s", parsed.Path)
	}

	repo := Repository{
		Server:    strings.TrimSuffix(parsed.ResolveReference(&url.URL{Path: "/"}).String(), "/"),
		Namespace: matches[1],
		Project:   matches[2],
		Ref:       rf.DefaultRef,
	}

	if strings.HasPrefix(repo.Server, "//") {
		// TODO fix weird edge case of default non-desired schema factory claiming generic gitlab service config
		repo.Server = fmt.Sprintf("%s:%s", strings.SplitN(rf.Server, "://", 2)[0], repo.Server)
	}

	if len(matches[5]) > 0 {
		repo.Ref = matches[5]
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

	if val, ok := props["namespace"]; ok {
		repo.Namespace = val
	} else {
		return nil, fmt.Errorf("missing property: namespace")
	}

	if val, ok := props["project"]; ok {
		repo.Project = val
	} else {
		return nil, fmt.Errorf("missing property: project")
	}

	if val, ok := props["ref"]; ok {
		repo.Ref = val
	}

	return repo, nil
}
