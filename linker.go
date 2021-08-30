package importshttp

// Linker
type Linker interface {
	Link(pkg Package) *Link
}

// LinkerList adds utility to a list of linkers.
type LinkerList []Linker

// GetPackageLinks will return a list of any generated links for the package.
func (ll LinkerList) GetPackageLinks(pkg Package) LinkList {
	var links LinkList

	for _, l := range ll {
		if link := l.Link(pkg); link != nil {
			links = append(links, *link)
		}
	}

	return links
}

func (ll LinkerList) MapPackageList(pkgs PackageList) PackageList {
	for mIdx := range pkgs {
		pkgs[mIdx].Links = append(pkgs[mIdx].Links, ll.GetPackageLinks(pkgs[mIdx])...)
		pkgs[mIdx].Links.SortByOrdering()
	}

	return pkgs
}
