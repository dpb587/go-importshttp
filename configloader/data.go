package configloader

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go.dpb.io/importshttp"
	"go.dpb.io/importshttp/themeless"
	"go.dpb.io/importshttp/themepro"
)

type Data struct {
	RepositoryFactory importshttp.RepositoryFactory `yaml:"-"`
	Linkers           importshttp.LinkerList        `yaml:"-"` // TODO customizable

	Server   DataServer      `yaml:"server"`
	Site     DataSite        `yaml:"site"`
	Theme    DataTheme       `yaml:"theme"`
	Packages DataPackageList `yaml:"packages"`
}

type DataServer struct {
	Bind string `yaml:"bind"`
}

type DataSite struct {
	URL             string                 `yaml:"url"`
	Title           string                 `yaml:"title"`
	Generator       string                 `yaml:"generator"`
	ContentLanguage string                 `yaml:"content_language"`
	Links           DataLinkList           `yaml:"links"`
	Metadata        map[string]interface{} `yaml:"metadata"`
}

func (ds DataSite) AsSite() importshttp.Site {
	return importshttp.Site{
		URL:             ds.URL,
		Title:           ds.Title,
		Generator:       ds.Generator,
		ContentLanguage: ds.ContentLanguage,
		Links:           ds.Links.AsLinkList(),
		Metadata:        ds.Metadata,
	}
}

type DataLink struct {
	Ordering int    `yaml:"ordering"`
	Label    string `yaml:"label"`
	URL      string `yaml:"url"`
}

type DataLinkList []DataLink

func (dll DataLinkList) AsLinkList() importshttp.LinkList {
	var res importshttp.LinkList

	for dlIdx, dl := range dll {
		if dl.Ordering == 0 {
			// assume semi-intentional ordering
			dl.Ordering = dlIdx
		}

		res = append(
			res,
			importshttp.Link(dl),
		)
	}

	return res
}

type DataTheme struct {
	importshttp.Theme
}

func getThemeFromString(in string) (importshttp.Theme, error) {
	switch in {
	case "pro", "default", "true":
		return themepro.Theme, nil
	case "less", "false":
		return themeless.Theme, nil
	}

	return importshttp.Theme{}, fmt.Errorf("invalid theme name: %s", in)
}

func (s *DataTheme) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var attemptedString string

	if err := unmarshal(&attemptedString); err == nil {
		var err error

		s.Theme, err = getThemeFromString(attemptedString)

		return err
	}

	var attemptedBool bool
	if err := unmarshal(&attemptedBool); err == nil {
		var err error

		s.Theme, err = getThemeFromString(strconv.FormatBool(attemptedBool))

		return err
	}

	return errors.New("TODO")
}

type DataPackage struct {
	Import           string                 `yaml:"import"`
	ImportSubpackage string                 `yaml:"import_subpackage"`
	Repository       DataPackageRepository  `yaml:"repository"`
	Deprecated       bool                   `yaml:"deprecated"`
	Unlisted         bool                   `yaml:"unlisted"`
	Metadata         map[string]interface{} `yaml:"metadata"`
	Links            DataLinkList           `yaml:"links"`
}

type DataPackageList []DataPackage

func (dpl DataPackageList) AsPackageList(factory importshttp.RepositoryFactory) (importshttp.PackageList, error) {
	var res importshttp.PackageList

	for _, dp := range dpl {
		repo, err := factory.NewRepository(dp.Repository.RepositoryConfig)
		if err != nil {
			return nil, fmt.Errorf("getting repository of %s: %v", dp.Import, err)
		}

		res = append(
			res,
			importshttp.Package{
				Import:           dp.Import,
				ImportSubpackage: dp.ImportSubpackage,
				Repository:       repo,
				Deprecated:       dp.Deprecated,
				Unlisted:         dp.Unlisted,
				Metadata:         dp.Metadata,
				Links:            dp.Links.AsLinkList(),
			},
		)
	}

	return res, nil
}

type DataPackageRepository struct {
	importshttp.RepositoryConfig
}

func (s *DataPackageRepository) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var attemptedURL string

	err := unmarshal(&attemptedURL)
	if err == nil {
		if !strings.Contains(attemptedURL, "://") {
			attemptedURL = fmt.Sprintf("//%s", attemptedURL)
		}

		parsedURL, err := url.Parse(attemptedURL)
		if err != nil {
			return err
		}

		s.RepositoryConfig = importshttp.NewRepositoryConfigURL(
			importshttp.UnknownVCS,
			parsedURL,
		)

		return nil
	}

	var attemptedMap map[string]string

	err = unmarshal(&attemptedMap)
	if err == nil {
		service := attemptedMap["service"]
		delete(attemptedMap, "service")

		s.RepositoryConfig = importshttp.NewRepositoryConfigProperties(
			importshttp.VCS(service),
			attemptedMap,
		)

		return nil
	}

	return errors.New("expected repository to be one of: URL, map[string]string")
}
