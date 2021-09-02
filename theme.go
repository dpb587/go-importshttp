package importshttp

import (
	"fmt"
	"html/template"
	"io/fs"
)

// ThemeFilePrefix is the location relative to Handler where custom files are available from.
var ThemeFilePrefix = "/_theme/"

// Theme contains the settings
type Theme struct {
	// Version optional value which will be appended to theme files (i.e. for cache-busting).
	Version string

	// Files contains any static files used by the theme.
	Files fs.FS

	// ErrorTemplate is used for HTTP error pages.
	ErrorTemplate *template.Template

	// PackageTemplate is used when a package has been found.
	PackageTemplate *template.Template

	// PackageListTemplate is used when there is a list of packages available.
	PackageListTemplate *template.Template
}

// FileURL accepts a path relative to the theme files directory, prepends it with the prefix, and appends a version.
func (t Theme) FileURL(file string) string {
	if len(t.Version) == 0 {
		return file
	}

	return fmt.Sprintf("%s%s?v=%s", ThemeFilePrefix, file, t.Version)
}
