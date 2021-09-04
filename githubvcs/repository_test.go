package githubvcs

import (
	"testing"

	"go.dpb.io/importshttp"
)

func Test_RepositoryRef_Explicit(t *testing.T) {
	subject := RepositoryRef{
		Repository: Repository{
			Host:       "test-github-1.com",
			Owner:      "test-owner-1",
			Repository: "test-repository-1",
		},
		Ref: "test-ref-1",
	}

	if _e, _a := importshttp.GitVCS, subject.RepositoryVCS(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "https://test-github-1.com/test-owner-1/test-repository-1", subject.RepositoryRoot(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "https://test-github-1.com/test-owner-1/test-repository-1", subject.SourceURL(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "https://test-github-1.com/test-owner-1/test-repository-1/tree/test-ref-1{/dir}", subject.SourceDirTemplateURL(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "https://test-github-1.com/test-owner-1/test-repository-1/blob/test-ref-1{/dir}/{file}#L{line}", subject.SourceFileTemplateURL(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func Test_Repository_Implicit(t *testing.T) {
	subject := Repository{
		Owner:      "test-owner-1",
		Repository: "test-repository-1",
	}

	if _e, _a := importshttp.GitVCS, subject.RepositoryVCS(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "https://github.com/test-owner-1/test-repository-1", subject.RepositoryRoot(); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}
