package configdefaults

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"go.dpb.io/importshttp/config"
	"go.dpb.io/importshttp/defaults"
)

func New() config.Raw {
	config := config.Raw{
		RepositoryFactory: defaults.RepositoryFactory,
		Linkers:           defaults.Linkers,

		Server: config.RawServer{
			Bind: "0.0.0.0:8080",
		},
		Site: config.RawSite{
			ContentLanguage: "en",
		},
		Theme: config.RawTheme{
			Theme: defaults.Theme,
		},
	}

	if v := os.Getenv("PORT"); len(v) > 0 {
		config.Server.Bind = fmt.Sprintf("0.0.0.0:%s", v)
	}

	if bi, ok := debug.ReadBuildInfo(); ok {
		config.Site.Generator = strings.TrimSpace(fmt.Sprintf("%s %s", bi.Path, bi.Main.Version))
	}

	return config
}
