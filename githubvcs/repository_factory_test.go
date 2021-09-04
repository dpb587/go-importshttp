package githubvcs

import (
	"strings"
	"testing"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/internal/urlutil"
)

func Test_RepositoryFactory_ErrNotDetected_WrongVCS(t *testing.T) {
	_, err := RepositoryFactory{
		Host: "test-github-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.FossilVCS, urlutil.MustParse("https://test-fossil-1.com/test-owner-1/test-repository-1")))
	if _e, _a := importshttp.ErrRepositoryConfigNotSupported, err; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_ErrNotDetected_WrongServer(t *testing.T) {
	_, err := RepositoryFactory{
		Host: "test-github-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.UnknownVCS, urlutil.MustParse("https://test-fossil-1.com/test-owner-1")))
	if _e, _a := importshttp.ErrRepositoryConfigNotSupported, err; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_URL_ErrPathFormat(t *testing.T) {
	for subtestNameURL, subtestValueURL := range map[string]string{
		"TooShort":   "https://test-github-1.com/test-owner-1",
		"TooLong":    "https://test-github-1.com/test-owner-1/test-repository-1/extra",
		"NonTree":    "https://test-github-1.com/test-owner-1/test-repository-1/blob/main/README.md",
		"TreeExcess": "https://test-github-1.com/test-owner-1/test-repository-1/tree/main/dir",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			_, err := RepositoryFactory{
				Host: "test-github-1.com",
			}.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.UnknownVCS, urlutil.MustParse(subtestValueURL)))
			if err == nil {
				t.Fatal("expected error but got: nil")
			} else if _e, _a := "expected github-style path", err.Error(); !strings.Contains(_a, _e) {
				t.Fatalf("expected string to contain `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_URL_ServiceOwnerRepository(t *testing.T) {
	for subtestNameURL, subtestValueURL := range map[string]string{
		"FullURL":    "https://other-github-1.com/test-owner-1/test-repository-1",
		"PackageURL": "other-github-1.com/test-owner-1/test-repository-1",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			repo, err := RepositoryFactory{
				Host: "test-github-1.com",
			}.NewRepository(importshttp.NewRepositoryConfigURL(RepositoryService, urlutil.MustParse(subtestValueURL)))
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			repoT, ok := repo.(Repository)
			if !ok {
				t.Fatalf("assertion failed on value type: %T", repo)
			}

			if _e, _a := false, repoT.Insecure; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "other-github-1.com", repoT.Host; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-owner-1", repoT.Owner; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-repository-1", repoT.Repository; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_URL_OwnerRepository(t *testing.T) {
	for subtestNameURL, subtestValueURL := range map[string]string{
		"FullURL":    "https://test-github-1.com/test-owner-1/test-repository-1",
		"PackageURL": "test-github-1.com/test-owner-1/test-repository-1",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			repo, err := RepositoryFactory{
				Host: "test-github-1.com",
			}.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.UnknownVCS, urlutil.MustParse(subtestValueURL)))
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			repoT, ok := repo.(Repository)
			if !ok {
				t.Fatalf("assertion failed on value type: %T", repo)
			}

			if _e, _a := false, repoT.Insecure; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-github-1.com", repoT.Host; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-owner-1", repoT.Owner; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-repository-1", repoT.Repository; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_URL_OwnerRepositoryRef(t *testing.T) {
	for subtestNameURL, subtestValueURL := range map[string]string{
		"FullURL": "https://test-github-1.com/test-owner-1/test-repository-1/tree/test-customref-1",
	} {
		t.Run(subtestNameURL, func(t *testing.T) {
			repo, err := RepositoryFactory{
				Host: "test-github-1.com",
			}.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.UnknownVCS, urlutil.MustParse(subtestValueURL)))
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			repoT, ok := repo.(RepositoryRef)
			if !ok {
				t.Fatalf("assertion failed on value type: %T", repo)
			}

			if _e, _a := false, repoT.Repository.Insecure; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-github-1.com", repoT.Repository.Host; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-owner-1", repoT.Repository.Owner; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-repository-1", repoT.Repository.Repository; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			} else if _e, _a := "test-customref-1", repoT.Ref; _e != _a {
				t.Fatalf("expected `%v` but got: %v", _e, _a)
			}
		})
	}
}

func Test_RepositoryFactory_Properties_Default(t *testing.T) {
	repo, err := RepositoryFactory{
		Host: "test-github-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		RepositoryService,
		map[string]string{
			"owner":      "test-owner-1",
			"repository": "test-repository-1",
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
	} else if _e, _a := "test-github-1.com", repoT.Host; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-owner-1", repoT.Owner; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-repository-1", repoT.Repository; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func exampleConfigPropertiesValid() map[string]string {
	return map[string]string{
		"host":       "test-github-1.com",
		"owner":      "test-owner-1",
		"repository": "test-repository-1",
		"ref":        "test-customref-1",
	}
}

func Test_RepositoryFactory_Properties_NonDefault(t *testing.T) {
	repo, err := RepositoryFactory{
		Host: "test-github-1.com",
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
	} else if _e, _a := "test-github-1.com", repoT.Repository.Host; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-owner-1", repoT.Repository.Owner; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-repository-1", repoT.Repository.Repository; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "test-customref-1", repoT.Ref; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_Properties_MissingOwner(t *testing.T) {
	config := exampleConfigPropertiesValid()
	delete(config, "owner")

	_, err := RepositoryFactory{
		Host: "test-github-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		RepositoryService,
		config,
	))
	if err == nil {
		t.Fatal("expected error but got: nil")
	} else if _e, _a := "missing property: owner", err.Error(); !strings.Contains(_a, _e) {
		t.Fatalf("expected string to contain `%v` but got: %v", _e, _a)
	}
}

func Test_RepositoryFactory_Properties_MissingRepository(t *testing.T) {
	config := exampleConfigPropertiesValid()
	delete(config, "repository")

	_, err := RepositoryFactory{
		Host: "test-github-1.com",
	}.NewRepository(importshttp.NewRepositoryConfigProperties(
		RepositoryService,
		config,
	))
	if err == nil {
		t.Fatal("expected error but got: nil")
	} else if _e, _a := "missing property: repository", err.Error(); !strings.Contains(_a, _e) {
		t.Fatalf("expected string to contain `%v` but got: %v", _e, _a)
	}
}
