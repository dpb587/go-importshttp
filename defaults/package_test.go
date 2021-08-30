package defaults

import (
	"strings"
	"testing"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/githubvcs"
	"go.dpb.io/importshttp/gitlabvcs"
	"go.dpb.io/importshttp/internal/urlutil"
)

func Test_RepositoryFactory_GitHub(t *testing.T) {
	repo, err := RepositoryFactory.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.GitVCS, urlutil.MustParse("//github.com/test-owner-1/test-repository-1/tree/main")))
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	} else if _e, _a := importshttp.GitVCS, repo.RepositoryVCS(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _v := "github.com", repo.RepositoryRoot(); !strings.Contains(_v, _e) {
		t.Fatalf("expected string containing `%v` but got: %v", _e, _v)
	} else if _, ok := repo.(githubvcs.Repository); !ok {
		t.Fatalf("unexpected type: %T", repo)
	}
}

func Test_RepositoryFactory_GitLab(t *testing.T) {
	repo, err := RepositoryFactory.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.GitVCS, urlutil.MustParse("//gitlab.com/test-owner-1/test-subgroup-1/test-repository-1/tree/main")))
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	} else if _e, _a := importshttp.GitVCS, repo.RepositoryVCS(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _v := "gitlab.com", repo.RepositoryRoot(); !strings.Contains(_v, _e) {
		t.Fatalf("expected string containing `%v` but got: %v", _e, _v)
	} else if _, ok := repo.(gitlabvcs.Repository); !ok {
		t.Fatalf("unexpected type: %T", repo)
	}
}

func Test_RepositoryFactory_Fallback(t *testing.T) {
	repo, err := RepositoryFactory.NewRepository(importshttp.NewRepositoryConfigURL(importshttp.FossilVCS, urlutil.MustParse("//fossil.example.com/somewhere/else/entirely")))
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	} else if _e, _a := importshttp.FossilVCS, repo.RepositoryVCS(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "fossil.example.com/somewhere/else/entirely", repo.RepositoryRoot(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _, ok := repo.(importshttp.CustomRepository); !ok {
		t.Fatalf("unexpected type: %T", repo)
	}
}
