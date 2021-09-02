package gitlabvcs

import (
	"strings"
	"testing"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/internal/urlutil"
)

func Test_RepositoryFactory_ErrNotSupported(t *testing.T) {
	_, err := RepositoryFactory{
		Server:     "https://test-gitlab-1.com",
		DefaultRef: "test-ref-1",
	}.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.FossilVCS, urlutil.MustParse("https://test-fossil-1.com/test-owner-1/test-project-1")))
	if _e, _a := importshttp.ErrRepositoryConfigNotSupported, err; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_ErrNotDetected_WrongServer(t *testing.T) {
	_, err := RepositoryFactory{
		Server:     "https://test-gitlab-1.com",
		DefaultRef: "test-ref-1",
	}.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.UnknownVCS, urlutil.MustParse("https://test-fossil-1.com/test-owner-1")))
	if _e, _a := importshttp.ErrRepositoryConfigNotSupported, err; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_URL_ErrPathFormat(t *testing.T) {
	for subtestNameURL, subtestValueURL := range map[string]string{
		"TooShort": "https://test-gitlab-1.com/test-owner-1",
		"NonTree":  "https://test-gitlab-1.com/test-owner-1/test-project-1/-/blob/main/README.md",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			_, err := RepositoryFactory{
				Server:     "https://test-gitlab-1.com",
				DefaultRef: "test-ref-1",
			}.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.UnknownVCS, urlutil.MustParse(subtestValueURL)))
			if err == nil {
				t.Fatal("expected error but got: nil")
			} else if _e, _a := "expected gitlab-style path", err.Error(); !strings.Contains(_a, _e) {
				t.Fatalf("expected string to contain `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_URL_NamespaceRepository(t *testing.T) {
	for subtestNameURL, subtestValueURL := range map[string]string{
		"FullURL": "https://test-gitlab-1.com/test-owner-1/test-project-1",
		"LazyURL": "//test-gitlab-1.com/test-owner-1/test-project-1",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			repo, err := RepositoryFactory{
				Server:     "https://test-gitlab-1.com",
				DefaultRef: "test-ref-1",
			}.NewRepository(importshttp.NewRepositoryConfigURL("", urlutil.MustParse(subtestValueURL)))
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			repoT, ok := repo.(Repository)
			if !ok {
				t.Fatalf("assertion failed on value type: %T", repo)
			}

			if _e, _a := "https://test-gitlab-1.com", repoT.Server; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-owner-1", repoT.Namespace; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-project-1", repoT.Project; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-ref-1", repoT.Ref; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_URL_NamespaceRepositoryRef(t *testing.T) {
	for subtestNameURL, subtestValueURL := range map[string]string{
		"FullURL": "https://test-gitlab-1.com/test-owner-1/test-project-1/-/tree/test-customref-1",
		"LazyURL": "//test-gitlab-1.com/test-owner-1/test-project-1/-/tree/test-customref-1",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			repo, err := RepositoryFactory{
				Server:     "https://test-gitlab-1.com",
				DefaultRef: "test-ref-1",
			}.NewRepository(importshttp.NewRepositoryConfigURL("", urlutil.MustParse(subtestValueURL)))
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			repoT, ok := repo.(Repository)
			if !ok {
				t.Fatalf("assertion failed on value type: %T", repo)
			}

			if _e, _a := "https://test-gitlab-1.com", repoT.Server; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-owner-1", repoT.Namespace; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-project-1", repoT.Project; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-customref-1", repoT.Ref; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_URL_NamespaceSubgroupRepositoryRef(t *testing.T) {
	for subtestNameURL, subtestValueURL := range map[string]string{
		"FullURL": "https://test-gitlab-1.com/test-owner-1/test-subgroup-1/test-subgroup-2/test-project-1/-/tree/test-customref-1",
		"LazyURL": "//test-gitlab-1.com/test-owner-1/test-subgroup-1/test-subgroup-2/test-project-1/-/tree/test-customref-1",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			repo, err := RepositoryFactory{
				Server:     "https://test-gitlab-1.com",
				DefaultRef: "test-ref-1",
			}.NewRepository(importshttp.NewRepositoryConfigURL("", urlutil.MustParse(subtestValueURL)))
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			repoT, ok := repo.(Repository)
			if !ok {
				t.Fatalf("assertion failed on value type: %T", repo)
			}

			if _e, _a := "https://test-gitlab-1.com", repoT.Server; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-owner-1/test-subgroup-1/test-subgroup-2", repoT.Namespace; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-project-1", repoT.Project; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-customref-1", repoT.Ref; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_Properties_Default(t *testing.T) {
	repo, err := RepositoryFactory{
		Server:     "https://test-gitlab-1.com",
		DefaultRef: "test-ref-1",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		importshttp.GitVCS,
		map[string]string{
			"namespace": "test-owner-1",
			"project":   "test-project-1",
		},
	))
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	repoT, ok := repo.(Repository)
	if !ok {
		t.Fatalf("assertion failed on value type: %T", repo)
	}

	if _e, _a := "https://test-gitlab-1.com", repoT.Server; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-owner-1", repoT.Namespace; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-project-1", repoT.Project; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-ref-1", repoT.Ref; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func exampleConfigPropertiesValid() map[string]string {
	return map[string]string{
		"server":    "testproto://test-gitlab-1.com",
		"namespace": "test-owner-1/test-subgroup-1",
		"project":   "test-project-1",
		"ref":       "test-customref-1",
	}
}

func Test_RepositoryFactory_Properties_NonDefault(t *testing.T) {
	repo, err := RepositoryFactory{
		Server:     "https://test-gitlab-1.com",
		DefaultRef: "test-ref-1",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		importshttp.GitVCS,
		exampleConfigPropertiesValid(),
	))
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	repoT, ok := repo.(Repository)
	if !ok {
		t.Fatalf("assertion failed on value type: %T", repo)
	}

	if _e, _a := "testproto://test-gitlab-1.com", repoT.Server; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-owner-1/test-subgroup-1", repoT.Namespace; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-project-1", repoT.Project; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-customref-1", repoT.Ref; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_Properties_MissingNamespace(t *testing.T) {
	config := exampleConfigPropertiesValid()
	delete(config, "namespace")

	_, err := RepositoryFactory{
		Server:     "https://test-gitlab-1.com",
		DefaultRef: "test-ref-1",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		importshttp.GitVCS,
		config,
	))
	if err == nil {
		t.Fatal("expected error but got: nil")
	} else if _e, _a := "missing property: namespace", err.Error(); !strings.Contains(_a, _e) {
		t.Fatalf("expected string to contain `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_Properties_MissingRepository(t *testing.T) {
	config := exampleConfigPropertiesValid()
	delete(config, "project")

	_, err := RepositoryFactory{
		Server:     "https://test-gitlab-1.com",
		DefaultRef: "test-ref-1",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		importshttp.GitVCS,
		config,
	))
	if err == nil {
		t.Fatal("expected error but got: nil")
	} else if _e, _a := "missing property: project", err.Error(); !strings.Contains(_a, _e) {
		t.Fatalf("expected string to contain `%v` but got: %v", _e, _a)
	}
}
