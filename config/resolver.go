package config

import (
	"fmt"
	"net/http"
	"strings"

	"go.dpb.io/importshttp"
)

type Resolved struct {
	pkgs    importshttp.PackageList
	site    importshttp.Site
	theme   importshttp.Theme
	handler http.Handler
}

func (rd *Resolved) PackageList() importshttp.PackageList {
	return rd.pkgs
}

func (rd *Resolved) Site() importshttp.Site {
	return rd.site
}

func (rd *Resolved) Theme() importshttp.Theme {
	return rd.theme
}

func (rd *Resolved) Handler() http.Handler {
	return rd.handler
}

func (d Raw) Resolve() (*Resolved, error) {
	resolved := &Resolved{
		theme: d.Theme.Theme,
	}

	linkers := d.Linkers

	{ // pkgs
		pkgs, err := d.Packages.AsPackageList(d.RepositoryFactory)
		if err != nil {
			return nil, fmt.Errorf("loading packages: %s", err)
		}

		resolved.pkgs = linkers.MapPackageList(pkgs)
	}

	{ // site
		site := d.Site.AsSite()
		site.PackageLinkers = linkers
		site.Links.SortByOrdering()

		{
			firstSegments := map[string]struct{}{}

			for _, pkg := range resolved.pkgs {
				firstSegments[strings.SplitN(pkg.Path(), "/", 2)[0]] = struct{}{}
			}

			if len(firstSegments) > 1 {
				return nil, fmt.Errorf("found multiple root package path segments: %v", firstSegments)
			}

			var firstSegments0 string
			for k := range firstSegments {
				firstSegments0 = k
			}

			if len(site.PackagePathPrefix) == 0 {
				site.PackagePathPrefix = firstSegments0
			}

			if site.PackagePathPrefix != firstSegments0 {
				// TODO figure out more complex routing goals
				return nil, fmt.Errorf("undetermined behavior with configured root package path (%s) not matching package list (%s)", site.PackagePathPrefix, firstSegments0)
			}
		}

		resolved.site = site
	}

	resolved.handler = importshttp.NewHandler(resolved.site, resolved.theme, resolved.pkgs)

	return resolved, nil
}
