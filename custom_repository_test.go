package importshttp

import (
	"testing"

	"go.dpb.io/importshttp/internal/urlutil"
)

func Test_CustomRepositoryFactory_URL(t *testing.T) {
	repo, err := CustomRepositoryFactory{}.NewRepository(NewRepositoryConfigURL(SubversionVCS, urlutil.MustParse("//svn.example.com/repository")))
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	repoT, ok := repo.(CustomRepository)
	if !ok {
		t.Fatalf("assertion failed on value type: %T", repo)
	}

	if _e, _a := SubversionVCS, repoT.VCS; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "svn.example.com/repository", repoT.Root; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_CustomRepositoryFactory_URL_SuffixVCS(t *testing.T) {
	repo, err := CustomRepositoryFactory{}.NewRepository(NewRepositoryConfigURL(UnknownVCS, urlutil.MustParse("//svn.example.com/repository.svn")))
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	repoT, ok := repo.(CustomRepository)
	if !ok {
		t.Fatalf("assertion failed on value type: %T", repo)
	}

	if _e, _a := SubversionVCS, repoT.VCS; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "svn.example.com/repository", repoT.Root; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}
