package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/config"
	"go.dpb.io/importshttp/config/configdefaults"
)

func main() {
	var outdir string

	flagSet := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	flagSet.StringVar(&outdir, "out", "public", "directory to write site")

	rawconfig := configdefaults.New()
	err := config.ParseFlags(flagSet, os.Args[1:], &rawconfig)
	if err != nil {
		if err == flag.ErrHelp {
			return
		}

		panic(err)
	}

	resolved, err := rawconfig.Resolve()
	if err != nil {
		panic(err)
	}

	site := resolved.Site()
	handler := resolved.Handler()
	theme := resolved.Theme()

	mirror := func(remote, localfile string, expectedStatusCode int) error {
		req := httptest.NewRequest(http.MethodGet, remote, nil)

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		if recorder.Code != expectedStatusCode {
			return fmt.Errorf("unexpected status: %v", recorder.Code)
		}

		abslocalfile := filepath.Join(outdir, localfile)

		err := os.MkdirAll(filepath.Dir(abslocalfile), 0700)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(abslocalfile, recorder.Body.Bytes(), 0600)
		if err != nil {
			return err
		}

		return nil
	}

	err = mirror("/", "index.html", http.StatusOK)
	if err != nil {
		panic(fmt.Errorf("saving /: %s", err))
	}

	err = mirror(fmt.Sprintf("/magic-nonexistant-page-%d.html", time.Now().UnixNano()), "404.html", http.StatusNotFound)
	if err != nil {
		panic(fmt.Errorf("saving 404.html: %s", err))
	}

	for _, pkg := range resolved.PackageList() {
		err := mirror(site.PackageURL(pkg.Path()), filepath.Join(site.TrimPackagePath(pkg.Path()), "index.html"), http.StatusOK)
		if err != nil {
			panic(fmt.Errorf("saving %s: %s", site.PackageURL(pkg.Path()), err))
		}
	}

	err = fs.WalkDir(theme.Files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if d.IsDir() {
			return nil
		}

		relpath := filepath.Join(importshttp.ThemeFilePrefix, path)

		err = mirror(relpath, relpath, http.StatusOK)
		if err != nil {
			panic(fmt.Errorf("saving %s: %s", theme.FileURL(path), err))
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
