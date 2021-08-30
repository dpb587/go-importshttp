package importshttp

import "strings"

// SourceRepositoryLinker generates a link based on a Package's SourceRepository, as available.
type SourceRepositoryLinker struct {
	Ordering int
	Label    string
}

func (l SourceRepositoryLinker) Link(pkg Package) *Link {
	rs, ok := pkg.Repository.(SourceRepository)
	if !ok {
		return nil
	}

	if sourceDirTemplateURL := rs.SourceDirTemplateURL(); len(sourceDirTemplateURL) > 0 {
		return &Link{
			Ordering: l.Ordering,
			Label:    l.Label,
			URL:      strings.Replace(sourceDirTemplateURL, "{/dir}", pkg.ImportSubpackage, -1),
		}
	} else if sourceURL := rs.SourceURL(); len(sourceURL) > 0 {
		return &Link{
			Ordering: l.Ordering,
			Label:    l.Label,
			URL:      sourceURL,
		}
	}

	return nil
}
