package configloader

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"go.dpb.io/importshttp"
)

// ParseFlags supports extracting basic configuration from command line options. For full configuration use a YAML file.
func ParseFlags(flagSet *flag.FlagSet, args []string, config *Data) error {
	var bind, configFile, theme string
	var pkgs flagPackageList

	flagSet.StringVar(&configFile, "config", "", "config file (yaml)")
	flagSet.StringVar(&bind, "bind", "", "bind ({address}:{port})")
	flagSet.StringVar(&theme, "theme", "", "theme (true/pro/default, false/less)")
	flagSet.Var(&pkgs, "pkg", "import and repository ({import}={repository-url}; {import}={vcs}={repo-root})") // TODO
	err := flagSet.Parse(args)
	if err != nil {
		return err
	} else if flagSet.NArg() > 0 {
		return fmt.Errorf("extra arguments found: %s", strings.Join(flagSet.Args(), " "))
	}

	if len(configFile) > 0 {
		fh, err := os.OpenFile(configFile, os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("opening config file: %s", err)
		}

		defer fh.Close()

		err = ParseYAML(fh, config)
		if err != nil {
			return fmt.Errorf("parsing config file: %s", err)
		}
	}

	if len(bind) > 0 {
		config.Server.Bind = bind
	}

	if len(theme) > 0 {
		config.Theme.Theme, err = getThemeFromString(theme)
		if err != nil {
			return err
		}
	}

	config.Packages = append(config.Packages, pkgs...)

	return nil
}

type flagPackageList DataPackageList

func (i *flagPackageList) String() string {
	return ""
}

func (i *flagPackageList) Set(value string) error {
	valueSplit := strings.SplitN(value, "=", 2)
	if len(valueSplit) != 2 {
		return fmt.Errorf("expected pkg in format of {import}={repository}")
	}

	if !strings.Contains(valueSplit[1], "://") {
		valueSplit[1] = fmt.Sprintf("//%s", valueSplit[1])
	}

	parsedURL, err := url.Parse(valueSplit[1])
	if err != nil {
		return err
	}

	dm := DataPackage{
		Import: valueSplit[0],
		Repository: DataPackageRepository{
			RepositoryConfig: importshttp.NewRepositoryConfigURL("", parsedURL),
		},
	}

	*i = append(*i, dm)

	return nil
}
