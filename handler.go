package importshttp

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Handler struct {
	site              Site
	theme             Theme
	themeFilesHandler http.Handler
	packageList       PackageList
}

var _ http.Handler = Handler{}

func NewHandler(site Site, theme Theme, packageList PackageList) *Handler {
	// TODO optimizations on packages

	return &Handler{
		site:              site,
		theme:             theme,
		themeFilesHandler: http.StripPrefix(ThemeFilePrefix, http.FileServer(http.FS(theme.Files))),
		packageList:       packageList,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		h.render(w, r, h.theme.PackageListTemplate, map[string]interface{}{
			"PackageList": h.packageList.FilterByListed(),
		})

		return
	} else if strings.HasPrefix(r.URL.Path, ThemeFilePrefix) {
		if h.themeFilesHandler != nil {
			if len(h.theme.Version) > 0 {
				w.Header().Set("Cache-Control", "max-age=86400")
			}

			h.themeFilesHandler.ServeHTTP(w, r)

			return
		}
	} else if pkg, ok := h.packageList.Find(fmt.Sprintf("%s%s", h.site.PackagePathPrefix, r.URL.Path)); ok {
		h.render(w, r, h.theme.PackageTemplate, map[string]interface{}{
			"Package":        pkg,
			"SubpackageList": h.packageList.FilterByParent(pkg.Path()).FilterByListed(),
		})

		return
	}

	w.WriteHeader(http.StatusNotFound)
	h.render(w, r, h.theme.ErrorTemplate, map[string]interface{}{
		"StatusCode": http.StatusNotFound,
		"StatusText": http.StatusText(http.StatusNotFound),
	})
}

func (h Handler) render(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data map[string]interface{}) {
	data["Site"] = h.site
	data["Theme"] = h.theme

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
