package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NYTimes/gziphandler"
	"go.dpb.io/importshttp/config"
	"go.dpb.io/importshttp/config/configdefaults"
)

func main() {
	var exit bool

	flagSet := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	flagSet.BoolVar(&exit, "exit", false, "exit before server start (validates config)")

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

	http.Handle("/", gziphandler.GzipHandler(resolved.Handler()))

	if exit {
		return
	}

	http.ListenAndServe(rawconfig.Server.Bind, nil)
}
