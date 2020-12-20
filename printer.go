package crawler

import (
	"fmt"
	"io"
	"sort"
)

// siteLinks is a container to store visited links and
// all links those were found in it
type siteLinks map[string]map[string]struct{}

// addLinkToSite stores links in the container
func (sl siteLinks) addLinkToSite(site, link string) {
	if _, ok := sl[site]; !ok {
		sl[site] = make(map[string]struct{})
	}
	sl[site][link] = struct{}{}
}

// Print returns sorted links/sublinks to Writer
func (sl siteLinks) Print(wr io.Writer) error {
	roots := make([]string, 0, len(sl))
	for r := range sl {
		roots = append(roots, r)
	}
	sort.Strings(roots)
	var err error
	for _, root := range roots {
		if _, err = fmt.Fprintf(wr, "%s\n", root); err != nil {
			return err
		}
		for _, link := range sortLinks(sl[root]) {
			if _, err = fmt.Fprintf(wr, "\t%s\n", link); err != nil {
				return err
			}
		}
	}
	return nil
}

// sortLinks sorts links by alphabet ascending
func sortLinks(links map[string]struct{}) []string {
	out := make([]string, 0, len(links))
	for link := range links {
		out = append(out, link)
	}
	sort.Strings(out)
	return out
}
