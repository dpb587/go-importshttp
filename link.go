package importshttp

import (
	"sort"
	"strings"
)

// Link is an address for finding more information.
type Link struct {
	// Ordering is used to sort multiple links - lower values imply more importance.
	Ordering int

	// Label is the text shown for the link.
	Label string

	// URL is the destination address.
	URL string
}

// LinkList adds utility to a list of links.
type LinkList []Link

// SortByOrdering reorders the list by Ordering in ascending, numeric order.
func (ll LinkList) SortByOrdering() {
	sort.Slice(
		ll,
		func(i, j int) bool {
			ip, jp := ll[i].Ordering, ll[j].Ordering
			if ip == jp {
				return strings.Compare(ll[i].Label, ll[j].Label) < 0
			}

			return ip < jp
		},
	)
}
