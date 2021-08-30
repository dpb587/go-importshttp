package importshttp

import (
	"fmt"
	"sort"
	"strings"
)

// PackageList adds utility to a list of packages.
type PackageList []Package

// Find will search the list for a package by the expected import. If an exact match is not found, the longest matching
// package will be returned (with the ImportSubpackage field including the unmatched path, assuming it is a subpackage).
// If still no match can be found, Package will be empty and bool will be false.
func (pl PackageList) Find(expected string) (Package, bool) {
	var longestok bool
	var longestpkg Package
	var longestpkgpath string

	for _, pkg := range pl {
		pkgpath := pkg.Path()

		if pkgpath == expected {
			return pkg, true
		} else if strings.HasPrefix(expected, fmt.Sprintf("%s/", pkgpath)) && len(pkgpath) > len(longestpkgpath) {
			longestok = true
			longestpkg = pkg
			longestpkgpath = pkg.Path()
			longestpkg.ImportSubpackage = strings.TrimPrefix(expected, fmt.Sprintf("%s/", longestpkgpath))
		}
	}

	return longestpkg, longestok
}

// FilterByParent returns all packages which appear as a direct subpackage of the input.
func (pl PackageList) FilterByParent(parentpkgpath string) PackageList {
	var pkgs PackageList

	prefix := fmt.Sprintf("%s/", parentpkgpath)

	for _, pkg := range pl {
		pkgpath := pkg.Path()

		if strings.HasPrefix(pkgpath, prefix) && !strings.Contains(strings.TrimPrefix(pkgpath, prefix), "/") {
			pkgs = append(pkgs, pkg)
		}
	}

	return pkgs
}

// FilterByListed will drop any packages whose Unlisted field is true.
func (pl PackageList) FilterByListed() PackageList {
	var pkgs PackageList

	for _, pkg := range pl {
		if !pkg.Unlisted {
			pkgs = append(pkgs, pkg)
		}
	}

	return pkgs
}

// SortByImport reorders the list by Path in ascending, lexicographic order.
func (pl PackageList) SortByImport() {
	sort.Slice(
		pl,
		func(i, j int) bool {
			// TODO try and sort major version suffixes?
			return strings.Compare(pl[i].Path(), pl[j].Path()) < 0
		},
	)
}
