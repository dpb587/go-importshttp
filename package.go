package importshttp

import "fmt"

// Package represents a Go package which can be imported and includes details about its repository, deprecation status,
// and useful metadata for rendering information to users.
//
// These values require more explicit configuration since this package does not perform any of the dynamic resolution
// steps of, for example, the package argument given to "go get".
type Package struct {
	// Import should be a module path - the canonical name for a module. For older, non-module repos this is typically the
	// repository root path. There must not be a leading or trailing slash.
	Import string

	// ImportSubpackage is an optional subpackage within the import. There must not be a leading or trailing slash.
	//
	// Typically this should be empty when configuring Package - this is mostly used when subpackage routing is enabled to
	// support dynamic subpackage information.
	ImportSubpackage string

	// Repository provides repository details for the Import.
	//
	// Note that ImportSubpackage would be considered a value for the "{/dir}" attribute and Repository should not already
	// be including it.
	Repository Repository

	// Deprecated indicates the go.mod file has the deprecated directive.
	Deprecated bool

	// Unlisted indicates this package should not be shown in any lists when rendering. Direct routes to this package or
	// its subpackages continue to work.
	Unlisted bool

	// Links are a list of relevant resources for this package. The lowest-order link is considered primary and may be
	// featured more prominently.
	Links LinkList

	// Metadata is arbitrary data intended for themes.
	Metadata map[string]interface{}
}

// Path is the package path - the combination of Import and ImportSubpackage.
func (i Package) Path() string {
	result := i.Import

	if len(i.ImportSubpackage) > 0 {
		result = fmt.Sprintf("%s/%s", result, i.ImportSubpackage)
	}

	return result
}

// GoGetImport builds the
func (i Package) GoGetImport() *GoGetImport {
	return &GoGetImport{
		Prefix:   i.Import,
		VCS:      i.Repository.RepositoryVCS(),
		RepoRoot: i.Repository.RepositoryRoot(),
	}
}

func (i Package) GoGetSource() *GoGetSource {
	sr, ok := i.Repository.(SourceRepository)
	if !ok {
		return nil
	}

	return &GoGetSource{
		RepoRootPrefix: i.Import,
		RepoURL:        sr.SourceURL(),
		DirTemplate:    sr.SourceDirTemplateURL(),
		FileTemplate:   sr.SourceFileTemplateURL(),
	}
}
