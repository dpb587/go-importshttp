package gitlabvcs

import (
	"strings"
	"testing"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/internal/urlutil"
)

func Test_RepositoryFactory_ErrNotSupported(t *testing.T) {
	_, err := RepositoryFactory{
		Host: "test-gitlab-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.FossilVCS, urlutil.MustParse("https://test-fossil-1.com/test-owner-1/test-project-1")))
	if _e, _a := importshttp.ErrRepositoryConfigNotSupported, err; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_ErrNotDetected_WrongServer(t *testing.T) {
	_, err := RepositoryFactory{
		Host: "test-gitlab-1.com",
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
				Host: "test-gitlab-1.com",
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
		"LazyURL": "test-gitlab-1.com/test-owner-1/test-project-1",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			repo, err := RepositoryFactory{
				Host: "test-gitlab-1.com",
			}.NewRepository(importshttp.NewRepositoryConfigURL("", urlutil.MustParse(subtestValueURL)))
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			repoT, ok := repo.(Repository)
			if !ok {
				t.Fatalf("assertion failed on value type: %T", repo)
			}

			if _e, _a := false, repoT.Insecure; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-gitlab-1.com", repoT.Host; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-owner-1", repoT.Namespace; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-project-1", repoT.Project; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_URL_NamespaceRepositoryRef(t *testing.T) {
	for subtestNameURL, subtestValueURL := range map[string]string{
		"FullURL": "https://test-gitlab-1.com/test-owner-1/test-project-1/-/tree/test-customref-1",
		"LazyURL": "test-gitlab-1.com/test-owner-1/test-project-1/-/tree/test-customref-1",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			repo, err := RepositoryFactory{
				Host: "test-gitlab-1.com",
			}.NewRepository(importshttp.NewRepositoryConfigURL("", urlutil.MustParse(subtestValueURL)))
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			repoT, ok := repo.(RepositoryRef)
			if !ok {
				t.Fatalf("assertion failed on value type: %T", repo)
			}

			if _e, _a := false, repoT.Repository.Insecure; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-gitlab-1.com", repoT.Repository.Host; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-owner-1", repoT.Repository.Namespace; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-project-1", repoT.Repository.Project; _e != _a {
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
		"LazyURL": "test-gitlab-1.com/test-owner-1/test-subgroup-1/test-subgroup-2/test-project-1/-/tree/test-customref-1",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			repo, err := RepositoryFactory{
				Host: "test-gitlab-1.com",
			}.NewRepository(importshttp.NewRepositoryConfigURL("", urlutil.MustParse(subtestValueURL)))
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			repoT, ok := repo.(RepositoryRef)
			if !ok {
				t.Fatalf("assertion failed on value type: %T", repo)
			}

			if _e, _a := false, repoT.Repository.Insecure; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-gitlab-1.com", repoT.Repository.Host; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-owner-1/test-subgroup-1/test-subgroup-2", repoT.Repository.Namespace; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-project-1", repoT.Repository.Project; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-customref-1", repoT.Ref; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_Properties_Default(t *testing.T) {
	repo, err := RepositoryFactory{
		Host: "test-gitlab-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		RepositoryService,
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

	if _e, _a := false, repoT.Insecure; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-gitlab-1.com", repoT.Host; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-owner-1", repoT.Namespace; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-project-1", repoT.Project; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func exampleConfigPropertiesValid() map[string]string {
	return map[string]string{
		"server":    "test-gitlab-1.com",
		"namespace": "test-owner-1/test-subgroup-1",
		"project":   "test-project-1",
		"ref":       "test-customref-1",
	}
}

func Test_RepositoryFactory_Properties_NonDefault(t *testing.T) {
	repo, err := RepositoryFactory{
		Host: "test-gitlab-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		RepositoryService,
		exampleConfigPropertiesValid(),
	))
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	repoT, ok := repo.(RepositoryRef)
	if !ok {
		t.Fatalf("assertion failed on value type: %T", repo)
	}

	if _e, _a := false, repoT.Repository.Insecure; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-gitlab-1.com", repoT.Repository.Host; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-owner-1/test-subgroup-1", repoT.Repository.Namespace; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-project-1", repoT.Repository.Project; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-customref-1", repoT.Ref; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_Properties_MissingNamespace(t *testing.T) {
	config := exampleConfigPropertiesValid()
	delete(config, "namespace")

	_, err := RepositoryFactory{
		Host: "test-gitlab-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		RepositoryService,
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
		Host: "test-gitlab-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		RepositoryService,
		config,
	))
	if err == nil {
		t.Fatal("expected error but got: nil")
	} else if _e, _a := "missing property: project", err.Error(); !strings.Contains(_a, _e) {
		t.Fatalf("expected string to contain `%v` but got: %v", _e, _a)
	}
}
