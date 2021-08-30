package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NYTimes/gziphandler"
	"go.dpb.io/importshttp/configloader"
)

func main() {
	var exit bool

	flagSet := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	flagSet.BoolVar(&exit, "exit", false, "exit before server start (validates config)")

	config := configloader.New()
	err := configloader.ParseFlags(flagSet, os.Args[1:], &config)
	if err != nil {
		if err == flag.ErrHelp {
			return
		}

		panic(err)
	}

	resolved, err := config.Resolve()
	if err != nil {
		panic(err)
	}

	http.Handle("/", gziphandler.GzipHandler(resolved.Handler()))

	if exit {
		return
	}

	http.ListenAndServe(config.Server.Bind, nil)
}
