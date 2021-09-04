package gitlabvcs

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
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
//     gitlab.com/my-awesome-group/my-subgroup/my-project
//
// For property-based configuration, the lowercase-form of fields are acceptable. Example:
//
//     { "host": "gitlab.com",
//       "namespace": "my-awesome-group/my-subgroup",
//       "project": "my-project",
//       "ref": "main" }
type RepositoryFactory struct {
	Insecure bool
	Host     string
}

func (rf RepositoryFactory) NewRepository(config importshttp.RepositoryConfig) (importshttp.Repository, error) {
	vcs, vcsknown := config.VCS()
	if vcsknown && vcs != RepositoryService && vcs != importshttp.GitVCS {
		return nil, importshttp.ErrRepositoryConfigNotSupported
	}

	url, urlknown := config.URL()
	if urlknown {
		normurl, urlmatch := rf.matchURL(vcs, url)
		if !urlmatch {
			return nil, importshttp.ErrRepositoryConfigNotSupported
		}

		return rf.newFromURL(normurl)
	} else if vcs != RepositoryService {
		return nil, importshttp.ErrRepositoryConfigNotSupported
	}

	props, propsknown := config.Properties()
	if !propsknown {
		return nil, fmt.Errorf("received invalid %T", config)
	}

	return rf.newFromProperties(props)
}

func (rf RepositoryFactory) matchURL(vcs importshttp.VCS, parsed *url.URL) (*url.URL, bool) {
	normurl := parsed.ResolveReference(&url.URL{})

	if len(parsed.Host) == 0 {
		pathSplit := strings.SplitN(parsed.Path, "/", 2)
		normurl.Host = pathSplit[0]
		if len(pathSplit) == 2 {
			normurl.Path = fmt.Sprintf("/%s", pathSplit[1])
		} else {
			normurl.Path = ""
		}
	}

	if len(parsed.Scheme) == 0 {
		if rf.Insecure {
			normurl.Scheme = "http:"
		} else {
			normurl.Scheme = "https:"
		}
	}

	if rf.Host == normurl.Host {
		return normurl, true
	} else if vcs == RepositoryService {
		return normurl, true
	}

	return nil, false
}

var rePathPreMatch = regexp.MustCompile(`^/([^/]+((?U)/[^/]+)*)/([^/]+)$`)

func (rf RepositoryFactory) newFromURL(parsed *url.URL) (importshttp.Repository, error) {
	pathSlashSplit := strings.SplitN(parsed.Path, "/-/", 2)

	matchPre := rePathPreMatch.FindStringSubmatch(pathSlashSplit[0])
	if len(matchPre) == 0 {
		return nil, fmt.Errorf("expected gitlab-style path of `/{owner}(/{subgroup})...?/{repository}(/-/tree/{ref})?` but got %s", parsed.Path)
	}

	repo := Repository{
		Insecure:  parsed.Scheme == "http:",
		Host:      parsed.Host,
		Namespace: matchPre[1],
		Project:   strings.TrimPrefix(matchPre[3], "/"),
	}

	var res importshttp.Repository = repo

	if len(pathSlashSplit) > 1 {
		pathPostSplit := strings.SplitN(pathSlashSplit[1], "/", 2)
		if len(pathPostSplit) != 2 || pathPostSplit[0] != "tree" || len(pathPostSplit[1]) == 0 {
			return nil, fmt.Errorf("expected gitlab-style path of `/{owner}(/{subgroup})...?/{repository}(/-/tree/{ref})?` but got %s", parsed.Path)
		}

		res = RepositoryRef{
			Repository: repo,
			Ref:        pathPostSplit[1],
		}
	}

	return res, nil
}

func (rf RepositoryFactory) newFromProperties(props map[string]string) (importshttp.Repository, error) {
	var repo = Repository{
		Host: rf.Host,
	}

	if val, ok := props["insecure"]; ok {
		valBool, err := strconv.ParseBool(val)
		if err != nil {
			return nil, errors.New("invalid property: insecure")
		}

		repo.Insecure = valBool
	}

	if val, ok := props["host"]; ok {
		repo.Host = val
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

	var res importshttp.Repository = repo

	if val, ok := props["ref"]; ok {
		res = RepositoryRef{
			Repository: repo,
			Ref:        val,
		}
	}

	return res, nil
}
