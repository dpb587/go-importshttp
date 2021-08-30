// themepro is a minimal, professional theme.
package themepro

import (
	"embed"
	"html/template"
	"io/fs"

	"go.dpb.io/importshttp"
)

//go:generate sh -c "rm -fr files && cd build && NODE_ENV=production npx webpack && cd .. && cat files/* | sha1sum | cut -c-10 | tr -d '\\n' > files/version"

//go:embed files/*
var filesFS embed.FS

//go:embed files/version
var filesVersion string

//go:embed error.html
var errorTemplateData string

//go:embed package.html
var packageTemplateData string

//go:embed package_list.html
var packageListTemplateData string

var Theme importshttp.Theme

func init() {
	var themeFiles, _ = fs.Sub(filesFS, "files")

	Theme = importshttp.Theme{
		Version:             filesVersion,
		Files:               themeFiles,
		ErrorTemplate:       template.Must(template.New("html").Parse(errorTemplateData)),
		PackageTemplate:     template.Must(template.New("html").Parse(packageTemplateData)),
		PackageListTemplate: template.Must(template.New("html").Parse(packageListTemplateData)),
	}
}
