package configloader

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"go.dpb.io/importshttp/defaults"
)

func New() Data {
	config := Data{
		RepositoryFactory: defaults.RepositoryFactory,
		Linkers:           defaults.Linkers,

		Server: DataServer{
			Bind: "0.0.0.0:8080",
		},
		Site: DataSite{
			ContentLanguage: "en",
		},
		Theme: DataTheme{
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
