package importshttp

import (
	"fmt"
	"strings"
)

// Site contains global settings that apply to the entire site.
type Site struct {
	// URL allows providing an absolute address (instead of using /).
	URL string

	// Title overrides the default import-host value on some pages.
	Title string

	// Generator is a value to be shown in meta tags.
	Generator string

	// ContentLanguage is the primary language value used by themes and packages.
	ContentLanguage string

	// PackagePathPrefix determines the part of the package path which should be removed when building site links.
	// Typically this is the first segment of the remote package imports - the hostname.
	PackagePathPrefix string

	// PackageLinkers are used for generating links for dynamically-discovered subpackages.
	PackageLinkers LinkerList

	// Links apply to all pages.
	Links LinkList

	// Metadata is arbitrary data which may be used by themes.
	Metadata map[string]interface{}
}

// AbsoluteURL always prepends the input with the site URL.
func (s Site) AbsoluteURL(rawurl string) string {
	return fmt.Sprintf("%s/%s", s.URL, strings.TrimPrefix(rawurl, "/"))
}

// PackageURL returns the AbsoluteURL for the given package path, respecting the site's PackagePathPrefix.
func (s Site) PackageURL(pkgpath string) string {
	return s.AbsoluteURL(s.TrimPackagePath(pkgpath))
}

func (s Site) TrimPackagePath(pkgpath string) string {
	return strings.TrimPrefix(pkgpath, s.PackagePathPrefix)
}
