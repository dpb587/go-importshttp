// themeless is a very minimal, unstyled theme.
package themeless

import (
	"html/template"

	"go.dpb.io/importshttp"
)

var Theme = importshttp.Theme{
	ErrorTemplate: template.Must(template.New("html").Parse(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta name="robots" content="noindex">
		<title>Error - {{ .StatusText }}</title>
	</head>
	<body>
		<h1>Error - {{ .StatusText }}</h1>
		<p>HTTP {{ .StatusCode }}</p>
	</body>
</html>
`)),
	PackageTemplate: template.Must(template.New("html").Parse(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="go-import" content="{{ .Package.GoImport.MetaContent }}" />
		{{ with .Package.GoSource -}}
		<meta name="go-source" content="{{ .MetaContent }}" />
		{{- end }}
		{{- if .Package.Links }}
		{{- with index .Package.Links 0 }}
		<meta http-equiv="refresh" content="3;url={{ .URL }}" />
		{{- end }}
		{{- end }}
		<title>{{ .Package.Import }}</title>
		{{- if .Package.ImportSubpackage }}
		<meta name="robots" content="nofollow noindex">
		{{- end }}
	</head>
	<body>
		<h1>{{ .Package.Path }}</h1>
		{{- if .Package.Links }}
		{{- with index .Package.Links 0 -}}
		<p>Nothing to see here. Try <a href="{{ .URL }}">{{ .URL }}</a>.</p>
		{{- end }}
		{{- else }}
		<p>Nothing to see here. Try <code>go get {{ .Package.Import }}</code>.</p>
		{{- end }}
	</body>
</html>
`)),
	PackageListTemplate: template.Must(template.New("html").Parse(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>
			{{- with .Site.Title -}}
				{{ . }}
			{{- else -}}
				{{ .Site.PackagePathPrefix }}
			{{- end -}}
		</title>
		{{- with .Package.Metadata.description }}
		<meta name="description" content="{{ . }}">
		{{- end }}
		{{- with .Site.Generator }}
		<meta name="generator" content="{{ . }}">
		{{- end }}
	</head>
	<body>
		<h1>
			{{- with .Site.Title -}}
				{{ . }}
			{{- else -}}
				{{ .Site.PackagePathPrefix }}
			{{- end -}}
		</h1>
		{{ with .PackageList }}
		<ul>
			{{ range . }}
			<li>
				<a href="{{ $.Site.AbsoluteURL .ImportPath }}">{{ .Import }}</a>
				{{- if .Deprecated }}
				<img alt="Deprecated icon" src="{{ $.Site.AbsoluteURL ( $.Theme.FileURL "deprecated.svg" ) }}" title="Deprecated" />
				{{- end }}
			</li>
			{{ end }}
		</ul>
		{{ else }}
		<p>Nothing to see here.</p>
		{{ end }}
	</body>
</html>
`)),
}
