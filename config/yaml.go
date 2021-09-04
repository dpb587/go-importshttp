package config

import (
	"io"

	yaml "gopkg.in/yaml.v2"
)

/*
ParseYAML supports extracting configuration from YAML using the following schema. All fields are optional unless noted.

	# configuration when running a server
	server:

		# bind specification for the server (format: {address}:{port})
		bind:

	# settings which apply to all pages on the site
	site:

		# base url for building absolute urls
		url:

		# title to show on the index page
		title:

		# value to embed in html meta tags
		generator:

		# value to embed in html for the language of apges
		content_language:

		# list of informational links (default theme renders these in a footer)
		links:
		-
			# used to reprioritize links (type: integer)
			ordering:

			# text used when showing the link
			label:

			# web address used for the link
			url:

		# arbitrary data
		metadata: {}

	# use the default theme
	theme: true

	# or: use the minimal theme
	theme: false

	# or: use a custom theme from local files
	theme:

		# path to the template used for error pages (required)
		error_template:

		# path to the template used for a single package (required)
		package_template:

		# path to the template used for listing packages (required)
		package_list_template:

		# path to a directory with static assets
		files_dir:

		# theme version tag to include in file URL query strings for cache-busting
		version:

	# list of packages to publish
	packages:
	-
		# the Go package import path (required)
		import:

		# a subdirectory within the import for the full package path, if needed; typically empty
		import_subpackage:

		# URL to a common repository root for the package source
		repository:

		# indicate the package is deprecated (type: boolean)
		deprecated:

		# hide the package in listings (type: boolean)
		unlisted:

		# or: object with config describing a supported repository root for the package
		repository:

			# supported service name (values: bitbucket, bzr, fossil, git, github, gitlab, hg, mod, svn; required)
			vcs:

			# service-specific keys for configuring the repository (required)
			*:

		# list of potential action links for the package
		links:
		-
			# used to reprioritize links (type: integer)
			ordering:

			# text used when showing the link
			label:

			# web address used for the link
			url:

		# arbitrary data (default theme shows a plan text `description`, if present)
		metadata: {}

	# customization of built-in link generators
	package_links:

		# customize godoc links
		godoc:

			# base url for the godoc server
			base_url:

			# ordering for links (type: integer)
			ordering:

		# or: disable godoc links
		godoc: false

		# customize source links
		source:

			# ordering for links (type: integer)
			ordering:

		# or: disable source links
		source: false

*/
func ParseYAML(r io.Reader, raw *Raw) error {
	d := yaml.NewDecoder(r)
	d.SetStrict(true)

	return d.Decode(raw)
}
